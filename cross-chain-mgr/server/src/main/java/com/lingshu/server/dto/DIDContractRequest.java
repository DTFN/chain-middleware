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
@ApiModel(description = "DID合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class DIDContractRequest implements Serializable {

    @ApiModelProperty(value = "DID合约名称", required = true)
    @NotBlank(message = "DID合约名称不能为空")
    @Length(max = 100, message = "DID合约名称不能超过1024")
    private String contractName;
}
