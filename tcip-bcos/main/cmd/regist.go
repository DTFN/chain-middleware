/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/request"

	"github.com/spf13/pflag"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"github.com/spf13/cobra"
)

// RegisterCMD 生成网关注册信息命令
//  @return *cobra.Command
func RegisterCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "register",
		Short: "Out put register tcip-bcos info files",
		Long:  "Out put register tcip-bcos info files",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			Register()
			fmt.Println("tcip-bcos exit")
			return nil
		},
	}
	registerAttachFlags(startCmd, []string{flagNameOfConfigFilepath, flagNameOfObjectPathPath})
	return startCmd
}

// Register 生成网关注册信息
func Register() {
	cliLog = logger.GetLogger(logger.ModuleRegister)
	config, _ := json.Marshal(conf.Config)
	cliLog.Debug(string(config))

	_ = request.InitRequestManager()
	err := request.RequestV1.GatewayRegister(RegisterInfoPath)
	if err != nil {
		fmt.Println("register error: ", err)
	}
	fmt.Println("success")
	//resStr, _ := json.Marshal(res)
	//if res.Code == common.Code_GATEWAY_SUCCESS {
	//	fmt.Println("register success: ", string(resStr))
	//} else {
	//	fmt.Println("register error: ", string(resStr))
	//}
}

func registerFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilePath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath,
		conf.ConfigFilePath, "specify config file path, if not set, default use ./tcip_bcos.yml")
	flags.StringVarP(&RegisterInfoPath, flagNameOfObjectPathPath, flagNameShortHandOfObjectPathPath,
		RegisterInfoPath, "object path, default is ./register.json")
	return flags
}

func registerAttachFlags(cmd *cobra.Command, flagNames []string) {
	flags := registerFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
