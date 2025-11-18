/**
 * Copyright 2014-2021  the original author or authors.
 * <p>
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * <p>
 * http://www.apache.org/licenses/LICENSE-2.0
 * <p>
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package com.lingshu.fabric.agent.config.properties;

import lombok.Data;
import lombok.extern.log4j.Log4j2;
import org.apache.commons.lang3.StringUtils;
import org.apache.commons.lang3.tuple.MutablePair;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

import java.io.File;
import java.math.BigInteger;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;

import static java.io.File.separator;

/**
 * constants in yml and static constants
 */
@Log4j2
@Data
@Component
@ConfigurationProperties(prefix = ConstantProperties.CONSTANT_PREFIX)
public class ConstantProperties {

    // constant
    public static final String CONSTANT_PREFIX = "constant";
    public static final String CONTRACT_NAME_ZERO = "0x00000000";
    public static final String ADDRESS_DEPLOY = "0x0000000000000000000000000000000000000000";
    public static final String ADDRESS_BUILTIN = "0x000000000000000000000000000000000000";
    public static final String ADDRESS_BUILTIN_NS = "0x0000000000000000000000000000000000001004";
    public static final int PUBLICKEY_LENGTH = 130;
    public static final int ADDRESS_LENGTH = 42;
    public static final String HAS_ROLE_ADMIN = "hasRole('admin')";
    public static final String HAS_ROLE_ADMIN_OR_DEVELOPER = "hasRole('admin') or hasRole('developer')";

    public static final String PARAM_APP_KEY = "appKey";
    public static final String PARAM_APP_SECRET = "appSecret";
    public static final String PARAM_TIMESTAMP = "timestamp";
    public static final String PARAM_SIGNATURE = "signature";
    public static final String PARAM_IS_TRANSFER_ENCRYPT = "isTransferEncrypt";
    public static final String ENCRYPT_FALSE = "false";

    // 保存证书的 organizations 目录
    private String organizationsPath;

    private boolean developerModeEnable = false;
    private boolean deployedModifyEnable = true;
    private BigInteger transRetainMax = new BigInteger("10000");
    /**
     * y:year, M:month, d:day of month, h:hour, m:minute, n:forever valid
     */
    private String ledgerInvalidGrayscaleValue = "1M";
    private String notSupportFrontIp;

    /**
     * block into
     */
    private BigInteger blockRetainMax = new BigInteger("10000");
    private BigInteger pullBlockInitCnts = new BigInteger("100");
    /**
     * 20 mills
     */
    private Long pullBlockSleepTime = 20L;
    private Boolean isBlockPullFromZero = true;

    /**
     * receive http request
     * unit: seconds
     */
    private Integer authTokenMaxAge = 900;
    private Boolean isUseSecurity = true;
    private String ignoreCheckFront = null;
    /**
     * verification code settings
     * unit: seconds
     */
    private Integer verificationCodeMaxAge = 300;
    private Boolean enableVerificationCode = true;
    private String verificationCodeValue = "8888";

    /**
     * front http request
     */
    private String frontUrl = "http://%1s:%2d/front/%3s";
    private Integer contractDeployTimeOut = 30000;
    private Integer httpTimeOut = 10000;
    private Integer interactTimeOut = 8000;
    private Boolean isPrivateKeyEncrypt = true;
    private Integer maxRequestFail = 3;
    private Long sleepWhenHttpMaxFail = 60000L;

    /**
     * transaction monitor
     * unit: seconds
     */
    private Long transMonitorTaskFixedRate = 60000L;
    private Integer monitorInfoRetainMax = 10000;
    private Long analysisSleepTime = 200L;
    private Boolean isMonitorIgnoreUser = false;
    private Boolean isMonitorIgnoreContract = false;
    private Integer monitorUnusualMaxCount = 20;

    /**
     * alert mail interval
     */
    private Integer auditMonitorTaskFixedDelay = 300000;
    private Integer nodeStatusMonitorTaskFixedDelay = 60000;
    private Integer certMonitorTaskFixedDelay = 300000;
    /**
     * application integration
     */
    private long appRequestTimeOut = 300000;
    private boolean isTransferEncrypt = true;

    /**
     * default resetLedgerListCycle interval gap, default 15000ms(15s)
     */
    private long resetLedgerListCycle = 15000;
    /**
     * pull block statistic interval: ms
     */
    private BigInteger statBlockRetainMax = new BigInteger("100000");
    private Integer statBlockFixedDelay = 5000;
    private Integer statBlockPageSize = 10;
    /**
     * enable pull external account(user address) and contract from block
     */
    private Boolean enableExternalFromBlock = true;

    //******************* Add in v1.4.0 start. *******************
    public static final boolean RETURN_EXECUTE_LOG = true;
    public static final int LEAST_VALIDATOR_TWO = 2;

    private int deployMode = 0;
    private int validInvitationCodeDay = 7;
    private int deployType = 0;
    private String bspSignAddress = "127.0.0.1:5004";

    /**
     * shell script
     */
    private String nodeOperateShell = "./script/deploy/host_operate.sh";
    private String buildChainShell = "./script/deploy/chain_tools.sh";
    // private String genAgencyShell = "./script/deploy/gen_agency_cert.sh";
    private String genAgencyMemberKeyShell = "./script/deploy/gen_agency_member_key.sh";
    private String genAgencyMemberCertShell = "./script/deploy/gen_agency_member_cert.sh";
    // private String genNodeShell = "./script/deploy/gen_node_cert.sh";
    private String scpShell = "./script/deploy/file_trans_util.sh";
    private String hostCheckShell = "./script/deploy/host_check.sh";
    private String dockerCheckShell = "./script/deploy/host_docker_check.sh";
    private String dockerPullCdnShell = "./script/deploy/host_docker_cdn.sh";
    private String hostDockerTcpShell = "./script/deploy/host_docker_tcp.sh";
    private String hostInitShell = "./script/deploy/host_init_shell.sh";
    private String hostCheckPortShell = "./script/deploy/host_check_port.sh";
    private String hostMetricShell = "./script/deploy/host_metric.sh";
    /**
     * to support | & > $
     */
    private String ansibleImageCheckShell = "./script/deploy/check_image_exist.sh";
    private String ansibleContainerCheckShell = "./script/deploy/check_container_exist.sh";
    private String hostCheckIpShell = "./script/deploy/host_check_ifconfig.sh";

    private String lingShuChainBinary =  "";

    // 节点信息
    private String netInfo = "./script/util/netInfo.sh";
    private String allInfo = "./script/util/allInfo.sh";


    // fabric
    private String caHost = "127.0.0.1";
    private Integer caPort = 7054;
    private String caAdmin = "admin";
    private String caPassword = "adminpw";
    private String metadataSource = "meta-infs/end2endit";
    private String chainCodePolicy = "chaincodeendorsementpolicy.yaml";
    private String initClientEnvShell = "./script/deploy/init_client_env.sh";
    private String initOrgConfigShell = "./script/deploy/init_org_config.sh";
    private String initPeerConfigShell = "./script/deploy/init_peer_config.sh";
    private String initOrdererConfigShell = "./script/deploy/init_orderer_config.sh";
    private String initGenesisBlockShell = "./script/deploy/init_genesis_block.sh";
    private String createAppChannelTxShell = "./script/deploy/create_app_channel_tx.sh";
    private String createAppChannelBlockShell = "./script/deploy/create_app_channel_block.sh";
    private String joinAppChannelShell = "./script/deploy/join_app_channel.sh";
    private String fetchAppChannelConfigShell = "./script/deploy/fetch_app_channel_config.sh";
    private String configAppChannelShell = "./script/deploy/config_app_channel.sh";
    private String updateAppChannelShell = "./script/deploy/update_app_channel.sh";
    private String packageChainCodeShell = "./script/chaincode/package_chain_code.sh";
    private String installChainCodeShell = "./script/chaincode/install_chain_code.sh";
    private String approveAndCommitChainCodeShell = "./script/chaincode/approve_and_commit_chain_code.sh";
    private String approveChainCodeShell = "./script/chaincode/approve_chain_code.sh";
    private String commitChainCodeShell = "./script/chaincode/commit_chain_code.sh";
    private String invokeChainCodeShell = "./script/chaincode/invoke_chain_code.sh";
    private String getChainCodeSequenceShell = "./script/chaincode/get_chain_code_sequence.sh";

    // default port
    private int defaultChainId = 1;

    // timeout config (ms)
    // check docker installed and active 1min
    private long execDockerCheckTimeout = 55 * 1000L;
    // check memory dependency, check container exist, check image exist
    private long execHostCheckTimeout = 55 * 1000L;
    // check port in use
    private long execHostCheckPortTimeout = 50 * 1000L;
    // async init host time out. 5min
    private long execHostInitTimeout = 5 * 60 * 1000L;
    // generate chain config and scp to host
    private long execHostConfigTimeout = 40 * 1000L;
    // generate chain config
    private long execBuildChainTimeout = 40 * 1000L;
    // docker command time out
    private long dockerRestartPeriodTime = 30 * 1000L;
    // common shell exec time out
    private long execShellTimeout = 10 * 60 * 1000L;
    // scp command concurrent await time
    private long execScpTimeout = 10 * 1000L;
    // add node concurrent await time
    private long execAddNodeTimeout = 40 * 1000L;
    private String installDir = "/home/songzehao/opt";
    // 所有节点主机上的节点并行启动的最大超时时间
    private long startNodesTimeout = 40 * 1000L;
    private Integer frontPortStart = 5001;
    private Integer rpcPortStart = 56656;
    private Integer p2pPortStart = 46656;

    private String[] permitUrlArray = new String[]{"/account/login", "/account/pictureCheckCode", "/login","/user/privateKey/**", "/encrypt", "/version"};
    private String dockerRepository= "registry.cn-hangzhou.aliyuncs.com/energy/chainfront";
    // private String imageTagUpdateUrl = "https://registry.hub.docker.com/v1/repositories/%s/tags";
    private String dockerRegistryMirror = "";
    private String nodesRootDir = "NODES_ROOT";
    private String nodesRootTmpDir = "NODES_ROOT_TMP";

    // 2022-10-25
    // 主机和节点 监控数据接口
    private String ratioInterval = "5m";
    private String prometheusRestQueryRangeApiUrl = "http://192.168.1.66:9090/api/v1/query_range?query={query}&start={start}&end={end}&step={step}";
    private String prometheusRestQueryApiUrl = "http://192.168.1.66:9090/api/v1/query?query={query}&time={time}";

    private Integer dayLimit = 31;
    private Integer metricLimit = 500;

    /**
     * Docker client connect daemon ip with proxy ip.
     */
    private Map<String, MutablePair<String, Integer>> dockerProxyMap = new ConcurrentHashMap<>();

    /**
     *
     * @param injectedValue
     * @param defaultValue
     * @return
     */
    private static String initDirectory(String injectedValue, String defaultValue){
        String newDirectory = injectedValue;

        if (StringUtils.isBlank(newDirectory)) {
            newDirectory = defaultValue;
        }

        if (newDirectory.trim().endsWith(separator)) {
            // ends with separator
            newDirectory = newDirectory.trim();
        } else {
            // append a separator
            newDirectory = String.format("%s%s", newDirectory.trim(), separator);
        }

        if (! newDirectory.startsWith("/")){
            // not an absolute path
            return String.format("%s/%s",new File(".").toPath().toAbsolutePath().toString(), newDirectory);
        }
        return newDirectory;
    }
}