package com.lingshu.server.common.api;

import org.apache.commons.lang3.StringUtils;
import org.apache.poi.ss.formula.functions.T;

/**
 * @Author wang jian
 * @Date 2022/4/25 15:53
 */
public class OpenAPIRespBuilder {

    public static <T> OpenAPIResp<T> success() {
        return build(ApiErrorCode.SUCCESS);
    }

    public static <T> OpenAPIResp<T> success(T data) {

        OpenAPIResp<T> resp = success();
        resp.setData(data);

        return resp;
    }

    public static OpenAPIResp success(ApiErrorCode respCode) {
        return build(respCode);
    }

    public static <T> OpenAPIResp<T> failure(ApiErrorCode respCode) {
        return build(respCode);
    }

    public static <T> OpenAPIResp<T> failure(ApiErrorCode respCode, String additionalErrMsg) {
        OpenAPIResp<T> resp = failure(respCode);
        setAdditionalErrMsg(resp, additionalErrMsg);
        return resp;
    }


    private static OpenAPIResp build(ApiErrorCode apiErrorCode) {
        return OpenAPIResp.builder()
            .code(apiErrorCode.getCode())
            .msg(apiErrorCode.getMsg())
            .build();
    }

    private static void setAdditionalErrMsg(OpenAPIResp resp, String additionalErrMsg) {
        if (StringUtils.isNotEmpty(additionalErrMsg)) {
            resp.setMsg(resp.getMsg() + ", " + additionalErrMsg);
        }
    }

    public static <T> OpenAPIResp failure(ApiErrorCode respCode, T data) {
        OpenAPIResp openAPIResp = build(respCode);
        openAPIResp.setData(data);
        return openAPIResp;
    }
}
