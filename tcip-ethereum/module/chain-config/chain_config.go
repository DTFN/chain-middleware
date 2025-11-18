/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_config

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/db"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/conf"
)

const (
	chainConfigKey = "chain#config"
	nilStr         = ""
	maxString      = "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
)

// ChainConfigManager 链配置管理模块全局对象
var ChainConfigManager *ChainConfig

// ChainConfig 链配置管理模块结构体
type ChainConfig struct {
	log *zap.SugaredLogger
}

// NewChainConfig 初始化链配置管理模块
func NewChainConfig() {
	ChainConfigManager = &ChainConfig{
		log: logger.GetLogger(logger.ModuleChainConfig),
	}
}

// Save 保存链配置
//
//	@receiver c
//	@param ethConfig
//	@param operate
//	@return error
func (c *ChainConfig) Save(ethConfig *conf.EthConfig, operate common.Operate) error {
	err := checkEthConfig(ethConfig)
	if err != nil {
		return err
	}
	setDefault(ethConfig)
	chainConfigOperate := &utils.ChainConfigOperate{
		ChainRid: ethConfig.ChainRid,
		Operate:  common.Operate_SAVE,
	}
	key := parseChainKey(ethConfig.ChainRid)
	bytes, err := proto.Marshal(ethConfig)
	if err != nil {
		msg := fmt.Sprintf("[Save] Marshal ethConfig error: %v, %s", ethConfig, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	has, err := db.Db.Has(key)
	if err != nil {
		msg := fmt.Sprintf("[Save] Check key error: %s, %s", key, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	if operate == common.Operate_SAVE {
		if has {
			msg := fmt.Sprintf("[Save] %s is existed",
				ethConfig.ChainRid)
			c.log.Error(msg)
			return errors.New(msg)
		}
		err = db.Db.Put(key, bytes)
		if err != nil {
			msg := fmt.Sprintf("[Save] Put ethConfig error: %s, %s",
				ethConfig.ChainRid, err.Error())
			c.log.Error(msg)
			return errors.New(msg)
		}
		c.AllChainRidsAdd(ethConfig.ChainRid)
	} else if operate == common.Operate_UPDATE {
		if !has {
			msg := fmt.Sprintf("[Save] key \"%s\" not existed, Create it first", key)
			c.log.Error(msg)
			return errors.New(msg)
		}
		err = db.Db.Put(key, bytes)
		if err != nil {
			msg := fmt.Sprintf("[Save] Put lingshuConfig error: %s, %s",
				ethConfig.ChainRid, err.Error())
			c.log.Error(msg)
			return errors.New(msg)
		}
		c.AllChainRidsAdd(ethConfig.ChainRid)
	} else {
		return fmt.Errorf("[Save] unsupported operate: %s, %s", ethConfig.ChainRid, operate)
	}
	utils.UpdateChainConfigChan <- chainConfigOperate
	return nil
}

func (c *ChainConfig) AllChainRidsAdd(chainResourceId string) {
	key := []byte("ALL_CHAIN_RID")
	chainRidsBytes, _ := db.Db.Get(key)
	var chainRids []string
	json.Unmarshal(chainRidsBytes, &chainRids)

	// 检查是否已经存在
	for _, chainRid := range chainRids {
		if chainRid == chainResourceId {
			return
		}
	}

	chainRids = append(chainRids, chainResourceId)
	jsonBytes, _ := json.Marshal(chainRids)
	db.Db.Put(key, jsonBytes)
}

func (c *ChainConfig) AllChainRidsDelete(chainResourceId string) {
	key := []byte("ALL_CHAIN_RID")
	chainRidsBytes, _ := db.Db.Get(key)
	var chainRids []string
	json.Unmarshal(chainRidsBytes, &chainRids)

	// 检查是否已经存在
	nweChainRids := make([]string, 0)
	for _, chainRid := range chainRids {
		if chainRid != chainResourceId {
			nweChainRids = append(nweChainRids, chainRid)
		}
	}

	jsonBytes, _ := json.Marshal(nweChainRids)
	db.Db.Put(key, jsonBytes)
}

func (c *ChainConfig) AllChainRids() []string {
	key := []byte("ALL_CHAIN_RID")
	chainRidsBytes, _ := db.Db.Get(key)
	var chainRids []string
	json.Unmarshal(chainRidsBytes, &chainRids)

	return chainRids
}

// Delete 删除链配置
//
//	@receiver c
//	@param chainResourceId
//	@return error
func (c *ChainConfig) Delete(chainResourceId string) error {
	key := parseChainKey(chainResourceId)
	has, err := db.Db.Has(key)
	if err != nil {
		msg := fmt.Sprintf("[Delete] Check key error: %s, %s", key, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	if !has {
		msg := fmt.Sprintf("[Delete] key \"%s\" not existed", key)
		c.log.Error(msg)
		return errors.New(msg)
	}
	err = db.Db.Delete(key)
	if err != nil {
		msg := fmt.Sprintf("[Delete] Delete key error: %s, %s", key, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	chainConfigOperate := &utils.ChainConfigOperate{
		ChainRid: chainResourceId,
		Operate:  common.Operate_UPDATE,
	}
	utils.UpdateChainConfigChan <- chainConfigOperate
	return nil
}

// Get 获取链配置
//
//	@receiver c
//	@param chainResourceId
//	@return []*common.BcosConfig
//	@return error
func (c *ChainConfig) Get(chainResourceId string) ([]*conf.EthConfig, error) {
	chainConfigs := make([]*conf.EthConfig, 0)
	if chainResourceId != nilStr {
		key := parseChainKey(chainResourceId)
		has, err := db.Db.Has(key)
		if err != nil {
			msg := fmt.Sprintf("[Get] Check key error: %s, %s", key, err.Error())
			c.log.Error(msg)
			return nil, errors.New(msg)
		}
		if !has {
			msg := fmt.Sprintf("[Get] key \"%s\" not existed", key)
			c.log.Error(msg)
			return nil, errors.New(msg)
		}
		chainConfigByte, err := db.Db.Get(key)
		if err != nil {
			msg := fmt.Sprintf("[Get] Check key error: %s, %s", key, err.Error())
			c.log.Error(msg)
			return nil, errors.New(msg)
		}
		var chainConfig conf.EthConfig
		err = proto.Unmarshal(chainConfigByte, &chainConfig)
		if err != nil {
			msg := fmt.Sprintf("[Get] Unmarshal chainConfig error: %s, %s", key, err.Error())
			c.log.Error(msg)
			return nil, errors.New(msg)
		}
		return append(chainConfigs, &chainConfig), nil
	}
	res, err := db.Db.NewIteratorWithRange(parseChainKey(nilStr), parseChainKey(maxString))
	if err != nil {
		msg := fmt.Sprintf("[Get] NewIteratorWithRange chainConfig error: %s", err.Error())
		c.log.Error(msg)
		return nil, errors.New(msg)
	}
	for res.Next() {
		chainConfigByte := res.Value()
		var chainConfig conf.EthConfig
		err := proto.Unmarshal(chainConfigByte, &chainConfig)
		if err != nil {
			msg := fmt.Sprintf("[Get] Unmarshal chainConfig error: %s, %s", res.Value(), err.Error())
			c.log.Error(msg)
			return nil, errors.New(msg)
		}
		chainConfigs = append(chainConfigs, &chainConfig)
	}
	return chainConfigs, nil
}

// SetState 设置链配置状态
//
//	@receiver c
//	@param chainConfig
//	@param state
//	@param stateMessage
//	@return error
func (c *ChainConfig) SetState(chainConfig *conf.EthConfig, state bool, stateMessage string) error {
	chainConfig.State = state
	chainConfig.StateMessage = stateMessage

	key := parseChainKey(chainConfig.ChainRid)
	bytes, err := proto.Marshal(chainConfig)
	if err != nil {
		msg := fmt.Sprintf("[SetState] Marshal fabricConfig error: %v, %s", chainConfig, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	err = db.Db.Put(key, bytes)
	if err != nil {
		msg := fmt.Sprintf("[SetState] Put fabricConfig error: %s, %s", chainConfig.ChainRid, err.Error())
		c.log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

// checkBcosConfig 检查配置的合法性
//
//	@param ethConfig
//	@return error
func checkEthConfig(ethConfig *conf.EthConfig) error {
	var typeInfo = reflect.TypeOf(*ethConfig)
	var valInfo = reflect.ValueOf(*ethConfig)
	num := typeInfo.NumField()
	errMsg := ""
	for i := 0; i < num; i++ {
		key := typeInfo.Field(i).Name
		val := valInfo.Field(i).Interface()
		if fmt.Sprintf("%T", val) == "string" {
			if val == "" && key != "StateMessage" {
				errMsg += fmt.Sprintf("%s/", key)
			}
		}
	}
	if errMsg != "" {
		errMsg = errMsg[:len(errMsg)-1]
		errMsg += " can't be empty"
		return errors.New(errMsg)
	}

	return nil
}

// setDefault 设置默认值
//
//	@param ethConfig
func setDefault(ethConfig *conf.EthConfig) {
	ethConfig.State = false
	ethConfig.StateMessage = "wait verify"
}

// parseChainKey 生成链配置的key
//
//	@param chainRid
//	@return []byte
func parseChainKey(chainRid string) []byte {
	return []byte(fmt.Sprintf("%s#%s", chainConfigKey, chainRid))
}
