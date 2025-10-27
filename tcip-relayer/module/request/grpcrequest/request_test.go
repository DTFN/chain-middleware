/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package grpcrequest

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"github.com/stretchr/testify/assert"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

func testInit() *GrpcRequest {
	logger.InitLogConfig(log)
	return NewGrpcRequest(logger.GetLogger(logger.ModuleRequest))
}

func TestGrpcRequest_CrossChainCancel(t *testing.T) {
	grpcRequest := testInit()
	req := &cross_chain.CrossChainCancelRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := grpcRequest.CrossChainCancel(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "missing address"))
}

func TestGrpcRequest_CrossChainConfirm(t *testing.T) {
	grpcRequest := testInit()
	req := &cross_chain.CrossChainConfirmRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := grpcRequest.CrossChainConfirm(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "missing address"))
}

func TestGrpcRequest_CrossChainTry(t *testing.T) {
	grpcRequest := testInit()
	req := &cross_chain.CrossChainTryRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := grpcRequest.CrossChainTry(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "missing address"))
}

func TestGrpcRequest_IsCrossChainSuccess(t *testing.T) {
	grpcRequest := testInit()
	req := &cross_chain.IsCrossChainSuccessRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := grpcRequest.IsCrossChainSuccess(req, 1, gatewayInfo)
	assert.Nil(t, res)
	assert.True(t, strings.Contains(err.Error(), "missing address"))
}

func TestGrpcRequest_PingPong(t *testing.T) {
	grpcRequest := testInit()
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	_, err := grpcRequest.PingPong(1, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "missing address"))
}

func TestGrpcRequest_VerifyTx(t *testing.T) {
	grpcRequest := testInit()
	_, err := grpcRequest.VerifyTx(nil, "")
	assert.True(t, strings.Contains(err.Error(), "implement this method"))
}
