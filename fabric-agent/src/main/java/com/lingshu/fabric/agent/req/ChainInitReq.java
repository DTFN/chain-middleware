package com.lingshu.fabric.agent.req;

import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

/**
 * DeployChainReq
 *
 * @author XuHang
 * @since 2023/11/13
 **/
@Data
@Accessors(chain = true)
public class ChainInitReq {
    private String chainNameCn;

    private String detail;

    private String chainVersion;

    private String chainName;

    private String agency;

    private String agentUrl;

    private List<Node> orderers;

    private List<Node> peers;

    @Data
    @Accessors(chain = true)
    public static class Node {
        private String ip;

        private int number;

        private String nodeName;

        private List<Integer> ports;
    }
}
