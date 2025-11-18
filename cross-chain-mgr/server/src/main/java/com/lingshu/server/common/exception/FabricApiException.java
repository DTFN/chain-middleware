package com.lingshu.server.common.exception;

import com.lingshu.server.common.api.RetCode;

/**
 * @author lin
 * @since 2025-09-16
 */
public class FabricApiException extends RuntimeException {

    private static final long serialVersionUID = 1L;
    private RetCode retCode;

    /**
     * init by RetCode.
     */
    public FabricApiException(RetCode retCode) {
        super(retCode.getMessage());
        this.retCode = retCode;
    }

    /**
     * init by RetCode and Throwable.
     */
    public FabricApiException(RetCode retCode, Throwable cause) {
        super(retCode.getMessage(), cause);
        retCode.setMessage(cause.getMessage());
        this.retCode = retCode;
    }

    /**
     * init by code and msg.
     */
    public FabricApiException(int code, String msg) {
        super(msg);
        this.retCode = new RetCode(code, msg);
    }

    /**
     * init by code 、 msg and Throwable.
     */
    public FabricApiException(int code, String msg, Throwable cause) {
        super(msg, cause);
        this.retCode = new RetCode(code, msg);
    }

    public RetCode getRetCode() {
        return retCode;
    }
}
