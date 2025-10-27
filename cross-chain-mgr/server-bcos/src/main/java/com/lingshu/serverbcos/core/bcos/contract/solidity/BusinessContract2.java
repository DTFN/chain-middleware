package com.lingshu.serverbcos.core.bcos.contract.solidity;

import org.fisco.bcos.sdk.abi.TypeReference;
import org.fisco.bcos.sdk.abi.datatypes.Function;
import org.fisco.bcos.sdk.abi.datatypes.Type;
import org.fisco.bcos.sdk.abi.datatypes.Utf8String;
import org.fisco.bcos.sdk.client.Client;
import org.fisco.bcos.sdk.contract.Contract;
import org.fisco.bcos.sdk.crypto.CryptoSuite;
import org.fisco.bcos.sdk.crypto.keypair.CryptoKeyPair;
import org.fisco.bcos.sdk.model.CryptoType;
import org.fisco.bcos.sdk.model.TransactionReceipt;
import org.fisco.bcos.sdk.transaction.model.exception.ContractException;
import org.web3j.protocol.core.RemoteFunctionCall;

import java.io.File;
import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;

@SuppressWarnings("unchecked")
public class BusinessContract2 extends Contract {
    protected BusinessContract2(String binary, String contractAddress, Client client, CryptoKeyPair credential) {
        super(binary, contractAddress, client, credential);
    }

//    public static String getBinary(CryptoSuite cryptoSuite) {
////        return (cryptoSuite.getCryptoTypeConfig() == CryptoType.ECDSA_TYPE ? BINARY : SM_BINARY);
//        return (cryptoSuite.getCryptoTypeConfig() == CryptoType.ECDSA_TYPE ? BINARY : BINARY);
//    }

    public static BusinessContract2 deploy(Client client, CryptoKeyPair credential, String binary) throws ContractException {
//        String binary = Files.readString(new File(filePath));
//        log.info("binary: {}", binary);
        return deploy(BusinessContract2.class, client, credential, binary, "");
    }

    public static BusinessContract2 load(/*String binary,*/ String contractAddress, Client client, CryptoKeyPair credential) {
        return new BusinessContract2(/*binary*/null, contractAddress, client, credential);
    }

    public TransactionReceipt call(String methodName, String paramValue0) {
        final Function function = new Function(
                methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public TransactionReceipt call(String methodName, String paramValue0, String paramValue1) {
        final Function function = new Function(
                methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0), new Utf8String(paramValue1)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public TransactionReceipt call(Function function) {
        return executeTransaction(function);
    }

    public String get(String methodName, String paramValue0) throws ContractException {
        final Function function = new Function(methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0)),
                Arrays.asList(new TypeReference<Utf8String>() {
                })
        );
        return executeCallWithSingleValueReturn(function, String.class);
    }

    public String get(Function function) throws ContractException {
        return executeCallWithSingleValueReturn(function, String.class);
    }

}
