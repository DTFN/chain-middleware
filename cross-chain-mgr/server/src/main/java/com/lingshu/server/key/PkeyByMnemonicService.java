package com.lingshu.server.key;

import cn.hutool.crypto.digest.DigestUtil;
import com.lingshu.server.key.util.CertUtils;
import com.lingshu.server.key.util.MasterKeyGenerator;
import com.lingshu.server.key.util.SecureRandomUtils;
import com.lingshu.server.key.util.SeedGenerator;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.bouncycastle.jcajce.provider.asymmetric.ec.BCECPrivateKey;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.springframework.stereotype.Component;
import org.web3j.crypto.ECKeyPair;
import org.web3j.crypto.Keys;
import org.web3j.crypto.MnemonicUtils;
import org.web3j.utils.Numeric;

import java.io.ByteArrayInputStream;
import java.io.ByteArrayOutputStream;
import java.security.*;
import java.security.cert.Certificate;

@Slf4j
@Component
public class PkeyByMnemonicService {
    
    private static final SecureRandom secureRandom = SecureRandomUtils.secureRandom();

    private static final SeedGenerator seedGenerator = new SeedGenerator();
    private static final MasterKeyGenerator keyGenerator = new MasterKeyGenerator();

    private static Certificate dummyCert;

    static {
        try {
            dummyCert = CertUtils.generateDummyCertificate();
        } catch (Exception e) {
            log.error("Error generating dummy cert, so p12 cannot be used",e);
        }
    }


    /**
     * create mnemonic by entropyStr, entropyStr can be null. If entropy is null, a 128 bit random entropy will be used.
     * 
     * @param entropyStr: random  entropy
     * @return mnemonic. mnemonic length is determined by entrophy length
     */
    public String createMnemonic(String entropyStr) {
        byte[] initialEntropy;
        if (StringUtils.isEmpty(entropyStr)) {
            initialEntropy = new byte[16];
            secureRandom.nextBytes(initialEntropy);
        }else{
            initialEntropy = Numeric.hexStringToByteArray(entropyStr);
        }
        String mnemonic = MnemonicUtils.generateMnemonic(initialEntropy);
        return mnemonic;
    }

    /**
     * create mnemonic. A 128 bit random entropy will be used.
     *
     * @return mnemonic. mnemonic length is 12
     */
    public String createMnemonic() {
        byte[] initialEntropy = new byte[16];
        secureRandom.nextBytes(initialEntropy);
        String mnemonic = MnemonicUtils.generateMnemonic(initialEntropy);
        return mnemonic;
    }
    /**    
     * generate PkeyInfo by mnemonic
     *
     * @param mnemonic
     * @param passphrase: password for create seed, it can be null.  
     * @return PkeyInfo
     */
    public ECKeyPair generatePrivateKeyByMnemonic(String mnemonic, String passphrase)
            throws InvalidAlgorithmParameterException, NoSuchAlgorithmException, NoSuchProviderException {
    	byte[] seed = seedGenerator.generateSeed(mnemonic, DigestUtil.sha256Hex(passphrase));
        return Keys.createEcKeyPair(new SecureRandom(seed));
    }

    public String encryptPrivateKeyByP12(String password, byte[] privateKey) throws Exception {
        if(privateKey == null || privateKey.length != 32) {
            throw new IllegalArgumentException("privateKey");
        }
        // 将字节数组转换为 BCECPrivateKey 对象
        BCECPrivateKey pk = new BCECPrivateKey("EC",
                new org.bouncycastle.jce.spec.ECPrivateKeySpec(
                        new java.math.BigInteger(1, privateKey),
                        org.bouncycastle.jce.ECNamedCurveTable.getParameterSpec("secp256k1")
                ),
                BouncyCastleProvider.CONFIGURATION
        );
        char[] passCharArray = DigestUtil.sha256Hex(password).toCharArray();
        ByteArrayOutputStream os = new ByteArrayOutputStream();
        Certificate[] certs = new Certificate[] {dummyCert};
        KeyStore ks = KeyStore.getInstance("PKCS12", "BC");
        ks.load(null);
        ks.setKeyEntry("key", pk, passCharArray, certs);
        ks.store(os, passCharArray);
        os.close();
        return Numeric.toHexStringNoPrefix(os.toByteArray());

    }

    public String decryptPrivateKeyByP12(String password, String encryptPrivateKey) throws Exception {
        byte[] toByteArray = Numeric.hexStringToByteArray(encryptPrivateKey);
        char[] passCharArray = DigestUtil.sha256Hex(password).toCharArray();
        KeyStore ks = KeyStore.getInstance("PKCS12", "BC");
        ks.load(new ByteArrayInputStream(toByteArray), passCharArray);
        BCECPrivateKey k = (BCECPrivateKey)ks.getKey("key", passCharArray);
        return Numeric.toHexStringNoPrefix(Numeric.toBytesPadded(k.getD(), 32));
    }
}
