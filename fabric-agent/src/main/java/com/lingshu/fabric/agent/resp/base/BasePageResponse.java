package com.lingshu.fabric.agent.resp.base;

import lombok.Data;
import lombok.experimental.Accessors;

import java.util.Collections;

/**
 * Entity class of page response info.
 */
@Data
@Accessors(chain = true)
public class BasePageResponse {

    private int code;
    private String message;
    private Object data = Collections.emptyList();
    private int totalCount;

    public BasePageResponse() {
    }

    public BasePageResponse(RetCode retcode) {
        this.code = retcode.getCode();
        this.message = retcode.getMessage();
    }

    public BasePageResponse(RetCode retcode, Object data, int totalCount) {
        this.code = retcode.getCode();
        this.message = retcode.getMessage();
        this.data = data;
        this.totalCount = totalCount;
    }

    public boolean isSuccess(){
        return ConstantCode.SUCCESS.getCode().equals(code);
    }
}
