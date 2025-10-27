package com.lingshu.server.core.lingshu.service;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.crypto.key.CryptoKeyPair;
import com.lingshu.chain.sdk.tx.common.exception.ContractException;
import com.lingshu.server.core.lingshu.contract.Project;
import org.springframework.stereotype.Service;

import javax.annotation.Resource;

/**
 * @author gongrui.wang
 * @since 2023/8/31
 */
@Service
public class ProjectService {
    @Resource
    IClient client;

    public Project getProject(String address, String userPriKey) {
        CryptoSuite cryptoSuite = client.getCryptoSuite();
        CryptoKeyPair cryptoKeyPair = cryptoSuite.loadKeyPair(userPriKey);
        return getProject(address, cryptoKeyPair);
    }

    public Project getProject(String address, CryptoKeyPair cryptoKeyPair) {
        Project project = Project.load(address, client, cryptoKeyPair);
        return project;
    }

    public Project deploy() {
        Project project = null;
        try {
            project = Project.deploy(client, client.getCryptoSuite().getKeyPair());
        } catch (ContractException e) {
            e.printStackTrace();
        }
        return project;
    }
}
