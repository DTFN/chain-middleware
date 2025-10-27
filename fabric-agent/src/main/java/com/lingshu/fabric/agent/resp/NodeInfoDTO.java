package com.lingshu.fabric.agent.resp;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * NodeInfoDTO
 *
 * @author XuHang
 * @since 2023/11/21
 **/
@Data
@Accessors(chain = true)
public class NodeInfoDTO {
    private Integer cpuNum = 1;
    private Float cpuRatio = 0f;
    private Long maxMem = 0L;
    private Long usedMem = 0L;
    private Long maxDisk = 0L;
    private Long usedDisk = 0L;
    private Long netIn = 0L;
    private Long netOut = 0L;
}
