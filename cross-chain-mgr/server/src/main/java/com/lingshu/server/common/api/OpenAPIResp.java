package com.lingshu.server.common.api;

import io.swagger.annotations.ApiModelProperty;
import lombok.*;

import java.io.Serializable;

/**
 * @Author wang jian
 * @Date 2022/4/25 15:51
 */
@Setter
@Getter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class OpenAPIResp<T> implements Serializable {

    @ApiModelProperty(value = "响应数据")
    private T data;

    @ApiModelProperty(value = "响应消息")
    private String msg;

    @ApiModelProperty(value = "响应码")
    private Integer code;
}