package com.lingshu.server.controller;

import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import com.lingshu.server.common.metrics.CreateAccountMetric;
import com.lingshu.server.core.business.AccountService;
import com.lingshu.server.dto.BackupAccountRequest;
import com.lingshu.server.dto.CreateAccountRequest;
import com.lingshu.server.dto.RestoreAccountRequest;
import com.lingshu.server.dto.UpdateAccountRequest;
import com.lingshu.server.dto.resp.account.CreateAccountResp;
import com.lingshu.server.dto.resp.account.DidDocumentResp;
import com.lingshu.server.key.PkeyByMnemonicService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.annotation.Resource;
import javax.validation.Valid;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@RestController
@RequestMapping("/v1/openapi/account/")
@Slf4j
public class AccountController {

    @Autowired
    private AccountService accountService;

    @Resource
    private PkeyByMnemonicService pkeyByMnemonicService;

    @Autowired
    private CreateAccountMetric createAccountMetric;

    @PostMapping("create")
    public OpenAPIResp<CreateAccountResp> create(@Valid @RequestBody CreateAccountRequest request) {
        try {
            CreateAccountResp createAccountResp = accountService.create(request.getPassword(), request.getNationality());
            createAccountMetric.incrementSuccess("create");
            return OpenAPIRespBuilder.success(createAccountResp);
        } catch (Exception e) {
            log.error("create account error:{}", e.getMessage(), e);
            createAccountMetric.incrementFail("create");
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("update")
    public OpenAPIResp<Boolean> update(@Valid @RequestBody UpdateAccountRequest request) {
        try {
            Boolean result = accountService.update(request);
            createAccountMetric.incrementSuccess("update");
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("update error:{}", e.getMessage(), e);
            createAccountMetric.incrementFail("update");
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    @PostMapping("getDidDocument")
    public OpenAPIResp<DidDocumentResp> getDidDocument(@RequestParam String did) {
        try {
            DidDocumentResp result = accountService.getDidDocument(did);
            return OpenAPIRespBuilder.success(result);
        } catch (Exception e) {
            log.error("get error:{}", e.getMessage(), e);
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    /**
     * 密钥导出
     * @param request
     * @return
     */
    @PostMapping("backup")
    public OpenAPIResp<String> backupAccount(@Valid @RequestBody BackupAccountRequest request) {
        try {
            String encrypt = accountService.backupAccount(request);
            return OpenAPIRespBuilder.success(encrypt);
        } catch (Exception e) {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }

    /**
     * 密钥导入
     * @param request
     * @return
     */
    @PostMapping("restore")
    public OpenAPIResp<CreateAccountResp> restoreAccount(@Valid @RequestBody RestoreAccountRequest request) {
        try {
            CreateAccountResp createAccountResp = accountService.importAccount(request);
            createAccountMetric.incrementSuccess("restore");
            return OpenAPIRespBuilder.success(createAccountResp);
        } catch (Exception e) {
            createAccountMetric.incrementFail("restore");
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, e.getMessage());
        }
    }
}
