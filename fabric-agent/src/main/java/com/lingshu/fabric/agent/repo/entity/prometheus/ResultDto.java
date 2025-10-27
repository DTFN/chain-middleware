package com.lingshu.fabric.agent.repo.entity.prometheus;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * ResultDto
 *
 * @author XuHang
 * @since 2023/12/5
 **/
@Data
@Accessors(chain = true)
public class ResultDto {
    private String status;
    private ResultDataDto data;
}
