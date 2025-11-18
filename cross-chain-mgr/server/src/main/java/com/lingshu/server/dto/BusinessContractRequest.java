package com.lingshu.server.dto;

import com.lingshu.server.core.enums.ContractTypeEnum;
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
@ApiModel(description = "BusinessContract合约请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class BusinessContractRequest implements Serializable {

    @ApiModelProperty(value = "BusinessContract合约类型", required = true)
    @NotNull(message = "BusinessContract合约类型不能为空")
    private ContractTypeEnum contractType;

    @ApiModelProperty(value = "BusinessContract合约名称", required = true)
    @NotBlank(message = "BusinessContract合约名称不能为空")
    @Length(max = 100, message = "BusinessContract合约名称不能超过1024")
    private String contractName;

    @ApiModelProperty(value = "BusinessContract合约version", required = true)
    private String version;

    @ApiModelProperty(value = "BusinessContract合约文件路径", required = true)
    @NotBlank(message = "BusinessContract合约文件路径不能为空")
//    @Length(max = 100, message = "BusinessContract合约文件路径不能超过1024")
    private String filePath;

//    @ApiModelProperty(value = "BusinessContract合约类型", required = true)
////    @NotNull(message = "BusinessContract合约类型不能为空")
//    private ContractOuterClass.RuntimeType runtimeType;
}
