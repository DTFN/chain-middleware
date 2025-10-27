/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package relayChainManager

import (
	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

// RelayChainManager 中继链接口
type RelayChainManager interface {
	InitContract(
		contractName, version, byteCodeBase64 string,
		kvJsonStr string,
		withSyncResult bool,
		timeout int64,
		runtime common.ChainmakerRuntimeType) error
	UpdateContract(
		contractName, version, byteCodeBase64 string,
		kvJsonStr string,
		withSyncResult bool,
		timeout int64,
		runtime common.ChainmakerRuntimeType) error
	InvokeContract(contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error)
	QueryContract(contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error)
}
