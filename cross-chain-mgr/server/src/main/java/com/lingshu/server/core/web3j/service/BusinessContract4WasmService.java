package com.lingshu.server.core.web3j.service;

import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil4Busi;
import com.lingshu.server.core.web3j.contract.BusinessContract4Wasm;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import com.lingshu.server.core.web3j.core.TransactionReceiptExt;
import com.lingshu.server.dto.BusinessContract4WasmCallRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnBean;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.util.HashMap;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-26
 */
@Slf4j
@Service
public class BusinessContract4WasmService {
    @Resource(name = "chainmakerChainClient4Busi")
    @Lazy
    private ChainmakerChainClient blockChainClient;
    @Resource(name = "chainmakerAccountUtil4Busi")
    private ChainmakerAccountUtil4Busi chainmakerAccountUtil;

    private BusinessContract4Wasm businessContract;

//    @Value("${contract-name.business-wasm}")
//    private String contractName;

//    @PostConstruct
//    @ConditionalOnBean(BlockChainClient.class) // 确保依赖的Bean存在
//    public void init() {
//        try {
//            this.businessContract = load();
//            log.info("Chainmaker-BusinessContract4Wasm合约加载成功");
//        } catch (Exception e) {
//            log.error("Chainmaker-BusinessContract4Wasm合约加载失败", e);
//            throw new RuntimeException("Chainmaker-BusinessContract4Wasm合约初始化失败", e);
//        }
//    }

    public BusinessContract4Wasm deploy(String contractName, String version, String filePath) throws Exception {
        try {
            BusinessContract4Wasm businessContract = BusinessContract4Wasm.deploy(chainmakerAccountUtil, blockChainClient,
                    contractName, version, filePath);
            log.info("Chainmaker-BusinessContract4Wasm合约部署成功:{}", businessContract);
            return businessContract;
        } catch (Exception e) {
            log.error("Chainmaker-BusinessContract4Wasm合约部署失败", e);
            throw new ApiException("合约部署失败,原因:" + e.getMessage());
        }
    }

    public BusinessContract4Wasm load(String contractName) {
        try {
            BusinessContract4Wasm businessContract = BusinessContract4Wasm.load(contractName, blockChainClient, chainmakerAccountUtil);
            log.info("Chainmaker-BusinessContract4Wasm合约加载成功:{}", businessContract);
            return businessContract;
        } catch (Exception e) {
            log.error("Chainmaker-BusinessContract4Wasm合约加载失败", e);
            throw new ApiException("合约加载失败,原因:" + e.getMessage());
        }
    }

    // 默认加载合约
//    public BusinessContract4Wasm load() {
//        if (StringUtils.isBlank(contractName)) {
//            log.error("请配置contract-name.business");
//            throw new ApiException("请配置Chainmaker-BusinessContract4Wasm合约");
//        }
//        return load(contractName);
//    }

    public TransactionReceipt call(BusinessContract4WasmCallRequest request) {
        Map<String, byte[]> paramsMap = new HashMap<>();
        paramsMap.put(request.getParamName(), request.getParamValue().getBytes());
        TransactionReceipt transactionReceipt = businessContract.call(request.getMethodName(), paramsMap);
        log.info("保存成功: {}", transactionReceipt);
        return transactionReceipt;
    }
}
