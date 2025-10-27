package com.lingshu.server.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import lombok.EqualsAndHashCode;
import org.hibernate.validator.constraints.Length;

import javax.validation.constraints.NotBlank;
import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
@ApiModel(description = "BusinessContract4WasmCall请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class BusinessContract4WasmCallRequest implements Serializable {
    @ApiModelProperty(value = "BusinessContract4Wasm methodName", required = true)
    @NotBlank(message = "BusinessContract4Wasm methodName不能为空")
    private String methodName;

    @ApiModelProperty(value = "BusinessContract4Wasm paramName", required = true)
    @NotBlank(message = "BusinessContract4Wasm paramName不能为空")
    private String paramName;

    @ApiModelProperty(value = "BusinessContract4Wasm paramValue", required = true)
    @NotBlank(message = "BusinessContract4Wasm paramValue不能为空")
    private String paramValue;

}
