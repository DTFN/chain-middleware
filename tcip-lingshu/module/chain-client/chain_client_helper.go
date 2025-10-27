/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/conf"
	"go.uber.org/zap"
)

const (
	keySuffix = ".key"
	crtSuffix = ".crt"
	pemSuffix = ".pem"
)

// BCOSConfigTemplate bcos配置结构体
type LingShuConfigTemplate struct {
	ChainID int64
}

// LingShuConfigTemplate  新建配置结构
//
//	@param chainID
//	@param groupIDStr
//	@param address
//	@param isSMCrypto
//	@return *LingShuConfigTemplate
//	@return error
func NewLingShuConfigTemplate(chainID int64) (*LingShuConfigTemplate, error) {
	return &LingShuConfigTemplate{
		ChainID: chainID,
	}, nil
}

// createSDK 创建dk
//
//	@param lingshuConfig
//	@param log
//	@return *client.Client
//	@return error
func createSDK(lingshuConfig *conf.LingShuConfig,
	log *zap.SugaredLogger) (bool, error) {
	// 此处使用的链ID
	_, err := NewLingShuConfigTemplate(lingshuConfig.ChainId)
	if err != nil {
		log.Errorf("[createSDK]NewEthConfigTemplate error %v", err)
		return false, err
	}
	log.Infof("[createSDK]Use chain config %v", lingshuConfig)
	return true, err
}
