package com.lingshu.server.key.util;

import java.nio.charset.StandardCharsets;

public class MasterKeyGenerator {

    public static final String Seed = "seed";


    public byte[] toMasterKey(byte[] data, String prefix){
        byte[] shakey = prefix.getBytes(StandardCharsets.UTF_8);
        return HmacSha512.INSTANCE.macHash(shakey, data);
    }


    public byte[] toMasterKey(byte[] data){
        return this.toMasterKey(data, Seed);
    }
}
