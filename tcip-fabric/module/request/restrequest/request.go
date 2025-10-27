/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package restrequest

import (
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"go.uber.org/zap"
)

// RestRequest rest请求结构体
type RestRequest struct {
	log *zap.SugaredLogger
}

// NewRestRequest 初始化rest请求
//  @param log
//  @return *RestRequest
func NewRestRequest(log *zap.SugaredLogger) *RestRequest {
	return &RestRequest{
		log: log,
	}
}

// BeginCrossChain 调用跨链接口
//  @receiver r
//  @param req
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func (r *RestRequest) BeginCrossChain(
	req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error) {
	panic("error")
}
