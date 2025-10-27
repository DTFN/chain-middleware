package com.lingshu.fabric.agent.req.chaincode;

import lombok.Data;

@Data
public class FabricCCApproveReq {
    private String packageId;
    private String channelId;//通道
    private String chainCodeName;//链码名称
    private String chainCodeVersion;//链码版本
    private boolean initRequired;
}
