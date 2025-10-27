package com.lingshu.fabric.agent.resp;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * 网络速度，单位 byte/s
 *
 * @author XuHang
 * @since 2023/11/21
 **/
@Data
@Accessors(chain = true)
public class NetInfoDTO {
    private String ip;
    private Long speedIn = 0L;
    private Long speedOut = 0L;
}
