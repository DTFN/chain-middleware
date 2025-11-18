package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

/**
 * @author: zehao.song
 */
@Data
@Accessors(chain = true)
@TableName("deploy_result_node")
public class DeployResultNodeDo {
    @TableId(value = "id", type = IdType.AUTO)
    private long id;
    private long resultId;
    private String nodeFullName;
    private Integer nodeType;
    private String networkName;
    private Integer nodePort;
    private String nodeIp;
    private String orgFullName;
    private String mspId;
    private Integer chainCodePort;
    private Integer operationsPort;
}
