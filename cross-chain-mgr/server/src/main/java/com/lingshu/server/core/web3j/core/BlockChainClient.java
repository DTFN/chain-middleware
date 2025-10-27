package com.lingshu.server.core.web3j.core;

import org.chainmaker.sdk.model.BlockInfo;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @author gongrui.wang
 * @since 2025/2/10
 */
public interface BlockChainClient {

    /**
     * 最新的区块高度
     * @return
     */
    long blockNumber();

    /**
     * 获取区块链节点数量
     * @return
     */
    Integer peersNum();


    /**
     * 根据Hash获取交易回执
     * @param txHash
     * @return
     */
    TransactionReceipt receipt(String txHash);

    /**
     * 查询区块信息
     * @param blockNo
     * @return
     */
    BlockInfo blockByNum(Long blockNo);
}
