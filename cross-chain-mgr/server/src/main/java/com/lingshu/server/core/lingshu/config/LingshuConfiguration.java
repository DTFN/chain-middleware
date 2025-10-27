package com.lingshu.server.core.lingshu.config;

import cn.hutool.extra.spring.SpringUtil;
import com.lingshu.chain.sdk.LingShuChainSDK;
import com.lingshu.chain.sdk.client.IClient;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.util.ResourceUtils;

import java.io.File;

/**
 * @author gongrui.wang
 * @since 2025/2/10
 */
@Slf4j
@ConditionalOnProperty(prefix = "lingshu", name = "enable", havingValue = "true")
@Configuration
public class LingshuConfiguration implements CommandLineRunner {
    @Value("${lingshu-config-file-path}")
    private String lingshuConfigFilePath;

    @Bean
    public IClient createLingshuChainClient() throws Exception {
        File configFile = ResourceUtils.getFile(lingshuConfigFilePath);
        LingShuChainSDK lingShuChainSDK = LingShuChainSDK.build(configFile.getPath());
        IClient client = lingShuChainSDK.getClient(1);
        return client;
    }

    @Override
    public void run(String... args) throws Exception {
        IClient client = SpringUtil.getBean(IClient.class);
        log.info("lingshu block height: {}", client.blockNum());
    }
}
