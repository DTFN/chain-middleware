/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_config

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/utils"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/db"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

var (
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	chainConfig = &common.ChainmakerConfig{
		ChainRid: "chain001",
		ChainId:  "mychannel",
		OrgId:    "org1",
		SignKey:  "signkey",
		SignCert: "signcert",
		TlsKey:   "sdkKeyText",
		TlsCert:  "sdkCrtText",
		Node: []*common.Node{
			{
				NodeAddr: "127.0.0.1:12301",
				TrustRoot: []string{
					"123",
				},
				ConnCnt:     10,
				EnableTls:   true,
				TlsHostName: "chainmaker.org",
			},
		},
	}
)

func initTest() {
	conf.Config.BaseConfig = &conf.BaseConfig{
		GatewayID:    "0",
		GatewayName:  "test",
		Address:      "https://127.0.0.1:19999",
		ServerName:   "chainmaker.org",
		Tlsca:        "../../config/cert/client/ca.crt",
		ClientKey:    "../../config/cert/client/client.key",
		ClientCert:   "../../config/cert/client/client.crt",
		TxVerifyType: "notneed",
		CallType:     "grpc",
	}
	conf.Config.DbPath = path.Join(os.TempDir(), time.Now().String())
	logger.InitLogConfig(log)
	db.NewDbHandle()
	utils.UpdateChainConfigChan = make(chan *utils.ChainConfigOperate)
	go listenChan(utils.UpdateChainConfigChan)
	NewChainConfig()
}

func listenChan(updateChan chan *utils.ChainConfigOperate) {
	for {
		<-updateChan
	}
}

func TestSave(t *testing.T) {
	initTest()
	NewChainConfig()
	err := ChainConfigManager.Save(chainConfig, common.Operate_UPDATE)
	assert.NotNil(t, err)

	err = ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.Nil(t, err)

	err = ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.NotNil(t, err)

	err = ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.NotNil(t, err)

	err = ChainConfigManager.Save(chainConfig, common.Operate_UPDATE)
	assert.Nil(t, err)

	err = ChainConfigManager.Save(chainConfig, common.Operate_DELETE)
	assert.NotNil(t, err)

	_ = db.Db.Close()
	err = ChainConfigManager.Save(chainConfig, common.Operate_UPDATE)
	assert.NotNil(t, err)

	chainConfig.ChainRid = ""
	err = ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	chainConfig.ChainRid = "chain001"
	assert.NotNil(t, err)

	err = ChainConfigManager.Save(&common.ChainmakerConfig{
		ChainRid: "chain002",
	}, common.Operate_SAVE)
	assert.NotNil(t, err)

	err = ChainConfigManager.Save(&common.ChainmakerConfig{
		ChainRid: "chain002",
		Node: []*common.Node{
			{}, {},
		},
	}, common.Operate_SAVE)
	assert.NotNil(t, err)
}

func TestDelete(t *testing.T) {
	initTest()
	NewChainConfig()
	err := ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.Nil(t, err)

	err = ChainConfigManager.Delete(chainConfig.ChainRid)
	assert.Nil(t, err)

	err = ChainConfigManager.Delete(chainConfig.ChainRid)
	assert.NotNil(t, err)

	_ = db.Db.Close()
	err = ChainConfigManager.Delete(chainConfig.ChainRid)
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	initTest()
	NewChainConfig()
	err := ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.Nil(t, err)

	res, err := ChainConfigManager.Get(chainConfig.ChainRid)
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	_, err = ChainConfigManager.Get("213")
	assert.NotNil(t, err)

	res, err = ChainConfigManager.Get("")
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	_ = db.Db.Close()
	_, err = ChainConfigManager.Get("213")
	assert.NotNil(t, err)

	res, err = ChainConfigManager.Get("")
	assert.Nil(t, err)
	assert.Equal(t, len(res), 0)
}

func TestSetState(t *testing.T) {
	initTest()
	NewChainConfig()
	err := ChainConfigManager.Save(chainConfig, common.Operate_SAVE)
	assert.Nil(t, err)

	err = ChainConfigManager.SetState(chainConfig, true, "success")
	assert.Nil(t, err)

	_ = db.Db.Close()
	err = ChainConfigManager.SetState(chainConfig, true, "success")
	assert.NotNil(t, err)
}
