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

// NewRestRequest restrequest新建
//  @param log
//  @return *RestRequest
func NewRestRequest(log *zap.SugaredLogger) *RestRequest {
	return &RestRequest{
		log: log,
	}
}

// BeginCrossChain 开始跨链
//  @receiver r
//  @param req
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func (r *RestRequest) BeginCrossChain(
	req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error) {
	panic("error")
}

// SyncBlockHeader 同步区块头
//  @receiver r
//  @param req
//  @return *relay_chain.SyncBlockHeaderResponse
//  @return error
func (r *RestRequest) SyncBlockHeader(
	req *relay_chain.SyncBlockHeaderRequest) (*relay_chain.SyncBlockHeaderResponse, error) {
	panic("error")
}

// InitSpvContract 初始化spv合约
//  @receiver r
//  @param req
//  @return *relay_chain.InitContractResponse
//  @return error
func (r *RestRequest) InitSpvContract(
	req *relay_chain.InitContractRequest) (*relay_chain.InitContractResponse, error) {
	panic("error")
}

// UpdateSpvContract 更新spv合约
//  @receiver r
//  @param req
//  @return *relay_chain.UpdateContractResponse
//  @return error
func (r *RestRequest) UpdateSpvContract(
	req *relay_chain.UpdateContractRequest) (*relay_chain.UpdateContractResponse, error) {
	panic("error")
}
