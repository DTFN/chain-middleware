/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package accesscontrol

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
)

var (
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	srcGw = &common.GatewayInfo{
		GatewayId:     "0",
		ToGatewayList: []string{"1"},
	}
	destGw = &common.GatewayInfo{
		GatewayId:       "1",
		FromGatewayList: []string{"0"},
	}
)

func TestAccessControl_GatewayPermissionsCheck(t *testing.T) {
	logger.InitLogConfig(log)
	err := InitAccessControl()
	assert.Nil(t, err)
	res := AccessV1.GatewayPermissionsCheck(srcGw, destGw)
	assert.True(t, res)

	destGw.GatewayId = "3"
	res = AccessV1.GatewayPermissionsCheck(srcGw, destGw)
	assert.False(t, res)

	destGw.FromGatewayList = []string{}
	srcGw.ToGatewayList = []string{}
	res = AccessV1.GatewayPermissionsCheck(srcGw, destGw)
	assert.True(t, res)
}

func TestInitAccessControl(t *testing.T) {
	logger.InitLogConfig(log)
	err := InitAccessControl()
	assert.Nil(t, err)
}
