package com.lingshu.fabric.agent.req;

import lombok.Data;

/**
 * @author lin
 * @since 2023-12-12
 */
@Data
public class GenCaServerReq {
    private String env;
    private String network;
    private String caName;
    private String containerName;
    private String caAdmin;
    private String caAdminPw;
    private String caDir;
    private String caHost;
    private Integer caPort;
    private Integer caListenPort;
}
