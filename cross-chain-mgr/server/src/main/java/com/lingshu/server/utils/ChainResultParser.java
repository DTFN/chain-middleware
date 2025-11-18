package com.lingshu.server.utils;

import com.google.protobuf.ByteString;
import java.nio.charset.StandardCharsets;

public class ChainResultParser {

    /**
     * 解析长安链返回结果，智能截断无效尾随数据
     */
    @SuppressWarnings("unchecked")
    public static <T> T parseResult(ByteString result, Class<T> tClass) {
        // 步骤1：截取业务数据段（假设前64字节为头部，需跳过）
        ByteString dataBytes = result;//.substring(64);
        byte[] rawBytes = dataBytes.toByteArray();

        // 步骤2：智能截断无效尾随数据（核心逻辑）
        int validLength = findValidLength(rawBytes);
        byte[] validBytes = new byte[validLength];
        System.arraycopy(rawBytes, 0, validBytes, 0, validLength);

        // 步骤3：按目标类型返回
        if (tClass.equals(String.class)) {
            return (T) new String(validBytes, StandardCharsets.UTF_8);
        }
        // 可扩展其他类型解析（如字节数组、自定义对象等）
        return null;
    }

    /**
     * 从后往前找，确定有效字节的长度（非0x00的部分为有效）
     */
    private static int findValidLength(byte[] bytes) {
        for (int i = bytes.length - 1; i >= 0; i--) {
            // 遇到第一个非0x00的字节，认为有效长度到i+1（包含i）
            if (bytes[i] != 0x00) {
                return i + 1;
            }
        }
        return 0; // 极端情况：所有字节都是0x00，返回空
    }
}