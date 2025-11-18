package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.IdType;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

@Data
@Accessors(chain = true)
@TableName("chain_code")
public class ChainCodeDo {
    @TableId(value = "id", type = IdType.AUTO)
    private long id;
    private String channelId;
    private String chainCodeName;
    private String lang;
}
