package com.lingshu.fabric.agent.resp;

import lombok.Data;

import java.util.List;

/**
 * IpAndPortsDTO
 *
 * @author XuHang
 * @since 2023/11/17
 **/
@Data
public class IpAndPortsDTO {
    private String ip;
    private List<Integer> ports;
}
