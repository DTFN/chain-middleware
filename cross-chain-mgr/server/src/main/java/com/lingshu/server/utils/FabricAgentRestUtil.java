package com.lingshu.server.utils;

import cn.hutool.http.HttpUtil;
import cn.hutool.json.JSONObject;
import cn.hutool.json.JSONUtil;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.lingshu.server.common.api.R;
import com.lingshu.server.common.exception.FabricApiException;
import com.lingshu.server.common.api.ConstantCode;
import com.lingshu.server.common.api.OpenAPIResp;
import lombok.extern.slf4j.Slf4j;

import java.util.Map;
import java.util.Optional;

@Slf4j
public class FabricAgentRestUtil {

    public static final String URI_CHAIN_INFO = "/fabric-agent/deploy/chainInfo";
    public static final String URI_DELETE_CHAIN = "/fabric-agent/deploy/delete";
    public static final String URI_CHAIN_CODE_UPLOAD = "/fabric-agent/code/upload";
    public static final String URI_CHAIN_CODE_PACKAGE = "/fabric-agent/code/package";
    public static final String URI_CHAIN_CODE_INSTALL = "/fabric-agent/code/install";
    public static final String URI_CHAIN_CODE_APPROVE = "/fabric-agent/code/approve";
    public static final String URI_CHAIN_CODE_COMMIT = "/fabric-agent/code/commit";
    public static final String URI_CHAIN_CODE_INVOKE = "/fabric-agent/code/invoke";
    public static final String URI_CHAIN_CODE_QUERY = "/fabric-agent/code/query";
    public static final String URI_ADD_APP_CHANNEL = "/fabric-agent/deploy/addAppChannel";
    public static final String URI_ADD_PEERS_INTO_APP_CHANNEL = "/fabric-agent/deploy/addPeersIntoAppChannel";
    public static final String URI_ADD_NODE = "/fabric-agent/deploy/addNode";
    public static final String URI_AVAILABLE_PORT = "/fabric-agent/deploy/getAvailablePort?ip=%s&startPort=%s&number=%s";
    public static final String URI_DEPLOY = "/fabric-agent/deploy";
    public static final String URI_DEPLOY_RESULT = "/fabric-agent/deploy/getDeployResult?chainUid=%s";
    public static final String URI_ADD_NODE_RESULT = "/fabric-agent/deploy/getAddNodeResult?chainUid=%s&requestId=%s";
    public static final String URI_INIT_FABRIC_DEPLOY = "/fabric-agent/deploy/initHostAndDocker";
    public static final String URI_HEALTHZ = "/fabric-agent/nodeInfo/healthz/%s/%s";
    public static final String URI_NODE_BASE_INFO = "/fabric-agent/nodeInfo/basic?ip=%s";
    public static final String URI_NODE_OPERATION = "/fabric-agent/deploy/nodeOperation";
    public static final String URI_NODE_PERFORMANCE = "/fabric-agent/nodeInfo/performance/byTime";

    private static ObjectMapper om = new ObjectMapper();

    public static <T> T get(String fabricAgentUrl, String uri, Class<T> clazz) {
        String body = null;
        try {
            body = HttpUtil.createGet(fabricAgentUrl + uri).execute().body();
        } catch (Exception e) {
            log.error("failed to invoke fabric agent,url:{},uri:{}", fabricAgentUrl, uri, e);
            throw new FabricApiException(ConstantCode.AGENT_CONN_REFUSED);
        }

        return getResult(body, clazz);
    }

    public static <T> T post(String fabricAgentUrl, String uri, String json, Class<T> clazz) {
        String body = null;
        String requestUrl = fabricAgentUrl + uri;
        try {
            byte[] bts = HttpUtil.post(requestUrl, json).getBytes();
            body = new String(bts);
        } catch (Exception e) {
            log.error("request agent failed: url={}", requestUrl, e);
            throw new FabricApiException(ConstantCode.AGENT_CONN_REFUSED);
        }
        log.info("request agent success: url={}, body={}", requestUrl, body);
        return getResult(body, clazz);
    }

    public static void post(String fabricAgentUrl, String uri, String json) {
        String body = null;
        String requestUrl = fabricAgentUrl + uri;
        try {
            byte[] bts = HttpUtil.post(requestUrl, json).getBytes();
            body = new String(bts);
        } catch (Exception e) {
            log.error("request agent failed: url={}", requestUrl, e);
            throw new FabricApiException(ConstantCode.AGENT_CONN_REFUSED);
        }
        log.info("request agent success: url={}, body={}", requestUrl, body);
        parseResult(JSONUtil.parseObj(body));
    }

    public static <T> T post(String fabricAgentUrl, String uri, Map<String, Object> requestBody, Class<T> clazz) {
        String body = null;
        try {
            byte[] bts = HttpUtil.post(fabricAgentUrl + uri, requestBody).getBytes();
            body = new String(bts);
        } catch (Exception e) {
            log.error("failed to invoke fabric agent,url:{},uri:{}", fabricAgentUrl, uri, e);
            throw new FabricApiException(ConstantCode.AGENT_CONN_REFUSED);
        }
        return getResult(body, clazz);
    }

    private static <T> T getResult(String result, Class<T> clazz) {
        if (result == null) {
            throw new FabricApiException(ConstantCode.REQUEST_FABRIC_AGENT_FAIL);
        }

        R<T> tr = null;
        try {
            tr = om.readValue(result, om.getTypeFactory().constructParametricType(R.class, clazz));
        } catch (JsonProcessingException e) {
            log.error("failed to read value, result:{}", result);
            throw new FabricApiException(ConstantCode.REQUEST_FABRIC_AGENT_FAIL);
        }

        if (!Integer.valueOf(0).equals(tr.getCode())) {
            String msg = Optional.ofNullable(tr).map(R::getMsg).map(String::valueOf).orElse("");
            log.error("the response of fabric agent is not expected, result:{}", result);
            throw new FabricApiException(ConstantCode.REQUEST_FABRIC_AGENT_FAIL.attach(msg));
        }
        return tr.getData();
    }

    private static void parseResult(JSONObject result) {
        if (result == null) {
            throw new FabricApiException(ConstantCode.REQUEST_FABRIC_AGENT_FAIL);
        }
        Object code = result.get("code");
        if (!"0".equals(String.valueOf(code))) {
            throw new FabricApiException(ConstantCode.REQUEST_FABRIC_AGENT_FAIL.attach(String.valueOf(result.get("msg"))));
        }
    }
}
