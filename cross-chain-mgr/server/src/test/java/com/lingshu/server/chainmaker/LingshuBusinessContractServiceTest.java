package com.lingshu.server.chainmaker;

import com.lingshu.server.ServerApplication;
import com.lingshu.server.core.web3j.contract.BusinessContract;
import com.lingshu.server.core.web3j.contract.BusinessContract4Wasm;
import com.lingshu.server.core.web3j.service.BusinessContract4WasmService;
import com.lingshu.server.core.web3j.service.BusinessContractService;
import org.chainmaker.pb.common.ContractOuterClass;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.Resource;

@SpringBootTest(classes = ServerApplication.class)
public class LingshuBusinessContractServiceTest {

    @Resource
    private BusinessContractService businessContractService;

    @Test
    public void deploy() throws Exception {
        BusinessContract businessContract = businessContractService.deploy("BUSINESS_CONTRACT9",
                "1.0", "/Users/apple/workspace/lingshu/crossChain/dev/cross-chain-mgr/chainmaker/contract/BusinessContract.bin");
        System.out.println(businessContract);
    }

    @Test
    public void test() throws Exception {
//        BusinessContract businessContract1 = businessContractService.load("test001");
//        System.out.println(businessContract1);

//        BusinessContract businessContract = businessContractService.load("BUSINESS_CONTRACT9");
//        System.out.println(businessContract);
//
//        String details = businessContract.getBusinessDetails("business-test002");
//        System.out.println(details);

//        TransactionReceipt transactionReceipt = businessContract.saveBusinessDetails("business-test002", "1234567890");
//        System.out.println(transactionReceipt);


//        String details1 = businessContract.getBusinessDetails("business-test002");
//        System.out.println(details1);
    }
}
