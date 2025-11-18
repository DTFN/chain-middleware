package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.core.ethereum.contract.BusinessContract2;
import com.lingshu.server.core.ethereum.service.EthereumBusinessContractService2;
import com.lingshu.server.dto.BusinessContractRequest;
import com.lingshu.server.dto.LingshuBusinessContractRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/v1/openapi/ethereum-business/")
@Slf4j
public class EthereumBusinessContractController {
    @Autowired
    private EthereumBusinessContractService2 ethereumBusinessContractService;

    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody BusinessContractRequest request) {
        try {
            TransactionReceipt transactionReceipt = ethereumBusinessContractService.deploy(request.getFilePath());
            Map<String, Object> map = new HashMap<>();
            map.put("address", transactionReceipt.getContractAddress());
            map.put("txHash", transactionReceipt.getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("deploy error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("load")
    public OpenAPIResp load(@Valid @RequestBody LingshuBusinessContractRequest request) {
        try {
            BusinessContract2 businessContract = ethereumBusinessContractService.load(request.getContractAddress());
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
//            TransactionReceipt receipt = ethereumBusinessContractService.saveBusinessDetails(request);
//            Map<String, Object> map = new HashMap<>();
//            map.put("txHash", receipt.getTransactionHash());
//            return OpenAPIRespBuilder.success(map);
//        } catch (Exception e) {
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }
//
//    @GetMapping("get")
//    public OpenAPIResp getBusinessDetails(@RequestParam String businessAddress) {
//        try {
//            String businessDetails = ethereumBusinessContractService.getBusinessDetails(businessAddress);
//            log.info("businessDetails:{}", businessDetails);
//            return OpenAPIRespBuilder.success(businessDetails);
//        } catch (Exception e) {
//            e.printStackTrace();
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }

}
