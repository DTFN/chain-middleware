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
@Data
@EqualsAndHashCode(callSuper = false)
public class CrossChainRequest implements Serializable {

    private String privateKeyHex;
    private CredentialSubject origin;
    private CredentialSubject target;

    private ResourceDomainInfo originResourceDomainInfo;
}