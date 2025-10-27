package com.lingshu.server.utils;

import cn.hutool.json.JSONUtil;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.TreeMap;

public class JsonOrderProcessor {
    /**
     * 递归将Map按键的字典顺序排序，包括嵌套的Map结构
     */
    private static Map<String, Object> sortMapByKey(Map<String, Object> map) {
        // TreeMap会自动按键的自然顺序（字典顺序）排序
        TreeMap<String, Object> sortedMap = new TreeMap<>(map);

        // 处理嵌套的Map
        for (Map.Entry<String, Object> entry : sortedMap.entrySet()) {
            Object value = entry.getValue();
            if (value instanceof Map) {
                // 递归排序嵌套的Map
                @SuppressWarnings("unchecked")
                Map<String, Object> nestedMap = (Map<String, Object>) value;
                entry.setValue(sortMapByKey(nestedMap));
            }
        }

        return sortedMap;
    }

    public static String convert(String vc) {
        try {
            // 创建ObjectMapper实例
            ObjectMapper objectMapper = new ObjectMapper();
            // 配置为 pretty print 格式，便于查看
            //objectMapper.enable(SerializationFeature.INDENT_OUTPUT);

            // 1. 将JSON字符串解析为有序Map（LinkedHashMap会保持插入顺序）
            Map<String, Object> originalMap = objectMapper.readValue(vc,
                    new TypeReference<LinkedHashMap<String, Object>>() {});

            // 2. 将Map转换为TreeMap以实现按字典顺序排序键
            // 递归处理嵌套的Map结构
            Map<String, Object> sortedMap = sortMapByKey(originalMap);

            // 3. 将排序后的Map序列化为JSON字符串（不进行格式化）
            String sortedJsonString = objectMapper.writeValueAsString(sortedMap);

            return sortedJsonString;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public static String convert(Object obj) {
        try {
            // 反序列化
            String vc  = JSONUtil.toJsonStr(obj);

            // 创建ObjectMapper实例
            ObjectMapper objectMapper = new ObjectMapper();
            // 配置为 pretty print 格式，便于查看
            //objectMapper.enable(SerializationFeature.INDENT_OUTPUT);

            // 1. 将JSON字符串解析为有序Map（LinkedHashMap会保持插入顺序）
            Map<String, Object> originalMap = objectMapper.readValue(vc,
                    new TypeReference<LinkedHashMap<String, Object>>() {});

            // 2. 将Map转换为TreeMap以实现按字典顺序排序键
            // 递归处理嵌套的Map结构
            Map<String, Object> sortedMap = sortMapByKey(originalMap);

            // 3. 将排序后的Map序列化为JSON字符串（不进行格式化）
            String sortedJsonString = objectMapper.writeValueAsString(sortedMap);

            return sortedJsonString;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }
}
    