package com.lingshu.bsp.front.rpcapi.bsc;

import cn.hutool.core.bean.BeanUtil;
import cn.hutool.core.collection.CollectionUtil;
import cn.hutool.json.JSONUtil;
import com.lingshu.bsp.front.rpcapi.bsc.contract.BusiCenterBsc;
import com.lingshu.bsp.front.rpcapi.bsc.contract.DIDManagerBsc;
import com.lingshu.bsp.front.rpcapi.eth.contract.CrossSaveEth;
import com.lingshu.chain.sdk.client.protocol.model.JsonTransactionResponse;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import lombok.extern.slf4j.Slf4j;
import okhttp3.OkHttpClient;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.web3j.crypto.CipherException;
import org.web3j.crypto.Credentials;
import org.web3j.crypto.WalletUtils;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.protocol.core.methods.response.Transaction;
import org.web3j.protocol.http.HttpService;
import org.web3j.tx.RawTransactionManager;
import org.web3j.tx.gas.StaticGasProvider;
import org.web3j.tx.response.PollingTransactionReceiptProcessor;
import org.web3j.tx.response.TransactionReceiptProcessor;

import javax.annotation.PostConstruct;
import java.io.IOException;
import java.math.BigInteger;
import java.security.Security;
import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * @author gongrui.wang
 * @since 2025/8/27
 */
@Slf4j
@Service
public class BscApiService {

    private static Web3j web3j;
    private static Credentials defaultCredentials = null;
    private static RawTransactionManager rawTransactionManager;

    @Value("${bscEndpoint:http://192.168.1.144:8645}")
    private String bscEndpoint;

    static {
        Security.addProvider(new BouncyCastleProvider());
    }

    @PostConstruct
    public void init() throws CipherException, IOException {
        log.info("bsc url: {}", bscEndpoint);

        // 加载客户端
        OkHttpClient client = HttpService.getOkHttpClientBuilder()
                .connectTimeout(3, TimeUnit.MINUTES)
                .readTimeout(3, TimeUnit.MINUTES)
                .writeTimeout(3, TimeUnit.MINUTES)
                .build();
        web3j = Web3j.build(new HttpService(bscEndpoint, client));

        // 加载钱包
        String walletFile = "{\"address\":\"b185d579023ad5584907a46c51fd22d244b8f549\",\"crypto\":{\"cipher\":\"aes-128-ctr\",\"ciphertext\":\"65974eb531b945a39345bbcdad263fbc1d99c1534514fc73f8aeeeb41fa3e252\",\"cipherparams\":{\"iv\":\"bfb5a3c705adb08ab9b0cad7bef3a70f\"},\"kdf\":\"scrypt\",\"kdfparams\":{\"dklen\":32,\"n\":262144,\"p\":1,\"r\":8,\"salt\":\"252e0b3ead0cb9de0e03d65f3b9332e97069162aa2ee0999255ce0fa366e1015\"},\"mac\":\"ecb46ddbf80a28d1b6f3adff5b241b789753b311a1d1e2e2fab1f786872e3eb4\"},\"id\":\"87475d1b-6bed-4cbf-afdc-77596116e3cf\",\"version\":3}";
        defaultCredentials = WalletUtils.loadJsonCredentials("123456789", walletFile);
        log.info("private key, {}", defaultCredentials.getEcKeyPair().getPrivateKey().toString(16));

        // 设置钱包
        TransactionReceiptProcessor receiptProcessor =
                new PollingTransactionReceiptProcessor(web3j, 50, 10000);
        rawTransactionManager = new RawTransactionManager(
                web3j,
                defaultCredentials,
                888,
                receiptProcessor
        );
    }

    public Web3j web3j() {
        return web3j;
    }

    public RawTransactionManager transactionManager() {
        return rawTransactionManager;
    }


    public JsonTransactionResponse getTransactionByHash(String transHash) {
        try {
            Transaction transaction = web3j.ethGetTransactionByHash(transHash).send().getResult();
            JsonTransactionResponse response = BeanUtil.copyProperties(transaction, JsonTransactionResponse.class, "blockNumber", "transactionIndex");
            response.setBlockNumber(transaction.getBlockNumberRaw());
            response.setTransactionIndex(transaction.getTransactionIndexRaw());
            return response;
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    public TransactionReceipt crossSave(String contractAddress) {
        CrossSaveEth crossSaveEth = CrossSaveEth.load(contractAddress, web3j,
                transactionManager(), new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        try {
            org.web3j.protocol.core.methods.response.TransactionReceipt receipt = crossSaveEth.CrossChainSave("hello", "hello").send();
            TransactionReceipt result = BeanUtil.copyProperties(receipt, TransactionReceipt.class);
            result.setStatus(receipt.isStatusOK() ? 0 : 1);
            return result;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }

    public BigInteger getBlockNumber() {
        try {
            return web3j.ethBlockNumber().send().getBlockNumber();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }

    public String deploy() {
        RemoteCall<CrossSaveEth> deploy = CrossSaveEth.deploy(web3j,
                transactionManager(),
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        String contractAddress = null;
        try {
            contractAddress = deploy.send().getContractAddress();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy address: {}", contractAddress);
        return contractAddress;
    }

    public Web3j getWeb3j() {
        return web3j;
    }

    public RawTransactionManager getRawTransactionManager() {
        return rawTransactionManager;
    }

    public String deployDidManager() {
        // check ledgerId
        RemoteCall<DIDManagerBsc> deploy = DIDManagerBsc.deploy(web3j,
                transactionManager(),
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        String contractAddress = null;
        try {
            contractAddress = deploy.send().getContractAddress();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy address: {}", contractAddress);
        return contractAddress;
    }

    public String deployBusiCenter() {
        // check ledgerId
        RemoteCall<BusiCenterBsc> deploy = BusiCenterBsc.deploy(web3j,
                transactionManager(),
                new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        String contractAddress = null;
        try {
            contractAddress = deploy.send().getContractAddress();
        } catch (Exception e) {
            e.printStackTrace();
        }
        log.info("deploy address: {}", contractAddress);
        return contractAddress;
    }
    public TransactionReceipt invokeBusiCenter(String address, List<Object> args) {
        // check ledgerId
        if (CollectionUtil.isEmpty(args)) {
            throw new RuntimeException("args are empty.");
        }
        BusiCenterBsc contract = BusiCenterBsc.load(address, web3j,
                transactionManager(), new StaticGasProvider(BigInteger.valueOf(20_000_000_000L), BigInteger.valueOf(2_800_000)));
        try {
            org.web3j.protocol.core.methods.response.TransactionReceipt receipt = contract.erc20MintVcs(JSONUtil.toJsonStr(args.get(0))).send();
            TransactionReceipt result = BeanUtil.copyProperties(receipt, TransactionReceipt.class);
            result.setStatus(receipt.isStatusOK() ? 0 : 1);
            return result;
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }
}
