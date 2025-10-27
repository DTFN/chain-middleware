/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package accesscontrol

import (
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"go.uber.org/zap"
)

// AccessV1 鉴权全局变量
var AccessV1 *AccessControl

// AccessControl 鉴权结构体
type AccessControl struct {
	log *zap.SugaredLogger
}

// InitAccessControl 初始化鉴权模块
func InitAccessControl() error {
	AccessV1 = &AccessControl{
		log: logger.GetLogger(logger.ModuleAccessControl),
	}
	return nil
}

// GatewayPermissionsCheck gateway鉴权
//  @receiver a
//  @param srcGatewayInfo
//  @param destGatewayInfo
//  @return bool
func (a *AccessControl) GatewayPermissionsCheck(srcGatewayInfo, destGatewayInfo *common.GatewayInfo) bool {
	srcAllow := false
	destAllow := false
	if len(srcGatewayInfo.ToGatewayList) == 0 && len(destGatewayInfo.FromGatewayList) == 0 {
		return true
	}
	for _, v := range srcGatewayInfo.ToGatewayList {
		if v == destGatewayInfo.GatewayId {
			srcAllow = true
		}
	}
	if srcAllow && len(destGatewayInfo.FromGatewayList) == 0 {
		return true
	}
	for _, v := range destGatewayInfo.FromGatewayList {
		if v == srcGatewayInfo.GatewayId {
			destAllow = true
		}
	}
	if destAllow && len(srcGatewayInfo.ToGatewayList) == 0 {
		return true
	}
	if srcAllow && destAllow {
		return true
	}

	return false
}
