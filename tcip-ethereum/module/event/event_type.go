/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import "fmt"

// EventInfo 产生事件的结构
type EventInfo struct {
	Topic        string
	ChainRid     string
	ContractName string
	TxProve      string
	Data         []string
	Tx           []byte
	TxId         string
	BlockHeight  int64
}

// ToString 转为string展示
//  @receiver e
//  @return string
func (e *EventInfo) ToString() string {
	return fmt.Sprintf("Topic: %s, ChainRid: %s, ContractName: %s, TxProve: %s,"+
		" Data: %v, txId: %s, BlockHeight: %d", e.Topic, e.ChainRid, e.ContractName,
		e.TxProve, e.Data, e.TxId, e.BlockHeight)
}
