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
@ApiModel(description = "BusinessContract创建请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class BusinessContractSaveRequest implements Serializable {
    @ApiModelProperty(value = "BusinessContract businessAddress", required = true)
    @NotBlank(message = "BusinessContract businessAddress不能为空")
    @Length(max = 64, message = "BusinessContract businessAddress长度不能超过64")
    private String businessAddress;

    @ApiModelProperty(value = "ResourceDomain内容", required = true)
    @NotBlank(message = "ResourceDomain内容不能为空")
    @Length(max = 1024, message = "ResourceDomain内容长度不能超过1024")
    private String details;
}
