/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"os"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/conf"
	"github.com/spf13/cobra"
)

const (
	flagNameOfConfigFilepath          = "conf-file"
	flagNameShortHandOfConfigFilepath = "c"
	flagNameOfObjectPathPath          = "object-path"
	flagNameShortHandOfObjectPathPath = "o"

	spvContractVersion          = "version"
	spvContractVShortersion     = "v"
	spvContractPath             = "path"
	spvContractShortPath        = "p"
	spvContractRuntimeType      = "runtime-type"
	spvContractShortRuntimeType = "r"
	spvContractParam            = "parameters"
	spvContractShortParam       = "P"
	spvContractChainId          = "chain-id"
	spvContractShortChainId     = "C"
	spvContractOption           = "option"
	spvContractShortOption      = "O"
)

var (
	// RegisterInfoPath 默认注册信息地址
	RegisterInfoPath = "register.json"
	// UpdateInfoPath 默认更新信息地址
	UpdateInfoPath = "update.json"
	// SpvVersion 默认合约版本号
	SpvVersion = "1.0"
	// SpvPath 默认合约地址
	SpvPath = "./spv1chain1.7z"
	// SpvRuntimeType 默认合约运行时
	SpvRuntimeType = "DOCKER_GO"
	// SpvParam 默认合约参数
	SpvParam = "{}"
	// SpvChainRid 默认chain id
	SpvChainRid = "chain1"
	// SpvOption 默认操作
	SpvOption = "install"
)

func initLocalConfig(cmd *cobra.Command) {
	if err := conf.InitLocalConfig(cmd); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
