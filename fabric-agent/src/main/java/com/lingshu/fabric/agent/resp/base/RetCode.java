package com.lingshu.fabric.agent.resp.base;

import lombok.Data;

/**
 * class about exception code and message
 * @related: ConstantCode
 */
@Data
public class RetCode {

    private Integer code;
    private String message;
    private String attachment;

    public RetCode() {}

    public RetCode(int code, String message) {
        this.code = code;
        this.message = message;
    }

    public RetCode attach(Object attachment) {
        if (attachment != null){
            this.attachment = String.valueOf(attachment);
        }
        return this;
    }
    public static RetCode mark(int code, String message) {
        return new RetCode(code, message);
    }

    public static RetCode mark(Integer code) {
        return new RetCode(code, null);
    }
}
