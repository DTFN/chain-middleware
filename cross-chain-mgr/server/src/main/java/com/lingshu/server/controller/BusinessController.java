package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.common.metrics.CrossChainTxMetric;
import com.lingshu.server.core.business.BusinessService;
import com.lingshu.server.dto.*;
import com.lingshu.server.dto.resp.busi.TxHashResp;
import io.swagger.annotations.ApiOperation;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import javax.validation.Valid;
import java.util.Optional;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@RestController
@RequestMapping("/v1/openapi/business/")
@Slf4j
public class BusinessController {

    @Autowired
    private BusinessService businessService;
    @Autowired
    private CrossChainTxMetric crossChainTxMetric;

    @ApiOperation("调用合约（跨链）")
    @PostMapping("crossChain")
    public OpenAPIResp<TxHashResp> crossChain(@Valid @RequestBody CrossChainRequest request) {
        try {
            TxHashResp result = businessService.crossChain(request);
            crossChainTxMetric.incrementSuccess(
                    request.getDid(),
                    Optional.ofNullable(request.getOrigin()).map(CredentialSubject::getResourceName).orElse("empty"),
                    Optional.ofNullable(request.getTarget()).map(CredentialSubject::getResourceName).orElse("empty")
            );
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("crossChain error:{}", e.getMessage(), e);
            crossChainTxMetric.incrementFail(
                    request.getDid(),
                    Optional.ofNullable(request.getOrigin()).map(CredentialSubject::getResourceName).orElse("empty"),
                    Optional.ofNullable(request.getTarget()).map(CredentialSubject::getResourceName).orElse("empty")
            );
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @ApiOperation("调用合约（非跨链）")
    @PostMapping("call")
    public OpenAPIResp<TxHashResp> call(@Valid @RequestBody CrossChainRequest request) {
        try {
            TxHashResp result = businessService.call(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("get")
    public OpenAPIResp get(@Valid @RequestBody CrossChainRequest request) {
        try {
            String result = businessService.get(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @ApiOperation("链账户VC签名")
    @PostMapping("vcSignView")
    public OpenAPIResp<String> vcSignView(@Valid @RequestBody CrossChainRequest request) {
        try {
            String result = businessService.genCrossChainVcListStr(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @ApiOperation("直接访问区块链")
    @PostMapping("crossChainDirect")
    public OpenAPIResp<TxHashResp> crossChainDirect(@Valid @RequestBody VerifiableCredentialListDto vcl) {
        try {
            TxHashResp result = businessService.crossChainDirect(vcl);
            crossChainTxMetric.incrementSuccess(vcl.getOrigin().getIssuer(), "empty", "empty");
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            crossChainTxMetric.incrementFail(vcl.getOrigin().getIssuer(), "empty", "empty");
            log.error("call error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("getDidDocument")
    public OpenAPIResp getDidDocument(@Valid @RequestBody GetBusinessDidDocumentRequest request) {
        try {
            String result = businessService.getDidDocument(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("getBalance")
    public OpenAPIResp<String> getBalance(@Valid @RequestBody GetBalanceRequest request) {
        try {
            String result = businessService.getBalance(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("mint")
    public OpenAPIResp<TxHashResp> mint(@Valid @RequestBody MintRequest request) {
        try {
            TxHashResp result = businessService.mint(request);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }
}
