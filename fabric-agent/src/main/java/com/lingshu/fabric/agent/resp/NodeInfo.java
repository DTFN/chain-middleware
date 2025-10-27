package com.lingshu.fabric.agent.resp;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * @author: zehao.song
 */
@Data
@Accessors(chain = true)
public class NodeInfo {

    private String nodeFullName;
    private String nodeIp;
    private Integer nodePort;

}
