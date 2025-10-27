package com.lingshu.fabric.agent.bo;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * @author: zehao.song
 */
@Data
@Accessors(chain = true)
public class NodeDetail {
    private boolean isPeer;
    private String nodeFullName;
    private String networkName;
    private Integer nodePort;
    private String nodeIp;
    private String orgFullName;
    private String mspId;
    private Integer chainCodePort;
    private Integer operationsPort;
}
