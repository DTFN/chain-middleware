package com.lingshu.fabric.agent.resp;

import lombok.Data;

import java.io.Serializable;
import java.util.List;
import java.util.Map;

@Data
public class ChainInfoDTO implements Serializable {

    // 共识类型（etcdraft、solo...）
    private String consensusType;
    // 链版本
    private String version;
    // 链名称
    private String chainName;

    // connection-tls.json中配置的通道节点信息
    private Map<String, List<NodeInfo>> channelNodeMap;

}
