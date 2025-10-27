package com.lingshu.server.chainmaker;

import com.lingshu.server.ServerApplication;
import com.lingshu.server.core.account.entity.ChainAccountDO;
//import com.lingshu.server.core.account.service.ChainAccountService;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.DIDManager;
import com.lingshu.server.core.web3j.service.DIDManagerService;
import org.chainmaker.sdk.User;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.Resource;

@SpringBootTest(classes = ServerApplication.class)
public class DIDManagerServiceTest {

    @Resource
    private DIDManagerService didManagerService;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

//    @Test
//    public void transfer() {
//
//        nftTokenService.transfer("0xd717e889592a949b72f4f7e7c782b6aa77306e81", "0xa77ae5c0d9ebed0a1708988f015c5cb1da557de9",
//                34L, "85d837ba29d62e16eee6e2296a17b381b8fb9a65");
//    }

    @Test
    public void getUserAccount() {
//        ChainAccountDO chainAccountDO = chainAccountService.findByAccountAddress("0xd717e889592a949b72f4f7e7c782b6aa77306e81");
//        User user = chainmakerAccountUtil.toUser(chainAccountDO);
//        System.out.println(user);
    }

    @Test
    public void deploy() throws Exception {
        DIDManager didManager = didManagerService.deploy("test001");
        System.out.println(didManager);
    }

    @Test
    public void test() throws Exception {
        DIDManager didManager1 = didManagerService.load("test001");
        System.out.println(didManager1);

        DIDManager didManager = didManagerService.load("DID_MANAGER");
        System.out.println(didManager);

//        TransactionReceipt receipt = didManager.createDID("did-test001", "1234567");
//        System.out.println(receipt);

        String didDetails = didManager.getDIDDetails("did-test001");
        System.out.println(didDetails);

        TransactionReceipt transactionReceipt = didManager.updateDID("did-test001", "1234567890");
        System.out.println(transactionReceipt);

        String didDetails1 = didManager.getDIDDetails("did-test001");
        System.out.println(didDetails1);
    }
}