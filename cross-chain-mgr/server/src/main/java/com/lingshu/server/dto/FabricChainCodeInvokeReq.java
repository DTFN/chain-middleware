package com.lingshu.server.dto;

import lombok.Data;

import java.util.List;

/**
 * @author lin
 * @since 2025-09-16
 */
@Data
public class FabricChainCodeInvokeReq {
    private String chainCodeName;
    private String function;

    private String channelId;
    private String chainCodeVersion;
    private String lang;
    private boolean init;
    private List<String> args;
}
