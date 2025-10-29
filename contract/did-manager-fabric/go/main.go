/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var DID_PREFIX = "DID_PREFIX"
var ERC20_PREFIX = "ERC20_PREFIX"

// SaveContract provides functions for save data
type CrossChain struct {
	contractapi.Contract
}

// InitLedger adds a base set of cars to the ledger
func (c *CrossChain) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 初始化did文档
	ctx.GetStub().PutState(DID_PREFIX, []byte("{}"))
	// 初始化ERC20
	ctx.GetStub().PutState(ERC20_PREFIX, []byte("{}"))

	return nil
}

// echo
func (c *CrossChain) Echo(ctx contractapi.TransactionContextInterface, content string) (string, error) {
	return content, nil
}

/**
下面是DID逻辑
*/

// 更新用户DID
func (c *CrossChain) UpdateDID(ctx contractapi.TransactionContextInterface, did string, didDoc string) (string, error) {
	didMapBytes, _ := ctx.GetStub().GetState(DID_PREFIX)
	var didMap map[string]string
	json.Unmarshal(didMapBytes, &didMap)

	// 插入用户did
	didMap[did] = didDoc

	// 序列化并保存
	didMapNew, _ := json.Marshal(didMap)
	ctx.GetStub().PutState(DID_PREFIX, didMapNew)

	return "ok", nil
}

// 获取用户DID
func (c *CrossChain) GetDIDDetails(ctx contractapi.TransactionContextInterface, did string) (string, error) {
	didMapBytes, _ := ctx.GetStub().GetState(DID_PREFIX)
	var didMap map[string]string
	json.Unmarshal(didMapBytes, &didMap)

	didDoc, exists := didMap[did]

	if exists {
		return didDoc, nil
	} else {
		return "", nil
	}
}

/**
下面是测试
*/

// 测试
func (c *CrossChain) Busi(ctx contractapi.TransactionContextInterface, vcsJson string) (string, error) {
	// jsonStr := "{\"origin\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"pay\",\"contract_type\":\"1\",\"func_params\":\"[{\\\"string\\\":\\\"c1\\\"},{\\\"string\\\":\\\"121.11\\\"}]\",\"gateway_id\":\"0\",\"param_name\":\"content\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-01T19:26:03.775721138+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"7112e69aadb1c24c9633ba5515da788dc3074df2f8e87b95bd1dbb2dd65269ee411bf0f659d87acccf9b83207a4998204229dc025b45de7662741d3b5d55ec09\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:example:foo#key4\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"contract_type\":\"1\",\"func_params\":\"[{\\\"string\\\":\\\"c1\\\"},{\\\"string\\\":\\\"12.1\\\"}]\",\"gateway_id\":\"0\",\"param_name\":\"content\",\"reource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-01T19:26:03.792124261+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"966b670a1419f36391532ce9e0da50217aef6d4e912cd515b8b04e4fb8d1a2f3927a86ffe70d9932786315e954b274901efcbaf0804e565b14a97d14a8ba300d\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:example:foo#key4\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"

	var vcs map[string]interface{}
	err := json.Unmarshal([]byte(vcsJson), &vcs)
	if err != nil {
		return fmt.Sprintf("%+v", err), nil;
	}

	// 提取VC
	originVc, _ := vcs["origin"].(map[string]interface{})
	// targetVc, _ := vcs["target"].(map[string]interface{})

	// 验签
	originVcResult := c.verifyVc(ctx, originVc)
	// targetVcResult := c.verifyVc(targetVc)
	if (!originVcResult) {
		fmt.Printf("originVcResult: %t\n", originVcResult)
		return "varify fail", nil
	}

	// todo 业务代码

	return "success", nil
}

func (c *CrossChain) getCredentialSubject(ctx contractapi.TransactionContextInterface, vcsJson string) (cs map[string]interface{}, resultErr error) {
	var vcs map[string]interface{}
	json.Unmarshal([]byte(vcsJson), &vcs)

	// 提取VC
	originVc, _ := vcs["origin"].(map[string]interface{})
	targetVc, _ := vcs["target"].(map[string]interface{})
	usedVc := originVc
	if originVc == nil {
		usedVc = targetVc
	}

	// 验签
	usedVcResult := c.verifyVc(ctx, usedVc)
	if (!usedVcResult) {
		fmt.Printf("originVcResult: %t\n", usedVcResult)
		return map[string]interface{}{}, errors.New("verify fail")
	}

	return usedVc["credentialSubject"].(map[string]interface{}), nil
}

func (c *CrossChain) verifyVc(ctx contractapi.TransactionContextInterface, vcs map[string]interface{}) bool {
	// 深度复制
	jsonStr, _ := json.Marshal(vcs)
	var vcsCopy map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &vcsCopy)
	if err != nil {
		return false;
	}

	// 获取签名
	proof, _ := vcsCopy["proof"].(map[string]interface{})
	signatureHex, _ := proof["signature"].(string)
	fmt.Printf("signatureHex: %s\n", signatureHex)
	signatureBytes, _ := hex.DecodeString(signatureHex)
	verificationMethod, _ := proof["verificationMethod"].(string)
	verificationMethodParts := strings.Split(verificationMethod, "#")
	did := verificationMethodParts[0]

	// 查询did文档
	didMapBytes, _ := ctx.GetStub().GetState(DID_PREFIX)
	var didMap map[string]string
	json.Unmarshal(didMapBytes, &didMap)
	didDocStr, exists := didMap[did]
	if !exists {
		didDocStr = "{}"
	}

	// 解析did文档
	var didDoc map[string]interface{}
	json.Unmarshal([]byte(didDocStr), &didDoc)

	// 遍历选择验签方法
	verificationMethodList := didDoc["verificationMethod"].([]interface{})
	publicHex := ""
	for _, verificationMethodItem := range(verificationMethodList) {
		verificationMethodItemReal := verificationMethodItem.(map[string]interface{})
		if verificationMethodItemReal["id"].(string) == verificationMethod {
			publicHex = verificationMethodItemReal["publicKeyHex"].(string)
			break
		}
	}
	// publicHex = "702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a"
	fmt.Printf("publicHex: %s\n", publicHex)
	publicBytes, _ := hex.DecodeString(publicHex)

	// 删除proof字段
	delete(vcsCopy, "proof")

	// 序列化json（1.12之后为按键的 ASCII 码升序序列化）
	vcWithoutProof, _ := json.Marshal(vcsCopy)

	// 验签
	valid := ed25519.Verify(publicBytes, []byte(vcWithoutProof), signatureBytes)
	fmt.Printf("valid: %t\n", valid)

	return valid
}

// 触发跨链保存
func (c *CrossChain) CrossChainSave(ctx contractapi.TransactionContextInterface, key, value string) (string, error) {
	_ = ctx.GetStub().PutState(key, []byte(value))
	_ = ctx.GetStub().PutState(key+"ready", []byte("false"))

	event := fmt.Sprintf("[\"%s\",\"%s\"]", key, value)
	err := ctx.GetStub().SetEvent("test", []byte(event))
	return "success", err
}

func (c *CrossChain) Query(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	value, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err.Error(), nil
	}
	// ready, err := ctx.GetStub().GetState(key + "ready")
	// if err != nil {
	// 	return err.Error(), nil
	// }
	return string(value), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(CrossChain))

	if err != nil {
		fmt.Printf("Error create test chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting test chaincode: %s", err.Error())
	}
}

// func main() {
// 	vcsJson := "{\"origin\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chainmaker001\",\"contract_address\":\"busiWasmCenter3513\",\"contract_func\":\"erc20MintVcs\",\"contract_type\":\"4\",\"func_params\":{\"amount\":\"100\",\"to\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-14T14:53:45.689709257+00:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"784cee1d1c9b2f7059d39c09d75c94061e28b8754af1e7acb76606728b48d7f56b6c92eea05ebe502153f27dab95b058743fb51ac439ce393b3e54528530db06\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"DID_MANAGER\",\"contract_func\":\"Erc20MintVcs\",\"contract_type\":\"2\",\"func_params\":{\"amount\":\"10\",\"to\":\"did:1\"},\"gateway_id\":\"2\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\",\"contentHash\":\"ac58bc50dd7bdd742494a46fc3672809718711e25fc66a819aced14136c0114e\",\"r\":\"302d654c988375ae83a0a6d8c09f0e679babebd85dbc58281263bef952dc9f0a\",\"s\":\"50f0b6cd7b3d3437c4ef688586e6049170fa9a31488d74020d623b32e4ecf697\",\"v\":27,\"verificationMethod\":\"did:1#eth\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}}"

// 	cc := new(CrossChain)
// 	cc.Erc20MintVcs(nil, vcsJson)
// }

// func main() {
// 	vcsJson := "{\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"DID_MANAGER\",\"contract_func\":\"Erc20MintVcs\",\"contract_type\":\"2\",\"func_params\":{\"amount\":\"10\",\"to\":\"did:1\"},\"gateway_id\":\"2\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":null,\"contentHash\":null,\"created\":\"2025-09-14T16:39:11.968247442+00:00\",\"jws\":null,\"proofPurpose\":\"assertionMethod\",\"r\":null,\"s\":null,\"signature\":\"b9f05613a1d5a5c100d734ff353a022f08261011200ba14fbb6ed50a8e6bdd5421cffcca0e16997bf99036fab2b0fd9e74712e77488c6fdbbe225384664f6802\",\"type\":\"JsonWebSignature2020\",\"v\":null,\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}}"

// 	cc := new(CrossChain)
// 	cc.Erc20MintVcs(nil, vcsJson)
// }


// func main() {
// 	oldAmount, err := strconv.ParseInt(string("0"), 10, 32)
// 	fmt.Printf("1111 %d, %+v", oldAmount, err)
// }