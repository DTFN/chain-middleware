package com.lingshu.fabric.agent.exception;

import com.lingshu.fabric.agent.config.code.RetCode;

/**
 * business exception.
 */
public class AgentException extends RuntimeException {

    private static final long serialVersionUID = 1L;
    private RetCode retCode;

    /**
     * init by RetCode.
     */
    public AgentException(RetCode retCode) {
        super(retCode.getMessage());
        this.retCode = retCode;
    }

    /**
     * init by RetCode and Throwable.
     */
    public AgentException(RetCode retCode, Throwable cause) {
        super(retCode.getMessage(), cause);
        retCode.setMessage(cause.getMessage());
        this.retCode = retCode;
    }

    /**
     * init by code and msg.
     */
    public AgentException(int code, String msg) {
        super(msg);
        this.retCode = new RetCode(code, msg);
    }

    /**
     * init by code 、 msg and Throwable.
     */
    public AgentException(int code, String msg, Throwable cause) {
        super(msg, cause);
        this.retCode = new RetCode(code, msg);
    }

    public RetCode getRetCode() {
        return retCode;
    }
}
