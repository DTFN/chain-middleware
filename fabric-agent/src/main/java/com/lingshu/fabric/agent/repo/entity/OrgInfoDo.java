package com.lingshu.fabric.agent.repo.entity;

import com.baomidou.mybatisplus.annotation.TableName;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Date;

@Data
@Accessors(chain = true)
@TableName("org_info")
public class OrgInfoDo {
    private String chainUid;
    private String orgName;
    private String adminKey;
    private String adminCrt;
    private String caCrt;
    private Date modifyTime;
    private Date createTime;
}
