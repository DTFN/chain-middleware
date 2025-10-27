package com.lingshu.fabric.agent.controller;

import com.lingshu.fabric.agent.api.R;
import com.lingshu.fabric.agent.req.NodePerformanceByTimeReq;
import com.lingshu.fabric.agent.resp.HealthzResult;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;
import com.lingshu.fabric.agent.service.NodeInfoService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import javax.validation.Valid;

/**
 * NodeInfoController
 *
 * @author XuHang
 * @since 2023/11/17
 **/
@RestController
@RequestMapping("nodeInfo")
@RequiredArgsConstructor
public class NodeInfoController {
    private final NodeInfoService service;

    @GetMapping("/healthz/{ip}/{port}")
    public R<HealthzResult> healthyCheck(@PathVariable String ip, @PathVariable int port) {
        return R.ok(service.healthyCheck(ip, port));
    }

    // 基础信息
    @GetMapping("/basic")
    public R<NodeInfoDTO> nodeInfo(@RequestParam String ip) {
        return R.ok(service.nodeInfo(ip));
    }

    // 性能历史记录
    @PostMapping("/performance/byTime")
    public R<UnStableNodeInfosDTO> performanceByTime(@RequestBody @Valid NodePerformanceByTimeReq req) {
        return R.ok(service.nodeInfoByTime(req));
    }
}
