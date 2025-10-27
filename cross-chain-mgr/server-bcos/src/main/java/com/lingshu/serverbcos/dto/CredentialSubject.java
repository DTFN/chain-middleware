package com.lingshu.serverbcos.dto;

import lombok.Data;

/**
 * @author: derrick
 * @since: 2025-09-08
 */
@Data
public class CredentialSubject {
    private String resourceName;
    private String contractFunc;
    private Object funcParams;
    private String paramName;
    private String key;//todo 测试
    private String did;
}