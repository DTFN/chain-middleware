package com.lingshu.fabric.agent.resp.base;

import lombok.Data;

/**
 * Entity class of response info.
 */
@Data
public class BaseResponse {

    private int code;
    private String msg;
    private Object data;
    private String attachment;

    public BaseResponse() {}

    public BaseResponse(RetCode retcode) {
        this.code = retcode.getCode();
        this.msg = retcode.getMessage();
        this.attachment = retcode.getAttachment();
    }

    public BaseResponse(RetCode retcode, Object data) {
        this.code = retcode.getCode();
        this.msg = retcode.getMessage();
        this.attachment = retcode.getAttachment();
        this.data = data;
    }

    public boolean isSuccess() {
        return ConstantCode.SUCCESS.getCode().equals(code);
    }
}
