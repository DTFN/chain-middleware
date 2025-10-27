/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package restrequest

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"go.uber.org/zap"
)

func initTest() *zap.SugaredLogger {
	log := []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	logger.InitLogConfig(log)
	return logger.GetLogger(logger.ModuleRequest)
}

func TestRestRequest_BeginCrossChain(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	log := initTest()
	req := NewRestRequest(log)
	_, _ = req.BeginCrossChain(nil)
}

func TestRestRequest_InitSpvContract(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	log := initTest()
	req := NewRestRequest(log)
	_, _ = req.InitSpvContract(nil)
}

func TestRestRequest_SyncBlockHeader(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	log := initTest()
	req := NewRestRequest(log)
	_, _ = req.SyncBlockHeader(nil)
}

func TestRestRequest_UpdateSpvContract(t *testing.T) {
	defer func() {
		err := recover()
		assert.NotNil(t, err)
	}()
	log := initTest()
	req := NewRestRequest(log)
	_, _ = req.UpdateSpvContract(nil)
}
