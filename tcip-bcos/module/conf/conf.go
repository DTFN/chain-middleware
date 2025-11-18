/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
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
	ConfigFilePath = "./tcip_bcos.yml"
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
//  @param cmd
//  @return error
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
//  @param cmd
//  @return *LocalConfig
//  @return error
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
//  @param ymlFile
//  @return string
func GetAbsPath(ymlFile string) string {
	if !filepath.IsAbs(ymlFile) {
		ymlFile, _ = filepath.Abs(ymlFile)
	}
	return ymlFile
}
