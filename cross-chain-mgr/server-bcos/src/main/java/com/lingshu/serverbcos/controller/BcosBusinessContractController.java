package com.lingshu.serverbcos.controller;

import com.lingshu.serverbcos.common.api.ApiErrorCode;
import com.lingshu.serverbcos.common.api.OpenAPIResp;
import com.lingshu.serverbcos.common.api.OpenAPIRespBuilder;
import com.lingshu.serverbcos.core.bcos.contract.solidity.HelloWorld;
import com.lingshu.serverbcos.core.bcos.service.BcosBusinessContractService3;
import com.lingshu.serverbcos.dto.BusinessContractSaveRequest;
import com.lingshu.serverbcos.dto.ContractRequest;
import com.lingshu.serverbcos.dto.CrossChainRequest;
import com.lingshu.serverbcos.dto.LingshuBusinessContractRequest;
import lombok.extern.slf4j.Slf4j;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/v1/openapi/bcos-business/")
@Slf4j
public class BcosBusinessContractController {
    @Autowired
    private BcosBusinessContractService3 businessService;

//    @PostMapping("crossChain")
//    public OpenAPIResp crossChain(@Valid @RequestBody CrossChainRequest request) {
//        try {
//            Map<String, Object> result = businessService.crossChain(request);
//            return OpenAPIRespBuilder.success(result);
//        } catch (Exception e) {
//            log.error("crossChain error:{}", e.getMessage(), e);
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
//        }
//    }

    @PostMapping("call")
    public OpenAPIResp call(@Valid @RequestBody ContractRequest request) {
        try {
            Map<String, Object> result = businessService.call(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("get")
    public OpenAPIResp get(@Valid @RequestBody ContractRequest request) {
        try {
            String result = businessService.get(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

}
