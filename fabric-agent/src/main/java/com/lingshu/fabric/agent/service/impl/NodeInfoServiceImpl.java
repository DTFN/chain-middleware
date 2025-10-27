package com.lingshu.fabric.agent.service.impl;

import com.lingshu.fabric.agent.req.NodePerformanceByTimeReq;
import com.lingshu.fabric.agent.resp.HealthzResult;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;
import com.lingshu.fabric.agent.service.MonitorService;
import com.lingshu.fabric.agent.service.NodeInfoService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;

/**
 * NodeInfoServiceImpl
 *
 * @author XuHang
 * @since 2023/11/17
 **/
@Slf4j
@Service
@RequiredArgsConstructor
public class NodeInfoServiceImpl implements NodeInfoService {
    private final RestTemplate healthzRestTemplate;
    private final MonitorService monitorService;

    @Override
    public HealthzResult healthyCheck(String ip, int port) {
        String url = String.format("http://%s:%s/healthz", ip, port);
        try {
            ResponseEntity<HealthzResult> entity = healthzRestTemplate.getForEntity(url, HealthzResult.class);
            return entity.getBody();
        } catch (Exception e) {
            log.warn("get healthz error, ip:{}, port:{}", ip, port);
            return null;
        }
    }

    @Override
    public NodeInfoDTO nodeInfo(String ip) {
        // 查询历史数据
        return monitorService.hostLatestInfo(ip);
    }

    @Override
    public UnStableNodeInfosDTO nodeInfoByTime(NodePerformanceByTimeReq req) {
        UnStableNodeInfosDTO r = monitorService.hostInfoHistory(req.getHostIp(), req.getStart(), req.getEnd());
        return r;
    }
}
