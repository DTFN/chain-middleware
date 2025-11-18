/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"encoding/base64"
	"fmt"
	"strings"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"

	"go.uber.org/zap"
)

// TbisEventParser tbis事件解析器
type TbisEventParser struct {
	log *zap.SugaredLogger
}

// NewTbisEventParser 新建tbis事件解析器
func NewTbisEventParser(log *zap.SugaredLogger) *TbisEventParser {
	return &TbisEventParser{
		log: log,
	}
}

// ParseListenEvent 格式化事件
//  @receiver t
//  @param chainRID
//  @param data
//  @return *ContractListenData
//  @return error
func (t *TbisEventParser) ParseListenEvent(chainRID, data string) (*ContractListenData, error) {
	/**
	CrossID:${CrossID};
	EventName:${EventName};
	Business:${Business};
	CommitContractName:${CommitContractName},CommitContractVersion:${CommitContractVersion},CommitContractMethod:
	${CommitContractMethod},CommitContractParamsFormat:${CommitContractParamsFormat},CommitContractParams:
	${CommitContractParams},CommitIdentity:${CommitIdentity},CommitExtraData:${CommitExtraData};
	ChainRID:${ChainID},Identity:${Identity},ContractName:${ContractName},ContractVersion:${ContractVersion},
	ContractMethod:${ContractMethod},ContractParamsFormat:${ContractParamsFormat},ContractParams:${ContractParams},
	ContractExtraData:${ContractExtraData};
	ChainRID:${ChainID},Identity:${Identity},ContractName:${ContractName},ContractVersion:${ContractVersion},
	ContractMethod:${ContractMethod},ContractParamsFormat:${ContractParamsFormat},ContractParams:${ContractParams},
	ContractExtraData:${ContractExtraData};
	ChainRID:${ChainID},Identity:${Identity},ContractName:${ContractName},ContractVersion:${ContractVersion},
	ContractMethod:${ContractMethod},ContractParamsFormat:${ContractParamsFormat},ContractParams:${ContractParams},
	ContractExtraData:${ContractExtraData};
	*/
	// 对字符串进行处理
	if len(data) == 0 {
		return nil, fmt.Errorf("empty listen event data")
	}
	t.log.Infof("parse event: %s", data)
	eventData := replace(data)
	// 分割字符串
	eventDatas := strings.Split(eventData, groupSplit)
	// 创建空的CrossEvent，后续进行填充
	listenData := NewEmptyContractListenData()
	listenData.EventData = eventData
	// 遍历datas，处理基础配置信息
	for i := 0; i < len(eventDatas); i++ {
		eventData := eventDatas[i]
		if isEmpty(eventData) {
			continue
		}
		kvs := NewNewKeyValuesByString(eventData)
		if crossID, ok := isCrossID(kvs); ok {
			listenData.CrossID = crossID
		} else if eventName, ok := isEventName(kvs); ok {
			listenData.EventName = eventName
		} else if business, ok := isBusiness(kvs); ok {
			listenData.Business = business
		}
	}
	// 重新遍历获取对应的合约及子链操作信息
	for i := 0; i < len(eventDatas); i++ {
		eventData := eventDatas[i]
		if isEmpty(eventData) {
			continue
		}
		kvs := NewNewKeyValuesByString(eventData)
		if isCommitContract(kvs) {
			contract, err := t.loadCommitContract(chainRID, listenData.Business, kvs)
			if err != nil {
				return nil, err
			}
			listenData.CommitCtt = contract
		}
		if isChainContract(kvs) {
			crossSubEvent, err := t.loadChainContract(listenData.Business, kvs)
			if err != nil {
				return nil, err
			}
			listenData.AppendSubEvent(crossSubEvent)
		}
	}
	// 校验listenData
	if err := listenData.CheckAndInit(); err != nil {
		return nil, err
	}
	return listenData, nil
}

// loadChainContract 提取合约内容
//  @receiver t
//  @param business
//  @param keyValues
//  @return *CrossSubEvent
//  @return error
func (t *TbisEventParser) loadChainContract(business string, keyValues *KeyValues) (*CrossSubEvent, error) {
	subEvent := NewEmptyCrossSubEvent()
	if chainRID, ok := keyValues.getValue(chainRIDKey); ok {
		subEvent.ChainRID = chainRID
	}
	if identity, ok := keyValues.getValue(identityKey); ok {
		subEvent.ChainCtt.Identity = identity
	}
	if contractName, ok := keyValues.getValue(contractNameKey); ok {
		subEvent.ChainCtt.Ctt.Name = contractName
	}
	if contractVersion, ok := keyValues.getValue(contractVersionKey); ok {
		subEvent.ChainCtt.Ctt.Version = contractVersion
	}
	if contractMethod, ok := keyValues.getValue(contractMethodKey); ok {
		subEvent.ChainCtt.Ctt.Method = contractMethod
	}
	if contractExtraData, ok := keyValues.getValue(contractExtraDataKey); ok {
		// 扩展数据，需要先进行base64解码（解码后数据为map对应的json）
		extraData := selfProtocolToMap(contractExtraData)
		subEvent.ChainCtt.Ctt.Puts(extraData)
	}
	if contractParamsFormat, ok := keyValues.getValue(contractParamsFormatKey); ok {
		srcFormat, err := base64.URLEncoding.DecodeString(contractParamsFormat)
		if err != nil {
			return nil, err
		}
		subEvent.ChainCtt.ParamsFormat = string(srcFormat)
	} else {
		//if business == "" {
		//	subEvent.ChainCtt.ParamsFormat = "" // 采用空即可
		//} else {
		//	bizResource, err := t.rm.GetBizContractByChainRIDAndBizKey(subEvent.ChainRID, business)
		//	if err != nil {
		//		return nil, err
		//	}
		//	subEvent.ChainCtt.ParamsFormat = bizResource.GetBusinessConfig().ContractParamsFormat
		//}
		subEvent.ChainCtt.ParamsFormat = ""
	}
	// 允许不设置参数，则对端链不需要处理该参数
	if contractParams, ok := keyValues.getValue(contractParamsKey); ok {
		if contractParams != "" {
			// 有数据，则进行解析
			paramArray, err := convertToParamArray(contractParams, true)
			if err != nil {
				return nil, err
			}
			subEvent.ChainCtt.ParamsData = paramArray
		}
	}
	return subEvent, nil
}

// loadCommitContract 提取提交内容
//  @receiver t
//  @param chainRID
//  @param business
//  @param keyValues
//  @return *CommitContract
//  @return error
func (t *TbisEventParser) loadCommitContract(
	chainRID, business string, keyValues *KeyValues) (*CommitContract, error) {
	// 表示该组配置的Commit信息
	cct := NewEmptyCommitContract()
	if cctName, ok := keyValues.getValue(commitContractNameKey); ok {
		// bcos环境下需要设置合约名称为合约的地址信息
		cct.Ctt.Name = cctName
	}
	if cctMethod, ok := keyValues.getValue(commitContractMethodKey); ok {
		cct.Ctt.Method = cctMethod
	}
	if cctVersion, ok := keyValues.getValue(commitContractVersionKey); ok {
		cct.Ctt.Version = cctVersion
	}
	if cctIdentity, ok := keyValues.getValue(commitIdentityKey); ok {
		cct.Identity = cctIdentity
	}
	if cctVal, ok := keyValues.getValue(commitExtraDataKey); ok {
		// 扩展数据，需要先进行base64解码（解码后数据为map对应的json）
		extraData := selfProtocolToMap(cctVal)
		cct.Ctt.Puts(extraData)
	}
	// 合约参数格式必须要进行配置
	if cctParamsFormat, ok := keyValues.getValue(commitContractParamsFormatKey); ok {
		// format采用的是base64编码格式，因此需要对其进行解码操作
		srcFormat, err := base64.URLEncoding.DecodeString(cctParamsFormat)
		if err != nil {
			return nil, err
		}
		cct.ParamsFormat = string(srcFormat)
	} else {
		//if business == "" {
		//	// 允许用户不填，采用空
		//	cct.ParamsFormat = ""
		//} else {
		//	bizResource, err := m.rm.GetBizContractByChainRIDAndBizKey(chainRID, business)
		//	if err != nil {
		//		return nil, err
		//	}
		//	cct.ParamsFormat = bizResource.GetBusinessConfig().ContractParamsFormat
		//}
		cct.ParamsFormat = ""
	}
	if cctParams, ok := keyValues.getValue(commitContractParamsKey); ok {
		if cctParams != "" {
			// 有数据，则进行解析
			paramArray, err := convertToParamArray(cctParams, true)
			if err != nil {
				return nil, err
			}
			// 表明有参数
			cct.ParamsData = paramArray
		}
	}
	return cct, nil
}

// replace 替换空格和换行符
//  @param data
//  @return string
func replace(data string) string {
	// 去掉空格
	data = strings.Replace(data, " ", "", -1)
	// 去掉换行符
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)
	return data
}

// isCrossID 获取跨链ID
//  @param keyValues
//  @return string
//  @return bool
func isCrossID(keyValues *KeyValues) (string, bool) {
	return keyValues.getValue(crossIDKey)
}

// isEventName 获取event名称
//  @param keyValues
//  @return string
//  @return bool
func isEventName(keyValues *KeyValues) (string, bool) {
	return keyValues.getValue(eventNameKey)
}

// isBusiness 获取Business字段
//  @param keyValues
//  @return string
//  @return bool
func isBusiness(keyValues *KeyValues) (string, bool) {
	return keyValues.getValue(businessKey)
}

// isCommitContract 获取提交的身份签名
//  @param keyValues
//  @return bool
func isCommitContract(keyValues *KeyValues) bool {
	return keyValues.hasKey(commitIdentityKey)
}

// isChainContract 获取链资源id
//  @param keyValues
//  @return bool
func isChainContract(keyValues *KeyValues) bool {
	return keyValues.hasKey(chainRIDKey)
}

// BuildBeginCrossChainRequestFromTbis 生成跨链请求参数
//  @param event
//  @param eventInfo
//  @param log
//  @param gatewayId
//  @return req
//  @return err
func BuildBeginCrossChainRequestFromTbis(event common.CrossChainEvent, eventInfo *EventInfo,
	log *zap.SugaredLogger, gatewayId string) (req *relay_chain.BeginCrossChainRequest, err error) {
	tbisEventParser := NewTbisEventParser(log)
	var (
		tbisEvent *ContractListenData
	)
	for _, data := range eventInfo.Data {
		tbisEvent, err = tbisEventParser.ParseListenEvent(eventInfo.ChainRid, data)
		if err != nil {
			continue
		}
	}
	if err != nil {
		return nil, err
	}
	subExtraData := tbisEvent.SubEvents[0].ChainCtt.Ctt.ExtraData
	var subAbi string
	if abi, ok := subExtraData[aBI]; ok {
		subAbi = abi
	}
	subParamType := make([]common.EventDataType, 0)
	for i := 0; i < len(tbisEvent.SubEvents[0].ChainCtt.ParamsData); i++ {
		subParamType = append(subParamType, common.EventDataType_STRING)
	}
	commitParamType := make([]common.EventDataType, 0)
	for i := 0; i < len(tbisEvent.CommitCtt.ParamsData)+1; i++ {
		commitParamType = append(commitParamType, common.EventDataType_STRING)
	}
	majorExtraData := tbisEvent.CommitCtt.Ctt.ExtraData
	var majorAbi string
	if abi, ok := majorExtraData[aBI]; ok {
		majorAbi = abi
	}
	log.Debugf("[BuildBeginCrossChainRequestFromTbis]tbis event %v", tbisEvent)
	return &relay_chain.BeginCrossChainRequest{
		Version:        common.Version_V1_0_0,
		CrossChainName: tbisEvent.Business,
		CrossChainFlag: TbisFlag,
		CrossChainMsg: []*common.CrossChainMsg{
			{
				GatewayId:    event.CrossChainCreate[0].GatewayId,
				ChainRid:     tbisEvent.SubEvents[0].ChainRID,
				ContractName: tbisEvent.SubEvents[0].ChainCtt.Ctt.Name,
				Method:       tbisEvent.SubEvents[0].ChainCtt.Ctt.Method,
				//Identity:     tbisEvent.SubEvents[0].ChainCtt.Identity,
				Parameter: fmt.Sprintf(tbisEvent.SubEvents[0].ChainCtt.ParamsFormat,
					convertToInterfaceArray(tbisEvent.SubEvents[0].ChainCtt.ParamsData, false)...),
				ParamDataType: subParamType,
				Abi:           subAbi,
			},
		},
		TxContent: &common.TxContent{
			TxId:        eventInfo.TxId,
			Tx:          eventInfo.Tx,
			GatewayId:   gatewayId,
			ChainRid:    eventInfo.ChainRid,
			TxProve:     eventInfo.TxProve,
			BlockHeight: eventInfo.BlockHeight,
		},
		From:    gatewayId,
		Timeout: event.Timeout,
		ConfirmInfo: &common.ConfirmInfo{
			ChainRid:     event.ChainRid,
			ContractName: tbisEvent.CommitCtt.Ctt.Name,
			Method:       tbisEvent.CommitCtt.Ctt.Method,
			Parameter: fmt.Sprintf(tbisEvent.CommitCtt.ParamsFormat,
				convertToInterfaceArray(tbisEvent.CommitCtt.ParamsData, true)...),
			Abi:           majorAbi,
			ParamDataType: commitParamType,
		},
		CancelInfo: &common.CancelInfo{
			ChainRid:     event.ChainRid,
			ContractName: tbisEvent.CommitCtt.Ctt.Name,
			Method:       tbisEvent.CommitCtt.Ctt.Method,
			Parameter: fmt.Sprintf(tbisEvent.CommitCtt.ParamsFormat,
				convertToInterfaceArray(tbisEvent.CommitCtt.ParamsData, true)...),
			Abi:           majorAbi,
			ParamDataType: commitParamType,
		},
		CrossType: event.CrossType,
	}, nil
}

// GetCommitParam 获取提交参数
//  @param chainRID
//  @param proveStatus
//  @param contractStatus
//  @param contractResult
//  @return string
func GetCommitParam(chainRID string, proveStatus, contractStatus int, contractResult string) string {
	chainTxContexts := NewChainTxContexts()
	contractResult = base64.StdEncoding.EncodeToString([]byte(contractResult))
	chainTxContexts.Append(NewChainTxContext(chainRID, proveStatus, contractStatus, contractResult))
	return chainTxContexts.ToContractParam()
}

func convertToInterfaceArray(srcData []string, nilStart bool) []interface{} {
	var (
		dstData []interface{}
	)
	if nilStart {
		dstData = make([]interface{}, 1)
		dstData[0] = "%s"
	} else {
		dstData = make([]interface{}, 0)
	}
	for _, data := range srcData {
		dataTmp, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			dstData = append(dstData, "")
		} else {
			dstData = append(dstData, dataTmp)
		}
	}
	return dstData
}
