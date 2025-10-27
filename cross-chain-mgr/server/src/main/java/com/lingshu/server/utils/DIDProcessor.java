package com.lingshu.server.utils;

import cn.hutool.json.JSONUtil;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.lingshu.server.dto.DidDto;

import java.util.Arrays;

/**
 * VC0模板填充与contentHash计算工具
 * 功能：1. 填充VC0中的%s占位符 2. 标准化JSON格式 3. 计算符合预期的contentHash
 */
public class DIDProcessor {

    // 1. VC0模板（含%s占位符，注意修正reource_name为resource_name）
    public static final String DID_TEMPLATE = "{\n" +
            "    \"@context\":\n" +
            "    [\n" +
            "        \"https://www.w3.org/ns/did/v1\",\n" +
            "        \"https://w3c-ccg.github.io/lds-jws2020/contexts/lds-jws2020-v1.json\",\n" +
            "        {\n" +
            "            \"Ed25519VerificationKey2018\": \"https://w3id.org/security#Ed25519VerificationKey2018\"\n" +
//            "            \"publicKeyBase58\": \"https://w3id.org/security#publicKeyBase58\",\n" +
//            "            \"publicKeyJwk\":\n" +
//            "            {\n" +
//            "                \"@id\": \"https://w3id.org/security#publicKeyJwk\",\n" +
//            "                \"@type\": \"@json\"\n" +
//            "            }\n" +
            "        }\n" +
            "    ],\n" +
            "    \"id\": \"%s\",\n" +
            "    \"verificationMethod\":\n" +
            "    [\n" +
            "        {\n" +
            "            \"id\": \"%s\",\n" +
            "            \"type\": \"Ed25519VerificationKey2018\",\n" +
            "            \"controller\": \"%s\",\n" +
            "            \"publicKeyHex\": \"%s\"\n" +
            "        },\n" +
            "        {\n" +
            "            \"id\": \"%s\",\n" +
            "            \"type\": \"EcdsaSecp256k1VerificationKey2019\",\n" +
            "            \"controller\": \"%s\",\n" +
            "            \"address\": \"%s\"\n" +
            "        }\n" +
            "    ],\n" +
//            "    ]\n" +
            "    \"assertionMethod\":\n" +
            "    [\n" +
            "        \"%s\",\n" +
            "        \"%s\"\n" +
            "    ],\n" +
            "    \"authentication\":\n" +
            "    [\n" +
            "        \"%s\",\n" +
            "        \"%s\"\n" +
            "    ],\n" +
            "    \"capabilityDelegation\":\n" +
            "    [\n" +
            "        \"%s\",\n" +
            "        \"%s\"\n" +
            "    ],\n" +
            "    \"capabilityInvocation\":\n" +
            "    [\n" +
            "        \"%s\",\n" +
            "        \"%s\"\n" +
            "    ]\n" +
            "}";

    // 2. JSON工具（处理func_params序列化和VC标准化）
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper()
            .disable(SerializationFeature.INDENT_OUTPUT) // 禁用缩进，生成压缩JSON
            .disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);

    public static String fillAndStandardizeDID(
            String did,
            String publicKeyHex,
            String address
    ) throws JsonProcessingException {
        String mothod0 = did + "#key0";
        String mothod1 = did + "#key1";
        // 填充模板并标准化VC
        String filledWithIndent = String.format(DID_TEMPLATE,
                did,
                mothod0,
                did,
                publicKeyHex,
                mothod1,
                did,
                address,
                mothod0,
                mothod1,
                mothod0,
                mothod1,
                mothod0,
                mothod1,
                mothod0,
                mothod1
        );

        // 标准化VC：去除缩进和换行，确保与ether-signer一致
        return OBJECT_MAPPER.readTree(filledWithIndent).toString();
    }

    public static String fillAndStandardizeDIDObj(
            String did,
            String publicKeyHex,
            String address,
            String nationality
    ) throws JsonProcessingException {
        String ethFunc = "#key1";
        String ethId = did + ethFunc;

        DidDto didDto = new DidDto();
        didDto
                .setContext(Arrays.asList("https://www.w3.org/ns/did/v1"))
                .setId(did)
//                .setAlsoknownas(nationality)
                .setVerificationMethod(Arrays.asList(
                        new DidDto.VerificationMethodEth()
                                .setId(ethId)
                                .setType("EcdsaSecp256k1VerificationKey2019")
                                .setController(did)
                                .setAddress(address)
                ))
                .setAssertionMethod(Arrays.asList(
                        ethId
                ))
        ;

        // 标准化VC：去除缩进和换行，确保与ether-signer一致
        return OBJECT_MAPPER.readTree(JSONUtil.toJsonStr(didDto)).toString();
    }
}