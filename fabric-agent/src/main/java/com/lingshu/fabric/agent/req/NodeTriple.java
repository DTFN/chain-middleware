package com.lingshu.fabric.agent.req;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import javax.validation.constraints.NotEmpty;

@Data
@NoArgsConstructor
@AllArgsConstructor
public class NodeTriple {
    @NotEmpty
    private String nodeName;
    @NotEmpty
    private String nodeOrgName;
    @NotEmpty
    private String nodeEndpoint;
}
