package com.lingshu.fabric.agent.service;

import com.lingshu.fabric.agent.req.NodePerformanceByTimeReq;
import com.lingshu.fabric.agent.resp.HealthzResult;
import com.lingshu.fabric.agent.resp.NodeInfoDTO;
import com.lingshu.fabric.agent.resp.UnStableNodeInfosDTO;

/**
 * NodeInfoService
 *
 * @author XuHang
 * @Date 2023/11/17 上午9:49
 **/
public interface NodeInfoService {
    HealthzResult healthyCheck(String ip, int port);
    NodeInfoDTO nodeInfo(String ip);
    UnStableNodeInfosDTO nodeInfoByTime(NodePerformanceByTimeReq red);
}
