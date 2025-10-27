-- 部署chain，添加节点
drop table `deploy_result` if exists;
create table `deploy_result` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `chain_uid` VARCHAR(100) NOT NULL comment '链hash',
    `request_id` VARCHAR(100) NOT NULL comment '请求Id,部署链时为chain_uid',
    `channel_id` VARCHAR(100) comment '通道id',
    `stage` INTEGER DEFAULT 0 comment '运行阶段',
    `error` TEXT comment '失败原因'
);
CREATE INDEX deploy_result_chain_uid_request_id ON deploy_result(`chain_uid`, `request_id`);

drop table `deploy_result_node` if exists;
create table `deploy_result_node` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `result_id` BIGINT comment 'deploy_result.id',
    `node_full_name` VARCHAR(100) comment '节点名',
    `node_type` INTEGER comment '节点类型，1-orderer，2-peer',
    `network_name` VARCHAR(100) comment '网络名称',
    `node_port` INTEGER comment '节点端口',
    `node_ip` VARCHAR(20) comment '节点Ip',
    `org_full_name` VARCHAR(100) comment '组织名',
    `msp_id` VARCHAR(100) comment 'MSPID',
    `chain_code_port` INTEGER comment '链码端口',
    `operations_port` INTEGER comment '操作端口'
);
CREATE INDEX deploy_result_id_ip_port ON deploy_result_node(`result_id`, `node_ip`, `node_port`);

-- 部署chain，添加节点
drop table `chain_info` if exists;
create table `chain_info` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `chain_uid` VARCHAR(100) NOT NULL comment '链uid',
    `ca_host` VARCHAR(20) comment 'ca主机ip',
    `ca_name` VARCHAR(60) comment 'ca名称',
    `ca_port` INTEGER NOT NULL comment 'CA端口',
    `modify_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);

-- 部署链码
drop table `chain_code` if exists;
create table `chain_code` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `channel_id` VARCHAR(100) comment '通道id',
    `chain_code_name` VARCHAR(100) NOT NULL comment '链码名称',
    `lang` VARCHAR(10) comment '链码语言'
);
CREATE INDEX chain_code_channel ON chain_code(`channel_id`, `chain_code_name`);

-- 部署链码-节点
drop table `chain_code_peer` if exists;
create table `chain_code_peer` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `chain_code_id` BIGINT comment '链码id',
    `peer_name` VARCHAR(100) NOT NULL comment '节点名称',
    `peer_url` VARCHAR(100) comment '节点url'
);
CREATE INDEX chain_code_peer_name_code_id ON chain_code_peer(`chain_code_id`, `peer_name`);

-- 组织信息
drop table `org_info` if exists;
create table `org_info` (
    `chain_uid` VARCHAR(100) NOT NULL comment '链hash',
    `org_name` VARCHAR(100) NOT NULL comment '组织agency',
    `admin_key` TEXT comment '管理员私钥，pem格式',
    `admin_crt` TEXT comment '管理员证书，pem格式',
    `ca_crt` TEXT comment 'ca证书，pem格式',
    `modify_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);
CREATE INDEX org_info_chain_uid_org_name ON org_info(`chain_uid`, `org_name`);


--- 节点监控
drop table `monitor_node` if exists;
create table `monitor_node` (
    `id` BIGINT PRIMARY KEY AUTO_INCREMENT comment 'id',
    `ip` VARCHAR(20) comment '节点Ip',
    `cpu_usage` DECIMAL(6, 3) comment '使用率',
    `mem_usage` BIGINT comment '使用的字节',
    `disk_usage` BIGINT comment '使用的字节',
    `net_in` BIGINT comment '字节/秒',
    `net_out` BIGINT comment '字节/秒',
    `create_time` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间'
);
CREATE INDEX monitor_node_ip_create_time ON monitor_node(`ip`, `create_time`);

--- 节点监控(不常变化的属性)
drop table `monitor_node_stable` if exists;
create table `monitor_node_stable` (
    `ip` VARCHAR(20) PRIMARY KEY comment '节点Ip',
    `cpu_number` INTEGER comment 'CPU数量',
    `max_memory` BIGINT comment '内存容量',
    `max_disk` BIGINT comment '硬盘容量',
    `modify_time` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间'
);
