/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/accesscontrol"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/crosschaintx"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/prove"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"
)

// InitServer 初始化服务
//  @param errorC
func InitServer(errorC chan error) {
	// 初始化 relay chain manager
	if err := relay_chain_chainmaker.InitRelayChain(conf.Config.RelayChain); err != nil {
		errorC <- err
		return
	}
	// 初始化access contol
	if err := request.InitRequestManager(); err != nil {
		errorC <- err
		return
	}
	// 初始化gateway manager
	if err := gateway.InitGatewayManager(); err != nil {
		errorC <- err
		return
	}
	// 初始化prove manager
	if err := prove.InitProveManager(); err != nil {
		errorC <- err
		return
	}
	// 初始化access contol
	if err := accesscontrol.InitAccessControl(); err != nil {
		errorC <- err
		return
	}
	// 初始化cross chain manager
	if err := crosschaintx.InitCrossChainManager(); err != nil {
		errorC <- err
		return
	}
}
