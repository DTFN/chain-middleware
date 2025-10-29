package com.lingshu.bsp.front.base.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurerAdapter;

@Configuration
public class RequestInterceptorConfig extends WebMvcConfigurerAdapter {

    // 注解拦截器
    @Override
    public void addInterceptors(InterceptorRegistry registry) {

    }

}
