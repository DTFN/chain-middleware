/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/event"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/request/grpcrequest"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/request/restrequest"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"go.uber.org/zap"
)

// Request 请求模块接口
type Request interface {
	BeginCrossChain(req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error)
}

// RequestManager 请求管理模块结构体
type RequestManager struct {
	log     *zap.SugaredLogger
	request Request
}

// RequestV1 请求管理模块对象
var RequestV1 *RequestManager

// InitRequestManager 初始化请求管理模块
//  @return error
func InitRequestManager() error {
	log := logger.GetLogger(logger.ModuleRequest)
	var request Request
	if conf.Config.Relay.CallType == conf.GrpcCallType {
		request = grpcrequest.NewGrpcRequest(log)
	} else if conf.Config.Relay.CallType == conf.RestCallType {
		request = restrequest.NewRestRequest(log)
	} else {
		panic("unsupport call_type:" + conf.Config.Relay.CallType)
	}
	RequestV1 = &RequestManager{
		request: request,
		log:     log,
	}
	return nil
}

// GatewayRegister 获取gateway注册信息
//  @receiver r
//  @param objectPath
//  @return error
func (r *RequestManager) GatewayRegister(objectPath string) error {
	tlsca, err := ioutil.ReadFile(conf.Config.BaseConfig.Tlsca)
	if err != nil {
		errMsg := fmt.Sprintf("[GatewayRegister] read tlsca error: %s", conf.Config.BaseConfig.Tlsca)
		r.log.Error(errMsg)
		panic(errMsg)
	}
	clientCert, err := ioutil.ReadFile(conf.Config.BaseConfig.ClientCert)
	if err != nil {
		errMsg := fmt.Sprintf("[GatewayRegister] read client_cert error: %s", conf.Config.BaseConfig.ClientCert)
		r.log.Error(errMsg)
		panic(errMsg)
	}

	req := &relay_chain.GatewayRegisterRequest{
		Version: common.Version_V1_0_0,
		GatewayInfo: buildGatewayInfo("", conf.Config.BaseConfig.GatewayName, conf.Config.BaseConfig.Address,
			conf.Config.BaseConfig.ServerName, string(tlsca), string(clientCert), conf.Config.BaseConfig.ToGatewayList,
			conf.Config.BaseConfig.FromGatewayList, r.log),
	}
	reqByte, _ := json.Marshal(req)
	return ioutil.WriteFile(objectPath, reqByte, 0600)
}

// GatewayUpdate 获取gateway更新信息
//  @receiver r
//  @param objectPath
//  @return error
func (r *RequestManager) GatewayUpdate(objectPath string) error {
	tlsca, err := ioutil.ReadFile(conf.Config.BaseConfig.Tlsca)
	if err != nil {
		errMsg := fmt.Sprintf("[GatewayUpdate] read tlsca error: %s", conf.Config.BaseConfig.Tlsca)
		r.log.Error(errMsg)
		panic(errMsg)
	}
	clientCert, err := ioutil.ReadFile(conf.Config.BaseConfig.ClientCert)
	if err != nil {
		errMsg := fmt.Sprintf("[GatewayUpdate] read client_cert error: %s", conf.Config.BaseConfig.ClientCert)
		r.log.Error(errMsg)
		panic(errMsg)
	}

	req := &relay_chain.GatewayUpdateRequest{
		Version: common.Version_V1_0_0,
		GatewayInfo: buildGatewayInfo(conf.Config.BaseConfig.GatewayID, conf.Config.BaseConfig.GatewayName,
			conf.Config.BaseConfig.Address, conf.Config.BaseConfig.ServerName, string(tlsca), string(clientCert),
			conf.Config.BaseConfig.ToGatewayList, conf.Config.BaseConfig.FromGatewayList, r.log),
	}
	reqByte, _ := json.Marshal(req)
	return ioutil.WriteFile(objectPath, reqByte, 0600)
}

// BeginCrossChain 开始跨链
//  @receiver r
//  @param eventInfo
// 如果要保存跨链信息的话，在这个函数里面实现就可以,可以根据结果写入数据库或者文件什么的都可以
func (r *RequestManager) BeginCrossChain(eventInfo *event.EventInfo) {
	beginCrossChainRequest, err := event.EventManagerV1.BuildCrossChainMsg(eventInfo)
	if err != nil {
		r.log.Errorf("[BeginCrossChain] %s", err.Error())
		return
	}
	if beginCrossChainRequest == nil {
		r.log.Warnf("[BeginCrossChain] build beginCrossChainRequest failed: topic %s", eventInfo.Topic)
		return
	}
	r.log.Infof("[BeginCrossChain] Call tcip-relayer BeginCrossChain method start: topic %s, request %+v",
		eventInfo.Topic, beginCrossChainRequest)
	res, err := r.request.BeginCrossChain(beginCrossChainRequest)
	//这里出错的话应该有一个回滚的操作，这个操作需要用户根据自己的实际情况来处理
	if err != nil {
		r.log.Errorf("[BeginCrossChain] Call tcip-relayer BeginCrossChain method error: topic %s, error %s",
			eventInfo.Topic, err.Error())
		return
	}
	resString, _ := json.Marshal(res)
	if res.Code != common.Code_GATEWAY_SUCCESS {
		r.log.Errorf("[BeginCrossChain] Call tcip-relayer BeginCrossChain method error: topic %s, response %s",
			eventInfo.Topic, string(resString))
		return
	}
	r.log.Info("[BeginCrossChain] Call tcip-relayer BeginCrossChain method success: topic %s, response %s",
		eventInfo.Topic, string(resString))
}

func getTxVerifyType(log *zap.SugaredLogger) common.TxVerifyType {
	var txVerifyType common.TxVerifyType
	switch conf.Config.BaseConfig.TxVerifyType {
	case conf.RpcTxVerify:
		txVerifyType = common.TxVerifyType_RPC_INTERFACE
	case conf.SpvTxVerify:
		txVerifyType = common.TxVerifyType_SPV
	case conf.NotNeedTxVerify:
		txVerifyType = common.TxVerifyType_NOT_NEED
	default:
		errMsg := fmt.Sprintf("unsupported tx_verify_type: %s", conf.Config.BaseConfig.TxVerifyType)
		log.Errorf(errMsg)
		panic(errMsg)
	}
	return txVerifyType
}

func getTxVerifyInterface(log *zap.SugaredLogger) *common.TxVerifyInterface {
	var txVerifyInterface *common.TxVerifyInterface
	if conf.Config.BaseConfig.TxVerifyType == conf.RpcTxVerify {
		v := conf.Config.BaseConfig.TxVerifyInterface
		tlscaTxVerify, err := ioutil.ReadFile(v.Tlsca)
		if err != nil {
			errMsg := fmt.Sprintf("read tx_verify_interface tlsca error: %s", v.Tlsca)
			log.Error(errMsg)
			panic(errMsg)
		}
		clientCertTxVerify, err := ioutil.ReadFile(v.ClientCert)
		if err != nil {
			errMsg := fmt.Sprintf("read tx_verify_interface client_cert error: %s", v.ClientCert)
			log.Errorf(errMsg)
			panic(errMsg)
		}
		txVerifyInterface = &common.TxVerifyInterface{
			Address:    v.Address,
			TlsEnable:  v.TlsEnable,
			Tlsca:      string(tlscaTxVerify),
			ClientCert: string(clientCertTxVerify),
			HostName:   v.HostName,
			//ClientKey:  v.ClientKey,
		}
	}
	return txVerifyInterface
}

func getCallType(log *zap.SugaredLogger) common.CallType {
	switch conf.Config.BaseConfig.CallType {
	case conf.GrpcCallType:
		return common.CallType_GRPC
	case conf.RestCallType:
		return common.CallType_REST
	default:
		errMsg := fmt.Sprintf("Error call_type: need grpc or restful, not %s", conf.Config.BaseConfig.CallType)
		log.Error(errMsg)
		panic(errMsg)
	}
}

func buildGatewayInfo(gatewayId, gatewayName, address, serverName, tlsca, clientCert string,
	toGateWayList, fromGatewayList []string, log *zap.SugaredLogger) *common.GatewayInfo {
	return &common.GatewayInfo{
		GatewayId:   gatewayId,
		GatewayName: gatewayName,
		Address:     address,
		ServerName:  serverName,
		Tlsca:       tlsca,
		ClientCert:  clientCert,
		//ClientKey:         conf.Config.BaseConfig.ClientKey,
		ToGatewayList:     toGateWayList,
		FromGatewayList:   fromGatewayList,
		TxVerifyType:      getTxVerifyType(log),
		TxVerifyInterface: getTxVerifyInterface(log),
		CallType:          getCallType(log),
	}
}
