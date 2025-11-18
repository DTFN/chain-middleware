package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Date;

@Data
@Accessors(chain = true)
@TableName("chain_info")
public class ChainInfoDo {
    @TableId(value = "id", type = IdType.AUTO)
    private long id;
    private String chainUid;
    private String caHost;
    private String caName;
    private Integer caPort;
    private Date modifyTime;
    private Date createTime;
}
