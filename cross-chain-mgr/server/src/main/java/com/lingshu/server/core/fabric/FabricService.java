package com.lingshu.server.core.fabric;

import cn.hutool.json.JSONUtil;
import com.lingshu.server.common.api.ConstantCode;
import com.lingshu.server.common.exception.FabricApiException;
import com.lingshu.server.dto.FabricChainCodeInvokeReq;
import com.lingshu.server.dto.FabricChainCodeInvokeResponse;
import com.lingshu.server.utils.FabricAgentRestUtil;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

/**
 * @author lin
 * @since 2025-09-16
 */
@Slf4j
@Service
public class FabricService {

    @Value("${fabric.agent}")
    private String agentUrl;

    public FabricChainCodeInvokeResponse invoke(FabricChainCodeInvokeReq dto) {
        try {
            log.info("invoke chain code start. agentUrl={} dto={}", agentUrl, dto);
            return FabricAgentRestUtil.post(agentUrl, FabricAgentRestUtil.URI_CHAIN_CODE_INVOKE, JSONUtil.toJsonStr(dto), FabricChainCodeInvokeResponse.class);
        } catch (Exception e) {
            log.error("invoke chain code failed. agentUrl={} dto={}", agentUrl, dto, e);
            throw new FabricApiException(ConstantCode.CHAIN_CODE_INVOKE_FAILED);
        }
    }
}
