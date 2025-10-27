package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.common.metrics.ResourceDomainMetric;
import com.lingshu.server.core.web3j.contract.ResourceDomain;
import com.lingshu.server.core.web3j.service.ResourceDomainService;
import com.lingshu.server.dto.ResourceDomainContractRequest;
import com.lingshu.server.dto.ResourceDomainSaveRequest;
import com.lingshu.server.dto.resp.busi.TxHashResp;
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
@RequestMapping("/v1/openapi/resource-domain/")
@Slf4j
public class ResourceDomainController {
    @Autowired
    private ResourceDomainService resourceDomainService;

    @Autowired
    private ResourceDomainMetric resourceDomainMetric;

    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody ResourceDomainContractRequest request) throws Exception {
        try {
            ResourceDomain resourceDomain = resourceDomainService.deploy(request.getContractName());
            Map<String, Object> map = new HashMap<>();
            map.put("address", resourceDomain.getDeployInfo().getContract().getAddress());
            map.put("txHash", resourceDomain.getDeployInfo().getTransactionReceipt().getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("deploy error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("load")
    public OpenAPIResp<Void> load(@Valid @RequestBody ResourceDomainContractRequest request) {
        try {
            ResourceDomain resourceDomain = resourceDomainService.load(request.getContractName());
            return OpenAPIRespBuilder.success();
        } catch (Exception e) {
            log.error("load error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("create")
    public OpenAPIResp<TxHashResp> createDomainDetails(@Valid @RequestBody ResourceDomainSaveRequest request) {
        try {
            TransactionReceipt receipt = resourceDomainService.saveDomainDetails(request);
            TxHashResp txHashResp = new TxHashResp();
            txHashResp.setTxHash(receipt.getTransactionHash());
            resourceDomainMetric.incrementSuccess(request.getGatewayId(), request.getChainRid(), request.getContractType());
            return OpenAPIRespBuilder.success(txHashResp);
        } catch (Exception e) {
            log.error("createDomainDetails error:{}", e.getMessage(), e);
            resourceDomainMetric.incrementFail(request.getGatewayId(), request.getChainRid(), request.getContractType());
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("update")
    public OpenAPIResp<TxHashResp> updateDomainDetails(@Valid @RequestBody ResourceDomainSaveRequest request) {
        try {
            TransactionReceipt receipt = resourceDomainService.updateDomainDetails(request);
            TxHashResp txHashResp = new TxHashResp();
            txHashResp.setTxHash(receipt.getTransactionHash());
            return OpenAPIRespBuilder.success(txHashResp);
        } catch (Exception e) {
            log.error("updateDomainDetails error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @GetMapping("get")
    public OpenAPIResp<String> getDomainDetails(@RequestParam String domainAddress) {
        try {
            String result = resourceDomainService.getDomainDetails(domainAddress);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("getDomainDetails error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    // 是否存在
    @GetMapping("exists")
    public OpenAPIResp<Boolean> doesDomainExist(@RequestParam String domainAddress) {
        try {
            boolean result = resourceDomainService.doesDomainExist(domainAddress);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("exists error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }
}
