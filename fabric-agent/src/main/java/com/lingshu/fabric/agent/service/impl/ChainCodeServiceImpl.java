package com.lingshu.fabric.agent.service.impl;

import com.alibaba.fastjson2.JSON;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.event.GatewayInitedEvent;
import com.lingshu.fabric.agent.exception.AgentException;
import com.lingshu.fabric.agent.repo.ChainCodeRepo;
import com.lingshu.fabric.agent.req.chaincode.*;
import com.lingshu.fabric.agent.resp.FabricChainCodeInvokeResp;
import com.lingshu.fabric.agent.resp.PackageInfoDTO;
import com.lingshu.fabric.agent.service.ChainCodeService;
import com.lingshu.fabric.agent.service.FabricSdkService;
import com.lingshu.fabric.agent.service.PathService;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.FileUtils;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationListener;
import org.springframework.stereotype.Service;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import java.io.File;
import java.io.IOException;
import java.util.*;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;
import java.util.stream.Collectors;

@Service
@Slf4j
public class ChainCodeServiceImpl implements ChainCodeService, ApplicationListener<GatewayInitedEvent> {

    @Value("${chaincode.timeout:60000}")
    private long TIMEOUT = 60 * 1000;
    private static final String TEMP_DIR = "temp/";
    private static final String SUFFIX = ".tar.gz";

    @Autowired
    private ConstantProperties constantProperties;
    @Autowired
    private PathService pathService;
    @Autowired(required = false)
    private Gateway fabricGateway;
    @Autowired
    private FabricSdkService fabricSdkService;
    @Autowired
    private ChainCodeRepo chainCodeRepo;

    @Override
    public void onApplicationEvent(GatewayInitedEvent event) {
        if (fabricGateway != null){
            fabricGateway.close();
        }
        fabricGateway = event.getGateway();
        log.info("FabricSdkService found new gateway: {}", fabricGateway);
    }

    @Override
    public PackageInfoDTO resolveChaincodePackage(MultipartFile chainCodeFile){
        LifecycleChaincodePackage lifecycleChaincodePackage = null;
        try {
            lifecycleChaincodePackage = LifecycleChaincodePackage.fromBytes(chainCodeFile.getBytes());
            TransactionRequest.Type type = lifecycleChaincodePackage.getType();
            String label = lifecycleChaincodePackage.getLabel();
            String filePath = saveChainCodeFile(chainCodeFile.getBytes(), label);
            return new PackageInfoDTO(type.name(), label, filePath);
        } catch (IOException | InvalidArgumentException e) {
            log.error("resolveChaincodePackage failed", e);
            throw new AgentException(ConstantCode.INVALID_CHAIN_CODE_PACKAGE);
        }
    }

    @Override
    public String installChainCode(FabricCCInstallReq req) throws InvalidArgumentException, ProposalException, IOException {
        Network network = fabricGateway.getNetwork(req.getChannelId());
        HFClient client = fabricGateway.getClient();
        Channel channel = network.getChannel();
        Collection<Peer> peers = channel.getPeers();

        LifecycleChaincodePackage lifecycleChaincodePackage = LifecycleChaincodePackage.fromFile(new File(req.getChainCodePath()));
        String label = lifecycleChaincodePackage.getLabel();
        Map<String, String> installedChainCodes = fabricSdkService.lifecycleQueryInstalledChainCodes(client, label, peers);
        String packageId = null;
        if (installedChainCodes != null && !installedChainCodes.isEmpty()){
            Set<String> peerSet = installedChainCodes.keySet();
            peers = peers.stream().filter(peer -> !peerSet.contains(peer.getName())).collect(Collectors.toList());
            Optional<String> first = installedChainCodes.values().stream().findFirst();
            if (first.isPresent()){
                packageId = first.get();
            }
            log.info("label [{}] has installed on peers: {}, packageId={}", label, peerSet, packageId);
        }

        if (!CollectionUtils.isEmpty(peers)){
            log.info("label [{}] start to install on peers: {}", label, peers.stream().map(Peer::getName).collect(Collectors.toList()));
            packageId = fabricSdkService.lifecycleInstallChaincode(client, peers, lifecycleChaincodePackage);
        }
        return packageId;
    }

    @Override
    public long approveChainCode(FabricCCApproveReq req) throws InvalidArgumentException, ProposalException, ExecutionException, InterruptedException, TimeoutException {
        Network network = fabricGateway.getNetwork(req.getChannelId());
        HFClient client = fabricGateway.getClient();
        Channel channel = network.getChannel();
        Collection<Peer> peers = channel.getPeers();
        String packageId = req.getPackageId();
        long sequence = fabricSdkService.getSequence(client, channel, peers, req.getChainCodeName());
        LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy = null;

        List<String> approvedPeers = fabricSdkService.queryApprovedPeers(client, channel, sequence, req.getChainCodeName(), req.getChainCodeVersion(),
                chaincodeEndorsementPolicy, null, req.isInitRequired(), peers);
        if (!CollectionUtils.isEmpty(approvedPeers)){
            if (peers.size() != approvedPeers.size()){
                List<String> allPeers = peers.stream().map(Peer::getName).collect(Collectors.toList());
                log.warn("commitChainCode failed approvedPeers={}, allPeers={}", approvedPeers, allPeers);
                throw new AgentException(ConstantCode.CHAIN_CODE_APPROVE_FAILED);
            }
            return sequence;
        }

        try {
            BlockEvent.TransactionEvent transactionEvent = fabricSdkService.lifecycleApproveChaincodeDefinitionForMyOrg(client, channel, peers, sequence, req.getChainCodeName(), req.getChainCodeVersion(),
                            chaincodeEndorsementPolicy, null, req.isInitRequired(), packageId)
                    .get(TIMEOUT, TimeUnit.MILLISECONDS);
            if (!transactionEvent.isValid()) {
                log.error("lifecycleApproveChaincodeDefinitionForMyOrg failed: {}", transactionEvent);
                throw new AgentException(ConstantCode.CHAIN_CODE_APPROVE_FAILED);
            }
        } catch (TimeoutException e){
            log.warn("approveChainCode timeout req={}", req);
        }
        return sequence;
    }

    @Override
    public void commitChainCode(FabricCCCommitReq req) throws InvalidArgumentException, ProposalException, ExecutionException, InterruptedException, TimeoutException {
        Network network = fabricGateway.getNetwork(req.getChannelId());
        HFClient client = fabricGateway.getClient();
        Channel channel = network.getChannel();
        Collection<Peer> peers = channel.getPeers();
        long sequence = req.getSequence();
        LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy = null;

        List<String> approvedPeers = fabricSdkService.queryApprovedPeers(client, channel, sequence, req.getChainCodeName(), req.getChainCodeVersion(),
                chaincodeEndorsementPolicy, null, req.isInitRequired(), peers);
        if (CollectionUtils.isEmpty(approvedPeers)){
            throw new AgentException(ConstantCode.CHAIN_CODE_NOT_APPROVED);
        }
        if (peers.size() > approvedPeers.size()){
            List<String> allPeers = peers.stream().map(Peer::getName).collect(Collectors.toList());
            log.warn("commitChainCode failed approvedPeers={}, allPeers={}", approvedPeers, allPeers);
            throw new AgentException(ConstantCode.CHAIN_CODE_NOT_APPROVED);
        }
        log.info("commitChainCode approvedPeers={}", approvedPeers);

        try {
            BlockEvent.TransactionEvent transactionEvent = fabricSdkService.commitChaincodeDefinitionRequest(client, channel, sequence, req.getChainCodeName(), req.getChainCodeVersion(),
                            chaincodeEndorsementPolicy, null, req.isInitRequired(), peers)
                    .get(TIMEOUT, TimeUnit.MILLISECONDS);
            if (!transactionEvent.isValid()) {
                log.error("commitChaincodeDefinitionRequest failed: {}", transactionEvent);
                throw new AgentException(ConstantCode.CHAIN_CODE_COMMIT_FAILED);
            }
        } catch (TimeoutException e){
            log.warn("commitChainCode timeout req={}", req);
        }
        fabricSdkService.verifyByQueryChaincodeDefinition(client, channel, req.getChainCodeName(), peers, sequence);
        chainCodeRepo.save(req.getChannelId(), req.getChainCodeName(), req.getLang(), peers);
    }

    @Override
    public FabricChainCodeInvokeResp invokeChainCode(FabricCCInvokeReq req) {
        Network network = fabricGateway.getNetwork(req.getChannelId());
        HFClient client = fabricGateway.getClient();
        Channel channel = network.getChannel();
        TransactionRequest.Type chainCodeType = TransactionRequest.Type.fromPackageName(req.getLang());
        List<String> args = req.getArgs();
        String[] params;
        if (CollectionUtils.isEmpty(args)) {
            params = new String[]{""};
        } else {
            params = args.toArray(new String[args.size()]);
        }
        try {
            BlockEvent.TransactionEvent transactionEvent = fabricSdkService.executeChaincode(client, client.getUserContext(), channel, req.getFunction(),
                    req.isInit(), req.getChainCodeName(), chainCodeType, params).get(TIMEOUT, TimeUnit.MILLISECONDS);
            BlockInfo.TransactionEnvelopeInfo.TransactionActionInfo transactionActionInfo = transactionEvent.getTransactionActionInfo(0);
            int status = transactionActionInfo.getProposalResponseStatus();
            byte[] payload = transactionActionInfo.getProposalResponsePayload();
            String result = new String(payload);
            return new FabricChainCodeInvokeResp(status, result);
        } catch (Exception e) {
            log.error("invokeChainCode failed", e);
            throw new AgentException(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
        }
    }

    @Override
    public Object queryChainCode(FabricCCInvokeReq req) {
        Network network = fabricGateway.getNetwork(req.getChannelId());
        HFClient client = fabricGateway.getClient();
        Channel channel = network.getChannel();
        TransactionRequest.Type chainCodeType = TransactionRequest.Type.fromPackageName(req.getLang());
        List<String> args = req.getArgs();
        String[] params;
        if (CollectionUtils.isEmpty(args)) {
            params = new String[]{""};
        } else {
            params = args.toArray(new String[args.size()]);
        }
        try {
            byte[] bytes = fabricSdkService.queryChainCode(channel, client.getUserContext(), channel.getPeers(), req.getChainCodeName(),
                    req.getChainCodeVersion(), chainCodeType, req.isInit(), req.getFunction(), params);
            String result = new String(bytes);
            return JSON.parse(result);
        } catch (Exception e) {
            log.error("queryChainCode failed", e);
            throw new AgentException(ConstantCode.EXEC_INVOKE_CHAIN_CODE_ERROR);
        }
    }

    private String saveChainCodeFile(byte[] chaincodeBytes, String chainCodeName){
        String fileName = chainCodeName + SUFFIX;
        File directory = new File(TEMP_DIR);
        if (!directory.exists()) {
            directory.mkdirs();
        }

        File tempFile = new File(TEMP_DIR + fileName);
        try {
            FileUtils.writeByteArrayToFile(tempFile, chaincodeBytes);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
        return tempFile.getAbsolutePath();
    }
}
