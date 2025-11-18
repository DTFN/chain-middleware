package com.lingshu.server.core.web3j.contract;

import com.lingshu.server.core.web3j.chainmaker.client.ChainmakerChainClient;
import com.lingshu.server.core.web3j.chainmaker.util.ChainmakerAccountUtil;
import lombok.Getter;
import org.web3j.abi.TypeReference;
import org.web3j.abi.datatypes.Bool;
import org.web3j.abi.datatypes.Function;
import org.web3j.abi.datatypes.Type;
import org.web3j.abi.datatypes.Utf8String;
import org.web3j.protocol.core.methods.response.TransactionReceipt;

import java.util.Arrays;
import java.util.Collections;

/**
 * @author: derrick
 * @since: 2025-08-25
 */
public class ResourceDomain extends ChainmakerContract {

//    public static final String CONTRACT_NAME = "RESOURCE_DOMAIN";

    public static final String BINARY = "608060405234801561001057600080fd5b5061068e806100206000396000f300608060405260043610610062576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680632e2cf984146100675780635574e0d2146100ba578063c38fb8361461010d578063d7b7e4cd146101c1575b600080fd5b34801561007357600080fd5b506100b8600480360381019080803590602001908201803590602001919091929391929390803590602001908201803590602001919091929391929390505050610214565b005b3480156100c657600080fd5b5061010b6004803603810190808035906020019082018035906020019190919293919293908035906020019082018035906020019190919293919293905050506102ff565b005b34801561011957600080fd5b506101466004803603810190808035906020019082018035906020019190919293919293905050506104b8565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561018657808201518184015260208101905061016b565b50505050905090810190601f1680156101b35780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101cd57600080fd5b506101fa600480360381019080803590602001908201803590602001919091929391929390505050610580565b604051808215151515815260200191505060405180910390f35b81816000868660405180838380828437820191505092505050908152602001604051809103902091906102489291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507ff5e947bc58a45ca3f7c47f7a942b67608bb05338693c9dd564f093e315b8962e848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6001848460405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff161515156103cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260258152602001807f5265736f75726365446f6d61696e3a20646f6d61696e20616c7265616479206581526020017f786973747300000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b81816000868660405180838380828437820191505092505050908152602001604051809103902091906104019291906105bd565b50600180858560405180838380828437820191505092505050908152602001604051809103902060006101000a81548160ff0219169083151502179055507ff5e947bc58a45ca3f7c47f7a942b67608bb05338693c9dd564f093e315b8962e848484846040518080602001806020018381038352878782818152602001925080828437820191505083810382528585828181526020019250808284378201915050965050505050505060405180910390a150505050565b6060600083836040518083838082843782019150509250505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105735780601f1061054857610100808354040283529160200191610573565b820191906000526020600020905b81548152906001019060200180831161055657829003601f168201915b5050505050905092915050565b60006001838360405180838380828437820191505092505050908152602001604051809103902060009054906101000a900460ff16905092915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106105fe57803560ff191683800117855561062c565b8280016001018555821561062c579182015b8281111561062b578235825591602001919060010190610610565b5b509050610639919061063d565b5090565b61065f91905b8082111561065b576000816000905550600101610643565b5090565b905600a165627a7a72305820ef13e280805350eb3aca9dd10150f9950a27a89da1828d3a74214e76a7036dc10029";

    public static final String FUNC_SAVE_DOMAIN_DETAILS = "saveDomainDetails";
    public static final String FUNC_UPDATE_DOMAIN_DETAILS = "updateDomainDetails";

    public static final String FUNC_GET_DOMAIN_DETAILS = "getDomainDetails";

    public static final String FUNC_DOES_DOMAIN_EXIST = "doesDomainExist";

    @Getter
    private DeployInfo deployInfo;

    public ResourceDomain(String contractName, ChainmakerChainClient chainmakerChainClient, ChainmakerAccountUtil chainmakerAccountUtil) {
        super(contractName, chainmakerChainClient, chainmakerAccountUtil);
    }

    public static ResourceDomain load(String contractName, ChainmakerChainClient chainmakerChainClient, ChainmakerAccountUtil chainmakerAccountUtil) {
        return new ResourceDomain(contractName, chainmakerChainClient, chainmakerAccountUtil);
    }

    public static ResourceDomain deploy(ChainmakerAccountUtil chainmakerAccountUtil, ChainmakerChainClient chainmakerChainClient, String contractName) throws Exception {
        Function function = new Function("",
//                Arrays.<Type>asList(new Utf8String(name), new Utf8String(symbol)),
                Collections.emptyList(), Collections.emptyList());
        DeployInfo deployInfo = ChainmakerContract.deployContract(chainmakerChainClient, contractName, function,
                null, BINARY.getBytes());
        ResourceDomain resourceDomain = load(contractName, chainmakerChainClient, chainmakerAccountUtil);
        resourceDomain.deployInfo = deployInfo;
        return resourceDomain;
    }

    public TransactionReceipt saveDomainDetails(String domainAddress, String details) {
        final Function function = new Function(FUNC_SAVE_DOMAIN_DETAILS,
                Arrays.<Type>asList(new Utf8String(domainAddress), new Utf8String(details)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public TransactionReceipt updateDomainDetails(String domainAddress, String details) {
        final Function function = new Function(FUNC_UPDATE_DOMAIN_DETAILS,
                Arrays.<Type>asList(new Utf8String(domainAddress), new Utf8String(details)),
                Collections.<TypeReference<?>>emptyList());
        return executeRemoteCallTransaction(function);
    }

    public String getDomainDetails(String domainAddress) /*throws ContractException*/ {
        final Function function = new Function(FUNC_GET_DOMAIN_DETAILS,
                Arrays.<Type>asList(new Utf8String(domainAddress)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Utf8String>() {
                }));
        return executeRemoteCallSingleValueReturn(function, String.class);
    }

    public Boolean doesDomainExist(String domainAddress) /*throws ContractException*/ {
        final Function function = new Function(FUNC_DOES_DOMAIN_EXIST,
                Arrays.<Type>asList(new Utf8String(domainAddress)),
                Arrays.<TypeReference<?>>asList(new TypeReference<Bool>() {
                }));
//        return executeRemoteCallSingleValueReturn(function, Boolean.class);
        return executeRemoteCallSingleValueReturn(function, Boolean.class);
    }
}
