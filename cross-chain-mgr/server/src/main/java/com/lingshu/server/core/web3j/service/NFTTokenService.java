package com.lingshu.server.core.web3j.service;

import com.lingshu.server.common.api.ApiException;
import com.lingshu.server.core.account.entity.ChainAccountDO;
//import com.lingshu.server.core.account.service.ChainAccountService;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.contract.NFTToken;
import lombok.extern.slf4j.Slf4j;
import org.chainmaker.pb.common.ChainmakerTransaction;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.User;
import org.chainmaker.sdk.model.BlockInfo;
import org.springframework.stereotype.Service;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import javax.annotation.Resource;
import java.math.BigInteger;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @Author wang jian
 * @Date 2022/6/17 18:12
 */
@Slf4j
@Service
public class NFTTokenService {

    @Resource(name = "chainmakerChainClientRelayer")
    private ChainmakerChainClient blockChainClient;
//    @Resource
//    private ChainAccountService chainAccountService;
    @Resource
    private ChainmakerAccountUtil chainmakerAccountUtil;

    private static ConcurrentHashMap<String, ConcurrentHashMap<String, NFTToken>> maps = new ConcurrentHashMap<>();


    public NFTToken deploy(String contractAddress, String name, String symbol) throws Exception {
        NFTToken nftToken = NFTToken.deploy(chainmakerAccountUtil, blockChainClient, contractAddress, name, symbol);
        return nftToken;
    }

    public String tokenURI(Long tokenId, String contractAddress) {
        NFTToken nftToken = NFTToken.load(contractAddress, blockChainClient, chainmakerAccountUtil);
        String tokenURI = nftToken.tokenURI(BigInteger.valueOf(tokenId));
        return tokenURI;
    }

    public TransactionReceipt genTransaction(String contractAddress, String txId, long gasUse) {
        ChainClient chainClient = blockChainClient.getChainClient();
        String address = chainmakerAccountUtil.accountAddress(chainClient.getClientUser());

        TransactionReceipt transactionReceipt = new TransactionReceipt();
        transactionReceipt.setTransactionHash(txId);
        // 获取调用合约的账户地址
        transactionReceipt.setFrom(address);
        transactionReceipt.setTo(contractAddress);
        transactionReceipt.setContractAddress(contractAddress);
        transactionReceipt.setGasUsed(String.valueOf(gasUse));
        return transactionReceipt;
    }

    public ChainmakerTransaction.TransactionInfo getTransactionByTxId(String txId) {
        ChainClient chainClient = blockChainClient.getChainClient();
        try {
            ChainmakerTransaction.TransactionInfo txByTxId = chainClient.getTxByTxId(txId, ChainmakerChainClient.REQUEST_TIMEOUT);
            return txByTxId;
        } catch (Exception e) {
            log.info("get TransactionInfo fail, txId:{}", txId);
        }
        return null;
    }

    public String mintAsync(Long tokenId, String toAddress, String contractAddress, String tokenURI) {
        NFTToken nftToken = NFTToken.load(contractAddress, blockChainClient, chainmakerAccountUtil);
        String txId = nftToken.mintAsync(toAddress, BigInteger.valueOf(tokenId), tokenURI);
        return txId;
    }

    public TransactionReceipt mint(Long tokenId, String toAddress, String contractAddress, String tokenURI) {
        NFTToken nftToken = NFTToken.load(contractAddress, blockChainClient, chainmakerAccountUtil);
        TransactionReceipt transactionReceipt =
            nftToken.mint(toAddress, BigInteger.valueOf(tokenId), tokenURI);
        return transactionReceipt;
    }

    public String ownerOf(Long tokenId, String contractAddress) throws Exception {
        NFTToken nftToken = NFTToken.load(contractAddress, blockChainClient, chainmakerAccountUtil);
        String address = nftToken.ownerOf(BigInteger.valueOf(tokenId));
        return address;
    }

    public TransactionReceipt transfer(String fromAddress, String toAddress, Long tokenId, String tokenAddress) {
        ChainAccountDO chainAccountDO = null;//chainAccountService.findByAccountAddress(fromAddress);
        if (chainAccountDO == null) {
            throw new ApiException("未找到发送者的账户");
        }
        User fromUser = chainmakerAccountUtil.toUser(chainAccountDO);
        NFTToken nftToken = NFTToken.load(tokenAddress, blockChainClient, chainmakerAccountUtil);
        TransactionReceipt txRecpt = nftToken.transferFrom(
                fromAddress, toAddress, BigInteger.valueOf(tokenId), fromUser);
        return txRecpt;
    }

    public BlockInfo getBlockByNumber(BigInteger blockNumber) {
        try {
            return blockChainClient.blockByNum(blockNumber.longValue());
        } catch (Exception e) {
            throw new ApiException("获取区块失败");
        }
    }

    public TransactionReceipt getTransactionReceiptByTxHash(String transactionHash) {
        try {
            return blockChainClient.receipt(transactionHash);
        } catch (Exception e) {
            log.error("fail to query receipt.txHash:{}", transactionHash, e);
            throw new ApiException("查询交易回执失败");
        }
    }
}
