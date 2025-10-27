package com.lingshu.server.dto;

import lombok.Data;

import javax.validation.constraints.NotNull;

/**
 * @author gongrui.wang
 * @since 2025/9/19
 */
@Data
public class BackupAccountRequest {
    @NotNull
    private String did;
    @NotNull
    private String mnemonic;
    @NotNull
    private String password;
}
