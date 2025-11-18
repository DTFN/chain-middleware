/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
)

// ChainClientMock 链客户端结构体Mock
type ChainClientMock struct {
	//client *fabsdk.FabricSDK
	log *zap.SugaredLogger
}

// InitChainClientMock 初始化链客户端
//  @return error
func InitChainClientMock() error {
	chainmakerClient := &ChainClientMock{
		log: logger.GetLogger(logger.ModuleChainClient),
	}
	ChainClientV1 = chainmakerClient
	return nil
}

// InvokeContract invoke合约
//  @receiver c
//  @param chainId
//  @param contractName
//  @param method
//  @param args
//  @param needTx
//  @return string
//  @return string
//  @return *peer.ProcessedTransaction
//  @return error
func (c *ChainClientMock) InvokeContract(
	chainId, contractName, method string,
	args [][]byte, needTx bool) (string, string, *peer.ProcessedTransaction, error) {
	return "", "123", &peer.ProcessedTransaction{
		TransactionEnvelope: &common.Envelope{
			Payload:   []byte("123"),
			Signature: []byte("321"),
		},
		ValidationCode: 10,
	}, nil
}

// GetTxProve 获取交易证明
//  @receiver c
//  @param tx
//  @param txId
//  @param chainId
//  @return string
func (c *ChainClientMock) GetTxProve(tx *peer.ProcessedTransaction, txId, chainId string) string {
	return "{}"
}

// TxProve 交易认证
//  @receiver c
//  @param txProve
//  @return bool
func (c *ChainClientMock) TxProve(txProve string) bool {
	return txProve != ""
}

// CheckChain 检查链的连通性
//  @receiver c
//  @return bool
func (c *ChainClientMock) CheckChain() bool {
	return true
}
