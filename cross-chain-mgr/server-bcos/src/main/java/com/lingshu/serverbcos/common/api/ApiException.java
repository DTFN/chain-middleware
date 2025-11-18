package com.lingshu.serverbcos.common.api;


import com.lingshu.serverbcos.common.api.ApiErrorCode;

/**
 * REST API 请求异常类
 *
 * @author Brian
 * @since 2021-03-13
 */
public class ApiException extends RuntimeException {

    private static final long serialVersionUID = -5885155226898287919L;

    /**
     * 错误码
     */
    private ApiErrorCode errorCode;

    public ApiException(ApiErrorCode errorCode) {
        super(errorCode.getMsg());
        this.errorCode = errorCode;
    }

    public ApiException(String message) {
        super(message);
    }

    public ApiException(Throwable cause) {
        super(cause);
    }

    public ApiException(String message, Throwable cause) {
        super(message, cause);
    }

    public ApiErrorCode getErrorCode() {
        return errorCode;
    }
}
