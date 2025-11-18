package com.lingshu.serverbcos.utils;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;

/**
 * VC0模板填充与contentHash计算工具
 * 功能：1. 填充VC0中的%s占位符 2. 标准化JSON格式 3. 计算符合预期的contentHash
 */
public class Vc0Processor {

    // 1. VC0模板（含%s占位符，注意修正reource_name为resource_name）
    public static final String VC0_TEMPLATE = "{\n" +
            "  \"@context\": [\n" +
            "    \"https://www.w3.org/2018/credentials/v1\"\n" +
            "  ],\n" +
            "  \"credentialSubject\": {\n" +
            "    \"chain_rid\": \"%s\",\n" +
            "    \"contract_address\": \"%s\",\n" +
            "    \"contract_func\": \"%s\",\n" +
            "    \"contract_type\": \"%s\",\n" +
            //"    \"func_params\": %s,\n" + // 无引号：填充JSON字符串
            "    \"gateway_id\": \"%s\",\n" +
//            "    \"param_name\": \"%s\",\n" +
            "    \"resource_name\": \"%s\"\n" + // 修正拼写：reource_name → resource_name
            "  },\n" +
            "  \"id\": \"https://uni.example.com/credentials/12345\",\n" +
            "  \"issuanceDate\": \"%s\",\n" +
            "  \"issuer\": \"https://uni.example.com/organization/1\",\n" +
            "  \"type\": [\n" +
//            "    \"VerifiableCredential\",\n" +
            "    \"VerifiableCredential\"\n" +
//            "    \"UniversityDegreeCredential\"\n" +
            "  ]\n" +
            "}";

    // 2. JSON工具（处理func_params序列化和VC标准化）
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper()
            .disable(SerializationFeature.INDENT_OUTPUT) // 禁用缩进，生成压缩JSON
            .disable(SerializationFeature.FAIL_ON_EMPTY_BEANS);

    /**
     * 步骤1：填充VC0模板的占位符
     *
     * @param chainRid        链标识（如chain1）
     * @param contractAddress 合约地址（如contract01-daiding）
     * @param contractFunc    合约方法（如erc20MintVcs）
     * @param contractType    合约类型（如1，字符串格式）
     * @param gatewayId       网关ID（如0，字符串格式）
     * @param resourceName    资源名（如chain1:contract01）
     * @param issuanceDate    签发日期（如2025-05-20）
     * @return 标准化VC字符串（无缩进、无换行）
     */
    public static String fillAndStandardizeVc0(
            String chainRid,
            String contractAddress,
            String contractFunc,
            String contractType,
            String gatewayId,
            String resourceName,
            String issuanceDate
    ) throws JsonProcessingException {
        // 填充模板并标准化VC
        String filledWithIndent = String.format(VC0_TEMPLATE,
                chainRid,
                contractAddress,
                contractFunc,
                contractType,
                gatewayId,
                resourceName,
                issuanceDate
        );

        // 标准化VC：去除缩进和换行，确保与ether-signer一致
        return OBJECT_MAPPER.readTree(filledWithIndent).toString();
    }
}