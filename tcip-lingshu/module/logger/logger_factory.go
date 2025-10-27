/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package logger

import (
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	// ModuleDefault 默认模块
	ModuleDefault = "[DEFAULT]"
	// ModuleStart 启动模块
	ModuleStart = "[START]"
	// ModuleRegister 注册模块
	ModuleRegister = "[REGISTER]"
	// ModuleRequest 请求管理模块模块
	ModuleRequest = "[REQUEST_MANAGER]"
	// ModuleChainmakerClient chainmaker链交互模块
	ModuleChainmakerClient = "[CHAINMAKER_CLIENT]"
	// ModuleRpcServer rpc模块
	ModuleRpcServer = "[RPC_SERVER]"
	// ModuleChainClient 链交互模块
	ModuleChainClient = "[CHAIN_CLIENT]"
	// ModuleHandler handler模块
	ModuleHandler = "[HANDLER]"
	// ModuleDb 数据存储模块
	ModuleDb = "[DB]"
	// ModuleChainConfig 链配置模块
	ModuleChainConfig

	defaultLogPath = "./logs/default.log" // release struct need this path
)

var (
	defaultLogConfig *Config
	loggers          = make(map[string]*zap.SugaredLogger)
	loggerMutex      sync.Mutex
	logInitialized   = false
)

// InitLogConfig set the config of logger module, called in initialization of config module
//  @param config
func InitLogConfig(config []*LogModuleConfig) {
	// 初始化loggers
	for _, logModuleConfig := range config {
		logPrintName := logPrintName(logModuleConfig.ModuleName)
		config := &Config{
			Module:       logPrintName,
			LogPath:      logModuleConfig.FilePath,
			LogLevel:     GetLogLevel(logModuleConfig.LogLevel),
			MaxAge:       logModuleConfig.MaxAge,
			RotationTime: logModuleConfig.RotationTime,
			JsonFormat:   false,
			ShowLine:     true,
			LogInConsole: logModuleConfig.LogInConsole,
			ShowColor:    logModuleConfig.ShowColor,
		}
		logger, _ := InitSugarLogger(config)
		loggers[logPrintName] = logger
	}
	// 最后添加"ModuleDefault"
	if _, exist := loggers[ModuleDefault]; !exist {
		// 创建默认的logger
		loggers[ModuleDefault] = getLogDefaultModuleConfig()
	}
	logInitialized = true
}

// GetLogger return the instance of SugaredLogger
//  @param name
//  @return *zap.SugaredLogger
func GetLogger(name string) *zap.SugaredLogger {
	if !logInitialized {
		panic("log has not been initialized")
	}
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	logHeader := name
	logger, ok := loggers[logHeader]
	if !ok {
		logger = copyLogger(name)
		loggers[name] = logger
	}
	return logger
}

func copyLogger(module string) *zap.SugaredLogger {
	defaultLogger := loggers[ModuleDefault]
	return zap.New(defaultLogger.Desugar().Core()).Named(module).WithOptions(zap.AddCaller()).Sugar()
}

func getLogDefaultModuleConfig() *zap.SugaredLogger {
	if defaultLogConfig == nil {
		defaultLogConfig = &Config{
			Module:       ModuleDefault,
			LogPath:      defaultLogPath,
			LogLevel:     LEVEL_INFO,
			MaxAge:       DEFAULT_MAX_AGE,
			RotationTime: DEFAULT_ROTATION_TIME,
			JsonFormat:   false,
			ShowLine:     true,
			LogInConsole: true,
			ShowColor:    true,
		}
		logger, _ := InitSugarLogger(defaultLogConfig)
		return logger
	}
	logger, _ := InitSugarLogger(defaultLogConfig)
	return logger
}

//func getLogModuleConfig(moduleName string) *zap.SugaredLogger {
//	innerLogConfig := &Config{
//		Module:       moduleName,
//		LogPath:      defaultLogPath,
//		LogLevel:     LEVEL_INFO,
//		MaxAge:       DEFAULT_MAX_AGE,
//		RotationTime: DEFAULT_ROTATION_TIME,
//		JsonFormat:   false,
//		ShowLine:     true,
//		LogInConsole: true,
//		ShowColor:    true,
//	}
//	logger, _ := InitSugarLogger(innerLogConfig)
//	return logger
//}

func logPrintName(moduleName string) string {
	if moduleName == "" {
		return ModuleDefault
	}
	return "[" + strings.ToUpper(moduleName) + "]"
}
