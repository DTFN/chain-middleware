/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
)

var (
	log []*logger.LogModuleConfig

	key1   = []byte("key1")
	value1 = []byte("value1")
	key2   = []byte("key2")
	value2 = []byte("value2")
	key3   = []byte("key3")
	value3 = []byte("value3")
	key4   = []byte("key4")
	value4 = []byte("value4")
)

func initTest() {
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	conf.Config.DbPath = path.Join(os.TempDir(), time.Now().String())
	logger.InitLogConfig(log)
}

func TestPut(t *testing.T) {
	initTest()
	NewDbHandle()

	err := Db.Put(key1, value1)
	assert.Nil(t, err)

	err = Db.Put(key1, nil)
	assert.NotNil(t, err)

	_ = Db.Close()

	err = Db.Put(key1, value1)
	assert.NotNil(t, err)
}

func TestGet(t *testing.T) {
	initTest()
	NewDbHandle()

	err := Db.Put(key1, value1)
	assert.Nil(t, err)

	res, err := Db.Get(key1)
	assert.Equal(t, res, value1)
	assert.Nil(t, err)

	res, err = Db.Get(key2)
	assert.Nil(t, res)
	assert.Nil(t, err)

	_ = Db.Close()

	res, err = Db.Get(key2)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestHas(t *testing.T) {
	initTest()
	NewDbHandle()

	err := Db.Put(key1, value1)
	assert.Nil(t, err)

	res, err := Db.Has(key1)
	assert.True(t, res)
	assert.Nil(t, err)

	res, err = Db.Has(key2)
	assert.False(t, res)
	assert.Nil(t, err)

	_ = Db.Close()

	res, err = Db.Has(key2)
	assert.False(t, res)
	assert.NotNil(t, err)
}

func TestDelete(t *testing.T) {
	initTest()
	NewDbHandle()

	err := Db.Put(key1, value1)
	assert.Nil(t, err)

	err = Db.Delete(key1)
	assert.Nil(t, err)

	err = Db.Delete(key2)
	assert.Nil(t, err)

	_ = Db.Close()

	err = Db.Delete(key2)
	assert.NotNil(t, err)
}

func TestNewIteratorWithRange(t *testing.T) {
	initTest()
	NewDbHandle()

	err := Db.Put(key1, value1)
	assert.Nil(t, err)

	err = Db.Put(key2, value2)
	assert.Nil(t, err)

	err = Db.Put(key3, value3)
	assert.Nil(t, err)

	err = Db.Put(key4, value4)
	assert.Nil(t, err)

	_, err = Db.NewIteratorWithRange(nil, nil)
	assert.NotNil(t, err)

	res, err := Db.NewIteratorWithRange(key1, key2)
	assert.Nil(t, err)
	count := 0
	for res.Next() {
		count++
	}
	assert.Equal(t, count, 1)

	res, err = Db.NewIteratorWithRange(key1, key4)
	assert.Nil(t, err)
	count = 0
	for res.Next() {
		count++
	}
	assert.Equal(t, count, 3)
}
