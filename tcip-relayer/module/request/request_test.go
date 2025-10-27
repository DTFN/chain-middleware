/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"os"
	"path"
	"strings"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"

	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"github.com/stretchr/testify/assert"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

func testInit() {
	logger.InitLogConfig(log)
	_ = relay_chain_chainmaker.InitRelayChainMock(nil)
	_ = InitRequestManagerMock()
}

func TestInitRequestManager(t *testing.T) {
	logger.InitLogConfig(log)
	err := InitRequestManager()
	assert.Nil(t, err)
}

func TestRequestManager_CrossChainTry(t *testing.T) {
	testInit()
	req := &cross_chain.CrossChainTryRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := RequestV1.CrossChainTry(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "Unsupported call type"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_GRPC
	res, err = RequestV1.CrossChainTry(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_REST
	res, err = RequestV1.CrossChainTry(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	req.Version = common.Version_V1_0_0
	res, err = RequestV1.CrossChainTry(req, 10, gatewayInfo)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestRequestManager_CrossChainConfirm(t *testing.T) {
	testInit()
	req := &cross_chain.CrossChainConfirmRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := RequestV1.CrossChainConfirm(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "Unsupported call type"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_GRPC
	res, err = RequestV1.CrossChainConfirm(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_REST
	res, err = RequestV1.CrossChainConfirm(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	req.Version = common.Version_V1_0_0
	res, err = RequestV1.CrossChainConfirm(req, 10, gatewayInfo)
	assert.Nil(t, err)
	assert.Equal(t, res, &cross_chain.CrossChainConfirmResponse{})
}

func TestRequestManager_CrossChainCancel(t *testing.T) {
	testInit()
	req := &cross_chain.CrossChainCancelRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := RequestV1.CrossChainCancel(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "Unsupported call type"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_GRPC
	res, err = RequestV1.CrossChainCancel(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_REST
	res, err = RequestV1.CrossChainCancel(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	req.Version = common.Version_V1_0_0
	res, err = RequestV1.CrossChainCancel(req, 10, gatewayInfo)
	assert.Nil(t, err)
	assert.Equal(t, res, &cross_chain.CrossChainCancelResponse{})
}

func TestRequestManager_IsCrossChainSuccess(t *testing.T) {
	testInit()
	req := &cross_chain.IsCrossChainSuccessRequest{
		Version: common.Version(10),
	}
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	res, err := RequestV1.IsCrossChainSuccess(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "Unsupported call type"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_GRPC
	res, err = RequestV1.IsCrossChainSuccess(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	gatewayInfo.CallType = common.CallType_REST
	res, err = RequestV1.IsCrossChainSuccess(req, 10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "unsupported version"))
	assert.Nil(t, res)

	req.Version = common.Version_V1_0_0
	res, err = RequestV1.IsCrossChainSuccess(req, 10, gatewayInfo)
	assert.Nil(t, err)
	assert.Equal(t, res, &cross_chain.IsCrossChainSuccessResponse{})
}

func TestRequestManager_PingPong(t *testing.T) {
	testInit()
	gatewayInfo := &common.GatewayInfo{
		CallType: common.CallType(10),
	}
	_, err := RequestV1.PingPong(10, gatewayInfo)
	assert.True(t, strings.Contains(err.Error(), "Unsupported call type"))

	gatewayInfo.CallType = common.CallType_GRPC
	_, err = RequestV1.PingPong(10, gatewayInfo)
	assert.Nil(t, err)

	gatewayInfo.CallType = common.CallType_REST
	_, err = RequestV1.PingPong(10, gatewayInfo)
	assert.Nil(t, err)
}

func TestRequestManager_VerifyTx(t *testing.T) {
	testInit()
	txVerifyInterface := &common.TxVerifyInterface{}
	res, err := RequestV1.VerifyTx(txVerifyInterface, "{}")
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(res), "true"))

	res, err = RequestV1.VerifyTx(txVerifyInterface, "")
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(res), "false"))
}
