package com.lingshu.server.core.web3j.chainmaker;

import lombok.Data;
import org.chainmaker.sdk.config.ChainClientConfig;

/**
 * @author gongrui.wang
 * @since 2025/2/26
 */
@Data
public class ChainmakerProperties extends ChainClientConfig {
    /**
     * 是否使用国密算法
     */
    private Boolean isGm;

    /**
     * ca配置文件路径
     */
    private String cryptoConfigPath;
}
