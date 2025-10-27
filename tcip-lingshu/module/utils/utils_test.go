/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	test = "test"
)

func TestCreateCertainTempFile(t *testing.T) {
	_, err := CreateCertainTempFile(os.TempDir(), test)
	assert.Nil(t, err)

	_, err = CreateCertainTempFile(path.Join(os.TempDir(), time.Now().String(), "test1"), test)
	assert.Nil(t, err)
}

func TestCreateTempFile(t *testing.T) {
	_, err := CreateTempFile("123")
	assert.Nil(t, err)
}

func TestDeepCopy(t *testing.T) {
	type testStruct struct {
		A int32
		B int32
	}
	testS := testStruct{
		A: 10,
		B: 20,
	}
	testS1 := testStruct{}
	err := DeepCopy(&testS1, testS)
	assert.Nil(t, err)
	assert.Equal(t, testS, testS1)
}

func TestWriteTempFile(t *testing.T) {
	_, err := WriteTempFile([]byte("123"), "test123")
	assert.Nil(t, err)
}
