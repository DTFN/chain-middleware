package com.lingshu.fabric.agent.req;

import lombok.Data;
import lombok.experimental.Accessors;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;

/**
 * 链的停止，启动，删除
 *
 * @author XuHang
 * @since 2023/11/22
 **/
@Data
@Accessors(chain = true)
public class NodeOperationReq {
    @NotNull
    private Operation operation;

    @NotNull
    private String nodeName;

    // 下个版本删除
    @NotEmpty
    private String hostIp;

    public enum Operation{
        START,
        STOP,
        RESTART,
        ;
    }
}
