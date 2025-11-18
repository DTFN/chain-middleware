package com.lingshu.server.core.ethereum.service;

import cn.hutool.core.bean.BeanUtil;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.ethereum.contract.BusinessContract2;
import com.lingshu.server.dto.BusinessContractSaveRequest;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.chainmaker.sdk.utils.FileUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.Request;
import org.web3j.protocol.core.methods.response.EthTransaction;
import org.web3j.protocol.core.methods.response.Transaction;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.StaticGasProvider;
import org.web3j.utils.Files;

import javax.annotation.PostConstruct;
import java.io.File;
import java.io.IOException;
import java.math.BigInteger;
import java.nio.file.Paths;
import java.util.List;

@Slf4j
@Service
public class EthereumBusinessContractService2 {
    @Autowired
    @Qualifier("ethereumWeb3j")
    private Web3j web3j;
    @Autowired
    @Qualifier("ethereumRawTransactionManager")
    private RawTransactionManager rawTransactionManager;

    @Value("${contract.business.eth.business-center}")
    private String contractAddress;

    private BusinessContract2 businessContract;

    @PostConstruct
    public void init() {
        try {
            this.businessContract = load();
            log.info("Ethereum-BusinessContract2合约加载成功:{}", businessContract);
        } catch (Exception e) {
            log.error("Ethereum-BusinessContract2合约加载失败", e);
            throw new RuntimeException("Ethereum-BusinessContract2合约初始化失败", e);
        }
    }

    // 默认加载合约
    public BusinessContract2 load() {
        if (StringUtils.isBlank(contractAddress)) {
            log.error("请配置contract-address.ethereum-business");
            throw new ApiException("请配置Ethereum-BusinessContract2业务合约地址");
        }
        return load(contractAddress);
    }

    public TransactionReceipt deploy(String filePath) throws Exception {
        String binary = Files.readString(new File(filePath));
        log.info("binary: {}", binary);
        RemoteCall<BusinessContract2> deploy = BusinessContract2.deploy(binary, web3j, rawTransactionManager,
                BigInteger.valueOf(1_000_000_000L), BigInteger.valueOf(20_000_000));
        TransactionReceipt result = null;
        try {
            result = deploy.send().getTransactionReceipt().get();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy info: {}", result);
        return result;
    }

    public BusinessContract2 load(String contractAddress) {
        BusinessContract2 businessContract = BusinessContract2.load(contractAddress, web3j, rawTransactionManager,
                BigInteger.valueOf(1_000_000_000L), BigInteger.valueOf(20_000_000));
        return businessContract;
    }

    public void ethGetTransactionByHash(String transactionHash) throws IOException {
        EthTransaction ethTransaction = web3j.ethGetTransactionByHash(transactionHash).send();
        Transaction transaction = ethTransaction.getResult();
        log.info("transaction: {}", transaction);
    }

//    public TransactionReceipt saveBusinessDetails(BusinessContract2SaveRequest request) {
//        TransactionReceipt transactionReceipt = null;
//        try {
//            transactionReceipt = businessContract
//                    .saveBusinessDetails(request.getBusinessAddress(), request.getDetails())
//                    .send();
//            log.info("保存成功: {}", transactionReceipt);
//
//        } catch (Exception e) {
//            e.printStackTrace();
//            log.error("保存失败: {}", e.getMessage());
//        }
//        return transactionReceipt;
//    }
//
//    public String getBusinessDetails(String businessAddress) throws Exception {
//
//        String details = businessContract.getBusinessDetails(businessAddress);
//        log.info("详情: {}", details);
//        return details;
//    }
//
//    public JsonTransactionResponse getTransactionByHash(String transHash) {
//        try {
//            Transaction transaction = web3j.ethGetTransactionByHash(transHash).send().getResult();
//            JsonTransactionResponse response = BeanUtil.copyProperties(transaction, JsonTransactionResponse.class,
//                    "blockNumber", "transactionIndex");
//            response.setBlockNumber(transaction.getBlockNumberRaw());
//            response.setTransactionIndex(transaction.getTransactionIndexRaw());
//            return response;
//        } catch (Exception e) {
//            e.printStackTrace();
//            return null;
//        }
//    }
//
//    public BigInteger getBlockNumber() {
//        try {
//            return web3j.ethBlockNumber().send().getBlockNumber();
//        } catch (IOException e) {
//            e.printStackTrace();
//        }
//        return null;
//    }

//    public void get(String transactionHash) {
//        Request<?, EthTransaction> ethTransactionRequest = web3j.ethGetTransactionByHash(transactionHash);
//        EthTransaction ethTransaction = ethTransactionRequest.send();
//        Transaction transaction = ethTransaction.getTransaction().get();
//        log.info("transaction info: {}", transaction);
//    }
}
