/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"os"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagNameOfConfigFilepath          = "conf-file"
	flagNameShortHandOFConfigFilepath = "c"
)

func initLocalConfig(cmd *cobra.Command) {
	if err := conf.InitLocalConfig(cmd); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}

func initFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilePath, flagNameOfConfigFilepath, flagNameShortHandOFConfigFilepath,
		conf.ConfigFilePath, "specify config file path, if not set, default use ./tcip_relayer.yml")
	return flags
}

func attachFlags(cmd *cobra.Command, flagNames []string) {
	flags := initFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
