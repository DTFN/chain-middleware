package com.lingshu.serverbcos.dto;

import com.lingshu.serverbcos.core.enums.ContractTypeEnum;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import lombok.Data;
import lombok.EqualsAndHashCode;
import org.hibernate.validator.constraints.Length;

import javax.validation.constraints.NotBlank;
import javax.validation.constraints.NotNull;
import java.io.Serializable;

/**
 * @author: derrick
 * @since: 2025-08-22
 */
@ApiModel(description = "BusinessContract合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class BusinessContractRequest implements Serializable {
    private ContractTypeEnum contractType;

    private String contractName;

    private String version;

    private String filePath;

}
