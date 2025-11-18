package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Date;

@Data
@Accessors(chain = true)
@TableName("monitor_node_stable")
public class MonitorNodeStableDo {
    @TableId(value = "ip")
    private String ip;
    private Integer cpuNumber;
    private Long maxMemory;
    private Long maxDisk;
    private Date modifyTime;
}
