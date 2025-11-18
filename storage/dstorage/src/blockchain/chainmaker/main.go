/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

/*
#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>

*/
import "C"

import (
	// "C"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	"bytes"
	"errors"
	"encoding/json"
	"os"

	"chainmaker.org/chainmaker/common/v2/evmutils/abi"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	"chainmaker.org/chainmaker/common/v2/crypto"

	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
    ethcommon "github.com/ethereum/go-ethereum/common"
	"gopkg.in/yaml.v3"
)

const (
	createContractTimeout = 5
	contractVersion        = "1.0.0"
)
var users = map[string]*User{}
type User struct {
    TlsKeyPath  string `yaml:"tls_key"`
    TlsCrtPath  string `yaml:"tls_crt"`
    SignKeyPath string `yaml:"sign_key"`
    SignCrtPath string `yaml:"sign_crt"`
}

type UserConfig struct {
    Users map[string]*User `yaml:"users"`
}

var sdkConfigOrg1Client1Path = "/bigdata/drw/storage/dstorage/src/blockchain/chainmaker/chainmaker/sdk_config.yml"

func main() {
	// fileStoreRet := deployContractEVM(sdkConfigOrg1Client1Path, "FileStorechainmaker4", "chainmaker/testdata/solidity/FileStore.bin")
	// fmt.Println(fileStoreRet)
	// // userInfo := invokeGetUser(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/FileStore.abi", "FileStore113", "getUser")
	// // fmt.Println(userInfo)
	// time.Sleep(5 * time.Second)

	
	// storeFileJson := invokeStoreFile(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/FileStore.abi", "FileStorechainmaker4", "storeFile", "aeec10fd-8242-4933-bf3c-32db001a9da5", "config.yaml", 
	// 532, 6, 4, []string{"http://127.0.0.1:5001", "http://127.0.0.1:5002", "http://127.0.0.1:5003", "http://127.0.0.1:5004", "http://127.0.0.1:5005", "http://127.0.0.1:5006"}, 
	// []string{"20410459132366525691697648220546981452656983018310847854076540929744258596752", "18117839898962812797923565871502147358498725488213329967250145524164426259245", "12922240549736046966005255618044387774511305058347430063372002921177566368446", 
	// "11365769016959303865992615604602367532759447157249400316645834812718765714672", "17866157197121340131504842200818932895967418319644693022523263379793733069052", "13668038742036366777173906959202607494519602953874157332009846655339718184455"})
	// fmt.Println(storeFileJson)

	// time.Sleep(5 * time.Second)


	// metadata := invokeFileStoreGetMetaData(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/FileStore.abi", "FileStorechainmaker4", "getFileMetaBasic", "aeec10fd-8242-4933-bf3c-32db001a9da5")
	// info := C.GoString(metadata)
	// fmt.Println("datainfo" + info)

	// ret := deployContractEVM(sdkConfigOrg1Client1Path, "Groth16Verifier112", "chainmaker/testdata/solidity/Groth16Verifier.bin")

	// fmt.Println(ret)

	// challengeAddress := deployContractEVM(sdkConfigOrg1Client1Path, "Challenge112", "chainmaker/testdata/solidity/Challenge.bin")

	// fmt.Println(challengeAddress)

	// // sdkPath string, abiPath string, contractName string, methodName string, address string
	// invokeChallengeInit(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challenge112", "init", "0x" + challengeAddress)

	// challengeId := invokeChallengeCreateChallenge(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challengechainmaker123", "createChallenge", "hash123", "test123", "cid123", "D21D3B948A4A7102F0E317B239EBFB06F8D950A8") 
	// info0 := C.GoString(challengeId)
	// fmt.Println("challengeId" + info0)

	// challengeInfo := invokeChallengeGetChallengesByUser(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challenge112", "getChallengesByUser")
	// fmt.Println("challengeInfo" + challengeInfo)

	// challengeInfo1 := invokeChallengeGetChallengeInfo(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challenge112", "getChallengeInfo", 21)
	// fmt.Println("challengeInfo" + challengeInfo1)
	// // storeFileJson := invokeStoreFile(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/FileStore.abi", "FileStore112", "storeFile", "cid123121", "example.txt", 
	// // 123456, 6, 4, []string{"ipfs://url1", "ipfs://url2", "ipfs://url2", "ipfs://url2", "ipfs://url2", "ipfs://url2"}, []string{"aaa", "bbb", "ccc", "ddd", "eee", "fff"})

	// // fmt.Println(storeFileJson)

	// 原始 JSON 数组字符串，例如 ["0xabc","0xdef"]
	// a := [2]string{"8417222678812533469924558530717241108748831227290678501684031769361826985668", "7582617418287375687851784485095765401224407760144740165472036189630951004519"}
	// b0 := [2]string{"2916128636534116228686483014118820460209874523813207432263458693650976140535", "10041088589092924572695552135236336577617425219240046437204716097746830167134"}
	// b1 := [2]string{"19825149063587233757127400911188256028691625025875966893663672624037668284334", "13806419363301945031604332327411738783524224697845758093049984139947327994997"}
	// c := [2]string{"13852641194570561383356361426255749834612314342438483188282603036290170591562", "5219680289597628299324940061170055425597438332581014145495078585389969430837"}
	// pub := [3]string{"0", "20410459132366525691697648220546981452656983018310847854076540929744258596752", "6541331497782555743317083588695942849605033550442254531334752835795293324955"}
	// challengeInfo1 := invokeChallengeGetChallengeInfo(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challengechainmaker123", "getChallengeInfo", 51)
	// info := C.GoString(challengeInfo1)
	// fmt.Println("challengeInfo" + info)
	
	// pa0 := "7237439946354208475566265764760674782027284456118678774937248367292666579116"
	// pa1 := "16888667091025092877627615507049549034487378360148098234450512032696140633306"
	// pb00 := "4492119549798414799580606400258184014364506224391211685230986049228169627765"
	// pb01 := "21691288014804180048031303967314807391323930071189816329227072094307953449306"
	// pb10 := "1504536702615250616681572387232872850321699446172256305660144146534100312500"
	// pb11 := "2835775901907586247758328335528292026903686963255348937646381903151382061448"
	// pc0 := "14006245309959381963894730208944021677446951661584660731114500688267609429104"
	// pc1 := "11315302815508421488070951551794367385437467530199132738230339990191844005724"
	// pub0 := "0"
	// pub1 := "20410459132366525691697648220546981452656983018310847854076540929744258596752"
	// pub2 := "20410459132366525691697648220546981452656983018310847854076540929744258596752"

	// invokeChallengeSubmitProof(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challengechainmaker123", "submitProof", 51, pa0, pa1, pb00, pb01, pb10, pb11, pc0, pc1, pub0, pub1, pub2);

	// challengeInfo2 := invokeChallengeGetChallengeInfo(sdkConfigOrg1Client1Path, "chainmaker/testdata/solidity/Challenge.abi", "Challengechainmaker123", "getChallengeInfo", 51)
	// info2 := C.GoString(challengeInfo2)
	// fmt.Println("challengeInfo" + info2)
	// admins := []string{"org1admin1", "org2admin1", "org3admin1", "org4admin1"}
	// initConfig("./test.yaml", admins)
}

func CreateChainClientWithSDKConf(sdkConfPath string) (*sdk.ChainClient, error) {
	return sdk.NewChainClient(
		sdk.WithConfPath(sdkConfPath),
	)
}

var usernames []string
//export initConfig
func initConfig(sdkPath string, adminIDs []string) C.bool {
	usernames = adminIDs
	// fmt.Printf("initConfig: %s usernames %+v\n", sdkPath, usernames)
	data, err := os.ReadFile(sdkPath)
    if err != nil {
        return C.bool(false)
    }
	var cfg UserConfig
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return C.bool(false)
    }
    users = cfg.Users
	return C.bool(true)
}

//export deployContractEVM
func deployContractEVM(sdkPath string, contractName string, byteCodePath string) *C.char {
	client, err := CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
		return C.CString("{\"error\":\"failed to read abi file\"}")
	}

	fmt.Printf("====================== 创建%s合约 ======================\n", contractName)
	// usernames := []string{UserNameOrg1Admin1, UserNameOrg2Admin1, UserNameOrg3Admin1, UserNameOrg4Admin1}
	resp, err := createUserContract(client, contractName, contractVersion,
		byteCodePath, common.RuntimeType_EVM, nil, true, usernames...)

	if err != nil {
		fmt.Printf("CREATE EVM contract resp: %+v\n", resp)
		fmt.Println("部署合约失败:", err)
	}

	resp1, err := client.GetContractInfo(contractName)
	if err != nil {
		fmt.Println("查询合约失败:", err)
		return C.CString("{\"error\":\"can not deploy contract\"}")
	} else {
		fmt.Println("合约存在 合约名称 ", resp1.Name, " 版本 ", resp1.Version, "合约地址:", resp1.Address)
		// fmt.Printf("CREATE EVM balance contract resp: %+v\n", resp1)
		return C.CString(resp1.Address)
	}
}

//export invokeStoreFile
func invokeStoreFile(sdkPath string, abiPath string, contractName string, methodName string, rawFileId string, pk string, fileId string, fileName string, author string, fileSize C.longlong) *C.char { 
	input := struct {
		RawFileId	   string
		Pk	           string
		FileId         string
		FileName       string
		Author         string
		FileSize       *big.Int
	}{
		RawFileId:		rawFileId,
		Pk:				pk,
		FileId:         fileId,// "cid123",
		FileName:       fileName, // "example.txt",
		Author:         author, 
		FileSize:       big.NewInt(int64(fileSize)),// big.NewInt(123456),
	}
    fmt.Printf("EVM参数: %+v\n", input)

	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		fmt.Println("加载 ABI 失败:", err)
		// return `{"error":"failed to read abi file"}`
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		fmt.Println("解析 ABI 失败:", err)
		// return `{"error":"failed to parse abi"}`
		return C.CString("")
	}

	// 编码参数
	// "storeFile"
	packed, err := parsedABI.Pack(methodName, &input)
	if err != nil {
		fmt.Println("encodeFileInput error ", err)
		// return `{"error":"failed to pack method arguments"}`
		return C.CString("")
	}

	dataString := hex.EncodeToString(packed)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}


	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		// return `{"error":"contract invocation failed"}`
		return C.CString("")
	} 

	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
	
	str := fmt.Sprint(out)
	return C.CString(str)
}

//export invokeStoreFileExtent
func invokeStoreFileExtent(sdkPath string, abiPath string, contractName string, methodName string, fileId string, totalShards C.longlong, dataShards C.longlong, ipfsUrls []string, cids []string) *C.char { 
	// Marshal 转 JSON
	ipfsUrlsBytes, err := json.Marshal(ipfsUrls)
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}

	ipfsUrlsStr := string(ipfsUrlsBytes)

	cidsBytes, err := json.Marshal(cids)
	if err != nil {
		log.Fatalf("json marshal failed: %v", err)
	}

	cidsStr := string(cidsBytes)

	input := struct {
		FileId         string
		TotalShards    *big.Int
		DataShards     *big.Int
		IpfsUrls       string
		Cids 		   string
	}{
		FileId:         fileId,// "cid123",
		TotalShards:    big.NewInt(int64(totalShards)),// big.NewInt(6),
		DataShards:     big.NewInt(int64(dataShards)),// big.NewInt(4),
		IpfsUrls:       ipfsUrlsStr,// []string{"ipfs://url1", "ipfs://url2", "ipfs://url2", "ipfs://url2", "ipfs://url2", "ipfs://url2"},
		Cids: 			cidsStr,// []string{"0xaaa111...", "0xbbb222...", "0xbbb222...", "0xbbb222...", "0xbbb222...", "0xbbb222..."},
	}
    fmt.Printf("EVM参数: %+v\n", input)

	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		fmt.Println("加载 ABI 失败:", err)
		// return `{"error":"failed to read abi file"}`
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		fmt.Println("解析 ABI 失败:", err)
		// return `{"error":"failed to parse abi"}`
		return C.CString("")
	}

	// 编码参数
	// "storeFile"
	packed, err := parsedABI.Pack(methodName, &input)
	if err != nil {
		fmt.Println("encodeFileInput error ", err)
		// return `{"error":"failed to pack method arguments"}`
		return C.CString("")
	}

	dataString := hex.EncodeToString(packed)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}


	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		// return `{"error":"contract invocation failed"}`
		return C.CString("")
	} 

	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
	
	str := fmt.Sprint(out)
	return C.CString(str)
}

func convertOutToJSONString(out []interface{}) *C.char {
    // 序列化为 JSON
    jsonBytes, err := json.Marshal(out)
    if err != nil {
        fmt.Println("JSON 序列化失败:", err)
        return C.CString("")
    }
    return C.CString(string(jsonBytes))
}

//export invokeHasFile
func invokeHasFile(sdkPath string, abiPath string, contractName string, methodName string, fileId string, fileName string, author string, pk string) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		fmt.Println("读取 ABI 文件失败:", err)
		return C.CString("")
	}

	parsedABI, err := ethabi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		fmt.Println("解析 ABI 失败:", err)
		return C.CString("")
	}

	method := parsedABI.Methods[methodName]
	fmt.Println("方法名称:", method.Name)
	for i, input := range method.Inputs {
		fmt.Printf("输入参数 %d 类型: %s\n", i, input.Type.String())
	}

	dataByte, err := parsedABI.Pack(methodName, fileId, fileName, author, pk)
	if err != nil {
		fmt.Println("打包数据失败:", err)
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	// resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	// if err1 != nil {
	// 	fmt.Println("调用合约失败:", err1)
	// 	return C.CString("")
	// }

	// // 解析 resp（字节数组）为字段
	// resultBytes := resp.ContractResult.Result
	// out, err := parsedABI.Unpack(methodName, resultBytes)
	// if err != nil {
	// 	fmt.Println("调用合约失败:", err)
	// 	return C.CString("")
	// }
	
	// str := fmt.Sprint(out)
	// return C.CString(str)
	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	}

	resultBytes := resp.ContractResult.Result
	if len(resultBytes) == 0 {
		fmt.Println("合约返回空结果")
		return C.CString("")
	}

	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		fmt.Println("ABI 解包失败:", err)
		return C.CString("")
	}

	if len(out) == 0 {
		fmt.Println("ABI 解包结果为空")
		return C.CString("")
	}


	return convertOutToJSONString(out)
}

//export invokeGetUser
func invokeGetUser(sdkPath string, abiPath string, contractName string, methodName string) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		fmt.Println("读取 abi 文件失败:", err)
		// return `{"error":"failed to read abi file"}`
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		fmt.Println("解析 ABI 失败:", err)
		// return `{"error":"failed to parse abi"}`
		return C.CString("")
	}

	// 编码参数
	packed, err := parsedABI.Pack(methodName)
	if err != nil {
		fmt.Println("encodeFileInput error ", err)
		return C.CString("")
		// return `{"error":"failed to pack method arguments"}`
	}

	dataString := hex.EncodeToString(packed)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}


	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		// return `{"error":"contract invocation failed"}`
		return C.CString("")
	} 
	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
	
	str := fmt.Sprint(out)
	return C.CString(str)
}

//export invokeChallengeInit
func invokeChallengeInit(sdkPath string, abiPath string, contractName string, methodName string, address string) bool { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		fmt.Println("读取 ABI 文件失败:", err)
		return false
	}

	parsedABI, err := ethabi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		fmt.Println("解析 ABI 失败:", err)
		return false
	}

	method := parsedABI.Methods[methodName]
	fmt.Println("方法名称:", method.Name)
	for i, input := range method.Inputs {
		fmt.Printf("输入参数 %d 类型: %s\n", i, input.Type.String())
	}

	verifierAddress := ethcommon.HexToAddress(address)
	fmt.Printf("获取地址成功，返回 verifierAddress %+v\n", verifierAddress)

	// addr := common.HexToAddress(address)
	// addr := evmutils.BigToAddress(evmutils.FromFromString(address))
	dataByte, err := parsedABI.Pack(methodName, verifierAddress)
	if err != nil {
		fmt.Println("打包数据失败:", err)
		return false
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return false
	}

	fmt.Printf("调用成功，返回 resp %+v\n", resp)
	return true
}

func goStringToCCharArray(goStr string, cStr **C.char) {
    *cStr = C.CString(goStr)
}

//export invokeChallengeGetChallengesByUser
func invokeChallengeGetChallengesByUser(sdkPath string, abiPath string, contractName string, methodName string) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return C.CString("")
	}

	dataByte, err := parsedABI.Pack(methodName)
	if err != nil {
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	} 

	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
	
	str := fmt.Sprint(out)
	// 去掉前后 []
	trimmed := strings.Trim(str, "[]")

	// 按空格分割
	parts := strings.Fields(trimmed)

	jsonData, err := json.Marshal(parts)
	if err != nil {
		panic(err)
	}
	return C.CString(string(jsonData))
}

//export invokeChallengeGetChallengeInfo
func invokeChallengeGetChallengeInfo(sdkPath string, abiPath string, contractName string, methodName string, challengeId uint64) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return C.CString("")
	}

	dataByte, err := parsedABI.Pack(methodName, big.NewInt(int64(challengeId)))
	if err != nil {
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	} 

	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
	
	str := fmt.Sprint(out)
	// 去掉前后 []
	trimmed := strings.Trim(str, "[]")

	// 按空格分割
	parts := strings.Fields(trimmed)
	
	jsonData, err := json.Marshal(parts)
	if err != nil {
		panic(err)
	}
	return C.CString(string(jsonData))
}

//export invokeFileStoreGetMetaData
func invokeFileStoreGetMetaData(sdkPath string, abiPath string, contractName string, methodName string, fileId string, pk string) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		return C.CString("")
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return C.CString("")
	}

	dataByte, err := parsedABI.Pack(methodName, fileId, pk)
	if err != nil {
		fmt.Println("调用合约失败1:", err)
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	} 


	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		fmt.Println("error")

		return C.CString("")
	}

	fmt.Println("resp", out)

	result := make([]interface{}, len(out))

	for i, v := range out {
        switch val := v.(type) {
        case string:
            // 尝试解析成 JSON，如果成功就转成数组或对象
            var temp interface{}
            if err := json.Unmarshal([]byte(val), &temp); err == nil {
                result[i] = temp
            } else {
                result[i] = val
            }
        default:
            // 保留原始类型（数字、[]interface{}等）
            result[i] = val
        }
    }

    finalJSON, err := json.MarshalIndent(result, "", "  ")
    if err != nil {
        panic(err)
    }
	return C.CString(string(finalJSON))
}

func parseSolidityStringArray(data []byte) ([]string, error) {
    if len(data) < 32 {
        return nil, errors.New("data too short")
    }

    // 数组长度
    arrLen := new(big.Int).SetBytes(data[:32]).Uint64()
    if arrLen > 1000000 {
        return nil, errors.New("array too long")
    }

    result := make([]string, arrLen)
    offsets := make([]uint64, arrLen)

    // 读取每个字符串偏移
    for i := uint64(0); i < arrLen; i++ {
        start := 32 + i*32
        if start+32 > uint64(len(data)) {
            return nil, errors.New("unexpected EOF reading offsets")
        }
        offsets[i] = new(big.Int).SetBytes(data[start : start+32]).Uint64()
    }

    for i := uint64(0); i < arrLen; i++ {
        off := offsets[i]
        if off+32 > uint64(len(data)) {
            return nil, errors.New("unexpected EOF reading string length")
        }
        strLen := new(big.Int).SetBytes(data[off : off+32]).Uint64()
        start := off + 32
        end := start + strLen
        if end > uint64(len(data)) {
            return nil, errors.New("unexpected EOF reading string data")
        }
        result[i] = string(data[start:end])
    }

    return result, nil
}


//export invokeFileStoreGetArray
func invokeFileStoreGetArray(sdkPath string, abiPath string, contractName string, methodName string, fileId string, pk string) *C.char { 
	// 加载 abi
	fmt.Println("abi path : ", abiPath)
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		return C.CString("")
	}

	parsedABI, err := ethabi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return C.CString("")
	}

	dataByte, err := parsedABI.Pack(methodName, fileId, pk)
	if err != nil {
		fmt.Println("调用合约失败1:", err)
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	} 

	var urls []string
	err2 := parsedABI.UnpackIntoInterface(&urls, "getFileMetaIpfsUrls", resp.ContractResult.Result)
	if err2 != nil {
		fmt.Println("ABI Unpack error:", err2)
		return C.CString("")
	}
	
	// 去掉每个字符串的多余引号
	clean := make([]string, len(urls))
	for i, v := range urls {
		clean[i] = strings.Trim(v, `"`)
	}
	// 包装成二维数组
    wrapped := [][]string{clean}

    // 序列化 JSON
    jsonBytes, err := json.Marshal(wrapped)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(jsonBytes))
	return C.CString(string(jsonBytes))
	// out, err := parsedABI.Unpack(methodName, resultBytes)
	// if err != nil {
	// 	fmt.Println("error", err)

	// 	return C.CString("")
	// }

	// fmt.Println("resp", out)

	// result := make([]interface{}, len(out))

	// for i, v := range out {
    //     switch val := v.(type) {
    //     case string:
    //         // 尝试解析成 JSON，如果成功就转成数组或对象
    //         var temp interface{}
    //         if err := json.Unmarshal([]byte(val), &temp); err == nil {
    //             result[i] = temp
    //         } else {
    //             result[i] = val
    //         }
    //     default:
    //         // 保留原始类型（数字、[]interface{}等）
    //         result[i] = val
    //     }
    // }

    // finalJSON, err := json.MarshalIndent(result, "", "  ")
    // if err != nil {
    //     panic(err)
    // }
	// return C.CString(string(finalJSON))
}

// 工具函数：将 uint64 C 数组转为 [N]*big.Int
func convertToBigIntArray(ptr *string, size int) []*big.Int {
	if ptr == nil {
		return nil
	}
	str := *ptr
	// 支持逗号或空格分隔
	tokens := strings.FieldsFunc(str, func(r rune) bool {
		return r == ' ' || r == ','
	})

	result := []*big.Int{}
	for i, token := range tokens {
		if size > 0 && i >= size {
			break
		}
		n := new(big.Int)
		n, ok := n.SetString(token, 10)
		if !ok {
			// 解析失败，跳过或你可以 log.Fatal
			continue
		}
		result = append(result, n)
	}
	return result
}

func convertStringToBigInt(str string) *big.Int {
	n := new(big.Int)
	n, ok := n.SetString(str, 10)
	if !ok {
		// 解析失败，你可以 log.Fatal
		return nil
	}
	return n
}

//export invokeChallengeSubmitProof
func invokeChallengeSubmitProof(sdkPath string, abiPath string, contractName string, methodName string, challengeId uint64, pA0 string, pA1 string,pB00 string,pB01 string,pB10 string,pB11 string,pC0 string,pC1 string,pubSignals0 string,pubSignals1 string,pubSignals2 string) bool { 
	// 转参数
	input := SubmitProofInput{
		ChallengeId: uint64(challengeId),
		PA:          [2]*big.Int{},
		PB0:         [2]*big.Int{},
		PB1:         [2]*big.Int{},
		PC:          [2]*big.Int{},
		PubSignals:  [3]*big.Int{},
	}
	input.PA[0] = convertStringToBigInt(pA0)
	input.PA[1] = convertStringToBigInt(pA1)
	input.PB0[0] = convertStringToBigInt(pB00)
	input.PB0[1] = convertStringToBigInt(pB01)
	input.PB1[0] = convertStringToBigInt(pB10)
	input.PB1[1] = convertStringToBigInt(pB11)
	input.PC[0] = convertStringToBigInt(pC0)
	input.PC[1] = convertStringToBigInt(pC1)
	input.PubSignals[0] = convertStringToBigInt(pubSignals0)
	input.PubSignals[1] = convertStringToBigInt(pubSignals1)
	input.PubSignals[2] = convertStringToBigInt(pubSignals2)
	// copy(input.PA[:], convertToBigIntArray(pA, 2))
	// copy(input.PB0[:], convertToBigIntArray(pB0, 2))
	// copy(input.PB1[:], convertToBigIntArray(pB1, 2))
	// copy(input.PC[:], convertToBigIntArray(pC, 2))
	// copy(input.PubSignals[:], convertToBigIntArray(pubSignals, 3))

	fmt.Printf("input %+v\n", input)
	// 执行调用
	ret, err := submitProofCall(sdkPath, abiPath, contractName, methodName, input)
	if err != nil {
		fmt.Println("submitProof error:", err)
		return ret
	}
	return ret
}


// 原始业务逻辑封装
type SubmitProofInput struct {
	ChallengeId uint64
	PA          [2]*big.Int
	PB0         [2]*big.Int
	PB1         [2]*big.Int
	PC          [2]*big.Int
	PubSignals  [3]*big.Int
}

func submitProofCall(sdkPath, abiPath, contractName string, methodName string, input SubmitProofInput) (bool, error) {
	// 读取 ABI
	abiJson, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return false, err
	}
	parsedABI, err := abi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return false, err
	}

	// Pack calldata
	packed, err := parsedABI.Pack(methodName,
		big.NewInt(int64(input.ChallengeId)),
		input.PA,
		input.PB0,
		input.PB1,
		input.PC,
		input.PubSignals,
	)
	if err != nil {
		return false, err
	}
	dataHex := hex.EncodeToString(packed)

	// 构造 KVP
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataHex),
		},
	}

	// 发送交易
	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Printf("resp err: %+v", resp)
		return false, err1
	}

	// 取合约返回值（二进制）
	resultBytes := resp.ContractResult.Result
	// 使用 ABI 解码返回值
	unpacked, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		fmt.Printf("unpack err: %+v", unpacked)
		return false, err
	}

	return true, nil
}

//export invokeChallengeCreateChallenge
func invokeChallengeCreateChallenge(sdkPath string, abiPath string, contractName string, methodName string, fileHash string, nonce string, cid string, prover string) *C.char { 
	// 加载 abi
	abiJson, err := ioutil.ReadFile(abiPath) // 你合约编译生成的 ABI 文件
	if err != nil {
		return C.CString("")
	}

	parsedABI, err := ethabi.JSON(bytes.NewReader(abiJson))
	if err != nil {
		return C.CString("")
	}

	proverAddr := ethcommon.HexToAddress(prover)
	fmt.Printf("获取地址成功，返回 proverAddr %+v\n", proverAddr)
	// addr := evmutils.BigToAddress(evmutils.FromDecimalString(prover))
	dataByte, err := parsedABI.Pack(methodName, fileHash, nonce, cid, proverAddr)
	if err != nil {
		fmt.Println("打包数据失败:", err)
		return C.CString("")
	}

	dataString := hex.EncodeToString(dataByte)
	kvs := []*common.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(dataString),
		},
	}

	resp, err1 := invokeEvmMethod(sdkPath, contractName, methodName, kvs, true)
	if err1 != nil {
		fmt.Println("调用合约失败:", err1)
		return C.CString("")
	} 

	// 解析 resp（字节数组）为字段
	resultBytes := resp.ContractResult.Result
	out, err := parsedABI.Unpack(methodName, resultBytes)
	if err != nil {
		return C.CString("")
	}
		
	str := fmt.Sprint(out)
	return C.CString(str)
}

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.Message == "" {
			resp.Message = resp.Code.String()
		}
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}

func GetEndorsersWithAuthType(hashType crypto.HashType, authType sdk.AuthType, payload *common.Payload, usernames ...string) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, name := range usernames {
		var entry *common.EndorsementEntry
		var err error
		switch authType {
		case sdk.PermissionedWithCert:
			u, ok := users[name]
			if !ok {
				return nil, errors.New("user not found")
			}
			if sdk.GetP11Handle() != nil || sdk.KMSEnabled() {
				entry, err = sdkutils.MakeEndorserWithPathAndP11Handle(u.SignKeyPath, u.SignCrtPath, sdk.GetP11Handle(),
					sdk.KMSEnabled(), payload)
				if err != nil {
					return nil, err
				}
			} else {
				entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, errors.New("invalid authType")
		}
		endorsers = append(endorsers, entry)
	}

	return endorsers, nil
}

func createUserContract(client *sdk.ChainClient, contractName, version, byteCodePath string,
	runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool, usernames ...string) (*common.TxResponse, error) {

	payload, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
	if err != nil {
		return nil, err
	}

	endorsers, err := GetEndorsersWithAuthType(client.GetHashType(),
		client.GetAuthType(), payload, usernames...)
	if err != nil {
		return nil, err
	}

	// 发送创建合约请求
	resp, err := client.SendContractManageRequest(payload, endorsers, createContractTimeout, withSyncResult)
	if err != nil {
		return nil, err
	}

	err = CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func invokeEvmMethod(sdkPath string, contractName string, methodName string, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {
	fmt.Println("====================== create client ======================")
	client, err := CreateChainClientWithSDKConf(sdkPath)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("====================== 调用%s合约 %s方法 ======================\n", contractName, methodName)
	return invokeUserContract(client, contractName, methodName, "", kvs, withSyncResult)
}


func invokeUserContract(client *sdk.ChainClient, contractName, method, txId string, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

	resp, err := client.InvokeContract(contractName, method, txId, kvs, -1, withSyncResult)
	if err != nil {
		return resp, err
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		return resp, fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	return resp, nil
}

