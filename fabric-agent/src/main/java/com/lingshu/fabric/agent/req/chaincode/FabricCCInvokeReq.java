package com.lingshu.fabric.agent.req.chaincode;

import lombok.Data;

import java.util.List;

@Data
public class FabricCCInvokeReq {

    private String channelId;
    private String chainCodeName;
    private String chainCodeVersion;
    private String lang;
    private boolean init;
    private String function;
    private List<String> args;
}
