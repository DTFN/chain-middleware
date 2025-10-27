package com.lingshu.fabric.agent.repo.entity;

import lombok.Data;
import lombok.experimental.Accessors;

/**
 * ShowTableDo
 *
 * @author XuHang
 * @since 2023/11/28
 **/
@Data
@Accessors(chain = true)
public class ShowTableDo {
    private String tableName;
    private String tableSchema;
}
