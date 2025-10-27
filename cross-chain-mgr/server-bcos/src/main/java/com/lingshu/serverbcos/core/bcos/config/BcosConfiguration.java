package com.lingshu.serverbcos.core.bcos.config;

import cn.hutool.extra.spring.SpringUtil;
import com.moandjiezana.toml.Toml;
import lombok.extern.slf4j.Slf4j;
import org.fisco.bcos.sdk.BcosSDK;
import org.fisco.bcos.sdk.client.Client;
import org.fisco.bcos.sdk.config.ConfigOption;
import org.fisco.bcos.sdk.config.exceptions.ConfigException;
import org.fisco.bcos.sdk.config.model.ConfigProperty;
import org.fisco.bcos.sdk.model.CryptoType;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.InputStreamSource;
import org.springframework.util.ResourceUtils;
import org.springframework.util.StreamUtils;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.nio.charset.StandardCharsets;

/**
 * @author: derrick
 * @since: 2025-09-04
 */
@Slf4j
@ConditionalOnProperty(prefix = "bcos", name = "enable", havingValue = "true")
@Configuration
public class BcosConfiguration implements CommandLineRunner {

    @Value("${bcos.config-path}")
    private String bcosConfigPath;

    @Bean
    public Client createBcosClient() throws IOException, ConfigException {
        File file = ResourceUtils.getFile(bcosConfigPath);
        FileInputStream fileInputStream = new FileInputStream(file);
        // 创建BcosSDK实例
        ConfigProperty configProperty = new Toml().read(fileInputStream).to(ConfigProperty.class);
        ConfigOption configOption = new ConfigOption(configProperty, CryptoType.ECDSA_TYPE);
        BcosSDK sdk = new BcosSDK(configOption);
        // 获取群组1的Client实例
        Client client = sdk.getClient(Integer.valueOf(1));
        return client;
    }

    @Override
    public void run(String... args) throws Exception {
        Client client = SpringUtil.getBean(Client.class);
        log.info("bcos block height: {}", client.getBlockNumber().getBlockNumber());
    }
}
