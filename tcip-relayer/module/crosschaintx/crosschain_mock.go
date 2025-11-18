/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package crosschaintx

import (
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
)

// InitCrossChainManagerMock 初始化跨链管理模块
//  @return error
func InitCrossChainManagerMock() error {
	CrossChainTxV1 = &CrossChainTxManager{
		log: logger.GetLogger(logger.ModuleCrossChainTx),
	}

	return nil
}
