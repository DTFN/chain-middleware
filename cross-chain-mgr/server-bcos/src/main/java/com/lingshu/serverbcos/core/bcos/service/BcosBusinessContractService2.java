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
public class BcosBusinessContractService2 {
    @Autowired
    private Client client;

    private HelloWorld businessContract;

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
