package com.lingshu.server.common.api;


import lombok.Data;
import lombok.experimental.Accessors;
import org.apache.commons.lang3.StringUtils;

import java.io.Serializable;
import java.util.Optional;

/**
 * REST API 返回结果
 *
 * @author Brian
 * @since 2021-03-13
 */
@Data
@Accessors(chain = true)
public class R<T> implements Serializable {

    /**
     * serialVersionUID
     */
    private static final long serialVersionUID = 1L;

    /**
     * 业务错误码
     */
    private Integer code;
    /**
     * 结果集
     */
    private T data;
    /**
     * 描述
     */
    private String msg;

    /**
     * 描述信息类型
     * 0：普通     用户操作成功
     * 1：提示错误  如用户输入的数据有误，没有权限操作一类
     * 2：失败错误  操作失败，系统错误
     */
    private Integer type = 0;

    public R() {
        // to do nothing
    }

    public R(RetCode errorCode) {
        errorCode = Optional.ofNullable(errorCode).orElse(ConstantCode.FAILED);
        this.code = errorCode.getCode();
        this.msg = errorCode.getMessage();
    }

    public static <T> R<T> ok(T data) {
        return restResult(data, ConstantCode.SUCCESS);
    }

    public static <T> R<T> ok(Integer code, String msg, T data) {
        return restResult(data, code, msg);
    }

    public static <T> R<T> ok() {
        return restResult(null, ConstantCode.SUCCESS.getCode(), ConstantCode.SUCCESS.getMessage());
    }

    public static <T> R<T> failed(String msg) {
        return restResult(null, ConstantCode.FAILED.getCode(), StringUtils.isEmpty(msg) ? ConstantCode.FAILED.getMessage() : msg);
    }

    public static <T> R<T> failed(Integer code, String msg) {
        return restResult(null, code, StringUtils.isEmpty(msg) ? ConstantCode.FAILED.getMessage() : msg);
    }

    public static <T> R<T> failed(RetCode errorCode) {
        return restResult(null, errorCode.getCode(), errorCode.getMessage());
    }

    public static <T> R<T> restResult(T data, RetCode errorCode) {
        return restResult(data, errorCode.getCode(), errorCode.getMessage());
    }

    private static <T> R<T> restResult(T data, Integer code, String msg) {
        R<T> apiResult = new R<>();
        apiResult.setCode(code);
        apiResult.setData(data);
        apiResult.setMsg(msg);
        return apiResult;
    }

    private static <T> R<T> restResult(T data, Integer code, String msg, Integer msgType) {
        R<T> apiResult = new R<>();
        apiResult.setCode(code);
        apiResult.setData(data);
        apiResult.setMsg(msg);
        apiResult.setType(msgType);
        return apiResult;
    }

    public boolean succeeded() {
        return ConstantCode.SUCCESS.getCode().equals(code);
    }

    public boolean failed() {
        return !succeeded();
    }

}
