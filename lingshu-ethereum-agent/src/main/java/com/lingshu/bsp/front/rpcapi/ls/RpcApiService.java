/**
 * Copyright 2014-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package com.lingshu.bsp.front.rpcapi.ls;

import cn.hutool.core.collection.CollectionUtil;
import cn.hutool.json.JSONUtil;
import com.lingshu.bsp.front.base.code.ConstantCode;
import com.lingshu.bsp.front.base.config.NodeConfig;
import com.lingshu.bsp.front.base.config.SdkConfig;
import com.lingshu.bsp.front.base.exception.FrontException;
import com.lingshu.bsp.front.rpcapi.ls.contract.BusiNormalLs;
import com.lingshu.bsp.front.rpcapi.ls.contract.CrossSaveLs;
import com.lingshu.bsp.front.util.CommonUtils;
import com.lingshu.bsp.front.util.JsonUtils;
import com.lingshu.chain.sdk.LingShuChainSDK;
import com.lingshu.chain.sdk.LingShuChainSDKException;
import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.chain.sdk.client.protocol.response.*;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.util.NumberUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.math.BigInteger;
import java.util.*;

/**
 * RpcApi manage.
 */
@Slf4j
@Service
public class RpcApiService {

    @Autowired
    private NodeConfig nodeConfig;
    @Autowired
    private SdkConfig sdkConfig;
    @Autowired
    private LingShuChainSDK chainSDK;

    /**
     * getBlockNumber.
     */
    public BigInteger getBlockNumber(int ledgerId) {
        BigInteger blockNumber = getSdkClient(ledgerId).blockNum().getBlockNum();
        return blockNumber;
    }

    /**
     * getBlockByNumber.
     *
     * @param blockNumber blockNumber
     */
    public ChainBlock.Block getBlockByNumber(int ledgerId, BigInteger blockNumber) {
        if (blockNumberCheck(ledgerId, blockNumber)) {
            throw new FrontException(ConstantCode.BLOCK_NUMBER_ERROR);
        }

        ChainBlock.Block chainBlock = getSdkClient(ledgerId).blockByNum(blockNumber, false, false).getBlock();
        CommonUtils.processBlockHexNumber(chainBlock);
        return chainBlock;
    }

    /**
     * getBlockByHash.
     *
     * @param blockHash blockHash
     */
    public ChainBlock.Block getBlockByHash(int ledgerId, String blockHash) {
        ChainBlock.Block chainBlock = getSdkClient(ledgerId).blockByHash(blockHash, false, false).getBlock();
        CommonUtils.processBlockHexNumber(chainBlock);
        return chainBlock;
    }

    /**
     * getBlockTransCntByNumber.
     *
     * @param blockNumber blockNumber
     */
    public int getBlockTransCntByNumber(int ledgerId, BigInteger blockNumber) {
        int transCnt;
        if (blockNumberCheck(ledgerId, blockNumber)) {
            throw new FrontException(ConstantCode.BLOCK_NUMBER_ERROR);
        }
        ChainBlock.Block chainBlock = getSdkClient(ledgerId).blockByNum(blockNumber, false, false).getBlock();
        transCnt = chainBlock.getTransactions().size();
        return transCnt;
    }

    /**
     * getPbftView.
     */
    public BigInteger getPbftView(int ledgerId) {

        BigInteger result;
        result = NumberUtil.toBigInt(getSdkClient(ledgerId).txMetrics().getTxMetrics().getPendingTxSize());
        return result;
    }

    /**
     * getTransactionReceipt.
     *
     * @param transHash transHash
     */
    public TransactionReceipt getTransactionReceipt(int ledgerId, String transHash) {
        TransactionReceipt  transactionReceipt = getSdkClient(ledgerId).receipt(transHash, true);
        CommonUtils.decodeReceipt(transactionReceipt, getSdkClient(ledgerId).getCryptoSuite());
        CommonUtils.processReceiptHexNumber(transactionReceipt);
        return transactionReceipt;
    }

    /**
     * getTransactionByHash.
     *
     * @param transHash transHash
     */
    public JsonTransactionResponse getTransactionByHash(int ledgerId, String transHash) {
        ChainTransaction chainTransaction = getSdkClient(ledgerId).transaction(transHash, true);
        Optional<JsonTransactionResponse> transactionResponseOpt = chainTransaction.getTransaction();
        JsonTransactionResponse transaction = null;
        if (transactionResponseOpt.isPresent()) {
            transaction = transactionResponseOpt.get();
        }
        CommonUtils.processTransHexNumber(transaction);
        return transaction;
    }


    /**
     * getCode.
     *
     * @param address address
     * @param blockNumber blockNumber
     */
    public String getCode(int ledgerId, String address, BigInteger blockNumber) {
        String code;
        if (blockNumberCheck(ledgerId, blockNumber)) {
            throw new FrontException(ConstantCode.BLOCK_NUMBER_ERROR);
        }
        code = getSdkClient(ledgerId)
                .bytecode(address).getBytecode();
        return code;
    }

    private boolean blockNumberCheck(int ledgerId, BigInteger blockNumber) {
        BigInteger currentNumber = null;
        currentNumber = getSdkClient(ledgerId).blockNum().getBlockNum();
        log.debug("**** currentNumber:{}", currentNumber);
        return (blockNumber.compareTo(currentNumber) > 0);

    }


    public List<String> getLedgerPeers(int ledgerId) {
        List<String> ledgerPeerList = getSdkClient(ledgerId).participantIds().getParticipantIds();
        return ledgerPeerList;
    }

    public List<String> getNodeIdList() {
        return chainSDK.getClient().peerIds().getResult();
    }

    // get all peers of chain
    public List<Peers.PeerInfo> getPeers(int ledgerId) {
        return getSdkClient(ledgerId)
                .peers()
                .getPeers();
    }

    /**
     * get BasicConsensusInfo and List of ViewInfo
     *
     * @param ledgerId
     * @return
     */
    public ConsensusState.ConsensusStateInfo getConsensusStatus(int ledgerId) {
        return getSdkClient(ledgerId).consensusState().getConsensusState();
    }

    public SyncState.SyncStateInfo getSyncState(int ledgerId) {
        return getSdkClient(ledgerId).syncState().getSyncState();
    }

    /**
     * get getSystemConfigByKey of tx_count_limit/tx_gas_limit
     * @param ledgerId
     * @param key
     * @return value for config, ex: return 1000
     */
//    public String getSystemConfigByKey(int ledgerId, String key) {
//        return getgetSdkClient(ledgerId)(ledgerId)
//                .config(key)
//                .getSystemConfig();
//    }

    /**
     * getNodeInfo.
     */
    public ChainSystemInfo.SystemInfo getNodeInfo() {
        return chainSDK.getClient().systemInfo().getResult();
    }

    /**
     * get node config info
     * @return
     */
    public Object getNodeConfig() {
        return JsonUtils.toJavaObject(nodeConfig.toString(), Object.class);
    }

//    public int getPendingTransactions(int ledgerId) {
//        return getSdkClient(ledgerId).pendingTxSize().getPendingTxSize().intValue();
//    }
//
//    public BigInteger getPendingTransactionsSize(int ledgerId) {
//        return getSdkClient(ledgerId).pendingTxSize().getPendingTxSize();
//    }

    public List<String> getValidatorList(int ledgerId) {
        return getSdkClient(ledgerId).validatorIds().getValidatorIds();
    }

    public List<String> getWitnessList(int ledgerId) {
        return getSdkClient(ledgerId).witnessIds().getWitnessIds();
    }

    /* above v2.6.1*/
    public ChainBlock.Block getBlockHeaderByHash(Integer ledgerId, String blockHash,
        boolean returnSealers) {
        ChainBlock  chainBlock = getSdkClient(ledgerId).blockByHash(blockHash, true, false);
        ChainBlock.Block  block = chainBlock.getBlock();
        CommonUtils.processBlockHexNumber(block);
        return block;
    }

    public ChainBlock.Block getBlockHeaderByNumber(Integer ledgerId, BigInteger blockNumber,
        boolean returnSealers) {
        ChainBlock  chainBlock = getSdkClient(ledgerId).blockByNum(blockNumber, true, false);
        ChainBlock.Block  block = chainBlock.getBlock();
        CommonUtils.processBlockHexNumber(block);
        return block;
    }

    /**
     * get batch receipt in one block
     * @param ledgerId
     * @param blockNumber
     * @param start start index
     * @param count cursor, if -1, return all
     */
    public List<TransactionReceipt> getBatchReceiptByBlockNumber(int ledgerId, BigInteger blockNumber, int start, int count) {
        return null;
    }

    public List<TransactionReceipt> getBatchReceiptByBlockHash(int ledgerId, String blockHash, int start, int count) {
        return null;
    }

    public IClient getSdkClient(Integer ledgerId) {
        //this.checkConnection();
        IClient client;
        try {
            client= chainSDK.getClient(ledgerId);
        } catch (LingShuChainSDKException e) {
            String errorMsg = e.getMessage();
            log.error("LingShuChain SDK getClient failed: {}", errorMsg);
            // check client error type
            if (errorMsg.contains("available peers")) {
                log.error("no available node to connect to");
                throw new FrontException(ConstantCode.SYSTEM_ERROR_NODE_INACTIVE.getCode(),
                    "no available node to connect to");
            }
            if (errorMsg.contains("existence of ledger")) {
                log.error("ledger: {} of the connected node not exist!", ledgerId);
                throw new FrontException(ConstantCode.SYSTEM_ERROR_RPC_NULL.getCode(),
                    "ledger: " + ledgerId + " of the connected node not exist!");
            }
            if (errorMsg.contains("no peers set up the ledger")) {
                log.error("no peers belong to this ledger: {}!", ledgerId);
                throw new FrontException(ConstantCode.SYSTEM_ERROR_NO_NODE_IN_LEDGER.getCode(),
                    "no peers belong to this ledger: " + ledgerId);
            }
            throw new FrontException(ConstantCode.RPC_CLIENT_IS_NULL);
            // refresh ledger list
            // getLedgerList();
        }
        return client;
    }

    private void checkConnection() {
        if (!SdkConfig.PEER_CONNECTED) {
            throw new FrontException(ConstantCode.SYSTEM_ERROR_NODE_INACTIVE);
        }
    }

    private String getNodeIpPort() {
        return sdkConfig.getIp() + ":" + sdkConfig.getRpcPort();
    }


    public TransactionReceipt crossSave(String address) {
        // check ledgerId
        IClient client = getSdkClient(1);
        CrossSaveLs crossSaveLs = CrossSaveLs.load(address, client, client.getCryptoSuite().getKeyPair());
        return crossSaveLs.CrossChainSave("1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8", "1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8");
    }

    public String deployCrossSave() {
        // check ledgerId
        IClient client = getSdkClient(1);
        try {
            CrossSaveLs contract = CrossSaveLs.deploy(client, client.getCryptoSuite().getKeyPair());
            return contract.getContractAddress();
        } catch (Exception e) {
            throw new RuntimeException(e.getMessage());
        }
    }

    public String deployBusiCenter() {
        // check ledgerId
        IClient client = getSdkClient(1);
        try {
            BusiNormalLs contract = BusiNormalLs.deploy(client, client.getCryptoSuite().getKeyPair());
            return contract.getContractAddress();
        } catch (Exception e) {
            throw new RuntimeException(e.getMessage());
        }
    }
    public TransactionReceipt invokeBusiCenter(String address, List<Object> args) {
        // check ledgerId
        if (CollectionUtil.isEmpty(args)) {
            throw new RuntimeException("args are empty.");
        }
        IClient client = getSdkClient(1);
        BusiNormalLs contract = BusiNormalLs.load(address, client, client.getCryptoSuite().getKeyPair());
        return contract.erc20MintVcs(JSONUtil.toJsonStr(args.get(0)));
    }

}
