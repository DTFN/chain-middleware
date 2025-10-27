package com.lingshu.chain.contract;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.codec.datatypes.Bool;
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
public class ResourceDomain extends Contract {
    public static final String FUNC_UPDATE_DOMAIN_DETAILS = "updateDomainDetails";

    public static final String FUNC_SAVE_DOMAIN_DETAILS = "saveDomainDetails";

    public static final String FUNC_GET_DOMAIN_DETAILS = "getDomainDetails";

    public static final String FUNC_DOES_DOMAIN_EXIST = "doesDomainExist";

    public static final String[] ABI_ARRAY = {"[{\"constant\":false,\"inputs\":[{\"name\":\"domainAddress\",\"type\":\"string\"},{\"name\":\"details\",\"type\":\"string\"}],\"name\":\"updateDomainDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"domainAddress\",\"type\":\"string\"},{\"name\":\"details\",\"type\":\"string\"}],\"name\":\"saveDomainDetails\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"domainAddress\",\"type\":\"string\"}],\"name\":\"getDomainDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"domainAddress\",\"type\":\"string\"}],\"name\":\"doesDomainExist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"domainAddress\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"details\",\"type\":\"string\"}],\"name\":\"DomainRegistered\",\"type\":\"event\"}]"};

    public static final String ABI = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", ABI_ARRAY);

    public static final String[] BINARY_ARRAY = {"608060405234801561001057600080fd5b5061068e806100206000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632e2cf984146100675780635574e0d2146100ba578063c38fb8361461010d578063d7b7e4cd146101c1575b600080fd5b34801561007357600080fd5b506100b8600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610214565b005b3480156100c657600080fd5b5061010b6004803603810190808035906020019082018035906020019190919293919293908035906020019082018035906020019190919293919293905050506102ff565b005b34801561011957600080fd5b506101466004803603810190808035906020019082018035906020019190919293919293905050506104b8565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561018657808201518184015260208101905061016b565b50505050905090810190601f1680156101b35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101cd57600080fd5b506101fa600480360381019080803590602001908201803590602001919091929391929390505050610580565b604051808215151515815260200191505060405180910390f35b81816000868660405180838380828437820191505092505050908152602001604051809103902091906102489291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507ff5e947bc58a45ca3f7c47f7a942b67608bb05338693c9dd564f093e315b8962e848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff161515156103cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001807f5265736f75726365446f6d61696e3a20646f6d61696e20616c7265616479206581526020017f786973747300000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b81816000868660405180838380828437820191505092505050908152602001604051809103902091906104019291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507ff5e947bc58a45ca3f7c47f7a942b67608bb05338693c9dd564f093e315b8962e848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105735780601f1061054857610100808354040283529160200191610573565b820191906000526020600020905b81548152906001019060200180831161055657829003601f168201915b5050505050905092915050565b60006001838360405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106105fe57803560ff191683800117855561062c565b8280016001018555821561062c579182015b8281111561062b578235825591602001919060010190610610565b5b509050610639919061063d565b5090565b61065f91905b8082111561065b576000816000905550600101610643565b5090565b905600a165627a7a72305820ef13e280805350eb3aca9dd10150f9950a27a89da1828d3a74214e76a7036dc10029"};

    public static final String BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", BINARY_ARRAY);

    public static final String[] SM_BINARY_ARRAY = {"608060405234801561001057600080fd5b5061068e806100206000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063781a16a614610067578063a59d36da146100ba578063d4d968951461010d578063f765f9ad146101c1575b600080fd5b34801561007357600080fd5b506100a0600480360381019080803590602001908201803590602001919091929391929390505050610214565b604051808215151515815260200191505060405180910390f35b3480156100c657600080fd5b5061010b600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610251565b005b34801561011957600080fd5b5061014660048036038101908080359060200190820180359060200191909192939192939050505061040a565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561018657808201518184015260208101905061016b565b50505050905090810190601f1680156101b35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101cd57600080fd5b506102126004803603810190808035906020019082018035906020019190919293919293908035906020019082018035906020019190919293919293905050506104d2565b005b60006001838360405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16905092915050565b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff1615151561031f576040517fc703cb120000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001807f5265736f75726365446f6d61696e3a20646f6d61696e20616c7265616479206581526020017f786973747300000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b81816000868660405180838380828437820191505092505050908152602001604051809103902091906103539291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507f894bacec2200a39552d3a7b98024096fdcfed79005b6f8ee1651b68c62bba2a8848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104c55780601f1061049a576101008083540402835291602001916104c5565b820191906000526020600020905b8154815290600101906020018083116104a857829003601f168201915b5050505050905092915050565b81816000868660405180838380828437820191505092505050908152602001604051809103902091906105069291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507f894bacec2200a39552d3a7b98024096fdcfed79005b6f8ee1651b68c62bba2a8848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106105fe57803560ff191683800117855561062c565b8280016001018555821561062c579182015b8281111561062b578235825591602001919060010190610610565b5b509050610639919061063d565b5090565b61065f91905b8082111561065b576000816000905550600101610643565b5090565b905600a165627a7a7230582037f7ef93352a4622feaf670d69e80fc6348c1f550fad9e99119634b80c10354d0029"};

    public static final String SM_BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", SM_BINARY_ARRAY);

    public static final Event DOMAIN_REGISTERED_EVENT = new Event("DomainRegistered", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
    ;

    protected ResourceDomain(String contractAddress, IClient client, CryptoKeyPair credential) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential);
    }

    protected ResourceDomain(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential, vmType);
    }

    protected ResourceDomain(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        super(binary, contractAddress, client, credential);
    }

    protected ResourceDomain(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        super(binary, contractAddress, client, credential, vmType);
    }

    public static ResourceDomain deploy(IClient client, CryptoKeyPair credential) throws ContractException {
        return deploy(ResourceDomain.class, client, credential, getBinary(client.getCryptoSuite()), "");
    }

    public static ResourceDomain deploy(IClient client, CryptoKeyPair credential, String binary) throws ContractException {
        return deploy(ResourceDomain.class, client, credential, binary, "");
    }

    public static ResourceDomain deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType) throws ContractException {
        return deploy(ResourceDomain.class, client, credential, getBinary(client.getCryptoSuite()), "", vmType);
    }

    public static ResourceDomain deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) throws ContractException {
        return deploy(ResourceDomain.class, client, credential, binary, "", vmType);
    }

    public static ResourceDomain load(String contractAddress, IClient client, CryptoKeyPair credential) {
        return new ResourceDomain(contractAddress, client, credential);
    }

    public static ResourceDomain load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        return new ResourceDomain(contractAddress, client, credential, vmType);
    }

    public static ResourceDomain load(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        return new ResourceDomain(contractAddress, client, credential, binary);
    }

    public static ResourceDomain load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        return new ResourceDomain(contractAddress, client, credential, vmType, binary);
    }

    public TransactionReceipt updateDomainDetails(String domainAddress, String details) {
        final Function function = new Function(
                FUNC_UPDATE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] updateDomainDetails(String domainAddress, String details, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_UPDATE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForUpdateDomainDetails(String domainAddress, String details) {
        final Function function = new Function(
                FUNC_UPDATE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getUpdateDomainDetailsInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_UPDATE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(), 
                (String) results.get(1).getValue()
                );
    }

    public TransactionReceipt saveDomainDetails(String domainAddress, String details) {
        final Function function = new Function(
                FUNC_SAVE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] saveDomainDetails(String domainAddress, String details, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_SAVE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForSaveDomainDetails(String domainAddress, String details) {
        final Function function = new Function(
                FUNC_SAVE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress), 
                new com.lingshu.chain.sdk.codec.datatypes.Utf8String(details)), 
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getSaveDomainDetailsInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_SAVE_DOMAIN_DETAILS, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(), 
                (String) results.get(1).getValue()
                );
    }

    public String getDomainDetails(String domainAddress) throws ContractException {
        final Function function = new Function(FUNC_GET_DOMAIN_DETAILS, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeCallWithSingleValueReturn(function, String.class);
    }

    public Boolean doesDomainExist(String domainAddress) throws ContractException {
        final Function function = new Function(FUNC_DOES_DOMAIN_EXIST, 
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(domainAddress)), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeCallWithSingleValueReturn(function, Boolean.class);
    }

    public List<DomainRegisteredEvtResp> getDomainRegisteredEvents(TransactionReceipt txReceipt) {
        List<Contract.EvtValuesWithLog> valueList = extractEventParametersWithLog(DOMAIN_REGISTERED_EVENT, txReceipt);
        ArrayList<DomainRegisteredEvtResp> responseList = new ArrayList<DomainRegisteredEvtResp>(valueList.size());
        for (Contract.EvtValuesWithLog eventValues : valueList) {
            DomainRegisteredEvtResp evtResp = new DomainRegisteredEvtResp();
            evtResp.log = eventValues.getLog();
            evtResp.domainAddress = (String) eventValues.getNonIndexedValues().get(0).getValue();
            evtResp.details = (String) eventValues.getNonIndexedValues().get(1).getValue();
            responseList.add(evtResp);
        }
        return responseList;
    }

    public void subscribeDomainRegisteredEvent(String fromBlock, String toBlock, List<String> otherTopics, EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(DOMAIN_REGISTERED_EVENT);
        subscribeEvent(ABI,BINARY,topic0,fromBlock,toBlock,otherTopics,callback);
    }

    public void subscribeDomainRegisteredEvent(EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(DOMAIN_REGISTERED_EVENT);
        subscribeEvent(ABI,BINARY,topic0,callback);
    }

    public static String getBinary(CryptoSuite cryptoSuite) {
        return (cryptoSuite.getCryptoType() == CryptoType.ECDSA_TYPE ? BINARY : SM_BINARY);
    }

    public static class DomainRegisteredEvtResp {
        public TransactionReceipt.Logs log;

        public String domainAddress;

        public String details;
    }
}
