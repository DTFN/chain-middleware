/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"strconv"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/FISCO-BCOS/go-sdk/client"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/utils"

	"github.com/FISCO-BCOS/go-sdk/conf"
)

const (
	keySuffix = ".key"
	crtSuffix = ".crt"
	pemSuffix = ".pem"
)

// BCOSConfigTemplate bcos配置结构体
type BCOSConfigTemplate struct {
	CaText       string
	CaFilePath   string
	KeyText      []byte
	KeyFilePath  string
	CertFilePath string
	CertText     []byte
	AccountText  []byte
	GroupID      int
	Address      string
	Http         bool
	ChainID      int64
	SMCrypto     bool
}

// NewBCOSConfigTemplate  新建bcos的配置结构
//  @param chainID
//  @param groupIDStr
//  @param address
//  @param isSMCrypto
//  @return *BCOSConfigTemplate
//  @return error
func NewBCOSConfigTemplate(chainID int64, groupIDStr, address string, isSMCrypto bool) (*BCOSConfigTemplate, error) {
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		return nil, err
	}
	return &BCOSConfigTemplate{
		ChainID:  chainID,
		GroupID:  groupID,
		Address:  address,
		Http:     false,
		SMCrypto: isSMCrypto,
	}, nil
}

// SetCa 设置ca
//  @receiver t
//  @param caText
//  @return error
func (t *BCOSConfigTemplate) SetCa(caText string) error {
	// 将ca写入文件
	// 首先将内容写入文件
	tempFile, err := utils.WriteTempFile([]byte(caText), crtSuffix)
	if err != nil {
		return err
	}
	t.CaFilePath = tempFile.Name()
	t.CaText = caText
	return nil
}

// SetKey 设置key
//  @receiver t
//  @param keyText
//  @return error
func (t *BCOSConfigTemplate) SetKey(keyText []byte) error {
	// 首先将内容写入文件
	tempFile, err := utils.WriteTempFile(keyText, keySuffix)
	if err != nil {
		return err
	}
	t.KeyFilePath = tempFile.Name()
	t.KeyText = keyText
	return nil
}

// SetCert 设置cert
//  @receiver t
//  @param certText
//  @return error
func (t *BCOSConfigTemplate) SetCert(certText []byte) error {
	// 首先将内容写入文件
	tempFile, err := utils.WriteTempFile(certText, crtSuffix)
	if err != nil {
		return err
	}
	t.CertFilePath = tempFile.Name()
	t.CertText = certText
	return nil
}

// SetAccount 设置account
//  @receiver t
//  @param accountText
//  @return error
func (t *BCOSConfigTemplate) SetAccount(accountText []byte) error {
	tempFile, err := utils.WriteTempFile(accountText, pemSuffix)
	if err != nil {
		return err
	}
	keyBytes, _, err := conf.LoadECPrivateKeyFromPEM(tempFile.Name())
	if err != nil {
		return err
	}
	t.AccountText = keyBytes
	return nil
}

// ToConfig 生成bcos config
//  @receiver t
//  @return *conf.Config
//  @return error
func (t *BCOSConfigTemplate) ToConfig() (*conf.Config, error) {
	return &conf.Config{
		IsHTTP:     t.Http,
		ChainID:    t.ChainID,
		CAFile:     t.CaFilePath,
		Key:        t.KeyFilePath,
		Cert:       t.CertFilePath,
		IsSMCrypto: t.SMCrypto,
		PrivateKey: t.AccountText,
		GroupID:    t.GroupID,
		NodeURL:    t.Address,
	}, nil
}

// createSDK 创建bcos的sdk
//  @param bcosConfig
//  @param log
//  @return *client.Client
//  @return error
func createSDK(bcosConfig *common.BcosConfig,
	log *zap.SugaredLogger) (*client.Client, error) {
	// 此处使用的链ID
	template, err := NewBCOSConfigTemplate(bcosConfig.ChainId, bcosConfig.GroupId,
		bcosConfig.Address, bcosConfig.IsSmCrypto)
	if err != nil {
		log.Errorf("[createSDK]NewBCOSConfigTemplate error %v", err)
		return nil, err
	}
	if err = template.SetCa(bcosConfig.Ca); err != nil {
		log.Errorf("[createSDK]SetCa error %v", err)
		return nil, err
	}
	if err = template.SetAccount([]byte(bcosConfig.PrivateKey)); err != nil {
		log.Errorf("[createSDK]SetAccount error %v", err)
		return nil, err
	}
	if err = template.SetKey([]byte(bcosConfig.TlsKey)); err != nil {
		log.Errorf("[createSDK]SetKey error %v", err)
		return nil, err
	}
	if err = template.SetCert([]byte(bcosConfig.TlsCert)); err != nil {
		log.Errorf("[createSDK]SetCert error %v", err)
		return nil, err
	}
	bcosConf, err := template.ToConfig()
	if err != nil {
		log.Errorf("[createSDK]ToConfig error %v", err)
		return nil, err
	}
	log.Infof("[createSDK]Use bcos config %v", bcosConfig)
	cli, err := client.Dial(bcosConf)
	return cli, err
}
