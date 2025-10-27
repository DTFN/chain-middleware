package com.lingshu.fabric.agent.repo.entity.prometheus;

import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Date;

/**
 * ValuesItemDto
 *
 * @author XuHang
 * @since 2023/12/5
 **/
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ValuesItemDto<T> {
    private T value;
    private Date time;
}
