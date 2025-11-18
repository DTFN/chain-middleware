package com.lingshu.fabric.agent.resp;

import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Date;

@Data
@Accessors(chain = true)
public class HealthzResult {
    private String status;
    private Date time;
}