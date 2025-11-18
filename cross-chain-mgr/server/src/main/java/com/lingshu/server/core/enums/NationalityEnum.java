package com.lingshu.server.core.enums;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;
import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.Arrays;
import java.util.Map;
import java.util.Optional;
import java.util.function.Function;
import java.util.stream.Collectors;

/**
 * 国籍关联
 */
@Getter
@AllArgsConstructor
public enum NationalityEnum {
    // CN - 中国，US - 美国，GB - 英国，JP - 日本，RU - 俄罗斯，CA - 加拿大，AU - 澳大利亚，IN - 印度，SG - 新加坡，HK - 中国香港，TW - 中国台湾，MO - 中国澳门
    CN("CN", "中国"),
    US("US", "美国"),
    GB("GB", "英国"),
    JP("JP", "日本"),
    RU("RU", "俄罗斯"),
    CA("CA", "加拿大"),
    AU("AU", "澳大利亚"),
    IN("IN", "印度"),
    SG("SG", "新加坡"),
    HK("HK", "中国香港"),
    TW("TW", "中国台湾"),
    MO("MO", "中国澳门"),
    ;

    private final String code;
    private final String name;

    public static final Map<String, NationalityEnum> BY_CODE;

    static {
        BY_CODE = Arrays.stream(values()).collect(Collectors.toMap(NationalityEnum::getCode, Function.identity(), (a, b) -> {
            throw new RuntimeException(String.format("[%s] [%s] has same code", a, b));
        }));
    }

    @JsonValue
    public String getCode() {
        return code;
    }

    @JsonCreator(mode = JsonCreator.Mode.DELEGATING)
    public static NationalityEnum byCode(Integer code) {
        return BY_CODE.getOrDefault(code, null);
    }

    public static String codeToName(Integer code) {
        return Optional.ofNullable(BY_CODE.get(code))
                .map(NationalityEnum::getName)
                .orElse(null);
    }
}
