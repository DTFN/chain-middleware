package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

@Data
@Accessors(chain = true)
@TableName("deploy_result")
public class DeployResultDo {
    @TableId(value = "id", type = IdType.AUTO)
    private long id;
    private String chainUid;
    private String requestId;
    private String channelId;
    private Integer stage;//部署阶段
    private String error;//部署失败原因
}
