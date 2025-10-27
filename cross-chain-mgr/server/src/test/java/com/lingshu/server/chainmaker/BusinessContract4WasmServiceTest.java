package com.lingshu.server.chainmaker;

import com.lingshu.server.ServerApplication;
import com.lingshu.server.core.web3j.contract.BusinessContract;
import com.lingshu.server.core.web3j.contract.BusinessContract4Wasm;
import com.lingshu.server.core.web3j.service.BusinessContract4WasmService;
import com.lingshu.server.core.web3j.service.BusinessContractService;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.Resource;

@SpringBootTest(classes = ServerApplication.class)
//使用main下resource的配置文件

public class BusinessContract4WasmServiceTest {
    @Resource
    private BusinessContract4WasmService businessContract4WasmService;

    @Test
    public void deploy() throws Exception {
        BusinessContract4Wasm businessContract = businessContract4WasmService.deploy("BUSINESS_CONTRACT-WASM",
                "1.0",
                "/Users/apple/workspace/lingshu/crossChain/dev/cross-chain-mgr/chainmaker/contract/busi_left.wasm");
        System.out.println(businessContract);
    }

    @Test
    public void test() throws Exception {
//        BusinessContract businessContract1 = businessContractService.load("test001");
//        System.out.println(businessContract1);

//        BusinessContract4Wasm businessContract = businessContract4WasmService.load("BUSINESS_CONTRACT12");
//        System.out.println(businessContract);
//
//        TransactionReceipt transactionReceipt = businessContract.call("echo", "content", "{\"content\":\"12345678\"}");
//        System.out.println(transactionReceipt);

    }
}
