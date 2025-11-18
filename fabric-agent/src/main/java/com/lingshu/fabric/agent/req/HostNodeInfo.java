package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;
import java.util.Map;

/**
 * @author: zehao.song
 */
@Data
public class HostNodeInfo {

    @NotEmpty
    private String ip;
    private Map<String, Integer> peerNameAndPorts; // peer0.org1.example.com,7051
    private Map<String, Integer> ordererNameAndPorts; // orderer0.org1.example.com,7050
}
