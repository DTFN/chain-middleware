/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

//EventOperate event更新结构体
type EventOperate struct {
	CrossChainEventID string
	ChainRid          string
	ContractName      string
	Operate           common.Operate
}

//ChainConfigOperate event更新结构体
type ChainConfigOperate struct {
	ChainRid string
	Operate  common.Operate
}

//EventChan event更新通道
var EventChan chan *EventOperate

// UpdateChainConfigChan chainconfig 更新通道
var UpdateChainConfigChan chan *ChainConfigOperate

// DeepCopy 结构体深拷贝
//  @param dst
//  @param src
//  @return error
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

// UnsupportVersion 不支持的版本打印
//  @param version
//  @return string
func UnsupportVersion(version common.Version) string {
	return fmt.Sprintf("Unsupported version: %d", version)
}

// ParseHeaderKey 生成存储同步区块头的key
//  @receiver c
//  @param chainRid
//  @return []byte
func ParseHeaderKey(chainRid string) []byte {
	return []byte(fmt.Sprintf("SYNCBH#%s", chainRid))
}
