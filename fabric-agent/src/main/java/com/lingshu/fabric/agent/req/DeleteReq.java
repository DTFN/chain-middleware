package com.lingshu.fabric.agent.req;

import lombok.Data;

import javax.validation.constraints.NotEmpty;
import javax.validation.constraints.NotNull;
import java.util.List;
import java.util.Map;

@Data
public class DeleteReq {
    @NotEmpty
    private String chainUid;
}
