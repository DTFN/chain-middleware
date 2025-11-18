/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/request"

	"github.com/spf13/pflag"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/logger"
	"github.com/spf13/cobra"
)

// SpvContractCMD spv合约管理命令
//  @return *cobra.Command
func SpvContractCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "spv",
		Short: "Install or update spv contract",
		Long:  "Install or update spv contract",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			SpvContract()
			fmt.Println("tcip-chainmaker exit")
			return nil
		},
	}
	SpvContractFlags(startCmd, []string{
		flagNameOfConfigFilepath,
		flagNameOfObjectPathPath,
		spvContractVersion,
		spvContractPath,
		spvContractRuntimeType,
		spvContractParam,
		spvContractChainId,
		spvContractOption,
	})
	return startCmd
}

// SpvContract spv合约管理
func SpvContract() {
	cliLog = logger.GetLogger(logger.ModuleRegister)
	config, _ := json.Marshal(conf.Config)
	cliLog.Debug(string(config))

	_ = request.InitRequestManager()
	switch SpvOption {
	case "install":
		err := request.RequestV1.InitSpvContracta(SpvVersion, SpvPath, SpvRuntimeType, SpvParam, SpvChainRid)
		if err != nil {
			fmt.Println("spv error: ", err)
			return
		}
	case "update":
		err := request.RequestV1.UpdateSpvContract(SpvVersion, SpvPath, SpvRuntimeType, SpvParam, SpvChainRid)
		if err != nil {
			fmt.Println("spv error: ", err)
			return
		}
	default:
		fmt.Printf("unsupport option: %s, need [install/update]", SpvOption)
	}
	fmt.Println("success")
}

// SpvContractFlagSet spv合约管理
//  @return *pflag.FlagSet
func SpvContractFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilePath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath,
		conf.ConfigFilePath, "specify config file path, if not set, default use ./tcip_chainmaker.yml")
	flags.StringVarP(&SpvVersion, spvContractVersion, spvContractVShortersion,
		SpvVersion, "spv contract version, default is 1.0")
	flags.StringVarP(&SpvPath, spvContractPath, spvContractShortPath,
		SpvPath, "spv contract path, default is"+
			" ./spv0chain1.7z, docker-go contract name is spv+gatewayId+chainId")
	flags.StringVarP(&SpvRuntimeType, spvContractRuntimeType, spvContractShortRuntimeType,
		SpvRuntimeType, "spv contract runtime type, default is DOCKER_GO, [DOCKER_GO/EVM/GASM/WXVM/WASMER], "+
			"docker-go contract name is spv+gatewayId+chainId")
	flags.StringVarP(&SpvParam, spvContractParam, spvContractShortParam,
		SpvParam, "spv contract parameter, json string, default is {}")
	flags.StringVarP(&SpvChainRid, spvContractChainId, spvContractShortChainId,
		SpvChainRid, "spv contract chain resource id, default is chain1")
	flags.StringVarP(&SpvOption, spvContractOption, spvContractShortOption,
		SpvOption, "spv contract option, default is install, [install/update]")
	return flags
}

// SpvContractFlags spv合约管理
//  @param cmd
//  @param flagNames
func SpvContractFlags(cmd *cobra.Command, flagNames []string) {
	flags := SpvContractFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
