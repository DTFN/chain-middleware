package com.lingshu.fabric.agent.service;

import com.lingshu.fabric.agent.req.*;
import com.lingshu.fabric.agent.resp.ChainInfoDTO;
import com.lingshu.fabric.agent.resp.ChainInitResp;
import com.lingshu.fabric.agent.resp.DeployResultDTO;
import com.lingshu.fabric.agent.resp.IpAndPortsDTO;
import org.springframework.web.bind.annotation.RequestParam;

import java.io.IOException;

public interface DeployService {
    ChainInfoDTO getChainInfo() throws IOException;

    ChainInitResp initHost(ChainInitReq req);
    IpAndPortsDTO getAvailablePort(@RequestParam String ip, @RequestParam int startPort, @RequestParam int number);

    void deploy(DeployReq req);

    void delete(DeleteReq req);

    void addNode(BaseDeployReq req);

    DeployResultDTO getDeployResult(String chainUid);
    DeployResultDTO getAddNodeResult(String chainUid, String requestId);

    void addAppChannel(AddAppChannelReq req);

    void addPeersIntoAppChannel(AddPeersIntoAppChannelReq req) throws Exception;

    void nodeOperation(NodeOperationReq req);
}
