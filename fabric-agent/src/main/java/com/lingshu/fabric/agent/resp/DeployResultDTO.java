package com.lingshu.fabric.agent.resp;

import com.lingshu.fabric.agent.bo.NodeDetail;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

@Data
@Accessors(chain = true)
public class DeployResultDTO {
    private String chainUid;
    private String requestId;
    private String channelId;
    private List<NodeDetail> hostNodes;
    private Integer stage;//部署阶段
    private String error;//部署失败原因
}
