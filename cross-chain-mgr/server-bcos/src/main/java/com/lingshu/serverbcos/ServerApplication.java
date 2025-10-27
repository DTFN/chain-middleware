package com.lingshu.serverbcos
        ;

import cn.hutool.extra.spring.EnableSpringUtil;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.ServletComponentScan;
import org.springframework.scheduling.annotation.EnableAsync;

// Spring Boot应用启动注解，启用自动配置和组件扫描
@SpringBootApplication
// 注释掉的Mapper扫描注解，用于指定MyBatis Mapper接口的扫描路径
//@MapperScan({"com.lingshu.server.core.*.dao", "com.chainable.infra.framework.security.mapper",
//        "com.lingshu.server.flow.dao", "com.lingshu.server.workbench.*.dao", "com.chainable.unitrust.mapper"})

//@MapperScan({"com.lingshu.server.core.*.dao"})
// 启用Servlet组件扫描，自动注册@WebServlet、@WebFilter、@WebListener等注解的组件
@ServletComponentScan
// 启用Spring的定时任务支持
//@EnableScheduling
// 启用Spring的异步方法执行支持
@EnableAsync
// 启用自定义的Spring工具类支持
@EnableSpringUtil
// 注释掉的组件扫描注解，用于指定Spring组件的扫描路径
//@ComponentScan(value = {"com.lingshu.*"})
//@ComponentScan(basePackages = {
//        "com.lingshu.server.controller", // Controller所在包
//        "com.lingshu.server.core" // 假设Configuration在这个包下
//})

public class ServerApplication {

    // 应用程序入口点
    public static void main(String[] args) {
        // 启动Spring Boot应用
        SpringApplication.run(ServerApplication.class, args);
    }
}
