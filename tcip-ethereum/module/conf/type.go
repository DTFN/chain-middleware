/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/logger"
)

// LocalConfig 本地配置信息
type LocalConfig struct {
	BaseConfig *BaseConfig               `mapstructure:"base"`  // 跨链网关基本配置
	RpcConfig  *RpcConfig                `mapstructure:"rpc"`   // Web监听配置
	Relay      *Relay                    `mapstructure:"relay"` // 中继网关信息
	DbPath     string                    `mapstructure:"db_path"`
	LogConfig  []*logger.LogModuleConfig `mapstructure:"log"` // 日志配置
}

// BaseConfig 跨链网关基本配置
type BaseConfig struct {
	GatewayID   string `mapstructure:"gateway_id"`   // 跨链网关ID，这里需要等待注册以后才能填写
	GatewayName string `mapstructure:"gateway_name"` // 跨链网关名称
	Address     string `mapstructure:"address"`      // 跨链网关地址
	ServerName  string `mapstructure:"server_name"`  // 跨链网关的server name
	Tlsca       string `mapstructure:"tls_ca"`       // 跨链网关的ca证书路径
	ClientCert  string `mapstructure:"client_cert"`  // 跨链网关的客户端证书路径
	ClientKey   string `mapstructure:"client_key"`   // 跨链网关的客户私钥在跨链网关服务器上的地址
	// 跨链网关允许向哪些网关发起请求，其他跨链网关的id，如果是空则不加限制
	ToGatewayList []string `mapstructure:"to_gateway_list"`
	// 跨链网关地址允许被哪些网关发起请求，其他跨链网关的id，如果是空则不加限制
	FromGatewayList []string `mapstructure:"from_gateway_list"`
	// 交易的验证方式，支持spv验证和rpc验证两种方式
	TxVerifyType string `mapstructure:"tx_verify_type"`
	// 如果是rpc验证方式，需要添加验证接口的信息,接口只支持restful的方式调用，具体调用方式请参考文档
	TxVerifyInterface *TxVerifyInterface `mapstructure:"tx_verify_interface"`
	CallType          string             `mapstructure:"call_type"`       // 跨链网关调用方式，现在只支持restful和grpc
	DefaultTimeout    uint32             `mapstructure:"default_timeout"` // 默认的全局超时时间
}

// RpcConfig rpc配置
type RpcConfig struct {
	Port           int          `mapstructure:"port"`      // 服务监听的端口号
	TLSConfig      TlsConfig    `mapstructure:"tls"`       // tls相关配置
	BlackList      []string     `mapstructure:"blacklist"` // 黑名单
	RestfulConfig  RstfulConfig `mapstructure:"restful"`   // resultful api 网关
	MaxSendMsgSize int          `mapstructure:"max_send_msg_size"`
	MaxRecvMsgSize int          `mapstructure:"max_recv_msg_size"`
}

// TlsConfig tls配置
type TlsConfig struct {
	CaFile     string `mapstructure:"ca_file"`
	KeyFile    string `mapstructure:"key_file"`
	CertFile   string `mapstructure:"cert_file"`
	ServerName string `mapstructure:"server_name"`
}

// RstfulConfig rest服务配置
type RstfulConfig struct {
	MaxRespBodySize int `mapstructure:"max_resp_body_size"`
}

// Relay 中继网关配置
type Relay struct {
	Address    string `mapstructure:"address"`     // 中继网关地址
	ServerName string `mapstructure:"server_name"` // 中继网关的server name
	Tlsca      string `mapstructure:"tls_ca"`      // 中继网关的ca证书路径
	ClientCert string `mapstructure:"client_cert"` // 中继网关的客户端证书路径
	ClientKey  string `mapstructure:"client_key"`  // 中继网关的客户端私钥
	CallType   string `mapstructure:"call_type"`   // 调用类型
}

// TxVerifyInterface 交易验证接口配置
type TxVerifyInterface struct {
	Address    string `mapstructure:"address"`     // 验证接口地址
	TlsEnable  bool   `mapstructure:"tls_enable"`  // 是否开启tls
	Tlsca      string `mapstructure:"tls_ca"`      // tls的ca证书路径
	ClientCert string `mapstructure:"client_cert"` // 客户端证书路径
	HostName   string `mapstructure:"host_name"`   // 服务名
	//ClientKey  string `mapstructure:"client_key"`  // 客户端私钥在中继网关服务器上的路径
}
