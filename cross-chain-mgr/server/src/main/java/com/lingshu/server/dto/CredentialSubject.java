package com.lingshu.server.dto;

import lombok.Data;

import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-09-08
 */
@Data
public class CredentialSubject {
    private String resourceName;
    private String contractFunc;
    private Object funcParams;
    private String paramName = "vcs";
    private String key;//todo 测试
}
