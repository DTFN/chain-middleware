package com.lingshu.server.dto;

import cn.hutool.core.annotation.Alias;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;
import lombok.experimental.Accessors;

import java.util.List;

@Data
@Accessors(chain = true)
public class DidDto {
    @JsonProperty("@context")
    @Alias("@context")
    private List<Object> context;  // 可以包含字符串和对象

    private String id;

    private List<Object> verificationMethod;

    private List<String> assertionMethod;

    // alsoknownas用于存储用户在其他系统中的别名，可以用来存储用户的国籍
    // https://www.w3.org/TR/did-extensions-properties/#alsoknownas
//    private String alsoknownas;

//    private List<String> authentication;

//    private List<String> capabilityDelegation;

//    private List<String> capabilityInvocation;

    @Data
    @Accessors(chain = true)
    public static class VerificationMethodEd25519 {
        private String id;

        private String type;

        private String controller;

        @JsonProperty("publicKeyHex")
        @Alias("publicKeyHex")
        private String publicKeyHex;
    }

    @Data
    @Accessors(chain = true)
    public static class VerificationMethodEth {
        private String id;

        private String type;

        private String controller;

        private String address;  // 用于EcdsaSecp256k1VerificationKey2019类型
    }

    @Data
    @Accessors(chain = true)
    public static class ContextObject {
        @JsonProperty("Ed25519VerificationKey2018")
        @Alias("Ed25519VerificationKey2018")
        private String ed25519VerificationKey2018;

        @JsonProperty("publicKeyBase58")
        @Alias("publicKeyBase58")
        private String publicKeyBase58;

        @JsonProperty("publicKeyJwk")
        @Alias("publicKeyJwk")
        private PublicKeyJwk publicKeyJwk;
    }

    @Data
    @Accessors(chain = true)
    public static class PublicKeyJwk {
        @JsonProperty("@id")
        @Alias("@id")
        private String id;

        @JsonProperty("@type")
        @Alias("@type")
        private String type;
    }
}
