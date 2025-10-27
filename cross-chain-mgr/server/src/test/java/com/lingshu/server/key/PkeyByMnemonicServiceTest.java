package com.lingshu.server.key;

import org.junit.jupiter.api.Test;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.ECKeyPair;
import org.web3j.utils.Numeric;

class PkeyByMnemonicServiceTest {

    @Test
    void generatePrivateKeyByMnemonic() throws Exception {

        PkeyByMnemonicService service = new PkeyByMnemonicService();

//        String mnemonic = service.createMnemonic();
        String mnemonic = "wreck edge until recycle clean benefit lunar simple announce vacuum flavor suspect";
        String password = "1234567";

        System.out.println("助记词:" + mnemonic);

        ECKeyPair keyPair = service.generatePrivateKeyByMnemonic(mnemonic, password);
        String privateKey = Numeric.toHexStringNoPrefix(keyPair.getPrivateKey());
        String publicKey = Numeric.toHexStringNoPrefix(keyPair.getPublicKey());
        String address = Credentials.create(keyPair).getAddress();
        System.out.println("私钥:" + privateKey);
        System.out.println("公钥:" + publicKey);
        System.out.println("地址:" + address);

        // 使用密码创建P12格式文件，这里转为二进制字符串
        String privateKeyByP12 = service.encryptPrivateKeyByP12(password, Numeric.hexStringToByteArray(privateKey));
        System.out.println("P12加密:" + privateKeyByP12);

        String decryptPrivateKeyByP12 = service.decryptPrivateKeyByP12(password, privateKeyByP12);
        System.out.println("解密后私钥:" + decryptPrivateKeyByP12);
    }
}