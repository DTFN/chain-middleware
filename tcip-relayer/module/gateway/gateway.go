/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	"github.com/emirpasic/gods/maps/treemap"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"github.com/gogo/protobuf/proto"
	"github.com/robfig/cron"
)

// GatewayV1 网关管理模块全局变量
var GatewayV1 *GatewayManager

//GatewayManager 网关管理模块结构体
type GatewayManager struct {
	gatewayInfo  *treemap.Map
	gateWayState *treemap.Map
	log          *zap.SugaredLogger
}

// GatewayConnState gateway的连接状态
type GatewayConnState struct {
	Gateway bool
	Chain   bool
}

// InitGatewayManager 网关管理初始化
//  @return error
func InitGatewayManager() error {
	GatewayV1 = &GatewayManager{
		gatewayInfo:  treemap.NewWithStringComparator(),
		gateWayState: treemap.NewWithStringComparator(),
		log:          logger.GetLogger(logger.ModuleGateway),
	}
	GatewayV1.log.Debug("[InitGatewayManager] init")
	_, _ = GatewayV1.GetGatewayInfoByRange("0", "99999999999")
	GatewayV1.pingPong()
	go func() {
		// 主要是为了将链上的gatewayId缓存到本地
		_, _ = GatewayV1.GetGatewayInfoByRange("0", "99999999999")
		cron2 := cron.New()

		GatewayV1.log.Infof("[InitGatewayManager] start pingPong")
		// 每五秒执行一次心跳
		err := cron2.AddFunc("@every 5s", GatewayV1.pingPong)
		if err != nil {
			fmt.Println(err)
		}
		cron2.Start()
	}()
	return nil
}

// GatewayRegister 网关注册
//  @receiver g
//  @param gatewayInfo
//  @return string
//  @return error
func (g *GatewayManager) GatewayRegister(gatewayInfo *common.GatewayInfo) (string, error) {
	err := checkGateway(gatewayInfo)
	g.log.Infof("[GatewayRegister] check gateway success: name: %s", gatewayInfo.GetGatewayName())
	if err != nil {
		msg := fmt.Sprintf("[GatewayRegister] gateway verify error: %s", err.Error())
		g.log.Error(msg)
		return "", errors.New(msg)
	}
	kv := make(map[string][]byte)
	kv[syscontract.SaveGateway_GATEWAY_INFO_BYTE.String()], _ = proto.Marshal(gatewayInfo)
	kvJsonStr, _ := json.Marshal(kv)

	gatewayIdByte, err := relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.SaveGateway, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GatewayRegister]: [%s] invoke contract error: %s", utils.SaveGateway, err.Error())
		return "", err
	}
	gatewayId := string(gatewayIdByte)
	g.log.Infof("[GatewayRegister] call cross chain contract method %s success, gatewayId: %s",
		utils.SaveGateway, gatewayId)

	gatewayInfo.GatewayId = gatewayId
	kv[syscontract.UpdateGateway_GATEWAY_ID.String()] = gatewayIdByte
	kv[syscontract.UpdateGateway_GATEWAY_INFO_BYTE.String()], _ = proto.Marshal(gatewayInfo)
	kvJsonStr, _ = json.Marshal(kv)
	_, err = relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.UpadteGateway, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GatewayRegister]: [%s] invoke contract error: %s", utils.UpadteGateway, err.Error())
		return "", err
	}
	g.gatewayInfo.Put(gatewayId, gatewayInfo)
	g.log.Infof("[GatewayRegister] call cross chain contract method %s success, gatewayId: %s",
		utils.UpadteGateway, gatewayId)
	g.log.Infof("[GatewayRegister] gateway register success, gatewayId: %s", gatewayId)
	return gatewayId, nil
}

// GatewayUpdate 网关更新
//  @receiver g
//  @param gatewayInfo
//  @return string
//  @return error
func (g *GatewayManager) GatewayUpdate(gatewayInfo *common.GatewayInfo) (string, error) {
	err := checkGateway(gatewayInfo)
	if err != nil {
		msg := fmt.Sprintf("[GatewayUpdate] gateway verify error: %s", err.Error())
		g.log.Error(msg)
		return "", errors.New(msg)
	}
	kv := make(map[string][]byte)
	kv[syscontract.UpdateGateway_GATEWAY_ID.String()] = []byte(gatewayInfo.GatewayId)
	kv[syscontract.UpdateGateway_GATEWAY_INFO_BYTE.String()], _ = proto.Marshal(gatewayInfo)
	kvJsonStr, _ := json.Marshal(kv)
	_, err = relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.UpadteGateway, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GatewayUpdate]: [%s] invoke contract error: %s", utils.UpadteGateway, err.Error())
		return "", err
	}
	g.gatewayInfo.Put(gatewayInfo.GatewayId, gatewayInfo)
	return gatewayInfo.GatewayId, nil
}

// GetGatewayInfo 获取素有网关ID
//  @return []string
func (g *GatewayManager) GetAllGatewayId() []string {
	keys := make([]string, 0, g.gatewayInfo.Size())
    it := g.gatewayInfo.Iterator()
    for it.Next() {
        // 假设 key 是 string 类型，进行类型断言
        if key, ok := it.Key().(string); ok {
            keys = append(keys, key)
        }
    }

	return keys
}

// GetGatewayInfo 获取网关
//  @receiver g
//  @param gatewayId
//  @return *common.GatewayInfo
//  @return error
func (g *GatewayManager) GetGatewayInfo(gatewayId string) (*common.GatewayInfo, error) {
	gatewayInfo, ok1 := g.gatewayInfo.Get(gatewayId)
	if ok1 {
		return gatewayInfo.(*common.GatewayInfo), nil
	}
	kv := make(map[string][]byte)
	kv[syscontract.GetGateway_GATEWAY_ID.String()] = []byte(gatewayId)
	kvJsonStr, _ := json.Marshal(kv)
	gatewayInfoByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetGateway, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GetGatewayInfo]: [%s] invoke contract error: %s", utils.GetGateway, err.Error())
		return nil, err
	}
	if gatewayInfoByte == nil {
		msg := fmt.Sprintf("[GetGatewayInfo]: no such gateway id: %s",
			gatewayId)
		g.log.Error(msg)
		return nil, errors.New(msg)
	}
	//var gatewayInfo common.GatewayInfo
	var chainGatewayInfo common.GatewayInfo
	err = proto.Unmarshal(gatewayInfoByte, &chainGatewayInfo)
	if err != nil {
		g.log.Errorf("[GetGatewayInfo]: [%s] unmarshal gateway info error: %s", utils.GetGateway, err.Error())
		return nil, err
	}
	g.gatewayInfo.Put(chainGatewayInfo.GatewayId, &chainGatewayInfo)
	return &chainGatewayInfo, nil
}

// GetGatewayInfoByRange 根据范围获取网关
//  @receiver g
//  @param startGatewayId
//  @param stopGatewayId
//  @return []*common.GatewayInfo
//  @return error
func (g *GatewayManager) GetGatewayInfoByRange(startGatewayId, stopGatewayId string) ([]*common.GatewayInfo, error) {
	kv := make(map[string][]byte)
	kv[syscontract.GetGatewayByRange_START_GATEWAY_ID.String()] = []byte(startGatewayId)
	kv[syscontract.GetGatewayByRange_STOP_GATEWAY_ID.String()] = []byte(stopGatewayId)
	kvJsonStr, _ := json.Marshal(kv)
	gatewayInfoListByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetGatewayByRange, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GetGatewayInfoByRange]: [%s] invoke contract error: %s",
			utils.GetGatewayByRange, err.Error())
		return nil, err
	}
	var gatewayInfoByteList [][]byte
	err = json.Unmarshal(gatewayInfoListByte, &gatewayInfoByteList)
	if err != nil {
		g.log.Errorf("[GetGatewayInfoByRange]: [%s] unmarshal gateway info list byte error: %s",
			utils.GetGatewayByRange, err.Error())
		return nil, err
	}
	var gatewayInfoList []*common.GatewayInfo
	for _, v := range gatewayInfoByteList {
		var gatewayInfo common.GatewayInfo
		err = proto.Unmarshal(v, &gatewayInfo)
		if err != nil {
			g.log.Errorf("[GetGatewayInfoByRange]: [%s] unmarshal gateway info error: %s",
				utils.GetGatewayByRange, err.Error())
			return nil, err
		}
		gatewayInfoList = append(gatewayInfoList, &gatewayInfo)
		g.gatewayInfo.Put(gatewayInfo.GatewayId, &gatewayInfo)
	}
	return gatewayInfoList, nil
}

// GetGatewayNum 获取网关数量
//  @receiver g
//  @return uint64
//  @return error
func (g *GatewayManager) GetGatewayNum() (uint64, error) {
	kv := make(map[string][]byte)
	kvJsonStr, _ := json.Marshal(kv)
	gatewayNumByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetGatewayNum, true, string(kvJsonStr), -1)
	if err != nil {
		g.log.Errorf("[GetGatewayNum]: [%s] invoke contract error: %s", utils.GetGatewayNum, err.Error())
		return 0, err
	}
	gatewayNum, err := strconv.Atoi(string(gatewayNumByte))
	if err != nil {
		g.log.Errorf("[GetGatewayNum]: [%s] invoke contract error: %s", utils.GetGatewayNum, err.Error())
		return 0, err
	}
	return uint64(gatewayNum), nil
}

// CheckGatewayConnect 检查gateway的连通性
//  @receiver g
//  @param gatewayId
//  @return bool
//  @return common.TxResultValue
func (g *GatewayManager) CheckGatewayConnect(gatewayId string) (bool, common.TxResultValue) {
	gatewayConnStateItfc, ok := g.gateWayState.Get(gatewayId)
	if !ok {
		return false, common.TxResultValue_GATEWAY_PINGPONG_ERROR
	}
	gatewayConnState, _ := gatewayConnStateItfc.(GatewayConnState)
	g.log.Debugf("[CheckGatewayConnect] gatewayConnState: %+v", gatewayConnState)
	if gatewayConnState.Gateway && gatewayConnState.Chain {
		return true, common.TxResultValue(math.MaxInt16)
	}
	if gatewayConnState.Gateway && !gatewayConnState.Chain {
		return false, common.TxResultValue_CHAIN_PING_ERROR
	}
	return false, common.TxResultValue_GATEWAY_PINGPONG_ERROR
}

// pingPong 检查网关的连通性
//  @receiver g
func (g *GatewayManager) pingPong() {
	for _, gatewayInfoItfe := range g.gatewayInfo.Values() {
		gatewayInfo, _ := gatewayInfoItfe.(*common.GatewayInfo)
		gatewayConnState := GatewayConnState{
			Gateway: false,
			Chain:   false,
		}
		res, err := request.RequestV1.PingPong(-1, gatewayInfo)
		if err != nil {
			msg := fmt.Sprintf("[pingPong] gateway error: id:%s error:%s", gatewayInfo.GatewayId, err.Error())
			g.log.Error(msg)
			g.gateWayState.Put(gatewayInfo.GatewayId, gatewayConnState)
			continue
		}
		gatewayConnState.Gateway = true
		gatewayConnState.Chain = res.ChainOk
		g.log.Debugf("[pingPong] success id:%s", gatewayInfo.GatewayId)
		g.gateWayState.Put(gatewayInfo.GatewayId, gatewayConnState)
	}
}

// checkGateway 检查网关的合法性
//  @param gatewayInfo
//  @return error
func checkGateway(gatewayInfo *common.GatewayInfo) error {
	var typeInfo = reflect.TypeOf(*gatewayInfo)
	var valInfo = reflect.ValueOf(*gatewayInfo)
	num := typeInfo.NumField()
	errMsg := ""
	for i := 0; i < num; i++ {
		key := typeInfo.Field(i).Name
		val := valInfo.Field(i).Interface()
		if fmt.Sprintf("%T", val) == "string" {
			if val == "" && key != "GatewayId" {
				errMsg += fmt.Sprintf("%s/", key)
			}
		}
	}
	if errMsg != "" {
		errMsg = errMsg[:len(errMsg)-1]
		errMsg += " can't be empty"
		return errors.New(errMsg)
	}
	// 这里不需要检查服务是否存在
	//err := request.RequestV1.PingPong(-1, gatewayInfo)
	//if err != nil {
	//	errMsg = fmt.Sprintf("gateway pingpong error: name:%s error:%s", gatewayInfo.GatewayName, err.Error())
	//	return errors.New(errMsg)
	//}
	return nil
}
