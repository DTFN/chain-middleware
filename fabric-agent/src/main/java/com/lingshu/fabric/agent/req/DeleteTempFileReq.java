package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;

@Data
public class DeleteTempFileReq {
    @NotEmpty
    private String chainCodeName;
}
