package com.lingshu.fabric.agent.controller;

import com.alibaba.fastjson2.JSON;
import com.lingshu.fabric.agent.api.R;
import com.lingshu.fabric.agent.req.*;
import com.lingshu.fabric.agent.resp.ChainInfoDTO;
import com.lingshu.fabric.agent.resp.ChainInitResp;
import com.lingshu.fabric.agent.resp.DeployResultDTO;
import com.lingshu.fabric.agent.resp.IpAndPortsDTO;
import com.lingshu.fabric.agent.resp.base.BaseResponse;
import com.lingshu.fabric.agent.resp.base.ConstantCode;
import com.lingshu.fabric.agent.service.DeployService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;
import java.io.IOException;

@RestController
@RequestMapping("deploy")
@Slf4j
public class DeployController {

    @Autowired
    private DeployService deployService;

    @GetMapping("chainInfo")
    public R<ChainInfoDTO> getChainInfo() throws IOException {
        return R.ok(deployService.getChainInfo());
    }

    @PostMapping("initHostAndDocker")
    public R<ChainInitResp> initHostAndDocker(@RequestBody @Validated ChainInitReq req) {
        return R.ok(deployService.initHost(req));
    }

    @GetMapping("/getAvailablePort")
    public R<IpAndPortsDTO> getAvailablePort(@RequestParam String ip, @RequestParam int startPort, @RequestParam int number) {
        return R.ok(deployService.getAvailablePort(ip, startPort, number));
    }

    @PostMapping("")
    public R deploy(@RequestBody @Validated DeployReq req) {
        log.info("deploy: req={}", JSON.toJSONString(req));
        deployService.deploy(req);
        return R.ok();
    }

    @PostMapping("/delete")
    public BaseResponse delete(@RequestBody @Validated DeleteReq req) {
        deployService.delete(req);
        return new BaseResponse(ConstantCode.SUCCESS);
    }
    @PostMapping("/addNode")
    public BaseResponse addNode(@RequestBody @Validated BaseDeployReq req) {
        deployService.addNode(req);
        return new BaseResponse(ConstantCode.SUCCESS);
    }

    @GetMapping("/getDeployResult")
    public R<DeployResultDTO> getDeployResult(@RequestParam String chainUid) {
        return R.ok(deployService.getDeployResult(chainUid));
    }

    @GetMapping("/getAddNodeResult")
    public R<DeployResultDTO> getAddNodeResult(@RequestParam String chainUid, @RequestParam String requestId) {
        return R.ok(deployService.getAddNodeResult(chainUid, requestId));
    }

    @PostMapping("/addAppChannel")
    public R<String> addAppChannel(@RequestBody @Validated AddAppChannelReq req) {
        log.info("addAppChannel: req={}", req);
        deployService.addAppChannel(req);
        return R.ok();
    }

    @PostMapping("/addPeersIntoAppChannel")
    public BaseResponse addPeersIntoAppChannel(@RequestBody @Validated AddPeersIntoAppChannelReq req) throws Exception {
        BaseResponse baseResponse = new BaseResponse(ConstantCode.SUCCESS);
        deployService.addPeersIntoAppChannel(req);
        return baseResponse;
    }

    @PostMapping("/nodeOperation")
    public R chainOperation(@RequestBody @Valid NodeOperationReq req) {
        log.info("node operation: req:{}", req);
        deployService.nodeOperation(req);
        return R.ok();
    }
}
