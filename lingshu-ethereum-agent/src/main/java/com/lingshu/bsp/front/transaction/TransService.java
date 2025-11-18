/**
 * Copyright 2014-2020 the original author or authors.
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package com.lingshu.bsp.front.transaction;


import cn.hutool.core.bean.BeanUtil;
import cn.hutool.core.collection.CollectionUtil;
import cn.hutool.json.JSONUtil;
import com.lingshu.bsp.front.base.config.SdkConfig;
import com.lingshu.bsp.front.base.properties.Constants;
import com.lingshu.bsp.front.event.WebSocketServerEthHandler;
import com.lingshu.bsp.front.event.WebSocketServerLsHandler;
import com.lingshu.bsp.front.rpcapi.eth.contract.DIDManagerEth;
import com.lingshu.bsp.front.rpcapi.eth.EthApiService;
import com.lingshu.bsp.front.rpcapi.eth.contract.BusiCenterEth;
import com.lingshu.bsp.front.rpcapi.ls.RpcApiService;
import com.lingshu.bsp.front.rpcapi.ls.contract.BusiNormalLs;
import com.lingshu.bsp.front.rpcapi.ls.contract.DIDManagerLs;
import com.lingshu.bsp.front.transaction.entity.DIDUpdate;
import com.lingshu.bsp.front.transaction.entity.VcCalReq;
import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import io.reactivex.functions.Consumer;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.bind.annotation.RequestBody;
import org.web3j.protocol.core.DefaultBlockParameterName;
import org.web3j.tx.gas.StaticGasProvider;

import javax.annotation.Resource;
import javax.validation.Valid;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * TransService. handle transactions of deploy/call contract
 */
@Slf4j
@Service
public class TransService {

    @Autowired
    private RpcApiService rpcApiService;
    @Autowired
    private Constants constants;
    @Autowired
    @Qualifier(value = "common")
    private CryptoSuite cryptoSuite;
    @Autowired
    private SdkConfig sdkConfig;
    @Autowired
    private EthApiService ethApiService;

    @Value("${ethGasLimit:30000000}")
    private Long ethGasLimit;

    public void addEthCrossChainMsgListener(String address) {
        log.info("listen eth {}", address);
        BusiCenterEth busiCenterEth = BusiCenterEth.load(address, ethApiService.web3j(),
                ethApiService.transactionManager(), new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        try {
            busiCenterEth.cROSS_CHAIN_VCEventFlowable(DefaultBlockParameterName.EARLIEST, DefaultBlockParameterName.LATEST).subscribe(new Consumer<BusiCenterEth.CROSS_CHAIN_VCEventResponse>() {
                @Override
                public void accept(BusiCenterEth.CROSS_CHAIN_VCEventResponse response) throws Exception {
                    List arrList = new ArrayList<>();
                    Map<String, Object> bean = BeanUtil.beanToMap(response.log);
                    arrList.add(bean);
                    log.info("eth event size: {}", arrList.size());
                    String jsonStr = JSONUtil.toJsonStr(arrList);
                    WebSocketServerEthHandler.broadcast(jsonStr);
                    log.info("{}", jsonStr);
                }
            });
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void addCrossChainMsgListener(String address) {
        log.info("listen {}", address);
        IClient sdkClient = rpcApiService.getSdkClient(1);
        BusiNormalLs busiCenter = BusiNormalLs.load(address, sdkClient, sdkClient.getCryptoSuite().getKeyPair());
        busiCenter.subscribeCROSS_CHAIN_VCEvent((s, i, list) -> {
            try {
                if (CollectionUtil.isEmpty(list)) {
                    return;
                }
                log.info("event size: {}", list.size());
                WebSocketServerLsHandler.broadcast(JSONUtil.toJsonStr(list));
            } catch (Exception e) {
                log.error("", e);
            }
        });
    }

    public TransactionReceipt ethVcCal(@Valid @RequestBody VcCalReq req) {
        BusiCenterEth didManagerEth = BusiCenterEth.load(req.getAddress(), ethApiService.getWeb3j(), ethApiService.getRawTransactionManager(),
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(ethGasLimit)));

        try {
            org.web3j.protocol.core.methods.response.TransactionReceipt receipt = didManagerEth.busiRun(req.getFuncName(), req.getVc()).send();
            TransactionReceipt result = BeanUtil.copyProperties(receipt, TransactionReceipt.class);
            result.setStatus(receipt.isStatusOK() ? 0 : 1);
            result.setContractAddress(req.getAddress());
            log.info("receipt: {}",JSONUtil.toJsonStr(result));
            return result;
        } catch (Exception e) {
            log.error("eth transaction fail", e);
            e.printStackTrace();
        }

        return null;
    }

    // ETH did更新
    public Object ethDidUpdate(@Valid @RequestBody DIDUpdate req) {
        DIDManagerEth didManagerEth = DIDManagerEth.load(req.getAddress(), ethApiService.getWeb3j(), ethApiService.getRawTransactionManager(),
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));

        try {
            org.web3j.protocol.core.methods.response.TransactionReceipt receipt = didManagerEth.updateDID(req.getDid(), req.getDidDoc()).send();
            TransactionReceipt result = BeanUtil.copyProperties(receipt, TransactionReceipt.class);
            result.setStatus(receipt.isStatusOK() ? 0 : 1);
            log.info("eth receipt: {}", JSONUtil.toJsonStr(result));
            return result;
        } catch (Exception e) {
            log.error("eth transaction fail", e);
            e.printStackTrace();
        }

        return null;
    }

    public TransactionReceipt vcCal(@Valid @RequestBody VcCalReq req) {
        if (req.getLedgerId() == null) {
            req.setLedgerId(sdkConfig.getDefaultLedgerId());
        }
        IClient client = rpcApiService.getSdkClient(req.getLedgerId());
        BusiNormalLs busiCenter = BusiNormalLs.load(req.getAddress(), client, client.getCryptoSuite().getKeyPair());
        TransactionReceipt tr = busiCenter.vcCal(req.getFuncName(), req.getVc());
        tr.setContractAddress(req.getAddress());
        log.info("result, {}", tr);
        return tr;
    }


    // 零数did更新
    public Object didUpdate(@Valid @RequestBody DIDUpdate req) {
        if (req.getLedgerId() == null) {
            req.setLedgerId(sdkConfig.getDefaultLedgerId());
        }
        IClient client = rpcApiService.getSdkClient(req.getLedgerId());
        DIDManagerLs didManagerLs = DIDManagerLs.load(req.getAddress(), client, client.getCryptoSuite().getKeyPair());
        TransactionReceipt tr = didManagerLs.updateDID(req.getDid(), req.getDidDoc());
        log.info("lingshu receipt: {}", JSONUtil.toJsonStr(tr));
        return tr;
    }

}


