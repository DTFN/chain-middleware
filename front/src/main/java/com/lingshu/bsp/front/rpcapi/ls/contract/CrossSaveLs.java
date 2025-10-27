package com.lingshu.bsp.front.rpcapi.ls.contract;

import com.lingshu.chain.sdk.client.IClient;
import com.lingshu.chain.sdk.codec.datatypes.Event;
import com.lingshu.chain.sdk.codec.datatypes.Function;
import com.lingshu.chain.sdk.codec.datatypes.Type;
import com.lingshu.chain.sdk.codec.datatypes.TypeReference;
import com.lingshu.chain.sdk.codec.datatypes.Utf8String;
import com.lingshu.chain.sdk.codec.datatypes.generated.tuples.generated.Tuple1;
import com.lingshu.chain.sdk.codec.datatypes.generated.tuples.generated.Tuple2;
import com.lingshu.chain.sdk.contract.Contract;
import com.lingshu.chain.sdk.crypto.CryptoSuite;
import com.lingshu.chain.sdk.crypto.key.CryptoKeyPair;
import com.lingshu.chain.sdk.evtsub.EvtSubCallback;
import com.lingshu.chain.sdk.model.CryptoType;
import com.lingshu.chain.sdk.model.TransactionCallback;
import com.lingshu.chain.sdk.model.TransactionReceipt;
import com.lingshu.chain.sdk.tx.common.exception.ContractException;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

@SuppressWarnings("unchecked")
public class CrossSaveLs extends Contract {
    public static final String FUNC_CROSS_CHAIN_TRY = "CrossChainTry";

    public static final String FUNC_CROSS_CHAIN_CANCEL = "CrossChainCancel";

    public static final String FUNC_QUERY = "query";

    public static final String FUNC_CROSS_CHAIN_SAVE = "CrossChainSave";

    public static final String FUNC_CROSS_CHAIN_CONFIRM = "CrossChainConfirm";

    public static final String[] ABI_ARRAY = {"[{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"string\"}],\"name\":\"CrossChainTry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"CrossChainCancel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"query\",\"outputs\":[{\"name\":\"value\",\"type\":\"string\"},{\"name\":\"flag\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"string\"}],\"name\":\"CrossChainSave\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"CrossChainConfirm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"string\"}],\"name\":\"Test\",\"type\":\"event\"}]"};

    public static final String ABI = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", ABI_ARRAY);

    public static final String[] BINARY_ARRAY = {"608060405234801561001057600080fd5b50610bc6806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063746182b71461007257806375116b13146101215780637c2619291461018a578063d1666246146102d8578063eac718c414610387575b600080fd5b34801561007e57600080fd5b5061011f600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506103f0565b005b34801561012d57600080fd5b50610188600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610529565b005b34801561019657600080fd5b506101f1600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506105e1565b604051808060200180602001838103835285818151815260200191508051906020019080838360005b8381101561023557808201518184015260208101905061021a565b50505050905090810190601f1680156102625780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b8381101561029b578082015181840152602081019050610280565b50505050905090810190601f1680156102c85780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b3480156102e457600080fd5b50610385600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506107fb565b005b34801561039357600080fd5b506103ee600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610a3d565b005b806000836040518082805190602001908083835b6020831015156104295780518252602082019150602081019050602083039250610404565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020908051906020019061046f929190610af5565b506040805190810160405280600581526020017f66616c73650000000000000000000000000000000000000000000000000000008152506001836040518082805190602001908083835b6020831015156104de57805182526020820191506020810190506020830392506104b9565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610524929190610af5565b505050565b6040805190810160405280600681526020017f6661696c656400000000000000000000000000000000000000000000000000008152506001826040518082805190602001908083835b6020831015156105975780518252602082019150602081019050602083039250610572565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906105dd929190610af5565b5050565b6060806000836040518082805190602001908083835b60208310151561061c57805182526020820191506020810190506020830392506105f7565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390206001846040518082805190602001908083835b6020831015156106875780518252602082019150602081019050602083039250610662565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020818054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561074f5780601f106107245761010080835404028352916020019161074f565b820191906000526020600020905b81548152906001019060200180831161073257829003601f168201915b50505050509150808054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156107eb5780601f106107c0576101008083540402835291602001916107eb565b820191906000526020600020905b8154815290600101906020018083116107ce57829003601f168201915b5050505050905091509150915091565b806000836040518082805190602001908083835b602083101515610834578051825260208201915060208101905060208303925061080f565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020908051906020019061087a929190610af5565b506040805190810160405280600581526020017f66616c73650000000000000000000000000000000000000000000000000000008152506001836040518082805190602001908083835b6020831015156108e957805182526020820191506020810190506020830392506108c4565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020908051906020019061092f929190610af5565b507f74cb234c0dd0ccac09c19041a69978ccb865f1f44a2877a009549898f6395b108282604051808060200180602001838103835285818151815260200191508051906020019080838360005b8381101561099757808201518184015260208101905061097c565b50505050905090810190601f1680156109c45780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156109fd5780820151818401526020810190506109e2565b50505050905090810190601f168015610a2a5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a15050565b6040805190810160405280600481526020017f74727565000000000000000000000000000000000000000000000000000000008152506001826040518082805190602001908083835b602083101515610aab5780518252602082019150602081019050602083039250610a86565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610af1929190610af5565b5050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b3657805160ff1916838001178555610b64565b82800160010185558215610b64579182015b82811115610b63578251825591602001919060010190610b48565b5b509050610b719190610b75565b5090565b610b9791905b80821115610b93576000816000905550600101610b7b565b5090565b905600a165627a7a72305820ba214fa9b5ca4bd1c7e96371a29c45f3795e2a338699919dbec571e889f5a7830029"};

    public static final String BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", BINARY_ARRAY);

    public static final String[] SM_BINARY_ARRAY = {"608060405234801561001057600080fd5b50610bc6806100206000396000f30060806040526004361061006d576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630958b9d4146100725780637135a360146100db5780639040597d1461018a578063acd6c172146101f3578063e2f2f80414610341575b600080fd5b34801561007e57600080fd5b506100d9600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506103f0565b005b3480156100e757600080fd5b50610188600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506104a8565b005b34801561019657600080fd5b506101f1600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506105e1565b005b3480156101ff57600080fd5b5061025a600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610699565b604051808060200180602001838103835285818151815260200191508051906020019080838360005b8381101561029e578082015181840152602081019050610283565b50505050905090810190601f1680156102cb5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156103045780820151818401526020810190506102e9565b50505050905090810190601f1680156103315780820380516001836020036101000a031916815260200191505b5094505050505060405180910390f35b34801561034d57600080fd5b506103ee600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506108b3565b005b6040805190810160405280600681526020017f6661696c656400000000000000000000000000000000000000000000000000008152506001826040518082805190602001908083835b60208310151561045e5780518252602082019150602081019050602083039250610439565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906104a4929190610af5565b5050565b806000836040518082805190602001908083835b6020831015156104e157805182526020820191506020810190506020830392506104bc565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610527929190610af5565b506040805190810160405280600581526020017f66616c73650000000000000000000000000000000000000000000000000000008152506001836040518082805190602001908083835b6020831015156105965780518252602082019150602081019050602083039250610571565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906105dc929190610af5565b505050565b6040805190810160405280600481526020017f74727565000000000000000000000000000000000000000000000000000000008152506001826040518082805190602001908083835b60208310151561064f578051825260208201915060208101905060208303925061062a565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610695929190610af5565b5050565b6060806000836040518082805190602001908083835b6020831015156106d457805182526020820191506020810190506020830392506106af565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390206001846040518082805190602001908083835b60208310151561073f578051825260208201915060208101905060208303925061071a565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108075780601f106107dc57610100808354040283529160200191610807565b820191906000526020600020905b8154815290600101906020018083116107ea57829003601f168201915b50505050509150808054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156108a35780601f10610878576101008083540402835291602001916108a3565b820191906000526020600020905b81548152906001019060200180831161088657829003601f168201915b5050505050905091509150915091565b806000836040518082805190602001908083835b6020831015156108ec57805182526020820191506020810190506020830392506108c7565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610932929190610af5565b506040805190810160405280600581526020017f66616c73650000000000000000000000000000000000000000000000000000008152506001836040518082805190602001908083835b6020831015156109a1578051825260208201915060208101905060208303925061097c565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906109e7929190610af5565b507ff163f45022c10df5cb23dbf517339aaad955630026d0acf2d8253aa935c7d7fd8282604051808060200180602001838103835285818151815260200191508051906020019080838360005b83811015610a4f578082015181840152602081019050610a34565b50505050905090810190601f168015610a7c5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b83811015610ab5578082015181840152602081019050610a9a565b50505050905090810190601f168015610ae25780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a15050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610b3657805160ff1916838001178555610b64565b82800160010185558215610b64579182015b82811115610b63578251825591602001919060010190610b48565b5b509050610b719190610b75565b5090565b610b9791905b80821115610b93576000816000905550600101610b7b565b5090565b905600a165627a7a723058209e80da5be4862486ca21d73932a68ff7b398ca6cb0d2f9789b4163c8c26035d10029"};

    public static final String SM_BINARY = com.lingshu.chain.sdk.codegen.util.GeneratorUtil.joinAll("", SM_BINARY_ARRAY);

    public static final Event TEST_EVENT = new Event("Test",
            Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
    ;

    protected CrossSaveLs(String contractAddress, IClient client, CryptoKeyPair credential) {
        super(getBinary(client.getCryptoSuite()), contractAddress, client, credential);
    }

    public static CrossSaveLs deploy(IClient client, CryptoKeyPair credential) throws ContractException {
        return deploy(CrossSaveLs.class, client, credential, getBinary(client.getCryptoSuite()), "");
    }

    public static CrossSaveLs load(String contractAddress, IClient client, CryptoKeyPair credential) {
        return new CrossSaveLs(contractAddress, client, credential);
    }

    public TransactionReceipt CrossChainTry(String key, String value) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_TRY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] CrossChainTry(String key, String value, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_TRY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForCrossChainTry(String key, String value) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_TRY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getCrossChainTryInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_CROSS_CHAIN_TRY,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(),
                (String) results.get(1).getValue()
        );
    }

    public TransactionReceipt CrossChainCancel(String key) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CANCEL,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] CrossChainCancel(String key, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CANCEL,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForCrossChainCancel(String key) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CANCEL,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple1<String> getCrossChainCancelInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_CROSS_CHAIN_CANCEL,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple1<String>(

                (String) results.get(0).getValue()
        );
    }

    public TransactionReceipt query(String key) {
        final Function function = new Function(
                FUNC_QUERY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] query(String key, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_QUERY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForQuery(String key) {
        final Function function = new Function(
                FUNC_QUERY,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple1<String> getQueryInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_QUERY,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple1<String>(

                (String) results.get(0).getValue()
        );
    }

    public Tuple2<String, String> getQueryOutput(TransactionReceipt transactionReceipt) {
        String data = transactionReceipt.getOutput();
        final Function function = new Function(FUNC_QUERY,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(),
                (String) results.get(1).getValue()
        );
    }

    public TransactionReceipt CrossChainSave(String key, String value) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_SAVE,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] CrossChainSave(String key, String value, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_SAVE,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForCrossChainSave(String key, String value) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_SAVE,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key),
                        new com.lingshu.chain.sdk.codec.datatypes.Utf8String(value)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple2<String, String> getCrossChainSaveInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_CROSS_CHAIN_SAVE,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}, new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple2<String, String>(

                (String) results.get(0).getValue(),
                (String) results.get(1).getValue()
        );
    }

    public TransactionReceipt CrossChainConfirm(String key) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CONFIRM,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return executeTransaction(function);
    }

    public byte[] CrossChainConfirm(String key, TransactionCallback callback) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CONFIRM,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return asyncExecuteTransaction(function, callback);
    }

    public String getSignedTxForCrossChainConfirm(String key) {
        final Function function = new Function(
                FUNC_CROSS_CHAIN_CONFIRM,
                Arrays.<Type>asList(new com.lingshu.chain.sdk.codec.datatypes.Utf8String(key)),
                Collections.<TypeReference<?>>emptyList());
        return createSignedTransaction(function);
    }

    public Tuple1<String> getCrossChainConfirmInput(TransactionReceipt txReceipt) {
        String data = txReceipt.getInput().substring(10);
        final Function function = new Function(FUNC_CROSS_CHAIN_CONFIRM,
                Arrays.<Type>asList(),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {}));
        List<Type> results = funcReturnDecoder.decode(data, function.getOutputParameters());
        return new Tuple1<String>(

                (String) results.get(0).getValue()
        );
    }

    public List<TestEvtResp> getTestEvents(TransactionReceipt txReceipt) {
        List<Contract.EvtValuesWithLog> valueList = extractEventParametersWithLog(TEST_EVENT, txReceipt);
        ArrayList<TestEvtResp> responseList = new ArrayList<TestEvtResp>(valueList.size());
        for (Contract.EvtValuesWithLog eventValues : valueList) {
            TestEvtResp evtResp = new TestEvtResp();
            evtResp.log = eventValues.getLog();
            evtResp.key = (String) eventValues.getNonIndexedValues().get(0).getValue();
            evtResp.value = (String) eventValues.getNonIndexedValues().get(1).getValue();
            responseList.add(evtResp);
        }
        return responseList;
    }

    public void subscribeTestEvent(String fromBlock, String toBlock, List<String> otherTopics, EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(TEST_EVENT);
        subscribeEvent(ABI,BINARY,topic0,fromBlock,toBlock,otherTopics,callback);
    }

    public void subscribeTestEvent(EvtSubCallback callback) {
        String topic0 = evtEncoder.encode(TEST_EVENT);
        subscribeEvent(ABI,BINARY,topic0,callback);
    }

    public static String getBinary(CryptoSuite cryptoSuite) {
        return (cryptoSuite.getCryptoType() == CryptoType.ECDSA_TYPE ? BINARY : SM_BINARY);
    }

    public static class TestEvtResp {
        public TransactionReceipt.Logs log;

        public String key;

        public String value;
    }
}
