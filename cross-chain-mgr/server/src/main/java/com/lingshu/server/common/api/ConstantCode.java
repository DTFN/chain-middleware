package com.lingshu.server.common.api;

/**
 * @author lin
 * @since 2025-09-16
 */
public class ConstantCode {
    /* return success */
    public static final RetCode SUCCESS = RetCode.mark(0, "success");
    public static final RetCode FAILED = RetCode.mark(500, "failed");

    public static final RetCode CHAIN_NO_FOUND = RetCode.mark(202599, "chain no found");
    public static final RetCode CHAIN_CODE_VERSION_EXISTS = RetCode.mark(202600, "chain code version exists");
    public static final RetCode CHAIN_CODE_NOT_EXISTS = RetCode.mark(202601, "chain code not exists");
    public static final RetCode CHAIN_CODE_DEPLOY_FAILED = RetCode.mark(202602, "chain code deploy failed");
    public static final RetCode CHANNEL_NO_FOUND = RetCode.mark(202603, "channel no found");
    public static final RetCode INVALID_LANG_TYPE = RetCode.mark(202604, "invalid lang type");
    public static final RetCode CHAIN_CODE_INVOKE_FAILED = RetCode.mark(202605, "chain code invoke failed");
    public static final RetCode DUPLICATED_CHANNEL_NAME = RetCode.mark(202606, "channel name duplicated in the same chain");
    public static final RetCode NO_USEFUL_PEER = RetCode.mark(202607, "no useful peer");
    public static final RetCode NO_USEFUL_ORDERER = RetCode.mark(202608, "no useful orderer");
    public static final RetCode ADD_APP_CHANNEL_FAILED = RetCode.mark(202609, "failed to add app channel");
    public static final RetCode ADD_PEERS_INTO_APP_CHANNEL_FAILED = RetCode.mark(202610, "failed to add peers into app channel");
    public static final RetCode CHAIN_CODE_FILE_NOT_EXISTS = RetCode.mark(202611, "chain code file not exists");
    public static final RetCode REQUEST_FABRIC_AGENT_FAIL = RetCode.mark(202612, "request fabric agent failed.");
    public static final RetCode EXPLORER_NOT_ACCESS = RetCode.mark(202613, "explorer not access.");
    public static final RetCode AGENT_IP_ERROR = RetCode.mark(202614, "agent ip error.");
    public static final RetCode AGENT_PORT_ERROR = RetCode.mark(202615, "agent port error.");
    public static final RetCode HOST_FREE_PORT_ERROR = RetCode.mark(202616, "get free host fail.");
    public static final RetCode DUPLICATED_CHAIN_NAME = RetCode.mark(202617, "chain name duplicated");
    public static final RetCode DEPLOY_PEER_FAIL = RetCode.mark(202618, "deploy peer fail");
    public static final RetCode DEPLOY_ORDERER_FAIL = RetCode.mark(202619, "deploy orderer fail");
    public static final RetCode INIT_CHAIN_FAIL = RetCode.mark(202620, "init chain fail");
    public static final RetCode NODE_NO_FOUND = RetCode.mark(202621, "node no found");
    public static final RetCode CHAIN_CODE_EXISTS = RetCode.mark(202622, "chain code already exists");
    public static final RetCode FABRIC_EXPLORER_CONN_REFUSED = RetCode.mark(202623, "connect to fabric explorer failed");
    public static final RetCode AGENT_CONN_REFUSED = RetCode.mark(202624, "connect to fabric agent failed");
    public static final RetCode FOUND_NO_NETWORK = RetCode.mark(202625, "found no networks");
    public static final RetCode INVALID_CHAIN_CODE = RetCode.mark(202630, "chain code invalid");
    public static final RetCode CHAIN_CODE_DEPLOYING = RetCode.mark(202631, "chain code is deploying");
    public static final RetCode FABRIC_NETWORK_ERROR = RetCode.mark(202632, "fabric network error");
    public static final RetCode MAX_CHANNEL_COUNT_ERROR = RetCode.mark(202641, "supports creating up to 10 channels");
    public static final RetCode CHAIN_MAKER_EXPLORER_CONN_REFUSED = RetCode.mark(202642, "connect to ChainMaker explorer failed");
}
