package com.lingshu.serverbcos.dto;

import lombok.Data;

@Data
public class ResourceDomainInfo {
        private String resourceName;
        private String gatewayId;
        private String chainRid;
        private String contractType;
        private String contractAddress;
}