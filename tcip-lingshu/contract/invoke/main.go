package main

import (
	"fmt"
	"github.com/FISCO-BCOS/go-sdk/abi"
	"github.com/FISCO-BCOS/go-sdk/abi/bind"
	"github.com/FISCO-BCOS/go-sdk/client"
	"github.com/FISCO-BCOS/go-sdk/conf"
	"github.com/FISCO-BCOS/go-sdk/core/types"
	"github.com/ethereum/go-ethereum/common"
	"os"
	"strings"
)

const (
	cabi = "[{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"CrossChainCancel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"CrossChainConfirm\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"string\"}],\"name\":\"CrossChainSave\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"value\",\"type\":\"string\"}],\"name\":\"CrossChainTry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"query\",\"outputs\":[{\"name\":\"value\",\"type\":\"string\"},{\"name\":\"flag\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"string\"}],\"name\":\"Test\",\"type\":\"event\"}]"
)

func main()  {
	invokeConfigs, err := conf.ParseConfigFile("sdk_config1.toml")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	invokeConfig := &invokeConfigs[0]
	invokeClient, err := client.Dial(invokeConfig)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}

	invokeParsed, err := abi.JSON(strings.NewReader(cabi))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return
	}
	_, receipt, err := bind.NewBoundContract(common.HexToAddress(os.Args[1]), invokeParsed, invokeClient, invokeClient, invokeClient).Transact(invokeClient.GetTransactOpts(),
		"CrossChainSave", os.Args[2], os.Args[3])

	if err != nil {
		fmt.Printf("invoke error\n")
		return
	}
	if receipt.Status != types.Success {
		fmt.Printf("invoke error state %v\n", receipt)
		return
	}
}