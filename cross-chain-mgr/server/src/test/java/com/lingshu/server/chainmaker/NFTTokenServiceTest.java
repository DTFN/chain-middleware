package com.lingshu.server.chainmaker;

import com.lingshu.server.ServerApplication;
import com.lingshu.server.core.account.entity.ChainAccountDO;
//import com.lingshu.server.core.account.service.ChainAccountService;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.service.NFTTokenService;
import org.chainmaker.sdk.User;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;

import javax.annotation.Resource;

@SpringBootTest(classes = ServerApplication.class)
public class NFTTokenServiceTest {

    @Resource
    private NFTTokenService nftTokenService;
//    @Resource
//    private ChainAccountService chainAccountService;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    @Test
    public void transfer() {

        nftTokenService.transfer("0xd717e889592a949b72f4f7e7c782b6aa77306e81", "0xa77ae5c0d9ebed0a1708988f015c5cb1da557de9",
                34L, "85d837ba29d62e16eee6e2296a17b381b8fb9a65");
    }

    @Test
    public void getUserAccount() {
//        ChainAccountDO chainAccountDO = chainAccountService.findByAccountAddress("0xd717e889592a949b72f4f7e7c782b6aa77306e81");
//        User user = chainmakerAccountUtil.toUser(chainAccountDO);
//        System.out.println(user);
    }
}