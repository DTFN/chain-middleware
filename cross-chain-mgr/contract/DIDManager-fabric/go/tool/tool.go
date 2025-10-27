/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tool

import (
	"encoding/json"
	"errors"

	"chainmaker.org/chainmaker/tcip-go/common"
	"chainmaker.org/chainmaker/tcip-go/common/relay_chain"
	"github.com/golang/protobuf/proto"
)

// json string, {gateway_id: "0", chain_id: "chain1", contract_name: "crosschain1", method:"cross_chain_try",parameter:"{}or[]",extra_data:"按需写，目标网关能解析就行"}
func CreateEvent(crossChainName, crossChainFlag, confirmInfoStr, cancelInfoStr, crossChainMsgs string) ([]byte, error) {
	var crossChainMsgsInfo []*common.CrossChainMsg
	var confirmInfo common.ConfirmInfo
	var cancelInfo common.CancelInfo
	err := json.Unmarshal([]byte(crossChainMsgs), &crossChainMsgsInfo)
	if err != nil {
		return nil, errors.New("Please enter the correct crossChainMsgs: " + crossChainMsgs + " " + err.Error())
	}
	if confirmInfoStr != "" {
		err = json.Unmarshal([]byte(confirmInfoStr), &confirmInfo)
		if err != nil {
			return nil, errors.New("Please enter the correct confirmInfo: " + confirmInfoStr + " " + err.Error())
		}
	}
	if cancelInfoStr != "" {
		err = json.Unmarshal([]byte(cancelInfoStr), &cancelInfo)
		if err != nil {
			return nil, errors.New("Please enter the correct cancelInfo: " + cancelInfoStr + " " + err.Error())
		}
	}
	beginCrossChainRequest := &relay_chain.BeginCrossChainRequest{
		Version:        common.Version_V1_0_0,
		CrossChainName: crossChainName,
		CrossChainFlag: crossChainFlag,
		CrossChainMsg:  crossChainMsgsInfo,
		ConfirmInfo:    &confirmInfo,
		CancelInfo:     &cancelInfo,
	}
	requestByte, err := proto.Marshal(beginCrossChainRequest)
	if err != nil {
		return nil, errors.New("Marshal BeginCrossChainRequest failed: " + err.Error())
	}
	return requestByte, nil
}
