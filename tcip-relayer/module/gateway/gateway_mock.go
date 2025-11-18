/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"github.com/emirpasic/gods/maps/treemap"
)

// InitGatewayManagerMock 网关mock
//  @return error
func InitGatewayManagerMock() error {
	GatewayV1 = &GatewayManager{
		gatewayInfo:  treemap.NewWithStringComparator(),
		gateWayState: treemap.NewWithStringComparator(),
		log:          logger.GetLogger(logger.ModuleGateway),
	}
	return nil
}
