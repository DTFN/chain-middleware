package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.core.web3j.contract.BusinessContract;
import com.lingshu.server.core.web3j.service.BusinessContractService;
import com.lingshu.server.dto.BusinessContractRequest;
import com.lingshu.server.dto.BusinessContractSaveRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@RestController
@RequestMapping("/v1/openapi/business-contract/")
@Slf4j
public class BusinessContractController {
    @Autowired
    private BusinessContractService businessContractService;

    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody BusinessContractRequest request) {
        try {
            //长安链 solidity
            BusinessContract businessContract = businessContractService.deploy(request.getContractName(), request.getVersion(), request.getFilePath());
            Map<String, Object> map = new HashMap<>();
            map.put("address", businessContract.getDeployInfo().getContract().getAddress());
            map.put("txHash", businessContract.getDeployInfo().getTransactionReceipt().getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("deploy error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("load")
    public OpenAPIResp load(@Valid @RequestBody BusinessContractRequest request) {
        try {
            BusinessContract businessContract = businessContractService.load(request.getContractName());
            log.info("businessContract:{}", businessContract);
            return OpenAPIRespBuilder.success();
        } catch (Exception e) {
            log.error("load error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

//    @PostMapping("save")
//    public OpenAPIResp saveBusinessDetails(@Valid @RequestBody BusinessContractSaveRequest request) {
//        try {
//            TransactionReceipt receipt = businessContractService.saveBusinessDetails(request);
//            Map<String, Object> map = new HashMap<>();
//            map.put("txHash", receipt.getTransactionHash());
//            return OpenAPIRespBuilder.success(map);
//        } catch (Exception e) {
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }

//    @GetMapping("get")
//    public OpenAPIResp getBusinessDetails(@RequestParam String businessAddress) {
//        try {
//            String result = businessContractService.getBusinessDetails(businessAddress);
//            return OpenAPIRespBuilder.success(result);
//        } catch (Exception e) {
//            log.error("getBusinessDetails error:{}", e.getMessage(), e);
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }

}
