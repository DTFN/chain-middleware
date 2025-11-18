package com.lingshu.server.dto;

import com.lingshu.server.core.enums.NationalityEnum;
import lombok.Data;

import javax.validation.constraints.NotNull;

/**
 * @author gongrui.wang
 * @since 2025/9/19
 */
@Data
public class CreateAccountRequest {
    @NotNull
    private String password;
    private NationalityEnum nationality = NationalityEnum.CN;
}
