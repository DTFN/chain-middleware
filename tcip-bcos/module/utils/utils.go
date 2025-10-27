/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

const (
	tempFileFormat = "temp-*"
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

// WriteTempFile 内容写入到临时文件中
//  @param content
//  @param suffix
//  @return *os.File
//  @return error
func WriteTempFile(content []byte, suffix string) (*os.File, error) {
	tempFile, err := CreateTempFile(suffix)
	if err != nil {
		return nil, err
	}
	// 将内容写入文件
	defer tempFile.Close()
	_, err = tempFile.Write(content)
	return tempFile, err
}

// CreateTempFile 创建临时文件
//  @param suffix
//  @return *os.File
//  @return error
func CreateTempFile(suffix string) (*os.File, error) {
	return CreateCertainTempFile("", suffix)
}

// CreateCertainTempFile 创建确定的临时文件（明确目录）
//  @param dir
//  @param suffix
//  @return *os.File
//  @return error
func CreateCertainTempFile(dir, suffix string) (*os.File, error) {
	if dir != "" {
		exist, err := FileIsExist(dir)
		if err != nil {
			return nil, err
		}
		if !exist {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}
	tempFile, err := os.CreateTemp(dir, tempFileFormat+suffix)
	if err != nil {
		return nil, err
	}
	return tempFile, nil
}

// FileIsExist 检测文件是否存在
//  @param path
//  @return bool
//  @return error
func FileIsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// UnsupportVersion 不支持的版本打印
//  @param version
//  @return string
func UnsupportVersion(version common.Version) string {
	return fmt.Sprintf("Unsupported version: %d", version)
}
