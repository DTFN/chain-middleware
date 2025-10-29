/*
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

import cn.hutool.core.bean.BeanUtil;
import com.lingshu.bsp.front.dto.*;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.chain.sdk.client.protocol.response.*;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.util.NumberUtil;
import io.swagger.annotations.Api;
import io.swagger.annotations.ApiImplicitParam;
import io.swagger.annotations.ApiImplicitParams;
import io.swagger.annotations.ApiOperation;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.math.BigInteger;
import java.util.List;

/**
 * RpcApiController.
 *
 */
@Api(value = "/rpc", tags = "rpc interface")
@RestController
@RequestMapping(value = "/{ledgerId}/rpc")
@ApiImplicitParams(@ApiImplicitParam(name = "ledgerId", value = "ledgerId",required = true,dataType = "Integer", paramType = "path"))
public class RpcApiController {

    @Autowired
    RpcApiService rpcApiService;

    @ApiOperation(value = "getBlockNumber", notes = "Get the latest block height of the node")
    @GetMapping("/blockNumber")
    public BigInteger getBlockNumber(@PathVariable int ledgerId) {
        return rpcApiService.getBlockNumber(ledgerId);
    }

    @ApiOperation(value = "getBlockByNumber", notes = "Get block information based on block height")
    @ApiImplicitParam(name = "blockNumber", value = "blockNumber", required = true,
            dataType = "BigInteger", paramType = "path")
    @GetMapping("/blockByNumber/{blockNumber}")
    public ChainBlock.Block getBlockByNumber(@PathVariable int ledgerId,
                                             @PathVariable BigInteger blockNumber) {
        return rpcApiService.getBlockByNumber(ledgerId, blockNumber);
    }

    @ApiOperation(value = "getBlockByHash", notes = "Get block information based on block hash")
    @ApiImplicitParam(name = "blockHash", value = "blockHash", required = true, dataType = "String",
            paramType = "path")
    @GetMapping("/blockByHash/{blockHash}")
    public ChainBlock.Block getBlockByHash(@PathVariable int ledgerId, @PathVariable String blockHash) {
        return rpcApiService.getBlockByHash(ledgerId, blockHash);
    }

    @ApiOperation(value = "getBlockTransCntByNumber",
            notes = "Get the number of transactions in the block based on the block height")
    @ApiImplicitParam(name = "blockNumber", value = "blockNumber", required = true,
            dataType = "BigInteger", paramType = "path")
    @GetMapping("/blockTransCnt/{blockNumber}")
    public int getBlockTransCntByNumber(@PathVariable int ledgerId,
            @PathVariable BigInteger blockNumber) {
        return rpcApiService.getBlockTransCntByNumber(ledgerId, blockNumber);
    }



    @ApiOperation(value = "getPbftView", notes = "Get PbftView")
    @GetMapping("/pbftView")
    public BigInteger getPbftView(@PathVariable int ledgerId) {
        return rpcApiService.getPbftView(ledgerId);
    }

    @ApiOperation(value = "getTransactionReceipt",
            notes = "Get a transaction receipt based on the transaction hash")
    @ApiImplicitParam(name = "transHash", value = "transHash", required = true, dataType = "String",
            paramType = "path")
    @GetMapping("/transactionReceipt/{transHash}")
    public TransactionReceipt getTransactionReceipt(@PathVariable int ledgerId,
            @PathVariable String transHash) {
        return rpcApiService.getTransactionReceipt(ledgerId, transHash);
    }

    @ApiOperation(value = "getTransactionByHash",
            notes = "Get transaction information based on transaction hash")
    @ApiImplicitParam(name = "transHash", value = "transHash", required = true, dataType = "String",
            paramType = "path")
    @GetMapping("/transaction/{transHash}")
    public TransactionDetail getTransactionByHash(@PathVariable int ledgerId,
                                                  @PathVariable String transHash) {
        JsonTransactionResponse response = rpcApiService.getTransactionByHash(ledgerId, transHash);
        TransactionDetail transactionDetail = BeanUtil.copyProperties(response, TransactionDetail.class);
        transactionDetail.setBlockNumber("0x"+NumberUtil.toHexStringNoPrefix(response.getBlockNumber()));
        return transactionDetail;
    }

    @ApiOperation(value = "getLedgerPeers", notes = "get list of ledger peers")
    @GetMapping("/ledgerPeers")
    public List<String> getLedgerPeers(@PathVariable int ledgerId) {
        return rpcApiService.getLedgerPeers(ledgerId);
    }

    @ApiOperation(value = "getNodeIDList", notes = "get list of node id")
    @GetMapping("/nodeIdList")
    public List<String> getNodeIDList() {
        return rpcApiService.getNodeIdList();
    }

    @ApiOperation(value = "getPeers", notes = "get list of peers")
    @GetMapping("/peers")
    public List<Peers.PeerInfo> getPeers(@PathVariable int ledgerId) {
        return rpcApiService.getPeers(ledgerId);
    }

    @ApiOperation(value = "getConsensusStatus", notes = "get consensus status of ledger")
    @GetMapping("/consensusStatus")
    public ConsensusState.ConsensusStateInfo getConsensusStatus(@PathVariable int ledgerId) {
        return rpcApiService.getConsensusStatus(ledgerId);
    }

    @ApiOperation(value = "getSyncState", notes = "get sync status of ledger")
    @GetMapping("/syncStatus")
    public SyncState.SyncStateInfo getSyncState(@PathVariable int ledgerId) {
        return rpcApiService.getSyncState(ledgerId);
    }

    @ApiOperation(value = "getNodeConfig", notes = "Get node config info")
    @GetMapping("/nodeConfig")
    public Object getNodeConfig() {
        return rpcApiService.getNodeConfig();
    }

    @ApiOperation(value = "validatorList", notes = "get list of ledger's validator list")
    @GetMapping("/validatorList")
    public List<String> getConsensusList(@PathVariable int ledgerId) {
        return rpcApiService.getValidatorList(ledgerId);
    }

    @ApiOperation(value = "witnessList", notes = "get list of ledger's witness list")
    @GetMapping("/witnessList")
    public List<String> getWitnessList(@PathVariable int ledgerId) {
        return rpcApiService.getWitnessList(ledgerId);
    }

    /* above 2.7.0 */
    @ApiOperation(value = "getBlockHeaderByHash", notes = "Get block header with validators based on block hash")
    @ApiImplicitParam(name = "blockHash", value = "blockHash", required = true,
        dataType = "String", paramType = "path")
    @GetMapping("/blockHeaderByHash/{blockHash}")
    public ChainBlock.Block getBlockHeaderByHash(@PathVariable int ledgerId,
        @PathVariable String blockHash) {
        return rpcApiService.getBlockHeaderByHash(ledgerId, blockHash, true);
    }

    @ApiOperation(value = "getBlockHeaderByNumber", notes = "Get block header with validators based on block height")
    @ApiImplicitParam(name = "blockNumber", value = "blockNumber", required = true,
        dataType = "BigInteger", paramType = "path")
    @GetMapping("/blockHeaderByNumber/{blockNumber}")
    public ChainBlock.Block getBlockHeaderByNumber(@PathVariable int ledgerId,
        @PathVariable BigInteger blockNumber) {
        return rpcApiService.getBlockHeaderByNumber(ledgerId, blockNumber, true);
    }


    /* above 2.7.0 */
    @ApiOperation(value = "getBatchReceiptByBlockNumber",
        notes = "Get the number of transactions in the block based on the block height")
    @ApiImplicitParams({
        @ApiImplicitParam(name = "blockNumber", value = "blockNumber", required = true,
            dataType = "BigInteger", paramType = "path"),
        @ApiImplicitParam(name = "start", value = "start", required = true,
            dataType = "int"),
        @ApiImplicitParam(name = "count", value = "count", required = true,
            dataType = "int")
    })
    @GetMapping("/transReceipt/batchByNumber/{blockNumber}")
    public List<TransactionReceipt> getBatchReceiptByBlockNumber(@PathVariable int ledgerId,
        @PathVariable BigInteger blockNumber,
        @RequestParam(value = "start", defaultValue = "0") int start,
        @RequestParam(value = "count", defaultValue = "-1") int count) {
        return rpcApiService.getBatchReceiptByBlockNumber(ledgerId, blockNumber, start, count);
    }

    @ApiOperation(value = "getBatchReceiptByBlockHash",
        notes = "Get the number of transactions in the block based on the block height")
    @ApiImplicitParams({
        @ApiImplicitParam(name = "blockHash", value = "blockHash", required = true,
            dataType = "String", paramType = "path"),
        @ApiImplicitParam(name = "start", value = "start", required = true,
            dataType = "int"),
        @ApiImplicitParam(name = "count", value = "count", required = true,
            dataType = "int")
    })
    @GetMapping("/transReceipt/batchByHash/{blockHash}")
    public List<TransactionReceipt> getBatchReceiptByBlockHash(@PathVariable int ledgerId,
        @PathVariable String blockHash,
        @RequestParam(value = "start", defaultValue = "0") int start,
        @RequestParam(value = "count", defaultValue = "-1") int count) {
        return rpcApiService.getBatchReceiptByBlockHash(ledgerId, blockHash, start, count);
    }

    @ApiOperation(value = "getNodeInfo", notes = "Get node information")
    @GetMapping("/nodeInfo")
    public ChainSystemInfo.SystemInfo getNodeInfo() {
        return rpcApiService.getNodeInfo();
    }

    @ApiOperation(value = "contract invoke", notes = "contract invoke")
    @PostMapping("/invoke")
    public TransactionReceipt invoke(@RequestBody RequestInvoke requestInvoke)
            throws Exception {
        TransactionReceipt transactionReceipt = rpcApiService.crossSave(requestInvoke.getContractAddress());
        transactionReceipt.setContractAddress(requestInvoke.getContractAddress());
        return transactionReceipt;
    }

    @ApiOperation(value = "contract deploy", notes = "contract deploy")
    @PostMapping("crossSave/deploy")
    public String deployCrossSave() throws Exception {
         return rpcApiService.deployCrossSave();
    }

    @ApiOperation(value = "contract deploy", notes = "contract deploy")
    @PostMapping("busi/deploy")
    public String deployBusiCenter() throws Exception {
         return rpcApiService.deployBusiCenter();
    }

    @ApiOperation(value = "contract invoke", notes = "contract invoke")
    @PostMapping("busi/invoke")
    public TransactionReceipt invokeBusiCenter(@RequestBody RequestInvoke requestInvoke)
            throws Exception {
        TransactionReceipt transactionReceipt = rpcApiService.invokeBusiCenter(requestInvoke.getContractAddress(), requestInvoke.getArgs());
        transactionReceipt.setContractAddress(requestInvoke.getContractAddress());
        return transactionReceipt;
    }
}
