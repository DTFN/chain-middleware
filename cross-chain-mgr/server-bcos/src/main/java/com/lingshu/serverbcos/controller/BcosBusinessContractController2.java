package com.lingshu.serverbcos.controller;

import com.lingshu.serverbcos.common.api.ApiErrorCode;
import com.lingshu.serverbcos.common.api.OpenAPIResp;
import com.lingshu.serverbcos.common.api.OpenAPIRespBuilder;
import com.lingshu.serverbcos.core.bcos.contract.solidity.HelloWorld;
import com.lingshu.serverbcos.core.bcos.service.BcosBusinessContractService;
import com.lingshu.serverbcos.dto.BusinessContractSaveRequest;
import com.lingshu.serverbcos.dto.LingshuBusinessContractRequest;
import lombok.extern.slf4j.Slf4j;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/v1/openapi/bcos-business2/")
@Slf4j
public class BcosBusinessContractController2 {
    @Autowired
    private BcosBusinessContractService bcosBusinessContractService;
    @PostMapping("deploy")
    public OpenAPIResp deploy() {
        try {
            TransactionReceipt transactionReceipt = bcosBusinessContractService.deploy();
            Map<String, Object> map = new HashMap<>();
            map.put("address", transactionReceipt.getContractAddress());
            map.put("txHash", transactionReceipt.getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("load")
    public OpenAPIResp load(@Valid @RequestBody LingshuBusinessContractRequest request) {
        try {
            HelloWorld businessContract = bcosBusinessContractService.load(request.getContractAddress());
            log.info("businessContract:{}", businessContract);
            String details = businessContract.get();
            log.info("details:{}", details);
            businessContract.set("1qazxsw2");
            details = businessContract.get();
            log.info("details:{}", details);
            return OpenAPIRespBuilder.success();
        } catch (Exception e) {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("save")
    public OpenAPIResp saveBusinessDetails(@Valid @RequestBody BusinessContractSaveRequest request) {
        try {
            TransactionReceipt receipt = bcosBusinessContractService.saveBusinessDetails(request);
            Map<String, Object> map = new HashMap<>();
            map.put("txHash", receipt.getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @GetMapping("get")
    public OpenAPIResp getBusinessDetails(@RequestParam String businessAddress) {
        try {
            String businessDetails = bcosBusinessContractService.getBusinessDetails(businessAddress);
            log.info("businessDetails:{}", businessDetails);
            return OpenAPIRespBuilder.success(businessDetails);
        } catch (Exception e) {
            e.printStackTrace();
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

}
