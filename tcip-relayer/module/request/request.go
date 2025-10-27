/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request/grpcrequest"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request/restrequest"
	"go.uber.org/zap"
)

//@todo 这个包有大量的重复代码，后续可以用空接口来优化结构

// RequestV1 请求全局变量
var RequestV1 *RequestManager

// Request 请求接口
type Request interface {
	CrossChainTry(
		txRequest *cross_chain.CrossChainTryRequest,
		timeout int64,
		destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainTryResponse, error)
	CrossChainConfirm(
		txRequest *cross_chain.CrossChainConfirmRequest,
		timeout int64,
		destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainConfirmResponse, error)
	CrossChainCancel(
		txRequest *cross_chain.CrossChainCancelRequest,
		timeout int64,
		destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainCancelResponse, error)
	IsCrossChainSuccess(
		txRequest *cross_chain.IsCrossChainSuccessRequest,
		timeout int64,
		destGatewayInfo *common.GatewayInfo) (*cross_chain.IsCrossChainSuccessResponse, error)
	PingPong(timeout int64, destGatewayInfo *common.GatewayInfo) (*cross_chain.PingPongResponse, error)
	VerifyTx(txVerifyInterface *common.TxVerifyInterface, txProve string) ([]byte, error)
}

// RequestManager 请求结构体
type RequestManager struct {
	log         *zap.SugaredLogger
	restRequest Request
	grpcRequest Request
}

// InitRequestManager 初始化请求
//  @return error
func InitRequestManager() error {
	log := logger.GetLogger(logger.ModuleRequest)
	RequestV1 = &RequestManager{
		restRequest: restrequest.NewRestRequest(log),
		grpcRequest: grpcrequest.NewGrpcRequest(log),
		log:         log,
	}
	return nil
}

// CrossChainTry 跨链执行
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainTryResponse
//  @return error
func (r *RequestManager) CrossChainTry(
	txRequest *cross_chain.CrossChainTryRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainTryResponse, error) {
	var (
		msg string
		res *cross_chain.CrossChainTryResponse
		err error
	)
	switch destGatewayInfo.CallType {
	case common.CallType_GRPC:
		res, err = r.grpcRequest.CrossChainTry(txRequest, timeout, destGatewayInfo)
	case common.CallType_REST:
		res, err = r.restRequest.CrossChainTry(txRequest, timeout, destGatewayInfo)
	default:
		msg = fmt.Sprintf("[CrossChainTry] Unsupported call type: %s", destGatewayInfo.CallType)
		err = errors.New(msg)
	}
	if err != nil {
		msg = fmt.Sprintf(
			"[crossChainConfirm] get gateway response error, gatewayName: %s, gatewayId: %s, error: %s",
			destGatewayInfo.GatewayName, destGatewayInfo.GatewayId, err.Error())
		r.log.Error(msg)
		return res, errors.New(msg)
	}
	return res, nil
}

// CrossChainConfirm 跨链提交
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainConfirmResponse
//  @return error
func (r *RequestManager) CrossChainConfirm(
	txRequest *cross_chain.CrossChainConfirmRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainConfirmResponse, error) {
	var (
		msg string
		res *cross_chain.CrossChainConfirmResponse
		err error
	)
	switch destGatewayInfo.CallType {
	case common.CallType_GRPC:
		res, err = r.grpcRequest.CrossChainConfirm(txRequest, timeout, destGatewayInfo)
	case common.CallType_REST:
		res, err = r.restRequest.CrossChainConfirm(txRequest, timeout, destGatewayInfo)
	default:
		msg = fmt.Sprintf("[crossChainConfirm] Unsupported call type: %s", destGatewayInfo.CallType)
		err = errors.New(msg)
	}
	if err != nil {
		msg = fmt.Sprintf(
			"[crossChainConfirm] get gateway response error, gatewayName: %s, gatewayId: %s, error: %s",
			destGatewayInfo.GatewayName, destGatewayInfo.GatewayId, err.Error())
		r.log.Error(msg)
		return nil, errors.New(msg)
	}
	return res, nil
}

// CrossChainCancel 跨链回滚
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainCancelResponse
//  @return error
func (r *RequestManager) CrossChainCancel(
	txRequest *cross_chain.CrossChainCancelRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainCancelResponse, error) {
	var (
		msg string
		res *cross_chain.CrossChainCancelResponse
		err error
	)
	switch destGatewayInfo.CallType {
	case common.CallType_GRPC:
		res, err = r.grpcRequest.CrossChainCancel(txRequest, timeout, destGatewayInfo)
	case common.CallType_REST:
		res, err = r.restRequest.CrossChainCancel(txRequest, timeout, destGatewayInfo)
	default:
		msg = fmt.Sprintf("[CrossChainCancel]Unsupported call type: %s", destGatewayInfo.CallType)
		err = errors.New(msg)
	}
	if err != nil {
		msg = fmt.Sprintf(
			"[CrossChainCancel] get gateway response error, gatewayName: %s, gatewayId: %s, error: %s",
			destGatewayInfo.GatewayName, destGatewayInfo.GatewayId, err.Error())
		r.log.Error(msg)
		return nil, errors.New(msg)
	}
	return res, nil
}

// IsCrossChainSuccess 询问跨链是否成功
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.IsCrossChainSuccessResponse
//  @return error
func (r *RequestManager) IsCrossChainSuccess(
	txRequest *cross_chain.IsCrossChainSuccessRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.IsCrossChainSuccessResponse, error) {
	var (
		msg string
		res *cross_chain.IsCrossChainSuccessResponse
		err error
	)
	switch destGatewayInfo.CallType {
	case common.CallType_GRPC:
		res, err = r.grpcRequest.IsCrossChainSuccess(txRequest, timeout, destGatewayInfo)
	case common.CallType_REST:
		res, err = r.restRequest.IsCrossChainSuccess(txRequest, timeout, destGatewayInfo)
	default:
		msg = fmt.Sprintf("[IsCrossChainSuccess] Unsupported call type: %s", destGatewayInfo.CallType)
		err = errors.New(msg)
	}
	if err != nil {
		msg = fmt.Sprintf(
			"[IsCrossChainSuccess] get gateway response error, gatewayName: %s, gatewayId: %s, error: %s",
			destGatewayInfo.GatewayName, destGatewayInfo.GatewayId, err.Error())
		r.log.Error(msg)
		return nil, errors.New(msg)
	}
	return res, nil
}

// PingPong 心跳
//  @receiver r
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.PingPongResponse
//  @return error
func (r *RequestManager) PingPong(timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.PingPongResponse, error) {
	var (
		msg string
		err error
		res *cross_chain.PingPongResponse
	)
	switch destGatewayInfo.CallType {
	case common.CallType_GRPC:
		res, err = r.grpcRequest.PingPong(timeout, destGatewayInfo)
	case common.CallType_REST:
		res, err = r.restRequest.PingPong(timeout, destGatewayInfo)
	default:
		msg = fmt.Sprintf("[PingPong] Unsupported call type: %s", destGatewayInfo.CallType)
		err = errors.New(msg)
	}
	if err != nil {
		msg = fmt.Sprintf(
			"[PingPong] can't connect to gateway, Name: %s, Id: %s, error: %s",
			destGatewayInfo.GatewayName, destGatewayInfo.GatewayId, err.Error())
		r.log.Error(msg)
		err = errors.New(msg)
		return nil, err
	}
	r.log.Debugf("[PingPong] res:%+v", res)
	return res, nil
}

// VerifyTx 交易验证
//  @receiver r
//  @param txVerifyInterface
//  @param txProve
//  @return []byte
//  @return error
func (r *RequestManager) VerifyTx(txVerifyInterface *common.TxVerifyInterface, txProve string) ([]byte, error) {
	// 当前仅支持restful的验证方式
	return r.restRequest.VerifyTx(txVerifyInterface, txProve)
}
