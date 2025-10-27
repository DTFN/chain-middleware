package com.lingshu.server.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

@Data
@Accessors(chain = true)
public class VerifiableCredentialListDto {
    private VerifiableCredentialDto origin;
    private VerifiableCredentialDto target;
}
