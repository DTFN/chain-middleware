/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
)

// LocalConfig 本地配置信息
type LocalConfig struct {
	BaseConfig *BaseConfig               `mapstructure:"base"`        // 中继网关基本配置，目前没什么用，留一个扩展
	RpcConfig  *RpcConfig                `mapstructure:"rpc"`         // Web监听配置
	RelayChain *RelayChain               `mapstructure:"relay_chain"` // Web监听配置
	LogConfig  []*logger.LogModuleConfig `mapstructure:"log"`         // 日志配置
}

// BaseConfig 中继网关基本配置
type BaseConfig struct {
	GatewayID      string `mapstructure:"gateway_id"`      // 中继网关ID
	GatewayName    string `mapstructure:"gateway_name"`    // 中继网关名称
	DefaultTimeout int64  `mapstructure:"default_timeout"` // 默认的全局超时时间
}

// RpcConfig rpc配置信息
type RpcConfig struct {
	Port           int           `mapstructure:"port"`      // 服务监听的端口号
	TLSConfig      tlsConfig     `mapstructure:"tls"`       // tls相关配置
	BlackList      []string      `mapstructure:"blacklist"` // 黑名单
	RestfulConfig  restfulConfig `mapstructure:"restful"`   // resultful api 网关
	MaxSendMsgSize int           `mapstructure:"max_send_msg_size"`
	MaxRecvMsgSize int           `mapstructure:"max_recv_msg_size"`
	ClientKey      string        `mapstructure:"client_key"`
	ClientCert     string        `mapstructure:"client_cert"`
}

type tlsConfig struct {
	Mode       string `mapstructure:"mode"`
	CaFile     string `mapstructure:"ca_file"`
	KeyFile    string `mapstructure:"key_file"`
	CertFile   string `mapstructure:"cert_file"`
	ServerName string `mapstructure:"server_name"`
}

type restfulConfig struct {
	MaxRespBodySize int `mapstructure:"max_resp_body_size"`
}

// RelayChain 中继链配置
type RelayChain struct {
	ChainmakerSdkConfigPath string  `mapstructure:"chainmaker_sdk_config_path"`
	Users                   []*User `mapstructure:"users"`
}

// User 中继链用户信息
type User struct {
	SignKeyPath string `mapstructure:"sign_key_path"`
	SignCrtPath string `mapstructure:"sign_crt_path"`
	OrgId       string `mapstructure:"org_id"`
}
