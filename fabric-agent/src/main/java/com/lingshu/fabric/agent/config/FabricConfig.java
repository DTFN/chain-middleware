package com.lingshu.fabric.agent.config;

import com.alibaba.fastjson2.JSON;
import com.alibaba.fastjson2.JSONObject;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.event.GatewayInitedEvent;
import com.lingshu.fabric.agent.exception.AgentException;
import lombok.extern.slf4j.Slf4j;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Identities;
import org.hyperledger.fabric.gateway.Wallet;
import org.hyperledger.fabric.gateway.Wallets;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationContext;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.InputStreamSource;
import org.springframework.util.StreamUtils;

import java.io.ByteArrayInputStream;
import java.io.File;
import java.io.IOException;
import java.io.Reader;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.InvalidKeyException;
import java.security.PrivateKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.util.List;
import java.util.stream.Collectors;

@Configuration
@Slf4j
//@ConditionalOnProperty(name = "fabric.mode.link", havingValue = "true")
public class FabricConfig {

    public static String networkConfigPath = "connection-tls.json";
    @Autowired
    private ConstantProperties constantProperties;
    @Value("${fabric.mode.link}")
    private boolean linkFlag;
    @Autowired
    private ApplicationContext context;

    @Bean
    public Gateway fabricGateway() throws Exception {
        return initGateway();
    }

    public void reinitializeFabricGatewayBean() throws Exception {
        // 手动调用获取新的实例
        Gateway newFabricGateway = initGateway();
        // 通知重新注入新的实例
        context.publishEvent(new GatewayInitedEvent(newFabricGateway));
    }

    private Gateway initGateway() throws Exception {
        // 根据配置计算实际配置文件路径
        InputStreamSource resource = new ClassPathResource(networkConfigPath);
        String jsonConfig = StreamUtils.copyToString(resource.getInputStream(), StandardCharsets.UTF_8);
        jsonConfig = jsonConfig.replaceAll("\\{organizationsPath}", constantProperties.getOrganizationsPath());
        ByteArrayInputStream jsonConfigInputStream = new ByteArrayInputStream(jsonConfig.getBytes(StandardCharsets.UTF_8));

        log.info("== use networkConfig: {}", jsonConfig);
        JSONObject connConf = JSON.parseObject(jsonConfig);
        JSONObject organizations = connConf.getJSONObject("organizations");
        String orgName = organizations.keySet().stream().findFirst().get();
        JSONObject org = organizations.getJSONObject(orgName);
        String mspId = org.getString("mspid");
        String certificatePath = org.getJSONObject("signedCertPEM").getString("path");
        String privateKeyPath = org.getJSONObject("adminPrivateKeyPEM").getString("path");

//        File certificateFile = new ClassPathResource(certificatePath).getFile();
//        File privateKeyFile = new ClassPathResource(privateKeyPath).getFile();

        // 使用org1中的admin初始化一个网关wallet账户用于连接网络
//        String certificatePath = "/home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/cert.pem";
        X509Certificate certificate = readX509Certificate(Paths.get(certificatePath));

//        String privateKeyPath = "/home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/6ac25f833930b9ac48f4598ccc2f0b3c30868296accfd0bf6ca4d27082c1cfea_sk";
        PrivateKey privateKey = getPrivateKey(Paths.get(privateKeyPath));

        Wallet wallet = Wallets.newInMemoryWallet();

        wallet.put("admin", Identities.newX509Identity(mspId, certificate, privateKey));

        // 根据connection-tls.json 获取Fabric网络连接对象
        Gateway.Builder builder = Gateway.createBuilder()
                .identity(wallet, "admin")
                .networkConfig(jsonConfigInputStream);

        // 连接网关
        Gateway gateway = builder.connect();
        log.info("== init gateway success");
        return gateway;
    }

    private static X509Certificate readX509Certificate(final Path certificatePath) throws IOException, CertificateException {
        try (Reader certificateReader = Files.newBufferedReader(certificatePath, StandardCharsets.UTF_8)) {
            return Identities.readX509Certificate(certificateReader);
        }
    }

    private static PrivateKey getPrivateKey(final Path privateKeyPath) throws IOException, InvalidKeyException {
        try (Reader privateKeyReader = Files.newBufferedReader(privateKeyPath, StandardCharsets.UTF_8)) {
            return Identities.readPrivateKey(privateKeyReader);
        }
    }
}
