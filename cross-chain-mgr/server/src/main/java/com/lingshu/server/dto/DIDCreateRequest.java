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
 * @since: 2025-08-22
 */
@ApiModel(description = "did创建请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class DIDCreateRequest implements Serializable {
    @ApiModelProperty(value = "DID ID", required = true)
    @NotBlank(message = "DID ID不能为空")
//    @Length(max = 64, message = "DID ID长度不能超过64")
    private String did;

    @ApiModelProperty(value = "DID内容", required = true)
    @NotBlank(message = "DID内容不能为空")
//    @Length(max = 1024, message = "DID内容长度不能超过1024")
    private String didDocument;
}
