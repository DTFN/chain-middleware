package com.lingshu.bsp.front.base.config;

import com.lingshu.bsp.front.base.config.entity.LingShuChainSDKWrapper;
import com.lingshu.chain.sdk.LingShuChainSDK;
import com.lingshu.chain.sdk.config.ConfigOption;
import com.lingshu.chain.sdk.config.exceptions.ConfigException;
import com.lingshu.chain.sdk.config.model.ConfigProperty;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.util.JsonUtil;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.HashMap;
import java.util.Map;

@Data
@Slf4j
@Configuration
@ConfigurationProperties(prefix = "sdk")
public class SdkConfig {
    public static boolean PEER_CONNECTED = true;
    public String certPath = "conf";
    public Boolean isGm;
    /* use String in java sdk*/
    private String corePoolSize;
    private String maxPoolSize;
    private String queueCapacity;
    /* use String in java sdk*/
    private String ip = "127.0.0.1";
    private String rpcPort = "56656";
    private String logLevel = "INFO";
    private Integer defaultLedgerId = 1;
    private Integer messageTimeout = 10000;
    private Integer receiptTimeout = 10000;
    private String keyStoreDir;
    private String fileFormat;

    @Bean
    public LingShuChainSDKWrapper getChainSDK() throws ConfigException {
        log.info("start init ConfigProperty");

        // common
        Map<String, Object> common = new HashMap<>();
        common.put("logLevel", logLevel);

        // account
        Map<String, Object> account = new HashMap<>();
        account.put("keyStoreDir", keyStoreDir);
        account.put("fileFormat", fileFormat);

        // cert config, encrypt type
        Map<String, Object> crypto = new HashMap<>();
        // cert use conf
        crypto.put("certDir", certPath);
        crypto.put("isGm", isGm);
        // user no need set this:crypto.put("sslCryptoType", encryptType);
        log.info("init cert crypto:{}, (using conf as cert path)", crypto);

        // peer, default one node in front
        Map<String, Object> network = new HashMap<>();
        String peer = ip + ":" + rpcPort;
        network.put("peer", peer);
        network.put("defaultLedgerId", defaultLedgerId);
        // network.put("messageTimeout", messageTimeout);
        network.put("receiptTimeout", receiptTimeout);
        log.info("init node network property :{}", peer);

        // thread pool config
//        log.info("init thread pool property");
//        Map<String, Object> threadPool = new HashMap<>();
//        threadPool.put("channelProcessorThreadSize", corePoolSize);
//        threadPool.put("receiptProcessorThreadSize", corePoolSize);
//        threadPool.put("maxBlockingQueueSize", queueCapacity);
//        log.info("init thread pool property:{}", threadPool);

        // init property
        ConfigProperty configProperty = new ConfigProperty();
        configProperty.setCrypto(crypto);
        configProperty.setNetwork(network);
//        configProperty.setThreadPool(threadPool);
        configProperty.setCommon(common);
        configProperty.setAccount(account);

        // init config option
        log.info("init configOption from configProperty: {}", JsonUtil.toJson(configProperty));
        ConfigOption configOption = new ConfigOption(configProperty);
        // init LingShuChain SDK
        log.info("init LingShuChain sdk instance, please check sdk.log");

        // 包装客户端
        LingShuChainSDKWrapper chainSDK = LingShuChainSDKWrapper.wrap(configOption);
        return chainSDK;
    }

    @Bean(name = "common")
    public CryptoSuite getCommonSuite(LingShuChainSDK chainSDK) {
        int encryptType = chainSDK.getConfig().getCryptoConfig().getSslCryptoType();
        log.info("getCommonSuite init encrypt type:{}", encryptType);
        return new CryptoSuite(encryptType);
    }
}

