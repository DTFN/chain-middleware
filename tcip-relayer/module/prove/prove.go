/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package prove

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"go.uber.org/zap"
)

// ProveV1 交易认证全局变量
var ProveV1 *ProveManager

// ProveManager 交易认证结构体
type ProveManager struct {
	log *zap.SugaredLogger
}

// InitProveManager 初始化交易认证
//  @return error
func InitProveManager() error {
	ProveV1 = &ProveManager{
		log: logger.GetLogger(logger.ModuleProve),
	}
	return nil
}

// ProveTx 交易认证
//  @receiver p
//  @param tx
//  @param gatewayId
//  @param chainRid
//  @return bool
//  @return error
func (p *ProveManager) ProveTx(tx *common.TxContentWithVerify, gatewayId, chainRid string) (bool, error) {
	gatewayInfo, err := gateway.GatewayV1.GetGatewayInfo(gatewayId)
	if err != nil {
		tx.TxVerifyResult = common.TxVerifyRsult_VERIFY_INVALID
		p.log.Errorf("[ProveTx] get gateway error: %s %s", gatewayId, err.Error())
		return false, err
	}
	p.log.Debugf("ProveTx, gatewayInfo: %+v", gatewayInfo)
	switch gatewayInfo.TxVerifyType {
	case common.TxVerifyType_SPV:
		var msg string
		tx.TxVerifyResult, msg = p.spvProve(tx, gatewayId, chainRid)
		if tx.TxVerifyResult != common.TxVerifyRsult_VERIFY_SUCCESS {
			tx.TryResult = []string{msg}
			return false, nil
		}
		return true, nil
	case common.TxVerifyType_RPC_INTERFACE:
		tx.TxVerifyResult = p.rpcProve(tx, gatewayInfo.TxVerifyInterface)
		if tx.TxVerifyResult != common.TxVerifyRsult_VERIFY_SUCCESS {
			return false, nil
		}
		return true, nil
	case common.TxVerifyType_NOT_NEED:
		tx.TxVerifyResult = common.TxVerifyRsult_VERIFY_NOT_NEED
		return true, nil
	default:
		return false, errors.New("TxVerifyType error, please update gateway info")
	}
}

func (p *ProveManager) spvProve(tx *common.TxContentWithVerify,
	gatewayId, chainId string) (common.TxVerifyRsult, string) {
	var (
		res []byte
		err error
	)
	spvName := utils.GetSpvContractName(gatewayId, chainId)
	getBlockHeaderParam := utils.GetBlockHeaderParam(tx.TxContent.BlockHeight)
	// 循环100秒钟，如果还没有同步到区块头，那表明系统出了问题，给他一个证明失败的结果
	for i := 0; i < 30; i++ {
		res, err = relay_chain_chainmaker.RelayChainV1.QueryContract(
			spvName, utils.GetBlockHeaderMethod,
			true, getBlockHeaderParam, -1)
		if err != nil {
			p.log.Warnf("[spvProve] get block header error: block height: %d, error: %s",
				tx.TxContent.BlockHeight, err.Error())
			if strings.Contains(err.Error(), "contractName not found") {
				return common.TxVerifyRsult_VERIFY_INVALID, "spv contract not found"
			}
		}
		if err == nil && res != nil {
			break
		}
		time.Sleep(time.Second)
	}
	res, err = relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.GetSpvContractName(gatewayId, chainId), utils.SpvTxVerifyMethod,
		true, tx.TxContent.TxProve, -1)
	if err != nil {
		msg := fmt.Sprintf("[spvProve] tx verify error: gatewayId: %s chainId: %s, err: %s",
			gatewayId, chainId, err.Error())
		p.log.Error(msg)
		return common.TxVerifyRsult_VERIFY_INVALID, ""
	}
	if string(res) == "true" {
		return common.TxVerifyRsult_VERIFY_SUCCESS, ""
	}
	p.log.Debugf("[spvProve] %s", string(res))
	return common.TxVerifyRsult_VERIFY_INVALID, ""
}

func (p *ProveManager) rpcProve(
	tx *common.TxContentWithVerify, txVerifyInterface *common.TxVerifyInterface) common.TxVerifyRsult {
	resByte, err := request.RequestV1.VerifyTx(txVerifyInterface, tx.TxContent.TxProve)
	if err != nil {
		msg := fmt.Sprintf("[rpcProve] tx verify error: gatewayId: %s txProve: %s",
			tx.TxContent.GatewayId, tx.TxContent.TxProve)
		p.log.Error(msg)
		return common.TxVerifyRsult_VERIFY_INVALID
	}
	var res cross_chain.TxVerifyResponse
	err = json.Unmarshal(resByte, &res)
	if err != nil {
		p.log.Errorf("[rpcProve] Unmarshal resByte error: %s", err.Error())
		return common.TxVerifyRsult_VERIFY_INVALID
	}
	p.log.Debugf("[rpcProve] VerifyTx result: %v", res)
	if res.TxVerifyResult {
		return common.TxVerifyRsult_VERIFY_SUCCESS
	}
	//resString := string(resByte)
	//if strings.Contains(resString, "true") {
	//	return common.TxVerifyRsult_VERIFY_SUCCESS
	//}
	return common.TxVerifyRsult_VERIFY_INVALID
}
