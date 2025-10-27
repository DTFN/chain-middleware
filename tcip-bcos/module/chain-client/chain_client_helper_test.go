/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"github.com/stretchr/testify/require"
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
	bcosConfig = &common.BcosConfig{
		ChainRid:   "chain001",
		ChainId:    chainID,
		TlsKey:     sdkKeyText,
		TlsCert:    sdkCrtText,
		PrivateKey: accountPem,
		GroupId:    groupID,
		Address:    address,
		Http:       false,
		IsSmCrypto: false,
	}
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
)

func TestBCOSTemplate(t *testing.T) {
	template, err := NewBCOSConfigTemplate(chainID, groupID, address, false)
	require.Nil(t, err)
	err = template.SetCa(caCrtText)
	require.Nil(t, err)
	err = template.SetKey([]byte(sdkKeyText))
	require.Nil(t, err)
	err = template.SetCert([]byte(sdkCrtText))
	require.Nil(t, err)
	err = template.SetAccount([]byte(accountPem))
	require.Nil(t, err)
	fmt.Println(template)
}

func TestCreateSDK(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			//fmt.Println(e)
			assert.Equal(t, e.(string), "failed to parse root certificate")
		}
	}()
	logger.InitLogConfig(log)
	sdk, _ := createSDK(bcosConfig, logger.GetLogger(logger.ModuleChainClient))
	assert.Nil(t, sdk)
	//fmt.Println(err)
}
