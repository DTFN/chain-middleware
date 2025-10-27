package com.lingshu.serverbcos.common.api;

import io.swagger.annotations.ApiParam;
import lombok.Data;

/**
 * @Author wang jian
 * @Date 2022/4/25 15:46
 */
@Data
public class OpenApiGeneralReq {

    private String sign;
    private String appId;
    private Integer timeStamp;
    @ApiParam(hidden = true)
    private Integer shopId;
}
