package com.lingshu.server.core.web3j.service;

import cn.hutool.json.JSONUtil;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.ResourceDomain;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import com.lingshu.server.dto.ResourceDomainSaveRequest;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Qualifier;
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
 * @since: 2025-08-25
 */
@Slf4j
@Service
public class ResourceDomainService {
    @Resource(name = "chainmakerChainClientRelayer")
    @Lazy
    private ChainmakerChainClient blockChainClient;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    private ResourceDomain resourceDomain;

    @Value("${contract.relayer.resource-domain}")
    private String contractName;

    private final ConcurrentHashMap<String, String> cache = new ConcurrentHashMap<>();

    @PostConstruct
    @ConditionalOnBean(BlockChainClient.class) // 确保依赖的Bean存在
    public void init() {
        try {
            this.resourceDomain = load();
            log.info("ResourceDomain合约加载成功:{}", resourceDomain);
        } catch (Exception e) {
            log.error("ResourceDomain合约加载失败", e);
            throw new RuntimeException("合约初始化失败", e);
        }
    }


    public ResourceDomain deploy(String contractName) throws Exception {
        try {
            ResourceDomain resourceDomain = ResourceDomain.deploy(chainmakerAccountUtil, blockChainClient, contractName);
            log.info("ResourceDomain合约部署成功:{}", resourceDomain);
            return resourceDomain;
        } catch (Exception e) {
            log.error("ResourceDomain合约部署失败", e);
            throw new ApiException("合约部署失败,原因:" + e.getMessage());
        }
    }

    public ResourceDomain load(String contractName) {
        try {
            ResourceDomain resourceDomain = ResourceDomain.load(contractName, blockChainClient, chainmakerAccountUtil);
            log.info("ResourceDomain合约加载成功:{}", resourceDomain);
            return resourceDomain;
        } catch (Exception e) {
            log.error("ResourceDomain合约加载失败", e);
            throw new ApiException("合约加载失败,原因:" + e.getMessage());
        }
    }

    // 默认加载合约
    public ResourceDomain load() {
        if (StringUtils.isBlank(contractName)) {
            log.error("请配置contract-name.resource-domain");
            throw new ApiException("请配置链上资源域名管理合约");
        }
        return load(contractName);
    }

    public TransactionReceipt saveDomainDetails(ResourceDomainSaveRequest request) {
        //需要将request拼装成
        ResourceDomainDto resourceDomainDto = new ResourceDomainDto();
        resourceDomainDto.setResourceName(request.getResourceName());
        resourceDomainDto.setGatewayId(request.getGatewayId());
        resourceDomainDto.setChainRid(request.getChainRid());
        resourceDomainDto.setContractType(request.getContractType());
        resourceDomainDto.setContractAddress(request.getContractAddress());

        String resourceDomainStr = JSONUtil.toJsonStr(resourceDomainDto);

        TransactionReceipt transactionReceipt = resourceDomain.saveDomainDetails(request.getResourceName(), resourceDomainStr);
        log.info("保存成功: {}", transactionReceipt);

        // 如果交易成功，更新缓存
        if (transactionReceipt != null) {
            cache.put(request.getResourceName(), resourceDomainStr);
            log.info("DID缓存已更新: {}", request.getResourceName());
        }
        return transactionReceipt;
    }

    @Data
    public class ResourceDomainDto {
        private String resourceName;
        private String gatewayId;
        private String chainRid;
        private String contractType;
        private String contractAddress;
    }

    public TransactionReceipt updateDomainDetails(ResourceDomainSaveRequest request) {
        //需要将request拼装成
        // "{\"resource_name\":\"chain1:contract-busiCenter-chainmaker-wasm-9105-chainmaker-2eth\",\"gateway_id\":\"0\",\"chain_rid\":\"chain_cm_01\",\"contract_type\":\"2\",\"contract_address\":\"BUSINESS_CENTER_WASM_9105_chainmaker_2_eth\"}"的字符串
        String details = "{\"resource_name\":\"" + request.getResourceName() +
                "\",\"gateway_id\":\"" + request.getGatewayId() +
                "\",\"chain_rid\":\"" + request.getChainRid() +
                "\",\"contract_type\":\"" + request.getContractType() +
                "\",\"contract_address\":\"" + request.getContractAddress() +
                "\"}";
        TransactionReceipt transactionReceipt = resourceDomain.updateDomainDetails(request.getResourceName(), details);
        log.info("更新成功: {}", transactionReceipt);

        // 如果交易成功，更新缓存
        if (transactionReceipt != null) {
            cache.put(request.getResourceName(), details);
            log.info("DID缓存已更新: {}", request.getResourceName());
        }
        return transactionReceipt;
    }

    public String getDomainDetails(String domainAddress) throws Exception {
        // 首先尝试从缓存获取
        String details = cache.get(domainAddress);
        if (details != null) {
            log.info("从缓存获取链上资源域名详情: {}", domainAddress);
            return details;
        }

        Boolean exist = resourceDomain.doesDomainExist(domainAddress);
        if (!exist) {
            log.info("链上不存在: {}", domainAddress);
            return null;
        }

        // 缓存中没有则从链上获取
        details = resourceDomain.getDomainDetails(domainAddress);
        if(StringUtils.isNotBlank(details)) {
            if (-1 != details.indexOf("{") && -1 != details.lastIndexOf("}")) {
                details = details.substring(details.indexOf("{"), details.lastIndexOf("}") + 1);
            } else {
                //如果不是json格式的，那么数据就不对了
                throw new ApiException("解析链上资源失败");
            }
        }
        log.info("从链上获取链上资源域名详情: {}", details);

        // 将结果存入缓存
        if (StringUtils.isNotBlank(details)) {
            cache.put(domainAddress, details);
        }
        return details;
    }

    public boolean doesDomainExist(String domainAddress) {
        Boolean res = resourceDomain.doesDomainExist(domainAddress);
        log.info("是否存在: {}", res);
        return res;
    }
}
