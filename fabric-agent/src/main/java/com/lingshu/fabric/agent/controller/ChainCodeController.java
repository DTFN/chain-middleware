package com.lingshu.fabric.agent.controller;

import com.lingshu.fabric.agent.api.R;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.req.chaincode.*;
import com.lingshu.fabric.agent.resp.FabricChainCodeInvokeResp;
import com.lingshu.fabric.agent.resp.PackageInfoDTO;
import com.lingshu.fabric.agent.service.ChainCodeService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import java.time.Duration;
import java.time.Instant;

@RestController
@RequestMapping("code")
@Slf4j
public class ChainCodeController {

    @Autowired
    private ChainCodeService chainCodeService;

    @PostMapping("/upload")
    public R<PackageInfoDTO> upload(@RequestPart("file") MultipartFile chainCodeFile){
        log.info("chain code upload");
        return R.ok(chainCodeService.resolveChaincodePackage(chainCodeFile));
    }

    @PostMapping("/install")
    public R<String> installChainCode(@RequestBody FabricCCInstallReq req){
        Instant startTime = Instant.now();
        log.info("chain code install: req={}", req);
        try {
            String packageId = chainCodeService.installChainCode(req);
            log.info("chain code install cost:{}ms", Duration.between(startTime, Instant.now()).toMillis());
            return R.ok(packageId);
        } catch (Exception e){
            log.error("installChainCode install, req={}", req, e);
            return R.failed(ConstantCode.CHAIN_CODE_INSTALL_FAILED);
        }
    }

    @PostMapping("/approve")
    public R<Long> approveChainCode(@RequestBody FabricCCApproveReq req){
        Instant startTime = Instant.now();
        log.info("chain code approve: req={}", req);
        try {
            long sequence = chainCodeService.approveChainCode(req);
            log.info("chain code approve cost:{}ms", Duration.between(startTime, Instant.now()).toMillis());
            return R.ok(sequence);
        } catch (Exception e){
            log.error("approveChainCode failed, req={}", req, e);
            return R.failed(ConstantCode.CHAIN_CODE_APPROVE_FAILED);
        }
    }

    @PostMapping("/commit")
    public R<String> commitChainCode(@RequestBody FabricCCCommitReq req){
        Instant startTime = Instant.now();
        log.info("chain code commit: req={}", req);
        try {
            chainCodeService.commitChainCode(req);
            log.info("chain code commit cost:{}ms", Duration.between(startTime, Instant.now()).toMillis());
            return R.ok();
        } catch (Exception e){
            log.error("commitChainCode failed, req={}", req, e);
            return R.failed(ConstantCode.CHAIN_CODE_COMMIT_FAILED);
        }
    }

    @PostMapping("/invoke")
    public R<Object> invokeChainCode(@RequestBody FabricCCInvokeReq req){
        Instant startTime = Instant.now();
        log.info("chain code invoke: req={}", req);
        try {
            FabricChainCodeInvokeResp resp = chainCodeService.invokeChainCode(req);
            log.info("chain code invoke cost:{}ms resp={}", Duration.between(startTime, Instant.now()).toMillis(), resp);
            return R.ok(resp);
        } catch (Exception e){
            log.error("invokeChainCode failed, req={}", req, e);
            return R.failed(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
        }
    }

    @PostMapping("/query")
    public R<Object> queryChainCode(@RequestBody FabricCCInvokeReq req){
        log.info("chain code query: req={}", req);
        try {
            return R.ok(chainCodeService.queryChainCode(req));
        } catch (Exception e){
            log.error("queryChainCode failed, req={}", req, e);
            return R.failed(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
        }
    }
}
