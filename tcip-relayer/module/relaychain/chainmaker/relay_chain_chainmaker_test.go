/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package relay_chain_chainmaker

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
)

var relayChainConfig = &conf.RelayChain{
	ChainmakerSdkConfigPath: "./test/sdk_config.yml",
	Users: []*conf.User{
		{
			SignCrtPath: "../../../config/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt",
			SignKeyPath: "../../../config/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key",
			OrgId:       "wx-org1.chainmaker.org",
		},
		{
			SignCrtPath: "../../../config/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt",
			SignKeyPath: "../../../config/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key",
			OrgId:       "wx-org2.chainmaker.org",
		},
		{
			SignCrtPath: "../../../config/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt",
			SignKeyPath: "../../../config/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key",
			OrgId:       "wx-org3.chainmaker.org",
		},
		{
			SignCrtPath: "../../../config/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt",
			SignKeyPath: "../../../config/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key",
			OrgId:       "wx-org4.chainmaker.org",
		},
	},
}

var logConfig = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		LogLevel:     "DEBUG",
		FilePath:     "/tmp/123456",
		LogInConsole: true,
	},
}

func TestInitRelayChain(t *testing.T) {
	logger.InitLogConfig(logConfig)
	err := InitRelayChain(relayChainConfig)
	assert.NotNil(t, err)
}

func TestGetKvsFromKvJsonStr(t *testing.T) {
	relayChain := RelayChainChainmaker{}
	logger.InitLogConfig(logConfig)
	relayChain.log = logger.GetLogger(logger.ModuleRelayChain)
	res, err := relayChain.getKvsFromKvJsonStr("{\"a\":\"1\",\"b\":\"2\"}")
	assert.Nil(t, res)
	assert.NotNil(t, err)

	res, err = relayChain.getKvsFromKvJsonStr("{\"a\":\"MQ==\",\"b\":\"MQ==\"}")
	assert.Equal(t, len(res), 2)
	assert.Nil(t, err)
}

//func TestGetEndorsersWithAuthType(t *testing.T) {
//	relayChain := RelayChainChainmaker{}
//	logger.InitLogConfig(logConfig)
//	relayChain.log = logger.GetLogger(logger.ModuleRelayChain)
//
//	endorsers, err := relayChain.getEndorsersWithAuthType(crypto.HASH_TYPE_SHA256, sdk.PermissionedWithCert,
//		&common.Payload{ChainId: "chain1"}, relayChainConfig.Users)
//	assert.Nil(t, err)
//	assert.Equal(t, len(endorsers), 4)
//
//	endorsers, err = relayChain.getEndorsersWithAuthType(crypto.HASH_TYPE_SHA256, sdk.PermissionedWithKey,
//		&common.Payload{ChainId: "chain1"}, relayChainConfig.Users)
//	assert.Nil(t, err)
//	assert.Equal(t, len(endorsers), 4)
//
//	endorsers, err = relayChain.getEndorsersWithAuthType(crypto.HASH_TYPE_SHA256, sdk.Public,
//		&common.Payload{ChainId: "chain1"}, relayChainConfig.Users)
//	assert.Nil(t, err)
//	assert.Equal(t, len(endorsers), 4)
//
//	endorsers, err = relayChain.getEndorsersWithAuthType(crypto.HASH_TYPE_SHA256, sdk.AuthType(100),
//		&common.Payload{ChainId: "chain1"}, relayChainConfig.Users)
//	assert.NotNil(t, err)
//	assert.Nil(t, endorsers)
//
//	relayChainConfig.Users[0].SignCrtPath = "abcd1234"
//	endorsers, err = relayChain.getEndorsersWithAuthType(crypto.HASH_TYPE_SHA256, sdk.PermissionedWithCert,
//		&common.Payload{ChainId: "chain1"}, relayChainConfig.Users)
//	assert.NotNil(t, err)
//	assert.Nil(t, endorsers)
//}

func TestInitRelayChainMock(t *testing.T) {
	err := InitRelayChainMock(nil)
	assert.Nil(t, err)
}

func TestInitContract(t *testing.T) {
	err := InitRelayChainMock(nil)
	assert.Nil(t, err)

	err = RelayChainV1.InitContract("", "", "", "",
		false, 100, tcipcommon.ChainmakerRuntimeType_DOCKER_GO)
	assert.Nil(t, err)
}

func TestUpdateContract(t *testing.T) {
	err := InitRelayChainMock(nil)
	assert.Nil(t, err)

	err = RelayChainV1.UpdateContract("", "", "", "",
		false, 100, tcipcommon.ChainmakerRuntimeType_DOCKER_GO)
	assert.Nil(t, err)
}

func TestInvokeContract(t *testing.T) {
	err := InitRelayChainMock(nil)
	assert.Nil(t, err)

	zeroByte := []byte("0")
	res, err := RelayChainV1.InvokeContract("", saveGatewayMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, zeroByte)

	res, err = RelayChainV1.InvokeContract("", updateGatewayMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, zeroByte)

	res, err = RelayChainV1.InvokeContract("", utils.SaveCrossChainInfo, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, zeroByte)

	res, err = RelayChainV1.InvokeContract("", "", false, "",
		100)
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestQueryContract(t *testing.T) {
	logger.InitLogConfig(logConfig)
	err := InitRelayChainMock(nil)
	assert.Nil(t, err)

	res, err := RelayChainV1.QueryContract("", utils.GetBlockHeaderMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, []byte("test"))

	res, err = RelayChainV1.QueryContract("", utils.SpvTxVerifyMethod, false, "{}",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, []byte("true"))

	kvJsonStr := fmt.Sprintf("{\"%s\":\"MA==\"}", syscontract.GetGateway_GATEWAY_ID.String())
	res, err = RelayChainV1.QueryContract("", utils.SpvTxVerifyMethod, false, kvJsonStr,
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, []byte("false"))

	res, err = RelayChainV1.QueryContract("", getGatewayMethod, false, kvJsonStr,
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	kvJsonStr = fmt.Sprintf("{\"%s\":\"OTk=\"}", syscontract.GetGateway_GATEWAY_ID.String())
	res, err = RelayChainV1.QueryContract("", getGatewayMethod, false, kvJsonStr,
		100)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	res, err = RelayChainV1.QueryContract("", getGatewayByRangeMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	res, err = RelayChainV1.QueryContract("", getGatewayNumMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, []byte("5"))

	res, err = RelayChainV1.QueryContract("", getCrossChainNumMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.Equal(t, res, []byte("1"))

	res, err = RelayChainV1.QueryContract("", getCrossChainMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	kvJsonStr = fmt.Sprintf("{\"%s\":\"OTk=\"}", syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String())
	res, err = RelayChainV1.QueryContract("", getCrossChainMethod, false, kvJsonStr,
		100)
	assert.NotNil(t, err)
	assert.Nil(t, res)

	kvJsonStr = fmt.Sprintf("{\"%s\":\"Mg==\"}", syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String())
	res, err = RelayChainV1.QueryContract("", getCrossChainMethod, false, kvJsonStr,
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	res, err = RelayChainV1.QueryContract("", utils.GetNotEndCrossChainIdList, false, "",
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	res, err = RelayChainV1.QueryContract("", getCrossChainByRangeMethod, false, "",
		100)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	res, err = RelayChainV1.QueryContract("", "", false, "",
		100)
	assert.Nil(t, err)
	assert.Nil(t, res)
}
