package com.lingshu.serverbcos.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import lombok.EqualsAndHashCode;

import javax.validation.constraints.NotBlank;
import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@ApiModel(description = "Contract合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class ContractRequest implements Serializable {

//    @ApiModelProperty(value = "合约参数", required = true)
//    @NotBlank(message = "合约参数不能为空")
    private String funcParams;

    private String contractFunc;

    private String contractAddress;
}
