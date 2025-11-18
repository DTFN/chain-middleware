package com.lingshu.server.core.bsc.config;

import cn.hutool.extra.spring.SpringUtil;
import lombok.extern.slf4j.Slf4j;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.WalletUtils;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.tx.response.TransactionReceiptProcessor;

import java.security.Security;

@Slf4j
@ConditionalOnProperty(prefix = "bsc", name = "enable", havingValue = "true")
@Configuration
public class BscConfiguration implements CommandLineRunner {
    static {
        Security.addProvider(new BouncyCastleProvider());
    }

    @Value("${bsc.ip}")
    private String bscIp;

    @Value("${bsc.port}")
    private String bscPort;

    @Bean("bscWeb3j")
    public Web3j createBscChainClient() {
        // 构建以太坊节点URL
        String bscUrl = "http://" + bscIp + ":" + bscPort;
        log.info("Connecting to bsc node: {}", bscUrl);

        // 加载客户端
        Web3j web3j = Web3j.build(new HttpService(bscUrl));
        return web3j;
    }

    @Bean("bscRawTransactionManager")
    public RawTransactionManager rawTransactionManager(@Qualifier("bscWeb3j")Web3j web3j) throws Exception {
        // 加载钱包
        String walletFile = "{\"address\":\"b185d579023ad5584907a46c51fd22d244b8f549\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"65974eb531b945a39345bbcdad263fbc1d99c1534514fc73f8aeeeb41fa3e252\",\"cipherparams\":{\"iv\":\"bfb5a3c705adb08ab9b0cad7bef3a70f\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"252e0b3ead0cb9de0e03d65f3b9332e97069162aa2ee0999255ce0fa366e1015\"},\"mac\":\"ecb46ddbf80a28d1b6f3adff5b241b789753b311a1d1e2e2fab1f786872e3eb4\"},\"id\":\"87475d1b-6bed-4cbf-afdc-77596116e3cf\",\"version\":3}";
        Credentials defaultCredentials = WalletUtils.loadJsonCredentials("123456789", walletFile);

        // 设置钱包
        TransactionReceiptProcessor receiptProcessor =
                new PollingTransactionReceiptProcessor(web3j, 50, 10000);
        RawTransactionManager rawTransactionManager = new RawTransactionManager(
                web3j,
                defaultCredentials,
                888,
                receiptProcessor
        );
        return rawTransactionManager;
    }

    @Override
    public void run(String... args) throws Exception {
        Web3j web3j = SpringUtil.getBean("bscWeb3j", Web3j.class);
        log.info("bsc block height: {}", web3j.ethBlockNumber().send().getBlockNumber());
    }
}
