/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package crosschaintx

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
)

const (
	zero  = "0"
	one   = "1"
	two   = "2"
	three = "3"
	four  = "4"
	five  = "5"
	six   = "6"
	seven = "7"
	eight = "8"
	nine  = "9"
)

var (
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	req = &relay_chain.BeginCrossChainRequest{
		Version:        common.Version_V1_0_0,
		CrossChainId:   "0",
		CrossChainName: "test",
		CrossChainFlag: "test",
		CrossChainMsg: []*common.CrossChainMsg{
			{
				GatewayId:    "0",
				ChainRid:     "chain1",
				ContractName: "test",
				Method:       "test",
				Parameter:    "{\"a\":\"MA==\"}",
				ConfirmInfo: &common.ConfirmInfo{
					ChainRid: "chain1",
				},
				CancelInfo: &common.CancelInfo{
					ChainRid: "chain1",
				},
			},
		},
		TxContent: &common.TxContent{},
		From:      "0",
		Timeout:   1000,
		ConfirmInfo: &common.ConfirmInfo{
			ChainRid: "chain1",
		},
		CancelInfo: &common.CancelInfo{
			ChainRid: "chain1",
		},
		CrossType: common.CrossType_INVOKE,
	}
)

func testInit() {
	logger.InitLogConfig(log)
	conf.Config.RelayChain = &conf.RelayChain{}
	conf.Config.BaseConfig = &conf.BaseConfig{
		GatewayID: zero,
	}
	_ = relay_chain_chainmaker.InitRelayChainMock(nil)
	_ = gateway.InitGatewayManagerMock()
	_ = request.InitRequestManagerMock()
}

func TestInitCrossChainManager(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)
}

func TestNewCrossChainInfo(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	crossChainInfo := CrossChainTxV1.BuildCrossChainInfoFromBeginCrossChainRequest(req)
	assert.NotNil(t, crossChainInfo)
	//	res, _ := proto.Marshal(crossChainInfo)
	//	fmt.Println(res)

	crossChainId, err := CrossChainTxV1.NewCrossChainInfo(crossChainInfo)
	assert.Nil(t, err)
	assert.Equal(t, crossChainId, "0")
}

func TestGetCrossChainInfo(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	crossChainInfo, err := CrossChainTxV1.GetCrossChainInfo("")
	assert.Nil(t, err)
	assert.NotNil(t, crossChainInfo)

	crossChainInfo, err = CrossChainTxV1.GetCrossChainInfo("100")
	assert.Nil(t, crossChainInfo)
	assert.NotNil(t, err)
}

func TestGetCrossChainInfoByRange(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	crossChainInfo, err := CrossChainTxV1.GetCrossChainInfoByRange("", "")
	assert.Nil(t, err)
	assert.Equal(t, len(crossChainInfo), 1)
}

func TestGetCrossChainNum(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	num, err := CrossChainTxV1.GetCrossChainNum()
	assert.Nil(t, err)
	assert.Equal(t, num, uint64(1))
}

func TestGetNotEndCrossChainIdList(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	crossChainInfo, err := CrossChainTxV1.GetNotEndCrossChainIdList()
	assert.Nil(t, err)
	assert.Equal(t, len(crossChainInfo), 1)
}

func TestCrossChainTry(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	CrossChainTxV1.crossChainTry(zero)
	CrossChainTxV1.crossChainTry("99")
	CrossChainTxV1.crossChainTry(four)
	CrossChainTxV1.crossChainTry(five)
	CrossChainTxV1.crossChainTry(six)
	CrossChainTxV1.crossChainTry(seven)
}

func TestCrossChainResult(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	CrossChainTxV1.crossChainResult(one)
	CrossChainTxV1.crossChainResult(three)
	CrossChainTxV1.crossChainResult("99")
	CrossChainTxV1.crossChainResult(six)
}

func TestCrosschainConfirm(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	CrossChainTxV1.crosschainConfirm(one)
	CrossChainTxV1.crosschainConfirm(two)
	CrossChainTxV1.crosschainConfirm("99")
	CrossChainTxV1.crosschainConfirm(eight)
	CrossChainTxV1.crosschainConfirm(nine)
}

func TestCrosschainSrcGatewayConfirm(t *testing.T) {
	testInit()
	err := InitCrossChainManager()
	assert.Nil(t, err)

	CrossChainTxV1.crosschainSrcGatewayConfirm(one)
	CrossChainTxV1.crosschainSrcGatewayConfirm(two)
	CrossChainTxV1.crosschainSrcGatewayConfirm("99")
	CrossChainTxV1.crosschainSrcGatewayConfirm(six)
}

func TestBuildTxContent(t *testing.T) {
	res := buildTxContent("123", []byte("123"), common.TxResultValue_TX_SUCCESS,
		"0", "chain1", "123")
	assert.NotNil(t, res)
}

func TestInitCrossChainManagerMock(t *testing.T) {
	testInit()
	err := InitCrossChainManagerMock()
	assert.Nil(t, err)
}
