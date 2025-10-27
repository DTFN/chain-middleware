package com.lingshu.server.utils;

import com.lingshu.chain.sdk.util.HexUtil;

import java.math.BigInteger;

public class EthDecoderUtil {
    public static String decodeOneString(String outputHex, boolean hasFuncName) {
        int startIndex = hasFuncName ? 8: 0;
        String outputHexNoPrefix = HexUtil.trimPrefix(outputHex);
        BigInteger index = new BigInteger(outputHexNoPrefix.substring(startIndex, startIndex + 64), 16);
        BigInteger stringLength = new BigInteger(outputHexNoPrefix.substring(startIndex + 64, startIndex + 64 + index.intValue() * 2), 16);
        String message = new String(cn.hutool.core.util.HexUtil.decodeHex(outputHexNoPrefix.substring(startIndex + 64 + index.intValue() * 2, startIndex + 64 + index.intValue() * 2 + stringLength.intValue() * 2)));
        return message;
    }
}
