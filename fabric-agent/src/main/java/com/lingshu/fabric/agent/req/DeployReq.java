package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;

@Data
public class DeployReq extends BaseDeployReq {

    @NotEmpty
    private String channelId;

}
