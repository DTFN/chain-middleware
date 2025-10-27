package com.lingshu.fabric.agent.req;

import lombok.Data;

import java.util.List;

/**
 * @author lin
 * @since 2023-11-16
 */
@Data
public class GenConfigTxReq {
    private String chainPath;
    private Org orderOrg;
    private Org peerOrg;
    private List<Node> ordererList;
    private List<Node> peerList;
    private List<Channel> channelList;


    @Data
    public static class Node {
        private String ip;
        private String port;
        private String domain;
        private String orgDomain;
    }

    @Data
    public static class Org {
        private String orgName;
        private String mspId;
        private String domain;
    }

    @Data
    public static class Channel {
        private String channelName;
    }
}
