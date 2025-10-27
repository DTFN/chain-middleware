package com.lingshu.fabric.agent;

import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import org.apache.commons.compress.utils.IOUtils;
import org.bouncycastle.asn1.pkcs.PrivateKeyInfo;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.openssl.PEMParser;
import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter;
import org.hyperledger.fabric.gateway.*;
import org.hyperledger.fabric.protos.common.Configtx;
import org.hyperledger.fabric.protos.common.Ledger;
import org.hyperledger.fabric.protos.orderer.Configuration;
import org.hyperledger.fabric.sdk.*;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.ProposalException;

import java.io.*;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.security.NoSuchProviderException;
import java.security.PrivateKey;
import java.security.cert.CertificateException;
import java.security.cert.X509Certificate;
import java.security.spec.InvalidKeySpecException;
import java.text.SimpleDateFormat;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.TimeoutException;

import static java.lang.String.format;
import static java.nio.charset.StandardCharsets.UTF_8;

/**
 * @Author: zehao.song
 * @Date: 2023/11/2 16:23
 */
public class Test {

    private static final long DEPLOYWAITTIME = 3 * 60 * 1000;

    public static void main(String[] args) throws Exception {
        // part1 网络连接
        String networkConfigPath = "/home/songzehao/codes/baas2/fabric-agent/src/main/resources/connection-tls.json";
        // 使用org1中的admin初始化一个网关wallet账户用于连接网络
        String certificatePath = "/home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/cert.pem";
        X509Certificate certificate = readX509Certificate(Paths.get(certificatePath));
        String certificateStr = new String(IOUtils.toByteArray(new FileInputStream(certificatePath)), "UTF-8");

        String privateKeyPath = "/home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/93c9b1f1f3671190e61c82c7dab24acfd59cf202e4e351c3caeba196b6760df4_sk";
        PrivateKey privateKey = getPrivateKey(Paths.get(privateKeyPath));

        Wallet wallet = Wallets.newInMemoryWallet();
        wallet.put("admin", Identities.newX509Identity("Org1MSP", certificate, privateKey));

        // 根据connection-tls.json 获取Fabric网络连接对象
        Gateway.Builder builder = Gateway.createBuilder()
                .identity(wallet, "admin")
                .networkConfig(Paths.get(networkConfigPath));
        // 连接网关
        Gateway gateway = builder.connect();
        // 获取通道
        String channelName = "channel1";
        Network network = gateway.getNetwork(channelName);
        HFClient hfClient = gateway.getClient();
        Channel channel1 = hfClient.getChannel("channel1");
        // 获取通道的配置块
        byte[] channelConfigBytes = channel1.getChannelConfigurationBytes();
        // 从配置块中提取共识配置信息
        Configuration.ConsensusType consensusType = Configuration.ConsensusType.parseFrom(Configtx.Config.parseFrom(channelConfigBytes)
                .getChannelGroup()
                .getGroupsMap()
                .get("Orderer")
                .getValuesMap()
                .get("ConsensusType")
                .getValue());
        System.out.println("consensusType: " + consensusType.getType());


        Peer peer = channel1.getPeers().stream().findFirst().get();
        // peer channel list
        Set<String> channelSet = hfClient.queryChannels(peer);


//        deployJavaChaincode(hfClient);

        // part2 链码交易
        // 获取合约对象
//        Contract contract = network.getContract("basic");
//        // 查询现有资产
//        // peer chaincode invoke -o 192.168.3.128:7050 -C channel1 -n basic --peerAddresses 192.168.3.128:7051 --tlsRootCertFiles /home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --tls --cafile /home/songzehao/fabric/organizations/ordererOrganizations/org1.example.com/orderers/orderer0.org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem -c '{"function":"GetAllAssets","Args":[]}'
//        byte[] queryAllAssets = contract.evaluateTransaction("GetAllAssets");
//        System.out.println("更新前所有资产: "+ new String(queryAllAssets, UTF_8));
//
//        // 增加新的资产
//        // peer chaincode invoke -o 192.168.3.128:7050 -C channel1 -n basic --peerAddresses 192.168.3.128:7051 --tlsRootCertFiles /home/songzehao/fabric/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --tls --cafile /home/songzehao/fabric/organizations/ordererOrganizations/org1.example.com/orderers/orderer0.org1.example.com/msp/tlscacerts/tlsca.org1.example.com-cert.pem -c '{"function":"CreateAsset","Args":["asset8", "gold", "88", "Song", "888"]}'
////        byte[] invokeResult = contract.createTransaction("CreateAsset")
////                .setEndorsingPeers(network.getChannel().getPeers(EnumSet.of(Peer.PeerRole.ENDORSING_PEER)))
////                .submit("asset8", "gold", "88", "Song", "888");
////        System.out.println("新资产: " + new String(invokeResult, StandardCharsets.UTF_8));
//
//        // 查询更新后的资产
//        byte[] queryAllAssetsAfter = contract.evaluateTransaction("GetAllAssets");
//        System.out.println("更新后所有资产: " + new String(queryAllAssetsAfter, UTF_8));
//
//
//        // part3 区块交易查询
//        System.out.println("channelName: " + channel1.getName());
//        System.out.println("orderers: " + channel1.getOrderers());
//        System.out.println("peers: " + channel1.getPeers());
//        System.out.println("peersOrganizationMSPIDs: " + channel1.getPeersOrganizationMSPIDs());
//
//        BlockchainInfo blockChainInfo = channel1.queryBlockchainInfo();
//        Ledger.BlockchainInfo blockchainInfo = blockChainInfo.getBlockchainInfo();
//        // 查询最新区块hash
//        ByteString currentBlockHash = blockchainInfo.getCurrentBlockHash();
//        System.out.println("currentBlockHash: " + bytesToHex(currentBlockHash.toByteArray()));
//        // 查询最新区块高度
//        long height = blockchainInfo.getHeight();
//        System.out.println("height: " + height);
//
//
//        // 获取区块链的第一个区块（即创世块）
//        BlockInfo genesisBlock = channel1.queryBlockByNumber(0);
//        System.out.println("genesisBlock dataHash: " + bytesToHex(genesisBlock.getDataHash()));
//        // 获取创世块的区块头信息，并从中获取创世时间戳
//        BlockInfo.EnvelopeInfo envelopeInfo = genesisBlock.getEnvelopeInfo(0);
//        long genesisTime = envelopeInfo.getTimestamp().getTime();
//        System.out.println("genesisTime: " + genesisTime + " -> " + new SimpleDateFormat("yyyy-MM-dd HH:mm:ss.SSS").format(new Date(genesisTime)));
//        String genesisBlockHash = bytesToHex(channel1.queryBlockByNumber(1).getPreviousHash());
//        System.out.println("genesisBlockHash: " + genesisBlockHash);
//
//
//        // 根据区块高度查询区块
//        BlockInfo latestBlockInfo = channel1.queryBlockByNumber(height - 1);
//        // 根据区块hash查询区块
//        latestBlockInfo = channel1.queryBlockByHash(currentBlockHash.toByteArray());
//        // 获取区块中的所有交易
//        for (BlockInfo.EnvelopeInfo envelope : latestBlockInfo.getEnvelopeInfos()) {
//            // 判断是否为交易类型的封装
//            if (envelope.getType() == BlockInfo.EnvelopeType.TRANSACTION_ENVELOPE) {
//                // 获取交易ID
//                String transactionId = envelope.getTransactionID();
//                System.out.println("Transaction ID: " + transactionId);
//                // 根据交易ID查区块详情
//                channel1.queryBlockByTransactionID(transactionId);
//                // 根据交易ID查交易详情
//                TransactionInfo transactionInfo = channel1.queryTransactionByID(transactionId);
//            }
//        }
        FabricUser peerAdmin = new FabricUser();
        peerAdmin.setName("admin");
        peerAdmin.setMspId("Org1MSP");
        Enrollment enrollment = new Enrollment() {
            @Override
            public PrivateKey getKey() {
                return privateKey;
            }

            @Override
            public String getCert() {
                return certificateStr;
            }
        };
        peerAdmin.setEnrollment(enrollment);
        Collection<Peer> peers = channel1.getPeers();

        TransactionRequest.Type chaincodeType = TransactionRequest.Type.JAVA;
        String chaincodeName = "basic1";
        String chaincodeVersion = "1.0";
        LifecycleChaincodePackage lifecycleChaincodePackage = createLifecycleChaincodePackage(
                "basic_1", // some label
                chaincodeType,
                Paths.get("src",  "test", "resources", "basic").toFile().getPath(),
                "",
                "src/test/resources/meta-infs/end2endit");
//        TransactionRequest.Type chaincodeType = TransactionRequest.Type.GO_LANG;
//        String chaincodeName = "samples";
//        String chaincodeVersion = "3.0";
//        LifecycleChaincodePackage lifecycleChaincodePackage = createLifecycleChaincodePackage(
//                "samples", // some label
//                chaincodeType,
//                Paths.get("src", "test", "resources", "gocc", "sample1").toFile().getPath(),
//                "github.com/example_cc",
//                "src/test/resources/meta-infs/end2endit");
        boolean initRequired = false;
        LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy = LifecycleChaincodeEndorsementPolicy.fromSignaturePolicyYamlFile(Paths.get("src/test/resources/chaincodeendorsementpolicy.yaml"));
        runChainCodeInChannel(hfClient, channel1, peerAdmin, peers, lifecycleChaincodePackage, chaincodeName, chaincodeVersion,
                chaincodeEndorsementPolicy, null, initRequired);

        BlockEvent.TransactionEvent transactionEvent;
        //            transactionEvent = executeChaincode(client, peerAdmin, channel, "init",
//                    initRequired ? true : null, // doInit don't even specify it has it should default to false
//                    chaincodeName, chaincodeType, "a,", "100", "b", "300").get(DEPLOYWAITTIME, TimeUnit.SECONDS);

//            if (initRequired) {
//                assertTrue(transactionEvent.getTransactionActionInfo(0).getChaincodeInputIsInit());
//            }

//            transactionEvent = executeChaincode(client, peerAdmin, channel, "move",
//                    initRequired, // doInit
//                    chaincodeName, chaincodeType, "a,", "b", "10").get(DEPLOYWAITTIME, TimeUnit.SECONDS);

        transactionEvent = executeChaincode(hfClient, peerAdmin, channel1, "InitLedger",
                initRequired, // doInit
                chaincodeName, chaincodeType, "").get(DEPLOYWAITTIME, TimeUnit.SECONDS);
        out(transactionEvent.toString());

        QueryByChaincodeRequest request = QueryByChaincodeRequest.newInstance(peerAdmin);
        request.setChaincodeName(chaincodeName);
        request.setChaincodeVersion(chaincodeVersion);
        request.setChaincodeLanguage(chaincodeType);
        request.setFcn("GetAllAssets");
        request.setArgs("");
        request.setInit(initRequired);
        request.setProposalWaitTime(DEPLOYWAITTIME);
        Collection<ProposalResponse> proposalResponses = channel1.queryByChaincode(request, peers);
        for (ProposalResponse response : proposalResponses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                out("Successful transaction proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
            } else {
                out("failed transaction proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
            }
        }
    }

    static void runChainCodeInChannel(HFClient client, Channel channel, FabricUser peerAdmin, Collection<Peer> peers,
                    LifecycleChaincodePackage lifecycleChaincodePackage, String chaincodeName,
                    String chaincodeVersion, LifecycleChaincodeEndorsementPolicy lifecycleChaincodeEndorsementPolicy,
                    ChaincodeCollectionConfiguration chaincodeCollectionConfiguration,  boolean initRequired) {

        try {
            //Should be no chaincode installed at this time.
            client.setUserContext(peerAdmin);

            final String chaincodeLabel = lifecycleChaincodePackage.getLabel();
            final TransactionRequest.Type chaincodeType = lifecycleChaincodePackage.getType();

            //Org1 installs the chaincode on its peers.
            out("Org1 installs the chaincode on its peers.");
            String org1ChaincodePackageID = lifecycleInstallChaincode(client, peers, lifecycleChaincodePackage);

            // Sequence  number increase with each change and is used to make sure you are referring to the same change.
            long sequence = -1L;
            final QueryLifecycleQueryChaincodeDefinitionRequest queryLifecycleQueryChaincodeDefinitionRequest = client.newQueryLifecycleQueryChaincodeDefinitionRequest();
            queryLifecycleQueryChaincodeDefinitionRequest.setChaincodeName(chaincodeName);

            Collection<LifecycleQueryChaincodeDefinitionProposalResponse> firstQueryDefininitions = channel.lifecycleQueryChaincodeDefinition(queryLifecycleQueryChaincodeDefinitionRequest, peers);

            for (LifecycleQueryChaincodeDefinitionProposalResponse firstDefinition : firstQueryDefininitions) {
                if (firstDefinition.getStatus() == ProposalResponse.Status.SUCCESS) {
                    sequence = firstDefinition.getSequence() + 1L; //Need to bump it up to the next.
                    break;
                } else { //Failed but why?
                    if (404 == firstDefinition.getChaincodeActionResponseStatus()) {
                        // not found .. done set sequence to 1;
                        sequence = 1;
                        break;
                    }
                }
            }

//            if (null != expected) {
//                assertEquals(expected.get("sequence"), sequence);
//            }

            //     ChaincodeCollectionConfiguration chaincodeCollectionConfiguration = collectionConfiguration == null ? null : ChaincodeCollectionConfiguration.fromYamlFile(new File(collectionConfiguration));
//            // ChaincodeCollectionConfiguration chaincodeCollectionConfiguration = ChaincodeCollectionConfiguration.fromYamlFile(new File("src/test/fixture/collectionProperties/PrivateDataIT.yaml"));
//            chaincodeCollectionConfiguration = null;
            final Peer anOrg1Peer = peers.iterator().next();
            out("Org1 approving chaincode definition for my org.");
            BlockEvent.TransactionEvent transactionEvent = lifecycleApproveChaincodeDefinitionForMyOrg(client, channel,
                    Collections.singleton(anOrg1Peer),  //support approve on multiple peers but really today only need one. Go with minimum.
                    sequence, chaincodeName, chaincodeVersion, lifecycleChaincodeEndorsementPolicy, chaincodeCollectionConfiguration, initRequired, org1ChaincodePackageID)
                    .get(DEPLOYWAITTIME, TimeUnit.SECONDS);
//            assertTrue(transactionEvent.isValid());

//            verifyByCheckCommitReadinessStatus(client, channel, sequence, chaincodeName, chaincodeVersion,
//                    lifecycleChaincodeEndorsementPolicy, chaincodeCollectionConfiguration, initRequired, peers,
//                    new HashSet<>(Arrays.asList(ORG_1_MSP)), // Approved
//                    new HashSet<>(Arrays.asList(ORG_2_MSP))); // Un approved.

            // Get collection of one of org2 orgs peers and one from the other.
            out("Org2 doing commit chaincode definition");
            transactionEvent = commitChaincodeDefinitionRequest(client, channel, sequence, chaincodeName, chaincodeVersion, lifecycleChaincodeEndorsementPolicy, chaincodeCollectionConfiguration, initRequired, peers)
                    .get(DEPLOYWAITTIME, TimeUnit.SECONDS);
//            assertTrue(transactionEvent.isValid());

            out("Org2 done with commit. block #%d!", transactionEvent.getBlockEvent().getBlockNumber());

//            verifyByQueryChaincodeDefinition(org2Client, org2Channel, chaincodeName, org2MyPeers, sequence, initRequired, chaincodeEndorsementPolicyAsBytes, chaincodeCollectionConfiguration);
//            verifyByQueryChaincodeDefinition(client, channel, chaincodeName, peers, sequence, initRequired, chaincodeEndorsementPolicyAsBytes, chaincodeCollectionConfiguration);
//
//            verifyByQueryChaincodeDefinitions(org2Client, org2Channel, org2MyPeers, chaincodeName);
//            verifyByQueryChaincodeDefinitions(client, channel, peers, chaincodeName);


        } catch (Exception e) {
            out("Caught an exception running channel %s", channel.getName());
            e.printStackTrace();
            fail("Test failed with error : " + e.getMessage());
        }
    }

    private static LifecycleChaincodePackage createLifecycleChaincodePackage(String chaincodeLabel, TransactionRequest.Type chaincodeType, String chaincodeSourceLocation, String chaincodePath, String metadadataSource) throws IOException, InvalidArgumentException {
        out("creating install package %s.", chaincodeLabel);

        Path metadataSourcePath = null;
        if (metadadataSource != null) {
            metadataSourcePath = Paths.get(metadadataSource);
        }
        LifecycleChaincodePackage lifecycleChaincodePackage = LifecycleChaincodePackage.fromSource(chaincodeLabel, Paths.get(chaincodeSourceLocation),
                chaincodeType,
                chaincodePath, metadataSourcePath);
//        assertEquals(chaincodeLabel, lifecycleChaincodePackage.getLabel()); // what we expect ?
//        assertEquals(chaincodeType, lifecycleChaincodePackage.getType());
//        assertEquals(chaincodePath, lifecycleChaincodePackage.getPath());

        return lifecycleChaincodePackage;
    }

    private static String lifecycleInstallChaincode(HFClient client, Collection<Peer> peers, LifecycleChaincodePackage lifecycleChaincodePackage) throws InvalidArgumentException, ProposalException, InvalidProtocolBufferException, ProposalException, InvalidArgumentException {

        int numInstallProposal = 0;

        numInstallProposal = numInstallProposal + peers.size();

        LifecycleInstallChaincodeRequest installProposalRequest = client.newLifecycleInstallChaincodeRequest();
        installProposalRequest.setLifecycleChaincodePackage(lifecycleChaincodePackage);
        installProposalRequest.setProposalWaitTime(DEPLOYWAITTIME);

        Collection<LifecycleInstallChaincodeProposalResponse> responses = client.sendLifecycleInstallChaincodeRequest(installProposalRequest, peers);
//        assertNotNull(responses);

        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();
        String packageID = null;
        for (LifecycleInstallChaincodeProposalResponse response : responses) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                out("Successful install proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
                if (packageID == null) {
                    packageID = response.getPackageId();
//                    assertNotNull(format("Hashcode came back as null from peer: %s ", response.getPeer()), packageID);
                } else {
//                    assertEquals("Miss match on what the peers returned back as the packageID", packageID, response.getPackageId());
                }
            } else {
                failed.add(response);
            }
        }

        //   }
        out("Received %d install proposal responses. Successful+verified: %d . Failed: %d", numInstallProposal, successful.size(), failed.size());

        if (failed.size() > 0) {
            ProposalResponse first = failed.iterator().next();
            fail("Not enough endorsers for install :" + successful.size() + ".  " + first.getMessage());
        }

        return packageID;

    }

    static CompletableFuture<BlockEvent.TransactionEvent> lifecycleApproveChaincodeDefinitionForMyOrg(HFClient client, Channel channel,
                                                                                               Collection<Peer> peers, long sequence,
                                                                                               String chaincodeName, String chaincodeVersion, LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy, ChaincodeCollectionConfiguration chaincodeCollectionConfiguration, boolean initRequired, String org1ChaincodePackageID) throws InvalidArgumentException, ProposalException {

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

        lifecycleApproveChaincodeDefinitionForMyOrgRequest.setPackageId(org1ChaincodePackageID);

        Collection<LifecycleApproveChaincodeDefinitionForMyOrgProposalResponse> lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse = channel.sendLifecycleApproveChaincodeDefinitionForMyOrgProposal(lifecycleApproveChaincodeDefinitionForMyOrgRequest,
                peers);

//        assertEquals(peers.size(), lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse.size());
        for (LifecycleApproveChaincodeDefinitionForMyOrgProposalResponse response : lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse) {
            final Peer peer = response.getPeer();

//            assertEquals(format("failure on %s  message is: %s", peer, response.getMessage()), ChaincodeResponse.Status.SUCCESS, response.getStatus());
//            assertFalse(peer + " " + response.getMessage(), response.isInvalid());
//            assertTrue(format("failure on %s", peer), response.isVerified());
        }

        return channel.sendTransaction(lifecycleApproveChaincodeDefinitionForMyOrgProposalResponse);

    }

    private static CompletableFuture<BlockEvent.TransactionEvent> commitChaincodeDefinitionRequest(HFClient client, Channel channel, long definitionSequence, String chaincodeName, String chaincodeVersion,
                                                                                            LifecycleChaincodeEndorsementPolicy chaincodeEndorsementPolicy,
                                                                                            ChaincodeCollectionConfiguration chaincodeCollectionConfiguration,
                                                                                            boolean initRequired, Collection<Peer> endorsingPeers) throws ProposalException, InvalidArgumentException, InterruptedException, ExecutionException, TimeoutException {
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

        for (LifecycleCommitChaincodeDefinitionProposalResponse resp : lifecycleCommitChaincodeDefinitionProposalResponses) {

            final Peer peer = resp.getPeer();
//            assertEquals(format("%s had unexpected status.", peer.toString()), ChaincodeResponse.Status.SUCCESS, resp.getStatus());
//            assertTrue(format("%s not verified.", peer.toString()), resp.isVerified());
        }

        return channel.sendTransaction(lifecycleCommitChaincodeDefinitionProposalResponses);

    }

    static CompletableFuture<BlockEvent.TransactionEvent> executeChaincode(HFClient client, User userContext, Channel channel, String fcn, Boolean doInit, String chaincodeName, TransactionRequest.Type chaincodeType, String... args) throws InvalidArgumentException, ProposalException {

        final ExecutionException[] executionExceptions = new ExecutionException[1];

        Collection<ProposalResponse> successful = new LinkedList<>();
        Collection<ProposalResponse> failed = new LinkedList<>();

        TransactionProposalRequest transactionProposalRequest = client.newTransactionProposalRequest();
        transactionProposalRequest.setChaincodeName(chaincodeName);
        transactionProposalRequest.setChaincodeLanguage(chaincodeType);
        transactionProposalRequest.setUserContext(userContext);

        transactionProposalRequest.setFcn(fcn);
        transactionProposalRequest.setProposalWaitTime(DEPLOYWAITTIME);
        transactionProposalRequest.setArgs(args);
        if (null != doInit) {
            transactionProposalRequest.setInit(doInit);
        }

        //  Collection<ProposalResponse> transactionPropResp = channel.sendTransactionProposalToEndorsers(transactionProposalRequest);
        Collection<ProposalResponse> transactionPropResp = channel.sendTransactionProposal(transactionProposalRequest, channel.getPeers());
        for (ProposalResponse response : transactionPropResp) {
            if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                out("Successful transaction proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
                successful.add(response);
            } else {
                failed.add(response);
            }
        }

        out("Received %d transaction proposal responses. Successful+verified: %d . Failed: %d",
                transactionPropResp.size(), successful.size(), failed.size());
        if (failed.size() > 0) {
            ProposalResponse firstTransactionProposalResponse = failed.iterator().next();
            fail("Not enough endorsers for executeChaincode(move a,b,100):" + failed.size() + " endorser error: " +
                    firstTransactionProposalResponse.getMessage() +
                    ". Was verified: " + firstTransactionProposalResponse.isVerified());
        }
        out("Successfully received transaction proposal responses.");

        //  System.exit(10);

        ////////////////////////////
        // Send Transaction Transaction to orderer
        out("Sending chaincode transaction(move a,b,100) to orderer.");
        return channel.sendTransaction(successful);

    }

    private static ChaincodeID deployJavaChaincode(HFClient client) {
        String channelId = "channel1";
        Channel channel = client.getChannel(channelId);
        out("deployChaincode - enter");
        ChaincodeID chaincodeID = null;
        String ccName = "basic";
        String ccVersion = "2.7";
        String ccPath = "github.com/example_cc";
        try {

            final String channelName = channel.getName();
            out("deployChaincode - channelName = " + channelName);

            Collection<Orderer> orderers = channel.getOrderers();
            Collection<ProposalResponse> responses;
            Collection<ProposalResponse> successful = new LinkedList<>();
            Collection<ProposalResponse> failed = new LinkedList<>();

            chaincodeID = ChaincodeID.newBuilder().setName(ccName)
                    .setVersion(ccVersion)
                    .setPath(ccPath).build();

            ////////////////////////////
            // Install Proposal Request
            //
            out("Creating install proposal");

            InstallProposalRequest installProposalRequest = client.newInstallProposalRequest();
            installProposalRequest.setChaincodeID(chaincodeID);

            ////For GO language and serving just a single user, chaincodeSource is mostly likely the users GOPATH
            installProposalRequest.setChaincodeSourceLocation(Paths.get("src",  "test", "resources", "gocc", "sample1").toFile());

            installProposalRequest.setChaincodeVersion(ccVersion);

            out("Sending install proposal");

            ////////////////////////////
            // only a client from the same org as the peer can issue an install request
            int numInstallProposal = 0;

            Collection<Peer> peersFromOrg = channel.getPeers();
            numInstallProposal = numInstallProposal + peersFromOrg.size();
            responses = client.sendInstallProposal(installProposalRequest, peersFromOrg);

            for (ProposalResponse response : responses) {
                if (response.getStatus() == ProposalResponse.Status.SUCCESS) {
                    out("Successful install proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
                    successful.add(response);
                } else {
                    failed.add(response);
                }
            }

            out("Received %d install proposal responses. Successful+verified: %d . Failed: %d", numInstallProposal, successful.size(), failed.size());

            if (failed.size() > 0) {
                ProposalResponse first = failed.iterator().next();
                fail("Not enough endorsers for install :" + successful.size() + ".  " + first.getMessage());
            }

            ///////////////
            //// Instantiate chaincode.
            //
            // From the docs:
            // The instantiate transaction invokes the lifecycle System Chaincode (LSCC) to create and initialize a chaincode on a channel
            // After being successfully instantiated, the chaincode enters the active state on the channel and is ready to process any transaction proposals of type ENDORSER_TRANSACTION

            InstantiateProposalRequest instantiateProposalRequest = client.newInstantiationProposalRequest();
            instantiateProposalRequest.setProposalWaitTime(3 * 60 * 1000);
            instantiateProposalRequest.setChaincodeID(chaincodeID);
            instantiateProposalRequest.setFcn("init");
            instantiateProposalRequest.setArgs("a", "500", "b", "999");

            Map<String, byte[]> tm = new HashMap<>();
            tm.put("HyperLedgerFabric", "InstantiateProposalRequest:JavaSDK".getBytes(UTF_8));
            tm.put("method", "InstantiateProposalRequest".getBytes(UTF_8));
            instantiateProposalRequest.setTransientMap(tm);

            /*
              policy OR(Org1MSP.member, Org2MSP.member) meaning 1 signature from someone in either Org1 or Org2
              See README.md Chaincode endorsement policies section for more details.
            */
            ChaincodeEndorsementPolicy chaincodeEndorsementPolicy = new ChaincodeEndorsementPolicy();
            chaincodeEndorsementPolicy.fromYamlFile(new File( "src/test/resources/chaincodeendorsementpolicy.yaml"));
            instantiateProposalRequest.setChaincodeEndorsementPolicy(chaincodeEndorsementPolicy);

            out("Sending instantiateProposalRequest to all peers...");
            successful.clear();
            failed.clear();

            responses = channel.sendInstantiationProposal(instantiateProposalRequest);

            for (ProposalResponse response : responses) {
                if (response.isVerified() && response.getStatus() == ProposalResponse.Status.SUCCESS) {
                    successful.add(response);
                    out("Succesful instantiate proposal response Txid: %s from peer %s", response.getTransactionID(), response.getPeer().getName());
                } else {
                    failed.add(response);
                }
            }
            out("Received %d instantiate proposal responses. Successful+verified: %d . Failed: %d", responses.size(), successful.size(), failed.size());
            if (failed.size() > 0) {
                ProposalResponse first = failed.iterator().next();
                fail("Not enough endorsers for instantiate :" + successful.size() + "endorser failed with " + first.getMessage() + ". Was verified:" + first.isVerified());
            }

            ///////////////
            /// Send instantiate transaction to orderer
            out("Sending instantiateTransaction to orderer...");
            CompletableFuture<BlockEvent.TransactionEvent> future = channel.sendTransaction(successful, orderers);

            out("calling get...");
            BlockEvent.TransactionEvent event = future.get(30, TimeUnit.SECONDS);
            out("get done...");

//            Assert(event.isValid()); // must be valid to be here.
            out("Finished instantiate transaction with transaction id %s", event.getTransactionID());

        } catch (Exception e) {
            e.printStackTrace();
            out("Caught an exception running channel %s", channel.getName());
            fail("Test failed with error : " + e.getMessage());
        }

        return chaincodeID;
    }

    private static void out(String format, Object... args) {

        System.err.flush();
        System.out.flush();

        System.out.println(format(format, args));
        System.err.flush();
        System.out.flush();

    }

    public static void fail(String message) {
        if (message == null) {
            throw new AssertionError();
        } else {
            throw new AssertionError(message);
        }
    }

    // 将字节数组转换为十六进制字符串
    private static String bytesToHex(byte[] bytes) {
        StringBuilder sb = new StringBuilder();
        for (byte b : bytes) {
            sb.append(String.format("%02x", b));
        }
        return sb.toString();
    }

    private static X509Certificate readX509Certificate(final Path certificatePath) throws IOException, CertificateException {
        try (Reader certificateReader = Files.newBufferedReader(certificatePath, UTF_8)) {
            return Identities.readX509Certificate(certificateReader);
        }
    }

    private static PrivateKey getPrivateKey(final Path privateKeyPath) throws IOException, InvalidKeyException {
        try (Reader privateKeyReader = Files.newBufferedReader(privateKeyPath, UTF_8)) {
            return Identities.readPrivateKey(privateKeyReader);
        }
    }


}
