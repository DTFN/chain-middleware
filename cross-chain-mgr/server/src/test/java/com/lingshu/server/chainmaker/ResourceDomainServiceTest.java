package com.lingshu.server.chainmaker;

import com.lingshu.server.ServerApplication;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.ResourceDomain;
import com.lingshu.server.core.web3j.service.ResourceDomainService;
import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.Resource;

@SpringBootTest(classes = ServerApplication.class)
public class ResourceDomainServiceTest {

    @Resource
    private ResourceDomainService resourceDomainService;

    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    @Test
    public void deploy() throws Exception {
        ResourceDomain resourceDomain = resourceDomainService.deploy("RESOURCE_DOMAIN");
        System.out.println(resourceDomain);
    }

    @Test
    public void test() throws Exception {
        ResourceDomain resourceDomain1 = resourceDomainService.load("test001");
        System.out.println(resourceDomain1);

        ResourceDomain resourceDomain = resourceDomainService.load("RESOURCE_DOMAIN");
        System.out.println(resourceDomain);

        Boolean domainExist = resourceDomain.doesDomainExist("domain-test001");
        System.out.println(domainExist);

        if (domainExist) {
            String details = resourceDomain.getDomainDetails("domain-test001");
            System.out.println(details);
        } else {
            TransactionReceipt transactionReceipt = resourceDomain.saveDomainDetails("domain-test001", "1234567890");
            System.out.println(transactionReceipt);
        }

        String details = resourceDomain.getDomainDetails("domain-test001");
        System.out.println(details);
    }
}
