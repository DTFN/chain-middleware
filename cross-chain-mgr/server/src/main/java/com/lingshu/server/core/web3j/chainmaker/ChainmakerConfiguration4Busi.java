package com.lingshu.server.core.web3j.chainmaker;

import cn.hutool.core.util.ReflectUtil;
import cn.hutool.extra.spring.SpringUtil;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import lombok.extern.slf4j.Slf4j;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.ChainManager;
import org.chainmaker.sdk.User;
import org.chainmaker.sdk.config.AuthType;
import org.chainmaker.sdk.config.NodeConfig;
import org.chainmaker.sdk.config.SdkConfig;
import org.chainmaker.sdk.utils.FileUtils;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

/**
 * @author gongrui.wang
 * @since 2025/2/10
 */
@Slf4j
@ConditionalOnProperty(prefix = "chainmaker-business", name = "enable", havingValue = "true")
@Configuration
public class ChainmakerConfiguration4Busi implements CommandLineRunner {

    @Bean("chainmakerPropertiesBusi")
    @ConfigurationProperties(prefix = "chainmaker-business")
    public ChainmakerProperties chainClientConfigBusi() {
        return new ChainmakerProperties();
    }

    @Bean("chainClientBusi")
    public ChainClient createChainmakerChainClient(@Qualifier("chainmakerPropertiesBusi") ChainmakerProperties chainmakerProperties) throws Exception {

        // 配置ca服务路径
        if (chainmakerProperties.getCryptoConfigPath() != null) {
            String parentPath = new File("").toURI().toString();
            // 检查 org.chainmaker.sdk.shanghai.crypto_env.CryptoEnv 是否存在，并给静态变量 CryptoConfigPath 赋值
            Class<?> cryptoEnvCls = Class.forName("org.chainmaker.sdk.shanghai.crypto_env.CryptoEnv");
            if (cryptoEnvCls != null) {
                ReflectUtil.setFieldValue(cryptoEnvCls, "CryptoConfigPath", parentPath + chainmakerProperties.getCryptoConfigPath());
            }
        }

        // 证书模式
        if (chainmakerProperties.getAuthType().equals(AuthType.PermissionedWithCert.getMsg())) {
            for (NodeConfig nodeConfig : chainmakerProperties.getNodes()) {
                List<byte[]> tlsCaCertList = new ArrayList<>();
                for (String rootPath : nodeConfig.getTrustRootPaths()) {
                    List<String> filePathList = FileUtils.getFilesByPath(rootPath);
                    for (String filePath : filePathList) {
                        tlsCaCertList.add(FileUtils.getFileBytes(filePath));
                    }
                }
                byte[][] tlsCaCerts = new byte[tlsCaCertList.size()][];
                tlsCaCertList.toArray(tlsCaCerts);
                nodeConfig.setTrustRootBytes(tlsCaCerts);
            }
        }

        ChainManagerPriv chainManager = ChainManagerPriv.getInstance();
        SdkConfig sdkConfig = new SdkConfig();
        sdkConfig.setChainClient(chainmakerProperties);
        ChainClient chainClient = chainManager.createChainClient(sdkConfig);
        // sdk的问题,需要自己设置私钥的字节信息,sdk内部调用会使用到
        User clientUser = chainClient.getClientUser();
//        clientUser.setPriBytes(chainmakerProperties.getUserSignKeyBytes());
        byte[] privateKeyBytes = FileUtils.getFileBytes(chainmakerProperties.getUserSignKeyFilePath());
        clientUser.setPriBytes(privateKeyBytes);

        return chainClient;
    }

    @Bean("chainmakerChainClient4Busi")
    public BlockChainClient chainmakerChainClient4Busi(@Qualifier("chainClientBusi") ChainClient chainClient) {
        return new ChainmakerChainClient(chainClient);
    }

    @Override
    public void run(String... args) throws Exception {
        //log.info(JSONUtil.toJsonPrettyStr(SpringUtil.getBean(ChainmakerProperties.class)));
        ChainmakerChainClient chainmakerChainClient = SpringUtil.getBean("chainmakerChainClient4Busi", ChainmakerChainClient.class);
        log.info("busi chainmaker block height: {}", chainmakerChainClient.blockNumber());
    }
}
