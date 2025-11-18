package com.lingshu.server.core.lingshu.service;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.crypto.key.CryptoKeyPair;
import com.lingshu.chain.sdk.tx.common.exception.ContractException;
import com.lingshu.server.core.lingshu.contract.LsBusinessContract2;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;

@Slf4j
@Service
public class LingshuBusinessContractService2 {
    @Resource
    IClient client;

    private LsBusinessContract2 businessContract;

    public LsBusinessContract2 load(String address) {
        return getBusinessContract(address, client.getCryptoSuite().getKeyPair());
    }

    public LsBusinessContract2 getBusinessContract(String address, String userPriKey) {
        CryptoSuite cryptoSuite = client.getCryptoSuite();
        CryptoKeyPair cryptoKeyPair = cryptoSuite.loadKeyPair(userPriKey);
        return getBusinessContract(address, cryptoKeyPair);
    }

    public LsBusinessContract2 getBusinessContract(String address, CryptoKeyPair cryptoKeyPair) {
        LsBusinessContract2 businessContract = LsBusinessContract2.load(address, client, cryptoKeyPair);
        return businessContract;
    }

    public LsBusinessContract2 deploy() {
        LsBusinessContract2 businessContract = null;
        try {
            businessContract = LsBusinessContract2.deploy(client, client.getCryptoSuite().getKeyPair());
        } catch (ContractException e) {
            e.printStackTrace();
        }
        return businessContract;
    }


}
