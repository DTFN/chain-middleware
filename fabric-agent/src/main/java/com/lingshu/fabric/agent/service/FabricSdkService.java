package com.lingshu.fabric.agent.service;

import com.alibaba.fastjson2.JSON;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.event.GatewayInitedEvent;
import com.lingshu.fabric.agent.exception.AgentException;
import com.lingshu.fabric.agent.repo.ChainCodeRepo;
import com.lingshu.fabric.agent.repo.entity.ChainCodePeerDo;
import com.lingshu.fabric.agent.req.chaincode.*;
import com.lingshu.fabric.agent.resp.FabricChainCodeInvokeResp;
import com.lingshu.fabric.agent.resp.PackageInfoDTO;
import lombok.extern.slf4j.Slf4j;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.ChaincodeEndorsementPolicyParseException;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Component;
import org.springframework.util.CollectionUtils;

import java.io.File;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.InvalidKeyException;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.stream.Collectors;

@Component
@Slf4j
public class FabricSdkService {

    @Value("${chaincode.timeout:60000}")
    private long TIMEOUT = 60 * 1000;
    @Autowired
    private ConstantProperties constantProperties;
    @Autowired
    private ChainCodeRepo chainCodeRepo;
    private static final String installedMsg = "chaincode already successfully installed";
    private static final String approvedMsg = "attempted to redefine uncommitted sequence";

    public LifecycleChaincodePackage createLifecycleChaincodePackage(String chaincodeLabel, TransactionRequest.Type chaincodeType,
                                                                     String chaincodeSourceLocation, String chaincodePath, String metadataSource) throws IOException, InvalidArgumentException {
        log.info("creating install package {}", chaincodeLabel);

        Path metadataSourcePath = null;
        if (metadataSource != null) {
            metadataSourcePath = Paths.get(metadataSource);
        }

        return LifecycleChaincodePackage.fromSource(chaincodeLabel, Paths.get(chaincodeSourceLocation),
                chaincodeType,
                chaincodePath, metadataSourcePath);
    }

    public String lifecycleInstallChaincode(HFClient client, Collection<Peer> peers, LifecycleChaincodePackage lifecycleChaincodePackage)
            throws ProposalException, InvalidArgumentException {

        int numInstallProposal = 0;

        numInstallProposal = numInstallProposal + peers.size();

        LifecycleInstallChaincodeRequest installProposalRequest = client.newLifecycleInstallChaincodeRequest();
        installProposalRequest.setLifecycleChaincodePackage(lifecycleChaincodePackage);
        installProposalRequest.setProposalWaitTime(TIMEOUT);

        Collection<LifecycleInstallChaincodeProposalResponse> responses = client.sendLifecycleInstallChaincodeRequest(installProposalRequest, peers);

        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();
        String packageID = null;
        for (LifecycleInstallChaincodeProposalResponse response : responses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                log.info("Successful install proposal response Txid: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
                if (packageID == null) {
                    packageID = response.getPackageId();
                }
            } else if (response.getMessage().contains(installedMsg)) {
                log.info(installedMsg + " proposal response Txid: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
            } else {

                log.error("lifecycleInstallChaincode failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
                        response.getPeer(), response.getStatus(), response.isInvalid(), response.isVerified(), response.getMessage());
                failed.add(response);
            }
        }

        log.info("Received {} install proposal responses. Successful. verified: {} . Failed: {}", numInstallProposal, successful.size(), failed.size());

        if (failed.size() > 0) {
//            ProposalResponse first = failed.iterator().next();
//            log.error("lifecycleInstallChaincode failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
//                    first.getPeer(), first.getStatus(), first.isInvalid(), first.isVerified(), first.getMessage());
            throw new AgentException(ConstantCode.CHAIN_CODE_INSTALL_FAILED);
        }
        return packageID;

    }

    public String lifecycleQueryInstalledChaincode(HFClient client, Collection<Peer> peers, String packageId)
            throws ProposalException, InvalidArgumentException {

        LifecycleQueryInstalledChaincodeRequest installProposalRequest = client.newLifecycleQueryInstalledChaincodeRequest();
        installProposalRequest.setPackageID(packageId);
        installProposalRequest.setProposalWaitTime(TIMEOUT);

        Collection<LifecycleQueryInstalledChaincodeProposalResponse> responses = client.sendLifecycleQueryInstalledChaincode(installProposalRequest, peers);

        String label = null;
        for (LifecycleQueryInstalledChaincodeProposalResponse response : responses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                if (label == null) {
                    label = response.getLabel();
                }
                log.info("Successful install proposal response txId: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
            }
        }
        return label;
    }

    public Map<String, String> lifecycleQueryInstalledChainCodes(HFClient client, String chaincodeLabel, Collection<Peer> peers)
            throws ProposalException, InvalidArgumentException {
        LifecycleQueryInstalledChaincodesRequest lifecycleQueryInstalledChaincodesRequest = client.newLifecycleQueryInstalledChaincodesRequest();
        lifecycleQueryInstalledChaincodesRequest.setProposalWaitTime(TIMEOUT);

        Collection<LifecycleQueryInstalledChaincodesProposalResponse> responses = client.sendLifecycleQueryInstalledChaincodes(lifecycleQueryInstalledChaincodesRequest, peers);

        Map<String, String> map = new HashMap<>();
        for (LifecycleQueryInstalledChaincodesProposalResponse response : responses) {

            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                Collection<LifecycleQueryInstalledChaincodesProposalResponse.LifecycleQueryInstalledChaincodesResult> lifecycleQueryInstalledChaincodesResult = response.getLifecycleQueryInstalledChaincodesResult();
                if (lifecycleQueryInstalledChaincodesResult == null) {
                    continue;
                }
                for (LifecycleQueryInstalledChaincodesProposalResponse.LifecycleQueryInstalledChaincodesResult result : lifecycleQueryInstalledChaincodesResult) {
                    String label = result.getLabel();
                    String packageId = result.getPackageId();
                    if (label.equals(chaincodeLabel)) {
                        map.put(response.getPeer().getName(), packageId);
                        break;
                    }
                }
                log.info("Successful query installed chaincodes proposal response txId: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
            }
        }
        return map;

    }

    public Map<String, String> lifecycleQueryInstalledChainCodes(HFClient client, Collection<Peer> peers)
            throws ProposalException, InvalidArgumentException {
        LifecycleQueryInstalledChaincodesRequest lifecycleQueryInstalledChaincodesRequest = client.newLifecycleQueryInstalledChaincodesRequest();
        lifecycleQueryInstalledChaincodesRequest.setProposalWaitTime(TIMEOUT);

        Collection<LifecycleQueryInstalledChaincodesProposalResponse> responses = client.sendLifecycleQueryInstalledChaincodes(lifecycleQueryInstalledChaincodesRequest, peers);

        Map<String, String> map = new HashMap<>();
        for (LifecycleQueryInstalledChaincodesProposalResponse response : responses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                Collection<LifecycleQueryInstalledChaincodesProposalResponse.LifecycleQueryInstalledChaincodesResult> lifecycleQueryInstalledChaincodesResult = response.getLifecycleQueryInstalledChaincodesResult();
                if (lifecycleQueryInstalledChaincodesResult == null) {
                    continue;
                }
                for (LifecycleQueryInstalledChaincodesProposalResponse.LifecycleQueryInstalledChaincodesResult result : lifecycleQueryInstalledChaincodesResult) {
                    String label = result.getLabel();
                    String packageId = result.getPackageId();
                    map.put(label, packageId);
                }
                log.info("Successful query installed chaincodes proposal response txId: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
            }
        }
        return map;

    }

    public List<String> queryApprovedPeers(HFClient client, Channel channel, long definitionSequence, String chaincodeName,
             String chaincodeVersion, LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy,
             ChaincodeCollectionConfiguration chaincodeCollectionConfiguration, boolean initRequired, Collection<Peer> org1MyPeers) throws InvalidArgumentException, ProposalException {
        LifecycleCheckCommitReadinessRequest lifecycleCheckCommitReadinessRequest = client.newLifecycleSimulateCommitChaincodeDefinitionRequest();
        lifecycleCheckCommitReadinessRequest.setSequence(definitionSequence);
        lifecycleCheckCommitReadinessRequest.setChaincodeName(chaincodeName);
        lifecycleCheckCommitReadinessRequest.setChaincodeVersion(chaincodeVersion);
        if (null != chaincodeEndorsementPolicy) {
            lifecycleCheckCommitReadinessRequest.setChaincodeEndorsementPolicy(chaincodeEndorsementPolicy);
        }
        if (null != chaincodeCollectionConfiguration) {
            lifecycleCheckCommitReadinessRequest.setChaincodeCollectionConfiguration(chaincodeCollectionConfiguration);
        }
        lifecycleCheckCommitReadinessRequest.setInitRequired(initRequired);

        List<String> approvedPeers = new ArrayList<>();
        Collection<LifecycleCheckCommitReadinessProposalResponse> lifecycleSimulateCommitChaincodeDefinitionProposalResponse = channel.sendLifecycleCheckCommitReadinessRequest(lifecycleCheckCommitReadinessRequest, org1MyPeers);
        for (LifecycleCheckCommitReadinessProposalResponse resp : lifecycleSimulateCommitChaincodeDefinitionProposalResponse) {
            final Peer peer = resp.getPeer();
            if (!CollectionUtils.isEmpty(resp.getApprovedOrgs())) {
                approvedPeers.add(peer.getName());
            }
        }
        return approvedPeers;
    }

    public CompletableFuture<BlockEvent.TransactionEvent> lifecycleApproveChaincodeDefinitionForMyOrg(
            HFClient client, Channel channel, Collection<Peer> peers, long sequence,
            String chaincodeName, String chaincodeVersion, LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy,
            ChaincodeCollectionConfiguration chaincodeCollectionConfiguration, boolean initRequired, String packageId)
            throws InvalidArgumentException, ProposalException {

        LifecycleApproveChaincodeDefinitionForMyOrgRequest lifecycleApproveChaincodeDefinitionForMyOrgRequest = client.newLifecycleApproveChaincodeDefinitionForMyOrgRequest();
        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setSequence(sequence);
        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setChaincodeName(chaincodeName);
        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setChaincodeVersion(chaincodeVersion);
        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setInitRequired(initRequired);

        if (null != chaincodeCollectionConfiguration) {
            lifecycleApproveChaincodeDefinitionForMyOrgRequest.setChaincodeCollectionConfiguration(chaincodeCollectionConfiguration);
        }

        if (null != chaincodeEndorsementPolicy) {
            lifecycleApproveChaincodeDefinitionForMyOrgRequest.setChaincodeEndorsementPolicy(chaincodeEndorsementPolicy);
        }

        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setPackageId(packageId);

        Collection<LifecycleApproveChaincodeDefinitionForMyOrgProposalResponse> lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse = channel.sendLifecycleApproveChaincodeDefinitionForMyOrgProposal(lifecycleApproveChaincodeDefinitionForMyOrgRequest,
                peers);

        int numInstallProposal = lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse.size();
        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();
        for (LifecycleApproveChaincodeDefinitionForMyOrgProposalResponse response : lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                log.info("Successful approve proposal response Txid: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
            } else {
                log.error("lifecycleApproveChaincodeDefinitionForMyOrg failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
                        response.getPeer(), response.getStatus(), response.isInvalid(), response.isVerified(), response.getMessage());
                failed.add(response);
            }

        }

        log.info("Received {} approve proposal responses. Successful. verified: {} . Failed: {}", numInstallProposal, successful.size(), failed.size());

        if (failed.size() > 0) {
//            ProposalResponse first = failed.iterator().next();
//            log.error("lifecycleApproveChaincodeDefinitionForMyOrg failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
//                    first.getPeer(), first.getStatus(), first.isInvalid(), first.isVerified(), first.getMessage());
            throw new AgentException(ConstantCode.CHAIN_CODE_APPROVE_FAILED);
        }
        return channel.sendTransaction(lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse);

    }

    public CompletableFuture<BlockEvent.TransactionEvent> commitChaincodeDefinitionRequest(
            HFClient client, Channel channel, long definitionSequence, String chaincodeName, String chaincodeVersion,
            LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy,
            ChaincodeCollectionConfiguration chaincodeCollectionConfiguration,
            boolean initRequired, Collection<Peer> endorsingPeers) throws ProposalException, InvalidArgumentException {
        LifecycleCommitChaincodeDefinitionRequest lifecycleCommitChaincodeDefinitionRequest = client.newLifecycleCommitChaincodeDefinitionRequest();

        lifecycleCommitChaincodeDefinitionRequest.setSequence(definitionSequence);
        lifecycleCommitChaincodeDefinitionRequest.setChaincodeName(chaincodeName);
        lifecycleCommitChaincodeDefinitionRequest.setChaincodeVersion(chaincodeVersion);
        if (null != chaincodeEndorsementPolicy) {
            lifecycleCommitChaincodeDefinitionRequest.setChaincodeEndorsementPolicy(chaincodeEndorsementPolicy);
        }
        if (null != chaincodeCollectionConfiguration) {
            lifecycleCommitChaincodeDefinitionRequest.setChaincodeCollectionConfiguration(chaincodeCollectionConfiguration);
        }
        lifecycleCommitChaincodeDefinitionRequest.setInitRequired(initRequired);

        Collection<LifecycleCommitChaincodeDefinitionProposalResponse> lifecycleCommitChaincodeDefinitionProposalResponses = channel.sendLifecycleCommitChaincodeDefinitionProposal(lifecycleCommitChaincodeDefinitionRequest,
                endorsingPeers);

        int numInstallProposal = lifecycleCommitChaincodeDefinitionProposalResponses.size();
        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();
        for (LifecycleCommitChaincodeDefinitionProposalResponse response : lifecycleCommitChaincodeDefinitionProposalResponses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                log.info("Successful commit proposal response Txid: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
            } else {
                log.error("commitChaincodeDefinitionRequest failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
                        response.getPeer(), response.getStatus(), response.isInvalid(), response.isVerified(), response.getMessage());
                failed.add(response);
            }
        }

        log.info("Received {} commit proposal responses. Successful. verified: {} . Failed: {}", numInstallProposal, successful.size(), failed.size());

        if (failed.size() > 0) {
//            ProposalResponse first = failed.iterator().next();
//            log.error("commitChaincodeDefinitionRequest failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
//                    first.getPeer(), first.getStatus(), first.isInvalid(), first.isVerified(), first.getMessage());
            throw new AgentException(ConstantCode.CHAIN_CODE_COMMIT_FAILED);
        }
        return channel.sendTransaction(lifecycleCommitChaincodeDefinitionProposalResponses);
    }

    public void verifyByQueryChaincodeDefinition(HFClient client, Channel channel, String chaincodeName, Collection<Peer> peers, long expectedSequence) throws ProposalException, InvalidArgumentException{

        final QueryLifecycleQueryChaincodeDefinitionRequest queryLifecycleQueryChaincodeDefinitionRequest = client.newQueryLifecycleQueryChaincodeDefinitionRequest();
        queryLifecycleQueryChaincodeDefinitionRequest.setChaincodeName(chaincodeName);

        Collection<LifecycleQueryChaincodeDefinitionProposalResponse> queryChaincodeDefinitionProposalResponses = channel.lifecycleQueryChaincodeDefinition(queryLifecycleQueryChaincodeDefinitionRequest, peers);

        for (LifecycleQueryChaincodeDefinitionProposalResponse response : queryChaincodeDefinitionProposalResponses) {
            if (!ChaincodeResponse.Status.SUCCESS.equals(response.getStatus()) || expectedSequence != response.getSequence()){
                throw new AgentException(ConstantCode.CHAIN_CODE_COMMIT_FAILED);
            }
        }
    }

    public long getSequence(HFClient client, Channel channel, Collection<Peer> peers, String chaincodeName) throws InvalidArgumentException, ProposalException {
        long sequence = -1L;
        final QueryLifecycleQueryChaincodeDefinitionRequest queryLifecycleQueryChaincodeDefinitionRequest = client.newQueryLifecycleQueryChaincodeDefinitionRequest();
        queryLifecycleQueryChaincodeDefinitionRequest.setChaincodeName(chaincodeName);

        Collection<LifecycleQueryChaincodeDefinitionProposalResponse> firstQueryDefininitions = channel.lifecycleQueryChaincodeDefinition(queryLifecycleQueryChaincodeDefinitionRequest, peers);

        for (LifecycleQueryChaincodeDefinitionProposalResponse firstDefinition : firstQueryDefininitions) {
            long seq = -1L;
            if (firstDefinition.getStatus() == ProposalResponse.Status.SUCCESS) {
                seq = firstDefinition.getSequence() + 1; //Need to bump it up to the next.
                log.info("Successful getSequence from peer:{} , seq={}", firstDefinition.getPeer().getName(), firstDefinition.getSequence());
            } else { //Failed but why?
                if (404 == firstDefinition.getChaincodeActionResponseStatus()) {
                    // not found .. done set sequence to 1;
                    sequence = 1;
                }
            }
            sequence = Math.max(seq, sequence);
        }
        return sequence;
    }

    public CompletableFuture<BlockEvent.TransactionEvent> executeChaincode(
            HFClient client, User userContext, Channel channel, String fcn, Boolean doInit, String chaincodeName,
            TransactionRequest.Type chaincodeType, String... args) throws InvalidArgumentException, ProposalException {

        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();

        TransactionProposalRequest transactionProposalRequest = client.newTransactionProposalRequest();
        transactionProposalRequest.setChaincodeName(chaincodeName);
        transactionProposalRequest.setChaincodeLanguage(chaincodeType);
        transactionProposalRequest.setUserContext(userContext);

        transactionProposalRequest.setFcn(fcn);
        transactionProposalRequest.setProposalWaitTime(TIMEOUT);
        transactionProposalRequest.setArgs(args);
        if (null != doInit) {
            transactionProposalRequest.setInit(doInit);
        }

        //  Collection<ProposalResponse> transactionPropResp = channel.sendTransactionProposalToEndorsers(transactionProposalRequest);
        Collection<Peer> peers = channel.getPeers();
        List<ChainCodePeerDo> chainCodePeers = chainCodeRepo.getChainCodePeers(channel.getName(), chaincodeName);
        if (chainCodePeers != null) {
            List<String> peerNames = chainCodePeers.stream().map(ChainCodePeerDo::getPeerName).collect(Collectors.toList());
            peers = peers.stream().filter(peer -> peerNames.contains(peer.getName())).collect(Collectors.toList());
        }
        Collection<ProposalResponse> transactionPropResp = channel.sendTransactionProposal(transactionProposalRequest, peers);
        for (ProposalResponse response : transactionPropResp) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                log.info("Successful transaction proposal response Txid: {} from peer {}", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
            } else {
                log.error("executeChaincode failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
                        response.getPeer(), response.getStatus(), response.isInvalid(), response.isVerified(), response.getMessage());
                failed.add(response);
            }
        }

        log.info("Received {} transaction proposal responses. Successful+verified: {} . Failed: {}",
                transactionPropResp.size(), successful.size(), failed.size());
        if (failed.size() > 0) {
//            ProposalResponse first = failed.iterator().next();
//            log.error("executeChaincode failed, peer={} status={} isInvalid={} isVerified={} errorMsg={}",
//                    first.getPeer(), first.getStatus(), first.isInvalid(), first.isVerified(), first.getMessage());
            throw new AgentException(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
        }
        return channel.sendTransaction(successful);

    }

    public byte[] queryChainCode(Channel channel, User peerAdmin, Collection<Peer> peers, String chaincodeName, String chaincodeVersion,
                                 TransactionRequest.Type chaincodeType, boolean initRequired, String function, String... args) throws InvalidArgumentException, ProposalException {
        QueryByChaincodeRequest request = QueryByChaincodeRequest.newInstance(peerAdmin);
        request.setChaincodeName(chaincodeName);
        request.setChaincodeVersion(chaincodeVersion);
        request.setChaincodeLanguage(chaincodeType);
        request.setFcn(function);
        request.setArgs(args);
        request.setInit(initRequired);
        request.setProposalWaitTime(TIMEOUT);
        Collection<ProposalResponse> proposalResponses = channel.queryByChaincode(request, peers);
        String errorMsg = "";
        for (ProposalResponse response : proposalResponses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                return response.getChaincodeActionResponsePayload();
            } else {
                errorMsg = String.format("failed transaction proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
            }
        }
        log.error("queryChainCode failed: {}", errorMsg);
        throw new AgentException(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
    }

}

