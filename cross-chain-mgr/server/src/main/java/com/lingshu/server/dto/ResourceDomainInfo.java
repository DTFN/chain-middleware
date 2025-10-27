package com.lingshu.server.dto;

import lombok.Data;

/**
 * @author: derrick
 * @since: 2025-09-08
 */
@Data
public class ResourceDomainInfo {
        private String resourceName;
        private String gatewayId;
        private String chainRid;
        private String contractType;
        private String contractAddress;
}
