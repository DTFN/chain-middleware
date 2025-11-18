package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;

@Data
public class AddAppChannelReq {
    @NotEmpty
    private String chainUid;
    @NotEmpty
    private String channelId;

    @NotEmpty
    private String peerName;
    @NotEmpty
    private String peerOrgName;
    @NotEmpty
    private String peerEndpoint;
    @NotEmpty
    private String ordererName;
    @NotEmpty
    private String ordererOrgName;
    @NotEmpty
    private String ordererEndpoint;

}
