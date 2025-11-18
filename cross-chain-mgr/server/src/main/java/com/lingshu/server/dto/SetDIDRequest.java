package com.lingshu.server.dto;

import com.lingshu.server.core.enums.ContractTypeEnum;
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
@Data
@EqualsAndHashCode(callSuper = false)
public class SetDIDRequest implements Serializable {
    private ContractTypeEnum contractType;
    private String busiContractAddress;
    private String didContractAddress;
}
