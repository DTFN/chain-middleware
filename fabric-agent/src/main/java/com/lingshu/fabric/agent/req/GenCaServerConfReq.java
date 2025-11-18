package com.lingshu.fabric.agent.req;

import lombok.Data;

/**
 * @author lin
 * @since 2023-11-16
 */
@Data
public class GenCaServerConfReq {
    private String caName;
    private String caIp;
    private Integer caPort;
    private String caAdmin;
    private String caPw;
    private String caDomain;
    private String orgDomain;
}
