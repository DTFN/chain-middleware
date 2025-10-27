package com.lingshu.serverbcos.utils;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.web3j.crypto.*;
import org.web3j.utils.Numeric;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.Security;
import java.util.Arrays;

public class EthVCSigner {

    // secp256k1曲线阶（固定值）
    private static final BigInteger N = new BigInteger("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16);
    private static final BigInteger HALF_N = N.shiftRight(1); // n/2，用于S值规范化

    static {
        Security.addProvider(new BouncyCastleProvider());
    }

    private static final ObjectMapper objectMapper = new ObjectMapper()
            .disable(SerializationFeature.INDENT_OUTPUT)
            .disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);

    // 计算contentHash（原始VC的Keccak-256）
    public static String calculateContentHash(String vcJson) throws JsonProcessingException{
        // 1. 标准化VC的JSON格式
        JsonNode vcNode = objectMapper.readTree(vcJson);
        String normalizedVc = objectMapper.writeValueAsString(vcNode);

        // 2. 以太坊签名前缀处理（EIP-191规范：\u0019Ethereum Signed Message:\n+消息长度）
        byte[] messageBytes = normalizedVc.getBytes(StandardCharsets.UTF_8);
        byte[] prefix = ("\u0019Ethereum Signed Message:\n" + messageBytes.length).getBytes(StandardCharsets.UTF_8);
        byte[] prefixedMessage = new byte[prefix.length + messageBytes.length];
        System.arraycopy(prefix, 0, prefixedMessage, 0, prefix.length);
        System.arraycopy(messageBytes, 0, prefixedMessage, prefix.length, messageBytes.length);
        byte[] messageHash = Hash.sha3(prefixedMessage);
        return Numeric.toHexStringNoPrefix(messageHash);
    }

    /**
     * 生成带前导零的64位十六进制地址（与ether-signer格式一致）
     * @param privateKeyHex 十六进制私钥
     * @return 64位十六进制地址（全小写，带前导零）
     */
    public static String getAddressFromPrivateKey(String privateKeyHex) throws Exception {
        String fullPrivateKey = privateKeyHex.length() < 64 ?
                String.format("%064x", new BigInteger(privateKeyHex, 16)) :
                privateKeyHex;
        // 2. 从私钥推导标准以太坊地址（40位，不带0x前缀）
        ECKeyPair ecKeyPair = ECKeyPair.create(new BigInteger(fullPrivateKey, 16));
        String standardAddress = Keys.getAddress(ecKeyPair);

        // 3. 补全为64位（前导补零）
        return String.format("%064x", new BigInteger(standardAddress, 16)).toLowerCase();
    }

    /**
     * 对VC进行以太坊签名（增加低S值规范化处理）
     */
    public static Sign.SignatureData signVC(String privateKeyHex, String vcJson) throws JsonProcessingException {
        try {
            // 1. 标准化VC的JSON格式
            JsonNode vcNode = objectMapper.readTree(vcJson);
            String normalizedVc = objectMapper.writeValueAsString(vcNode);

            // 2. 以太坊签名前缀处理（EIP-191规范：\u0019Ethereum Signed Message:\n+消息长度）
            byte[] messageBytes = normalizedVc.getBytes(StandardCharsets.UTF_8);
            byte[] prefix = ("\u0019Ethereum Signed Message:\n" + messageBytes.length).getBytes(StandardCharsets.UTF_8);
            byte[] prefixedMessage = new byte[prefix.length + messageBytes.length];
            System.arraycopy(prefix, 0, prefixedMessage, 0, prefix.length);
            System.arraycopy(messageBytes, 0, prefixedMessage, prefix.length, messageBytes.length);
            byte[] messageHash = Hash.sha3(prefixedMessage);

            // 3. 加载私钥（处理可能的前导0x前缀）
            String cleanPrivateKey = Numeric.cleanHexPrefix(privateKeyHex);
            ECKeyPair keyPair = ECKeyPair.create(Numeric.hexStringToByteArray(cleanPrivateKey));

            // 4. 生成ECDSA签名并处理低S值（符合EIP-2规范）
            ECDSASignature signature = keyPair.sign(messageHash);
            // 关键修正：确保s值小于n/2（符合以太坊规范）
            if (signature.s.compareTo(HALF_N) > 0) {
                signature = new ECDSASignature(signature.r, N.subtract(signature.s)); // 翻转S值
            }

            // 5. 标准化r和s：确保为32字节（移除前导零/补前导零）
            byte[] rBytes = standardizeSignatureBytes(signature.r.toByteArray());
            byte[] sBytes = standardizeSignatureBytes(signature.s.toByteArray());

            int recId = recoverRecId(keyPair, messageHash, signature);
            byte v = (byte) (recId + 27);

            return new Sign.SignatureData(v, rBytes, sBytes);
        } catch (Exception e) {
            throw new RuntimeException("签名生成失败: " + e.getMessage(), e);
        }
    }

    /**
     * 关键工具方法：标准化签名的r/s字节数组，确保为32字节（64位十六进制）
     * - 若输入超过32字节：移除前导零（如33字节的00fe... → 32字节的fe...）
     * - 若输入不足32字节：前导补零（如31字节的123... → 32字节的00123...）
     */
    private static byte[] standardizeSignatureBytes(byte[] inputBytes) {
        byte[] output = new byte[32]; // 固定32字节输出
        if (inputBytes.length > 32) {
            // 情况1：输入超过32字节（含前导零），截取后32字节（移除前导零）
            System.arraycopy(inputBytes, inputBytes.length - 32, output, 0, 32);
        } else {
            // 情况2：输入不足32字节，前导补零，数据从末尾开始填充
            System.arraycopy(inputBytes, 0, output, 32 - inputBytes.length, inputBytes.length);
        }
        return output;
    }

    /**
     * 恢复签名的恢复ID（recId），用于后续地址验证
     */
    private static int recoverRecId(ECKeyPair keyPair, byte[] messageHash, ECDSASignature signature) {
        for (int recId = 0; recId < 4; recId++) { // 恢复ID范围是0-3
            BigInteger recoveredPubKey = Sign.recoverFromSignature(recId, signature, messageHash);
            if (recoveredPubKey != null && recoveredPubKey.equals(keyPair.getPublicKey())) {
                return recId;
            }
        }
        throw new RuntimeException("无法生成有效的恢复ID（recId），签名可能无效");
    }

    /**
     * 验证以太坊签名
     */
    public static boolean verifyVC(String publicKeyHex, String vcJson, Sign.SignatureData signatureData) throws JsonProcessingException {
        try {
            // 1. 标准化VC的JSON格式
            JsonNode vcNode = objectMapper.readTree(vcJson);
            String normalizedVc = objectMapper.writeValueAsString(vcNode);
            byte[] messageBytes = normalizedVc.getBytes(StandardCharsets.UTF_8);

            // 2. 处理签名前缀
            byte[] prefix = ("\u0019Ethereum Signed Message:\n" + messageBytes.length).getBytes(StandardCharsets.UTF_8);
            byte[] prefixedMessage = new byte[prefix.length + messageBytes.length];
            System.arraycopy(prefix, 0, prefixedMessage, 0, prefix.length);
            System.arraycopy(messageBytes, 0, prefixedMessage, prefix.length, messageBytes.length);
            byte[] messageHash = Hash.sha3(prefixedMessage);

            // 3. 解析签名
            int recId = signatureData.getV()[0] - 27; // 转换为恢复ID
            ECDSASignature sig = new ECDSASignature(
                    new BigInteger(1, signatureData.getR()),
                    new BigInteger(1, signatureData.getS())
            );

            // 4. 恢复公钥
            BigInteger recoveredPublicKey = Sign.recoverFromSignature(recId, sig, messageHash);
            if (recoveredPublicKey == null) {
                return false;
            }

            // 5. 验证公钥
            byte[] expectedPublicKey = Numeric.hexStringToByteArray(publicKeyHex);
            byte[] recoveredPublicKeyBytes = Numeric.toBytesPadded(recoveredPublicKey, 64);
            return Arrays.equals(recoveredPublicKeyBytes, expectedPublicKey);
        } catch (Exception e) {
            e.printStackTrace();
            return false;
        }
    }

    /**
     * 从私钥推导公钥
     */
    public static String derivePublicKey(String privateKeyHex) {
        ECKeyPair keyPair = ECKeyPair.create(Numeric.hexStringToByteArray(privateKeyHex));
        return Numeric.toHexStringNoPrefix(Numeric.toBytesPadded(keyPair.getPublicKey(), 64));
    }

    public static void main(String[] args) throws Exception {
        // 1. 待签名的VC（压缩格式）
        String vcJson0 = "{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01-daiding\",\"contract_func\":\"erc20MintVcs\",\"contract_type\":\"1\",\"func_params\":{\"amount\":100,\"to\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}";
        // 示例VC内容
        String vc = "{\n" +
                "  \"@context\": [\n" +
                "    \"https://www.w3.org/2018/credentials/v1\"\n" +
                "  ],\n" +
                "  \"credentialSubject\": {\n" +
                "    \"chain_rid\": \"chain1\",\n" +
                "    \"contract_address\": \"contract01-daiding\",\n" +
                "    \"contract_func\": \"erc20MintVcs\",\n" +
                "    \"contract_type\": \"1\",\n" +
                "    \"func_params\": {\n" +
                "      \"amount\": 100,\n" +
                "      \"to\": \"did:1\"\n" +
                "    },\n" +
                "    \"gateway_id\": \"0\",\n" +
                "    \"param_name\": \"content\",\n" +
                "    \"reource_name\": \"chain1:contract01\"\n" +
                "  },\n" +
                "  \"id\": \"https://uni.example.com/credentials/12345\",\n" +
                "  \"issuanceDate\": \"2025-05-20\",\n" +
                "  \"issuer\": \"https://uni.example.com/organization/1\",\n" +
                "  \"type\": [\n" +
                "    \"VerifiableCredential\",\n" +
                "    \"UniversityDegreeCredential\"\n" +
                "  ]\n" +
                "}";

        String vcJson = JsonOrderProcessor.convert(vc);

        // 2. 私钥
        String privateKeyHex = "111111111";

        // 计算address
        String address = getAddressFromPrivateKey(privateKeyHex);
        System.out.println("address: " + address);
        // 计算contentHash
        String contentHash = calculateContentHash(vcJson);
        System.out.println("contentHash: " + contentHash);
        // 签名生成r/s/v
        // 3. 生成签名
        Sign.SignatureData signatureData = signVC(privateKeyHex, vcJson);
        String r = Numeric.toHexStringNoPrefix(signatureData.getR());
        String s = Numeric.toHexStringNoPrefix(signatureData.getS());
        int v = signatureData.getV()[0] & 0xFF; // 转换为无符号整数
        System.out.println("生成的签名:");
        System.out.println("r: " + r);
        System.out.println("s: " + s);
        System.out.println("v: " + v);

        // 4. 验证签名
        String publicKey = derivePublicKey(privateKeyHex);
        boolean isValid = verifyVC(publicKey, vcJson, signatureData);
        System.out.println("签名验证结果: " + (isValid ? "有效" : "无效"));
    }
}
    
