package com.lingshu.serverbcos.core.bcos.service;

import com.lingshu.serverbcos.common.api.ApiException;
import com.lingshu.serverbcos.core.bcos.contract.solidity.HelloWorld;
import com.lingshu.serverbcos.dto.BusinessContractSaveRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.fisco.bcos.sdk.client.Client;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import javax.annotation.PostConstruct;

/**
 * @author: derrick
 * @since: 2025-09-04
 */
@Slf4j
@Service
public class BcosBusinessContractService {
    @Autowired
    private Client client;
    @Value("${contract-address.bcos-business}")
    private String contractAddress;
    private HelloWorld businessContract;

    @PostConstruct
    public void init() {
        try {
            this.businessContract = load();
            log.info("BCOS-BusinessContract合约加载成功:{}", businessContract);
        } catch (Exception e) {
            log.error("BCOS-BusinessContract合约加载失败", e);
            throw new RuntimeException("BCOS-BusinessContract合约初始化失败", e);
        }
    }

    // 默认加载合约
    public HelloWorld load() {
        if (StringUtils.isBlank(contractAddress)) {
            log.error("请配置contract-address.bcos-business");
            throw new ApiException("请配置BCOS-BusinessContract业务合约地址");
        }
        return load(contractAddress);
    }

    public TransactionReceipt deploy() {
        TransactionReceipt result = null;
        try {
            HelloWorld businessContract = HelloWorld.deploy(client, client.getCryptoSuite().getCryptoKeyPair());
            result = businessContract.getDeployReceipt();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy info: {}", result);
        return result;
    }

    public HelloWorld load(String contractAddress) {
        HelloWorld businessContract = HelloWorld.load(contractAddress, client, client.getCryptoSuite().getCryptoKeyPair());
        return businessContract;
    }

    public TransactionReceipt saveBusinessDetails(BusinessContractSaveRequest request) {
        TransactionReceipt transactionReceipt = null;
        try {
            transactionReceipt = businessContract.set(request.getDetails());
            log.info("保存成功: {}", transactionReceipt);
        } catch (Exception e) {
            e.printStackTrace();
            log.error("保存失败: {}", e.getMessage());
        }
        return transactionReceipt;
    }

    public String getBusinessDetails(String businessAddress) throws Exception {

        String details = businessContract.get();
        //String details = businessContract.getBusinessDetails(businessAddress);
        log.info("详情: {}", details);
        return details;
    }
}
