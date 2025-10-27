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

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/utils"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/db"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/logger"
)

const (
	sdkKeyText = "-----BEGIN PRIVATE KEY-----\nMIGEAgEAMBAGByqGSM49AgEGBSuBBAAKBG0wawIBAQQgLGUmixHrD7qjlFeQYUVt\nTqAcwPd6YemZqF5bz/YzkyehRANCAAT98Wd9JW1Fv7xAOyN5S+GQREij4McJcc+H\njcHwt9gG6vj2MLkeF9iHNxDeD4WRihOfSSwGpe3v37qMm1yIE7OZ\n-----END PRIVATE KEY-----"
	sdkCrtText = "-----BEGIN CERTIFICATE-----\nMIIBgzCCASmgAwIBAgIUBn7qQz2uMAJw1osCUD9sjFBk91wwCgYIKoZIzj0EAwIw\nNzEPMA0GA1UEAwwGYWdlbmN5MRMwEQYDVQQKDApmaXNjby1iY29zMQ8wDQYDVQQL\nDAZhZ2VuY3kwIBcNMjExMjE2MDk1NDA2WhgPMjEyMTExMjIwOTU0MDZaMDExDDAK\nBgNVBAMMA3NkazETMBEGA1UECgwKZmlzY28tYmNvczEMMAoGA1UECwwDc2RrMFYw\nEAYHKoZIzj0CAQYFK4EEAAoDQgAE/fFnfSVtRb+8QDsjeUvhkERIo+DHCXHPh43B\n8LfYBur49jC5HhfYhzcQ3g+FkYoTn0ksBqXt79+6jJtciBOzmaMaMBgwCQYDVR0T\nBAIwADALBgNVHQ8EBAMCBeAwCgYIKoZIzj0EAwIDSAAwRQIgHxz9ZQMgic52HvML\nt8AmSlGMo33nDpV6Nz7SuiezdqECIQCYLIP6nN7W/aj6+eqhjcKn5XJAypgGuI5y\nqEOBFLer3w==\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBezCCASGgAwIBAgIUd/vq/b9+CYOdr2a4lGqjd1s8PiAwCgYIKoZIzj0EAwIw\nNTEOMAwGA1UEAwwFY2hhaW4xEzARBgNVBAoMCmZpc2NvLWJjb3MxDjAMBgNVBAsM\nBWNoYWluMB4XDTIxMTIxNjA5NTQwNloXDTMxMTIxNDA5NTQwNlowNzEPMA0GA1UE\nAwwGYWdlbmN5MRMwEQYDVQQKDApmaXNjby1iY29zMQ8wDQYDVQQLDAZhZ2VuY3kw\nVjAQBgcqhkjOPQIBBgUrgQQACgNCAAQ270cEs1AcnLtARy8WYcVjgP7HCfA+GeEN\nniwbMU8er4IOZ9WM6ihaeHUNt/TkOgo7Xc4Mw1IBwN/k1q2GlpycoxAwDjAMBgNV\nHRMEBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIFlTs/ZmN1qvTGQiBBQelCY2gi96\n5STdrm4La0ENQOcSAiEA3DDubVr/Y/9BBO9eyI12w6PmK+3J5xxC/rUHmIjlDMc=\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIBvjCCAWSgAwIBAgIUBCDLHI2oWBXSywRsYGWPn3zokK8wCgYIKoZIzj0EAwIw\nNTEOMAwGA1UEAwwFY2hhaW4xEzARBgNVBAoMCmZpc2NvLWJjb3MxDjAMBgNVBAsM\nBWNoYWluMCAXDTIxMTIxNjA5NTQwNloYDzIxMjExMTIyMDk1NDA2WjA1MQ4wDAYD\nVQQDDAVjaGFpbjETMBEGA1UECgwKZmlzY28tYmNvczEOMAwGA1UECwwFY2hhaW4w\nVjAQBgcqhkjOPQIBBgUrgQQACgNCAARlf+1VJLYJyjuNVnw9rXQ4zNB+Sucix2vJ\n7bviXgyuvtu2cZHC5/BZ8l5ODMqSlPpKn9qWJUmxi3vC8szWXZcqo1MwUTAdBgNV\nHQ4EFgQUMBD7X1irOaZIPCvyaquGVSyHzyQwHwYDVR0jBBgwFoAUMBD7X1irOaZI\nPCvyaquGVSyHzyQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiEA\n2lpAxoB/kWnD6Mv/4Q/hVSby3U/6BM/gTlOq/kTQuoECIHI2Yi0CnyqZciUujliY\nbRRI7XZWJ41h6KE7B4qkzB0T\n-----END CERTIFICATE-----"
	caCrtText  = "-----BEGIN CERTIFICATE-----\nMIIBvjCCAWSgAwIBAgIUBCDLHI2oWBXSywRsYGWPn3zokK8wCgYIKoZIzj0EAwIw\nNTEOMAwGA1UEAwwFY2hhaW4xEzARBgNVBAoMCmZpc2NvLWJjb3MxDjAMBgNVBAsM\nBWNoYWluMCAXDTIxMTIxNjA5NTQwNloYDzIxMjExMTIyMDk1NDA2WjA1MQ4wDAYD\nVQQDDAVjaGFpbjETMBEGA1UECgwKZmlzY28tYmNvczEOMAwGA1UECwwFY2hhaW4w\nVjAQBgcqhkjOPQIBBgUrgQQACgNCAARlf+1VJLYJyjuNVnw9rXQ4zNB+Sucix2vJ\n7bviXgyuvtu2cZHC5/BZ8l5ODMqSlPpKn9qWJUmxi3vC8szWXZcqo1MwUTAdBgNV\nHQ4EFgQUMBD7X1irOaZIPCvyaquGVSyHzyQwHwYDVR0jBBgwFoAUMBD7X1irOaZI\nPCvyaquGVSyHzyQwDwYDVR0TAQH/BAUwAwEB/zAKBggqhkjOPQQDAgNIADBFAiEA\n2lpAxoB/kWnD6Mv/4Q/hVSby3U/6BM/gTlOq/kTQuoECIHI2Yi0CnyqZciUujliY\nbRRI7XZWJ41h6KE7B4qkzB0T\n-----END CERTIFICATE-----"
	accountPem = "-----BEGIN PRIVATE KEY-----\nMIGNAgEAMBAGByqGSM49AgEGBSuBBAAKBHYwdAIBAQQgMJGybVfcQv5XaWyqBU+N\n3hRcdJMxkqLpChBwzspnM06gBwYFK4EEAAqhRANCAAThfDgQMhjaf0mxRCb2oOOC\n8CYjxMNNHu37T+uRzeewz4Af/02qXB+fst5tSAw6rMKtUe7xBL4H+RXRk8/GN8yU\n-----END PRIVATE KEY-----"
	chainID    = 1
	groupID    = "1"
	address    = "127.0.0.1:8080"
)

var (
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	chainConfig = &conf.LingShuConfig{
		ChainRid: "chain001",
		ChainId:  chainID,
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
