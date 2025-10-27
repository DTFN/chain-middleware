package com.lingshu.server.utils;

import org.bouncycastle.crypto.CryptoException;
import org.bouncycastle.crypto.Signer;
import org.bouncycastle.crypto.params.Ed25519PrivateKeyParameters;
import org.bouncycastle.crypto.params.Ed25519PublicKeyParameters;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.util.encoders.Hex;

import java.nio.charset.StandardCharsets;
import java.security.Security;

public class Ed25519Signer {

    static {
        // 添加BouncyCastle提供器
        Security.addProvider(new BouncyCastleProvider());
    }

    /**
     * 使用ED25519算法对VC进行签名
     *
     * @param privateKeyHex 十六进制格式的私钥
     * @param vcContent     要签名的VC内容（JSON字符串）
     * @return 十六进制格式的签名结果
     * @throws CryptoException 加密异常
     */
    public static String signVC(String privateKeyHex, String vcContent) throws CryptoException {
        // 将十六进制私钥转换为字节数组
        byte[] privateKeyBytes = Hex.decode(privateKeyHex);

        // 创建ED25519私钥参数
        Ed25519PrivateKeyParameters privateKey = new Ed25519PrivateKeyParameters(privateKeyBytes, 0);

        // 初始化签名器
        Signer signer = new org.bouncycastle.crypto.signers.Ed25519Signer();
        signer.init(true, privateKey);

        // 4. 对VC内容进行签名
        byte[] vcBytes = vcContent.getBytes(StandardCharsets.UTF_8);
        signer.update(vcBytes, 0, vcBytes.length);
        byte[] signature = signer.generateSignature();

        // 5. 将签名转换为十六进制字符串返回
        return Hex.toHexString(signature);
    }

    /**
     * 验证VC签名
     *
     * @param publicKeyHex 十六进制格式的公钥
     * @param vcContent    原始VC内容
     * @param signatureHex 十六进制格式的签名
     * @return 验证结果（true=有效，false=无效）
     */
    public static boolean verifyVC(String publicKeyHex, String vcContent, String signatureHex) {
        try {
            // 1. 转换公钥和签名为字节数组
            byte[] publicKeyBytes = Hex.decode(publicKeyHex);
            byte[] signatureBytes = Hex.decode(signatureHex);
            byte[] vcBytes = vcContent.getBytes(StandardCharsets.UTF_8);

            // 2. 创建公钥参数
            Ed25519PublicKeyParameters publicKey = new Ed25519PublicKeyParameters(publicKeyBytes, 0);

            // 3. 初始化验证器
            Signer verifier = new org.bouncycastle.crypto.signers.Ed25519Signer();
            verifier.init(false, publicKey);

            // 4. 验证签名
            verifier.update(vcBytes, 0, vcBytes.length);
            return verifier.verifySignature(signatureBytes);
        } catch (Exception e) {
            e.printStackTrace();
            return false;
        }
    }

    /**
     * 从私钥推导公钥（ED25519公钥可由私钥直接推导）
     *
     * @param privateKeyHex 十六进制私钥
     * @return 十六进制公钥
     */
    public static String derivePublicKey(String privateKeyHex) {
        byte[] privateKeyBytes = Hex.decode(privateKeyHex);
        Ed25519PrivateKeyParameters privateKey = new Ed25519PrivateKeyParameters(privateKeyBytes, 0);
        Ed25519PublicKeyParameters publicKey = privateKey.generatePublicKey();
        return Hex.toHexString(publicKey.getEncoded());
    }
}
