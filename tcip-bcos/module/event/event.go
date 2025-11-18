/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	// "chainmaker.org/chainmaker/common/v2/evmutils/abi"
	chain_config "chainmaker.org/chainmaker/tcip-bcos/v2/module/chain-config"

	//tbis_event "chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"

	bcosabi "github.com/FISCO-BCOS/go-sdk/abi"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"

	"github.com/Knetic/govaluate"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"github.com/emirpasic/gods/maps/treemap"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/db"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

const (
	eventKey    = "event"
	maxParamLen = 11
)

// eventData index key
var eventDataIndexKey = []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten"}

// EventManager 跨链触发器结构体
type EventManager struct {
	log             *zap.SugaredLogger
	crossChainEvent *treemap.Map
}

// EventManagerV1 跨链触发器
var EventManagerV1 *EventManager

// InitEventManager 初始化跨链触发器
func InitEventManager() {
	EventManagerV1 = &EventManager{
		log:             logger.GetLogger(logger.ModuleDb),
		crossChainEvent: treemap.NewWithStringComparator(),
	}
	// 建立缓存
	_, _ = EventManagerV1.GetEvent("")
}

// SaveEvent 保存event
//  @receiver e
//  @param event
//  @param isNew
//  @return error
func (e *EventManager) SaveEvent(event *common.CrossChainEvent, isNew bool) error {
	err := checkEvent(event)
	if err != nil {
		e.log.Errorf("[SaveEvent] %s", err.Error())
		return err
	}
	event = setState(event, false, "Wait Verify")
	has, err := db.Db.Has([]byte(eventKey))
	if err != nil {
		msg := fmt.Sprintf("[SaveEvent] %s", err.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	eventData := make(map[string]*common.CrossChainEvent)
	if has {
		eventDataJson, err1 := db.Db.Get([]byte(eventKey))
		if err1 != nil {
			msg := fmt.Sprintf("[SaveEvent] Get event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}
		err1 = json.Unmarshal(eventDataJson, &eventData)
		if err1 != nil {
			msg := fmt.Sprintf("[SaveEvent] Unmarshal event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}
	}
	event.CrossChainEventId = string(EventKey(event.EventName, event.ContractName, event.ChainRid))
	if _, ok := eventData[event.CrossChainEventId]; ok && isNew {
		msg := fmt.Sprintf("[SaveEvent]New event error, id[%s] already existed", event.CrossChainEventId)
		e.log.Warn(msg)
		return errors.New(msg)
	}
	if _, ok := eventData[event.CrossChainEventId]; !ok && !isNew {
		msg := fmt.Sprintf("[SaveEvent]Update event error, id[%s] not existed", event.CrossChainEventId)
		e.log.Warn(msg)
		return errors.New(msg)
	}
	eventData[event.CrossChainEventId] = event
	eventDataByte, err := json.Marshal(eventData)
	if err != nil {
		msg := fmt.Sprintf("[SaveEvent] Marshal event list failed: %s", err.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	e.crossChainEvent.Put(eventKey, eventData)
	err = db.Db.Put([]byte(eventKey), eventDataByte)
	if err != nil {
		msg := fmt.Sprintf("[SaveEvent] Save event list failed: %s", err.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	utils.EventChan <- &utils.EventOperate{
		CrossChainEventID: event.CrossChainEventId,
		ChainRid:          event.ChainRid,
		ContractName:      event.ContractName,
		Operate:           common.Operate_SAVE,
	}
	return nil
}

// DeleteEvent 删除event
//  @receiver e
//  @param event
//  @return error
func (e *EventManager) DeleteEvent(event *common.CrossChainEvent) error {
	event.CrossChainEventId = string(EventKey(event.EventName,
		event.ContractName, event.ChainRid))
	if event.CrossChainEventId == "" {
		return errors.New("no cross chain event found")
	}
	has, err := db.Db.Has([]byte(eventKey))
	if err != nil {
		msg := fmt.Sprintf("[DeleteEvent] %s", err.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	eventData := make(map[string]*common.CrossChainEvent)
	if has {
		eventDataJson, err1 := db.Db.Get([]byte(eventKey))
		if err1 != nil {
			msg := fmt.Sprintf("[DeleteEvent] Get event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}
		err1 = json.Unmarshal(eventDataJson, &eventData)
		if err1 != nil {
			msg := fmt.Sprintf("[DeleteEvent] Unmarshal event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}
		if _, ok := eventData[event.CrossChainEventId]; !ok {
			return errors.New("no cross chain event found")
		}
		eventOperate := &utils.EventOperate{
			ChainRid:     eventData[event.CrossChainEventId].ChainRid,
			ContractName: eventData[event.CrossChainEventId].ContractName,
			Operate:      common.Operate_DELETE,
		}
		delete(eventData, event.CrossChainEventId)
		eventDataByte, err1 := json.Marshal(eventData)
		if err1 != nil {
			msg := fmt.Sprintf("[DeleteEvent] Marshal event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}
		e.crossChainEvent.Put(eventKey, eventData)
		err1 = db.Db.Put([]byte(eventKey), eventDataByte)
		if err1 != nil {
			msg := fmt.Sprintf("[DeleteEvent] Save event list failed: %s", err1.Error())
			e.log.Error(msg)
			return errors.New(msg)
		}

		utils.EventChan <- eventOperate
		return nil
	}
	return errors.New("no cross chain event found")
}

// GetEvent 获取跨链事件
//  @receiver e
//  @param crossChainEventId
//  @return []*common.CrossChainEvent
//  @return error
func (e *EventManager) GetEvent(crossChainEventId string) ([]*common.CrossChainEvent, error) {
	has, err := db.Db.Has([]byte(eventKey))
	if err != nil || !has {
		msg := fmt.Sprintf("[GetEvent] No cross chain event found: id %s", crossChainEventId)
		e.log.Error(msg)
		return nil, errors.New(msg)
	}
	eventData := make(map[string]*common.CrossChainEvent)
	eventDataJson, err := db.Db.Get([]byte(eventKey))
	if err != nil {
		msg := fmt.Sprintf("[GetEvent] Get event list failed: %s", err.Error())
		e.log.Error(msg)
		return nil, errors.New(msg)
	}
	err = json.Unmarshal(eventDataJson, &eventData)
	if err != nil {
		msg := fmt.Sprintf("[GetEvent] Unmarshal event list failed: %s", err.Error())
		e.log.Error(msg)
		return nil, errors.New(msg)
	}
	e.crossChainEvent.Put(eventKey, eventData)
	var eventList []*common.CrossChainEvent
	if crossChainEventId != "" {
		event, ok := eventData[crossChainEventId]
		if ok {
			return append(eventList, event), nil
		}
		msg := fmt.Sprintf("[GetEvent] No cross chain event found: id %s", crossChainEventId)
		e.log.Error(msg)
		return nil, errors.New(msg)
	}
	for _, e := range eventData {
		eventList = append(eventList, e)
	}
	return eventList, nil
}

// SetEventState 设置event状态
//  @receiver e
//  @param crossChainEventId
//  @param state
//  @param stateMessages
//  @return error
func (e *EventManager) SetEventState(crossChainEventId string, state bool, stateMessages string) error {
	eventInfos, err := e.GetEvent(crossChainEventId)
	if err != nil {
		e.log.Error(err.Error())
		return err
	}
	eventInfo := eventInfos[0]
	eventInfo = setState(eventInfo, state, stateMessages)
	eventData := make(map[string]*common.CrossChainEvent)
	eventDataJson, err1 := db.Db.Get([]byte(eventKey))
	if err1 != nil {
		msg := fmt.Sprintf("[SetEventState] Get event list failed: %s", err1.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	err1 = json.Unmarshal(eventDataJson, &eventData)
	if err1 != nil {
		msg := fmt.Sprintf("[SetEventState] Unmarshal event list failed: %s", err1.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	eventData[eventInfo.CrossChainEventId] = eventInfo

	eventDataByte, err1 := json.Marshal(eventData)
	if err1 != nil {
		msg := fmt.Sprintf("[SetEventState] Marshal event list failed: %s", err1.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	e.crossChainEvent.Put(eventKey, eventData)
	err1 = db.Db.Put([]byte(eventKey), eventDataByte)
	if err1 != nil {
		msg := fmt.Sprintf("[SetEventState] Save event list failed: %s", err1.Error())
		e.log.Error(msg)
		return errors.New(msg)
	}
	return nil
}

// BuildCrossChainMsg 创建跨链信息
//  @receiver e
//  @param eventInfo
//  @return req
//  @return err
func (e *EventManager) BuildCrossChainMsg(eventInfo *EventInfo) (req *relay_chain.BeginCrossChainRequest, err error) {
	var eventData map[string]interface{}
	e.log.Debugf("into BuildCrossChainMsg")
	eventListItfc, ok := e.crossChainEvent.Get(eventKey)
	if !ok {
		_, err = e.GetEvent("")
		if err != nil {
			e.log.Warn(err.Error())
			return nil, err
		}
		eventListItfc, ok = e.crossChainEvent.Get(eventKey)
		if !ok {
			e.log.Warnf("[BuildCrossChainMsg]no cross chain event found, unsupported this eventInfo: %v", eventInfo)
			return nil, nil
		}
	}
	eventMap, _ := eventListItfc.(map[string]*common.CrossChainEvent)

	for _, eventSrc := range eventMap {
		event := common.CrossChainEvent{}
		_ = utils.DeepCopy(&event, eventSrc)
		if event.EventName != eventInfo.Topic ||
			event.ChainRid != eventInfo.ChainRid ||
			event.ContractName != eventInfo.ContractName {
			continue
		}
		if event.TriggerCondition == common.TriggerCondition_COMPLETELY_CONTRACT_EVENT {
			return e.buildBeginCrossChainRequestFromEvent(event, eventInfo)
		} else if event.TriggerCondition == common.TriggerCondition_TBIS_CONTRACT {
			//chainMakerEventInfo := &tbis_event.EventInfo{
			//	Topic:        eventInfo.Topic,
			//	ChainRid:     eventInfo.ChainRid,
			//	ContractName: eventInfo.ContractName,
			//	TxProve:      eventInfo.TxProve,
			//	Data:         eventInfo.Data,
			//	Tx:           eventInfo.Tx,
			//	TxId:         eventInfo.TxId,
			//	BlockHeight:  eventInfo.BlockHeight,
			//}
			//return tbis_event.BuildBeginCrossChainRequestFromTbis(event, chainMakerEventInfo,
			//	e.log, conf.Config.BaseConfig.GatewayID)
			e.log.Error("not support")
		}
		eventData, err = e.eventDataParse(event.BcosEventDataType, eventInfo.Data)
		if err != nil {
			msg := fmt.Sprintf("[BuildCrossChainMsg] eventDataParse error: %s, topic: %s, crossChainEventId %s",
				err.Error(), eventInfo.Topic, event.CrossChainEventId)
			e.log.Warnf(msg)
			continue
		}
		if !e.checkIsCrossChain(event.IsCrossChain, eventData) {
			continue
		}
		// 普通跨链交易走这里
		req, err = e.buildBeginCrossChainRequest(event, eventData, eventInfo)
		if err != nil {
			msg := fmt.Sprintf("[BuildCrossChainMsg] buildBeginCrossChainRequest error: %s, topic: %s, "+
				"crossChainEventId %s", err.Error(), eventInfo.Topic, event.CrossChainEventId)
			e.log.Warnf(msg)
			continue
		}
		return req, nil
	}
	e.log.Infof("[BuildCrossChainMsg] This topic is not a cross chain event: topic %s", eventInfo.Topic)
	return nil, nil
}

// buildBeginCrossChainRequestFromEvent 创建跨链请求，从事件配置
//  @receiver e
//  @param event
//  @param eventInfo
//  @return *relay_chain.BeginCrossChainRequest
//  @return error
func (e *EventManager) buildBeginCrossChainRequestFromEvent(
	event common.CrossChainEvent,
	eventInfo *EventInfo) (*relay_chain.BeginCrossChainRequest, error) {
	if len(eventInfo.Data) < 1 {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This topic is not a "+
			"cross chain event: topic %s", eventInfo.Topic)
		e.log.Warnf(msg)
		return nil, errors.New(msg)
	}
	eventData, err := base64.StdEncoding.DecodeString(eventInfo.Data[1])
	if err != nil {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This topic is not a "+
			"cross chain event: topic %s, need base64: %s", eventInfo.Topic, eventInfo.Data[1])
		e.log.Warnf(msg)
		return nil, errors.New(msg)
	}
	var beginCrossChainRequest relay_chain.BeginCrossChainRequest
	err = proto.Unmarshal([]byte(eventData), &beginCrossChainRequest)
	// 反序列化请求失败，表明这个不是一个跨链事件，直接忽略
	if err != nil {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This topic is not a "+
			"cross chain event: topic %s, error %s", eventInfo.Topic, err.Error())
		e.log.Warnf(msg)
		return nil, errors.New(msg)
	}
	// 如果tx是空的，表明接收到事件的时候，没有正确的获取到tx内容，需要检查一下具体的问题
	if eventInfo.Tx == nil {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This tx is nil, event %+v", eventInfo)
		e.log.Warnf(msg)
		return nil, errors.New(msg)
	}
	beginCrossChainRequest.TxContent = &common.TxContent{
		TxId:        eventInfo.TxId,
		Tx:          eventInfo.Tx,
		TxResult:    common.TxResultValue_TX_SUCCESS,
		GatewayId:   conf.Config.BaseConfig.GatewayID,
		ChainRid:    eventInfo.ChainRid,
		TxProve:     eventInfo.TxProve,
		BlockHeight: eventInfo.BlockHeight,
	}
	beginCrossChainRequest.From = conf.Config.BaseConfig.GatewayID
	beginCrossChainRequest.Timeout = event.Timeout
	beginCrossChainRequest.CrossType = event.CrossType
	return &beginCrossChainRequest, nil
}

// buildBeginCrossChainRequest 创建跨链请求，从事件内容
//  @receiver e
//  @return unc
func (e *EventManager) buildBeginCrossChainRequest(
	event common.CrossChainEvent,
	eventData map[string]interface{},
	eventInfo *EventInfo) (*relay_chain.BeginCrossChainRequest, error) {
	e.log.Debugf("buildBeginCrossChainRequest eventData: %+v", eventData)
	crossChainMsgArray := make([]*common.CrossChainMsg, len(event.CrossChainCreate))
	for i, crossChainCreate := range event.CrossChainCreate {
		e.log.Debugf("crossChainCreate.ContractName %s", crossChainCreate.ContractName)
			crossChainCreate.ConfirmInfo = e.buildConfirm(crossChainCreate.ConfirmInfo, eventData)
			crossChainCreate.CancelInfo = e.buildCancel(crossChainCreate.CancelInfo, eventData)

		if crossChainCreate.ContractName == "$VC_CONTRACT" {
			// 解析VC,并编码
			e.parseVcAndEncodeEvmParam(crossChainCreate, eventData)
		} else {
			if crossChainCreate.Parameter != "" {
				paramData := make([]interface{}, len(crossChainCreate.ParamData))
				for j, v := range crossChainCreate.ParamData {
					paramData[j] = eventData[eventDataIndexKey[v]]
				}
				if len(paramData) > 0 {
					crossChainCreate.Parameter = fmt.Sprintf(crossChainCreate.Parameter,
						paramData...)
				}
			}
		}
		crossChainMsgArray[i] = crossChainCreate
	}
	return &relay_chain.BeginCrossChainRequest{
		Version:        common.Version_V1_0_0,
		CrossChainName: event.CrossChainName,
		CrossChainFlag: event.CrossChainFlag,
		CrossChainMsg:  crossChainMsgArray,
		TxContent: &common.TxContent{
			TxId:        eventInfo.TxId,
			Tx:          eventInfo.Tx,
			TxResult:    common.TxResultValue_TX_SUCCESS,
			GatewayId:   conf.Config.BaseConfig.GatewayID,
			ChainRid:    eventInfo.ChainRid,
			TxProve:     eventInfo.TxProve,
			BlockHeight: eventInfo.BlockHeight,
		},
		From:        conf.Config.BaseConfig.GatewayID,
		Timeout:     event.Timeout,
		ConfirmInfo: e.buildConfirm(event.ConfirmInfo, eventData),
		CancelInfo:  e.buildCancel(event.CancelInfo, eventData),
		CrossType:   event.CrossType,
	}, nil
}

func (e *EventManager) parseVcAndEncodeEvmParam(
	crossChainMsg *common.CrossChainMsg,
	eventData map[string]interface{},
) {
	e.log.Debugf("parseVcAndEncodeEvmParam %+v", eventData)

	// 提取VC字符串（VC字符串是合约事件抛出的第一个参数,hex编码后的string）
	vcStr, _ :=  eventData[eventDataIndexKey[1]].(string)
	e.log.Debugf("vcStr %s", vcStr)

	// 解析VC
	gatewayId, contractType, chainRid, contractName, method := e.parseVc(vcStr)

	// 参数赋值
	// var gatewayId = "1"
	// var chainRid = "chainmaker002"
	// var contractName = "kvStore"
	// var method = "put"
	// var parameterStr = "[{\"string\": \"aabb\"},{\"string\": \"{}\"}]"
	// paramData := []map[string]string{
	// 	{"string": param.DID},
	// 	{"string": param.DIDDocument},
	// }
	// parameterStr, err := json.Marshal(paramData)

	// todo EVM入参编码
	// 构造input
	// {
	// 	"name": "didDocument",
	// 	"type": "string",
	// },

	dataVcStr, _ := json.Marshal(map[string]interface{}{
		"data": vcStr,
	})

	// 填入属性
	crossChainMsg.GatewayId = gatewayId
	crossChainMsg.ChainRid = chainRid
	crossChainMsg.ContractName = contractName
	crossChainMsg.Method = method
	crossChainMsg.Parameter = string(dataVcStr)
	crossChainMsg.ExtraData = "$VC_CROSS"
	e.log.Debugf("parseVcAndEncodeEvmParam gatewayId: %s, contractType: %s, chainRid: %s, contractName: %s, method: %s, encodedParameter: %s", gatewayId, contractType, chainRid, contractName, method, dataVcStr)
}

func (e *EventManager) parseVc(vcJsonStr string) (gatewayId string, contractType string, chainRid string, contractName string, method string) {
	// vcJsonStrT := `{
	// 	"credentialSubject": {
	// 		"gateway_id": "1",
	// 		"chain_rid": "chainmaker002",
	// 		"contract_address": "kvStore",
	// 		"contract_func": "put",
	// 		"func_params": "[{\"string\": \"demo1\"},{\"string\": \"{}\"}]"
	// 	}
	// }`
	// s := "{\"credentialSubject\":{\"gateway_id\":\"1\",\"chain_rid\":\"chainmaker002\",\"contract_address\":\"kvStore\",\"contract_func\":\"put\",\"func_params\":\"[{\\\"string\\\": \\\"demo1\\\"},{\\\"string\\\": \\\"{}\\\"}]\"}}"

	var result map[string]interface{}
	json.Unmarshal([]byte(vcJsonStr), &result)
	e.log.Debugf("vcJsonStr %s, result %+v", vcJsonStr, result)

	var targetVc map[string]interface{} = result["target"].(map[string]interface{})
	e.log.Debugf("vcJsonStr %s, targetVc %+v", vcJsonStr, targetVc)

	credentialSubject, _ := targetVc["credentialSubject"].(map[string]interface{})
	e.log.Debugf("credentialSubject %+v", credentialSubject)

	return credentialSubject["gateway_id"].(string), 
	credentialSubject["contract_type"].(string), 
	credentialSubject["chain_rid"].(string), 
	credentialSubject["contract_address"].(string), 
	credentialSubject["contract_func"].(string)
}

func (r *EventManager) parseOneString(rawHexStr string) (string, error) {
	// 移除可能存在的0x前缀
	hexStr := strings.TrimPrefix(rawHexStr, "0x")
	
	// 转换为字节数组
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return "", errors.New("decode fail")
	}

	r.log.Debugf("[didManagerEvent]1 %+v", data)
	didIdentifierOffset := bytesToUint256(data[0:32])
	r.log.Debugf("[didManagerEvent]2 %s %+v", didIdentifierOffset, data[0:32])

	didIdentifierLen := bytesToUint256(data[didIdentifierOffset:didIdentifierOffset + 32])
	r.log.Debugf("[didManagerEvent]3")
	didIdentifier := string(data[didIdentifierOffset + 32 : didIdentifierOffset + 32 + didIdentifierLen])

	return didIdentifier, err
}

// 将32字节转换为uint64整数,只取后8位
func bytesToUint256(b []byte) uint64 {
	bs := b[24: 32]
	var result uint64 = binary.BigEndian.Uint64(bs)
	return result
}

// buildConfirm 创建提交信息
//  @receiver e
//  @param confirmInfo
//  @param eventData
//  @return *common.ConfirmInfo
func (e *EventManager) buildConfirm(confirmInfo *common.ConfirmInfo,
	eventData map[string]interface{}) *common.ConfirmInfo {
	if confirmInfo == nil || confirmInfo.ChainRid == "" {
		return &common.ConfirmInfo{}
	}
	if confirmInfo.Parameter != "" {
		paramData := make([]interface{}, len(confirmInfo.ParamData))
		for i, v := range confirmInfo.ParamData {
			paramData[i] = eventData[eventDataIndexKey[v]]
		}
		if len(paramData) > 0 {
			confirmInfo.Parameter = fmt.Sprintf(confirmInfo.Parameter,
				paramData...)
		}
	}
	return confirmInfo
}

// buildCancel 创建回滚信息
//  @receiver e
//  @param cancelInfo
//  @param eventData
//  @return *common.CancelInfo
func (e *EventManager) buildCancel(cancelInfo *common.CancelInfo,
	eventData map[string]interface{}) *common.CancelInfo {
	if cancelInfo == nil || cancelInfo.ChainRid == "" {
		return &common.CancelInfo{}
	}
	paramData := make([]interface{}, len(cancelInfo.ParamData))
	for i, v := range cancelInfo.ParamData {
		paramData[i] = eventData[eventDataIndexKey[v]]
	}
	if cancelInfo.Parameter != "" {
		paramData := make([]interface{}, len(cancelInfo.ParamData))
		for i, v := range cancelInfo.ParamData {
			paramData[i] = eventData[eventDataIndexKey[v]]
		}
		if len(paramData) > 0 {
			cancelInfo.Parameter = fmt.Sprintf(cancelInfo.Parameter,
				paramData...)
		}
	}
	return cancelInfo
}

// eventDataParse 处理事件参数
//  @receiver e
//  @param dataType
//  @param eventData
//  @return map[string]interface{}
//  @return error
func (e *EventManager) eventDataParse(dataType []common.EventDataType,
	eventData []string) (map[string]interface{}, error) {
	var (
		eventDataMap = make(map[string]interface{})
	)
	for index, data := range dataType {
		eventDataMap[eventDataIndexKey[index]] = eventData[index]
		switch data {
		case common.EventDataType_BYTE:
		case common.EventDataType_INT:
		case common.EventDataType_FLOAT:
		case common.EventDataType_BOOL:
		case common.EventDataType_MAP:
		case common.EventDataType_STRING:
		case common.EventDataType_ARRAY:
		case common.EventDataType_ADDRESS:
		case common.EventDataType_HASH:
		default:
			return nil, fmt.Errorf("Unsupported data type: %s ", data)
		}
	}
	return eventDataMap, nil
}

// checkIsCrossChain 判断是否为跨链请求
//  @receiver e
//  @param isCrossChain
//  @param eventData
//  @return bool
func (e *EventManager) checkIsCrossChain(isCrossChain string, eventData map[string]interface{}) bool {
	expression, _ := govaluate.NewEvaluableExpression(isCrossChain)

	e.log.Debugf("[checkIsCrossChain] eventData: %+v\n", eventData)
	res, err := expression.Evaluate(eventData)
	if err != nil {
		e.log.Warnf("[checkIsCrossChain] Evaluate error: %s, isCrossChain: %s",
			err.Error(), isCrossChain)
		return false
	}
	return res.(bool)
}

// checkEvent 判断事件配置是否合法
//  @param event
//  @return error
func checkEvent(event *common.CrossChainEvent) error {
	if event.EventName == "" {
		return errors.New("EventName can't be empty")
	}
	if event.ChainRid == "" {
		return errors.New("ChainId can't be empty")
	}
	res, err := chain_config.ChainConfigManager.Get(event.ChainRid)
	if err != nil || len(res) == 0 {
		return errors.New("ChainRid not found")
	}
	if event.ContractName == "" {
		return errors.New("ContractName can't be empty")
	}
	if event.TriggerCondition == common.TriggerCondition_CONTRACT_EVENT {
		if len(event.BcosEventDataType) == 0 {
			return errors.New("EventDataType is required")
		}
		if event.IsCrossChain == "" {
			return errors.New("IsCrossChain is required")
		}
		if len(event.CrossChainCreate) == 0 {
			return errors.New("CrossChainCreate is required")
		}
		_, err := govaluate.NewEvaluableExpression(event.IsCrossChain)
		if err != nil {
			return fmt.Errorf("IsCrossChain error: %s", err.Error())
		}
		err = checkConfirm(event.ConfirmInfo, "ConfirmInfo", true)
		if err != nil {
			return err
		}
		err = checkCancel(event.CancelInfo, "CancelInfo", true)
		if err != nil {
			return err
		}
		for i, c := range event.CrossChainCreate {
			err2 := checkChainCreate(c, i)
			if err2 != nil {
				return err2
			}
		}
		for i, a := range event.Abi {
			arg := &bcosabi.Argument{}
			err2 := arg.UnmarshalJSON([]byte(a))
			if err2 != nil {
				return fmt.Errorf("Abi[%d] abi: %s read error: %s ", i, a, err2.Error())
			}
		}
	} else if event.TriggerCondition == common.TriggerCondition_COMPLETELY_CONTRACT_EVENT {
		return nil
	} else if event.TriggerCondition == common.TriggerCondition_TBIS_CONTRACT {
		return checkTbis(event)
	} else {
		return errors.New("TriggerCondition error, need " +
			"common.TriggerCondition_CONTRACT_EVENT or common.TriggerCondition_COMPLETELY_CONTRACT_EVENT")
	}
	return nil
}

// checkTbis 检查tbis的事件触发器
//  @param event
//  @return error
func checkTbis(event *common.CrossChainEvent) error {
		if len(event.CrossChainCreate) != 1 {
			return errors.New("event.CrossChainCreate len must 1")
		}
		if event.CrossChainCreate[0].GatewayId == "" {
			return errors.New("event.CrossChainCreate[0].GatewayId can't be empty")
		}
		if event.CrossType != common.CrossType_INVOKE {
			return errors.New("event.CrossType only support 1")
	}
	return nil
}

// checkChainCreate 检查事件触发器是否合约
//  @param c
//  @param i
//  @return error
func checkChainCreate(c *common.CrossChainMsg, i int) error {
	if c.GatewayId == "" {
		return fmt.Errorf("CrossChainCreate[%d].GatewayId is required", i)
	}
	if c.ChainRid == "" {
		return fmt.Errorf("CrossChainCreate[%d].ChainId is required", i)
	}
	if c.ContractName == "" {
		return fmt.Errorf("CrossChainCreate[%d].ContractName is required", i)
	}
	if c.Method == "" {
		return fmt.Errorf("CrossChainCreate[%d].Method is required", i)
	}
	if c.Parameter != "" {
		if strings.Count(c.Parameter, "%") != len(c.ParamData) {
			return fmt.Errorf("CrossChainCreate[%d].Parameter and "+
				"len(CrossChainCreate[%d].ParamData) not match", i, i)
		}
	}
	if len(c.ParamDataType) > 0 {
		if len(c.ParamDataType) != len(c.ParamData) {
			return fmt.Errorf("CrossChainCreate[%d].Parameter and "+
				"len(CrossChainCreate[%d].ParamData) not match", i, i)
		}
		if len(c.ParamDataType) > maxParamLen || len(c.ParamData) > maxParamLen {
			return fmt.Errorf("CrossChainCreate[%d].ParamDataType and "+
				"len(CrossChainCreate[%d].ParamData) cannot exceed 10. If the value exceeds 10, modify the code", i, i)
		}

		err := checkAbi(c.Abi, "CrossChainCreate")
		if err != nil {
			return err
		}
	}

	err := checkConfirm(c.ConfirmInfo, fmt.Sprintf("CrossChainCreate[%d].ConfirmInfo", i), false)
	if err != nil {
		return err
	}
	err = checkCancel(c.CancelInfo, fmt.Sprintf("CrossChainCreate[%d].CancelInfo", i), false)
	if err != nil {
		return err
	}
	return nil
}

// checkConfirm 检查提交信息是否合法
//  @param confirm
//  @param errorFill
//  @param checkChain
//  @return error
func checkConfirm(confirm *common.ConfirmInfo, errorFill string, checkChain bool) error {
	if confirm == nil || confirm.ChainRid == "" {
		return nil
	}
	if checkChain {
		res, err := chain_config.ChainConfigManager.Get(confirm.ChainRid)
		if err != nil || len(res) == 0 {
			return fmt.Errorf("%s.ChainRid not found", errorFill)
		}
	}
	if confirm.ContractName == "" {
		return fmt.Errorf("%s.ContractName is required", errorFill)
	}
	if confirm.Method == "" {
		return fmt.Errorf("%s.Method is required", errorFill)
	}
	if len(confirm.ParamDataType) > 0 {
		if len(confirm.ParamDataType) > maxParamLen || len(confirm.ParamData) > maxParamLen {
			return fmt.Errorf("%s.ParamDataType and "+
				"len(%s.ParamData) cannot exceed 10. If the value exceeds 10, modify the code", errorFill, errorFill)
	}
	if confirm.Parameter != "" {
			if strings.Count(confirm.Parameter, "%")+
				strings.Count(confirm.Parameter, common.TryResult_TRY_RESULT.String()) !=
				len(confirm.ParamDataType) {
			return fmt.Errorf("%s.Parameter and "+
					"len(%s.ParamDataType) not match", errorFill, errorFill)
		}
	}
		return checkAbi(confirm.Abi, errorFill)
	}
	return nil
}

// checkCancel 检查回滚信息是否合法
//  @param cancel
//  @param errorFill
//  @param checkChain
//  @return error
func checkCancel(cancel *common.CancelInfo, errorFill string, checkChain bool) error {
	if cancel == nil || cancel.ChainRid == "" {
		return nil
	}
	if checkChain {
		res, err := chain_config.ChainConfigManager.Get(cancel.ChainRid)
		if err != nil || len(res) == 0 {
			return fmt.Errorf("%s.ChainRid not found", errorFill)
		}
	}
	if cancel.ContractName == "" {
		return fmt.Errorf("%s.ContractName is required", errorFill)
	}
	if cancel.Method == "" {
		return fmt.Errorf("%s.Method is required", errorFill)
	}
	if len(cancel.ParamDataType) > 0 {
		if len(cancel.ParamDataType) != len(cancel.ParamData) {
			return fmt.Errorf("%s.ParamDataType and "+
				"len(%s.ParamData) not match", errorFill, errorFill)
		}
		if len(cancel.ParamDataType) > maxParamLen || len(cancel.ParamData) > maxParamLen {
			return fmt.Errorf("%s.ParamDataType and "+
				"len(%s.ParamData) cannot exceed 10. If the value exceeds 10, modify the code", errorFill, errorFill)
	}
	if cancel.Parameter != "" {
		if strings.Count(cancel.Parameter, "%") != len(cancel.ParamData) {
			return fmt.Errorf("%s.Parameter and "+
				"len(%s.ParamData) not match", errorFill, errorFill)
		}
	}
		return checkAbi(cancel.Abi, errorFill)
	}
	return nil
}

// checkAbi 检查abi是否合法
//  @param abiStr
//  @param errorFill
//  @return error
func checkAbi(abiStr string, errorFill string) error {
	if abiStr == "" {
		return fmt.Errorf("%s.Abi abi can not be nil", errorFill)
	}
	_, err := bcosabi.JSON(strings.NewReader(abiStr))
	if err != nil {
		return fmt.Errorf("%s.Abi abi %s read error: %s", errorFill, abiStr, err.Error())
	}
	return nil
}

// EventKey 生成eventKey
//  @param eventName
//  @param contractName
//  @param chainRid
//  @return []byte
func EventKey(eventName, contractName, chainRid string) []byte {
	return []byte(fmt.Sprintf("%s#%s#%s", eventName, contractName, chainRid))
}

// setState 设置状态值
//  @param event
//  @return *common.CrossChainEvent
func setState(event *common.CrossChainEvent, state bool, stateMessage string) *common.CrossChainEvent {
	event.State = state
	event.StateMessage = stateMessage
	return event
}
