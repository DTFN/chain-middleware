package com.lingshu.server.core.web3j.contract;

import cn.hutool.core.convert.Convert;
import cn.hutool.core.util.HexUtil;
import com.lingshu.server.common.exception.ChangAnTxException;
import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.google.common.collect.ImmutableMap;
import com.google.protobuf.ByteString;
import com.lingshu.server.core.web3j.core.TransactionReceiptExt;
import com.lingshu.server.utils.ChainResultParser;
import com.lingshu.server.utils.EthDecoderUtil;
import lombok.Data;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.chainmaker.pb.common.ChainmakerTransaction;
import org.chainmaker.pb.common.ContractOuterClass;
import org.chainmaker.pb.common.Request;
import org.chainmaker.pb.common.ResultOuterClass;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.User;
import org.chainmaker.sdk.utils.SdkUtils;
import org.chainmaker.sdk.utils.Utils;
import org.web3j.abi.FunctionEncoder;
import org.web3j.abi.datatypes.Address;
import org.web3j.abi.datatypes.Function;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.utils.Numeric;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;

/**
 * @author gongrui.wang
 * @since 2025/2/12
 */
@Slf4j
public abstract class ChainmakerContract {

    private static final String CONTRACT_DEFAULT_VERSION = "1";
    private static final String CONTRACT_ARGS_EVM_PARAM = "data";

    @Getter
    private String contractName;
    private ChainmakerChainClient chainmakerChainClient;
    private ChainmakerAccountUtil chainmakerAccountUtil;

    public ChainmakerContract(String contractName, ChainmakerChainClient chainmakerChainClient, ChainmakerAccountUtil chainmakerAccountUtil) {
        this.contractName = contractName;
        this.chainmakerChainClient = chainmakerChainClient;
        this.chainmakerAccountUtil = chainmakerAccountUtil;
    }

    protected static DeployInfo deployContract(ChainmakerChainClient chainmakerChainClient, String contractName,
                                               Function function, String version, byte[] byteCode) throws Exception {
        String encodedConstructor = FunctionEncoder.encode(function);
        ImmutableMap<String, byte[]> paramsMap = ImmutableMap.of(
                CONTRACT_ARGS_EVM_PARAM, encodedConstructor.substring(10).getBytes()
        );
        ChainClient chainClient = chainmakerChainClient.getChainClient();
        if (StringUtils.isBlank(version)) {
            version = CONTRACT_DEFAULT_VERSION;
        }
        Request.Payload payload = chainClient.createContractCreatePayload(contractName, version, byteCode, ContractOuterClass.RuntimeType.EVM, paramsMap);
        Request.EndorsementEntry[] endorsers = SdkUtils.getEndorsers(payload, new User[]{chainClient.getClientUser()}, chainClient.getHash());
        ResultOuterClass.TxResponse txResponse = chainClient.sendContractManageRequest(payload, endorsers,
                ChainmakerChainClient.REQUEST_TIMEOUT, ChainmakerChainClient.REQUEST_TIMEOUT);
        if (txResponse.getCode() != ResultOuterClass.TxStatusCode.SUCCESS) {
            throw new RuntimeException("Deploy Contract Request Failed, result = " + txResponse);
        }
        ResultOuterClass.ContractResult contractResult = txResponse.getContractResult();
        if (contractResult.getCode() != 0) {
            throw new RuntimeException("Deploy Contract Failed, result = " + contractResult);
        }
        ContractOuterClass.Contract contract = ContractOuterClass.
                Contract.newBuilder().
                mergeFrom(txResponse.getContractResult().getResult().toByteArray()).
                build();

        DeployInfo deployInfo = new DeployInfo();
        deployInfo.setTxResponse(txResponse);
        deployInfo.setContract(contract);

        /**
         * 部署后实际上获取不到块高，只有交易哈希
         *
         */
        ChainmakerTransaction.TransactionInfo txInfo = chainClient.getTxByTxId(txResponse.getTxId(), ChainmakerChainClient.REQUEST_TIMEOUT);
        TransactionReceipt transactionReceipt = new TransactionReceipt();
        transactionReceipt.setBlockNumber(String.valueOf(txInfo.getBlockHeight()));
        transactionReceipt.setTransactionHash(txResponse.getTxId());
        deployInfo.setTransactionReceipt(transactionReceipt);
        deployInfo.setTransactionInfo(txInfo);
        return deployInfo;
    }

    /**
     * 执行合约返回数据
     *
     * @param function
     * @param tClass
     * @param <T>
     * @return
     */
    protected <T> T executeRemoteCallSingleValueReturn(Function function, Class<T> tClass) {
        String methodDataStr = FunctionEncoder.encode(function);
        ImmutableMap<String, byte[]> paramsMap = ImmutableMap.of(
                CONTRACT_ARGS_EVM_PARAM, methodDataStr.getBytes()
        );
        String method = methodDataStr.substring(0, 10);
        ResultOuterClass.TxResponse txResponse = null;
        try {
            ChainClient chainClient = chainmakerChainClient.getChainClient();
            txResponse = chainClient.queryContract(contractName, method, null, paramsMap, ChainmakerChainClient.REQUEST_TIMEOUT);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        if (txResponse.getCode() != ResultOuterClass.TxStatusCode.SUCCESS) {
            throw new RuntimeException("Send Query Contract Request Failed, result = " + txResponse);
        }
        ResultOuterClass.ContractResult contractResult = txResponse.getContractResult();
        if (contractResult.getCode() != 0) {
            throw new RuntimeException("Query Contract Failed, result = " + contractResult);
        }
        ByteString result = contractResult.getResult();
        if (tClass.equals(Address.class)) {
            return (T) new Address(Numeric.toHexString(result.toByteArray()));
        }
        if (tClass.equals(BigInteger.class)) {
            return (T) Numeric.toBigInt(result.toByteArray());
        }
        //Boolean
        if (tClass.equals(Boolean.class)) {
            byte[] bytes = result.toByteArray();
            // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0] or
            // [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1]
            return (T) Boolean.valueOf(bytes[31] == 1);
            //return (T) Boolean.valueOf(result.toByteArray());
        }
//        ByteString byteString = result.substring(64);
        String resultStr = ChainResultParser.parseResult(result, String.class);
        if (tClass.equals(String.class)) {
            return (T) resultStr;
        }
        return Convert.convert(tClass, resultStr);
    }

    /**
     * 异步调用执行合约
     *
     * @param function
     * @return txId
     */
    protected String executeRemoteCallTransactionAsync(Function function) {
        String methodDataStr = FunctionEncoder.encode(function);
        ImmutableMap<String, byte[]> paramsMap = ImmutableMap.of(
                CONTRACT_ARGS_EVM_PARAM, methodDataStr.getBytes()
        );
        String method = methodDataStr.substring(0, 10);
        ResultOuterClass.TxResponse txResponse = null;
        String txId = Utils.generateTxId();
        try {
            ChainClient chainClient = chainmakerChainClient.getChainClient();
            txResponse = chainClient.invokeContract(contractName, method, txId, paramsMap, ChainmakerChainClient.REQUEST_TIMEOUT, -1);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        if (txResponse.getCode() != ResultOuterClass.TxStatusCode.SUCCESS) {
            throw new ChangAnTxException(txResponse);
        }
        return txId;
    }

    /**
     * 调用执行合约
     *
     * @param function
     * @return
     */
    protected TransactionReceipt executeRemoteCallTransaction(Function function) {
        String methodDataStr = FunctionEncoder.encode(function);
        ImmutableMap<String, byte[]> paramsMap = ImmutableMap.of(
                CONTRACT_ARGS_EVM_PARAM, methodDataStr.getBytes()
        );
        String method = methodDataStr.substring(0, 10);
        ResultOuterClass.TxResponse txResponse = null;
        try {
            ChainClient chainClient = chainmakerChainClient.getChainClient();
            String txId = Utils.generateTxId();
            log.info("executeRemoteCallTransaction txId:{}", txId);
            txResponse = chainClient.invokeContract(contractName, method, txId, paramsMap,
                    ChainmakerChainClient.REQUEST_TIMEOUT, ChainmakerChainClient.REQUEST_TIMEOUT);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        // 判断是否成功
        if (txResponse.getCode() != ResultOuterClass.TxStatusCode.SUCCESS) {
            throw new ChangAnTxException(txResponse);
        }

        // 获取返回值
        ResultOuterClass.ContractResult contractResult = txResponse.getContractResult();
        if (contractResult.getCode() != 0) {
            ByteString messageBytes = contractResult.getMessageBytes();

            // 解析返回值
            String message = EthDecoderUtil.decodeOneString(messageBytes.toStringUtf8(), true);

            throw new RuntimeException("Invoker Contract Failed, code = " + contractResult.getCode() + ", message: " + message);
        }

        // 参数
        String txId = txResponse.getTxId();
        ChainClient chainClient = chainmakerChainClient.getChainClient();
        String address = chainmakerAccountUtil.accountAddress(chainClient.getClientUser());

        TransactionReceipt transactionReceipt = new TransactionReceipt();
        transactionReceipt.setTransactionHash(txId);
        // 获取调用合约的账户地址
        transactionReceipt.setFrom(address);
        transactionReceipt.setTo(contractName);
        transactionReceipt.setContractAddress(contractName);
        transactionReceipt.setGasUsed(String.valueOf(contractResult.getGasUsed()));

        log.info("executeRemoteCallTransaction transactionHash:{}, transactionReceipt:{}", transactionReceipt.getTransactionHash(), transactionReceipt);
        return transactionReceipt;
    }

    /**
     * 指定账户方式调用
     *
     * @param function
     * @param user
     * @return
     */
    protected TransactionReceipt executeRemoteCallTransactionWithUser(Function function, User user) {
        String methodDataStr = FunctionEncoder.encode(function);
        ImmutableMap<String, byte[]> paramsMap = ImmutableMap.of(
                CONTRACT_ARGS_EVM_PARAM, methodDataStr.getBytes()
        );
        String method = methodDataStr.substring(0, 10);
        ResultOuterClass.TxResponse txResponse = null;
        try {
            ChainClient chainClient = chainmakerChainClient.getChainClient();
            Request.Payload payload = chainClient.invokeContractPayload(
                    contractName, method, null, paramsMap);
            Request.EndorsementEntry[] endorsers = SdkUtils.getEndorsers(payload,
                    new User[]{user}, chainClient.getHash());
            log.info("send request endorsers: {}", endorsers);
            txResponse = chainClient.sendContractRequest(payload, endorsers,
                    ChainmakerChainClient.REQUEST_TIMEOUT, ChainmakerChainClient.REQUEST_TIMEOUT, user);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        if (txResponse.getCode() != ResultOuterClass.TxStatusCode.SUCCESS) {
            throw new RuntimeException("Send Invoker Contract Request Failed, result = " + txResponse);
        }
        ResultOuterClass.ContractResult contractResult = txResponse.getContractResult();
        if (contractResult.getCode() != 0) {
            ByteString messageBytes = contractResult.getMessageBytes();//.substring(64 * 2 + 8);
//            String message = ChainResultParser.parseResult(messageBytes, String.class);
            String message = new String(Numeric.hexStringToByteArray(messageBytes.toStringUtf8()));
            throw new RuntimeException("Invoker Contract Failed, code = " + contractResult.getCode() + ", message: " + message);
        }

        // 参数
        String txId = txResponse.getTxId();
        String address = chainmakerAccountUtil.accountAddress(user);

        TransactionReceipt transactionReceipt = new TransactionReceipt();
        transactionReceipt.setTransactionHash(txId);
        // 获取调用合约的账户地址
        transactionReceipt.setFrom(address);
        transactionReceipt.setTo(contractName);
        transactionReceipt.setContractAddress(contractName);
        transactionReceipt.setGasUsed(String.valueOf(contractResult.getGasUsed()));
        return transactionReceipt;
    }

    /**
     * 部署合约的信息
     */
    @Data
    public static class DeployInfo {
        private ResultOuterClass.TxResponse txResponse;
        private ContractOuterClass.Contract contract;
        private TransactionReceipt transactionReceipt;
        private ChainmakerTransaction.TransactionInfo transactionInfo;
    }
}
