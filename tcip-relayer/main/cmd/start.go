/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/server"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/rpcserver"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/common/v2/json"

	_ "net/http/pprof"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"github.com/spf13/cobra"
)

var cliLog *zap.SugaredLogger

// StartCMD 启动命令
func StartCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Startup tcip-relayer",
		Long:  "Startup tcip-relayer",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			mainStart()
			fmt.Println("tcip-relayer exit")
			return nil
		},
	}
	attachFlags(startCmd, []string{flagNameOfConfigFilepath})
	return startCmd
}

func mainStart() {
	cliLog = logger.GetLogger(logger.ModuleCli)
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

	printLogo()

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

func printLogo() {
	cliLog.Infof(logo())
}
