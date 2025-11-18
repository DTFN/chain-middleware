/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"chainmaker.org/chainmaker/pb-go/v2/common"

	"go.uber.org/zap"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/logger"
)

// ChainClientMock 链客户端结构体Mock
type ChainClientMock struct {
	client map[string]*sdk.ChainClient
	log    *zap.SugaredLogger
}

// InitChainClientMock 初始化链客户端
//  @return error
func InitChainClientMock() error {
	chainmakerClient := &ChainClientMock{
		client: make(map[string]*sdk.ChainClient),
		log:    logger.GetLogger(logger.ModuleChainClient),
	}
	ChainClientV1 = chainmakerClient
	return nil
}

// InvokeContract invoke合约
//  @receiver c
//  @param chainId
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return *common.TransactionInfo
//  @return error
func (c *ChainClientMock) InvokeContract(chainId, contractName,
	method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, *common.TransactionInfo, error) {
	return []byte("123"), &common.TransactionInfo{
		Transaction: &common.Transaction{
			Payload: &common.Payload{
				TxId: "1234567890",
			},
		},
		BlockHeight: 10,
	}, nil
}

// GetTxProve 获取交易证明
//  @receiver c
//  @param blockHeight
//  @param chainId
//  @param tx
//  @return string
func (c *ChainClientMock) GetTxProve(blockHeight uint64, chainId string, tx *common.TransactionInfo) string {
	return "{}"
}

// CheckChain 检查链的连通性
//  @receiver c
//  @return bool
func (c *ChainClientMock) CheckChain() bool {
	return true
}
