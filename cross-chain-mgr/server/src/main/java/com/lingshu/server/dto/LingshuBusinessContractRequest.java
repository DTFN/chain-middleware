package com.lingshu.server.dto;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import lombok.EqualsAndHashCode;
import org.chainmaker.pb.common.ContractOuterClass;
import org.hibernate.validator.constraints.Length;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@ApiModel(description = "LingshuBusinessContract合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class LingshuBusinessContractRequest implements Serializable {

    @ApiModelProperty(value = "LingshuBusinessContract合约地址", required = true)
    @NotBlank(message = "LingshuBusinessContract合约地址不能为空")
    @Length(max = 100, message = "LingshuBusinessContract合约地址不能超过1024")
    private String contractAddress;
}
