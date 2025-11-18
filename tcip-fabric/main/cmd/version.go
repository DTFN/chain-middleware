/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

// VersionCMD 版本命令
//  @return *cobra.Command
func VersionCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show tcip-fabric version",
		Long:  "Show ChainMaker version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			PrintVersion()
			return nil
		},
	}
}

func logo() string {
	fig := figure.NewFigure("Tcip-Fabric", "slant", true)
	s := fig.String()
	fragment := "================================================================================================="
	versionInfo := fmt.Sprintf("Tcip-Fabric Version: %s\n", conf.CurrentVersion)
	if conf.CurrentBranch != "" {
		conf.CurrentBranch = fmt.Sprintf("Tcip-Fabric Branch: %s\n", conf.CurrentBranch)
	}
	if conf.CurrentCommit != "" {
		conf.CurrentCommit = fmt.Sprintf("Tcip-Fabric Commit: %s\n", conf.CurrentCommit)
	}
	if conf.BuildTime != "" {
		conf.BuildTime = fmt.Sprintf("Tcip-Fabric Build Time: %s\n", conf.BuildTime)
	}
	return fmt.Sprintf("\n%s\n%s%s\n%s\n%s\n%s\n%s\n", fragment, s, fragment, versionInfo,
		conf.CurrentBranch, conf.CurrentCommit, conf.BuildTime)
}

// PrintVersion 打印版本信息
func PrintVersion() {
	fmt.Println(logo())
	fmt.Println()
}
