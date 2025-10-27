package com.lingshu.server.core.web3j.contract;

import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil4Busi;
import com.lingshu.server.core.web3j.core.TransactionReceiptExt;
import lombok.Getter;
import org.chainmaker.sdk.ChainClient;
import org.chainmaker.sdk.utils.FileUtils;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.Arrays;
import java.util.Collections;
import java.util.Map;

/**
 * @author: derrick
 * @since: 2025-08-26
 */
public class BusinessContract4Wasm extends ChainmakerContract4Wasm {
    //    public static final String CONTRACT_NAME = "BUSINESS_CONTRACT_WASM";
//    public static final String BINARY = "";
//    public static final String FUNC_GET_BUSINESS_DETAILS = "getBusinessDetails";
    @Getter
    private DeployInfo deployInfo;

    public BusinessContract4Wasm(String contractName,
                                 ChainmakerChainClient chainmakerChainClient,
                                 ChainmakerAccountUtil4Busi chainmakerAccountUtil) {
        super(contractName, chainmakerChainClient, chainmakerAccountUtil);
    }

    public static BusinessContract4Wasm load(String contractName,
                                             ChainmakerChainClient chainmakerChainClient,
                                             ChainmakerAccountUtil4Busi chainmakerAccountUtil) {
        return new BusinessContract4Wasm(contractName, chainmakerChainClient, chainmakerAccountUtil);
    }

    public static BusinessContract4Wasm deploy(ChainmakerAccountUtil4Busi chainmakerAccountUtil,
                                               ChainmakerChainClient chainmakerChainClient,
                                               String contractName, String version, String filePath) throws Exception {
        Function function = new Function("",
//                Arrays.<Type>asList(new Utf8String(name), new Utf8String(symbol)),
                Collections.emptyList(),
                Collections.emptyList());
        byte[] byteCode = FileUtils.getFileBytes(filePath);
        DeployInfo deployInfo = ChainmakerContract4Wasm.deployContract(chainmakerChainClient, contractName, function, version, byteCode);
        BusinessContract4Wasm businessContract = load(contractName, chainmakerChainClient, chainmakerAccountUtil);
        businessContract.deployInfo = deployInfo;
        return businessContract;
    }


    public TransactionReceipt call(String methodName, Map<String, byte[]> paramsMap) /*throws ContractException*/ {
        return executeRemoteCallTransaction(methodName, paramsMap);
    }

    public String get(String methodName, Map<String, byte[]> paramsMap) /*throws ContractException*/ {
        return executeRemoteCallSingleValueReturn(methodName, paramsMap, String.class);
    }
}
