package com.lingshu.server.core.web3j.service;

import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.BusinessContract;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import com.lingshu.server.dto.BusinessContractSaveRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnBean;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;

/**
 * @author: derrick
 * @since: 2025-08-26
 */
@Slf4j
@Service
public class BusinessContractService {
    @Resource(name = "chainmakerChainClientRelayer")
    @Lazy
    private ChainmakerChainClient blockChainClient;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    private BusinessContract businessContract;

//    @Value("${contract-name.business}")
//    private String contractName;
//    @PostConstruct
//    @ConditionalOnBean(BlockChainClient.class) // 确保依赖的Bean存在
//    public void init() {
//        try {
//            this.businessContract = load();
//            log.info("Chainmaker-BusinessContract合约加载成功:{}", businessContract);
//        } catch (Exception e) {
//            log.error("Chainmaker-BusinessContract合约加载失败", e);
//            throw new RuntimeException("Chainmaker-BusinessContract合约初始化失败", e);
//        }
//    }

    public BusinessContract deploy(String contractName, String version, String filePath) throws Exception {
        try {
            BusinessContract businessContract = BusinessContract.deploy(chainmakerAccountUtil, blockChainClient,
                    contractName, version, filePath);
            log.info("Chainmaker-BusinessContract合约部署成功:{}", businessContract);
            return businessContract;
        } catch (Exception e) {
            log.error("Chainmaker-BusinessContract合约部署失败", e);
            throw new ApiException("合约部署失败,原因:" + e.getMessage());
        }
    }

    public BusinessContract load(String contractName) {
        try {
            BusinessContract businessContract = BusinessContract.load(contractName, blockChainClient, chainmakerAccountUtil);
            log.info("Chainmaker-BusinessContract合约加载成功:{}", businessContract);
            return businessContract;
        } catch (Exception e) {
            log.error("Chainmaker-BusinessContract合约加载失败", e);
            throw new ApiException("合约加载失败,原因:" + e.getMessage());
        }
    }

//    // 默认加载合约
//    public BusinessContract load() {
//        if (StringUtils.isBlank(contractName)) {
//            log.error("请配置contract-name.business");
//            throw new ApiException("请配置业务合约");
//        }
//        return load(contractName);
//    }

//    public TransactionReceipt saveBusinessDetails(BusinessContractSaveRequest request) {
//        TransactionReceipt transactionReceipt = businessContract.saveBusinessDetails(request.getBusinessAddress(), request.getDetails());
//        log.info("保存成功: {}", transactionReceipt);
//        return transactionReceipt;
//    }

//    public String getBusinessDetails(String businessAddress) throws Exception {
//        //todo test
//        BusinessContract businessContract1 = load("BUSINESS_CONTRACT");
//        String details = businessContract1.getBusinessDetails(businessAddress);
//        log.info("详情: {}", details);
//        return details;
//    }
}
