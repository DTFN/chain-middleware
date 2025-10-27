/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/request"

	"github.com/spf13/pflag"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/logger"
	"github.com/spf13/cobra"
)

// UpdateCMD 生成网关注更新信息
//
//	@return *cobra.Command
func UpdateCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "update",
		Short: "Out put Update tcip-bcos info files",
		Long:  "Out put Update tcip-bcos info files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			Update()
			fmt.Println("tcip-bcos exit")
			return nil
		},
	}
	updateAttachFlags(startCmd, []string{flagNameOfConfigFilepath, flagNameOfObjectPathPath})
	return startCmd
}

// Update 生成网关注更新信息
func Update() {
	cliLog = logger.GetLogger(logger.ModuleRegister)
	config, _ := json.Marshal(conf.Config)
	cliLog.Debug(string(config))

	_ = request.InitRequestManager()
	err := request.RequestV1.GatewayUpdate(UpdateInfoPath)
	if err != nil {
		fmt.Println("update error: ", err)
	}
	fmt.Println("success")
	//resStr, _ := json.Marshal(res)
	//if res.Code == common.Code_GATEWAY_SUCCESS {
	//	fmt.Println("update success: ", string(resStr))
	//} else {
	//	fmt.Println("update error: ", string(resStr))
	//}
}

func updateFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilePath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath,
		conf.ConfigFilePath, "specify config file path, if not set, default use ./tcip_ethereum.yml")
	flags.StringVarP(&UpdateInfoPath, flagNameOfObjectPathPath, flagNameShortHandOfObjectPathPath,
		UpdateInfoPath, "object path, default is ./update.json")
	return flags
}

func updateAttachFlags(cmd *cobra.Command, flagNames []string) {
	flags := updateFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
