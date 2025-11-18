package com.lingshu.serverbcos.core.bcos.service;

import cn.hutool.core.date.DateUtil;
import cn.hutool.core.lang.Assert;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.alibaba.fastjson.JSON;
import com.lingshu.serverbcos.common.api.ApiException;
import com.lingshu.serverbcos.core.bcos.contract.solidity.BusinessContract2;
import com.lingshu.serverbcos.core.enums.ContractTypeEnum;
import com.lingshu.serverbcos.dto.*;
import com.lingshu.serverbcos.utils.Ed25519Signer;
import com.lingshu.serverbcos.utils.EthVCSigner;
import com.lingshu.serverbcos.utils.JsonOrderProcessor;
import com.lingshu.serverbcos.utils.Vc0Processor;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.fisco.bcos.sdk.client.Client;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.fisco.bcos.sdk.utils.Numeric;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.web3j.crypto.Sign;
import org.web3j.crypto.Sign;

import java.util.HashMap;
import java.util.Map;
import java.util.Objects;

import static com.lingshu.serverbcos.core.enums.ContractTypeEnum.CHAINMAKER_WASM;

/**
 * @author: derrick
 * @since: 2025-09-04
 */
@Slf4j
@Service
public class BcosBusinessContractService3 {
    @Autowired
    private Client client;

//    @Autowired
//    private ResourceDomainService resourceDomainService;

//    private BusinessContract2 businessContract;
    public Map<String, Object> call(ContractRequest request) throws Exception {
//        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
//        ResourceDomainInfo originResourceDomainInfo = request.getOriginResourceDomainInfo();
//        //将ResourceDomain合约上的信息，写入origin/target vc中
//        JSONObject originObject1 = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin());
//        JSONObject vcObject = new JSONObject();
//        vcObject.put("origin", originObject1);

        String contractAddress = request.getContractAddress();
        String contractFunc = request.getContractFunc();
        String param = request.getFuncParams();
        log.info("contractAddress: {}, contractFunc: {}, param: {}", contractAddress, contractFunc, param);
//        String paramsStr = JsonOrderProcessor.convert(param);
//        log.info("vc: {}", paramsStr);
//        //调用origin合约方法
//        //Map<String, Object> result = contractService.call(paramsStr);
//
//        JSONObject jsonObject = JSONUtil.parseObj(paramsStr);
//        if (!jsonObject.containsKey("origin")) {
//            throw new ApiException("非法参数");
//        }
//        JSONObject originObject = jsonObject.getJSONObject("origin");
//        if (!originObject.containsKey("credentialSubject")) {
//            throw new ApiException("非法参数");
//        }
//        JSONObject credentialSubject = originObject.getJSONObject("credentialSubject");
//        String resourceName = credentialSubject.getStr("resource_name");
//        String gatewayId = credentialSubject.getStr("gateway_id");
//        String chainRid = credentialSubject.getStr("chain_rid");
//        String contractType = credentialSubject.getStr("contract_type");
//        String contractAddress = credentialSubject.getStr("contract_address");
//        String contractFunc = credentialSubject.getStr("contract_func");
//        String funcParams = credentialSubject.getStr("func_params");
//        String key = credentialSubject.getStr("key");
//        String paramName = credentialSubject.getStr("param_name");


        BusinessContract2 businessContract2 = BusinessContract2.load(/*request.getBinary(),*/
                request.getContractAddress(), client, client.getCryptoSuite().getCryptoKeyPair());
        TransactionReceipt transactionReceipt2 = null;
        if ("createDID".equals(contractFunc) || "updateDID".equals(contractFunc)) {
//            param
            JSONObject funcParams = JSONUtil.parseObj(param);
            String did = funcParams.getStr("did");
            String didDocument = funcParams.getStr("didDocument");
            transactionReceipt2 = businessContract2.call(contractFunc, did, didDocument);
        } else {
            transactionReceipt2 = businessContract2.call(contractFunc, param);
        }

        // 校验status(0x开头的16进制)
        String status = transactionReceipt2.getStatus();
        String statusWithoutPrefix = status.toLowerCase().replaceAll("0x", "");
        Assert.isTrue("0".equals(statusWithoutPrefix), () -> new ApiException("合约执行异常, 合约地址: " + request.getContractAddress() + ", 错误码: " + status));

        log.info("保存成功: {}", transactionReceipt2);
        Map<String, Object> map = new HashMap<>();
        map.put("txHash", transactionReceipt2.getTransactionHash());

        return map;
    }

    public String get(ContractRequest request) throws Exception {
//        //根据resource_name,获取ResourceDomain合约上的信息（orign/target）
//        ResourceDomainInfo originResourceDomainInfo = request.getOriginResourceDomainInfo();
//        //将ResourceDomain合约上的信息，写入origin/target vc中
//        JSONObject originObject = dealVc(originResourceDomainInfo, request.getPrivateKeyHex(), request.getOrigin());
//        JSONObject vcObject = new JSONObject();
//        vcObject.put("origin", originObject);

        String contractAddress = request.getContractAddress();
        String contractFunc = request.getContractFunc();
        String param = request.getFuncParams();
        log.info("contractAddress: {}, contractFunc: {}, param: {}", contractAddress, contractFunc, param);
//        String paramsStr = JsonOrderProcessor.convert(param);
//        log.info("vc: {}", paramsStr);
////        //调用origin合约方法
////        String result = contractService.get(paramsStr);
//
//        JSONObject jsonObject = JSONUtil.parseObj(paramsStr);
//        if (!jsonObject.containsKey("origin")) {
//            throw new ApiException("非法参数");
//        }
//        JSONObject originObject = jsonObject.getJSONObject("origin");
//        if (!originObject.containsKey("credentialSubject")) {
//            throw new ApiException("非法参数");
//        }
//        JSONObject credentialSubject = originObject.getJSONObject("credentialSubject");
//        String resourceName = credentialSubject.getStr("resource_name");
//        String gatewayId = credentialSubject.getStr("gateway_id");
//        String chainRid = credentialSubject.getStr("chain_rid");
//        String contractType = credentialSubject.getStr("contract_type");
//        String contractAddress = credentialSubject.getStr("contract_address");
//        String contractFunc = credentialSubject.getStr("contract_func");
//        String funcParams = credentialSubject.getStr("func_params");
//        String key = credentialSubject.getStr("key");
//        String paramName = credentialSubject.getStr("param_name");


        BusinessContract2 businessContract2 = BusinessContract2.load(/*request.getBinary(),*/
                contractAddress, client, client.getCryptoSuite().getCryptoKeyPair());

        String
            result = businessContract2.get(contractFunc, param);

        log.info("查询成功: {}", result);

//        TransactionReceipt transactionReceipt2 = null;
//        if ("createDID".equals(contractFunc) || "updateDID".equals(contractFunc)) {
//            JSONObject funcParams1 = credentialSubject.getJSONObject("func_params");
//            String did = funcParams1.getStr("did");
//            String didDocument = funcParams1.getStr("didDocument");
//            transactionReceipt2 = businessContract2.call(contractFunc, did, didDocument);
//        } else {
//            transactionReceipt2 = businessContract2.call(contractFunc, paramsStr);
//        }
//        log.info("保存成功: {}", transactionReceipt2);
//        Map<String, Object> map = new HashMap<>();
//        map.put("txHash", transactionReceipt2.getTransactionHash());
        return result;
    }
}
