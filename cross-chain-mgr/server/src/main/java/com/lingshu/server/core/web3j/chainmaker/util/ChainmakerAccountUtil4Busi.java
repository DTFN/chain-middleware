package com.lingshu.server.core.web3j.chainmaker.util;

import cn.hutool.core.util.StrUtil;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.account.entity.ChainAccountDO;
import com.lingshu.server.core.web3j.chainmaker.ChainmakerProperties;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import lombok.extern.slf4j.Slf4j;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.crypto.digests.SM3Digest;
import org.bouncycastle.jcajce.provider.asymmetric.ec.BCECPublicKey;
import org.bouncycastle.jcajce.provider.digest.Keccak;
import org.bouncycastle.jce.ECNamedCurveTable;
import org.bouncycastle.jce.interfaces.ECPublicKey;
import org.bouncycastle.jce.spec.ECNamedCurveParameterSpec;
import org.bouncycastle.jce.spec.ECParameterSpec;
import org.bouncycastle.jce.spec.ECPublicKeySpec;
import org.bouncycastle.math.ec.ECPoint;
import org.bouncycastle.openssl.PEMKeyPair;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;
import org.bouncycastle.openssl.jcajce.JcaPEMWriter;
import org.bouncycastle.util.encoders.Hex;
import org.bouncycastle.util.io.pem.PemObject;
import org.bouncycastle.util.io.pem.PemWriter;
import org.chainmaker.pb.config.ChainConfigOuterClass;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.User;
import org.chainmaker.sdk.config.AuthType;
import org.chainmaker.sdk.crypto.ChainMakerCryptoSuiteException;
import org.chainmaker.sdk.utils.CryptoUtils;
import org.springframework.stereotype.Component;

import javax.annotation.Resource;
import java.io.StringReader;
import java.io.StringWriter;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.security.*;
import java.security.cert.Certificate;
import java.security.interfaces.ECPrivateKey;
import java.security.spec.InvalidKeySpecException;

/**
 * @author gongrui.wang
 * @since 2025/2/13
 */
@Slf4j
@Component
public class ChainmakerAccountUtil4Busi {
    @Resource(name = "chainmakerChainClient4Busi")
    private ChainmakerChainClient chainmakerChainClient;
    @Resource(name = "chainmakerPropertiesBusi")
    private ChainmakerProperties chainmakerProperties;

    /**
     * 获取用户账号
     * @param accountDO
     * @return
     */
    public User toUser(ChainAccountDO accountDO) {
        ChainClient chainClient = chainmakerChainClient.getChainClient();
        if (accountDO == null) {
            return null;
        }
        if (StrUtil.isEmpty(accountDO.getPrivateKey())) {
            throw new ApiException("用户私钥为空");
        }
        if (StrUtil.isEmpty(accountDO.getPublicKey())) {
            throw new ApiException("用户公钥为空");
        }
        try {

            ECPrivateKey privateKey = convertPemToPrivateKey(accountDO.getPrivateKey());
            ECPublicKey publicKey = convertPemToPublicKey(accountDO.getPublicKey());

            User user = new User(chainClient.getClientUser().getOrgId());
            user.setPrivateKey(privateKey);
            user.setPriBytes(accountDO.getPrivateKey().getBytes());
            user.setPublicKey(publicKey);
            user.setPukBytes(accountDO.getPublicKey().getBytes());
            user.setCryptoSuite(chainClient.getClientUser().getCryptoSuite());
            user.setAuthType(chainClient.getClientUser().getAuthType());
            return user;
        } catch (Exception e) {
            log.error(e.getMessage(), e);
            throw new ApiException("获取用户信息失败");
        }
    }


    /**
     * 转换证书为PEM格式
     * @param certificate
     * @return
     */
    private String convertCertificateToPem(Certificate certificate) {
        StringWriter sw = new StringWriter();
        PemWriter pw = new PemWriter(sw);
        try {
            pw.writeObject(new PemObject("CERTIFICATE", certificate.getEncoded()));
            pw.flush();
        } catch (Exception e) {
            throw new RuntimeException("Failed to convert certificate to PEM format", e);
        } finally {
            try {
                pw.close();
            } catch (Exception e) {
            }
        }
        return sw.toString();
    }

    /**
     * 转换私钥为PEM格式
     * @param privateKey
     * @return
     */
    private String convertKeyToPem(ECPrivateKey privateKey) {
        StringWriter sw = new StringWriter();
        JcaPEMWriter pw = new JcaPEMWriter(sw);
        try {
            pw.writeObject(privateKey);
            pw.flush();
        } catch (Exception e) {
            throw new RuntimeException("Failed to convert key to PEM format", e);
        } finally {
            try {
                pw.close();
            } catch (Exception e) {
                // Ignore
            }
        }
        return sw.toString();
    }

    /**
     * pem格式私钥字符串转为ECPrivate
     * @param pem
     * @return
     */
    public static ECPrivateKey convertPemToPrivateKey(String pem) {
        try (StringReader sr = new StringReader(pem)) {
            PEMParser pemParser = new PEMParser(sr);
            PEMKeyPair object = (PEMKeyPair) pemParser.readObject();
            JcaPEMKeyConverter converter = new JcaPEMKeyConverter().setProvider("BC");
            PrivateKey privateKey = converter.getPrivateKey(object.getPrivateKeyInfo());
            return (ECPrivateKey) privateKey;
        } catch (Exception e) {
            throw new RuntimeException("Failed to convert PEM to ECPrivateKey", e);
        }
    }

    private String convertPublicKeyToPem(ECPublicKey publicKey) {
        StringWriter sw = new StringWriter();
        JcaPEMWriter pw = new JcaPEMWriter(sw);
        try {
            pw.writeObject(publicKey);
            pw.flush();
        } catch (Exception e) {
            throw new RuntimeException("Failed to convert key to PEM format", e);
        }
        try {
            pw.close();
        } catch (Exception e) {
            // Ignore
        }
        return sw.toString();
    }

    public static ECPublicKey convertPemToPublicKey(String pem) {
        try (StringReader sr = new StringReader(pem)) {
            PEMParser pemParser = new PEMParser(sr);
            SubjectPublicKeyInfo subjectPublicKeyInfo = (SubjectPublicKeyInfo) pemParser.readObject();
            JcaPEMKeyConverter converter = new JcaPEMKeyConverter().setProvider("BC");
            PublicKey publicKey = converter.getPublicKey(subjectPublicKeyInfo);
            return (ECPublicKey) publicKey;
        } catch (Exception e) {
            throw new RuntimeException("Failed to convert PEM to ECPublicKey", e);
        }
    }

    /**
     * 生成用户
     * @return
     */
    public User generateUser() {
        return generateUser(chainmakerProperties.getIsGm());
    }

    /**
     * 生成用户账号
     * @param isGm
     * @return
     */
    public User generateUser(boolean isGm) {
        // 新建一个账户的公私钥
        try {
            KeyPairGenerator keyPairGenerator = null;
            if (isGm) {
                ECParameterSpec spec = ECNamedCurveTable.getParameterSpec("sm2p256v1");
                keyPairGenerator = KeyPairGenerator.getInstance("EC", "BC");
                keyPairGenerator.initialize(spec, new SecureRandom());
            } else {
                ECParameterSpec spec = ECNamedCurveTable.getParameterSpec("secp256r1");
                keyPairGenerator = KeyPairGenerator.getInstance("ECDSA", "BC");
                keyPairGenerator.initialize(spec, new SecureRandom());
            }
            KeyPair keyPair = keyPairGenerator.generateKeyPair();
            ECPrivateKey privateKey = (ECPrivateKey) keyPair.getPrivate();
            ECPublicKey publicKey = (ECPublicKey) keyPair.getPublic();

            String priPem = convertKeyToPem(privateKey);
            String pubPem = convertPublicKeyToPem(publicKey);
            User user = new User("", priPem.getBytes(), new byte[]{}, pubPem.getBytes(), AuthType.Public.getMsg());
            user.setPriBytes(priPem.getBytes());
            user.setPrivateKey(privateKey);
            user.setPublicKey(publicKey);
            return user;
        } catch (Exception e) {
            throw new RuntimeException("Failed to generate user", e);
        }
    }

    /**
     * 计算账户地址
     * @return
     */
    public String accountAddress(User user) {
        // 是否国密
        if (chainmakerProperties.getIsGm()) {
            return accountAddress(user.getPublicKey());
        }

        if (user.getPublicKey() != null) {
            return "0x" + CryptoUtils.pkToAddrStr(user.getPublicKey(), ChainConfigOuterClass.AddrType.ETHEREUM, "sha256");
        }
        return "0x" + CryptoUtils.getEVMAddressFromPrivateKeyBytes(user.getPriBytes(), "sha256");
    }

    // 根据私钥pem生成地址
    public static String accountAddress(boolean isGm, String privateKeyPem) {
        try {
            // 从pem读取私钥并生成公钥
            PublicKey publicKey = getPublicKeyFromPrivateKey(isGm, CryptoUtils.getPrivateKeyFromBytes(privateKeyPem.getBytes(StandardCharsets.UTF_8)));

            // 长安链国密和非国密都使用ETHEREUM
            ChainConfigOuterClass.AddrType addrType = ChainConfigOuterClass.AddrType.ETHEREUM;

            BCECPublicKey pubKey = (BCECPublicKey)publicKey;
            ECParameterSpec spec = pubKey.getParameters();
            byte[] result = publicKeyToByte(pubKey, ((ECNamedCurveParameterSpec)spec).getName());

            byte[] dest = new byte[result.length - 1];
            System.arraycopy(result, 1, dest, 0, result.length - 1);

            byte[] bytesAddr = keccak256(dest);
            String addrStr = Hex.toHexString(bytesAddr).substring(24);

            if (addrStr.startsWith("0x") || addrStr.startsWith("0X")) {
                return addrStr;
            } else {
                return "0x" + addrStr;
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    public static PublicKey getPublicKeyFromPrivateKey(boolean isGm, PrivateKey privateKey)throws ChainMakerCryptoSuiteException, NoSuchAlgorithmException, InvalidKeySpecException, NoSuchProviderException {
        org.bouncycastle.jce.interfaces.ECPrivateKey ecPrivateKey = (org.bouncycastle.jce.interfaces.ECPrivateKey)privateKey;
        ECNamedCurveParameterSpec spec = isGm? ECNamedCurveTable.getParameterSpec("sm2p256v1"): ECNamedCurveTable.getParameterSpec("secp256r1");
        ECPoint Q = spec.getG().multiply(ecPrivateKey.getD());
        ECPublicKeySpec pubSpec = new ECPublicKeySpec(Q, spec);
        KeyFactory keyFactory = KeyFactory.getInstance("EC", "BC");
        return keyFactory.generatePublic(pubSpec);
    }

    // 根据私钥pem生成地址
    public static String accountAddress(PublicKey publicKey) {
        try {
            // 长安链国密和非国密都使用ETHEREUM
            ChainConfigOuterClass.AddrType addrType = ChainConfigOuterClass.AddrType.ETHEREUM;

            BCECPublicKey pubKey = (BCECPublicKey)publicKey;
            ECParameterSpec spec = pubKey.getParameters();
            byte[] result = publicKeyToByte(pubKey, ((ECNamedCurveParameterSpec)spec).getName());

            byte[] dest = new byte[result.length - 1];
            System.arraycopy(result, 1, dest, 0, result.length - 1);

            byte[] bytesAddr = keccak256(dest);
            String addrStr = Hex.toHexString(bytesAddr).substring(24);

            if (addrStr.startsWith("0x") || addrStr.startsWith("0X")) {
                return addrStr;
            } else {
                return "0x" + addrStr;
            }
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    private static byte[] keccak256(byte[] data) {
        Keccak.Digest256 keccakDigest = new Keccak.Digest256();
        return keccakDigest.digest(data);
    }

    private static String zxAddress(byte[] data) {
        SM3Digest digest = new SM3Digest();
        digest.update(data, 0, data.length);
        byte[] result = new byte[digest.getDigestSize()];
        digest.doFinal(result, 0);
        if (result.length < 20) {
            return "";
        } else {
            byte[] dest = new byte[20];
            System.arraycopy(result, 0, dest, 0, 20);
            return "ZX" + Hex.toHexString(dest);
        }
    }

    private static byte[] publicKeyToByte(BCECPublicKey publicKey, String name) {
        return name.contains("secp") ? ecPublickeyToByte(publicKey) : sm2PublickeyToByte(publicKey);
    }

    // 未用到
    private static byte[] ecPublickeyToByte(BCECPublicKey publicKey) {
        int length = (publicKey.getParams().getCurve().getField().getFieldSize() + 7) / 8;
        byte[] x = publicKey.getQ().getXCoord().getEncoded();
        byte[] y = publicKey.getQ().getYCoord().getEncoded();
        byte[] result = new byte[1 + 2 * length];
        result[0] = 4;
        System.arraycopy(x, 0, result, 1, length);
        System.arraycopy(y, 0, result, length + 1, length);
        return result;
    }

    private static byte[] sm2PublickeyToByte(BCECPublicKey publicKey) {
        int length = (publicKey.getParams().getCurve().getField().getFieldSize() + 7) / 8;
        BigInteger x = publicKey.getQ().getXCoord().toBigInteger();
        BigInteger y = publicKey.getQ().getYCoord().toBigInteger();
        byte[] a = x.abs().toByteArray();
        byte[] b = y.abs().toByteArray();

        // BigIntger有时会在数组前加入0x00字节的情况
        boolean aStartZero = a[0] == 0;
        boolean bStartZero = b[0] == 0;

        byte[] result = new byte[1 + 2 * length];
        result[0] = 4;
        System.arraycopy(a, aStartZero?1:0, result, 1, length);
        System.arraycopy(b, bStartZero?1:0, result, length + 1, length);
        return result;
    }
}
