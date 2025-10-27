/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/db"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"

	chainmaker_common "chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/request/grpcrequest"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/request/restrequest"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

// Request 请求接口
type Request interface {
	BeginCrossChain(req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error)
	SyncBlockHeader(req *relay_chain.SyncBlockHeaderRequest) (*relay_chain.SyncBlockHeaderResponse, error)
	InitSpvContract(req *relay_chain.InitContractRequest) (*relay_chain.InitContractResponse, error)
	UpdateSpvContract(req *relay_chain.UpdateContractRequest) (*relay_chain.UpdateContractResponse, error)
}

// RequestManager 请求管理结构体
type RequestManager struct {
	log     *zap.SugaredLogger
	request Request
}

//RequestV1 rquest模块对象
var RequestV1 *RequestManager

// InitRequestManager 初始化request
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

// GatewayUpdate 获取更新gateway的信息
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

// BeginCrossChain 如果要保存跨链信息的话，在这个函数里面实现就可以,可以根据结果写入数据库或者文件什么的都可以
//  @receiver r
//  @param eventInfo
func (r *RequestManager) BeginCrossChain(eventInfo *event.EventInfo) {
	r.log.Debugf("[BeginCrossChain] eventInfo %+v", eventInfo)
	beginCrossChainRequest, err := event.EventManagerV1.BuildCrossChainMsg(eventInfo)
	if err != nil {
		r.log.Errorf("[BeginCrossChain] %s", err.Error())
		return
	}
	r.log.Debugf("[BeginCrossChain] beginCrossChainRequest %+v", beginCrossChainRequest)
	if beginCrossChainRequest == nil {
		r.log.Warnf("[BeginCrossChain] build beginCrossChainRequest failed: topic %s", eventInfo.Topic)
		return
	}
	r.log.Info("[BeginCrossChain] Call tcip-relayer BeginCrossChain method start: topic %s, request %+v",
		eventInfo.Topic, beginCrossChainRequest)

	res, err := r.request.BeginCrossChain(beginCrossChainRequest)
	//这里出错的话应该有一个回滚的操作，这个操作需要用户根据自己的实际情况来处理
	if err != nil {
		r.log.Errorf("[BeginCrossChain] Call tcip-relayer BeginCrossChain method "+
			"error: topic %s, error %s, txId: %s",
			eventInfo.Topic, err.Error(), eventInfo.TxId)
		return
	}
	resString, _ := json.Marshal(res)
	if res.Code != common.Code_GATEWAY_SUCCESS {
		r.log.Errorf("[BeginCrossChain] Call tcip-relayer BeginCrossChain method "+
			"error: topic %s, response %s, txId: %s",
			eventInfo.Topic, string(resString), eventInfo.TxId)
		return
	}
	r.log.Info("[BeginCrossChain] Call tcip-relayer BeginCrossChain method "+
		"success: topic %s, response %s, txId %s",
		eventInfo.Topic, string(resString), eventInfo.TxId)
}

// SyncBlockHeader 同步区块头
//  @receiver r
//  @param blockHeader
//  @param chainRid
func (r *RequestManager) SyncBlockHeader(blockHeader *chainmaker_common.BlockHeader, chainRid string) {
	beginTime := time.Now().Unix()
	blockHeaderByte, err := proto.Marshal(blockHeader)
	if err != nil {
		r.log.Errorf("[SyncBlockHeader]Marshal blockHeader failed: error: %s, chainId: %s",
			err.Error(), blockHeader.ChainId)
		return
	}
	request := &relay_chain.SyncBlockHeaderRequest{
		Version:     common.Version_V1_0_0,
		GatewayId:   conf.Config.BaseConfig.GatewayID,
		ChainRid:    chainRid,
		BlockHeight: blockHeader.BlockHeight,
		BlockHeader: blockHeaderByte,
	}
	for {
		res, err := r.request.SyncBlockHeader(request)
		if err != nil {
			r.log.Errorf("[SyncBlockHeader]Request SyncBlockHeader failed: error: %s, chainId: %s",
				err.Error(), blockHeader.ChainId)
			continue
		}
		if res.Code == common.Code_GATEWAY_SUCCESS {
			r.log.Infof("[SyncBlockHeader]SyncBlockHeader success: chainId: %s, blockHeight: %d,"+
				" message: %s, timeUsed: %d",
				blockHeader.ChainId, blockHeader.BlockHeight, res.Message, time.Now().Unix()-beginTime)
			_ = db.Db.Put(utils.ParseHeaderKey(chainRid), []byte(fmt.Sprintf("%d", blockHeader.BlockHeight)))
			return
		}
		r.log.Errorf("[SyncBlockHeader]Request SyncBlockHeader failed: code: %d, error: %s, chainId: %s",
			res.Code, res.Message, blockHeader.ChainId)
		time.Sleep(time.Second * 5)
	}
}

// InitSpvContracta 初始化spv合约
//  @receiver r
//  @param version
//  @param path
//  @param runtimeType
//  @param kvJsonStr
//  @param chainRid
//  @return error
func (r *RequestManager) InitSpvContracta(version, path, runtimeType, kvJsonStr, chainRid string) error {
	byteCode, keyValuePairs, err := dealContract(path, runtimeType, kvJsonStr)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	request := &relay_chain.InitContractRequest{
		Version:         common.Version_V1_0_0,
		ContractVersion: version,
		ByteCode:        []byte(base64.StdEncoding.EncodeToString(byteCode)),
		RuntimeType:     common.ChainmakerRuntimeType(common.ChainmakerRuntimeType_value[runtimeType]),
		KeyValuePairs:   keyValuePairs,
		GatewayId:       conf.Config.BaseConfig.GatewayID,
		ChainRid:        chainRid,
	}
	resp, err := r.request.InitSpvContract(request)
	if err != nil {
		msg := fmt.Sprintf("Init spv contract error: %s", err.Error())
		r.log.Error(msg)
		return errors.New(msg)
	}
	if resp.Code != common.Code_GATEWAY_SUCCESS {
		msg := fmt.Sprintf("Init spv contract error: %s", resp.Message)
		r.log.Error(msg)
		return errors.New(msg)
	}
	r.log.Infof("Init spv contract: %+v", resp)
	return nil
}

// UpdateSpvContract 更新spv合约
//  @receiver r
//  @param version
//  @param path
//  @param runtimeType
//  @param kvJsonStr
//  @param chainRid
//  @return error
func (r *RequestManager) UpdateSpvContract(version, path, runtimeType, kvJsonStr, chainRid string) error {
	byteCode, keyValuePairs, err := dealContract(path, runtimeType, kvJsonStr)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}
	request := &relay_chain.UpdateContractRequest{
		Version:         common.Version_V1_0_0,
		ContractVersion: version,
		ByteCode:        []byte(base64.StdEncoding.EncodeToString(byteCode)),
		RuntimeType:     common.ChainmakerRuntimeType(common.ChainmakerRuntimeType_value[runtimeType]),
		KeyValuePairs:   keyValuePairs,
		GatewayId:       conf.Config.BaseConfig.GatewayID,
		ChainRid:        chainRid,
	}
	resp, err := r.request.UpdateSpvContract(request)
	if err != nil {
		msg := fmt.Sprintf("Init spv contract error: %s", err.Error())
		r.log.Error(msg)
		return errors.New(msg)
	}
	if resp.Code != common.Code_GATEWAY_SUCCESS {
		msg := fmt.Sprintf("Init spv contract error: %s", resp.Message)
		r.log.Error(msg)
		return errors.New(msg)
	}
	r.log.Infof("Init spv contract: %+v", resp)
	return nil
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

func getKvsFromKvJsonStr(kvJsonStr string) ([]*common.ContractKeyValuePair, error) {
	var kvMap map[string]string
	err := json.Unmarshal([]byte(kvJsonStr), &kvMap)
	if err != nil {
		errStr := fmt.Sprintf("[getKvsFromKvJsonStr] kvJsonStr must be json string: %s -> %s",
			kvJsonStr, err.Error())
		return nil, errors.New(errStr)
	}
	kvs := make([]*common.ContractKeyValuePair, 0)
	for k, v := range kvMap {
		kv := &common.ContractKeyValuePair{
			Key:   k,
			Value: []byte(v),
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
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

func dealContract(path, runtimeType, kvJsonStr string) ([]byte, []*common.ContractKeyValuePair, error) {
	byteCode, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("Read File error: %s ", err.Error())
	}
	if _, ok := common.ChainmakerRuntimeType_value[runtimeType]; !ok {
		return nil, nil, fmt.Errorf("unsupported runtimeType: %s", runtimeType)
	}
	keyValuePairs, err := getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return nil, nil, fmt.Errorf("Error JSON string: %s ", err.Error())
	}
	return byteCode, keyValuePairs, nil
}
