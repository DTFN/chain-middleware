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
@ApiModel(description = "ResourceDomain创建请求")
@Data
@EqualsAndHashCode(callSuper = false)
public class ResourceDomainSaveRequest implements Serializable {
    private String resourceName;
    private String gatewayId;
    private String chainRid;
    private String contractType;
    private String contractAddress;
}
