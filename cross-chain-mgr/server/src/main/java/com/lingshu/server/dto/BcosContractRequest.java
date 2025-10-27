package com.lingshu.server.dto;

import io.swagger.annotations.ApiModel;
import lombok.Data;
import lombok.EqualsAndHashCode;

import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@ApiModel(description = "Contract合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class BcosContractRequest implements Serializable {

//    @ApiModelProperty(value = "合约参数", required = true)
//    @NotBlank(message = "合约参数不能为空")
//    private String param;
//
//    private String contractAddress;

    private String funcParams;

    private String contractFunc;

    private String contractAddress;

//    private String binary;
}
