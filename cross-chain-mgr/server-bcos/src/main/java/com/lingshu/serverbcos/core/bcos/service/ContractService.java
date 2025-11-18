package com.lingshu.serverbcos.core.bcos.service;

import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.lingshu.serverbcos.common.api.ApiException;
import com.lingshu.serverbcos.common.api.OpenAPIResp;
import com.lingshu.serverbcos.core.bcos.contract.solidity.BusinessContract;
import com.lingshu.serverbcos.core.bcos.contract.solidity.BusinessContract2;
import com.lingshu.serverbcos.core.bcos.contract.solidity.HelloWorld;
import com.lingshu.serverbcos.core.enums.ContractTypeEnum;
import com.lingshu.serverbcos.dto.BusinessContractRequest;
import lombok.extern.slf4j.Slf4j;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;


import java.io.File;
import java.util.HashMap;
import java.util.Map;
import org.apache.commons.io.FileUtils;
import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;

/**
 * @author: derrick
 * @since: 2025-08-26
 */
@Slf4j
@Service
public class ContractService {
    @Autowired
    private BcosBusinessContractService3 bcosBusinessContractService;

    public Map<String, Object> deploy(BusinessContractRequest request) throws Exception {
        ContractTypeEnum contractTypeEnum = request.getContractType();
        Map<String, Object> map = new HashMap<>();
//        switch (contractTypeEnum) {
////            case LINGSHU_SOLIDITY:
////                lingshuBusinessContractService.deploy(request.getContractName(), request.getVersion(), request.getFilePath());
////            break;
//
//            case BCOS_SOLIDITY:
//                String filePath = request.getFilePath();
//                String binary = FileUtils.readFileToString(new File(filePath), StandardCharsets.UTF_8);
//                log.info("binary: {}", binary);
//                TransactionReceipt transactionReceipt = bcosBusinessContractService.deploy(binary);
//                map.put("address", transactionReceipt.getContractAddress());
//                map.put("txHash", transactionReceipt.getTransactionHash());
//                break;
//            default:
//                throw new ApiException("非法参数");
//        }

        return map;
    }

    public Map<String, Object> call(String param) throws Exception {
        JSONObject jsonObject = JSONUtil.parseObj(param);
        if (!jsonObject.containsKey("origin")) {
            throw new ApiException("非法参数");
        }
        JSONObject originObject = jsonObject.getJSONObject("origin");
        if (!originObject.containsKey("credentialSubject")) {
            throw new ApiException("非法参数");
        }
        JSONObject credentialSubject = originObject.getJSONObject("credentialSubject");
        String resourceName = credentialSubject.getStr("resource_name");
        String gatewayId = credentialSubject.getStr("gateway_id");
        String chainRid = credentialSubject.getStr("chain_rid");
        String contractType = credentialSubject.getStr("contract_type");
        String contractAddress = credentialSubject.getStr("contract_address");
        String contractFunc = credentialSubject.getStr("contract_func");
        String funcParams = credentialSubject.getStr("func_params");
        String key = credentialSubject.getStr("key");

        ContractTypeEnum contractTypeEnum = ContractTypeEnum.getByCode(contractType);
        Map<String, Object> map = new HashMap<>();
//        Map<String, Object> map = new HashMap<>();
//        switch (contractTypeEnum) {
////            case LINGSHU_SOLIDITY:
////                lingshuBusinessContractService.deploy(request.getContractName(), request.getVersion(), request.getFilePath());
////            break;
//            case BCOS_SOLIDITY:
//                BusinessContract2 businessContract = bcosBusinessContractService.load(contractAddress);
//                TransactionReceipt transactionReceipt = businessContract.call(contractFunc, param);
//                log.info("保存成功: {}", transactionReceipt);
//                map.put("txHash", transactionReceipt.getTransactionHash());
////                String info = businessContract.get();
////                log.info("查询成功: {}", info);
//                break;
//            default:
//                throw new ApiException("非法参数");
//        }
        return map;
    }


}
