/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	bcostypes "github.com/FISCO-BCOS/go-sdk/core/types"

	"go.uber.org/zap"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/logger"
)

// ChainClientMock 链客户端结构体Mock
type ChainClientMock struct {
	client map[string]*sdk.ChainClient
	log    *zap.SugaredLogger
}

// InitChainClientMock 初始化链客户端
//
//	@return error
func InitChainClientMock() error {
	chainmakerClient := &ChainClientMock{
		client: make(map[string]*sdk.ChainClient),
		log:    logger.GetLogger(logger.ModuleChainClient),
	}
	ChainClientV1 = chainmakerClient
	return nil
}

// InvokeContract 调用合约
//
//	@receiver c
//	@param hainId
//	@param contractName
//	@param method
//	@param abiStr
//	@param args
//	@param needTx
//	@param paramType
//	@return []string
//	@return *bcostypes.TransactionDetail
//	@return error
func (c *ChainClientMock) InvokeContract(hainId, contractName, method, abiStr string, args string,
	needTx bool, paramType []tcipcommon.EventDataType) ([]string, *bcostypes.TransactionDetail, error) {
	return []string{"123"}, &bcostypes.TransactionDetail{
		Hash:        "234567890",
		BlockNumber: "10",
	}, nil
}

func (c *ChainClientMock) DIDManagerUpdate(ledgerId int, contractAddress string, did string, didDoc string) (*bcostypes.Receipt, error) {
	return nil, nil
}

func (c *ChainClientMock) InvokeContractByVc(ledgerId int, contractAddress string, funcName string, vc string) ([]string, *bcostypes.TransactionDetail, error) {
	return nil, nil, nil
}

// GetTxProve 获取交易证明
//
//	@receiver c
//	@param tx
//	@param chainId
//	@return string
func (c *ChainClientMock) GetTxProve(tx *bcostypes.TransactionDetail, chainId string) string {
	return "{}"
}

// TxProve 交易认证
//
//	@receiver c
//	@param txProve
//	@return bool
func (c *ChainClientMock) TxProve(txProve string) bool {
	return true
}

// CheckChain 检查链
//
//	@receiver c
//	@return bool
func (c *ChainClientMock) CheckChain() bool {
	return true
}
