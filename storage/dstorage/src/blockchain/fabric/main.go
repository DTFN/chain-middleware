package main
import "C"

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"time"
	"strconv"
	"encoding/json"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	mspID        = "Org1MSP"
	cryptoPath   = "/bigdata/drw/storage/fabric/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/Admin@org1.example.com/msp/signcerts/cert.pem"
	keyPath      = cryptoPath + "/users/Admin@org1.example.com/msp/keystore/"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
	channelName  = "mychannel"
	chaincodeName = "my_new_storage"
)

func main() {
	// ret := invokeFabricStoreFile("FileStoreContract", "StoreFile", "cid1234571", "sample.pdf", 123456, 4, 3, []string{"ipfs1", "ipfs2", "ipfs3", "ipfs4"}, []string{"cid1", "cid2", "cid3", "cid4"})
	// fmt.Println(ret)

	// time.Sleep(5 * time.Second)

	// ret1 := invokeFabricGetFileMetaData("FileStoreContract", "GetFileMeta", "cid1234571")
	// fmt.Println(ret1)

	// // ret := invokeFabricStoreFile("FileStoreContract", "StoreFile", "cid1234566", "sample.pdf", 123456, 4, 3, []string{"ipfs1", "ipfs2", "ipfs3", "ipfs4"}, []string{"hash1", "hash2", "hash3", "hash4"})
	// // fmt.Println(ret)
	// ret2 := invokeFabricGetUser("FileStoreContract", "GetUser")
	// fmt.Println(ret2)
	// ret := invokeFabricChallengeCreateChallenge("ChallengeContract", "CreateChallenge", "cid123456", "nonce123456", "cid123456", "eDUwOTo6Q049QWRtaW5Ab3JnMS5leGFtcGxlLmNvbSxPVT1hZG1pbixMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVTOjpDTj1jYS5vcmcxLmV4YW1wbGUuY29tLE89b3JnMS5leGFtcGxlLmNvbSxMPVNhbiBGcmFuY2lzY28sU1Q9Q2FsaWZvcm5pYSxDPVVT")
	// fmt.Println(ret)

	// ret := invokeChallengeGetChallengesByUser("ChallengeContract", "GetMyChallenges")
	// fmt.Println(ret)

	// ret := invokeChallengeGetChallengeInfo("ChallengeContract", "GetChallenge", 1)
	// fmt.Println(ret)

	// ret1 := invokeChallengeSubmitProof("ChallengeContract", "SubmitProof", 1, "proofCID", "proofURL", "proofHash", "publicCID", "publicURL", "publicHash")
	// fmt.Println(ret1)
	
	// ret2 := invokeChallengeGetProof("ChallengeContract", "GetProof", 1)
	// fmt.Println(ret2)

	// ret3 := invokeChallengeUpdateChallenge("ChallengeContract", "UpdateChallenge", 1, true)
	// fmt.Println(ret3)

	// ret4:= invokeChallengeGetChallengeInfo("ChallengeContract", "GetChallenge", 1)
	// fmt.Println(ret4)
}

//export initFabricConfig
func initFabricConfig(m string, c string, p string, g string, cn string, ccn string, cp string, kp string, tlscp string) {
	// 设置全局变量
	mspID = m
	cryptoPath = c
	certPath     = cryptoPath + cp
	keyPath      = cryptoPath + kp
	tlsCertPath  = cryptoPath + tlscp
	peerEndpoint = p
	gatewayPeer = g
	channelName = cn
	chaincodeName = ccn
}


//export invokeFabricStoreFile
func invokeFabricStoreFile(contractName string, methodName string, rawFileId, pk, fileId string, fileName, author string, fileSize C.longlong, totalShards C.longlong, dataShards C.longlong, ipfsUrls []string, cids []string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	fileSizeStr := strconv.FormatInt(int64(fileSize), 10)
	totalShardsStr := strconv.FormatInt(int64(totalShards), 10)
	dataShardsStr := strconv.FormatInt(int64(dataShards), 10)
	ipfsUrlsStr, _ := json.Marshal(ipfsUrls) // 或 JSON 序列化
	cidsStr, _ := json.Marshal(cids)

	// 3. 准备调用参数
	args := []string{
		rawFileId,
		pk,
		fileId,
		fileName,
		author,
		fileSizeStr,
		totalShardsStr,
		dataShardsStr,
		string(ipfsUrlsStr),
		string(cidsStr),
	}

	fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: FileStoreContract:StoreFile")
	result, err := contract.SubmitTransaction(realName, args...)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricHasFile
func invokeFabricHasFile(contractName string, methodName string, fileId, fileName, author, pk string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)


	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: FileStoreContract:HasFile")
	respBytes, err := contract.EvaluateTransaction(realName, fileId, fileName, author, pk)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}

	fmt.Printf("Raw bytes: %s\n", string(respBytes))
	// 链码返回的是 JSON，你需要解析
	// 尝试先当作 JSON 字符串解析
	var inner string
	if err := json.Unmarshal(respBytes, &inner); err == nil {
		// 内层是字符串，再解一次
		respBytes = []byte(inner)
	}

	// 解析成结构体
	var result struct {
		IDs    []string `json:"ids"`
		Access []bool   `json:"access"`
	}
	if err := json.Unmarshal(respBytes, &result); err != nil {
		log.Fatalf("Failed to unmarshal result: %v", err)
	}

	// 组装成二维数组 [[ids], [access]]
	final := []interface{}{result.IDs, result.Access}

	// 输出 JSON
	jsonBytes, err := json.Marshal(final)
	if err != nil {
		log.Fatalf("JSON marshal failed: %v", err)
	}

	fmt.Println(string(jsonBytes)) // 输出 [["a","b","c"],[true,true,true]]

	return C.CString(string(jsonBytes))
}

//export invokeFabricGetFileMetaData
func invokeFabricGetFileMetaData(contractName string, methodName string, fileId, pk string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)


	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: FileStoreContract:GetFileMeta")
	result, err := contract.EvaluateTransaction(realName, fileId, pk)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricGetUser
func invokeFabricGetUser(contractName string, methodName string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)
	// 3. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: FileStoreContract:GetUser")
	result, err := contract.SubmitTransaction(realName)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeGetChallengesByUser
func invokeFabricChallengeGetChallengesByUser(contractName string, methodName string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:GetMyChallenges")
	result, err := contract.EvaluateTransaction(realName)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeGetChallengeInfo
func invokeFabricChallengeGetChallengeInfo(contractName string, methodName string, challengeId uint64) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)


	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:GetChallenge")
	result, err := contract.EvaluateTransaction(realName, strconv.FormatInt(int64(challengeId), 10))
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeSubmitProof
func invokeFabricChallengeSubmitProof(contractName string, methodName string, challengeId uint64, 
	proofCID, proofURL, proofHash string,
	publicCID, publicURL, publicHash string,
	) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	challengeIdStr := strconv.FormatInt(int64(challengeId), 10)
	// 3. 准备调用参数
	args := []string{
		challengeIdStr,
		proofCID,
		proofURL,
		proofHash,
		publicCID,
		publicURL,
		publicHash,
		"",
		"",
		"",
	}

	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:SubmitProof")
	result, err := contract.SubmitTransaction(realName, args...)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeGetProof
func invokeFabricChallengeGetProof(contractName string, methodName string, challengeId uint64) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:GetProof")
	result, err := contract.SubmitTransaction(realName, strconv.FormatInt(int64(challengeId), 10))
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeUpdateChallenge
func invokeFabricChallengeUpdateChallenge(contractName string, methodName string, challengeId uint64, ProofVerified bool) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:UpdateChallenge")
	result, err := contract.SubmitTransaction(realName, strconv.FormatInt(int64(challengeId), 10), strconv.FormatBool(ProofVerified))
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricChallengeCreateChallenge
func invokeFabricChallengeCreateChallenge(contractName string, methodName string, fileHash string, nonce string, cid string, prover string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// 3. 准备调用参数
	args := []string{
		fileHash,
		nonce,
		cid,
		prover,
	}

	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: ChallengeContract:CreateChallenge")
	result, err := contract.SubmitTransaction(realName, args...)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

//export invokeFabricFileStoreGetMetaData
func invokeFabricFileStoreGetMetaData(contractName string, methodName string, fileId, pk string) *C.char { 
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	// fmt.Printf("--> Query Result: %+v\n", args)

	// 4. 调用链码
	realName := methodName
	fmt.Println("--> Submit Transaction: FileStoreContract:getFileMeta")
	result, err := contract.SubmitTransaction(realName, fileId, pk)
	if err != nil {
		fmt.Printf("Failed to submit transaction: %v", err)
		return C.CString("")
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)

	return C.CString(string(result))
}

func callGetPoseidonHash() {
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, "FileStoreContract")

	// 3. 执行查询
	fmt.Println("--> Evaluate Transaction: FileStoreContract:GetFile")
	result, err := contract.EvaluateTransaction(
		"FileStoreContract:GetFileMetaPoseidonHashes", // 根据实际链码调整函数名
		"file1234",                  // 查询参数（文件ID）
	)
	if err != nil {
		log.Fatalf("Failed to evaluate transaction: %v", err)
	}

	fmt.Printf("--> Query Result: %s\n", string(result))
}

func callStoreFile() {
	// 1. 设置 gRPC 连接
	clientConnection := newGrpcConnection()
	defer clientConnection.Close()

	// 2. 创建 Gateway 客户端
	id := newIdentity()
	sign := newSign()
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, "FileStoreContract")

	// 3. 准备调用参数
	args := []string{
		"file12345",
		"sample.pdf",
		"123456",
		"4",
		"3",
		`["ipfs1","ipfs2","ipfs3","ipfs4"]`,
		`["hash11","hash21","hash31","hash41"]`,
	}

	// 4. 调用链码
	fmt.Println("--> Submit Transaction: FileStoreContract:StoreFile")
	result, err := contract.SubmitTransaction("StoreFile", args...)
	if err != nil {
		log.Fatalf("Failed to submit transaction: %v", err)
	}

	fmt.Printf("--> Transaction submitted successfully\nTxID: %s\nResult: %s\n", 
		string(result), result)
}

// 创建 gRPC 连接
func newGrpcConnection() *grpc.ClientConn {
	fmt.Printf("[DEBUG] creating gRPC connection to: %q\n", peerEndpoint)
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)

	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)
	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// 加载 TLS 证书
func loadCertificate(filename string) (*x509.Certificate, error) {
	certPEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certPEM)
}

// 创建身份
func newIdentity() *identity.X509Identity {
	// 1. 读取证书文件
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	// 2. 将PEM格式证书转换为x509.Certificate
	cert, err := identity.CertificateFromPEM(certPEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse certificate: %w", err))
	}

	// 3. 创建身份
	id, err := identity.NewX509Identity(mspID, cert)
	if err != nil {
		panic(err)
	}
	return id
}

// 创建签名器
func newSign() identity.Sign {
	files, err := os.ReadDir(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(keyPath, files[0].Name()))
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}
	return sign
}