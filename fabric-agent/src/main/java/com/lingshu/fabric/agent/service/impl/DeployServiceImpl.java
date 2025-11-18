package com.lingshu.fabric.agent.service.impl;

import com.alibaba.fastjson2.JSON;
import com.alibaba.fastjson2.JSONArray;
import com.alibaba.fastjson2.JSONObject;
import com.baomidou.mybatisplus.core.conditions.query.LambdaQueryWrapper;
import com.baomidou.mybatisplus.core.conditions.update.LambdaUpdateWrapper;
import com.google.common.collect.Lists;
import com.google.protobuf.InvalidProtocolBufferException;
import com.lingshu.fabric.agent.bo.NetworkConfig;
import com.lingshu.fabric.agent.bo.NodeDetail;
import com.lingshu.fabric.agent.config.FabricConfig;
import com.lingshu.fabric.agent.config.code.ConstantCode;
import com.lingshu.fabric.agent.config.properties.ConstantProperties;
import com.lingshu.fabric.agent.enums.FabricDeployNodeStage;
import com.lingshu.fabric.agent.enums.FabricDeployStage;
import com.lingshu.fabric.agent.enums.FabricNodeType;
import com.lingshu.fabric.agent.enums.ScpTypeEnum;
import com.lingshu.fabric.agent.event.GatewayInitedEvent;
import com.lingshu.fabric.agent.exception.AgentException;
import com.lingshu.fabric.agent.repo.ChainInfoRepo;
import com.lingshu.fabric.agent.repo.DeployResultNodeRepo;
import com.lingshu.fabric.agent.repo.DeployResultRepo;
import com.lingshu.fabric.agent.repo.OrgInfoRepo;
import com.lingshu.fabric.agent.repo.entity.ChainInfoDo;
import com.lingshu.fabric.agent.repo.entity.DeployResultDo;
import com.lingshu.fabric.agent.repo.entity.DeployResultNodeDo;
import com.lingshu.fabric.agent.req.*;
import com.lingshu.fabric.agent.resp.*;
import com.lingshu.fabric.agent.service.*;
import com.lingshu.fabric.agent.util.ExecuteResult;
import com.lingshu.fabric.agent.util.JavaCommandExecutor;
import com.lingshu.fabric.agent.util.TemplateUtil;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.io.FileUtils;
import org.apache.commons.lang3.tuple.Pair;
import org.bouncycastle.asn1.ASN1Integer;
import org.bouncycastle.asn1.DEROctetString;
import org.bouncycastle.asn1.DERSequenceGenerator;
import org.bouncycastle.crypto.Digest;
import org.bouncycastle.crypto.digests.SHA256Digest;
import org.bouncycastle.crypto.digests.SHA3Digest;
import org.hyperledger.fabric.gateway.Gateway;
import org.hyperledger.fabric.gateway.Network;
import org.hyperledger.fabric.protos.common.Common;
import org.hyperledger.fabric.protos.common.Configtx;
import org.hyperledger.fabric.protos.orderer.Configuration;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.HFClient;
import org.hyperledger.fabric.sdk.exception.InvalidArgumentException;
import org.hyperledger.fabric.sdk.exception.TransactionException;
import org.jetbrains.annotations.NotNull;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.ApplicationListener;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.InputStreamSource;
import org.springframework.scheduling.annotation.Async;
import org.springframework.scheduling.concurrent.ThreadPoolTaskScheduler;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.Assert;
import org.springframework.util.CollectionUtils;
import org.springframework.util.StreamUtils;
import org.yaml.snakeyaml.DumperOptions;
import org.yaml.snakeyaml.Yaml;

import java.io.*;
import java.lang.reflect.Field;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@Service
@Slf4j
public class DeployServiceImpl implements DeployService, ApplicationListener<GatewayInitedEvent> {
    @Autowired(required = false)
    private Gateway gateway;
    @Autowired
    private AnsibleService ansibleService;
    @Autowired
    private MonitorService monitorService;
    @Autowired
    private PathService pathService;
    @Autowired
    private ConstantProperties constantProperties;
    @Autowired
    private DeployResultRepo deployResultRepo;
    @Autowired
    private DeployResultNodeRepo deployResultNodeRepo;
    @Autowired
    private ChainInfoRepo chainInfoRepo;
    @Autowired
    private NodeInfoService nodeInfoService;
    @Autowired
    private OrgInfoRepo orgInfoRepo;

    @Autowired
    @Qualifier(value = "initAsyncScheduler")
    private ThreadPoolTaskScheduler threadPoolTaskScheduler;

    @Value("${fabricImage:busybox}")
    private String DOCKER_IMAGE;

    @Value("${fabric.log-level:INFO}")
    private String fabricLogLevel;
    @Value("${fabric.version:2.2}")
    private String fabricVersion;

    // 初始化超时时间（秒）
    @Value("${constant.initHostTimeOut:300}")
    private long initHostTimeOut;

    // 节点资源检查
    @Value("${constant.memPerNode:#{1e6}}")
    private long memPerNode;

    @Value("${constant.hashAlgorithm:SHA2}")
    private String hashAlgorithm;
    @Autowired
    private FabricConfig fabricConfig;
    @Value("${fabric.mode.link}")
    private boolean linkFlag;
    private ConcurrentHashMap<String, Set<Integer>> usedPortsMap = new ConcurrentHashMap<>();

    @Override
    public void onApplicationEvent(GatewayInitedEvent event) {
        gateway = event.getGateway();
        log.info("DeployServiceImpl found new gateway: {}", gateway);
    }

    @Override
    public ChainInfoDTO getChainInfo() throws IOException {
        ChainInfoDTO chainInfoDTO = new ChainInfoDTO();
        InputStreamSource resource = null;
        // 获取通道
        if (linkFlag) {
            resource = new ClassPathResource(FabricConfig.networkConfigPath);
        } else {
            Path nodesRoot = Paths.get(constantProperties.getNodesRootDir());
            List<Path> chains = Files.list(nodesRoot).filter(Files::isDirectory).collect(Collectors.toList());
            File file = chains.get(0).resolve(FabricConfig.networkConfigPath).toFile();
            resource = new FileSystemResource(file);
        }
        String connConfStr = null;
        try {
            connConfStr = StreamUtils.copyToString(resource.getInputStream(), StandardCharsets.UTF_8);
        } catch (IOException e) {
            log.error("read conf failed: ", e);
            throw new AgentException(ConstantCode.SYSTEM_EXCEPTION);
        }
        JSONObject connConf = JSON.parseObject(connConfStr);
        chainInfoDTO.setChainName(connConf.getString("name"));

        JSONObject channels = connConf.getJSONObject("channels");
        // 遍历每一个通道，否则client.channels为空
        for (String channelName : channels.keySet()) {
            gateway.getNetwork(channelName);
        }
        String channelName = channels.keySet().stream().findFirst().get();

        Network network = gateway.getNetwork(channelName);
        Channel channel = network.getChannel();
        // 获取通道的配置块
        byte[] channelConfigBytes = new byte[0];
        try {
            channelConfigBytes = channel.getChannelConfigurationBytes();
        } catch (InvalidArgumentException e) {
            throw new AgentException(ConstantCode.SYSTEM_EXCEPTION);
        } catch (TransactionException e) {
            throw new AgentException(ConstantCode.SYSTEM_EXCEPTION);
        }
        // 从配置块中提取共识配置信息
        Configuration.ConsensusType consensusType = null;
        try {
            consensusType = Configuration.ConsensusType.parseFrom(Configtx.Config.parseFrom(channelConfigBytes)
                    .getChannelGroup()
                    .getGroupsMap()
                    .get("Orderer")
                    .getValuesMap()
                    .get("ConsensusType")
                    .getValue());
        } catch (InvalidProtocolBufferException e) {
            throw new AgentException(ConstantCode.SYSTEM_EXCEPTION);
        }
        log.info("consensusType: {}", consensusType.getType());
        chainInfoDTO.setConsensusType(consensusType.getType());

        // 获取链版本
        chainInfoDTO.setVersion(fabricVersion);

        Map<String, List<NodeInfo>> channelNodeMap = fillChannelNodeMap();
        chainInfoDTO.setChannelNodeMap(channelNodeMap);
        log.info("== chain info: {}", JSON.toJSONString(chainInfoDTO));
        return chainInfoDTO;
    }

    private Map<String, List<NodeInfo>> fillChannelNodeMap() {
        Map<String, List<NodeInfo>> channelNodeMap = new HashMap<>();

        HFClient hfClient = gateway.getClient();
        Map<String, Channel> channels = null;
        try {
            Field channelsField = HFClient.class.getDeclaredField("channels");
            channelsField.setAccessible(true);
            channels = (Map<String, Channel>) channelsField.get(hfClient);
            log.debug("== channelIds from gateway: {}", channels.keySet());
            //log.debug("== channels from gateway: {}", JSON.toJSONString(channels));
        } catch (Exception e) {
            log.error("== fillChannelNodeMap error:", e);
            throw new AgentException(ConstantCode.SYSTEM_EXCEPTION);
        }
        channels.forEach((channelId, channel) -> {
            log.debug("== channel [{}] from gateway: {}", channelId, JSON.toJSONString(channel));
            List<NodeInfo> nodeInfos = new ArrayList<>();
            channel.getPeers().stream().forEach(peer -> {
                log.debug("== peer properties: peerName:{},{}", peer.getName(), peer.getProperties());
                NodeInfo nodeInfo = new NodeInfo();
                String url = peer.getUrl();
                nodeInfo.setNodeIp(url.substring(url.lastIndexOf("/") + 1, url.lastIndexOf(":")));
                nodeInfo.setNodePort(Integer.valueOf(url.substring(url.lastIndexOf(":") + 1)));
                nodeInfo.setNodeFullName(peer.getName());
                nodeInfos.add(nodeInfo);
            });
            channel.getOrderers().stream().forEach(orderer -> {
                NodeInfo nodeInfo = new NodeInfo();
                String url = orderer.getUrl();
                nodeInfo.setNodeIp(url.substring(url.lastIndexOf("/") + 1, url.lastIndexOf(":")));
                nodeInfo.setNodePort(Integer.valueOf(url.substring(url.lastIndexOf(":") + 1)));
                nodeInfo.setNodeFullName(orderer.getName());
                nodeInfos.add(nodeInfo);
            });
            channelNodeMap.put(channelId, nodeInfos);
        });
        return channelNodeMap;
    }

    @Override
    public ChainInitResp initHost(ChainInitReq dto) {
        // limit: single fabric to agent
        limitSingleChain();
        // 校验
        Assert.isTrue(!CollectionUtils.isEmpty(dto.getPeers()), "peers not allow empty");
        Assert.isTrue(!CollectionUtils.isEmpty(dto.getOrderers()), "orderers not allow empty");

        // 检查资源
        checkResource(dto);

        // 并行检查
        List<CompletableFuture<ChainInitResp.Result>> fs = Stream.concat(
                        dto.getPeers().stream()
                                .map(n -> new ChainInitResp.Result(n.getIp(), FabricNodeType.PEER, n.getNumber(), Boolean.TRUE, null)),
                        dto.getOrderers().stream()
                                .map(n -> new ChainInitResp.Result(n.getIp(), FabricNodeType.ORDERER, n.getNumber(), Boolean.TRUE, null))
                )
                .map(r ->
                        // 插入线程池,并行检查
                        CompletableFuture.supplyAsync(() -> {
                            try {
                                // ping
                                ansibleService.execPing(r.getIp());

                                // 检查该机构是否已经部署
                                boolean containerExists = ansibleService.checkContainerExists(r.getIp(), dto.getAgency() + "$");
                                Assert.isTrue(!containerExists, () -> String.format("%s 容器已存在", dto.getAgency()));
                            } catch (Exception e) {
                                log.error("init host:{} fail", r.getIp(), e);
                                r.setResult(Boolean.FALSE).setErrorReason(e.getMessage());
                            }
                            return r;
                        }, threadPoolTaskScheduler)
                )
                .collect(Collectors.toList());

        // 等待所有结果执行完
        try {
            CompletableFuture.allOf(fs.toArray(new CompletableFuture[0])).get(initHostTimeOut, TimeUnit.SECONDS);
        } catch (Exception e) {
            log.error("init hosts fail", e);
            throw new RuntimeException("initHostAndDocker fail");
        }

        // 返回结果
        return new ChainInitResp(
                fs.stream()
                        .map(c -> {
                            try {
                                return c.get();
                            } catch (Exception e) {
                                return new ChainInitResp.Result();
                            }
                        })
                        .collect(Collectors.toList())
        );
    }

    /**
     * 检查资源
     */
    private void checkResource(ChainInitReq dto) {
        Map<String, Integer> hosts = Stream.concat(
                        dto.getPeers().stream(),
                        dto.getOrderers().stream()
                )
                .collect(
                        Collectors.toMap(
                                ChainInitReq.Node::getIp,
                                n -> Optional.ofNullable(n.getNumber()).orElse(0),
                                (a, b) -> a + b
                        )
                );

        // 检查节点资源
        List<String> errors = hosts.entrySet().stream()
                .flatMap(e -> {
                    long needMem = e.getValue() * memPerNode;
                    int needCpu = e.getValue() >= 3 ? 2 : 1;

                    NodeInfoDTO allInfo = monitorService.hostInfo(e.getKey());

                    List<String> hostErrors = new LinkedList<>();
                    if (allInfo.getMaxMem() - allInfo.getUsedMem() < needMem) {
                        hostErrors.add(String.format("%s 需要内存 %s, 实际内存 %s", e.getKey(), needMem, allInfo.getMaxMem() - allInfo.getUsedMem()));
                    }
                    if (allInfo.getCpuNum() < needCpu) {
                        hostErrors.add(String.format("%s 需要cpu数量 %s, 实际数量 %s", e.getKey(), needCpu, allInfo.getCpuNum()));
                    }
                    return hostErrors.stream();
                })
                .collect(Collectors.toList());
        Assert.isTrue(errors.size() == 0, () -> String.join(",", errors));
    }

    @Override
    public IpAndPortsDTO getAvailablePort(String ip, int startPort, int number) {
        int start = startPort;
        IpAndPortsDTO ipAndPortsDTO = new IpAndPortsDTO();
        ipAndPortsDTO.setIp(ip);
        ipAndPortsDTO.setPorts(new LinkedList<>());

        for (int i = 0; i < number; i++) {
            int port = availablePort(ip, start);
            ipAndPortsDTO.getPorts().add(port);
            start = port + 1;
        }
        return ipAndPortsDTO;
    }

    private int availablePort(String ip, int startPort) {
        int port = startPort;
        while (port < 65536) {
            Map<Integer, Boolean> portUseMap = ansibleService.getPortUseMap(ip, port);
            if (portUseMap.values().contains(Boolean.TRUE)) {
                port++;
            } else {
                break;
            }
        }

        Assert.isTrue(port < 65536, () -> String.format("ip:%s not find available port", ip));
        return port;
    }

    private void pullFabricImage(String ip, String version) {
        String imageFullName = String.format("%s:%s", DOCKER_IMAGE, version);
        boolean isExist = ansibleService.checkImageExists(ip, imageFullName);
        if (isExist) {
            log.info("pullImage jump over for image:{} already exist.", imageFullName);
            return;
        }

        // 下载docker镜像
        log.info("docker pull {}", imageFullName);
        String dockerPullCommand = String.format("docker pull %s", imageFullName);
        ExecuteResult result = ansibleService.execDocker(ip, dockerPullCommand);
        if (result.failed()) {
            throw new RuntimeException(ConstantCode.ANSIBLE_PULL_DOCKER_HUB_ERROR.getMessage());
        }
    }

    @Override
    public DeployResultDTO getDeployResult(String chainUid) {
        return getDeployResultFromDB(chainUid, chainUid);
    }

    @Override
    public DeployResultDTO getAddNodeResult(String chainUid, String requestId) {
        return getDeployResultFromDB(chainUid, requestId);
    }

    private DeployResultDTO getDeployResultFromDB(String chainUid, String requestId) {
        DeployResultDo chain = deployResultRepo.getOne(
                new LambdaQueryWrapper<DeployResultDo>()
                        .eq(DeployResultDo::getChainUid, chainUid)
                        .eq(DeployResultDo::getRequestId, requestId)
        );
        Assert.notNull(chain, "not have this chain, chainUid:" + chainUid);
        List<DeployResultNodeDo> nodes = deployResultNodeRepo.list(
                new LambdaQueryWrapper<DeployResultNodeDo>()
                        .eq(DeployResultNodeDo::getResultId, chain.getId())
        );

        // 格式转换
        DeployResultDTO deployResultDTO = new DeployResultDTO();
        BeanUtils.copyProperties(chain, deployResultDTO);
        List<NodeDetail> collect = nodes.stream()
                .map(n -> {
                    NodeDetail nodeDetail = new NodeDetail();
                    BeanUtils.copyProperties(n, nodeDetail);
                    nodeDetail.setPeer(FabricNodeType.PEER.getIndex().equals(n.getNodeType())? true: false);
                    return nodeDetail;
                })
                .collect(Collectors.toList());
        deployResultDTO.setHostNodes(collect);
        return deployResultDTO;
    }

    @Override
    public void addAppChannel(AddAppChannelReq req) {
        // check param
        String chainUid = req.getChainUid();
        Path chainRoot = pathService.getChainRoot(chainUid);
        if (!chainRoot.toFile().exists()) {
            log.error("check your chainUid, invalid chain path: {}", chainRoot.toFile().getAbsolutePath());
            throw new AgentException(ConstantCode.CHAIN_NOT_EXIST);
        }
        // 3.1.创建应用通道tx交易文件
        // -c channel2 -p ./NODES_ROOT/1ef243fd87
        String createTxCommand = String.format("bash %s -c %s -p %s",
                // create_app_channel_tx.sh shell script
                constantProperties.getCreateAppChannelTxShell(),
                // app channel id
                req.getChannelId(),
                // chain nodes root path
                chainRoot.toFile().getAbsolutePath()
        );
        exec(createTxCommand);
        // 3.2.创建应用通道block区块文件
        // -c channel2 -p ./NODES_ROOT/1ef243fd87 -o 192.168.3.128:7050 -j orderer0.org1.example.com -r org1.example.com -m Org1MSP -n 192.168.3.128:7051 -s peer0.org1.example.com -g org1.example.com
        String peerName = req.getPeerName();
        String peerOrgName = req.getPeerOrgName();
        String peerEndpoint = req.getPeerEndpoint();
        String ordererName = req.getOrdererName();
        String ordererOrgName = req.getOrdererOrgName();
        String ordererEndpoint = req.getOrdererEndpoint();
        String peerMspId = genMspId(peerOrgName, true);
        String createBlockCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // create_app_channel_block.sh shell script
                constantProperties.getCreateAppChannelBlockShell(),
                // app channel id
                req.getChannelId(),
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(createBlockCommand);
    }

    @Override
    public void addPeersIntoAppChannel(AddPeersIntoAppChannelReq req) throws Exception {
        String chainUid = req.getChainUid();
        String channelId = req.getChannelId();
        String peerName = req.getPeerName();
        String peerOrgName = req.getPeerOrgName();
        String peerEndpoint = req.getPeerEndpoint();
        String peerMspId = genMspId(peerOrgName, true);
        String ordererName = req.getOrdererName();
        String ordererOrgName = req.getOrdererOrgName();
        String ordererEndpoint = req.getOrdererEndpoint();

        List<NodeTriple> peersNotInChannel = req.getPeersNotInChannel();
        List<NodeTriple> peersInChannel = req.getPeersInChannel();
        List<NodeTriple> allPeers = new ArrayList<>();
        allPeers.addAll(peersInChannel);
        allPeers.addAll(peersNotInChannel);

        // 3.3.peer节点加入通道
        for (NodeTriple peerTriple : peersNotInChannel) {
            String peerName0 = peerTriple.getNodeName();
            String peerOrgName0 = peerTriple.getNodeOrgName();
            String peerEndpoint0 = peerTriple.getNodeEndpoint();
            String peerMspId0 = genMspId(peerOrgName0, true);
            String joinCommand = String.format("bash %s -c %s -p %s -m %s -n %s -s %s -g %s",
                    // join_app_channel.sh shell script
                    constantProperties.getJoinAppChannelShell(),
                    // app channel id
                    channelId,
                    // chain nodes root path
                    pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                    peerMspId0,
                    peerEndpoint0,
                    peerName0,
                    peerOrgName0
            );
            exec(joinCommand);
        }
        // 3.4.获取应用通道最近的配置块
        String fetchCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // fetch_app_channel_config.sh shell script
                constantProperties.getFetchAppChannelConfigShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(fetchCommand);
        // 3.5.配置应用通道
        String anchorPeersStr = getAnchorPeersStr(allPeers);
        String configCommand = String.format("bash %s -c %s -p %s -m %s -x %s",
                // config_app_channel.sh shell script
                constantProperties.getConfigAppChannelShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                peerMspId,
                anchorPeersStr
        );
        exec(configCommand);
        // 3.6.更改应用通道
        String updateCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // update_app_channel.sh shell script
                constantProperties.getUpdateAppChannelShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(updateCommand);
        // 重新持久化配置并重新初始化gateway
        try {
            initGateway(req);
        } catch (Exception e) {
            log.error("addPeersIntoAppChannel.initGateway err: {}", e.getMessage());
            throw e;
        }
    }

    private void initGateway(AddPeersIntoAppChannelReq req) throws Exception {
        String chainUid = req.getChainUid();
        String channelId = req.getChannelId();
        List<NodeTriple> peersNotInChannel = req.getPeersNotInChannel();
        Path chainRoot = pathService.getChainRoot(chainUid);
        // rebuild connection-tls.json
        Path networkConfigPath = chainRoot.resolve(FabricConfig.networkConfigPath);
        File networkConfigFile = networkConfigPath.toFile();
        log.info("== Read networkConfig from: {}", networkConfigFile.getAbsolutePath());
        String networkConfigJsonStr = new String(Files.readAllBytes(networkConfigPath));
        log.info("== Read networkConfig content: {}", networkConfigJsonStr);
        NetworkConfig networkConfig = JSON.parseObject(networkConfigJsonStr, NetworkConfig.class);

        Map<String, Object> channels = networkConfig.getChannels();
        Map<String, Object> organizations = networkConfig.getOrganizations();
        Map<String, Object> peers = networkConfig.getPeers();
        if (channels.containsKey(channelId)) {
            JSONObject channelObj = JSON.parseObject(JSON.toJSONString(channels.get(channelId)));
            for (NodeTriple nodeTriple : peersNotInChannel) {
                JSONObject peerX = new JSONObject();
                String nodeName = nodeTriple.getNodeName();
                String peerIp = nodeTriple.getNodeEndpoint().split(":")[0];
                String peerPort = nodeTriple.getNodeEndpoint().split(":")[1];
                peerX.put("endorsingPeer", channelObj.getJSONObject("peers").isEmpty());
                peerX.put("chaincodeQuery", true);
                peerX.put("ledgerQuery", true);
                peerX.put("eventSource", true);
                channelObj.getJSONObject("peers").put(nodeName, peerX);

                String orgName = nodeTriple.getNodeOrgName();
                // $.organizations
                if (organizations.get(orgName) == null) {
                    JSONObject orgObj = new JSONObject();
                    orgObj.put("mspid", genMspId(orgName, true));
                    orgObj.put("peers", Lists.newArrayList(nodeName));
                    orgObj.put("certificateAuthorities", Lists.newArrayList(buildCaName(orgName)));
                    JSONObject adminPrivateKeyPEM = new JSONObject();
                    Path mspPath = chainRoot.resolve("organizations").resolve("peerOrganizations").resolve(orgName).resolve("users")
                            .resolve("Admin@" + orgName).resolve("msp");
                    Path keystorePath = mspPath.resolve("keystore");
                    // 过滤出普通文件
                    List<Path> files = Files.list(keystorePath).filter(Files::isRegularFile).collect(Collectors.toList());
                    String keystore = files.get(0).toFile().getAbsolutePath();
                    adminPrivateKeyPEM.put("path", keystore);
                    orgObj.put("adminPrivateKeyPEM", adminPrivateKeyPEM);
                    JSONObject signedCertPEM = new JSONObject();
                    signedCertPEM.put("path", mspPath.resolve("signcerts").resolve("cert.pem").toFile().getAbsolutePath());
                    orgObj.put("signedCertPEM", signedCertPEM);
                    organizations.put(orgName, orgObj);
                } else {
                    JSONObject orgObj = JSON.parseObject(JSON.toJSONString(organizations.get(orgName)));
                    if (!orgObj.getJSONArray("peers").contains(nodeName)) {
                        orgObj.getJSONArray("peers").add(nodeName);
                    }
                    organizations.put(orgName, orgObj);
                }

                // $.peers
                if (!peers.containsKey(nodeName)) {
                    JSONObject peer = new JSONObject();
                    peer.put("url", "grpcs://" + peerIp + ":" + peerPort);
                    JSONObject grpcOptions = new JSONObject();
                    grpcOptions.put("request-timeout", 120001);
                    peer.put("grpcOptions", grpcOptions);
                    JSONObject tlsCACerts = new JSONObject();
                    String tlsCa = chainRoot.resolve(String.format("organizations/peerOrganizations/%s/peers/%s/tls/ca.crt",
                            orgName, nodeName)).toFile().getAbsolutePath();
                    tlsCACerts.put("path", tlsCa);
                    peer.put("tlsCACerts", tlsCACerts);
                    peers.put(nodeName, peer);
                }
            }
            channels.put(channelId, channelObj);
        } else {
            JSONObject channelObj = new JSONObject();
            JSONArray orderers = JSON.parseObject(JSON.toJSONString(channels.get(channels.keySet().stream().findFirst().get())))
                    .getJSONArray("orderers");
            channelObj.put("orderers", orderers);
            channelObj.put("peers", new JSONObject());
            for (NodeTriple nodeTriple : peersNotInChannel) {
                JSONObject peerX = new JSONObject();
                String nodeName = nodeTriple.getNodeName();
                String peerIp = nodeTriple.getNodeEndpoint().split(":")[0];
                String peerPort = nodeTriple.getNodeEndpoint().split(":")[1];
                peerX.put("endorsingPeer", channelObj.getJSONObject("peers").isEmpty());
                peerX.put("chaincodeQuery", true);
                peerX.put("ledgerQuery", true);
                peerX.put("eventSource", true);
                channelObj.getJSONObject("peers").put(nodeName, peerX);

                String orgName = nodeTriple.getNodeOrgName();
                // $.organizations
                if (organizations.get(orgName) == null) {
                    JSONObject orgObj = new JSONObject();
                    orgObj.put("mspid", genMspId(orgName, true));
                    orgObj.put("peers", Lists.newArrayList(nodeName));
                    orgObj.put("certificateAuthorities", Lists.newArrayList(buildCaName(orgName)));
                    JSONObject adminPrivateKeyPEM = new JSONObject();
                    Path mspPath = chainRoot.resolve("organizations").resolve("peerOrganizations").resolve(orgName).resolve("users")
                            .resolve("Admin@" + orgName).resolve("msp");
                    Path keystorePath = mspPath.resolve("keystore");
                    // 过滤出普通文件
                    List<Path> files = Files.list(keystorePath).filter(Files::isRegularFile).collect(Collectors.toList());
                    String keystore = files.get(0).toFile().getAbsolutePath();
                    adminPrivateKeyPEM.put("path", keystore);
                    orgObj.put("adminPrivateKeyPEM", adminPrivateKeyPEM);
                    JSONObject signedCertPEM = new JSONObject();
                    signedCertPEM.put("path", mspPath.resolve("signcerts").resolve("cert.pem").toFile().getAbsolutePath());
                    orgObj.put("signedCertPEM", signedCertPEM);
                    organizations.put(orgName, orgObj);
                } else {
                    JSONObject orgObj = JSON.parseObject(JSON.toJSONString(organizations.get(orgName)));
                    if (!orgObj.getJSONArray("peers").contains(nodeName)) {
                        orgObj.getJSONArray("peers").add(nodeName);
                    }
                    organizations.put(orgName, orgObj);
                }

                // $.peers
                if (!peers.containsKey(nodeName)) {
                    JSONObject peer = new JSONObject();
                    peer.put("url", "grpcs://" + peerIp + ":" + peerPort);
                    JSONObject grpcOptions = new JSONObject();
                    grpcOptions.put("request-timeout", 120001);
                    peer.put("grpcOptions", grpcOptions);
                    JSONObject tlsCACerts = new JSONObject();
                    String tlsCa = chainRoot.resolve(String.format("organizations/peerOrganizations/%s/peers/%s/tls/ca.crt",
                            orgName, nodeName)).toFile().getAbsolutePath();
                    tlsCACerts.put("path", tlsCa);
                    peer.put("tlsCACerts", tlsCACerts);
                    peers.put(nodeName, peer);
                }
            }
            channels.put(channelId, channelObj);
        }

        String newNetworkConfigJsonStr = JSON.toJSONString(networkConfig);
        log.info("== Auto gen networkConfig: {}", newNetworkConfigJsonStr);
        log.info("Ready to re-persist to: {}", networkConfigFile.getAbsolutePath());
        FileUtils.write(networkConfigFile, newNetworkConfigJsonStr);

        // init new gateway
        fabricConfig.reinitializeFabricGatewayBean();
    }

    @Override
    public void nodeOperation(NodeOperationReq req) {
        String dockerComposeCmd = "";
        String nodeName = req.getNodeName();

        switch (req.getOperation()) {
            case START:
                dockerComposeCmd = String.format("docker start %s", nodeName);
                break;
            case STOP:
                dockerComposeCmd = String.format("docker stop %s", nodeName);
                break;
            case RESTART:
                dockerComposeCmd = String.format("docker restart %s", nodeName);
                break;
            default:
        }

        ansibleService.exec(req.getHostIp(), dockerComposeCmd);
    }

    @Async
    @Override
    public void deploy(DeployReq req) {
        log.info("deploy req:{}", req);

        String chainUid = req.getChainUid();
        // 链部署记录
        DeployResultDo deployResultDo = new DeployResultDo().setChainUid(chainUid).setRequestId(chainUid);
        deployResultRepo.save(deployResultDo);

        try {
            // part A
            if (CollectionUtils.isEmpty(req.getHostNodesInOrgMap())) {
                throw new RuntimeException("节点信息为空");
            }

            trackDeploymentStage(req.getChainUid(), FabricDeployStage.INIT_CA);
            initCaServerConf(req);
            initCaServerContainer(req);

            // 自动检测和分配节点可用端口
            HashMap<String, Set<Integer>> usedPortsByIp = checkAvailableNodePort(req);
            Path organizationsPath = pathService.getOrganizationsRoot(chainUid);


            // 1. init client env
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.INIT_CLIENT_ENV);
            intClientEnv(req);

            // 2. init org config
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.INIT_ORG);
            initOrg(req);

            // 3.init peer config
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.INIT_PEER);
            initPeer(req);
            // 4.init orderer config
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.INIT_ORDERER);
            initOrderer(req);

            // 5. generate config tx
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.GEN_CONF_TX);
            genConfigTx(req);

            // 6. init genesis config
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.GEN_GENESIS);
            initGenesisConf(req);

            // part B
            // 1.发送各个节点主机对应的链目录
            Set<String> hostIps = new HashSet<>();
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.SCP_CONF);
            transferFiles2NodeHosts(req, organizationsPath, hostIps);

            // 2.准备和发送docker-compose.yaml到各节点主机，并远程启动各个节点: docker-compose -f docker-compose.yaml up -d
            List<NodeTriple> orderers = new ArrayList<>();
            List<NodeTriple> peers = new ArrayList<>();
            trackDeploymentStage(req.getChainUid(), FabricDeployStage.LAUNCH_NODE);
            List<NodeDetail> hostNodes = readyDockerComposeAndLaunchNodes(req, hostIps, orderers, peers, usedPortsByIp);
            // 循环等待最大超时时间，检测所有节点启动成功与否，并返回一个可用的peer和orderer
            Pair<Integer, Integer> usefulNodeIndexes = waitForAllNodesStarted(hostNodes, orderers, peers);
            if (usefulNodeIndexes.getLeft() == -1 || usefulNodeIndexes.getRight() == -1) {
                throw new AgentException(ConstantCode.START_NODES_ERROR);
            }

            // 3.创建应用通道并将所有peer加入
            trackDeploymentStage(chainUid, FabricDeployStage.CREATE_CHANNEL);
            createAndJoinChannel(req, orderers, peers, usefulNodeIndexes);

            // 保存结果
            deployResultDo.setChannelId(req.getChannelId());
            deployResultDo.setStage(FabricDeployStage.DEPLOY_SUCCESS.getIndex());
            deployResultRepo.updateById(deployResultDo);
            List<DeployResultNodeDo> resultNodes = hostNodes.stream()
                    .map(h -> {
                                DeployResultNodeDo deployResultNodeDo = new DeployResultNodeDo();
                                BeanUtils.copyProperties(h, deployResultNodeDo);
                                deployResultNodeDo.setNodeType(h.isPeer() ? FabricNodeType.PEER.getIndex() : FabricNodeType.ORDERER.getIndex());
                                deployResultNodeDo.setResultId(deployResultDo.getId());
                                return deployResultNodeDo;
                            }
                    )
                    .collect(Collectors.toList());
            deployResultNodeRepo.saveBatch(resultNodes);

            // auto init gateway
            initGateway(req, hostNodes);
        } catch (Exception e) {
            log.error("deploy chain fail", e);
            deployResultDo.setStage(FabricDeployStage.DEPLOY_FAILED.getIndex());
            deployResultDo.setError(e.getMessage());
            deployResultRepo.updateById(deployResultDo);
        }
    }

    /**
     * 可以被纳管多次，因为agent不需要额外处理数据；
     * 纳管模式下，不允许再被创建链；非纳管模式下，不允许创建多个链。
     */
    private void limitSingleChain() {
        if (linkFlag) {
            throw new RuntimeException("该agent已纳管了一条fabric链，不允许再创建另外的fabric链");
        } else {
            // 非纳管的默认agent对单链，首先定位到链目录
            Path nodesRoot = Paths.get(constantProperties.getNodesRootDir());
            if (!nodesRoot.toFile().exists()) {
                log.info("= not found any chains in nodes_root");
                return;
            }
            List<Path> chains = null;
            try {
                chains = Files.list(nodesRoot).filter(Files::isDirectory).collect(Collectors.toList());
            } catch (IOException e) {
                throw new AgentException(ConstantCode.IO_EXCEPTION);
            }
            if (chains.isEmpty()) {
                log.info("= not found any chains in nodes_root");
            } else if (chains.size() >= 1) {
                log.error("= found more than one chain in nodes_root");
                throw new RuntimeException("该agent已创建了一条fabric链，不允许再创建另外的fabric链");
            }
        }
    }

    private void initGateway(DeployReq req, List<NodeDetail> hostNodes) throws Exception {
        String chainUid = req.getChainUid();
        // build connection-tls.json
        NetworkConfig networkConfig = new NetworkConfig();
        networkConfig.setName(chainUid);
        networkConfig.setVersion("1.0.0");
        Set<String> orgNames = req.getHostNodesInOrgMap().keySet();
        Map<String, Object> clientMap = new HashMap<>();
        clientMap.put("organization", orgNames.stream().findFirst().get()); // 只支持单机构连接
        JSONObject connectionPeer = new JSONObject();
        connectionPeer.put("endorser", "300");
        JSONObject connection = new JSONObject();
        JSONObject timeout = new JSONObject();
        timeout.put("peer", connectionPeer);
        timeout.put("orderer", "300");
        connection.put("timeout", timeout);
        clientMap.put("connection", connection);
        networkConfig.setClient(clientMap);

        Map<String, Object> channels = new HashMap<>();
        JSONObject channel = new JSONObject();
        List<String> ordererNames = hostNodes.stream().filter(nodeDetail -> !nodeDetail.isPeer())
                .map(NodeDetail::getNodeFullName).collect(Collectors.toList());
        channel.put("orderers", ordererNames);
        List<String> peerNames = hostNodes.stream().filter(NodeDetail::isPeer)
                .map(NodeDetail::getNodeFullName).collect(Collectors.toList());
        JSONObject channelPeers = new JSONObject();
        for (int i = 0; i < peerNames.size(); i++) {
            String peerName = peerNames.get(i);
            JSONObject peerX = new JSONObject();
            peerX.put("endorsingPeer", i == 0);
            peerX.put("chaincodeQuery", true);
            peerX.put("ledgerQuery", true);
            peerX.put("eventSource", true);
            channelPeers.put(peerName, peerX);
        }
        channel.put("peers", channelPeers);
        channels.put(req.getChannelId(), channel);
        networkConfig.setChannels(channels);

        Path chainRoot = pathService.getChainRoot(chainUid);
        Map<String, Object> organizations = new HashMap<>();
        for (String orgName : orgNames) {
            JSONObject org = new JSONObject();
            String mspId = genMspId(orgName, true);
            org.put("mspid", mspId);
            List<String> peerNameList = hostNodes.stream().filter(node -> node.isPeer() && orgName.equals(node.getOrgFullName()))
                    .map(NodeDetail::getNodeFullName).collect(Collectors.toList());
            if (CollectionUtils.isEmpty(peerNameList)) {
                continue;
            }
            org.put("peers", peerNameList);
            JSONArray certificateAuthorities = new JSONArray();
            certificateAuthorities.add(buildCaName(orgName));
            org.put("certificateAuthorities", certificateAuthorities);
            JSONObject adminPrivateKeyPEM = new JSONObject();
            Path mspPath = chainRoot.resolve("organizations").resolve("peerOrganizations").resolve(orgName).resolve("users")
                    .resolve("Admin@" + orgName).resolve("msp");
            Path keystorePath = mspPath.resolve("keystore");
            // 过滤出普通文件
            List<Path> files = Files.list(keystorePath).filter(Files::isRegularFile).collect(Collectors.toList());
            String keystore = files.get(0).toFile().getAbsolutePath();
            adminPrivateKeyPEM.put("path", keystore);
            org.put("adminPrivateKeyPEM", adminPrivateKeyPEM);
            JSONObject signedCertPEM = new JSONObject();
            signedCertPEM.put("path", mspPath.resolve("signcerts").resolve("cert.pem").toFile().getAbsolutePath());
            org.put("signedCertPEM", signedCertPEM);
            organizations.put(orgName, org);
        }
        networkConfig.setOrganizations(organizations);

        Map<String, Object> orderers = new HashMap<>();
        Map<String, Object> peers = new HashMap<>();
        for (NodeDetail nodeDetail : hostNodes) {
            String nodeFullName = nodeDetail.getNodeFullName();
            String orgFullName = nodeDetail.getOrgFullName();
            if (!nodeDetail.isPeer()) {
                JSONObject orderer = new JSONObject();
                orderer.put("url", "grpcs://" + nodeDetail.getNodeIp() + ":" + nodeDetail.getNodePort());
                orderer.put("mspid", nodeDetail.getMspId());
                orderer.put("grpcOptions", new JSONObject());
                JSONObject tlsCACerts = new JSONObject();
                String tlsCa = chainRoot.resolve(String.format("organizations/ordererOrganizations/%s/orderers/%s/tls/ca.crt",
                        orgFullName, nodeFullName)).toFile().getAbsolutePath();
                tlsCACerts.put("path", tlsCa);
                orderer.put("tlsCACerts", tlsCACerts);
                JSONObject adminPrivateKeyPEM = new JSONObject();
                Path mspPath = chainRoot.resolve("organizations").resolve("ordererOrganizations").resolve(orgFullName).resolve("users")
                        .resolve("Admin@" + orgFullName).resolve("msp");
                Path keystorePath = mspPath.resolve("keystore");
                // 过滤出普通文件
                List<Path> files = Files.list(keystorePath).filter(Files::isRegularFile).collect(Collectors.toList());
                String keystore = files.get(0).toFile().getAbsolutePath();
                adminPrivateKeyPEM.put("path", keystore);
                orderer.put("adminPrivateKeyPEM", adminPrivateKeyPEM);
                JSONObject signedCertPEM = new JSONObject();
                signedCertPEM.put("path", mspPath.resolve("signcerts").resolve("cert.pem").toFile().getAbsolutePath());
                orderer.put("signedCertPEM", signedCertPEM);
                orderers.put(nodeFullName, orderer);
            } else {
                JSONObject peer = new JSONObject();
                peer.put("url", "grpcs://" + nodeDetail.getNodeIp() + ":" + nodeDetail.getNodePort());
                JSONObject grpcOptions = new JSONObject();
                grpcOptions.put("request-timeout", 120001);
                peer.put("grpcOptions", grpcOptions);
                JSONObject tlsCACerts = new JSONObject();
                String tlsCa = chainRoot.resolve(String.format("organizations/peerOrganizations/%s/peers/%s/tls/ca.crt",
                        orgFullName, nodeFullName)).toFile().getAbsolutePath();
                tlsCACerts.put("path", tlsCa);
                peer.put("tlsCACerts", tlsCACerts);
                peers.put(nodeFullName, peer);
            }
        }
        networkConfig.setOrderers(orderers);
        networkConfig.setPeers(peers);

        Map<String, Object> certificateAuthorities = new HashMap<>();
        Map<String, Object> ca = new HashMap<>();
        ca.put("url", "https://" + constantProperties.getCaHost() + ":" + getCaServerPort(chainUid, orgNames.stream().findFirst().get()));
        JSONObject grpcOptions = new JSONObject();
        grpcOptions.put("verify", true);
        ca.put("grpcOptions", grpcOptions);
        JSONObject tlsCACerts = new JSONObject();
        tlsCACerts.put("path", chainRoot.resolve("ca-cert.pem").toFile().getAbsolutePath());
        ca.put("tlsCACerts", tlsCACerts);
        JSONArray registrar = new JSONArray();
        JSONObject registrarObj = new JSONObject();
        registrarObj.put("enrollId", constantProperties.getCaAdmin());
        registrarObj.put("enrollSecret", constantProperties.getCaPassword());
        registrar.add(registrarObj);
        ca.put("registrar", registrar);
        certificateAuthorities.put(buildCaName(orgNames.stream().findFirst().get()), ca);
        networkConfig.setCertificateAuthorities(certificateAuthorities);

        String networkConfigJsonStr = JSON.toJSONString(networkConfig);
        log.info("== Auto gen networkConfig: {}", networkConfigJsonStr);
        File networkConfigFile = chainRoot.resolve(FabricConfig.networkConfigPath).toFile();
        log.info("Ready to persist to: {}", networkConfigFile.getAbsolutePath());
        FileUtils.write(networkConfigFile, networkConfigJsonStr);

        // init new gateway
        fabricConfig.reinitializeFabricGatewayBean();
    }

    // pem去掉中间的换行符
    private String processPem(String pem) {
        if (pem == null) {
            return pem;
        }

        String[] split = pem.trim().split("\n");
        if (split.length < 2) {
            return pem;
        }

        String head = split[0];
        String foot = split[split.length - 1];
        List<String> content = Arrays.asList(Arrays.copyOfRange(split, 1, split.length - 1));

        return String.format("%s\n%s\n%s", head, String.join("", content), foot);
    }

    @Override
    public void delete(DeleteReq req) {
        String chainUid = req.getChainUid();
        Path chainDockerDir = pathService.getChainDockerDir(chainUid);
        String[] list = chainDockerDir.toFile().list();
        if (Objects.nonNull(list)) {
            for (String hostIp : Objects.requireNonNull(chainDockerDir.toFile().list())) {
                String dockerComposeFilePath = String.format("%s/%s/docker-compose.yaml", constantProperties.getInstallDir(), chainUid);
                String chainDir = String.format("%s/%s ", constantProperties.getInstallDir(), chainUid);
                String deleteDir = String.format("%s/%s", constantProperties.getInstallDir(), "deleted-tmp");
                String chainDeleteDir = String.format("%s/%s", deleteDir, chainUid);
                String deleteChainCmd = String.format("mv %s %s ", chainDir, chainDeleteDir);

                String dockerComposeCmd = String.format("docker-compose -f %s down --volumes --remove-orphans ", dockerComposeFilePath);
                // stop nodes with docker-compose
                log.info("stopping nodes on host: {}", hostIp);

                ansibleService.exec(hostIp, dockerComposeCmd, false);
                ansibleService.exec(hostIp, "mkdir -p " + deleteDir, false);
                // fix repeat chainUid deleted
                ansibleService.exec(hostIp, "sudo rm -rf " + chainDeleteDir, false);
                ansibleService.exec(hostIp, deleteChainCmd, false);
            }
        }

        killCaServer(chainUid);

        try {
            pathService.deleteChain(chainUid);
            chainInfoRepo.removeByChainUid(chainUid);

            // 删除数据
            List<DeployResultDo> rs = deployResultRepo.list(
                    new LambdaQueryWrapper<DeployResultDo>()
                            .eq(DeployResultDo::getChainUid, chainUid)
            );
            for (DeployResultDo r : rs) {
                deployResultRepo.remove(new LambdaQueryWrapper<DeployResultDo>().eq(DeployResultDo::getId, r.getId()));
                deployResultNodeRepo.remove(new LambdaQueryWrapper<DeployResultNodeDo>().eq(DeployResultNodeDo::getResultId, r.getId()));
            }
        } catch (IOException e) {
            log.error("Delete chain config files error:[]", e);
            log.error("Please delete/move chain config files manually");
            throw new AgentException(ConstantCode.DELETE_NODE_DIR_ERROR);
        }
    }

    private void killCaServer(String chainUid) {
//        Integer caServerPort;
//        try {
//            caServerPort = getCaServerPort(chainUid);
//        } catch (Exception e) {
//            return;
//        }
//        if (caServerPort != null && caServerPort > 0) {
//            JavaCommandExecutor.executeCommand("kill -9 $(lsof -t -i:" + caServerPort + ")", constantProperties.getExecShellTimeout());
//        }
        try {
            killCaServerContainer(chainUid);
        } catch (Exception e) {
            log.error("fail to kill ca container.");
        }
    }

    private void killCaServerContainer(String chainUid) {
        Path caServerDir = pathService.getCaServerDir(chainUid);
        File[] orgList = caServerDir.toFile().listFiles();
        if (Objects.isNull(orgList)) {
            return;
        }
        for (File orgDomainDir : orgList) {
            if (!orgDomainDir.isDirectory()) {
                continue;
            }
            if (orgDomainDir.list() == null) {
                return;
            }
            long dockerFileNum = Arrays.stream(Objects.requireNonNull(orgDomainDir.list())).filter(f -> f.endsWith("docker-compose.yaml")).count();
            if (dockerFileNum == 0) {
                return;
            }
            String dockerComposeFilePath = String.format("%s/docker-compose.yaml", orgDomainDir.getAbsolutePath());
            String dockerComposeCmd = String.format("docker-compose -f %s down --volumes --remove-orphans ", dockerComposeFilePath);
            exec(dockerComposeCmd);
        }

    }

    @Async
    @Override
    public void addNode(BaseDeployReq req) {
        try {
            assignNodePortBeforeAddingNode(req);

            trackDeployNodeStage(req.getChainUid(), req.getRequestId(), FabricDeployNodeStage.INIT_PEER);
            initPeer(req);
            trackDeployNodeStage(req.getChainUid(), req.getRequestId(), FabricDeployNodeStage.INIT_ORDERER);
            initOrderer(req);

            Path organizationsPath = pathService.getOrganizationsRoot(req.getChainUid());
            Set<String> hostIps = new HashSet<>();
            trackDeployNodeStage(req.getChainUid(), req.getRequestId(), FabricDeployNodeStage.SCP_CONF);
            transferFiles2NodeHosts(req, organizationsPath, hostIps);
            // 准备和发送docker-compose.yaml到各节点主机，并远程启动各个节点: docker-compose -f docker-compose.yaml up -d
            List<NodeTriple> orderers = new ArrayList<>();
            List<NodeTriple> peers = new ArrayList<>();
            trackDeployNodeStage(req.getChainUid(), req.getRequestId(), FabricDeployNodeStage.LAUNCH_NODE);
            List<NodeDetail> hostNodes = readyDockerComposeAndLaunchNodes(req, hostIps, orderers, peers, new HashMap<>(usedPortsMap));
            // 循环等待最大超时时间，检测所有节点启动成功与否，并返回一个可用的peer和orderer
            Pair<Integer, Integer> usefulNodeIndexes = waitForAllNodesStarted(hostNodes, orderers, peers);
            Integer startedOrdererIndex = usefulNodeIndexes.getLeft();
            Integer startedPeerIndex = usefulNodeIndexes.getRight();
            if (startedOrdererIndex == -1 && !CollectionUtils.isEmpty(orderers)) {
                throw new AgentException(ConstantCode.START_NODES_ERROR);
            }
            if (startedPeerIndex == -1 && !CollectionUtils.isEmpty(peers)) {
                throw new AgentException(ConstantCode.START_NODES_ERROR);
            }

            // 保存链
            DeployResultDo deployResultDTO = getDeployResultByRequestId(req.getChainUid(), req.getRequestId());
            deployResultRepo.update(null,
                    new LambdaUpdateWrapper<DeployResultDo>()
                            .eq(DeployResultDo::getChainUid, deployResultDTO.getChainUid())
                            .eq(DeployResultDo::getRequestId, deployResultDTO.getRequestId())
                            .set(DeployResultDo::getStage, FabricDeployNodeStage.DEPLOY_SUCCESS.getIndex())
            );
            // 保存节点
            List<DeployResultNodeDo> resultNodes = hostNodes.stream()
                    .map(h -> {
                                DeployResultNodeDo deployResultNodeDo = new DeployResultNodeDo();
                                BeanUtils.copyProperties(h, deployResultNodeDo);
                                deployResultNodeDo.setResultId(deployResultDTO.getId());
                                deployResultNodeDo.setNodeType(h.isPeer() ? FabricNodeType.PEER.getIndex() : FabricNodeType.ORDERER.getIndex());
                                return deployResultNodeDo;
                            }
                    )
                    .collect(Collectors.toList());
            deployResultNodeRepo.saveBatch(resultNodes);
            // auto init gateway
            try {
                initGateway(req, orderers);
            } catch (Exception e) {
                log.error("addNode.initGateway err: {}", e.getMessage());
                throw e;
            }
        } catch (Exception e) {
            log.error("fail to  add node", e);
            saveAddNodeFailedResult(req, e);
        }
    }

    private void initGateway(BaseDeployReq req, List<NodeTriple> ordererTriples) throws Exception {
        if (CollectionUtils.isEmpty(ordererTriples)) {
            return;
        }
        String chainUid = req.getChainUid();
        Path chainRoot = pathService.getChainRoot(chainUid);
        // rebuild connection-tls.json
        Path networkConfigPath = chainRoot.resolve(FabricConfig.networkConfigPath);
        File networkConfigFile = networkConfigPath.toFile();
        log.info("=== Read networkConfig from: {}", networkConfigFile.getAbsolutePath());
        String networkConfigJsonStr = new String(Files.readAllBytes(networkConfigPath));
        log.info("=== Read networkConfig content: {}", networkConfigJsonStr);
        NetworkConfig networkConfig = JSON.parseObject(networkConfigJsonStr, NetworkConfig.class);

        Map<String, Object> channels = networkConfig.getChannels();
        Map<String, Object> orderers = networkConfig.getOrderers();
        Iterator<String> iterator = channels.keySet().iterator();
        while (iterator.hasNext()) {
            String channelId = iterator.next();
            JSONObject channelObj = JSON.parseObject(JSON.toJSONString(channels.get(channelId)));
            JSONArray ordererArr = channelObj.getJSONArray("orderers");
            for (NodeTriple nodeTriple : ordererTriples) {
                String nodeName = nodeTriple.getNodeName();
                String orgFullName = nodeTriple.getNodeOrgName();
                String nodeIp = nodeTriple.getNodeEndpoint().split(":")[0];
                String nodePort = nodeTriple.getNodeEndpoint().split(":")[1];
                // 1.给$.channels的所有channel.orderers添加所有新orderer
                if (!ordererArr.contains(nodeName)) {
                    ordererArr.add(nodeName);
                }
                // 2.给$.orderers添加所有新orderer
                if (!orderers.containsKey(nodeName)) {
                    JSONObject orderer = new JSONObject();
                    orderer.put("url", "grpcs://" + nodeIp + ":" + nodePort);
                    orderer.put("mspid", genMspId(orgFullName, false));
                    orderer.put("grpcOptions", new JSONObject());
                    JSONObject tlsCACerts = new JSONObject();
                    String tlsCa = chainRoot.resolve(String.format("organizations/ordererOrganizations/%s/orderers/%s/tls/ca.crt",
                            orgFullName, nodeName)).toFile().getAbsolutePath();
                    tlsCACerts.put("path", tlsCa);
                    orderer.put("tlsCACerts", tlsCACerts);
                    JSONObject adminPrivateKeyPEM = new JSONObject();
                    Path mspPath = chainRoot.resolve("organizations").resolve("ordererOrganizations").resolve(orgFullName).resolve("users")
                            .resolve("Admin@" + orgFullName).resolve("msp");
                    Path keystorePath = mspPath.resolve("keystore");
                    // 过滤出普通文件
                    List<Path> files = Files.list(keystorePath).filter(Files::isRegularFile).collect(Collectors.toList());
                    String keystore = files.get(0).toFile().getAbsolutePath();
                    adminPrivateKeyPEM.put("path", keystore);
                    orderer.put("adminPrivateKeyPEM", adminPrivateKeyPEM);
                    JSONObject signedCertPEM = new JSONObject();
                    signedCertPEM.put("path", mspPath.resolve("signcerts").resolve("cert.pem").toFile().getAbsolutePath());
                    orderer.put("signedCertPEM", signedCertPEM);
                    orderers.put(nodeName, orderer);
                }
            }
            channels.put(channelId, channelObj);
        }

        String newNetworkConfigJsonStr = JSON.toJSONString(networkConfig);
        log.info("=== Auto gen networkConfig: {}", newNetworkConfigJsonStr);
        log.info("=== Ready to re-persist to: {}", networkConfigFile.getAbsolutePath());
        FileUtils.write(networkConfigFile, newNetworkConfigJsonStr);

        // init new gateway
        fabricConfig.reinitializeFabricGatewayBean();
    }

    private void initCaServerContainer(DeployReq req) {
        for (String orgDomain : req.getHostNodesInOrgMap().keySet()) {
            Path caServerPath = pathService.getOrgCaServerDir(req.getChainUid(), orgDomain);
            GenCaServerReq genCaServerConfReq = buildGenCaServerDockerReq(req);
            genCaServerConfReq.setNetwork("ca_" + orgDomain.replaceAll("\\.", "_") + req.getChainUid()) ;
            genCaServerConfReq.setCaName(buildCaName(orgDomain));
            genCaServerConfReq.setContainerName(buildCaName(orgDomain).replaceAll("\\-", "."));
            genCaServerConfReq.setCaDir(caServerPath.toFile().getAbsolutePath());
            TemplateUtil.newCaServerDockerYaml(caServerPath, genCaServerConfReq);

            String cmd = String.format("cd %s && docker-compose -f docker-compose.yaml up -d", genCaServerConfReq.getCaDir());
            exec(cmd);
        }

    }

    private GenCaServerReq buildGenCaServerDockerReq(DeployReq req) {
        String orgDomain = req.getHostNodesInOrgMap().keySet().stream().findAny().get();
        GenCaServerReq genCaServerReq  = new GenCaServerReq();
        genCaServerReq.setEnv("dev");
        genCaServerReq.setCaHost(constantProperties.getCaHost());
        genCaServerReq.setCaAdmin(constantProperties.getCaAdmin());
        genCaServerReq.setCaAdminPw(constantProperties.getCaPassword());
        Integer caPort = getCaServerPort(req.getChainUid(), orgDomain);
        Integer caListenPort = availablePort(constantProperties.getCaHost(), caPort  + 1000);
        genCaServerReq.setCaPort(caPort);
        genCaServerReq.setCaListenPort(caListenPort);

        return genCaServerReq;
    }

    private void initCaServerConf(DeployReq req) {
        GenCaServerConfReq genCaServerConfReq = buildGenCaServerConfReq(req);
        for (String orgDomainName : req.getHostNodesInOrgMap().keySet()) {
            Path caServerPath = pathService.getOrgCaServerDir(req.getChainUid(), orgDomainName);
            TemplateUtil.newCaServerConfigYaml(caServerPath, genCaServerConfReq);
        }
    }

    private void intClientEnv(DeployReq req) {
        for (String orgDomain : req.getHostNodesInOrgMap().keySet()) {
            String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
            String caName = buildCaName(orgDomain);
            String caIp = constantProperties.getCaHost();
            Integer caPort = getCaServerPort(req.getChainUid(), orgDomain);
            String caAdmin = constantProperties.getCaAdmin();
            String caPassword = constantProperties.getCaPassword();

            String initClientEnvCommand = String.format("bash %s " +
                            "-p %s -c %s -i %s -t %s -a %s -w %s -o %s",
                    // init_client_env.sh
                    constantProperties.getInitClientEnvShell(),
                    // chain nodes root path
                    chainPath,
                    caName,
                    caIp,
                    caPort,
                    caAdmin,
                    caPassword,
                    orgDomain
            );
            exec(initClientEnvCommand);
        }

    }

    private void initOrg(DeployReq req) {
        String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
        String caIp = constantProperties.getCaHost();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);
        String caName = buildCaName(orgDomain);
        Integer caPort = getCaServerPort(req.getChainUid(), orgDomain);

        GenMspConfReq genMspConfReq = new GenMspConfReq();
        genMspConfReq.setCaCertName(caIp.replaceAll("\\.", "-") + "-" + caPort + "-" + caName.replaceAll("\\.", "-"));
        TemplateUtil.newConfigYaml(pathService.getChainRoot(req.getChainUid()), genMspConfReq);

        String initOrgConfigCommand = String.format("bash %s " +
                        "-p %s -o %s",
                // init_org_config.sh
                constantProperties.getInitOrgConfigShell(),
                // chain nodes root path
                chainPath,
                orgDomain
        );
        exec(initOrgConfigCommand);
    }

    private void initPeer(BaseDeployReq req) {
        String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);

        req.getHostNodesInOrgMap().forEach((domain, hostIpNodeList) -> {
            hostIpNodeList.forEach(hostIpNode -> {
                Map<String, Integer> nodeMap = hostIpNode.getPeerNameAndPorts();
                if (CollectionUtils.isEmpty(nodeMap)) {
                    return;
                }
                nodeMap.forEach((node, port) -> {
                    String peerName = node.substring(0, node.indexOf("."));
                    String initNodeConfigCommand = String.format("bash %s " +
                                    "-p %s -o %s -c %s -i %s -t %s -n %s -e %s -g %s -v %s",
                            // init_peer_config.sh
                            constantProperties.getInitPeerConfigShell(),
                            // chain nodes root path
                            chainPath,
                            orgDomain,
                            buildCaName(orgDomain),
                            constantProperties.getCaHost(),
                            getCaServerPort(req.getChainUid(), orgDomain),
                            peerName,
                            hostIpNode.getIp(),
                            String.format("%s_%s", hostIpNode.getIp(), port),
                            node
                    );
                    exec(initNodeConfigCommand);
                });
            });
        });
    }

    private void initOrderer(BaseDeployReq req) {
        String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);

        req.getHostNodesInOrgMap().forEach((domain, hostIpNodeList) -> {
            hostIpNodeList.forEach(hostIpNode -> {
                Map<String, Integer> nodeMap = hostIpNode.getOrdererNameAndPorts();
                if (CollectionUtils.isEmpty(nodeMap)) {
                    return;
                }
                nodeMap.forEach((node, port) -> {
                    String nodeName = node.substring(0, node.indexOf("."));
                    String initNodeConfigCommand = String.format("bash %s " +
                                    "-p %s -o %s -c %s -i %s -t %s -n %s -e %s -v %s",
                            // init_orderer_config.sh
                            constantProperties.getInitOrdererConfigShell(),
                            // chain nodes root path
                            chainPath,
                            orgDomain,
                            buildCaName(orgDomain),
                            constantProperties.getCaHost(),
                            getCaServerPort(req.getChainUid(), orgDomain),
                            nodeName,
                            hostIpNode.getIp(),
                            node
                    );
                    exec(initNodeConfigCommand);
                });
            });
        });
    }

    private void genConfigTx(DeployReq req) {
        GenConfigTxReq genConfigTxReq = buildGenConfigReq(req);
        Path chainRoot = pathService.getChainRoot(req.getChainUid());
        TemplateUtil.newConfigTxYaml(chainRoot, genConfigTxReq);
    }

    private void initGenesisConf(DeployReq req) {
        String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);
        String initGenesisBlockCommand = String.format("bash %s " +
                        "-p %s -o %s",
                // init_peer_config.sh
                constantProperties.getInitGenesisBlockShell(),
                // chain nodes root path
                chainPath,
                orgDomain
        );
        exec(initGenesisBlockCommand);
    }

    private HashMap<String, Set<Integer>> checkAvailableNodePort(BaseDeployReq req) {
        HashMap<String, Set<Integer>> usedPortsByIp = new HashMap<>();
        req.getHostNodesInOrgMap().forEach((orgName, hostNodeInfoList) -> {
            hostNodeInfoList.stream().forEach(hostNodeInfo -> {
                String hostIp = hostNodeInfo.getIp();
                Set<Integer> portSet = usedPortsByIp.get(hostIp);
                if (portSet == null) {
                    usedPortsByIp.put(hostIp, new HashSet<>());
                }
                Map<String, Integer> peerNameAndPorts = hostNodeInfo.getPeerNameAndPorts();
                if (!CollectionUtils.isEmpty(peerNameAndPorts)) {
                    peerNameAndPorts.forEach((peerName, port) -> {
                        int availableNodePort = getAvailablePort(hostIp, port, usedPortsByIp.get(hostIp));
                        peerNameAndPorts.put(peerName, availableNodePort);
                    });
                }
                Map<String, Integer> ordererNameAndPorts = hostNodeInfo.getOrdererNameAndPorts();
                if (!CollectionUtils.isEmpty(ordererNameAndPorts)) {
                    ordererNameAndPorts.forEach((ordererName, port) -> {
                        int availableNodePort = getAvailablePort(hostIp, port, usedPortsByIp.get(hostIp));
                        ordererNameAndPorts.put(ordererName, availableNodePort);
                    });
                }
            });
        });
        log.info("after first check usedPortsByIp: {}", JSON.toJSONString(usedPortsByIp));
        return usedPortsByIp;
    }

    private void assignNodePortBeforeAddingNode(BaseDeployReq req) {
        req.getHostNodesInOrgMap().forEach((orgName, hostNodeInfoList) -> {
            hostNodeInfoList.forEach(hostNodeInfo -> {
                String hostIp = hostNodeInfo.getIp();
                Set<Integer> portSet = usedPortsMap.get(hostIp);
                if (portSet == null) {
                    portSet = queryUsedPortsByIp(hostIp);
                    usedPortsMap.put(hostIp, portSet);
                }
                Map<String, Integer> peerNameAndPorts = hostNodeInfo.getPeerNameAndPorts();
                if (!CollectionUtils.isEmpty(peerNameAndPorts)) {
                    peerNameAndPorts.forEach((peerName, port) -> {
                        int availableNodePort = getAvailablePort(hostIp, port, usedPortsMap.get(hostIp));
                        peerNameAndPorts.put(peerName, availableNodePort);
                    });
                }
                Map<String, Integer> ordererNameAndPorts = hostNodeInfo.getOrdererNameAndPorts();
                if (!CollectionUtils.isEmpty(ordererNameAndPorts)) {
                    ordererNameAndPorts.forEach((ordererName, port) -> {
                        int availableNodePort = getAvailablePort(hostIp, port, usedPortsMap.get(hostIp));
                        ordererNameAndPorts.put(ordererName, availableNodePort);
                    });
                }
            });
        });
        log.info("assignedPortsMap: {}", JSON.toJSONString(usedPortsMap));
    }

    private byte[] calBlockHash(String blockPath) {
        try {
            byte[] bytes = FileUtils.readFileToByteArray(new File(blockPath));
            Common.Block block = Common.Block.parseFrom(bytes);
            Common.BlockHeader header = block.getHeader();
            long blockNumber = header.getNumber();
            byte[] previousHash = header.getPreviousHash().toByteArray();
            byte[] dataHash = header.getDataHash().toByteArray();

            ByteArrayOutputStream s = new ByteArrayOutputStream();
            DERSequenceGenerator seq = new DERSequenceGenerator(s);
            seq.addObject(new ASN1Integer(blockNumber));
            seq.addObject(new DEROctetString(previousHash));
            seq.addObject(new DEROctetString(dataHash));
            seq.close();

            Digest digest;
            if ("SHA3".equals(hashAlgorithm)) {
                digest = new SHA3Digest();
            } else {
                // Default to SHA2
                digest = new SHA256Digest();
            }
            byte[] retValue = new byte[digest.getDigestSize()];
            digest.update(s.toByteArray(), 0, s.size());
            digest.doFinal(retValue, 0);
            return retValue;
        } catch (Exception e) {
            log.error("calculate block Hash fail, blockPath:{}", blockPath, e);
            return null;
        }
    }

    private GenCaServerConfReq buildGenCaServerConfReq(DeployReq req) {
        String orgName = req.getHostNodesInOrgMap().keySet().stream().findFirst().get();
        String caName = buildCaName(orgName);
        String caIp = constantProperties.getCaHost();
        Integer caPort = availablePort(caIp, constantProperties.getCaPort());
        String caAdmin = constantProperties.getCaAdmin();
        String caPassword = constantProperties.getCaPassword();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);

        // 保存caPort
        ChainInfoDo chainInfoDo = new ChainInfoDo()
                .setChainUid(req.getChainUid())
                .setCaHost(caIp)
                .setCaName(caName)
                .setCaPort(caPort);
        chainInfoRepo.saveOrUpdate(chainInfoDo);

        GenCaServerConfReq genCaServerConfReq = new GenCaServerConfReq();
        genCaServerConfReq.setCaName(caName.replaceAll("\\.", "-"));
        genCaServerConfReq.setCaIp(caIp);
        genCaServerConfReq.setCaPort(caPort);
        genCaServerConfReq.setCaAdmin(caAdmin);
        genCaServerConfReq.setCaPw(caPassword);
        genCaServerConfReq.setOrgDomain(orgDomain);
        genCaServerConfReq.setCaDomain(buildCaName(orgDomain));
        return genCaServerConfReq;
    }

    private GenConfigTxReq buildGenConfigReq(DeployReq req) {
        GenConfigTxReq genConfigTxReq = new GenConfigTxReq();
        String chainPath = pathService.getChainRoot(req.getChainUid()).toFile().getAbsolutePath();
        String orgDomain = Lists.newArrayList(req.getHostNodesInOrgMap().keySet()).get(0);
        GenConfigTxReq.Org orderOrg = new GenConfigTxReq.Org();
        orderOrg.setDomain(orgDomain);
        orderOrg.setOrgName("Orderer" + orgDomain.replaceAll("\\.", "").toUpperCase());
        orderOrg.setMspId(genMspId(orgDomain, false));
        GenConfigTxReq.Org peerOrg = new GenConfigTxReq.Org();
        peerOrg.setDomain(orgDomain);
        peerOrg.setOrgName("Peer" + orgDomain.replaceAll("\\.", "").toUpperCase());
        peerOrg.setMspId(genMspId(orgDomain, true));

        List<GenConfigTxReq.Node> ordererList = new ArrayList<>();
        req.getHostNodesInOrgMap().forEach((domain, hostIpNodeList) -> {
            hostIpNodeList.forEach(hostIpNode -> {
                Map<String, Integer> nodeMap = hostIpNode.getOrdererNameAndPorts();
                nodeMap.forEach((node, port) -> {
                    GenConfigTxReq.Node ordererNode = new GenConfigTxReq.Node();
                    ordererNode.setDomain(node);
                    ordererNode.setIp(hostIpNode.getIp());
                    ordererNode.setOrgDomain(orgDomain);
                    ordererNode.setPort(String.valueOf(port));

                    ordererList.add(ordererNode);
                });
            });
        });

        List<GenConfigTxReq.Node> peerList = new ArrayList<>();
        req.getHostNodesInOrgMap().forEach((domain, hostIpNodeList) -> {
            hostIpNodeList.forEach(hostIpNode -> {
                Map<String, Integer> nodeMap = hostIpNode.getPeerNameAndPorts();
                nodeMap.forEach((node, port) -> {
                    GenConfigTxReq.Node peerNode = new GenConfigTxReq.Node();
                    peerNode.setDomain(node);
                    peerNode.setIp(hostIpNode.getIp());
                    peerNode.setOrgDomain(orgDomain);
                    peerNode.setPort(String.valueOf(port));

                    peerList.add(peerNode);
                });
            });
        });
        List<GenConfigTxReq.Channel> channelList = new ArrayList<>();
        GenConfigTxReq.Channel channel = new GenConfigTxReq.Channel();
        channel.setChannelName(req.getChannelId());
        channelList.add(channel);

        genConfigTxReq.setOrderOrg(orderOrg);
        genConfigTxReq.setOrdererList(ordererList);
        genConfigTxReq.setPeerOrg(peerOrg);
        genConfigTxReq.setPeerList(peerList);
        genConfigTxReq.setChannelList(channelList);
        genConfigTxReq.setChainPath(chainPath);
        return genConfigTxReq;
    }

    private void transferFiles2NodeHosts(BaseDeployReq req, Path organizationsPath, Set<String> hostIps) {
        String chainUid = req.getChainUid();
        Set<String> orgNames = req.getHostNodesInOrgMap().keySet();
        for (String orgName : orgNames) {
            for (HostNodeInfo hostNodeInfo : req.getHostNodesInOrgMap().get(orgName)) {
                String nodeIp = hostNodeInfo.getIp();
                hostIps.add(nodeIp);
                Map<String, Integer> peerNameAndPorts = hostNodeInfo.getPeerNameAndPorts();
                if (!CollectionUtils.isEmpty(peerNameAndPorts)) {
                    Path orgPath = organizationsPath.resolve("peerOrganizations").resolve(orgName);

                    String orgMspPath = orgPath.resolve("msp").toAbsolutePath().toString();
                    String remoteDst = String.format("%s/%s/organizations/peerOrganizations/%s", constantProperties.getInstallDir(), chainUid, orgName);
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgMspPath, remoteDst);

                    String orgTlsCaPath = orgPath.resolve("tlsca").toAbsolutePath().toString();
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgTlsCaPath, remoteDst);

                    String orgUsersPath = orgPath.resolve("users").toAbsolutePath().toString();
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgUsersPath, remoteDst);

                    for (String peerNameStr : peerNameAndPorts.keySet()) {
                        String peerPath = orgPath.resolve("peers").resolve(peerNameStr).toAbsolutePath().toString();
                        String remotePeerDst = String.format("%s/%s/organizations/peerOrganizations/%s/peers", constantProperties.getInstallDir(), chainUid, orgName);
                        ansibleService.scp(ScpTypeEnum.UP, nodeIp, peerPath, remotePeerDst);
                    }
                }
                Map<String, Integer> ordererNameAndPorts = hostNodeInfo.getOrdererNameAndPorts();
                if (!CollectionUtils.isEmpty(ordererNameAndPorts)) {
                    Path orgPath = organizationsPath.resolve("ordererOrganizations").resolve(orgName);

                    String orgMspPath = orgPath.resolve("msp").toAbsolutePath().toString();
                    String remoteDst = String.format("%s/%s/organizations/ordererOrganizations/%s", constantProperties.getInstallDir(), chainUid, orgName);
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgMspPath, remoteDst);

                    String orgTlsCaPath = orgPath.resolve("tlsca").toAbsolutePath().toString();
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgTlsCaPath, remoteDst);

                    String orgUsersPath = orgPath.resolve("users").toAbsolutePath().toString();
                    ansibleService.scp(ScpTypeEnum.UP, nodeIp, orgUsersPath, remoteDst);

                    for (String ordererNameStr : ordererNameAndPorts.keySet()) {
                        String ordererPath = orgPath.resolve("orderers").resolve(ordererNameStr).toAbsolutePath().toString();
                        String remoteOrdererDst = String.format("%s/%s/organizations/ordererOrganizations/%s/orderers", constantProperties.getInstallDir(), chainUid, orgName);
                        ansibleService.scp(ScpTypeEnum.UP, nodeIp, ordererPath, remoteOrdererDst);
                    }
                }
            }
        }
    }

    private List<NodeDetail> readyDockerComposeAndLaunchNodes(BaseDeployReq req, Set<String> hostIps,
                                                              List<NodeTriple> orderers,
                                                              List<NodeTriple> peers,
                                                              HashMap<String, Set<Integer>> usedPortsByIp) {
        String chainUid = req.getChainUid();
        Path chainRoot = pathService.getChainRoot(chainUid);
        Path chainDockerDir = pathService.getChainDockerDir(chainUid);
        String chainDockerDirStr = chainDockerDir.toFile().getAbsolutePath();
        String localGenesisPath = chainRoot.resolve("config/system-genesis-block/genesis.block").toFile().getAbsolutePath();
        String remoteGenesisBlock = String.format("%s/%s/%s", constantProperties.getInstallDir(), chainUid, "config/system-genesis-block/genesis.block");
        // 为每个节点主机生成一份空的docker-compose.yaml
        hostIps.stream().forEach(hostIp -> {
            String createDockerComposeCmd = String.format("mkdir -p %s/%s && touch %s/%s/docker-compose.yaml",
                    chainDockerDirStr, hostIp,
                    chainDockerDirStr, hostIp);
            exec(createDockerComposeCmd);
        });

        HashMap<String, List<NodeDetail>> hostNodesMap = new HashMap<>();
        req.getHostNodesInOrgMap().forEach((orgName, hostIpNodeList) -> {
            hostIpNodeList.stream().forEach(node -> {
                String hostIp = node.getIp();
                Set<Integer> portSet = usedPortsByIp.get(hostIp);
                if (!CollectionUtils.isEmpty(node.getPeerNameAndPorts())) {
                    node.getPeerNameAndPorts().forEach((peerName, port) -> {
                        List<NodeDetail> nodeDetails = hostNodesMap.get(hostIp);
                        if (CollectionUtils.isEmpty(nodeDetails)) {
                            hostNodesMap.put(hostIp, new ArrayList<>());
                        }
                        NodeDetail nodeDetail = new NodeDetail();
                        nodeDetail.setPeer(true).setNodeFullName(peerName).setNetworkName(chainUid)
                                .setNodePort(port).setNodeIp(hostIp).setOrgFullName(orgName)
                                .setMspId(genMspId(orgName, true));
                        nodeDetail.setChainCodePort(getAvailablePort(hostIp, port, portSet));
                        nodeDetail.setOperationsPort(getAvailablePort(hostIp, port, portSet));
                        hostNodesMap.get(hostIp).add(nodeDetail);
                        peers.add(new NodeTriple(peerName, orgName, hostIp + ":" + nodeDetail.getNodePort()));
                    });
                }
                if (!CollectionUtils.isEmpty(node.getOrdererNameAndPorts())) {
                    node.getOrdererNameAndPorts().forEach((ordererName, port) -> {
                        List<NodeDetail> nodeDetails = hostNodesMap.get(hostIp);
                        if (CollectionUtils.isEmpty(nodeDetails)) {
                            hostNodesMap.put(hostIp, new ArrayList<>());
                        }
                        NodeDetail nodeDetail = new NodeDetail();
                        nodeDetail.setPeer(false).setNodeFullName(ordererName).setNetworkName(chainUid)
                                .setNodePort(port).setNodeIp(hostIp).setOrgFullName(orgName)
                                .setMspId(genMspId(orgName, false));
                        nodeDetail.setOperationsPort(getAvailablePort(hostIp, port, portSet));
                        hostNodesMap.get(hostIp).add(nodeDetail);
                        orderers.add(new NodeTriple(ordererName, orgName, hostIp + ":" + nodeDetail.getNodePort()));
                    });
                }
            });
        });
        log.info("hostNodesMap: {}", JSON.toJSONString(hostNodesMap));
        log.info("after final check usedPortsByIp: {}", JSON.toJSONString(usedPortsByIp));

        hostIps.parallelStream().forEach(hostIp -> {
            List<NodeDetail> nodeDetails = hostNodesMap.get(hostIp);
            DumperOptions dumperOptions = new DumperOptions();
            dumperOptions.setDefaultFlowStyle(DumperOptions.FlowStyle.BLOCK);
            Yaml yaml = new Yaml(dumperOptions);
            String dockerComposeFile = String.format("%s/%s/docker-compose.yaml", chainDockerDirStr, hostIp);
            File file = new File(dockerComposeFile);
            Map<String, Object> yamlData = null;
            try {
                yamlData = yaml.load(new FileInputStream(file));
            } catch (FileNotFoundException e) {
                throw new AgentException(ConstantCode.DOCKER_COMPOSE_YAML_ERROR);
            }
            for (NodeDetail nodeDetail : nodeDetails) {
                String networkName = nodeDetail.getNetworkName();
                String nodeFullName = nodeDetail.getNodeFullName();
                Integer nodePort = nodeDetail.getNodePort();
                String mspId = nodeDetail.getMspId();
                Integer operationsPort = nodeDetail.getOperationsPort();
                String orgFullName = nodeDetail.getOrgFullName();
                if (yamlData == null) {
                    yamlData = new LinkedHashMap<>();
                    yamlData.put("version", "3.7");

                    yamlData.put("volumes", new LinkedHashMap<>());

                    yamlData.put("networks", new LinkedHashMap<>());
                    Map<String, Object> networks = (Map<String, Object>) yamlData.get("networks");
                    networks.put("dev", new LinkedHashMap<>());
                    Map<String, Object> dev = (Map<String, Object>) networks.get("dev");
                    dev.put("name", networkName);

                    yamlData.put("services", new LinkedHashMap<>());
                }
                Map<String, Object> volumes = (Map<String, Object>) yamlData.get("volumes");
                volumes.put(nodeFullName, null);
                Map<String, Object> services = (Map<String, Object>) yamlData.get("services");
                if (nodeDetail.isPeer()) {
                    LinkedHashMap<Object, Object> node = new LinkedHashMap<>();
                    services.put(nodeFullName, node);
                    node.put("container_name", nodeFullName);
                    node.put("image", "hyperledger/fabric-peer:" + fabricVersion);
                    // environment
                    node.put("environment", Lists.newArrayList("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
                            "CORE_VM_DOCKER_HOSTCONFIG_NETWORKMODE=" + networkName,
                            "FABRIC_LOGGING_SPEC=" + fabricLogLevel,
                            "CORE_PEER_TLS_ENABLED=true",
                            "CORE_PEER_PROFILE_ENABLED=false",
                            "CORE_PEER_TLS_CERT_FILE=/etc/hyperledger/fabric/tls/server.crt",
                            "CORE_PEER_TLS_KEY_FILE=/etc/hyperledger/fabric/tls/server.key",
                            "CORE_PEER_TLS_ROOTCERT_FILE=/etc/hyperledger/fabric/tls/ca.crt",
                            "CORE_PEER_ID=" + nodeFullName,
                            "CORE_PEER_ADDRESS=0.0.0.0:" + nodePort,
                            "CORE_PEER_LISTENADDRESS=0.0.0.0:" + nodePort,
                            "CORE_PEER_CHAINCODELISTENADDRESS=0.0.0.0:" + nodeDetail.getChainCodePort(),
                            "CORE_PEER_GOSSIP_BOOTSTRAP=" + hostIp + ":" + nodePort,
                            // 保证浏览器可以发现该peer
                            "CORE_PEER_GOSSIP_EXTERNALENDPOINT=" + hostIp + ":" + nodePort,
                            "CORE_PEER_LOCALMSPID=" + mspId,
                            "CORE_OPERATIONS_LISTENADDRESS=0.0.0.0:" + operationsPort));
                    String peerPath = String.format("%s/%s/organizations/peerOrganizations/%s/peers/%s/", constantProperties.getInstallDir(), chainUid, orgFullName, nodeFullName);
                    String mspPath = peerPath + "msp";
                    String tlsPath = peerPath + "tls";
                    node.put("volumes", Lists.newArrayList("/var/run/docker.sock:/host/var/run/docker.sock",
                            mspPath + ":/etc/hyperledger/fabric/msp",
                            tlsPath + ":/etc/hyperledger/fabric/tls",
                            peerPath + ":/var/hyperledger/production"));
                    node.put("working_dir", "/opt/gopath/src/github.com/hyperledger/fabric/peer");
                    node.put("command", "peer node start");
                    node.put("ports", Lists.newArrayList(nodePort + ":" + nodePort,
                            operationsPort + ":" + operationsPort));
                    node.put("networks", Lists.newArrayList("dev"));
                } else {
                    LinkedHashMap<Object, Object> node = new LinkedHashMap<>();
                    services.put(nodeFullName, node);
                    node.put("container_name", nodeFullName);
                    node.put("image", "hyperledger/fabric-orderer:" + fabricVersion);
                    // environment
                    node.put("environment", Lists.newArrayList("CORE_VM_ENDPOINT=unix:///host/var/run/docker.sock",
                            "FABRIC_LOGGING_SPEC=" + fabricLogLevel,
                            "ORDERER_GENERAL_LISTENADDRESS=0.0.0.0",
                            "ORDERER_GENERAL_LISTENPORT=" + nodePort,
                            "ORDERER_GENERAL_GENESISMETHOD=file",
                            "ORDERER_GENERAL_GENESISFILE=/var/hyperledger/orderer/orderer.genesis.block",
                            "ORDERER_GENERAL_LOCALMSPID=" + mspId,
                            "ORDERER_GENERAL_LOCALMSPDIR=/var/hyperledger/orderer/msp",
                            "ORDERER_OPERATIONS_LISTENADDRESS=0.0.0.0:" + operationsPort,
                            "ORDERER_GENERAL_TLS_ENABLED=true",
                            "ORDERER_GENERAL_TLS_PRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
                            "ORDERER_GENERAL_TLS_CERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
                            "ORDERER_GENERAL_TLS_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]",
                            "ORDERER_KAFKA_TOPIC_REPLICATIONFACTOR=1",
                            "ORDERER_KAFKA_VERBOSE=true",
                            "ORDERER_GENERAL_CLUSTER_CLIENTCERTIFICATE=/var/hyperledger/orderer/tls/server.crt",
                            "ORDERER_GENERAL_CLUSTER_CLIENTPRIVATEKEY=/var/hyperledger/orderer/tls/server.key",
                            "ORDERER_GENERAL_CLUSTER_ROOTCAS=[/var/hyperledger/orderer/tls/ca.crt]"
                    ));
                    node.put("working_dir", "/opt/gopath/src/github.com/hyperledger/fabric");
                    node.put("command", "orderer");
                    String ordererPath = String.format("%s/%s/organizations/ordererOrganizations/%s/orderers/%s/",
                            constantProperties.getInstallDir(), chainUid, orgFullName, nodeFullName);
                    String mspPath = ordererPath + "msp";
                    String tlsPath = ordererPath + "tls";
                    node.put("volumes", Lists.newArrayList(remoteGenesisBlock + ":/var/hyperledger/orderer/orderer.genesis.block",
                            mspPath + ":/var/hyperledger/orderer/msp",
                            tlsPath + ":/var/hyperledger/orderer/tls",
                            ordererPath + ":/var/hyperledger/production/orderer"));
                    node.put("ports", Lists.newArrayList(nodePort + ":" + nodePort,
                            operationsPort + ":" + operationsPort));
                    node.put("networks", Lists.newArrayList("dev"));
                }
                FileWriter writer = null;
                try {
                    writer = new FileWriter(file);
                } catch (IOException e) {
                    throw new AgentException(ConstantCode.DOCKER_COMPOSE_YAML_ERROR);
                }
                yaml.dump(yamlData, writer);
            }
            // scp genesis.block
            ansibleService.exec(hostIp, String.format("mkdir -p `dirname %s`", remoteGenesisBlock));
            ansibleService.scp(ScpTypeEnum.UP, hostIp, localGenesisPath, remoteGenesisBlock);
            // scp docker-compose.yaml
            String remoteDst = String.format("%s/%s", constantProperties.getInstallDir(), chainUid);
            ansibleService.scp(ScpTypeEnum.UP, hostIp, dockerComposeFile, remoteDst);
            String dockerComposeCmd = String.format("docker-compose -f %s/%s/docker-compose.yaml up -d",
                    constantProperties.getInstallDir(), chainUid);
            // start nodes with docker-compose
            log.info("starting nodes on host: {}", hostIp);
            ansibleService.exec(hostIp, dockerComposeCmd);
        });

        // 展开每个节点
        List<NodeDetail> hostNodes = hostNodesMap.values().stream().flatMap(List::stream).collect(Collectors.toList());
        return hostNodes;
    }

    private void createAndJoinChannel(DeployReq req, List<NodeTriple> orderers,
                                      List<NodeTriple> peers,
                                      Pair<Integer, Integer> usefulNodeIndexes) {
        String chainUid = req.getChainUid();
        String channelId = req.getChannelId();
        Integer peerIndex = usefulNodeIndexes.getRight();
        String peerName = peers.get(peerIndex).getNodeName();
        String peerOrgName = peers.get(peerIndex).getNodeOrgName();
        String peerEndpoint = peers.get(peerIndex).getNodeEndpoint();
        Integer ordererIndex = usefulNodeIndexes.getLeft();
        String ordererName = orderers.get(ordererIndex).getNodeName();
        String ordererOrgName = orderers.get(ordererIndex).getNodeOrgName();
        String ordererEndpoint = orderers.get(ordererIndex).getNodeEndpoint();

        // 3.1.创建应用通道tx交易文件
        // -c channel1 -p ./NODES_ROOT/1ef243fd87
        String createTxCommand = String.format("bash %s -c %s -p %s",
                // create_app_channel_tx.sh shell script
                constantProperties.getCreateAppChannelTxShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath()
        );
        exec(createTxCommand);
        // 3.2.创建应用通道block区块文件
        // -c channel1 -p ./NODES_ROOT/1ef243fd87 -o 192.168.3.128:7050 -j orderer0.org1.example.com -r org1.example.com -m Org1MSP -n 192.168.3.128:7051 -s peer0.org1.example.com -g org1.example.com
        String peerMspId = genMspId(peerOrgName, true);
        String createBlockCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // create_app_channel_block.sh shell script
                constantProperties.getCreateAppChannelBlockShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(createBlockCommand);
        // 3.3.peer节点加入通道
        trackDeploymentStage(chainUid, FabricDeployStage.JOIN_CHANNEL);
        for (NodeTriple peerTriple : peers) {
            String peerName0 = peerTriple.getNodeName();
            String peerOrgName0 = peerTriple.getNodeOrgName();
            String peerEndpoint0 = peerTriple.getNodeEndpoint();
            String peerMspId0 = genMspId(peerOrgName0, true);
            String joinCommand = String.format("bash %s -c %s -p %s -m %s -n %s -s %s -g %s",
                    // join_app_channel.sh shell script
                    constantProperties.getJoinAppChannelShell(),
                    // app channel id
                    channelId,
                    // chain nodes root path
                    pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                    peerMspId0,
                    peerEndpoint0,
                    peerName0,
                    peerOrgName0
            );
            exec(joinCommand);
        }
        // 3.4.获取应用通道最近的配置块
        String fetchCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // fetch_app_channel_config.sh shell script
                constantProperties.getFetchAppChannelConfigShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(fetchCommand);
        // 3.5.配置应用通道
        String anchorPeersStr = getAnchorPeersStr(peers);
        String configCommand = String.format("bash %s -c %s -p %s -m %s -x %s",
                // config_app_channel.sh shell script
                constantProperties.getConfigAppChannelShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                peerMspId,
                anchorPeersStr
        );
        exec(configCommand);
        // 3.6.更改应用通道
        String updateCommand = String.format("bash %s -c %s -p %s -o %s -j %s -r %s -m %s -n %s -s %s -g %s",
                // update_app_channel.sh shell script
                constantProperties.getUpdateAppChannelShell(),
                // app channel id
                channelId,
                // chain nodes root path
                pathService.getChainRoot(chainUid).toFile().getAbsolutePath(),
                ordererEndpoint,
                ordererName,
                ordererOrgName,
                peerMspId,
                peerEndpoint,
                peerName,
                peerOrgName
        );
        exec(updateCommand);
    }

    private void exec(String command) {
        ExecuteResult result = JavaCommandExecutor.executeCommand(command, constantProperties.getExecShellTimeout());
        if (result.failed()) {
            throw new AgentException(ConstantCode.EXEC_ERROR.attach(result.getExecuteOut()));
        }
    }

    private String genMspId(String orgName, boolean forPeer) {
        if (forPeer) {
            return orgName.replaceAll("\\.", "") + "MSP";
        }
        return "Orderer" + orgName.replaceAll("\\.", "") + "MSP";
    }

    private String getAnchorPeersStr(List<NodeTriple> peers) {
        // {\"host\":\"192.168.3.128\",\"port\":7051}
        StringBuffer anchorPeersStr = new StringBuffer();
        peers.stream().forEach(peer -> {
            String peerEndpoint = peer.getNodeEndpoint();
            String peerIp = peerEndpoint.split(":")[0];
            String peerPort = peerEndpoint.split(":")[1];
            String peerStr = String.format("'{\"host\":\"%s\",\"port\":%s}", peerIp, peerPort);
            anchorPeersStr.append(peerStr).append(",'");
        });
        return anchorPeersStr.deleteCharAt(anchorPeersStr.length() - 2).toString();
    }

    private Pair<Integer, Integer> waitForAllNodesStarted(List<NodeDetail> hostNodes,
                                                          List<NodeTriple> orderers,
                                                          List<NodeTriple> peers) {
        long startNodesTimeout = constantProperties.getStartNodesTimeout();
        long startTime = System.currentTimeMillis();
        int startedOrdererIndex = -1;
        int startedPeerIndex = -1;
        List<NodeDetail> replicationNodeList = new ArrayList<>(hostNodes);
        Iterator<NodeDetail> iterator = replicationNodeList.iterator();
        while (true) {
            try {
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                throw new RuntimeException(e);
            }
            while (iterator.hasNext()) {
                NodeDetail nodeDetail = iterator.next();
                String hostIp = nodeDetail.getNodeIp();
                Integer port = nodeDetail.getNodePort();
                boolean notInUse = ansibleService.checkPortInUse(hostIp, port);
                // if started, skip and return useful indexes
                if (!notInUse) {
                    log.info("== node started already: {}:{}", hostIp, port);
                    if (startedOrdererIndex == -1) {
                        for (int i = 0; i < orderers.size(); i++) {
                            NodeTriple triple = orderers.get(i);
                            if ((hostIp + ":" + port).equals(triple.getNodeEndpoint())) {
                                startedOrdererIndex = i;
                                break;
                            }
                        }
                    }
                    if (startedPeerIndex == -1) {
                        for (int i = 0; i < peers.size(); i++) {
                            NodeTriple triple = peers.get(i);
                            if ((hostIp + ":" + port).equals(triple.getNodeEndpoint())) {
                                startedPeerIndex = i;
                                break;
                            }
                        }
                    }
                    iterator.remove();
                } else {
                    log.info("== node not started yet: {}:{}", hostIp, port);
                }
            }
            if (CollectionUtils.isEmpty(replicationNodeList)) {
                // removed all nodes
                log.info("== all nodes started successfully");
                break;
            }
            long endTime = System.currentTimeMillis();
            if (endTime - startTime >= startNodesTimeout) {
                log.warn("== nodes start time over");
                break;
            }
        }
        return Pair.of(startedOrdererIndex, startedPeerIndex);
    }

    private int getAvailablePort(String hostIp, int port, Set<Integer> portSet) {
        // port可能是已经占用的了，所以进行一次检测，返回一个可用的port
        port = ansibleService.getAvailablePort(hostIp, port);
        while (portSet.contains(port)) {
            port++;
            port = ansibleService.getAvailablePort(hostIp, port);
        }
        portSet.add(port);
        return port;
    }

    private void trackDeploymentStage(String chainUid, FabricDeployStage deployStage) {
        // 更新部署阶段
        deployResultRepo.update(null,
                new LambdaUpdateWrapper<DeployResultDo>()
                        .eq(DeployResultDo::getChainUid, chainUid)
                        .eq(DeployResultDo::getRequestId, chainUid)
                        .set(DeployResultDo::getStage, deployStage.getIndex())
        );
    }

    private void trackDeployNodeStage(String chainUid, String requestId, FabricDeployNodeStage deployStage) {
        DeployResultDo deployResultDo = getDeployResultByRequestId(chainUid, requestId);
        deployResultDo.setStage(deployStage.getIndex());
    }

    private DeployResultDo getDeployResultByRequestId(String chainUid, String requestId) {
        DeployResultDo deployResultDo = deployResultRepo.getOne(
                new LambdaQueryWrapper<DeployResultDo>()
                        .eq(DeployResultDo::getChainUid, chainUid)
                        .eq(DeployResultDo::getRequestId, requestId)
        );
        if (deployResultDo != null) {
            return deployResultDo;
        }
        deployResultDo = new DeployResultDo().setChainUid(chainUid).setRequestId(requestId);
        deployResultRepo.save(deployResultDo);
        return deployResultDo;
    }

    private Integer getCaServerPort(String chainUid, String orgDomain) {
        LambdaQueryWrapper<ChainInfoDo> caQuery = new LambdaQueryWrapper<>();
        caQuery.eq(ChainInfoDo::getChainUid, chainUid);
        caQuery.eq(ChainInfoDo::getCaName, buildCaName(orgDomain));
        List<ChainInfoDo> chainInfoList = chainInfoRepo.list(caQuery);
        if (chainInfoList == null || chainInfoList.size() == 0) {
            return readCaPortFromYaml(chainUid, orgDomain);
        }
        return chainInfoList.get(0).getCaPort();
    }

    private Integer readCaPortFromYaml(String chainUid, String orgDomain) {
        try {
            Path caServerPath = pathService.getOrgCaServerDir(chainUid, orgDomain);
            Path caServerConfigPath = caServerPath.resolve("fabric-ca-server-config.yaml");
            Yaml yaml = new Yaml();
            Map<String, Object> yamlMap = yaml.load(Files.newInputStream(caServerConfigPath));
            Object portObj = yamlMap.get("port");
            if (portObj == null) {
                throw new RuntimeException("ca server port is not set");
            }
            Integer port = Integer.parseInt(portObj.toString());
            return port;
        } catch (Exception e) {
            throw new RuntimeException("fail to read ca server port");
        }
    }

    private String buildDockerComposeDirectory(String chainUid, String hostIp) {
        Path chainRoot = pathService.getChainRoot(chainUid);
        String chainRootStr = chainRoot.toFile().getAbsolutePath();
        String directory = String.format("%s/%s", chainRootStr, hostIp);
        File dir = new File(directory);
        if (!dir.exists()) {
            dir.mkdirs();
        }
        String prefix = "deploy";
        File[] deployDirList = dir.listFiles(f -> f.isDirectory() && f.getName().startsWith(prefix));
        int deployIndex = 0;
        String lastIndexStr = null;
        if (deployDirList != null && deployDirList.length > 0) {
            String lastDeployDir = deployDirList[deployDirList.length - 1].getName();
            if (lastDeployDir.length() > prefix.length()) {
                lastIndexStr = lastDeployDir.substring(prefix.length());
            }
        }
        if (lastIndexStr != null) {
            Integer lastIndex = Integer.parseInt(lastIndexStr, 10);
            deployIndex = lastIndex + 1;
        }
        return prefix + String.format("%03d", deployIndex);
    }

    @Transactional(rollbackFor = Exception.class)
    private void saveAddNodeFailedResult(BaseDeployReq req, Exception e) {
        DeployResultDo deployResultDo = getDeployResultByRequestId(req.getChainUid(), req.getRequestId());
        deployResultDo.setStage(FabricDeployNodeStage.DEPLOY_FAILED.getIndex());
        deployResultDo.setError(e.getMessage());
        deployResultRepo.updateById(deployResultDo);
        List<DeployResultNodeDo> resultNodes = buildAllHostNodeResultList(req);
        resultNodes.forEach(node -> node.setResultId(deployResultDo.getId()));
        deployResultNodeRepo.saveBatch(resultNodes);
    }

    private List<DeployResultNodeDo> buildAllHostNodeResultList(BaseDeployReq req) {
        List<DeployResultNodeDo> resultNodes = new ArrayList<>();
        req.getHostNodesInOrgMap().forEach((domain, hostList) -> {
            hostList.forEach((host) -> {
                String hostIp = host.getIp();
                List<DeployResultNodeDo> orderNodeResultList = buildNodeResultList(host.getOrdererNameAndPorts(), domain, hostIp, FabricNodeType.ORDERER.getIndex());
                if (orderNodeResultList != null) {
                    resultNodes.addAll(orderNodeResultList);
                }
                List<DeployResultNodeDo> peerNodeResultList = buildNodeResultList(host.getPeerNameAndPorts(), domain, hostIp, FabricNodeType.PEER.getIndex());
                if (peerNodeResultList != null) {
                    resultNodes.addAll(peerNodeResultList);
                }
            });
        });
        return resultNodes;
    }

    private List<DeployResultNodeDo> buildNodeResultList(Map<String, Integer> nodeNamePortMap, String domain, String hostIp, Integer nodeType) {
        if (nodeNamePortMap == null) {
            return null;
        }
        List<DeployResultNodeDo> resultNodes = new ArrayList<>();
        nodeNamePortMap.forEach((nodeName, port) -> {
            DeployResultNodeDo deployResultNodeDo = new DeployResultNodeDo();
            deployResultNodeDo.setNodeIp(hostIp);
            deployResultNodeDo.setNodePort(port);
            deployResultNodeDo.setNodeFullName(nodeName);
            deployResultNodeDo.setNodeType(nodeType);
            deployResultNodeDo.setOrgFullName(domain);
            deployResultNodeDo.setMspId(genMspId(domain, FabricNodeType.PEER.getIndex().equals(nodeType)));
            resultNodes.add(deployResultNodeDo);
        });
        return resultNodes;
    }

    private Set<Integer> queryUsedPortsByIp(String ip) {
        Set<Integer> portSet =  queryNodeUsedPortsByIp(ip);
        Set<Integer> caPortSet =  queryCaServerUsedPortsByIp(ip);
        portSet.addAll(caPortSet);
        return portSet;
    }

    @NotNull
    private Set<Integer> queryNodeUsedPortsByIp(String ip) {
        List<DeployResultNodeDo> nodes = deployResultNodeRepo.list(
                new LambdaQueryWrapper<DeployResultNodeDo>()
                        .eq(DeployResultNodeDo::getNodeIp, ip)
        );
        Set<Integer> portSet = new HashSet<>();
        if (CollectionUtils.isEmpty(nodes)) {
            return portSet;
        }
        nodes.forEach(node -> {
           portSet.add(node.getNodePort());
           portSet.add(node.getOperationsPort());
           portSet.add(node.getChainCodePort());
        });
        return portSet;
    }

    @NotNull
    private Set<Integer> queryCaServerUsedPortsByIp(String ip) {
        LambdaQueryWrapper<ChainInfoDo> caQuery = new LambdaQueryWrapper<>();
        caQuery.eq(ChainInfoDo::getCaHost, ip);
        List<ChainInfoDo> chainInfoList = chainInfoRepo.list(caQuery);
        Set<Integer> portSet = new HashSet<>();
        if (CollectionUtils.isEmpty(chainInfoList)) {
            return portSet;
        }
        chainInfoList.forEach(node -> {
            portSet.add(node.getCaPort());
        });
        return portSet;
    }

    private String buildCaName(String orgDomain) {
        if (orgDomain == null) {
            throw new AgentException(ConstantCode.PARAM_EXCEPTION.attach("orgDomain cannot be null"));
        }
        return "ca-" + orgDomain.replaceAll("\\.", "-");
    }
}
