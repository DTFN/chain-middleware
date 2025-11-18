package com.lingshu.server.core.lingshu.contract;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.codec.datatypes.*;
import com.lingshu.chain.sdk.contract.Contract;
import com.lingshu.chain.sdk.crypto.key.CryptoKeyPair;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.tx.common.VmTypeEnum;
import com.lingshu.chain.sdk.tx.common.exception.ContractException;

import java.io.IOException;
import java.util.Arrays;
import java.util.Collections;

@SuppressWarnings("unchecked")
public class LsBusinessContract2 extends Contract {

    protected LsBusinessContract2(String contractAddress, IClient client, CryptoKeyPair credential) {
        super(/*getBinary(client.getCryptoSuite())*/ null, contractAddress, client, credential);
    }

    protected LsBusinessContract2(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        super(/*getBinary(client.getCryptoSuite())*/ null, contractAddress, client, credential, vmType);
    }

    protected LsBusinessContract2(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        super(binary, contractAddress, client, credential);
    }

    protected LsBusinessContract2(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        super(binary, contractAddress, client, credential, vmType);
    }

    public static LsBusinessContract2 deploy(IClient client, CryptoKeyPair credential) throws ContractException {
        return deploy(LsBusinessContract2.class, client, credential, null, "");
    }

    public static LsBusinessContract2 deploy(IClient client, CryptoKeyPair credential, String binary) throws ContractException {
        return deploy(LsBusinessContract2.class, client, credential, binary, "");
    }

    public static LsBusinessContract2 deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType) throws ContractException {
        return deploy(LsBusinessContract2.class, client, credential, null, "", vmType);
    }

    public static LsBusinessContract2 deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) throws ContractException {
        return deploy(LsBusinessContract2.class, client, credential, binary, "", vmType);
    }

    public static LsBusinessContract2 load(String contractAddress, IClient client, CryptoKeyPair credential) {
        return new LsBusinessContract2(contractAddress, client, credential);
    }

    public static LsBusinessContract2 load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        return new LsBusinessContract2(contractAddress, client, credential, vmType);
    }

    public static LsBusinessContract2 load(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        return new LsBusinessContract2(contractAddress, client, credential, binary);
    }

    public static LsBusinessContract2 load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        return new LsBusinessContract2(contractAddress, client, credential, vmType, binary);
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

    public String get(String methodName, String paramValue0) throws IOException {
        final Function function = new Function(methodName,
                Arrays.<Type>asList(new Utf8String(paramValue0)),
                Arrays.asList(new TypeReference<Utf8String>() {
                })
        );
        return createSignedTransaction(function);
    }

    public TransactionReceipt getResult(String methodName, String vcs) {
        final Function function = new Function(
                methodName,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(vcs)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public static class BusinessRegisteredEvtResp {
        public TransactionReceipt.Logs log;

        public String businessAddress;

        public String details;
    }
}
