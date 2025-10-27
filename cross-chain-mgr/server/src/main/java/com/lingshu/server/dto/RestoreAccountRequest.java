package com.lingshu.server.dto;

import com.lingshu.server.core.enums.NationalityEnum;
import lombok.Data;

/**
 * @author gongrui.wang
 * @since 2025/9/19
 */
@Data
public class RestoreAccountRequest {
    // 加密后的私钥
    private String encrypt;
    private String password;

    // 国籍
    private NationalityEnum nationality = NationalityEnum.CN;
}
