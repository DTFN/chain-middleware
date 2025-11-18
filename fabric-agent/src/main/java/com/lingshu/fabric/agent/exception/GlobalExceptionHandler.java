package com.lingshu.fabric.agent.exception;

import com.google.common.collect.Maps;
import com.lingshu.fabric.agent.api.R;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.validation.BindException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.util.Date;
import java.util.Map;

/**
 * @author Brian
 * @since 08/09/2020
 */
@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {


    @Value("${base.enable-exception-log:false}")
    private Boolean enableExceptionLog;

    @ExceptionHandler(value = {Exception.class})
    public R<Object> otherException(Exception ex) {
        log.error("catch other exception", ex);
        if (ex instanceof AgentException){
            return R.failed(((AgentException) ex).getRetCode());
        }
        if (ex instanceof IllegalArgumentException) {
            return R.failed(ex.getMessage());
        }
        if (ex instanceof MethodArgumentNotValidException) {
            return R.failed(((MethodArgumentNotValidException) ex).getBindingResult().getFieldError().getDefaultMessage());
        }
        if(ex instanceof BindException){
            String defaultMessage = ((BindException) ex).getBindingResult().getFieldError().getDefaultMessage();
            return R.failed(defaultMessage);
        }

        log.error("系统异常", ex);
        return R.failed(ConstantCode.FAILED);
    }
}