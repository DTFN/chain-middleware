/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/main/cmd"
	"github.com/spf13/cobra"
)

// ./tcip-lingshu start -c tcip-lingshu.yml
func main() {
	mainCmd := &cobra.Command{Use: "tcip-lingshu"}
	mainCmd.AddCommand(cmd.StartCMD())
	mainCmd.AddCommand(cmd.RegisterCMD())
	mainCmd.AddCommand(cmd.UpdateCMD())
	mainCmd.AddCommand(cmd.VersionCMD())
	mainCmd.AddCommand(cmd.SpvContractCMD())

	err := mainCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
