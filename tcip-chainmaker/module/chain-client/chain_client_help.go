/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"fmt"
	"math"

	"chainmaker.org/chainmaker/common/v2/crypto/hash"

	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

// newClient 新建链客户端
//  @param chainmakerConfig
//  @param log
//  @return *sdk.ChainClient
//  @return error
func newClient(chainmakerConfig *common.ChainmakerConfig, log *zap.SugaredLogger) (*sdk.ChainClient, error) {
	pkcs11Config := &sdk.Pkcs11Config{
		Enabled: chainmakerConfig.Pkcs11.Enable,
	}

	nodeCfgList := make([]sdk.ChainClientOption, 0, 20)
	nodeCfgList = append(nodeCfgList, sdk.WithChainClientOrgId(chainmakerConfig.OrgId))
	nodeCfgList = append(nodeCfgList, sdk.WithChainClientChainId(chainmakerConfig.ChainId))

	nodeCfgList = append(nodeCfgList, sdk.WithUserKeyBytes([]byte(chainmakerConfig.TlsKey)))
	nodeCfgList = append(nodeCfgList, sdk.WithUserCrtBytes([]byte(chainmakerConfig.TlsCert)))
	nodeCfgList = append(nodeCfgList, sdk.WithUserSignKeyBytes([]byte(chainmakerConfig.SignKey)))
	nodeCfgList = append(nodeCfgList, sdk.WithUserSignCrtBytes([]byte(chainmakerConfig.SignCert)))
	nodeCfgList = append(nodeCfgList, sdk.WithAuthType(sdk.AuthTypeToStringMap[sdk.PermissionedWithCert]))
	nodeCfgList = append(nodeCfgList, sdk.WithRetryInterval(500))
	nodeCfgList = append(nodeCfgList, sdk.WithRetryLimit(1000))

	nodeCfgList = append(nodeCfgList, sdk.WithPkcs11Config(pkcs11Config))

	for idx, node := range chainmakerConfig.Node {
		nodeClient := sdk.NewNodeConfig(
			sdk.WithNodeAddr(node.NodeAddr),
			sdk.WithNodeConnCnt(int(node.ConnCnt)),
			sdk.WithNodeUseTLS(node.EnableTls),
			sdk.WithNodeCACerts(node.TrustRoot),
			sdk.WithNodeTLSHostName(node.TlsHostName),
		)
		log.Infof("chain%v, node idx:%v, node addr:%v, tls hostname:%v",
			chainmakerConfig.ChainId, idx, node.NodeAddr, node.TlsHostName)
		nodeCfgList = append(nodeCfgList, sdk.AddChainClientNodeConfig(nodeClient))
	}

	rpcCliCfg := sdk.NewRPCClientConfig(
		sdk.WithRPCClientMaxSendMessageSize(100),
		sdk.WithRPCClientMaxReceiveMessageSize(100),
	)

	nodeCfgList = append(nodeCfgList, sdk.WithRPCClientConfig(rpcCliCfg))

	client, err := sdk.NewChainClient(
		nodeCfgList...,
	)

	if nil != err {
		return nil, err
	}

	return client, nil
}

// getListenKey 获取监听缓存的key
//  @param chainId
//  @param contracrName
//  @return string
func getListenKey(chainId, contracrName string) string {
	return fmt.Sprintf("%s#%s", chainId, contracrName)
}

// getMerkleProve 获取默克尔树验证路径
//  @param hashType
//  @param hashs
//  @param index
//  @return [][]byte
//  @return error
func getMerkleProve(hashType string, hashs [][]byte, index uint32) ([][]byte, error) {
	merkleList, err := hash.BuildMerkleTree(hashType, hashs)
	if err != nil {
		return nil, fmt.Errorf("[getMerkleProve] BuildMerkleTree error: %s", err.Error())
	}
	path := [][]byte{}
	dep := int(math.Log2(float64(len(merkleList) + 1)))

	merkleList = merkleList[:len(merkleList)-1]
	for i := 1; i < dep; i++ {
		levelCount := int(math.Pow(2, float64(i)))           //该层有多少个元素
		levelList := merkleList[len(merkleList)-levelCount:] //该层的元素列表
		merkleList = merkleList[:len(merkleList)-levelCount]
		mask := index >> (dep - i - 1)
		if mask%2 == 0 { //left
			mask++ //取右边那个
		} else {
			mask-- //取左边那个
		}
		if levelList[mask] == nil { //单数的情况，重复取左边元素
			mask--
		}
		path = leftJoin(levelList[mask], path)
	}
	return path, nil
}

func leftJoin(n []byte, list [][]byte) [][]byte {
	result := make([][]byte, len(list)+1)
	result[0] = n
	for i, x := range list {
		result[i+1] = x
	}
	return result
}

// prove spv验证
//  @param path
//  @param merkleRoot
//  @param txHash
//  @param index
//  @param hashType
//  @return bool
//func prove(path [][]byte, merkleRoot, txHash []byte, index uint32, hashType string) bool {
//	intermediateNodes := path
//	// Shortcut the empty-block case
//	if bytes.Equal(txHash[:], merkleRoot[:]) && index == 0 && len(intermediateNodes) == 0 {
//		return true
//	}
//
//	current := txHash
//	idx := index
//	proofLength := len(intermediateNodes)
//
//	numSteps := (proofLength)
//
//	for i := 0; i < numSteps; i++ {
//		next := intermediateNodes[i]
//		if idx%2 == 1 {
//			current, _ = hashMerkleBranches(hashType, next, current)
//		} else {
//			current, _ = hashMerkleBranches(hashType, current, next)
//		}
//		idx >>= 1
//	}
//
//	return bytes.Equal(current, merkleRoot)
//}
//
//func hashMerkleBranches(hashType string, left []byte, right []byte) ([]byte, error) {
//	data := make([]byte, len(left)+len(right))
//	copy(data[:len(left)], left)
//	copy(data[len(left):], right)
//	return hash.GetByStrType(hashType, data)
//}
