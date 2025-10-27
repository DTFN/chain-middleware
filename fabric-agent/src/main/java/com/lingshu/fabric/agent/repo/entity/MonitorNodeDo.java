package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Date;

@Data
@Accessors(chain = true)
@TableName("monitor_node")
public class MonitorNodeDo {
    @TableId(value = "id", type = IdType.AUTO)
    private Long id;
    private String ip;
    private float cpuUsage;
    private Long memUsage;
    private Long diskUsage;
    private Long netIn;
    private Long netOut;
    private Date createTime;
}
