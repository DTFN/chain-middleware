package com.lingshu.server.dto;

import cn.hutool.core.annotation.Alias;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;
import java.util.Map;

@Data
@Accessors(chain = true)
public class VerifiableCredentialDto {
    @JsonProperty("@context")
    @Alias("@context")
    private List<String> context;

    private CredentialSubject credentialSubject;

    private String id;

    private String issuanceDate;

    private String issuer;

    private Object proof;

    private List<String> type;

    @Data
    @Accessors(chain = true)
    public static class CredentialSubject {
        @JsonProperty("chain_rid")
        @Alias("chain_rid")
        private String chainRid;

        @JsonProperty("contract_address")
        @Alias("contract_address")
        private String contractAddress;

        @JsonProperty("contract_func")
        @Alias("contract_func")
        private String contractFunc;

        @JsonProperty("contract_type")
        @Alias("contract_type")
        private String contractType;

        @JsonProperty("func_params")
        @Alias("func_params")
        private Object funcParams;

        @JsonProperty("gateway_id")
        @Alias("gateway_id")
        private String gatewayId;

        @JsonProperty("param_name")
        @Alias("param_name")
        private String paramName;

        @JsonProperty("resource_name")
        @Alias("resource_name")
        private String resourceName;

        @JsonProperty("key")
        @Alias("key")
        private String key;
    }

    @Data
    @Accessors(chain = true)
    public static class ProofEd25519 {
        private String created;

        private String proofPurpose;

        private String signature;

        private String type;

        private String verificationMethod;
    }

    @Data
    @Accessors(chain = true)
    public static class ProofEth {
        private String address;

        private String contentHash;

        private String r;

        private String s;

        private Integer v;

        private String verificationMethod;
    }
}
