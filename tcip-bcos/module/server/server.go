/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	chain_client "chainmaker.org/chainmaker/tcip-bcos/v2/module/chain-client"
	chain_config "chainmaker.org/chainmaker/tcip-bcos/v2/module/chain-config"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/db"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/event"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/request"
)

// InitServer 初始化服务
//  @param errorC
func InitServer(errorC chan error) {
	// 初始化db
	db.NewDbHandle()
	// 初始化链配置管理模块
	chain_config.NewChainConfig()
	// 初始化跨链触发器
	event.InitEventManager()
	// 初始化 relay chain manager
	if err := chain_client.InitChainClient(); err != nil {
		errorC <- err
		return
	}
	// 初始化 request manager
	if err := request.InitRequestManager(); err != nil {
		errorC <- err
		return
	}
}
