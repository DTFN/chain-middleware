/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"encoding/json"
	"errors"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"go.uber.org/zap"
)

type requestMock struct {
}

// CrossChainTry mock request
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainTryResponse
//  @return error
func (r requestMock) CrossChainTry(
	txRequest *cross_chain.CrossChainTryRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainTryResponse, error) {
	switch txRequest.Version {
	case common.Version_V1_0_0:
		return &cross_chain.CrossChainTryResponse{
			CrossChainId:   "0",
			CrossChainName: "test",
			CrossChainFlag: "test",
			Code:           common.Code_GATEWAY_SUCCESS,
			Message:        common.Code_GATEWAY_SUCCESS.String(),
		}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

//CrossChainConfirm mock request
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainConfirmResponse
//  @return error
func (r requestMock) CrossChainConfirm(
	txRequest *cross_chain.CrossChainConfirmRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainConfirmResponse, error) {
	switch txRequest.Version {
	case common.Version_V1_0_0:
		return &cross_chain.CrossChainConfirmResponse{}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// CrossChainCancel mock request
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainCancelResponse
//  @return error
func (r requestMock) CrossChainCancel(
	txRequest *cross_chain.CrossChainCancelRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainCancelResponse, error) {
	switch txRequest.Version {
	case common.Version_V1_0_0:
		return &cross_chain.CrossChainCancelResponse{}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// IsCrossChainSuccess  mock request
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.IsCrossChainSuccessResponse
//  @return error
func (r requestMock) IsCrossChainSuccess(
	txRequest *cross_chain.IsCrossChainSuccessRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.IsCrossChainSuccessResponse, error) {
	switch txRequest.Version {
	case common.Version_V1_0_0:
		return &cross_chain.IsCrossChainSuccessResponse{}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// PingPong mock request
//  @receiver r
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.PingPongResponse
//  @return error
func (r requestMock) PingPong(timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.PingPongResponse, error) {
	return nil, nil
}

// VerifyTx  mock request
//  @receiver r
//  @param txVerifyInterface
//  @param txProve
//  @return []byte
//  @return error
func (r requestMock) VerifyTx(txVerifyInterface *common.TxVerifyInterface, txProve string) ([]byte, error) {
	if txVerifyInterface == nil || txProve == "" {
		return []byte("false"), nil
	}
	res := &cross_chain.TxVerifyResponse{
		TxVerifyResult: true,
	}
	return json.Marshal(res)
}

// NewRequest mock request
//  @param log
//  @return Request
func NewRequest(log *zap.SugaredLogger) Request {
	return &requestMock{}
}

// InitRequestManagerMock mock request
//  @return error
func InitRequestManagerMock() error {
	log := logger.GetLogger(logger.ModuleRequest)
	RequestV1 = &RequestManager{
		restRequest: NewRequest(log),
		grpcRequest: NewRequest(log),
		log:         log,
	}
	return nil
}
