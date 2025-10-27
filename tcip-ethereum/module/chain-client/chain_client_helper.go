/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/conf"
)

const (
	keySuffix = ".key"
	crtSuffix = ".crt"
	pemSuffix = ".pem"
)

// EthConfigTemplate bcos配置结构体
type EthConfigTemplate struct {
	ChainID int64
}

// NewEthConfigTemplate  新建bcos的配置结构
//
//	@param chainID
//	@param groupIDStr
//	@param address
//	@param isSMCrypto
//	@return *EthConfigTemplate
//	@return error
func NewEthConfigTemplate(chainID int64) (*EthConfigTemplate, error) {
	return &EthConfigTemplate{
		ChainID: chainID,
	}, nil
}

// createSDK 创建bcos的sdk
//
//	@param log
//	@return *client.Client
//	@return error
func createSDK(chainConfig *conf.EthConfig,
	log *zap.SugaredLogger) (bool, error) {
	// 此处使用的链ID
	//template, err := NewEthConfigTemplate(chainConfig.ChainId)
	_, err := NewEthConfigTemplate(chainConfig.ChainId)
	if err != nil {
		log.Errorf("[createSDK]NewEthConfigTemplate error %v", err)
		return false, err
	}
	log.Infof("[createSDK]Use chain config %v", chainConfig)
	return true, err
}
