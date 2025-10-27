package com.lingshu.fabric.agent.bo;

import lombok.Data;

import java.util.Map;

/**
 * @author: zehao.song
 */
@Data
public class NetworkConfig {
    private String name;
    private String version;
    private Map<String, Object> client;
    private Map<String, Object> channels;
    private Map<String, Object> organizations;
    private Map<String, Object> orderers;
    private Map<String, Object> peers;
    private Map<String, Object> certificateAuthorities;
}
