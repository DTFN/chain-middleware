package com.lingshu.bsp.front.rpcapi.ls.contract;

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
public class DIDManagerLs extends Contract {
    public static final String FUNC_CREATE_DID = "createDID";

    public static final String FUNC_GET_DID_DETAILS = "getDIDDetails";

    public static final String FUNC_UPDATE_DID = "updateDID";

    public static final String FUNC_DOES_DID_EXIST = "doesDIDExist";

    public static final String FUNC_ECHO = "echo";

    public static final String[] ABI_ARRAY = {"[{\"constant\":false,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"},{\"name\":\"didDocument\",\"type\":\"string\"}],\"name\":\"createDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"getDIDDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"},{\"name\":\"newDidDocument\",\"type\":\"string\"}],\"name\":\"updateDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"doesDIDExist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"echo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"did\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"didDocument\",\"type\":\"string\"}],\"name\":\"DIDModify\",\"type\":\"event\"}]"};

    public static final String ABI = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", ABI_ARRAY);

    public static final String[] BINARY_ARRAY = {"608060405234801561001057600080fd5b50610764806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680631539a50714610072578063a3a7c60b146100c5578063b42a02b514610179578063dd742fa5146101cc578063f15da7291461021f575b600080fd5b34801561007e57600080fd5b506100c36004803603810190808035906020019082018035906020019190919293919293908035906020019082018035906020019190919293919293905050506102d3565b005b3480156100d157600080fd5b506100fe600480360381019080803590602001908201803590602001919091929391929390505050610466565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561013e578082015181840152602081019050610123565b50505050905090810190601f16801561016b5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561018557600080fd5b506101ca60048036038101908080359060200190820180359060200191909192939192939080359060200190820180359060200191909192939192939050505061052e565b005b3480156101d857600080fd5b50610205600480360381019080803590602001908201803590602001919091929391929390505050610619565b604051808215151515815260200191505060405180910390f35b34801561022b57600080fd5b50610258600480360381019080803590602001908201803590602001919091929391929390505050610656565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561029857808201518184015260208101905061027d565b50505050905090810190601f1680156102c55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff1615151561037b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f4449444d616e616765723a2044494420616c726561647920657869737473000081525060200191505060405180910390fd5b81816000868660405180838380828437820191505092505050908152602001604051809103902091906103af929190610693565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507fbc4b9177d949f0995aac6fc1e139b1e35a50170b6e9ca10360f73c3ac103d905848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105215780601f106104f657610100808354040283529160200191610521565b820191906000526020600020905b81548152906001019060200180831161050457829003601f168201915b5050505050905092915050565b8181600086866040518083838082843782019150509250505090815260200160405180910390209190610562929190610693565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507fbc4b9177d949f0995aac6fc1e139b1e35a50170b6e9ca10360f73c3ac103d905848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b60006001838360405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16905092915050565b606082828080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106d457803560ff1916838001178555610702565b82800160010185558215610702579182015b828111156107015782358255916020019190600101906106e6565b5b50905061070f9190610713565b5090565b61073591905b80821115610731576000816000905550600101610719565b5090565b905600a165627a7a7230582077fbc1ae35dcd17a7766d823aacc4e08f63634daaea5b5cd27a47fc092a6f2600029"};

    public static final String BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", BINARY_ARRAY);

    public static final String[] SM_BINARY_ARRAY = {"608060405234801561001057600080fd5b50610764806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680634913a9a11461007257806374a27ec214610126578063b197713214610179578063c00059e5146101cc578063e1ed624c1461021f575b600080fd5b34801561007e57600080fd5b506100ab6004803603810190808035906020019082018035906020019190919293919293905050506102d3565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100eb5780820151818401526020810190506100d0565b50505050905090810190601f1680156101185780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561013257600080fd5b50610177600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610310565b005b34801561018557600080fd5b506101b26004803603810190808035906020019082018035906020019190919293919293905050506103fb565b604051808215151515815260200191505060405180910390f35b3480156101d857600080fd5b5061021d600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610438565b005b34801561022b57600080fd5b506102586004803603810190808035906020019082018035906020019190919293919293905050506105cb565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561029857808201518184015260208101905061027d565b50505050905090810190601f1680156102c55780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b606082828080601f016020809104026020016040519081016040528093929190818152602001838380828437820191505050505050905092915050565b8181600086866040518083838082843782019150509250505090815260200160405180910390209190610344929190610693565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507fc40bd4c8b1923a877fd83a34670f7018ddd78ca59c2e5358c09c2854402995cb848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b60006001838360405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16905092915050565b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff161515156104e0576040517fc703cb1200000000000000000000000000000000000000000000000000000000815260040180806020018281038252601e8152602001807f4449444d616e616765723a2044494420616c726561647920657869737473000081525060200191505060405180910390fd5b8181600086866040518083838082843782019150509250505090815260200160405180910390209190610514929190610693565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507fc40bd4c8b1923a877fd83a34670f7018ddd78ca59c2e5358c09c2854402995cb848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156106865780601f1061065b57610100808354040283529160200191610686565b820191906000526020600020905b81548152906001019060200180831161066957829003601f168201915b5050505050905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106d457803560ff1916838001178555610702565b82800160010185558215610702579182015b828111156107015782358255916020019190600101906106e6565b5b50905061070f9190610713565b5090565b61073591905b80821115610731576000816000905550600101610719565b5090565b905600a165627a7a72305820a9ad4ffac8838c625e098e02ad9ec14c4bd29417d441ed66c8895bb0d5d264410029"};

    public static final String SM_BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", SM_BINARY_ARRAY);

    public static final Event DID_MODIFY_EVENT = new Event("DIDModify", 
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
    ;

    protected DIDManagerLs(String contractAddress, IClient client, CryptoKeyPair credential) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential);
    }

    protected DIDManagerLs(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential, vmType);
    }

    protected DIDManagerLs(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        super(binary, contractAddress, client, credential);
    }

    protected DIDManagerLs(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        super(binary, contractAddress, client, credential, vmType);
    }

    public static DIDManagerLs deploy(IClient client, CryptoKeyPair credential) throws ContractException {
        return deploy(DIDManagerLs.class, client, credential, getBinary(client.getCryptoSuite()), "");
    }

    public static DIDManagerLs deploy(IClient client, CryptoKeyPair credential, String binary) throws ContractException {
        return deploy(DIDManagerLs.class, client, credential, binary, "");
    }

    public static DIDManagerLs deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType) throws ContractException {
        return deploy(DIDManagerLs.class, client, credential, getBinary(client.getCryptoSuite()), "", vmType);
    }

    public static DIDManagerLs deploy(IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) throws ContractException {
        return deploy(DIDManagerLs.class, client, credential, binary, "", vmType);
    }

    public static DIDManagerLs load(String contractAddress, IClient client, CryptoKeyPair credential) {
        return new DIDManagerLs(contractAddress, client, credential);
    }

    public static DIDManagerLs load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType) {
        return new DIDManagerLs(contractAddress, client, credential, vmType);
    }

    public static DIDManagerLs load(String contractAddress, IClient client, CryptoKeyPair credential, String binary) {
        return new DIDManagerLs(contractAddress, client, credential, binary);
    }

    public static DIDManagerLs load(String contractAddress, IClient client, CryptoKeyPair credential, VmTypeEnum vmType, String binary) {
        return new DIDManagerLs(contractAddress, client, credential, vmType, binary);
    }

    public TransactionReceipt createDID(String did, String didDocument) {
        final Function function = new Function(
                FUNC_CREATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(didDocument)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] createDID(String did, String didDocument, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_CREATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(didDocument)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForCreateDID(String did, String didDocument) {
        final Function function = new Function(
                FUNC_CREATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(didDocument)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getCreateDIDInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_CREATE_DID, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(), 
                (String) results.get(1).getValue()
                );
    }

    public String getDIDDetails(String did) throws ContractException {
        final Function function = new Function(FUNC_GET_DID_DETAILS, 
                Arrays.<Type>asList(new Utf8String(did)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeCallWithSingleValueReturn(function, String.class);
    }

    public TransactionReceipt updateDID(String did, String newDidDocument) {
        final Function function = new Function(
                FUNC_UPDATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(newDidDocument)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] updateDID(String did, String newDidDocument, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_UPDATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(newDidDocument)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForUpdateDID(String did, String newDidDocument) {
        final Function function = new Function(
                FUNC_UPDATE_DID, 
                Arrays.<Type>asList(new Utf8String(did),
                new Utf8String(newDidDocument)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getUpdateDIDInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_UPDATE_DID, 
                Arrays.<Type>asList(), 
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(), 
                (String) results.get(1).getValue()
                );
    }

    public Boolean doesDIDExist(String did) throws ContractException {
        final Function function = new Function(FUNC_DOES_DID_EXIST, 
                Arrays.<Type>asList(new Utf8String(did)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {}));
        return executeCallWithSingleValueReturn(function, Boolean.class);
    }

    public String echo(String did) throws ContractException {
        final Function function = new Function(FUNC_ECHO, 
                Arrays.<Type>asList(new Utf8String(did)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        return executeCallWithSingleValueReturn(function, String.class);
    }

    public List<DIDModifyEvtResp> getDIDModifyEvents(TransactionReceipt txReceipt) {
        List<EvtValuesWithLog> valueList = extractEventParametersWithLog(DID_MODIFY_EVENT, txReceipt);
        ArrayList<DIDModifyEvtResp> responseList = new ArrayList<DIDModifyEvtResp>(valueList.size());
        for (EvtValuesWithLog eventValues : valueList) {
            DIDModifyEvtResp evtResp = new DIDModifyEvtResp();
            evtResp.log = eventValues.getLog();
            evtResp.did = (String) eventValues.getNonIndexedValues().get(0).getValue();
            evtResp.didDocument = (String) eventValues.getNonIndexedValues().get(1).getValue();
            responseList.add(evtResp);
        }
        return responseList;
    }

    public void subscribeDIDModifyEvent(String fromBlock, String toBlock, List<String> otherTopics, EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(DID_MODIFY_EVENT);
        subscribeEvent(ABI,BINARY,topic0,fromBlock,toBlock,otherTopics,callback);
    }

    public void subscribeDIDModifyEvent(EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(DID_MODIFY_EVENT);
        subscribeEvent(ABI,BINARY,topic0,callback);
    }

    public static String getBinary(CryptoSuite cryptoSuite) {
        return (cryptoSuite.getCryptoType() == CryptoType.ECDSA_TYPE ? BINARY : SM_BINARY);
    }

    public static class DIDModifyEvtResp {
        public TransactionReceipt.Logs log;

        public String did;

        public String didDocument;
    }
}
