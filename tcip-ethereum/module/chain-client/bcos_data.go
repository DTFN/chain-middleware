/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

// BCOSBlock bcos的区块结构
type BCOSBlock struct {
	Transactions     []string `json:"transactions"`
	Number           string   `json:"number"`
	Hash             string   `json:"hash"`
	ParentHash       string   `json:"parentHash"`
	TransactionsRoot string   `json:"transactionsRoot"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	DbHash           string   `json:"dbHash"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
}

// TransactionHashs 获取交易哈希
//  @receiver b
//  @return []string
func (b *BCOSBlock) TransactionHashs() []string {
	return b.Transactions
}

// BCOSEvent bcos的event
type BCOSEvent struct {
	Topic        string
	ContractName string
	Data         []string
}
