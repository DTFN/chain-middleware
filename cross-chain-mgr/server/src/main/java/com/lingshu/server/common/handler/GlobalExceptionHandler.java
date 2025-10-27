package com.lingshu.server.common.handler;

import com.alibaba.csp.sentinel.slots.block.BlockException;
//import com.chainable.infra.framework.common.exception.SecurityAuthException;
//import com.chainable.nftserver.common.api.ApiErrorCode;
//import com.chainable.nftserver.common.api.ApiException;
//import com.chainable.nftserver.common.api.OpenAPIResp;
//import com.chainable.nftserver.common.api.OpenAPIRespBuilder;
import com.lingshu.server.common.api.ApiErrorCode;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.common.api.OpenAPIResp;
import com.lingshu.server.common.api.OpenAPIRespBuilder;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

import java.lang.reflect.UndeclaredThrowableException;

@Slf4j
@RestControllerAdvice
public class GlobalExceptionHandler {

    // 捕捉自定义业务异常
    @ExceptionHandler(value = {
        ApiException.class})
    public OpenAPIResp exception(ApiException ex) {
        log.error("api exception：{}", ex);
        if (null != ex.getErrorCode()) {
            return OpenAPIRespBuilder.failure(ex.getErrorCode());
        } else {
            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, ex.getMessage());
        }

    }

//    // 捕捉自定义业务异常
//    @ExceptionHandler(value = {com.lingshu.server.common.api.ApiException.class})
//    public OpenAPIResp apiException(ApiException ex) {
//        log.error("api exception：{}", ex);
//        if (null != ex.getErrorCode()) {
//            return OpenAPIRespBuilder.failure(ex.getErrorCode());
//        } else {
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, ex.getMessage());
//        }
//
//    }

    // 被限流的方法上未定义 throws BlockException，会被JVM包装一层UndeclaredThrowableException
    // 这里捕捉限流异常进行判断处理
    @ExceptionHandler(value = {BlockException.class, UndeclaredThrowableException.class})
    public OpenAPIResp flowBlockException(Exception ex) {
        if (ex instanceof UndeclaredThrowableException) {
            // 如果不是包装的BlockException异常，由otherException处理
            UndeclaredThrowableException throwableException = (UndeclaredThrowableException) ex;
            if (!(throwableException.getUndeclaredThrowable() instanceof BlockException)) {
                return otherException(ex);
            }
        }
        log.error("异常msg：{}", "请求被限流");
        return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, "请求被限制，请稍后再试");
    }

//    // 捕捉其他所有异常
//    @ExceptionHandler(value = {
//        SecurityAuthException.class})
//    public OpenAPIResp securityAuthException(Exception ex) {
//        log.info("security exception:{}", ex);
//        SecurityAuthException securityAuthException = (SecurityAuthException) ex;
//        Integer code = securityAuthException.getCode();
//        ApiErrorCode apiErrorCode = ApiErrorCode.fromCode(code);
//        if (apiErrorCode.getCode().equals(ApiErrorCode.FAILED.getCode())) {
//            return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED, ex.getMessage() == null ? "服务器异常" :
//                ex.getMessage());
//        }
//        return OpenAPIRespBuilder.failure(apiErrorCode);
//    }

    // 捕捉其他所有异常
    @ExceptionHandler(value = {
        IllegalArgumentException.class,
        RuntimeException.class,
        Exception.class})
    public OpenAPIResp otherException(Exception ex) {
        String msg = ex.getMessage();
        log.error("异常msg：{}", ex);
        return OpenAPIRespBuilder.failure(ApiErrorCode.FAILED,  "服务器异常");
    }


}
