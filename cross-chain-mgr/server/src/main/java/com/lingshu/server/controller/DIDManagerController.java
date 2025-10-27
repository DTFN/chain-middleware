package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.core.web3j.contract.DIDManager;
import com.lingshu.server.core.web3j.service.DIDManagerService;
import com.lingshu.server.dto.DIDContractRequest;
import com.lingshu.server.dto.DIDCreateRequest;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.validation.Valid;
import java.util.HashMap;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@RestController
@RequestMapping("/v1/openapi/did-manager/")
@Slf4j
public class DIDManagerController {
    @Autowired
    private DIDManagerService didManagerService;

    @PostMapping("deploy")
    public OpenAPIResp deploy(@Valid @RequestBody DIDContractRequest request) throws Exception {
        try {
            DIDManager manager = didManagerService.deploy(request.getContractName());
            Map<String, Object> map = new HashMap<>();
            map.put("address", manager.getDeployInfo().getContract().getAddress());
            map.put("txHash", manager.getDeployInfo().getTransactionReceipt().getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("deploy DIDManager error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("load")
    public OpenAPIResp load(@Valid @RequestBody DIDContractRequest request) {
        try {
            DIDManager manager = didManagerService.load(request.getContractName());
            log.info("DIDManager: {}", manager);
            return OpenAPIRespBuilder.success();
        } catch (Exception e) {
            log.error("load contract error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    // 查询DID信息
    @GetMapping("query")
    public OpenAPIResp queryDID(@RequestParam String did) {
        try {
            String result = didManagerService.getDIDDetails(did);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("queryDID error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    // 是否存在
    @GetMapping("exists")
    public OpenAPIResp existsDID(@RequestParam String did) {
        try {
            boolean result = didManagerService.existsDID(did);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("existsDID error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    //创建
    @PostMapping("create")
    public OpenAPIResp createDID(@Valid @RequestBody DIDCreateRequest request) {
        try {
            TransactionReceipt receipt = didManagerService.createDID(request);
            Map<String, Object> map = new HashMap<>();
            map.put("txHash", receipt.getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("create error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    //更新
    @PostMapping("update")
    public OpenAPIResp updateDID(@Valid @RequestBody DIDCreateRequest request) {
        try {
            TransactionReceipt receipt = didManagerService.updateDID(request);
            Map<String, Object> map = new HashMap<>();
            map.put("txHash", receipt.getTransactionHash());
            return OpenAPIRespBuilder.success(map);
        } catch (Exception e) {
            log.error("update error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @GetMapping("echo")
    public OpenAPIResp echo(@RequestParam String did) {
        try {
            String result = didManagerService.echo(did);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("echo error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }
}
