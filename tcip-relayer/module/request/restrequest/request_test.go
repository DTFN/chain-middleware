/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package restrequest

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "[test]",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

func testInit() *RestRequest {
	logger.InitLogConfig(log)
	return NewRestRequest(logger.GetLogger(logger.ModuleRequest))
}

func TestRestRequest_CrossChainCancel(t *testing.T) {
	restRequest := testInit()
	req := &cross_chain.CrossChainCancelRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := restRequest.CrossChainCancel(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "unsupported protocol scheme"))
}

func TestRestRequest_CrossChainConfirm(t *testing.T) {
	restRequest := testInit()
	req := &cross_chain.CrossChainConfirmRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := restRequest.CrossChainConfirm(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "unsupported protocol scheme"))
}

func TestRestRequest_CrossChainTry(t *testing.T) {
	restRequest := testInit()
	req := &cross_chain.CrossChainTryRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := restRequest.CrossChainTry(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "unsupported protocol scheme"))
}

func TestRestRequest_IsCrossChainSuccess(t *testing.T) {
	restRequest := testInit()
	req := &cross_chain.IsCrossChainSuccessRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := restRequest.IsCrossChainSuccess(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "unsupported protocol scheme"))
}

func TestRestRequest_PingPong(t *testing.T) {
	restRequest := testInit()
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	_, err := restRequest.PingPong(1, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported protocol scheme"))
}

func TestRestRequest_VerifyTx(t *testing.T) {
	conf.Config.BaseConfig = &conf.BaseConfig{
		DefaultTimeout: 1,
	}
	restRequest := testInit()
	txVerifyInterface := &common.TxVerifyInterface{}
	res, err := restRequest.VerifyTx(txVerifyInterface, "{}")
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "http: no Host in request URL"))
}
