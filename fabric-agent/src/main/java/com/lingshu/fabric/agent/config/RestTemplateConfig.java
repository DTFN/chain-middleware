package com.lingshu.fabric.agent.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.client.SimpleClientHttpRequestFactory;
import org.springframework.web.client.RestTemplate;

/**
 * RestTemplateConfig
 *
 * @author XuHang
 * @since 2023/11/17
 **/
@Configuration
public class RestTemplateConfig {
    @Value("${restTemplate.healthz.readTimeOut:1000}")
    private int readTimeOut;

    @Value("${restTemplate.healthz.connectTimeout:1000}")
    private int connectTimeout;

    /**
     * resttemplate for deploy contract.
     */
    @Bean(name = "healthzRestTemplate")
    public RestTemplate getFabricRestTemplate() {
        SimpleClientHttpRequestFactory factory = new SimpleClientHttpRequestFactory();
        // ms
        factory.setReadTimeout(readTimeOut);
        // ms
        factory.setConnectTimeout(connectTimeout);
        return new RestTemplate(factory);
    }
}
