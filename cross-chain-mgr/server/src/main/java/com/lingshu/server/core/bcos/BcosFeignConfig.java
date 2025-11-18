package com.lingshu.server.core.bcos;

import feign.RequestInterceptor;
import feign.RequestTemplate;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;

import java.nio.charset.StandardCharsets;

// 注意此类不要交给spring管理,不然所有feign调用都会使用该配置,除非你就是想这样.
public class BcosFeignConfig implements RequestInterceptor {

    @Override
    public void apply(RequestTemplate requestTemplate) {
//        ServletRequestAttributes attributes = (ServletRequestAttributes) RequestContextHolder.getRequestAttributes();
//        HttpServletRequest request = attributes.getRequest();

        requestTemplate.header(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE);
        requestTemplate.header(HttpHeaders.ACCEPT_CHARSET, StandardCharsets.UTF_8.name());
    }

//    @Bean(name = "notifyEncoder")
//    public Encoder encoder(ICompanyApiNotifySubscribeService companyApiNotifySubscribeService) {
//        return new NotifyEncoder(companyApiNotifySubscribeService);
//    }

//    @Bean(name="notifyDecoder")
//    public Decoder decoder() {
//        return new NotifyDecoder();
//    }
//
//    @Bean(name = "notifyErrorDecoder")
//    public ErrorDecoder errorDecoder() {
//        return (methodKey, response) -> {
//            try {
//                String body = Util.toString(response.body().asReader());
//                //logger.error("methodKey={}, responseBody={}", methodKey, body);
//                return new RRException(String.format("系统未知异常 response body: %s", body));
//            } catch (Exception e) {
//                return FeignException.errorStatus(methodKey, response);
//            }
//        };
//    }
}
