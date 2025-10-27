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
	//"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
    	// "github.com/ethereum/go-ethereum/common"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

var ERC20_PREFIX = "ERC20_PREFIX"

// SaveContract provides functions for save data
type CrossChain struct {
	contractapi.Contract
}

// InitLedger adds a base set of cars to the ledger
func (c *CrossChain) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// 初始化ERC20
	ctx.GetStub().PutState(ERC20_PREFIX, []byte("{}"))

	return nil
}

// 测试did-manager
func (c *CrossChain) GetDidDoc(ctx contractapi.TransactionContextInterface, did string) (string, error) {
	ccArgs := [][]byte{
		[]byte("GetDIDDetails"),
		[]byte(did),
	}
	response := ctx.GetStub().InvokeChaincode("DID_MANAGER", ccArgs, "")
	payload := response.Payload
	return string(payload), nil
}

// echo
func (c *CrossChain) Echo(ctx contractapi.TransactionContextInterface, content string) (string, error) {
	return content, nil
}

/**
测试方法
*/

func (c *CrossChain) VerifyVcs(ctx contractapi.TransactionContextInterface, vcsJson string) (string, error) {
	// 解析vc
	credentialSubject, credentialSubjectErr := c.getCredentialSubject(ctx, vcsJson)
	fmt.Printf("credentialSubject %+v", credentialSubject)
	if credentialSubjectErr != nil {
		return "fail", credentialSubjectErr
	}
	return "ok", nil
}

/**
下面是业务方法
*/

// erc20铸造
func (c *CrossChain) Erc20MintVcs(ctx contractapi.TransactionContextInterface, vcsJson string) (string, error) {
	// 解析vc
	credentialSubject, credentialSubjectErr := c.getCredentialSubject(ctx, vcsJson)
	fmt.Printf("credentialSubject %+v", credentialSubject)
	if credentialSubjectErr != nil {
		return "fail", credentialSubjectErr
	}

	// 获取参数
	func_params := credentialSubject["func_params"].(map[string]interface{})
	to := func_params["to"].(string)
	amountStr := func_params["amount"].(string)
	amount, _ := strconv.ParseInt(amountStr, 10, 32)

	// 查询状态
	c.Erc20Mint(ctx, to, amount)

	// 发送事件
	sendCrossChainEventIfNeed(ctx, vcsJson)

	return "ok", nil
}

// erc20铸造
func (c *CrossChain) Erc20Mint(ctx contractapi.TransactionContextInterface, to string, amount int64) (string, error) {
	// 查询状态
	erc20MapBytes, _ := ctx.GetStub().GetState(ERC20_PREFIX)
	var erc20Map map[string]string
	json.Unmarshal(erc20MapBytes, &erc20Map)

	// 查询用户原始账户
	oldAmountStr, exists := erc20Map[to]
	if !exists {
		oldAmountStr = "0"
	}
	oldAmount, _ := strconv.ParseInt(oldAmountStr, 10, 32)

	// 计算新值
	newAmount := oldAmount + amount
	erc20Map[to] = fmt.Sprintf("%d", newAmount)

	// 序列化并保存
	erc20MapNew, _ := json.Marshal(erc20Map)
	ctx.GetStub().PutState(ERC20_PREFIX, erc20MapNew)

	return "ok", nil
}

// erc20转移
func (c *CrossChain) Erc20TransferVcs(ctx contractapi.TransactionContextInterface, vcsJson string) (string, error) {
	// 解析vc
	credentialSubject, _ := c.getCredentialSubject(ctx, vcsJson)

	// 获取参数
	func_params := credentialSubject["func_params"].(map[string]interface{})
	from := func_params["from"].(string)
	to := func_params["to"].(string)
	amountStr := func_params["amount"].(string)
	amount, _ := strconv.ParseInt(amountStr, 10, 32)

	// 查询状态
	erc20MapBytes, _ := ctx.GetStub().GetState(ERC20_PREFIX)
	var erc20Map map[string]string
	json.Unmarshal(erc20MapBytes, &erc20Map)

	// 查询用户原始账户
	oldToAmountStr, toAmountExists := erc20Map[to]
	if !toAmountExists {
		oldToAmountStr = "0"
	}
	oldToAmount, _ := strconv.ParseInt(oldToAmountStr, 10, 32)
	oldFromAmountStr, fromAmountExists := erc20Map[from]
	if !fromAmountExists {
		oldFromAmountStr = "0"
	}
	oldFromAmount, _ := strconv.ParseInt(oldFromAmountStr, 10, 32)

	// 校验
	if oldFromAmount < amount {
		return fmt.Sprintf("error, from account not have enough, oldFromAmount: %d, amount: %d", oldFromAmount, amount), nil
	}

	// 计算新值
	newFromAmount := oldFromAmount - amount
	newToAmount := oldToAmount + amount
	erc20Map[from] = fmt.Sprintf("%d", newFromAmount)
	erc20Map[to] = fmt.Sprintf("%d", newToAmount)

	// 序列化并保存
	erc20MapNew, _ := json.Marshal(erc20Map)
	ctx.GetStub().PutState(ERC20_PREFIX, erc20MapNew)

	// 发送事件
	sendCrossChainEventIfNeed(ctx, vcsJson)

	return "ok", nil
}

// erc20账户余额
func (c *CrossChain) Erc20GetBalanceVcs(ctx contractapi.TransactionContextInterface, vcsJson string) (string, error) {
	// 解析vc
	credentialSubject, _ := c.getCredentialSubject(ctx, vcsJson)

	// 获取参数
	func_params := credentialSubject["func_params"].(map[string]interface{})
	account := func_params["account"].(string)

	// 查询余额
	result, err := c.Erc20GetBalance(ctx, account)

	return result, err
}

// erc20账户余额
func (c *CrossChain) Erc20GetBalance(ctx contractapi.TransactionContextInterface, did string) (string, error) {
	// 查询状态
	erc20MapBytes, _ := ctx.GetStub().GetState(ERC20_PREFIX)
	var erc20Map map[string]string
	json.Unmarshal(erc20MapBytes, &erc20Map)

	// 查询用户原始账户
	oldAmountStr, exists := erc20Map[did]
	if !exists {
		oldAmountStr = "0"
	}

	return oldAmountStr, nil
}

// erc20账户余额
func (c *CrossChain) GetErc20(ctx contractapi.TransactionContextInterface) (string, error) {
	// 查询状态
	erc20MapBytes, _ := ctx.GetStub().GetState(ERC20_PREFIX)

	return string(erc20MapBytes), nil
}

func sendCrossChainEventIfNeed(ctx contractapi.TransactionContextInterface, vcsJson string) {
	var vcs map[string]interface{}
	err := json.Unmarshal([]byte(vcsJson), &vcs)
	if err != nil {
		return
	}

	// 提取VC
	originVc, _ := vcs["origin"].(map[string]interface{})
	targetVc, _ := vcs["target"].(map[string]interface{})

	if originVc != nil && targetVc != nil {
		ccm := map[string]interface{}{
			"target": targetVc,
		}
		ccmStr, _ := json.Marshal(ccm)
		ccmArray := []string{string(ccmStr)}
		ccmArrayStr, _ := json.Marshal(ccmArray)
		ctx.GetStub().SetEvent("CROSS_CHAIN_VC", ccmArrayStr)
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
		return fmt.Sprintf("%+v", err), nil
	}

	// 提取VC
	originVc, _ := vcs["origin"].(map[string]interface{})
	// targetVc, _ := vcs["target"].(map[string]interface{})

	// 验签
	originVcResult := c.verifyVc(ctx, originVc)
	// targetVcResult := c.verifyVc(targetVc)
	if !originVcResult {
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
	if !usedVcResult {
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
		return false
	}

	// 获取签名
	proof, _ := vcsCopy["proof"].(map[string]interface{})
	verificationMethod, _ := proof["verificationMethod"].(string)
	verificationMethodParts := strings.Split(verificationMethod, "#")
	did := verificationMethodParts[0]

	// 跨链码查询did文档
	ccArgs := [][]byte{
		[]byte("GetDIDDetails"),
		[]byte(did),
	}
	response := ctx.GetStub().InvokeChaincode("DID_MANAGER", ccArgs, "")
	payload := response.Payload
	didDocStr := string(payload)
	if didDocStr == "" {
		didDocStr = "{}"
	}

	// 解析did文档
	var didDoc map[string]interface{}
	json.Unmarshal([]byte(didDocStr), &didDoc)

	// 遍历选择验签方法
	verificationMethodList := didDoc["verificationMethod"].([]interface{})
	verificationMethodItemReal := map[string]interface{}{}
	for _, verificationMethodItem := range verificationMethodList {
		verificationMethodItemRealTmp := verificationMethodItem.(map[string]interface{})
		if verificationMethodItemRealTmp["id"].(string) == verificationMethod {
			verificationMethodItemReal = verificationMethodItemRealTmp
			// publicHex = verificationMethodItemReal["publicKeyHex"].(string)
			break
		}
	}
	
	// 验证
	delete(vcsCopy, "proof")
	vcWithoutProof, _ := json.Marshal(vcsCopy)
	valid := false
	if verificationMethodItemReal["type"].(string) == "Ed25519VerificationKey2018" {
		valid = verifyByEd25519(verificationMethodItemReal, proof, string(vcWithoutProof))
	}
	if verificationMethodItemReal["type"].(string) == "EcdsaSecp256k1VerificationKey2019" {
		valid = verifyByEth(verificationMethodItemReal, proof, string(vcWithoutProof))
	}

	return valid
}

func verifyByEd25519(verificationMethodItemReal map[string]interface{}, proof map[string]interface{}, vcWithoutProof string) bool {
	publicHex := verificationMethodItemReal["publicKeyHex"].(string)
	publicBytes, _ := hex.DecodeString(publicHex)

	signatureHex, _ := proof["signature"].(string)
	signatureBytes, _ := hex.DecodeString(signatureHex)

	// 验签
	valid := ed25519.Verify(publicBytes, []byte(vcWithoutProof), signatureBytes)

	return valid
}

func verifyByEth(verificationMethodItemReal map[string]interface{}, proof map[string]interface{}, vcWithoutProof string) bool {
	r := proof["r"].(string)
	s := proof["s"].(string)
	v := int(proof["v"].(float64))
	addressInDid := verificationMethodItemReal["address"].(string)

	address, err := RecoverAddressFromRSV(vcWithoutProof, r, s, v)
	if err != nil {
		return false
	}

	// 比较地址
	addressInDidNormal := strings.ToLower(trimHex(addressInDid))
	addressInDidNormal = string([]byte(addressInDidNormal)[len(addressInDidNormal) - 40:])
	addressNormal := strings.ToLower(trimHex(address))
	if addressInDidNormal == addressNormal {
		return true
	}

	return false
}

func trimHex(hexString string) string {
	if strings.HasPrefix(strings.ToLower(hexString), "0x") {
		return hexString[2:]
	}
	return hexString
}

func RecoverAddressFromRSV(message, rHex, sHex string, v int) (string, error) {
	// 解析r, s, v值
	r, err := hex.DecodeString(rHex)
	if err != nil {
		return "", fmt.Errorf("无效的r值: %v", err)
	}
	s, err := hex.DecodeString(sHex)
	if err != nil {
		return "", fmt.Errorf("无效的s值: %v", err)
	}

	// 处理以太坊消息前缀
	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	// 组合签名
	signature := make([]byte, 65)
	copy(signature[:32], r)
	copy(signature[32:64], s)
	signature[64] = byte((v + 1) % 2)

	// 从签名中恢复公钥
	publicKey, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		return "", fmt.Errorf("恢复公钥失败: %v", err)
	}

	// 转换为公钥地址
	address := crypto.PubkeyToAddress(*publicKey)
	return address.Hex(), nil
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
// 	// 示例用法
// 	message := "123456"

//     // contentHash := "2912723b3ed60c075b271f075d881d82fa5de112b6c25f7dfa4cab85de25045a";
//     // address := "000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea";
	
// 	// 这里替换为实际的RSV值（十六进制）
// 	r := "68c33515e77b9f80404e57512ea112156a1ff058a364f7e96671bf62ca8a63eb"
// 	s := "4ccc92170b52985d35be0478786b08a92d4a81bb5186d69ebe4075a30c48b775"
// 	v := 27 // 通常是0x1b或0x1c
	
// 	address, err := RecoverAddressFromRSV(message, r, s, v)
// 	if err != nil {
// 		fmt.Printf("恢复地址失败: %v\n", err)
// 		return
// 	}
	
// 	fmt.Printf("恢复的地址: %s\n", address)
// }

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
