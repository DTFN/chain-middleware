/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package tool

import (
	"encoding/json"
	"fmt"
	"testing"

	"chainmaker.org/chainmaker/tcip-go/common"
)

func TestCreateEvent(t *testing.T) {
	var crossChainMsg []*common.CrossChainMsg
	err := json.Unmarshal([]byte("[{\"gateway_id\": \"0\",\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"CrossChainTry\",\"parameter\":\"[\\\"method\\\",\\\"CrossChainTry\\\"]\",\"extra_data\":\"按需写，目标网关能解析就行\",\"confirm_info\":{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"cross_chain_confirm\"},\"cancel_info\":{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"cross_chain_cancel\"}}]"), &crossChainMsg)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("%+v\n", crossChainMsg)
	}
}
