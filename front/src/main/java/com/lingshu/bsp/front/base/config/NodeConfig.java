/*
 * Copyright 2014-2020  the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 */
package com.lingshu.bsp.front.base.config;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.lingshu.bsp.front.util.CommonUtils;
import com.lingshu.bsp.front.base.properties.Constants;
import lombok.Data;
import java.util.List;
import lombok.extern.slf4j.Slf4j;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.InitializingBean;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;

/**
 * initial info of node's configuration
 */
@Slf4j
@Data
@Configuration
public class NodeConfig implements InitializingBean {
    @JsonIgnore
    @Autowired
    private Constants constants;
    @Deprecated
    private String orgName = "lingshu";
    private String p2pIp;
    private String listenIp;
    private String rpcPort;
    private String p2pPort;
    private String ledgerDataPath;
    private boolean enableStatistic;

    @Override
    public void afterPropertiesSet() throws Exception {
        if(StringUtils.isBlank(constants.getNodePath())) {
            return;
        }
        List<String> nodeInfos = CommonUtils.readFileToList(constants.getNodePath() + Constants.CONFIG_FILE);

        this.p2pIp = CommonUtils.getCurrentIp();
        if (nodeInfos == null || nodeInfos.isEmpty() ) {
          //  throw new FrontException(ConstantCode.GET_NODE_CONFIG_FAILE);
            log.warn("cannot read config.yaml");
            return;
        }

        for (String str : nodeInfos) {
            if (str.contains(" data_path")) {
                this.ledgerDataPath = str.substring(str.indexOf("\"") + 1, str.lastIndexOf("\""));
                continue;
            }
            if (str.contains(" listen_ip")) {
                this.listenIp = str.substring(str.indexOf("\"") + 1, str.lastIndexOf("\""));
                continue;
            }
            if (str.contains(" listen_port")) {
                this.p2pPort = str.substring(str.indexOf(": ") + 2);
                continue;
            }
            if (str.contains(" listen_port")) {
                this.rpcPort = str.substring(str.indexOf(": ") + 2);
                continue;
            }
//            if (str.contains(" enable_statistic")) {
//                this.enableStatistic = Boolean.parseBoolean(str.substring(str.indexOf("=") + 1));
//            }
        }
    }

    @Override
    public String toString() {
        final StringBuilder sb = new StringBuilder("{");
        sb.append("\"orgName\":\"").append(orgName).append('\"');
        sb.append(",\"p2pIp\":\"").append(p2pIp).append('\"');
        sb.append(",\"listenIp\":\"").append(listenIp).append('\"');
        sb.append(",\"rpcPort\":\"").append(rpcPort).append('\"');
        sb.append(",\"p2pPort\":\"").append(p2pPort).append('\"');
        sb.append(",\"ledgerDataPath\":\"").append(ledgerDataPath).append('\"');
        sb.append(",\"enableStatistic\":").append(enableStatistic);
        sb.append('}');
        return sb.toString();
    }
}
