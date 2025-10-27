/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

import (
	"errors"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
)

type requestMock struct {
}

// BeginCrossChain 发送跨链请求
//  @receiver r
//  @param req
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func (r *requestMock) BeginCrossChain(
	req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error) {
	switch req.Version {
	case common.Version_V1_0_0:
		return &relay_chain.BeginCrossChainResponse{
			CrossChainId: "0",
			Code:         common.Code_GATEWAY_SUCCESS,
			Message:      common.Code_GATEWAY_SUCCESS.String(),
		}, nil
	default:
		return nil, errors.New("unsupported version")
	}
}

// InitRequestManagerMock 初始化request模块
//  @return error
func InitRequestManagerMock() error {
	log := logger.GetLogger(logger.ModuleRequest)
	RequestV1 = &RequestManager{
		request: &requestMock{},
		log:     log,
	}
	return nil
}
