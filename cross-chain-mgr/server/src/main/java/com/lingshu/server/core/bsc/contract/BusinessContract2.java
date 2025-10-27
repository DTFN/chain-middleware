package com.lingshu.server.core.bsc.contract;

import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.RemoteFunctionCall;
import org.web3j.protocol.core.methods.response.BaseEventResponse;
import org.web3j.protocol.core.methods.response.TransactionReceipt;
import org.web3j.tx.Contract;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.gas.ContractGasProvider;

import java.io.IOException;
import java.math.BigInteger;
import java.util.Arrays;
import java.util.Collections;

@SuppressWarnings("unchecked")
public class BusinessContract2 extends Contract {

//    public static final String BINARY = "608060405234801561001057600080fd5b506104aa806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630f57eaf414610051578063a2681585146100a4575b600080fd5b34801561005d57600080fd5b506100a2600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610158565b005b3480156100b057600080fd5b506100dd600480360381019080803590602001908201803590602001919091929391929390505050610311565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561011d578082015181840152602081019050610102565b50505050905090810190601f16801561014a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16151515610226576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001807f427573696e657373436f6e74726163743a20627573696e65737320616c72656181526020017f647920657869737473000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b818160008686604051808383808284378201915050925050509081526020016040518091039020919061025a9291906103d9565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507f04f0a19b3458ac122557e5494f593ccf09abb67f5b8fd5ba741f795eeb547b01848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103cc5780601f106103a1576101008083540402835291602001916103cc565b820191906000526020600020905b8154815290600101906020018083116103af57829003601f168201915b5050505050905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061041a57803560ff1916838001178555610448565b82800160010185558215610448579182015b8281111561044757823582559160200191906001019061042c565b5b5090506104559190610459565b5090565b61047b91905b8082111561047757600081600090555060010161045f565b5090565b905600a165627a7a723058204a245c96a0cab76cb11d9d9a105cb6861eea75c800b256840934b89d3ea4c6440029";
//    public static final String FUNC_SAVE_BUSINESS_DETAILS = "saveBusinessDetails";

    public static final String FUNC_GET_BUSINESS_DETAILS = "getBusinessDetails";

//    public static final Event BUSINESS_REGISTERED_EVENT = new Event("BusinessRegistered",
//            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
//    ;

    protected BusinessContract2(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected BusinessContract2(String binary, String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        super(binary, contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    protected BusinessContract2(String binary, String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        super(binary, contractAddress, web3j, credentials, contractGasProvider);
    }

    protected BusinessContract2(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected BusinessContract2(String binary, String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        super(binary, contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

    protected BusinessContract2(String binary, String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        super(binary, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    protected BusinessContract2(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        super(contractAddress, web3j, transactionManager, contractGasProvider.getGasPrice(), contractGasProvider.getGasLimit());
    }

    // 存储字节码的实例变量
//    private final String binary;


    public static RemoteCall<BusinessContract2> deploy(String binary, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
        return deployRemoteCall(BusinessContract2.class, web3j, credentials, contractGasProvider, binary, "");
    }

    @Deprecated
    public static RemoteCall<BusinessContract2> deploy(String binary, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BusinessContract2.class, web3j, credentials, gasPrice, gasLimit, binary, "");
    }

    public static RemoteCall<BusinessContract2> deploy(String binary, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
        return deployRemoteCall(BusinessContract2.class, web3j, transactionManager, contractGasProvider, binary, "");
    }

    @Deprecated
    public static RemoteCall<BusinessContract2> deploy(String binary, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return deployRemoteCall(BusinessContract2.class, web3j, transactionManager, gasPrice, gasLimit, binary, "");
    }

    public static BusinessContract2 load(String contractAddress, Web3j web3j, Credentials credentials, BigInteger gasPrice, BigInteger gasLimit) {
        return new BusinessContract2(contractAddress, web3j, credentials, gasPrice, gasLimit);
    }

    public static BusinessContract2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, BigInteger gasPrice, BigInteger gasLimit) {
        return new BusinessContract2(contractAddress, web3j, transactionManager, gasPrice, gasLimit);
    }

//    public static BusinessContract2 load(String contractAddress, Web3j web3j, Credentials credentials, ContractGasProvider contractGasProvider) {
//        return new BusinessContract2(contractAddress, web3j, credentials, contractGasProvider);
//    }
//
//    public static BusinessContract2 load(String contractAddress, Web3j web3j, TransactionManager transactionManager, ContractGasProvider contractGasProvider) {
//        return new BusinessContract2(contractAddress, web3j, transactionManager, contractGasProvider);
//    }

//    public TransactionReceipt saveBusinessDetails(String businessAddress, String details) {
//        final Function function = new Function(
//                FUNC_SAVE_BUSINESS_DETAILS,
//                Arrays.<Type>asList(new Utf8String(businessAddress),
//                        new Utf8String(details)),
//                Collections.<TypeReference<?>>emptyList());
//        return executeTransaction(function);
//    }

    public RemoteFunctionCall<TransactionReceipt> call(String methodName, String paramValue0) {
        final Function function = new Function(
                methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public RemoteFunctionCall<TransactionReceipt> call(String methodName, String paramValue0, String paramValue1) {
        final Function function = new Function(
                methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0), new Utf8String(paramValue1)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

//    public RemoteFunctionCall<TransactionReceipt> saveBusinessDetails(String businessAddress, String details) {
//        final Function function = new Function(
//                FUNC_SAVE_BUSINESS_DETAILS,
//                Arrays.<Type>asList(new Utf8String(businessAddress), new Utf8String(details)),
//                Collections.<TypeReference<?>>emptyList());
//        return executeRemoteCallTransaction(function);
//    }
//
    public String getBusinessDetails(String businessAddress) throws IOException {
        final Function function = new Function(FUNC_GET_BUSINESS_DETAILS,
                Arrays.<Type>asList(new Utf8String(businessAddress)),
                Arrays.asList(new TypeReference<Utf8String>() {
                })
        );
        return executeCallSingleValueReturn(function, String.class);
    }
//
//    public String getBusinessDetails1(String businessAddress) throws IOException {
//        final Function function = new Function(FUNC_GET_BUSINESS_DETAILS,
//                Arrays.<Type>asList(new Utf8String(businessAddress)),
//                Arrays.asList(new TypeReference<Utf8String>() {
//                })
//        );
//        Utf8String result = executeCallSingleValueReturn(function);
//        return result.getValue();
//    }

//    public RemoteFunctionCall<TransactionReceipt> CrossChainConfirm(String key) {
//        final Function function = new Function(
//                FUNC_CROSSCHAINCONFIRM,
//                Arrays.<Type>asList(new org.web3j.abi.datatypes.Utf8String(key)),
//                Collections.<TypeReference<?>>emptyList());
//        return executeRemoteCallTransaction(function);
//    }

//    public static class TestEthEventResponse extends BaseEventResponse {
//        public byte[] key;
//
//        public byte[] value;
//    }
}
