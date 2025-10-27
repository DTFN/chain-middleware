/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"github.com/gogo/protobuf/proto"

	"github.com/stretchr/testify/assert"

	chain_config "chainmaker.org/chainmaker/tcip-bcos/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/db"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/utils"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
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
	log         []*logger.LogModuleConfig
	chainConfig *common.BcosConfig
	event       *common.CrossChainEvent
	eventInfo   *EventInfo
)

func initTest() {
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	chainConfig = &common.BcosConfig{
		ChainRid:   "chain001",
		ChainId:    chainID,
		TlsKey:     sdkKeyText,
		TlsCert:    sdkCrtText,
		PrivateKey: accountPem,
		GroupId:    groupID,
		Address:    address,
		Http:       false,
		IsSmCrypto: false,
		Ca:         caCrtText,
	}

	event = &common.CrossChainEvent{
		CrossChainEventId: "00001",
		ChainRid:          "chain001",
		ContractName:      "ADDRESS1",
		CrossChainName:    "test",
		CrossChainFlag:    "test",
		TriggerCondition:  1,
		EventName:         "Transfer(address,address,uint256)",
		Timeout:           1000,
		CrossType:         1,
		BcosEventDataType: []common.EventDataType{8, 8, 8, 4},
		IsCrossChain:      "two == \"0x0000000\"",
		Abi:               []string{"{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}", "{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"}", "{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}"},
		ConfirmInfo: &common.ConfirmInfo{
			ChainRid:      "chain001",
			ContractName:  "ADDRESS2",
			Method:        "transfer",
			Parameter:     "[\"%s\",\"%s\"]",
			ParamData:     []int32{1, 3},
			ParamDataType: []common.EventDataType{8, 4},
			Abi:           "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"minter\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]",
		},
		CancelInfo: &common.CancelInfo{
			ChainRid:      "chain001",
			ContractName:  "ADDRESS2",
			Method:        "burn",
			Parameter:     "[\"%s\"]",
			ParamData:     []int32{3},
			ParamDataType: []common.EventDataType{4},
			Abi:           "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"minter\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]",
		},
		CrossChainCreate: []*common.CrossChainMsg{
			{
				GatewayId:     "1",
				ChainRid:      "chain002",
				ContractName:  "ADDRESS2",
				Method:        "minter",
				Parameter:     "[\"%s\"]",
				ParamData:     []int32{3},
				ParamDataType: []common.EventDataType{4},
				Abi:           "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"minter\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]",
				ConfirmInfo: &common.ConfirmInfo{
					ChainRid:      "chain002",
					ContractName:  "ADDRESS2",
					Method:        "transfer",
					Parameter:     "[\"%s\",\"%s\"]",
					ParamData:     []int32{1, 3},
					ParamDataType: []common.EventDataType{8, 4},
					Abi:           "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"minter\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]",
				},
				CancelInfo: &common.CancelInfo{
					ChainRid:      "chain002",
					ContractName:  "ADDRESS2",
					Method:        "burn",
					Parameter:     "[\"%s\"]",
					ParamData:     []int32{3},
					ParamDataType: []common.EventDataType{4},
					Abi:           "[{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowed\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"minter\",\"outputs\":[{\"name\":\"success\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"initialSupply\",\"type\":\"uint256\"},{\"name\":\"tokenName\",\"type\":\"string\"},{\"name\":\"tokenSymbol\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Burn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"}]",
				},
			},
		},
	}
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
	eventInfo = &EventInfo{
		Topic:        "Transfer(address,address,uint256)",
		ChainRid:     "chain001",
		ContractName: "ADDRESS1",
		TxProve:      "{}",
		Data:         []string{"0x0000000", "0x0000000", "0x0000000", "10"},
		Tx:           []byte("qweqrt"),
		TxId:         "hfjdkhsalufksa",
		BlockHeight:  10,
	}
	conf.Config.DbPath = path.Join(os.TempDir(), time.Now().String())
	logger.InitLogConfig(log)
	db.NewDbHandle()
	utils.EventChan = make(chan *utils.EventOperate)
	utils.UpdateChainConfigChan = make(chan *utils.ChainConfigOperate)
	go listenEventChan(utils.EventChan)
	go listenConfigChan(utils.UpdateChainConfigChan)
	chain_config.NewChainConfig()
}

func listenEventChan(updateChan chan *utils.EventOperate) {
	for {
		<-updateChan
	}
}

func listenConfigChan(updateChan chan *utils.ChainConfigOperate) {
	for {
		<-updateChan
	}
}

func tetsCancel(t *testing.T) {
	event.CrossChainCreate[0].CancelInfo.Abi = "{\a\""
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.Abi = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.Parameter = "%s%s%s%s%s"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func testConfirm(t *testing.T) {
	event.CrossChainCreate[0].ConfirmInfo.Parameter = "%s%s%s%s"
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func testCrossChainCreate(t *testing.T) {
	event.CrossChainCreate[0].Parameter = "%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s"
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamDataType = make([]common.EventDataType, 15)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamData = make([]int32, 10)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ChainRid = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].GatewayId = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func Test_SaveEvent(t *testing.T) {
	initTest()
	InitEventManager()

	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.NotNil(t, err)

	event.Abi[0] = "\"{|a"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	tetsCancel(t)
	testConfirm(t)

	testCrossChainCreate(t)

	event.CancelInfo.ChainRid = "0099"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CancelInfo = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ConfirmInfo.ChainRid = "0099"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ConfirmInfo = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.IsCrossChain = "abcdqwetrtrewq=cdwer"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.IsCrossChain = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.BcosEventDataType = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.TriggerCondition = common.TriggerCondition_COMPLETELY_CONTRACT_EVENT
	err = EventManagerV1.SaveEvent(event, false)
	assert.Nil(t, err)

	event.TriggerCondition = common.TriggerCondition(10)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainEventId = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ChainRid = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.EventName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func Test_DeleteEvent(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	err := EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	err = EventManagerV1.DeleteEvent(event)
	assert.Nil(t, err)

	event.CrossChainEventId = ""
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	event.CrossChainEventId = "123"
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	db.Db.Close()
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)
}

func Test_GetEvent(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	res, err := EventManagerV1.GetEvent("1234")
	assert.NotNil(t, err)
	assert.Nil(t, res)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	res, err = EventManagerV1.GetEvent(event.CrossChainEventId)
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	res, err = EventManagerV1.GetEvent("")
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	err = EventManagerV1.DeleteEvent(event)
	assert.Nil(t, err)

	res, err = EventManagerV1.GetEvent(event.CrossChainEventId)
	assert.NotNil(t, err)
	assert.Equal(t, len(res), 0)
}

func Test_BuildCrossChainMsg(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	req, err := EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.NotNil(t, req)
	assert.Nil(t, err)

	eventInfo.Data[2] = "0xqwertyuiop"

	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.Nil(t, err)

	event.BcosEventDataType[0] = common.EventDataType(100)
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.Nil(t, err)

	event.TriggerCondition = common.TriggerCondition_COMPLETELY_CONTRACT_EVENT
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	reqByte, _ := proto.Marshal(&relay_chain.BeginCrossChainRequest{
		Version: common.Version_V1_0_0,
	})
	eventInfo.Data[1] = string(reqByte)
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.NotNil(t, req)
	assert.Nil(t, err)

	eventInfo.Data = []string{}
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	eventInfo.Tx = nil
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)
}

func Test_SetEventState(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	err := EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	err = EventManagerV1.SetEventState(string(EventKey(event.EventName, event.ContractName, event.ChainRid)),
		true, "success")
	assert.Nil(t, err)
}
