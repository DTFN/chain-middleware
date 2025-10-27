package com.lingshu.server.core.lingshu.contract;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.codec.datatypes.Event;
import com.lingshu.chain.sdk.codec.datatypes.Function;
import com.lingshu.chain.sdk.codec.datatypes.Type;
import com.lingshu.chain.sdk.codec.datatypes.TypeReference;
import com.lingshu.chain.sdk.codec.datatypes.Utf8String;
import com.lingshu.chain.sdk.codec.datatypes.generated.tuples.generated.Tuple2;
import com.lingshu.chain.sdk.contract.Contract;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.crypto.key.CryptoKeyPair;
import com.lingshu.chain.sdk.evtsub.EvtSubCallback;
import com.lingshu.chain.sdk.model.CryptoType;
import com.lingshu.chain.sdk.model.TransactionCallback;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.tx.common.VmTypeEnum;
import com.lingshu.chain.sdk.tx.common.exception.ContractException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

@SuppressWarnings("unchecked")
public class BusinessContract extends Contract {
    public static final String FUNC_SAVE_BUSINESS_DETAILS = "saveBusinessDetails";

    public static final String FUNC_GET_BUSINESS_DETAILS = "getBusinessDetails";

    public static final String[] ABI_ARRAY = {"[{\"constant\":false,\"inputs\":[{\"name\":\"businessAddress\",\"type\":\"string\"},{\"name\":\"details\",\"type\":\"string\"}],\"name\":\"saveBusinessDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"businessAddress\",\"type\":\"string\"}],\"name\":\"getBusinessDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"businessAddress\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"details\",\"type\":\"string\"}],\"name\":\"BusinessRegistered\",\"type\":\"event\"}]"};

    public static final String ABI = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", ABI_ARRAY);

    public static final String[] BINARY_ARRAY = {"608060405234801561001057600080fd5b506104aa806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630f57eaf414610051578063a2681585146100a4575b600080fd5b34801561005d57600080fd5b506100a2600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610158565b005b3480156100b057600080fd5b506100dd600480360381019080803590602001908201803590602001919091929391929390505050610311565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561011d578082015181840152602081019050610102565b50505050905090810190601f16801561014a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16151515610226576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001807f427573696e657373436f6e74726163743a20627573696e65737320616c72656181526020017f647920657869737473000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b818160008686604051808383808284378201915050925050509081526020016040518091039020919061025a9291906103d9565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507f04f0a19b3458ac122557e5494f593ccf09abb67f5b8fd5ba741f795eeb547b01848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103cc5780601f106103a1576101008083540402835291602001916103cc565b820191906000526020600020905b8154815290600101906020018083116103af57829003601f168201915b5050505050905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061041a57803560ff1916838001178555610448565b82800160010185558215610448579182015b8281111561044757823582559160200191906001019061042c565b5b5090506104559190610459565b5090565b61047b91905b8082111561047757600081600090555060010161045f565b5090565b905600a165627a7a723058204a245c96a0cab76cb11d9d9a105cb6861eea75c800b256840934b89d3ea4c6440029"};

    public static final String BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", BINARY_ARRAY);

    public static final String[] SM_BINARY_ARRAY = {"608060405234801561001057600080fd5b506104aa806100206000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063123105221461005157806394f22242146100a4575b600080fd5b34801561005d57600080fd5b506100a2600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610158565b005b3480156100b057600080fd5b506100dd600480360381019080803590602001908201803590602001919091929391929390505050610311565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561011d578082015181840152602081019050610102565b50505050905090810190601f16801561014a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16151515610226576040517fc703cb120000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001807f427573696e657373436f6e74726163743a20627573696e65737320616c72656181526020017f647920657869737473000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b818160008686604051808383808284378201915050925050509081526020016040518091039020919061025a9291906103d9565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507f78a718e416a2e932e844ac3c2b187091d5b790ec5e3ba57228ad375d2caa9b28848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103cc5780601f106103a1576101008083540402835291602001916103cc565b820191906000526020600020905b8154815290600101906020018083116103af57829003601f168201915b5050505050905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061041a57803560ff1916838001178555610448565b82800160010185558215610448579182015b8281111561044757823582559160200191906001019061042c565b5b5090506104559190610459565b5090565b61047b91905b8082111561047757600081600090555060010161045f565b5090565b905600a165627a7a723058206f65c215294c0d889a66fab028f94c250d21d0813fcdeec5326e9df0cb80cc960029"};

    public static final String SM_BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", SM_BINARY_ARRAY);

    public static final Event BUSINESS_REGISTERED_EVENT = new Event("BusinessRegistered",
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
    ;

    protected BusinessContract(String contractAddress, IClient client, CryptoKeyPair credential) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential);
    }

    protected BusinessContract(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential, vmType);
    }

    protected BusinessContract(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        super(binary, contractAddress, client, credential);
    }

    protected BusinessContract(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        super(binary, contractAddress, client, credential, vmType);
    }

    public static BusinessContract deploy(IClient client, CryptoKeyPair credential) throws ContractException {
        return deploy(BusinessContract.class, client, credential, getBinary(client.getCryptoSuite()), "");
    }

    public static BusinessContract deploy(IClient client, CryptoKeyPair credential, String binary) throws ContractException {
        return deploy(BusinessContract.class, client, credential, binary, "");
    }

    public static BusinessContract deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType) throws ContractException {
        return deploy(BusinessContract.class, client, credential, getBinary(client.getCryptoSuite()), "", vmType);
    }

    public static BusinessContract deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) throws ContractException {
        return deploy(BusinessContract.class, client, credential, binary, "", vmType);
    }

    public static BusinessContract load(String contractAddress, IClient client, CryptoKeyPair credential) {
        return new BusinessContract(contractAddress, client, credential);
    }

    public static BusinessContract load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        return new BusinessContract(contractAddress, client, credential, vmType);
    }

    public static BusinessContract load(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        return new BusinessContract(contractAddress, client, credential, binary);
    }

    public static BusinessContract load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        return new BusinessContract(contractAddress, client, credential, vmType, binary);
    }

    public TransactionReceipt saveBusinessDetails(String businessAddress, String details) {
        final Function function = new Function(
                FUNC_SAVE_BUSINESS_DETAILS,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(businessAddress),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] saveBusinessDetails(String businessAddress, String details, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_SAVE_BUSINESS_DETAILS,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(businessAddress),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForSaveBusinessDetails(String businessAddress, String details) {
        final Function function = new Function(
                FUNC_SAVE_BUSINESS_DETAILS,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(businessAddress),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getSaveBusinessDetailsInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_SAVE_BUSINESS_DETAILS,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(),
                (String) results.get(1).getValue()
        );
    }

    public String getBusinessDetails(String businessAddress) throws ContractException {
        final Function function = new Function(FUNC_GET_BUSINESS_DETAILS,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(businessAddress)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeCallWithSingleValueReturn(function, String.class);
    }

    public List<BusinessRegisteredEvtResp> getBusinessRegisteredEvents(TransactionReceipt txReceipt) {
        List<Contract.EvtValuesWithLog> valueList = extractEventParametersWithLog(BUSINESS_REGISTERED_EVENT, txReceipt);
        ArrayList<BusinessRegisteredEvtResp> responseList = new ArrayList<BusinessRegisteredEvtResp>(valueList.size());
        for (Contract.EvtValuesWithLog eventValues : valueList) {
            BusinessRegisteredEvtResp evtResp = new BusinessRegisteredEvtResp();
            evtResp.log = eventValues.getLog();
            evtResp.businessAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
            evtResp.details = (String) eventValues.getNonIndexedValues().get(1).getValue();
            responseList.add(evtResp);
        }
        return responseList;
    }

    public void subscribeBusinessRegisteredEvent(String fromBlock, String toBlock, List<String> otherTopics, EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(BUSINESS_REGISTERED_EVENT);
        subscribeEvent(ABI,BINARY,topic0,fromBlock,toBlock,otherTopics,callback);
    }

    public void subscribeBusinessRegisteredEvent(EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(BUSINESS_REGISTERED_EVENT);
        subscribeEvent(ABI,BINARY,topic0,callback);
    }

    public static String getBinary(CryptoSuite cryptoSuite) {
        return (cryptoSuite.getCryptoType() == CryptoType.ECDSA_TYPE ? BINARY : SM_BINARY);
    }

    public static class BusinessRegisteredEvtResp {
        public TransactionReceipt.Logs log;

        public String businessAddress;

        public String details;
    }
}
