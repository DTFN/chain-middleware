package com.lingshu.server.key.util;

public class SeedGenerator {

    private static final String MNEMONIC = "mnemonic";

    public byte[] generateSeed(String mnemonic, String passphrase){
        if(passphrase == null) {
            passphrase = "";
        }
        String salt = MNEMONIC + passphrase;
        byte[] saltBytes = salt.getBytes();
        byte[] seed = PBKDF2WithHmacSha512.INSTANCE.kdf(mnemonic.toCharArray(), saltBytes);
        return seed;
    }

}