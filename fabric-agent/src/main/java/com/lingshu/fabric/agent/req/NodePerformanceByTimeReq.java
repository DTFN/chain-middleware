package com.lingshu.fabric.agent.req;

import lombok.Data;
import lombok.experimental.Accessors;

import javax.validation.constraints.NotNull;
import java.util.Date;

/**
 * ChainOperation
 *
 * @author XuHang
 * @since 2023/11/22
 **/
@Data
@Accessors(chain = true)
public class NodePerformanceByTimeReq {
    @NotNull(message = "hostIp不可为空")
    private String hostIp;

    private Date start;

    private Date end;
}
