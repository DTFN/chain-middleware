package com.lingshu.server.core.web3j.service;

import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.DIDManager;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import com.lingshu.server.dto.DIDCreateRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.autoconfigure.condition.ConditionalOnBean;
import org.springframework.context.annotation.Lazy;
import org.springframework.stereotype.Service;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.PostConstruct;
import javax.annotation.Resource;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author: derrick
 * @since: 2025-08-21
 */
@Slf4j
@Service
public class DIDManagerService {
    @Resource(name = "chainmakerChainClientRelayer")
    @Lazy
    private ChainmakerChainClient blockChainClient;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    // 将DIDManager作为Bean注入
    private DIDManager didManager;

    @Value("${contract.relayer.did-manager}")
    private String contractName;

    private final ConcurrentHashMap<String, String> didCache = new ConcurrentHashMap<>();

    @PostConstruct
    @ConditionalOnBean(BlockChainClient.class) // 确保依赖的Bean存在
    public void init() {
        try {
            this.didManager = load();
            log.info("DIDManager合约加载成功:{}", didManager);
        } catch (Exception e) {
            log.error("DIDManager合约加载失败", e);
            throw new RuntimeException("合约初始化失败", e);
        }
    }

    public DIDManager deploy(String contractName) throws Exception {
        try {
            DIDManager didManager = DIDManager.deploy(chainmakerAccountUtil, blockChainClient, contractName);
            log.info("DIDManager合约部署成功:{}", didManager);
            return didManager;
        } catch (Exception e) {
            log.error("DIDManager合约部署失败", e);
            throw new ApiException("合约部署失败,原因:" + e.getMessage());
        }
    }

    public DIDManager load(String contractName) {
        try {
            DIDManager didManager = DIDManager.load(contractName, blockChainClient, chainmakerAccountUtil);
            log.info("DIDManager合约加载成功:{}", didManager);
            return didManager;
        } catch (Exception e) {
            log.error("DIDManager合约加载失败", e);
            throw new ApiException("合约加载失败,原因:" + e.getMessage());
        }
    }

    // 默认加载合约
    public DIDManager load() {
        if (StringUtils.isBlank(contractName)) {
            log.error("请配置contract-name.did-manager");
            throw new ApiException("请配置DID合约");
        }
        return load(contractName);
    }

    // 直接使用实例变量，无需重复加载
    public String getDIDDetails(String did) throws Exception {
        // 首先尝试从缓存获取
        String didDetails = didCache.get(did);
        if (didDetails != null) {
            log.info("从缓存获取DID详情: {}", did);
            return didDetails;
        }

        Boolean exist = didManager.doesDIDExist(did);
        if (!exist) {
            log.info("DID不存在: {}", did);
            return null;
        }

        // 缓存中没有则从链上获取
        didDetails = didManager.getDIDDetails(did);
        if (StringUtils.isNotBlank(didDetails)) {
            if (-1 != didDetails.indexOf("{") && -1 != didDetails.lastIndexOf("}")) {
                didDetails = didDetails.substring(didDetails.indexOf("{"), didDetails.lastIndexOf("}") + 1);
            }
        }
        log.info("从链上获取DID详情: {}", didDetails);

        // 将结果存入缓存
        if (StringUtils.isNotBlank(didDetails)) {
            didCache.put(did, didDetails);
        }
        return didDetails;
    }

    public boolean existsDID(String did) {
        Boolean res = didManager.doesDIDExist(did);
        log.info("DID是否存在: {}", res);
        return res;
    }

    public TransactionReceipt createDID(DIDCreateRequest request) {
//        if (existsDID(request.getDid())) {
//            log.error("DID已存在");
//            return "DID已存在";
//        }
        TransactionReceipt transactionReceipt = didManager.createDID(request.getDid(), request.getDidDocument());
        log.info("DID创建成功: {}", transactionReceipt);

        // 如果交易成功，更新缓存
        if (transactionReceipt != null) {
            didCache.put(request.getDid(), request.getDidDocument());
            log.info("DID缓存已更新: {}", request.getDid());
        }

        return transactionReceipt;
    }

    public TransactionReceipt updateDID(DIDCreateRequest request) {
//        if (!existsDID(request.getDid())) {
//            log.error("DID不存在");
//            return "DID不存在";
//        }
        TransactionReceipt transactionReceipt = didManager.updateDID(request.getDid(), request.getDidDocument());
        log.info("DID更新成功: {}", transactionReceipt);

        // 如果交易成功，更新缓存
        if (transactionReceipt != null) {
            didCache.put(request.getDid(), request.getDidDocument());
            log.info("DID缓存已更新: {}", request.getDid());
        }
        return transactionReceipt;
    }

    public String echo(String did) {
        String echo = didManager.echo(did);
        log.info("DID echo: {}", echo);
        return echo;
    }

    /**
     * 从缓存中移除指定的 DID
     *
     * @param did DID ID
     */
    public void removeFromCache(String did) {
        if (did != null) {
            didCache.remove(did);
            log.info("从缓存中移除DID: {}", did);
        }
    }

    /**
     * 清空整个 DID 缓存
     */
    public void clearCache() {
        didCache.clear();
        log.info("DID缓存已清空");
    }

    /**
     * 获取缓存大小
     *
     * @return 缓存中的 DID 数量
     */
    public int getCacheSize() {
        return didCache.size();
    }
}
