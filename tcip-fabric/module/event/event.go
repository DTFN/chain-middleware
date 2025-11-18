/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	//tbis_event "chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"
	chain_config "chainmaker.org/chainmaker/tcip-fabric/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"

	"github.com/Knetic/govaluate"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"github.com/emirpasic/gods/maps/treemap"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/db"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
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
	e.log.Debugf("[BuildCrossChainMsg] eventInfo: %+v", eventInfo)
	var eventData map[string]interface{}
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
			e.log.Errorf("not support")
			return nil, nil
		}
		eventData, err = e.eventDataParse(
			event.FabricEventDataType,
			event.FabricArrayEventDataType,
			event.FabricMapEventDataType,
			[]byte(eventInfo.Data[0]))
		e.log.Debugf("BuildCrossChainMsg, eventData: %+v", eventData)
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

// buildBeginCrossChainRequestFromEvent 从合约事件中创建跨链请求
//  @receiver e
//  @param event
//  @param eventInfo
//  @return *relay_chain.BeginCrossChainRequest
//  @return error
func (e *EventManager) buildBeginCrossChainRequestFromEvent(
	event common.CrossChainEvent,
	eventInfo *EventInfo) (*relay_chain.BeginCrossChainRequest, error) {
	eventData, err := base64.StdEncoding.DecodeString(eventInfo.Data[0])
	if err != nil {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This topic is not a "+
			"cross chain event: topic %s, need base64: %s", eventInfo.Topic, eventInfo.Data[0])
		e.log.Warnf(msg)
		return nil, errors.New(msg)
	}
	var beginCrossChainRequest relay_chain.BeginCrossChainRequest
	err = proto.Unmarshal(eventData, &beginCrossChainRequest)
	// 反序列化请求失败，表明这个不是一个跨链事件，直接忽略
	if err != nil {
		msg := fmt.Sprintf("[buildBeginCrossChainRequestFromEvent] This topic is not a "+
			"cross chain event: topic %s, error %s", eventInfo.Topic, err.Error())
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

// buildBeginCrossChainRequest 从事件配置中创建跨链请求
//  @receiver e
//  @return unc
func (e *EventManager) buildBeginCrossChainRequest(
	event common.CrossChainEvent,
	eventData map[string]interface{},
	eventInfo *EventInfo) (*relay_chain.BeginCrossChainRequest, error) {
	crossChainMsgArray := make([]*common.CrossChainMsg, len(event.CrossChainCreate))
	for i, crossChainCreate := range event.CrossChainCreate {
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
	vcStr, _ :=  eventData[eventDataIndexKey[0]].(string)
	e.log.Debugf("vcStr %s", vcStr)

	// 解析VC
	gatewayId, contractType, chainRid, contractName, method := e.parseVc(vcStr)

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
//  @param eventDataType
//  @param arrayDataType
//  @param mapDataType
//  @param eventData
//  @return map[string]interface{}
//  @return error
func (e *EventManager) eventDataParse(eventDataType common.EventDataType,
	arrayDataType []common.EventDataType, mapDataType []*common.MapEventDataType,
	eventData []byte) (map[string]interface{}, error) {
	e.log.Debugf("eventDataParse, arrayDataType: %+v, eventDataStr: %s", arrayDataType, string(eventData))
	var (
		eventDataMap = make(map[string]interface{})
		err          error
	)
	if eventDataType == common.EventDataType_STRING {
		eventDataMap[eventDataIndexKey[0]] = string(eventData)
	} else if eventDataType == common.EventDataType_MAP {
		mapData := make(map[string]interface{})
		err = json.Unmarshal(eventData, &mapData)
		if err != nil {
			return nil, fmt.Errorf("Map parse error, only support "+
				"map[string]interface{}: data %s error %s ",
				eventData, err.Error())
		}
		for index, dataType := range mapDataType {
			eventDataMap[eventDataIndexKey[index]] = e.parseEventToMap(dataType.DataType,
				-1, dataType.Key, mapData, nil)
		}
	} else if eventDataType == common.EventDataType_ARRAY {
		arrayData := make([]interface{}, 0)
		err = json.Unmarshal(eventData, &arrayData)
		if err != nil {
			return nil, fmt.Errorf("Array parse error, only support "+
				"[]interface{}: data %s error %s ",
				eventData, err.Error())
		}
		for index, dataType := range arrayDataType {
			eventDataMap[eventDataIndexKey[index]] = e.parseEventToMap(dataType,
				index, "", nil, arrayData)
		}
	} else {
		e.log.Error("[eventDataParse] Unsupported type: %s", eventDataType.String())
		return nil, fmt.Errorf("Unsupported data type: %s ", eventDataType)
	}
	return eventDataMap, nil
}

// checkIsCrossChain 跨链过滤检测
//  @receiver e
//  @param isCrossChain
//  @param eventData
//  @return bool
func (e *EventManager) checkIsCrossChain(isCrossChain string, eventData map[string]interface{}) bool {
	expression, _ := govaluate.NewEvaluableExpression(isCrossChain)

	res, err := expression.Evaluate(eventData)
	if err != nil {
		e.log.Warnf("[checkIsCrossChain] Evaluate error: %s, isCrossChain: %s",
			err.Error(), isCrossChain)
		return false
	}
	return res.(bool)
}

// parseEventToMap 不管怎么样，最后入参的时候一定是[]byte类型的，所以只要处理byte和string就可以了
//  @receiver e
//  @param dataType
//  @param index
//  @param key
//  @param mapData
//  @param arrData
//  @return interface{}
func (e *EventManager) parseEventToMap(dataType common.EventDataType,
	index int, key string, mapData map[string]interface{}, arrData []interface{}) interface{} {
	// 收集类型错误
	defer func() {
		err := recover()
		if err != nil {
			e.log.Error("[parseEventToMap] error occurred: %v", err)
		}
	}()
	//var err error
	switch dataType {
	case common.EventDataType_BYTE:
		if key != "" {
			return mapData[key].([]byte)
		}
		return arrData[index].([]byte)
	case common.EventDataType_INT:
		//tmpNum := 0
		//if key != "" {
		//	tmpNum, err = strconv.Atoi(string(mapData[key].([]byte)))
		//	if err != nil {
		//		e.log.Errorf("[parseEventToMap]Int parse error: error %s", err.Error())
		//		return nil
		//	}
		//} else {
		//	tmpNum, err = strconv.Atoi(string(arrData[index].([]byte)))
		//	if err != nil {
		//		e.log.Errorf("[parseEventToMap]Int parse error: error %s", err.Error())
		//		return nil
		//	}
		//}
		//return tmpNum
	case common.EventDataType_FLOAT:
		//tmpNum := 0.1
		//if key != "" {
		//	tmpNum, err = strconv.ParseFloat(string(mapData[key].([]byte)), 64)
		//	if err != nil {
		//		e.log.Errorf("[parseEventToMap]Float parse error: error %s", err.Error())
		//		return nil
		//	}
		//} else {
		//	tmpNum, err = strconv.ParseFloat(string(arrData[index].([]byte)), 64)
		//	if err != nil {
		//		e.log.Errorf("[parseEventToMap]Float parse error: error %s", err.Error())
		//		return nil
		//	}
		//}
		//return tmpNum
	case common.EventDataType_BOOL:
		//var (
		//	res interface{}
		//	err error
		//)
		//if key != "" {
		//	res, err = strconv.ParseBool(mapData[key].(string))
		//} else {
		//	res, err = strconv.ParseBool(arrData[index].(string))
		//}
		//if err != nil {
		//	e.log.Errorf("[parseEventToMap]Bool parse error: error %s", err.Error())
		//	return nil
		//}
		//return res
	case common.EventDataType_MAP:
		//var err error
		//data := make(map[string]interface{})
		//if key != "" {
		//	err = json.Unmarshal(mapData[key].([]byte), &data)
		//} else {
		//	err = json.Unmarshal(arrData[index].([]byte), &data)
		//}
		//if err != nil {
		//	e.log.Errorf("[parseEventToMap]Map parse error, only support "+
		//		"map[string]interface{}: error %s", err.Error())
		//	return nil
		//}
		//return data
	case common.EventDataType_ARRAY:
		//var err error
		//data := make([]interface{}, 0)
		//if key != "" {
		//	err = json.Unmarshal(mapData[key].([]byte), &data)
		//} else {
		//	err = json.Unmarshal(arrData[index].([]byte), &data)
		//}
		//if err != nil {
		//	e.log.Errorf("[parseEventToMap]Array parse error, only support "+
		//		"map[string]interface{}: error %s", err.Error())
		//	return nil
		//}
		//return data
	case common.EventDataType_STRING:
		if key != "" {
			return mapData[key].(string)
		}
		return arrData[index].(string)
	default:
		e.log.Errorf("[parseEventToMap]Unsupported data type: %s", dataType)
	}
	return mapData[key]
}

// checkEvent 检查event合法性
//  @param event
//  @return error
func checkEvent(event *common.CrossChainEvent) error {
	if event.EventName == "" {
		return errors.New("EventName can't be empty")
	}
	if event.ChainRid == "" {
		return errors.New("ChainId can't be empty")
	}
	if event.TriggerCondition == common.TriggerCondition_CONTRACT_EVENT {
		if event.FabricEventDataType != common.EventDataType_ARRAY &&
			event.FabricEventDataType != common.EventDataType_MAP {
			return errors.New("FabricEventDataType only support array")
		}
		if len(event.FabricArrayEventDataType) == 0 &&
			len(event.FabricMapEventDataType) == 0 {
			return errors.New("FabricArrayEventDataType or FabricMapEventDataType is required")
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
			return checkChainCreate(c, i)
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

// checkChainCreate 检查跨链事件触发器的合法性
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
	if len(c.ParamData) > maxParamLen {
		return fmt.Errorf("len(CrossChainCreate[%d].ParamData).ParamData cannot exceed %d. "+
			"If the value exceeds %d, modify the code", i, maxParamLen, maxParamLen)
	}
	err := checkConfirm(c.ConfirmInfo, fmt.Sprintf("CrossChainCreate[%d].ConfirmInfo", i), false)
	if err != nil {
		return err
	}
	err = checkCancel(c.CancelInfo, fmt.Sprintf("CrossChainCreate[%d].CancelInfo", i), false)
	if err != nil {
		return err
	}
	return err
}

// checkConfirm 检查提交参数的合法性
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
	if confirm.Parameter != "" {
		if strings.Count(confirm.Parameter, "%") != len(confirm.ParamData) {
			return fmt.Errorf("%s.Parameter and "+
				"len(%s.ParamData) not match", errorFill, errorFill)
		}
	}
	if len(confirm.ParamData) > maxParamLen {
		return fmt.Errorf("len(%s.ParamData).ParamData cannot exceed %d. "+
			"If the value exceeds %d, modify the code", errorFill, maxParamLen, maxParamLen)
	}
	return nil
}

// checkCancel 检查会滚参数的合法性
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
	if cancel.Parameter != "" {
		if strings.Count(cancel.Parameter, "%") != len(cancel.ParamData) {
			return fmt.Errorf("%s.Parameter and "+
				"len(%s.ParamData) not match", errorFill, errorFill)
		}
	}
	if len(cancel.ParamData) > maxParamLen {
		return fmt.Errorf("len(%s.ParamData).ParamData cannot exceed %d. "+
			"If the value exceeds %d, modify the code", errorFill, maxParamLen, maxParamLen)
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
