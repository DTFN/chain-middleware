/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package relay_chain_chainmaker

import (
	"encoding/json"
	"errors"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

var (
	getGatewayMethod        = utils.GetGateway
	saveGatewayMethod       = utils.SaveGateway
	updateGatewayMethod     = utils.UpadteGateway
	getGatewayByRangeMethod = utils.GetGatewayByRange
	getGatewayNumMethod     = utils.GetGatewayNum

	getCrossChainMethod        = utils.GetCrossChainInfo
	getCrossChainByRangeMethod = utils.GetCrossChainInfoByRange
	getCrossChainNumMethod     = utils.GetCrossChainNum
)

const (
	iAmErr = "I am error"
)

var (
	gatewayInfos = map[string]*tcipcommon.GatewayInfo{
		"0": {
			GatewayId:    "0",
			TxVerifyType: tcipcommon.TxVerifyType_SPV,
			RelayChainId: "0",
		},
		"1": {
			GatewayId:         "1",
			TxVerifyType:      tcipcommon.TxVerifyType_RPC_INTERFACE,
			TxVerifyInterface: &tcipcommon.TxVerifyInterface{},
		},
		"2": {
			GatewayId:    "2",
			TxVerifyType: tcipcommon.TxVerifyType_NOT_NEED,
		},
		"3": {
			GatewayId:    "1",
			TxVerifyType: tcipcommon.TxVerifyType_RPC_INTERFACE,
		},
		"4": {
			GatewayId:       "1",
			TxVerifyType:    tcipcommon.TxVerifyType_RPC_INTERFACE,
			RelayChainId:    "0",
			FromGatewayList: []string{"100"},
		},
	}
	crossChainInfos = map[string]*tcipcommon.CrossChainInfo{
		"0": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
		},
		"1": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_TX_SUCCESS,
					},
				},
			},
		},
		"2": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_TX_SUCCESS,
					},
				},
			},
			CrossChainResult: true,
		},
		"3": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_TX_SUCCESS,
					},
					TxVerifyResult: tcipcommon.TxVerifyRsult_VERIFY_SUCCESS,
				},
				{
					TxVerifyResult: tcipcommon.TxVerifyRsult_VERIFY_INVALID,
				},
			},
			CrossChainResult: true,
		},
		"4": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "99",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_TX_SUCCESS,
					},
					TxVerifyResult: tcipcommon.TxVerifyRsult_VERIFY_SUCCESS,
				},
				{
					TxVerifyResult: tcipcommon.TxVerifyRsult_VERIFY_INVALID,
				},
			},
			CrossChainResult: true,
		},
		"5": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "99",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType:        tcipcommon.CrossType_INVOKE,
			CrossChainResult: true,
		},
		"6": {
			CrossChainId:   "99",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "99",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType:        tcipcommon.CrossType_INVOKE,
			CrossChainResult: true,
		},
		"7": {
			CrossChainId:   "7",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "4",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType:        tcipcommon.CrossType_INVOKE,
			CrossChainResult: true,
		},
		"8": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "99",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_TX_SUCCESS,
					},
				},
			},
			CrossChainResult: true,
		},
		"9": {
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			CrossChainMsg: []*tcipcommon.CrossChainMsg{
				{
					GatewayId:    "0",
					ChainRid:     "chain1",
					ContractName: "test",
					Method:       "test",
					Parameter:    "{\"a\":\"MA==\"}",
					ConfirmInfo: &tcipcommon.ConfirmInfo{
						ChainRid: "chain1",
					},
					CancelInfo: &tcipcommon.CancelInfo{
						ChainRid: "chain1",
					},
				},
			},
			FirstTxContent: &tcipcommon.TxContentWithVerify{},
			From:           "0",
			Timeout:        1000,
			ConfirmInfo: &tcipcommon.ConfirmInfo{
				ChainRid: "chain1",
			},
			CancelInfo: &tcipcommon.CancelInfo{
				ChainRid: "chain1",
			},
			CrossType: tcipcommon.CrossType_INVOKE,
			CrossChainTxContent: []*tcipcommon.TxContentWithVerify{
				{
					TxContent: &tcipcommon.TxContent{
						TxResult: tcipcommon.TxResultValue_GATEWAY_PINGPONG_ERROR,
					},
				},
			},
			CrossChainResult: true,
		},
	}
)

// RelayChainChainmakerMock 中继链mock
type RelayChainChainmakerMock struct {
	config            *conf.RelayChain
	log               *zap.SugaredLogger
	gatewayMap        map[string]*tcipcommon.GatewayInfo
	crossChainInfoMap map[string]*tcipcommon.CrossChainInfo
}

// InitRelayChainMock 中继链mock
//  @param relayChainConfig
//  @return error
func InitRelayChainMock(relayChainConfig *conf.RelayChain) error {
	log := logger.GetLogger(logger.ModuleRelayChain)
	relayChainChainmaker := &RelayChainChainmakerMock{
		config:            relayChainConfig,
		log:               log,
		gatewayMap:        make(map[string]*tcipcommon.GatewayInfo),
		crossChainInfoMap: make(map[string]*tcipcommon.CrossChainInfo),
	}
	RelayChainV1 = relayChainChainmaker
	utils.CrossChainTryChan = make(chan string)
	utils.CrossChainResultChan = make(chan string)
	utils.CrossChainConfirmChan = make(chan string)
	utils.CrossChainSrcGatewayConfirmChan = make(chan string)
	utils.DIDManagerUpdateDIDChan = make(chan string)
	return nil
}

// InitContract 中继链mock
//  @receiver r
//  @param contractName
//  @param version
//  @param byteCodeBase64
//  @param kvJsonStr
//  @param withSyncResult
//  @param timeout
//  @param runtime
//  @return error
func (r *RelayChainChainmakerMock) InitContract(
	contractName, version, byteCodeBase64 string,
	kvJsonStr string,
	withSyncResult bool,
	timeout int64,
	runtime tcipcommon.ChainmakerRuntimeType) error {
	return nil
}

// UpdateContract 中继链mock
//  @receiver r
//  @param contractName
//  @param version
//  @param byteCodeBase64
//  @param kvJsonStr
//  @param withSyncResult
//  @param timeout
//  @param runtime
//  @return error
func (r *RelayChainChainmakerMock) UpdateContract(
	contractName, version, byteCodeBase64 string,
	kvJsonStr string,
	withSyncResult bool,
	timeout int64,
	runtime tcipcommon.ChainmakerRuntimeType) error {
	return nil
}

// InvokeContract 中继链mock
//  @receiver r
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return error
func (r *RelayChainChainmakerMock) InvokeContract(
	contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error) {
	kvMap := make(map[string][]byte)
	_ = json.Unmarshal([]byte(kvJsonStr), &kvMap)
	if method == saveGatewayMethod {
		return []byte("0"), nil
	}
	if method == updateGatewayMethod {
		return []byte("0"), nil
	}
	if method == utils.SaveCrossChainInfo {
		return []byte("0"), nil
	}
	return nil, nil
}

// QueryContract 中继链mock
//  @receiver r
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return error
func (r *RelayChainChainmakerMock) QueryContract(
	contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error) {
	kvMap := make(map[string][]byte)
	_ = json.Unmarshal([]byte(kvJsonStr), &kvMap)
	if method == utils.GetBlockHeaderMethod {
		return []byte("test"), nil
	}
	if method == utils.SpvTxVerifyMethod {
		if kvJsonStr == "{}" {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	}
	if method == getGatewayMethod {
		if string(kvMap[syscontract.GetGateway_GATEWAY_ID.String()]) == "99" {
			return nil, errors.New(iAmErr)
		}
		if string(kvMap[syscontract.GetGateway_GATEWAY_ID.String()]) == "100" {
			return nil, nil
		}
		return proto.Marshal(gatewayInfos[string(kvMap[syscontract.GetGateway_GATEWAY_ID.String()])])
	}
	if method == getCrossChainMethod {
		if string(kvMap[syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String()]) == "" {
			return json.Marshal(crossChainInfos["0"])
		}
		if string(kvMap[syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String()]) == "99" {
			return nil, errors.New(iAmErr)
		}
		if string(kvMap[syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String()]) == "100" {
			return nil, nil
		}
		crossChainInfo, _ := json.Marshal(
			crossChainInfos[string(kvMap[syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String()])])
		return crossChainInfo, nil
	}
	if method == utils.GetNotEndCrossChainIdList {
		return []byte("[\"0\"]"), nil
	}
	if method == getGatewayByRangeMethod {
		gatewayInfo0, _ := proto.Marshal(gatewayInfos["0"])
		gatewayInfo2, _ := proto.Marshal(gatewayInfos["1"])
		return json.Marshal([][]byte{gatewayInfo0, gatewayInfo2})
	}
	if method == getGatewayNumMethod {
		return []byte("5"), nil
	}
	if method == getCrossChainNumMethod {
		return []byte("1"), nil
	}
	if method == getCrossChainByRangeMethod {
		crossChainInfo0, _ := json.Marshal(crossChainInfos["0"])
		return json.Marshal([][]byte{crossChainInfo0})
	}
	return nil, nil
}
