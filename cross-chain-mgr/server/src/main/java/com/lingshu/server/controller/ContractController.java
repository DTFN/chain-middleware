package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.core.web3j.service.ContractService;
import com.lingshu.server.dto.BusinessContractRequest;
import com.lingshu.server.dto.BcosContractRequest;
import com.lingshu.server.dto.ContractRequest;
import com.lingshu.server.dto.SetDIDRequest;
import com.lingshu.server.dto.resp.busi.TxHashResp;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@RestController
@RequestMapping("/v1/openapi/contract/")
@Slf4j
public class ContractController {

    @Autowired
    private ContractService contractService;

    //合约部署
    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody BusinessContractRequest request) {
        try {
            Map<String, Object> result = contractService.deploy(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("deploy error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    //合约调用
    @PostMapping("call")
    public OpenAPIResp<TxHashResp> call(@Valid @RequestBody ContractRequest request) {
        try {
            TxHashResp result = contractService.call(request.getParam());
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    //测试
    @PostMapping("setDidAddress")
    public OpenAPIResp setDidAddress(@Valid @RequestBody SetDIDRequest request) {
        try {
            Map<String, Object> result = contractService.setDidAddress(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }
}
