package com.lingshu.bsp.front.dto;

import java.util.List;


public class TransactionDetail {
//
// Source code recreated from a .class file by IntelliJ IDEA
// (powered by FernFlower decompiler)
//

    private String blockHash;
    private String blockNumber;
    private String from;
    private String gas;
    private String hash;
    private String input;
    private String nonce;
    private String to;
    private String transactionIndex;
    private String value;
    private String gasPrice;
    private String blockLimit;
    private String chainId;
    private String ledgerId;
    private String extraData;
    private com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.SignatureResponse signature;
    private List<com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.TxProof> txProof;


    public String getBlockNumber() {
        return blockNumber;
    }

    public void setBlockHash(String blockHash) {
        this.blockHash = blockHash;
    }

    public void setBlockNumber(String blockNumber) {
        this.blockNumber = blockNumber;
    }

    public void setFrom(String from) {
        this.from = from;
    }

    public void setGas(String gas) {
        this.gas = gas;
    }

    public void setHash(String hash) {
        this.hash = hash;
    }

    public void setInput(String input) {
        this.input = input;
    }

    public void setNonce(String nonce) {
        this.nonce = nonce;
    }

    public void setTo(String to) {
        this.to = to;
    }

    public void setTransactionIndex(String transactionIndex) {
        this.transactionIndex = transactionIndex;
    }

    public void setValue(String value) {
        this.value = value;
    }

    public void setGasPrice(String gasPrice) {
        this.gasPrice = gasPrice;
    }

    public void setBlockLimit(String blockLimit) {
        this.blockLimit = blockLimit;
    }

    public void setChainId(String chainId) {
        this.chainId = chainId;
    }

    public void setLedgerId(String ledgerId) {
        this.ledgerId = ledgerId;
    }

    public void setExtraData(String extraData) {
        this.extraData = extraData;
    }

    public void setSignature(com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.SignatureResponse signature) {
        this.signature = signature;
    }

    public void setTxProof(List<com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.TxProof> txProof) {
        this.txProof = txProof;
    }

    public String getBlockHash() {
        return this.blockHash;
    }

    public String getFrom() {
        return this.from;
    }

    public String getGas() {
        return this.gas;
    }

    public String getHash() {
        return this.hash;
    }

    public String getInput() {
        return this.input;
    }

    public String getNonce() {
        return this.nonce;
    }

    public String getTo() {
        return this.to;
    }

    public String getTransactionIndex() {
        return this.transactionIndex;
    }

    public String getValue() {
        return this.value;
    }

    public String getGasPrice() {
        return this.gasPrice;
    }

    public String getBlockLimit() {
        return this.blockLimit;
    }

    public String getChainId() {
        return this.chainId;
    }

    public String getLedgerId() {
        return this.ledgerId;
    }

    public String getExtraData() {
        return this.extraData;
    }

    public com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.SignatureResponse getSignature() {
        return this.signature;
    }

    public List<com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse.TxProof> getTxProof() {
        return this.txProof;
    }

    public static class SignatureResponse {
        private String r;
        private String s;
        private String v;
        private String signature;

        public SignatureResponse() {
        }

        public void setR(String r) {
            this.r = r;
        }

        public void setS(String s) {
            this.s = s;
        }

        public void setV(String v) {
            this.v = v;
        }

        public void setSignature(String signature) {
            this.signature = signature;
        }

        public String getR() {
            return this.r;
        }

        public String getS() {
            return this.s;
        }

        public String getV() {
            return this.v;
        }

        public String getSignature() {
            return this.signature;
        }
    }

    public static class TxProof {
        private List<String> left;
        private List<String> right;

        public TxProof() {
        }

        public void setLeft(List<String> left) {
            this.left = left;
        }

        public void setRight(List<String> right) {
            this.right = right;
        }

        public List<String> getLeft() {
            return this.left;
        }

        public List<String> getRight() {
            return this.right;
        }
    }
}

