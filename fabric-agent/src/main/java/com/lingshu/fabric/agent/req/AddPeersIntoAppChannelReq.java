package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;

@Data
public class AddPeersIntoAppChannelReq {
    @NotEmpty
    private String chainUid;
    @NotEmpty
    private String channelId;
    // 本次需要一起加进channel的peers
    @NotNull
    private List<NodeTriple> peersNotInChannel;
    // 已经在该channel的peers
    @NotNull
    private List<NodeTriple> peersInChannel;

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
