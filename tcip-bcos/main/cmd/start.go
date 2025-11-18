/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/rpcserver"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/server"

	"github.com/spf13/pflag"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cliLog *zap.SugaredLogger

// StartCMD 启动命令
//  @return *cobra.Command
func StartCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Startup tcip-bcos",
		Long:  "Startup tcip-bcos",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			mainStart()
			fmt.Println("tcip-bcos exit")
			return nil
		},
	}
	startAttachFlags(startCmd, []string{flagNameOfConfigFilepath})
	return startCmd
}

func mainStart() {
	cliLog = logger.GetLogger(logger.ModuleStart)
	config, _ := json.Marshal(conf.Config)
	cliLog.Debug(string(config))

	rpcServer, err := rpcserver.NewRpcServer()
	if err != nil {
		cliLog.Errorf("rpc server init failed, %s", err.Error())
		return
	}

	// new an error channel to receive errors
	errorC := make(chan error, 1)

	server.InitServer(errorC)

	// start rpc server and listen in another go routine
	err = rpcServer.Start()
	if err != nil {
		errorC <- err
	}

	cliLog.Infof(logo())

	// handle exit signal in separate go routines
	go handleExitSignal(errorC)

	// listen error signal in main function
	err = <-errorC
	if err != nil {
		cliLog.Error("server encounters error ", err)
	}
	rpcServer.Stop()
	cliLog.Info("All is stopped!")
}

func startFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilePath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath,
		conf.ConfigFilePath, "specify config file path, if not set, default use ./tcip_bcos.yml")
	return flags
}

func startAttachFlags(cmd *cobra.Command, flagNames []string) {
	flags := startFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}

// handleExitSignal listen exit signal for process stop
func handleExitSignal(exitC chan<- error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	defer signal.Stop(signalChan)

	for sig := range signalChan {
		cliLog.Infof("received exit signal: %d (%s)", sig, sig)
		exitC <- nil
	}
}
