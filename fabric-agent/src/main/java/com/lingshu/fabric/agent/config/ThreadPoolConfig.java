package com.lingshu.fabric.agent.config;

import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;

import java.util.concurrent.ThreadPoolExecutor;

/**
 * ThreadPoolConfig
 *
 * @author XuHang
 * @since 2023/11/16
 **/
@Slf4j
@Configuration
public class ThreadPoolConfig {
    @Bean
    public ThreadPoolTaskScheduler initAsyncScheduler() {
        log.info("start initAsyncScheduler init...");
        ThreadPoolTaskScheduler scheduler = new ThreadPoolTaskScheduler();
        scheduler.setPoolSize(20);
        scheduler.afterPropertiesSet();
        scheduler.setThreadNamePrefix("ThreadPoolTaskScheduler-async-init:");
        scheduler.setRejectedExecutionHandler(new ThreadPoolExecutor.CallerRunsPolicy());
        return scheduler;
    }
}
