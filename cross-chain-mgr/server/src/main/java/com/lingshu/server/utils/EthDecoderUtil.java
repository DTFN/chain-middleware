package com.lingshu.server.utils;

import com.lingshu.chain.sdk.util.HexUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.util.StringUtils;

import java.math.BigInteger;

@Slf4j
public class EthDecoderUtil {
    public static String decodeOneString(String outputHex, boolean hasFuncName) {
        if (outputHex == null) {
            return null;
        }

        int startIndex = hasFuncName ? 8 : 0;
        String outputHexNoPrefix = HexUtil.trimPrefix(outputHex);

        //
        if (!StringUtils.hasText(outputHexNoPrefix)) {
            return "";
        }

        BigInteger index = new BigInteger(outputHexNoPrefix.substring(startIndex, startIndex + 64), 16);
        BigInteger stringLength = new BigInteger(outputHexNoPrefix.substring(startIndex + 64, startIndex + 64 + index.intValue() * 2), 16);
        log.info("startIndex: {}, index: {}, stringLength: {}, outputHexNoPrefix: {}", startIndex, index, stringLength, outputHexNoPrefix);
        String substring = outputHexNoPrefix.substring(
                startIndex + 64 + index.intValue() * 2,
                startIndex + 64 + index.intValue() * 2 + stringLength.intValue() * 2
        );
        if (substring == null) {
            return null;
        }
        if (substring.length() == 0) {
            return "";
        }
        String message = new String(cn.hutool.core.util.HexUtil.decodeHex(substring));
        return message;
    }
}
