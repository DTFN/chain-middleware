/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"fmt"
	"path/filepath"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/logger"
	proto "github.com/gogo/protobuf/proto"
)

var (
	// CurrentVersion 当前版本
	CurrentVersion = "1.0.0"
	// CurrentBranch 当前版本
	CurrentBranch = ""
	// CurrentCommit 当前版本
	CurrentCommit = ""
	// BuildTime 编译时间
	BuildTime = ""
	// ConfigFilePath 默认配置位置
	ConfigFilePath = "./tcip_ethereum.yml"
	// BaseConf 基础配置
	BaseConf = &BaseConfig{}
	// Config 全剧配置
	Config = &LocalConfig{}
)

const (
	// GrpcCallType grpc
	GrpcCallType = "grpc"
	// RestCallType restful
	RestCallType = "restful"

	// RpcTxVerify rpc
	RpcTxVerify = "rpc"
	// SpvTxVerify spv
	SpvTxVerify = "spv"
	// NotNeedTxVerify not need
	NotNeedTxVerify = "notneed"
)

// InitLocalConfig init local config
//
//	@param cmd
//	@return error
func InitLocalConfig(cmd *cobra.Command) error {
	// 1. init config
	config, err := initLocal(cmd)
	if err != nil {
		return err
	}
	// 处理 log config
	logModuleConfigs := config.LogConfig
	for i := 0; i < len(logModuleConfigs); i++ {
		logModuleConfig := logModuleConfigs[i]
		logModuleConfig.FilePath = GetAbsPath(logModuleConfig.FilePath)
	}
	// 2. set log config
	logger.InitLogConfig(config.LogConfig)
	// 3. set global config and export
	Config = config
	BaseConf = config.BaseConfig
	logger.GetLogger(logger.ModuleDefault).Info(fmt.Sprintf("Local config inited, GatewayID=[%s], Name=[%s]",
		BaseConf.GatewayID, BaseConf.GatewayName))
	return nil
}

// initLocal 初始化本地配置
//
//	@param cmd
//	@return *LocalConfig
//	@return error
func initLocal(cmd *cobra.Command) (*LocalConfig, error) {
	cmViper := viper.New()

	// 1. load the path of the config files
	ymlFile := ConfigFilePath
	ymlFile = GetAbsPath(ymlFile)
	ConfigFilePath = ymlFile

	// 2. load the config file
	cmViper.SetConfigFile(ymlFile)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, command := range cmd.Commands() {
		err := cmViper.BindPFlags(command.PersistentFlags())
		if err != nil {
			return nil, err
		}
	}

	// 3. create new CMConfig instance
	config := &LocalConfig{}
	if err := cmViper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetAbsPath 获取绝对路径
//
//	@param ymlFile
//	@return string
func GetAbsPath(ymlFile string) string {
	if !filepath.IsAbs(ymlFile) {
		ymlFile, _ = filepath.Abs(ymlFile)
	}
	return ymlFile
}

// eth配置信息
type EthConfig struct {
	// 链资源Id
	ChainRid string `protobuf:"bytes,1,opt,name=chain_rid,json=chainRid,proto3" json:"chain_rid,omitempty"`
	// chain id
	ChainId int64 `protobuf:"varint,2,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// 状态，由网关设置用户不需要关系
	State bool `protobuf:"varint,3,opt,name=state,proto3" json:"state,omitempty"`
	// 状态描述，由网关设置，状态为true，此处是成功信息，状态为false，此处有错误信息
	StateMessage string `protobuf:"bytes,4,opt,name=state_message,json=stateMessage,proto3" json:"state_message,omitempty"`
}

func (m *EthConfig) Reset()         { *m = EthConfig{} }
func (m *EthConfig) String() string { return proto.CompactTextString(m) }
func (*EthConfig) ProtoMessage()    {}

func (m *EthConfig) GetChainId() int64 {
	if m != nil {
		return m.ChainId
	}
	return 0
}

func (m *EthConfig) GetState() bool {
	if m != nil {
		return m.State
	}
	return false
}

func (m *EthConfig) GetStateMessage() string {
	if m != nil {
		return m.StateMessage
	}
	return ""
}

type ChainIdentityRequest struct {
	// 请求消息版本
	Version common.Version `protobuf:"varint,1,opt,name=version,proto3,enum=common.Version" json:"version,omitempty"`
	// 操作
	Operate common.Operate `protobuf:"varint,2,opt,name=operate,proto3,enum=common.Operate" json:"operate,omitempty"`
	// chainmaker配置，chainmaker网关专属
	ChainmakerConfig *common.ChainmakerConfig `protobuf:"bytes,4,opt,name=chainmaker_config,json=chainmakerConfig,proto3" json:"chainmaker_config,omitempty"`
	// fabric配置，fabric网关专属
	FabricConfig *common.FabricConfig `protobuf:"bytes,5,opt,name=fabric_config,json=fabricConfig,proto3" json:"fabric_config,omitempty"`
	// bcos配置，bcos网关专属
	BcosConfig *common.BcosConfig `protobuf:"bytes,6,opt,name=bcos_config,json=bcosConfig,proto3" json:"bcos_config,omitempty"`
	// ethereum配置, ethereum网关专属
	EthConfig *EthConfig `protobuf:"bytes,7,opt,name=eth_config,json=ethConfig,proto3" json:"eth_config,omitempty"`
	// lingshu配置, lingshu网关专属
	LsConfig *EthConfig `protobuf:"bytes,7,opt,name=ls_config,json=lsConfig,proto3" json:"ls_config,omitempty"`
}
