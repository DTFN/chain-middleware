package com.lingshu.server.core.ethereum.service;

import cn.hutool.core.bean.BeanUtil;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.ethereum.contract.CrossSaveEth;
import com.lingshu.server.dto.BusinessContractSaveRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.RemoteFunctionCall;
import org.web3j.protocol.core.methods.response.Transaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.StaticGasProvider;

import javax.annotation.PostConstruct;
import java.io.IOException;
import java.math.BigInteger;

@Slf4j
@Service
public class EthereumCrossSaveBusinessContractService {
    @Autowired
    @Qualifier("ethereumWeb3j")
    private Web3j web3j;

    @Autowired
    @Qualifier("ethereumRawTransactionManager")
    private RawTransactionManager rawTransactionManager;

    @Value("${contract-address.ethereum-cross-save-business}")
    private String contractAddress;

    private CrossSaveEth crossSaveEth;

    @PostConstruct
    public void init() {
        try {
            this.crossSaveEth = load();
            log.info("EthereumCrossSave-BusinessContract合约加载成功:{}", crossSaveEth);
        } catch (Exception e) {
            log.error("EthereumCrossSave-BusinessContract合约加载失败", e);
            throw new RuntimeException("EthereumCrossSave-BusinessContract合约初始化失败", e);
        }
    }

    // 默认加载合约
    public CrossSaveEth load() {
        if (StringUtils.isBlank(contractAddress)) {
            log.error("请配置contract-address.ethereum-cross-save-business");
            throw new ApiException("请配置EthereumCrossSave-BusinessContract业务合约地址");
        }
        return load(contractAddress);
    }

    public TransactionReceipt deploy() {
        RemoteCall<CrossSaveEth> deploy = CrossSaveEth.deploy(web3j, rawTransactionManager,
                new StaticGasProvider(BigInteger.valueOf(1_000_000_000L), BigInteger.valueOf(20_000_000)));
        TransactionReceipt result = null;
        try {
            result = deploy.send().getTransactionReceipt().get();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy info: {}", result);
        return result;
    }

    public CrossSaveEth load(String contractAddress) {
        CrossSaveEth crossSaveEth = CrossSaveEth.load(contractAddress, web3j, rawTransactionManager,
                new StaticGasProvider(BigInteger.valueOf(1_000_000_000L), BigInteger.valueOf(20_000_000)));
        return crossSaveEth;
    }

    public TransactionReceipt saveBusinessDetails(BusinessContractSaveRequest request) {
        TransactionReceipt transactionReceipt = null;
        try {
            transactionReceipt = crossSaveEth
                    .CrossChainSave(request.getBusinessAddress(), request.getDetails())
                    .send();
            log.info("保存成功: {}", transactionReceipt);

        } catch (Exception e) {
            e.printStackTrace();
            log.error("保存失败: {}", e.getMessage());
        }
        return transactionReceipt;
    }

    public String getBusinessDetails(String businessAddress) throws Exception {

        RemoteFunctionCall<TransactionReceipt> functionCall = crossSaveEth.query(businessAddress);
        TransactionReceipt transactionReceipt = functionCall.send();
        String details = transactionReceipt.getTransactionHash();
        log.info("详情: {}", details);
        return details;
    }

    public JsonTransactionResponse getTransactionByHash(String transHash) {
        try {
            Transaction transaction = web3j.ethGetTransactionByHash(transHash).send().getResult();
            JsonTransactionResponse response = BeanUtil.copyProperties(transaction, JsonTransactionResponse.class,
                    "blockNumber", "transactionIndex");
            response.setBlockNumber(transaction.getBlockNumberRaw());
            response.setTransactionIndex(transaction.getTransactionIndexRaw());
            return response;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public TransactionReceipt crossSave(String contractAddress) {
        CrossSaveEth crossSaveEth = CrossSaveEth.load(contractAddress, web3j, rawTransactionManager,
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        try {
            TransactionReceipt result = crossSaveEth.CrossChainSave("hello", "hello").send();
            return result;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public BigInteger getBlockNumber() {
        try {
            return web3j.ethBlockNumber().send().getBlockNumber();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }
}
