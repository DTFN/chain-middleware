package com.lingshu.server.core.web3j.chainmaker.client;

import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.web3j.core.BlockChainClient;
import lombok.AllArgsConstructor;
import lombok.Getter;
import org.chainmaker.pb.common.ChainmakerBlock;
import org.chainmaker.pb.common.ChainmakerTransaction;
import org.chainmaker.pb.discovery.Discovery;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.model.BlockInfo;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

/**
 * @author gongrui.wang
 * @since 2025/2/10
 */
@AllArgsConstructor
public class ChainmakerChainClient implements BlockChainClient {

    /**
     * 超时时间，默认3秒
     */
    public static final long REQUEST_TIMEOUT = 3 * 1000;

    @Getter
    private ChainClient chainClient;

    @Override
    public long blockNumber() {
        try {
            return chainClient.getCurrentBlockHeight(REQUEST_TIMEOUT);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Override
    public Integer peersNum() {
        try {
            Discovery.ChainInfo chainInfo = chainClient.getChainInfo(REQUEST_TIMEOUT);
            return chainInfo.getNodeListCount();
        } catch (Exception e) {
            throw new ApiException("获取链节点数量失败");
        }
    }

    @Override
    public TransactionReceipt receipt(String txHash) {
        try {
            ChainmakerTransaction.TransactionInfo transaction = chainClient.getTxByTxId(txHash, REQUEST_TIMEOUT);
            TransactionReceipt receipt = new TransactionReceipt();
            receipt.setBlockHash(transaction.getBlockHash().toStringUtf8());
            receipt.setBlockNumber(String.valueOf(transaction.getBlockHeight()));
            receipt.setTransactionHash(txHash);
            receipt.setTransactionIndex(String.valueOf(transaction.getTxIndex()));
            return receipt;
        } catch (Exception e) {
            throw new ApiException("查询交易回执失败");
        }
    }

    @Override
    public BlockInfo blockByNum(Long blockNo) {
        try {
            ChainmakerBlock.BlockInfo blockInfo =  chainClient.getBlockByHeight(blockNo, false, REQUEST_TIMEOUT);
            BlockInfo block = new BlockInfo();
            block.setFblockHeight(blockInfo.getBlock().getHeader().getBlockHeight());
            block.setFchainId(blockInfo.getBlock().getHeader().getChainId());
            block.setFcreateTime(String.valueOf(blockInfo.getBlock().getHeader().getBlockTimestamp()));
            return block;
        } catch (Exception e) {
            throw new ApiException("查询区块失败");
        }
    }
}
