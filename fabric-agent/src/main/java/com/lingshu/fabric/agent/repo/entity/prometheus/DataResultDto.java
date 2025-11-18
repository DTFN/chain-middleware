package com.lingshu.fabric.agent.repo.entity.prometheus;

import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

/**
 * ResultDto
 *
 * @author XuHang
 * @since 2023/12/5
 **/
@Data
@Accessors(chain = true)
public class DataResultDto {
    private DataResultMetricDto metric;
    private List<String> value;
    private List<List<String>> values;
}
