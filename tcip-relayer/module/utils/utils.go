/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"encoding/json"
	"fmt"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
)

const (
	// SpvTxVerifyMethod spv证明方法
	SpvTxVerifyMethod = "verify_tx"

	// SyncBlockHeaderMethod 同步区块头方法
	SyncBlockHeaderMethod = "sync_block_header"

	// GetBlockHeaderMethod 获取区块头方法
	GetBlockHeaderMethod = "get_block_header"
)

var (
	// CrossChainManager 跨链管理合约名称
	CrossChainManager = syscontract.SystemContract_RELAY_CROSS.String()

	// SaveGateway 保存网关方法名
	SaveGateway = syscontract.RelayCrossFunction_SAVE_GATEWAY.String()
	// UpadteGateway 更新网关方法名
	UpadteGateway = syscontract.RelayCrossFunction_UPDATE_GATEWAY.String()
	// GetGatewayNum 获取网关个数方法名
	GetGatewayNum = syscontract.RelayCrossFunction_GET_GATEWAY_NUM.String()
	// GetGateway 获取网关方法名
	GetGateway = syscontract.RelayCrossFunction_GET_GATEWAY.String()
	// GetGatewayByRange 批量获取网关方法名
	GetGatewayByRange = syscontract.RelayCrossFunction_GET_GATEWAY_BY_RANGE.String()
	// SaveCrossChainInfo 保存跨链交易方法名
	SaveCrossChainInfo = syscontract.RelayCrossFunction_SAVE_CROSS_CHAIN_INFO.String()
	// UpdateCrossChainTry 更新cross chain try方法名
	UpdateCrossChainTry = syscontract.RelayCrossFunction_UPDATE_CROSS_CHAIN_TRY.String()
	// UpdateCrossChainResult 更新跨链结果方法名
	UpdateCrossChainResult = syscontract.RelayCrossFunction_UPDATE_CROSS_CHAIN_RESULT.String()
	// UpdateCrossChainConfirm 更新目标网关confirm方法名
	UpdateCrossChainConfirm = syscontract.RelayCrossFunction_UPDATE_CROSS_CHAIN_CONFIRM.String()
	// UpdateSrcGatewayConfrim 更新源网关confirm方法名
	UpdateSrcGatewayConfrim = syscontract.RelayCrossFunction_UPDATE_SRC_GATEWAY_CONFIRM.String()
	// GetCrossChainNum 获取跨链交易条数方法名
	GetCrossChainNum = syscontract.RelayCrossFunction_GET_CROSS_CHAIN_NUM.String()
	// GetCrossChainInfo 获取跨链交易方法名
	GetCrossChainInfo = syscontract.RelayCrossFunction_GET_CROSS_CHAIN_INFO.String()
	// GetCrossChainInfoByRange 批量获取跨链交易方法名
	GetCrossChainInfoByRange = syscontract.RelayCrossFunction_GET_CROSS_CHAIN_INFO_BY_RANGE.String()
	// GetNotEndCrossChainIdList 获取未完成的跨链交易id方法名
	GetNotEndCrossChainIdList = syscontract.RelayCrossFunction_GET_NOT_END_CROSS_CHIAN_ID_LIST.String()
	// DIDManager did管理合约
	DIDManager = "DID_MANAGER"
)

var (
	// CrossChainTryChan 通知跨链交易模块调用crossChainTry
	CrossChainTryChan chan string
	// CrossChainResultChan 通知跨链交易模块调用跨链结果更新
	CrossChainResultChan chan string
	// CrossChainConfirmChan 通知跨链交易模块调用crossChainConfirm
	CrossChainConfirmChan chan string
	// CrossChainSrcGatewayConfirmChan 通知跨链交易模块调用src crossChainConfirm
	CrossChainSrcGatewayConfirmChan chan string
	// DIDManagerUpdateDIDChan 通知更新DID文档
	DIDManagerUpdateDIDChan chan string
)

// GetSpvContractName 获取spv交易证明的参数
//  @param gatewayId
//  @param chainRid
//  @return string
func GetSpvContractName(gatewayId, chainRid string) string {
	return "spv" + gatewayId + chainRid
}

// GetBlockHeaderParam 获取区块头的参数
//  @param blockHeight
//  @return string
func GetBlockHeaderParam(blockHeight int64) string {
	blockHeightByte := []byte(fmt.Sprintf("%d", blockHeight))
	param := make(map[string][]byte)
	param["block_height"] = blockHeightByte
	paramByte, _ := json.Marshal(param)
	return string(paramByte)
}

// GetSyncBlockHeaderParameter 获取同步区块头的参数
//  @param blockHeight
//  @param blockHeader
//  @return string
func GetSyncBlockHeaderParameter(blockHeight uint64, blockHeader []byte) string {
	res := make(map[string][]byte)
	res["block_height"] = []byte(fmt.Sprintf("%d", blockHeight))
	res["block_header"] = blockHeader
	resJson, _ := json.Marshal(res)
	return string(resJson)
}

// UnsupportVersion 不支持的版本打印
//  @param version
//  @return string
func UnsupportVersion(version common.Version) string {
	return fmt.Sprintf("Unsupported version: %d", version)
}
