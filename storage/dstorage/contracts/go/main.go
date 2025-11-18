/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/v2/contractapi"
	"github.com/hyperledger/fabric-samples/asset-transfer-basic/chaincode-go/chaincode"
)

func main() {
	// 创建两个合约实例
	challengeContract := new(chaincode.ChallengeContract)
	filestoreContract := new(chaincode.FileStoreContract)
	
	chaincode_contract, err := contractapi.NewChaincode(challengeContract, filestoreContract)
	if err != nil {
		log.Panicf("Error creating contract chaincode: %v", err)
	}

	if err := chaincode_contract.Start(); err != nil {
		log.Panicf("Error starting contract chaincode: %v", err)
	}
}
