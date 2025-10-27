package com.lingshu.fabric.agent.resp;

import com.lingshu.fabric.agent.repo.entity.prometheus.ValuesItemDto;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

@Data
@Accessors(chain = true)
public class UnStableNodeInfosDTO {
    private String ip;
    private List<ValuesItemDto<Float>> cpuUsage;
    private List<ValuesItemDto<Long>> memUsage;
    private List<ValuesItemDto<Long>> diskUsage;
    private List<ValuesItemDto<Long>> netIn;
    private List<ValuesItemDto<Long>> netOut;
}
