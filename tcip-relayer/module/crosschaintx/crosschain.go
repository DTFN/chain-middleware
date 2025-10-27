/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package crosschaintx

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/prove"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/accesscontrol"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"
)

// CrossChainTxV1 跨链管理全局变量
var CrossChainTxV1 *CrossChainTxManager

const (
	crossChainResultTrue  = "true"
	crossChainResultFalse = "false"
)

// CrossChainTxManager 跨链管理结构体
type CrossChainTxManager struct {
	log *zap.SugaredLogger
}

// InitCrossChainManager 初始化跨链管理模块
//
//	@return error
func InitCrossChainManager() error {
	CrossChainTxV1 = &CrossChainTxManager{
		log: logger.GetLogger(logger.ModuleCrossChainTx),
	}
	CrossChainTxV1.start()
	CrossChainTxV1.Recover()
	return nil
}

// start 监听跨链并进行处理
//
//	@receiver c
func (c *CrossChainTxManager) start() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case crossChainId, ok := <-utils.CrossChainTryChan:
				if ok {
					go c.crossChainTry(crossChainId)
				}
			case crossChainId, ok := <-utils.CrossChainResultChan:
				if ok {
					go c.crossChainResult(crossChainId)
				}
			case crossChainId, ok := <-utils.CrossChainConfirmChan:
				if ok {
					go c.crosschainConfirm(crossChainId)
				}
			case crossChainId, ok := <-utils.CrossChainSrcGatewayConfirmChan:
				if ok {
					go c.crosschainSrcGatewayConfirm(crossChainId)
				}
			case jsonStr, ok := <-utils.DIDManagerUpdateDIDChan:
				if ok {
					go c.sendCrossChain(jsonStr)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("have close event happened")
				close(utils.CrossChainTryChan)
				close(utils.CrossChainResultChan)
				close(utils.CrossChainConfirmChan)
				close(utils.CrossChainSrcGatewayConfirmChan)
				close(utils.DIDManagerUpdateDIDChan)
				return
			}
		}
	}()
}

// BuildCrossChainInfoFromBeginCrossChainRequest 通过跨链请求信息生成新的跨链信息
//
//	@receiver c
//	@param beginCrossChainRequest
//	@return *common.CrossChainInfo
func (c *CrossChainTxManager) BuildCrossChainInfoFromBeginCrossChainRequest(
	beginCrossChainRequest *relay_chain.BeginCrossChainRequest) *common.CrossChainInfo {
	crossChainTxInfo := &common.CrossChainInfo{
		CrossChainName: beginCrossChainRequest.CrossChainName,
		CrossChainFlag: beginCrossChainRequest.CrossChainFlag,
		From:           beginCrossChainRequest.From,
		CrossChainMsg:  beginCrossChainRequest.CrossChainMsg,
		FirstTxContent: &common.TxContentWithVerify{
			TxContent: beginCrossChainRequest.TxContent,
		},
		State:       common.CrossChainStateValue_NEW,
		ConfirmInfo: beginCrossChainRequest.ConfirmInfo,
		CancelInfo:  beginCrossChainRequest.CancelInfo,
		Timeout:     beginCrossChainRequest.Timeout,
		CrossType:   beginCrossChainRequest.CrossType,
	}
	return crossChainTxInfo
}

// NewCrossChainInfo 生成心的跨链信息
//
//	@receiver c
//	@param crossChainInfo
//	@return string
//	@return error
func (c *CrossChainTxManager) NewCrossChainInfo(crossChainInfo *common.CrossChainInfo) (string, error) {
	kv := make(map[string][]byte)
	kv[syscontract.SaveCrossChainInfo_CROSS_CHAIN_INFO_BYTE.String()], _ = json.Marshal(crossChainInfo)
	kvJsonStr, _ := json.Marshal(kv)

	crossChainIdByte, err := relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.SaveCrossChainInfo, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[NewCrossChainInfo]: [%s] invoke contract error: %s",
			utils.SaveCrossChainInfo, err.Error())
		return "", err
	}
	crossChainId := string(crossChainIdByte)

	crossChainInfo.CrossChainId = crossChainId
	return crossChainId, nil
}

// GetCrossChainInfo 获取跨链信息
//
//	@receiver c
//	@param crossChainId
//	@return *common.CrossChainInfo
//	@return error
func (c *CrossChainTxManager) GetCrossChainInfo(crossChainId string) (*common.CrossChainInfo, error) {
	kv := make(map[string][]byte)
	kv[syscontract.GetCrossChainInfo_CROSS_CHAIN_ID.String()] = []byte(crossChainId)
	kvJsonStr, _ := json.Marshal(kv)
	crossChainInfoByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetCrossChainInfo, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[GetCrossChainInfo]: [%s] invoke contract error: %s",
			utils.GetCrossChainInfo, err.Error())
		return nil, err
	}
	if crossChainInfoByte == nil {
		msg := fmt.Sprintf("[GetCrossChainInfo]: no such cross chain id: %s",
			crossChainId)
		c.log.Error(msg)
		return nil, errors.New(msg)
	}
	//var gatewayInfo common.GatewayInfo
	var tmpCrossChainInfo common.CrossChainInfo
	err = json.Unmarshal(crossChainInfoByte, &tmpCrossChainInfo)
	if err != nil {
		c.log.Errorf("[GetCrossChainInfo]: [%s] unmarshal cross chain info error: %s",
			utils.GetCrossChainInfo, err.Error())
		return nil, err
	}
	return &tmpCrossChainInfo, nil
}

// GetCrossChainInfoByRange 获取跨链信息列表
//
//	@receiver c
//	@param startCrossChainId
//	@param stopCrossChainId
//	@return []*common.CrossChainInfo
//	@return error
func (c *CrossChainTxManager) GetCrossChainInfoByRange(
	startCrossChainId, stopCrossChainId string) ([]*common.CrossChainInfo, error) {
	kv := make(map[string][]byte)
	kv[syscontract.GetCrossChainInfoByRange_START_CROSS_CHAIN_ID.String()] = []byte(startCrossChainId)
	kv[syscontract.GetCrossChainInfoByRange_STOP_CROSS_CHAIN_ID.String()] = []byte(stopCrossChainId)
	kvJsonStr, _ := json.Marshal(kv)
	crossChainInfoListByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetCrossChainInfoByRange, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[GetCrossChainInfoByRange]: [%s] invoke contract error: %s",
			utils.GetCrossChainInfoByRange, err.Error())
		return nil, err
	}
	var crossChainInfoByteList [][]byte
	err = json.Unmarshal(crossChainInfoListByte, &crossChainInfoByteList)
	if err != nil {
		c.log.Errorf("[GetCrossChainInfoByRange]: [%s] unmarshal cross chain info list byte error: %s",
			utils.GetCrossChainInfoByRange, err.Error())
		return nil, err
	}
	//crossChainInfoByteList := strings.Split(string(crossChainInfoListByte), "##")
	var crossChainInfoList []*common.CrossChainInfo
	for _, v := range crossChainInfoByteList {
		var crossChainInfo common.CrossChainInfo
		err = json.Unmarshal(v, &crossChainInfo)
		if err != nil {
			c.log.Errorf("[GetGatewayInfoByRange]: [%s] unmarshal cross chain info error: %s",
				utils.GetCrossChainInfoByRange, err.Error())
			return nil, err
		}
		crossChainInfoList = append(crossChainInfoList, &crossChainInfo)
	}
	return crossChainInfoList, nil
}

// GetCrossChainNum 获取跨链信息个数
//
//	@receiver c
//	@return uint64
//	@return error
func (c *CrossChainTxManager) GetCrossChainNum() (uint64, error) {
	kv := make(map[string][]byte)
	kvJsonStr, _ := json.Marshal(kv)
	crossChainNumByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetCrossChainNum, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[GetCrossChainNum]: [%s] invoke contract error: %s", utils.GetCrossChainNum, err.Error())
		return 0, err
	}
	crossChainNum, err := strconv.Atoi(string(crossChainNumByte))
	if err != nil {
		c.log.Errorf("[GetCrossChainNum]: [%s] invoke contract error: %s", utils.GetCrossChainNum, err.Error())
		return 0, err
	}
	return uint64(crossChainNum), nil
}

// GetNotEndCrossChainIdList 获取未结束的跨链id
//
//	@receiver c
//	@return []*common.CrossChainInfo
//	@return error
func (c *CrossChainTxManager) GetNotEndCrossChainIdList() ([]*common.CrossChainInfo, error) {
	kv := make(map[string][]byte)
	kvJsonStr, _ := json.Marshal(kv)
	notEndCrossChainIdListByte, err := relay_chain_chainmaker.RelayChainV1.QueryContract(
		utils.CrossChainManager, utils.GetNotEndCrossChainIdList, true, string(kvJsonStr), -1)
	if err != nil {
		msg := fmt.Sprintf("[GetNotEndCrossChainIdList]Call contract error：%s", err.Error())
		c.log.Error(msg)
		return nil, errors.New(msg)
	}
	var notEndCrossChainIdList []string
	err = json.Unmarshal(notEndCrossChainIdListByte, &notEndCrossChainIdList)
	if err != nil {
		msg := fmt.Sprintf("[GetNotEndCrossChainIdList]unmarshal not end cross chain id list error：%s",
			err.Error())
		c.log.Error(msg)
		return nil, errors.New(msg)
	}

	if len(notEndCrossChainIdList) == 0 {
		c.log.Infof("[GetNotEndCrossChainIdList] No not end cross chain")
		return nil, nil
	}

	crossChainList := make([]*common.CrossChainInfo, 0)
	for _, v := range notEndCrossChainIdList {
		crossChainInfo, err1 := c.GetCrossChainInfo(v)
		if err1 != nil {
			c.log.Warnf("[GetNotEndCrossChainIdList]cross chain id not found %s, error %s", v, err.Error())
			continue
		}
		if crossChainInfo.CrossChainId == "" {
			crossChainInfo.CrossChainId = v
		}
		crossChainList = append(crossChainList, crossChainInfo)
	}
	return crossChainList, nil
}

// Recover 断电恢复
//
//	@receiver c
func (c *CrossChainTxManager) Recover() {
	c.log.Infof("[Recover] recover start")
	notEndCrossChainList, err := c.GetNotEndCrossChainIdList()
	if err != nil {
		c.log.Warnf("[Recover]Can't get not end cross chain info list: %s", err.Error())
		return
	}
	if notEndCrossChainList == nil {
		c.log.Infof("[Recover]No not end cross chain info found, needn't recover")
		return
	}
	for _, crossChainInfo := range notEndCrossChainList {
		// 等待执行表示还没有执行，直接执行就可以了
		if crossChainInfo.State == common.CrossChainStateValue_WAIT_EXECUTE {
			c.log.Infof("[Recover] CrossChainStateValue_WAIT_EXECUTE: crossChainId: %s",
				crossChainInfo.CrossChainId)
			c.crossChainTry(crossChainInfo.CrossChainId)
			continue
		}
		// 等待提交，那么先获取结果，然后提交,更新完结果之后，会直接发送confirm，不需要我们再触发一次
		if crossChainInfo.State == common.CrossChainStateValue_WAIT_CONFIRM {
			c.log.Infof("[Recover] CrossChainStateValue_WAIT_CONFIRM: crossChainId: %s",
				crossChainInfo.CrossChainId)
			c.crossChainResult(crossChainInfo.CrossChainId)
			// 这个时候无论如何，结果已经有了，那么发送一次源网关的确认，但是上一步如果目标网关的confirm没有被发送，那么会自动发送一次源
			// 网关的确认，可能触发两次，但是发送之前会确认一次结果有没有被写入，所以触发两次也没关系
			c.crosschainConfirm(crossChainInfo.CrossChainId)
			c.crosschainSrcGatewayConfirm(crossChainInfo.CrossChainId)
		}
	}
}

// crossChainTry 跨链第二阶段，远程链执行
//
//	@receiver c
//	@param crossChainId
func (c *CrossChainTxManager) crossChainTry(crossChainId string) {
	crossChainInfo, err := c.GetCrossChainInfo(crossChainId)
	if err != nil {
		c.log.Errorf("[crossChainTry] can't get crossChainInfo: crossChainid: %s, error: %s",
			crossChainId, err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(crossChainInfo.CrossChainMsg))

	c.log.Debugf("[crossChainTry]start crossChainId: %s", crossChainInfo.CrossChainId)

	srcGatewayInfo, err := gateway.GatewayV1.GetGatewayInfo(crossChainInfo.From)
	if err != nil {
		msg := fmt.Sprintf("[crossChainTry] gateway not found: id: %s error: %s, crossChainId: %s",
			crossChainInfo.From, err.Error(), crossChainId)
		c.log.Error(msg)
	}

	var crossChainTxUpChain []*common.CrossChainTxUpChain
	for i, crossChainMsg := range crossChainInfo.CrossChainMsg {
		go func(index int, crossChainMsg *common.CrossChainMsg) {
			defer wg.Done()
			// 恢复的时候如果结果不为空，那表明已经发送过请求，不需要再发一遍了
			if len(crossChainInfo.CrossChainTxContent) > 0 &&
				crossChainInfo.CrossChainTxContent[index] != nil {
				return
			}
			// 获取目标网关信息
			destGatewayInfo, err2 := gateway.GatewayV1.GetGatewayInfo(crossChainMsg.GatewayId)
			if err2 != nil {
				msg := fmt.Sprintf("[crossChainTry] gateway not found: id:%s error:%s",
					crossChainMsg.GatewayId, err2.Error())
				c.log.Error(msg)
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", nil, common.TxResultValue_GATEWAY_NOT_FOUND,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{common.TxResultValue_GATEWAY_NOT_FOUND.String()},
					},
				})
				return
			}

			// 检查是否需要处理
			if destGatewayInfo.RelayChainId != conf.Config.BaseConfig.GatewayID {
				return
			}

			c.log.Infof("[crossChainTry]start deal crossChainTry: gateayId %s, crossChainId: %s",
				destGatewayInfo.GatewayId, crossChainInfo.CrossChainId)

			// src gateway获取失败，需要填充失败信息
			if err != nil {
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", nil, common.TxResultValue_SRC_GATEWAY_GET_ERROR,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{common.TxResultValue_SRC_GATEWAY_GET_ERROR.String()},
					},
				})
				return
			}

			// 判断网关权限
			if !accesscontrol.AccessV1.GatewayPermissionsCheck(srcGatewayInfo, destGatewayInfo) {
				msg := fmt.Sprintf("[crossChainTry] Please check gateway permissions, src: %s, dest %s",
					srcGatewayInfo.GatewayId, destGatewayInfo.GatewayId)
				c.log.Error(msg)
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", nil, common.TxResultValue_TX_NO_PERMISSIONS,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{common.TxResultValue_TX_NO_PERMISSIONS.String()},
					},
				})
				return
			}

			// 检查目标网关的连接
			connOk, txResult := gateway.GatewayV1.CheckGatewayConnect(crossChainMsg.GatewayId)
			if !connOk {
				msg := fmt.Sprintf("[crossChainTry] gateway connection error: id:%s",
					crossChainMsg.GatewayId)
				c.log.Error(msg)
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", nil, txResult,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{txResult.String()},
					},
				})
				return
			}

			// 发送跨链请求
			req := &cross_chain.CrossChainTryRequest{
				Version:        common.Version_V1_0_0,
				CrossChainId:   crossChainInfo.CrossChainId,
				CrossChainName: crossChainInfo.CrossChainName,
				CrossChainFlag: crossChainInfo.CrossChainFlag,
				CrossChainMsg:  crossChainMsg,
				TxContent:      crossChainInfo.FirstTxContent,
				From:           crossChainInfo.From,
			}
			res, err2 := request.RequestV1.CrossChainTry(req, crossChainInfo.Timeout, destGatewayInfo)
			if err2 != nil {
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", []byte(err2.Error()), common.TxResultValue_TX_TIMEOUT,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{common.TxResultValue_TX_TIMEOUT.String()},
					},
				})
				return
			}
			if res.Code != common.Code_GATEWAY_SUCCESS {
				crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
					Index: int32(index),
					TxContentWithVerify: &common.TxContentWithVerify{
						TxContent: buildTxContent("", []byte(res.Message), common.TxResultValue_TX_FAIL,
							crossChainMsg.GatewayId, crossChainMsg.ChainRid, ""),
						TxVerifyResult: common.TxVerifyRsult_VERIFY_INVALID,
						TryResult:      []string{common.TxResultValue_TX_FAIL.String()},
					},
				})
				return
			}
			crossChainTxContent := &common.TxContentWithVerify{
				TxContent: res.TxContent,
				TryResult: res.TryResult,
			}
			// query类型无交易可验
			if crossChainInfo.CrossType == common.CrossType_QUERY {
				crossChainTxContent.TxVerifyResult = common.TxVerifyRsult_VERIFY_NOT_NEED
			} else {
				_, err2 = prove.ProveV1.ProveTx(crossChainTxContent, destGatewayInfo.GatewayId,
					crossChainMsg.ChainRid)
			}
			crossChainTxUpChain = append(crossChainTxUpChain, &common.CrossChainTxUpChain{
				Index:               int32(index),
				TxContentWithVerify: crossChainTxContent,
			})
			if err2 != nil {
				return
			}
		}(i, crossChainMsg)
	}

	wg.Wait()

	if len(crossChainTxUpChain) > 0 {
		kv := make(map[string][]byte)
		kv[syscontract.UpdateCrossChainTry_CROSS_CHAIN_ID.String()] = []byte(crossChainId)
		kv[syscontract.UpdateCrossChainTry_CROSS_CHAIN_TX_BYTE.String()], _ = json.Marshal(crossChainTxUpChain)
		kvJsonStr, _ := json.Marshal(kv)
		_, err := relay_chain_chainmaker.RelayChainV1.InvokeContract(
			utils.CrossChainManager, utils.UpdateCrossChainTry, true, string(kvJsonStr), -1)
		if err != nil {
			c.log.Errorf("[crossChainTry]: [%s] invoke contract error: %s",
				utils.UpdateCrossChainTry, err.Error())
			return
		}
	}

	c.log.Infof("[crossChainTry]end crossChainId: %s", crossChainInfo.CrossChainId)
}

// sendCrossChain
//
//	@receiver c
//	@param crossChainId
func (c *CrossChainTxManager) sendCrossChain(jsonStr string) {
	msg := fmt.Sprintf("[sendCrossChain] jsonStr: %s", jsonStr)
	c.log.Debugf(msg)

	// 获取所有网关，批量推送
	gatewayIds := gateway.GatewayV1.GetAllGatewayId()

	for _, gatewayId := range gatewayIds {
		destGatewayInfo, err := gateway.GatewayV1.GetGatewayInfo(gatewayId)
		if err != nil {
			msg := fmt.Sprintf("[sendCrossChain] gateway not found: id: %s error: %s, jsonStr: %s",
				gatewayId, err.Error(), jsonStr)
			c.log.Error(msg)
			continue
		}

		// 检查是否需要处理
		if destGatewayInfo.RelayChainId != conf.Config.BaseConfig.GatewayID {
			continue
		}

		c.log.Infof("[sendCrossChain]start deal crossChainTry: gateayId %s, crossChainId: %s",
			destGatewayInfo.GatewayId, gatewayId)

		// 检查目标网关的连接
		connOk, vv := gateway.GatewayV1.CheckGatewayConnect(gatewayId)
		if !connOk {
			msg := fmt.Sprintf("[sendCrossChain] gateway connection error: id:%s, error:%+v", gatewayId, vv)
			c.log.Error(msg)
			continue
		}

		// 发送跨链请求
		crossChainMsg := &common.CrossChainMsg{
				GatewayId:    "",
				ChainRid:     "$allChain",
				ContractName: "DID_MANAGER",
				Method:       "*",
				Parameter:    jsonStr,
				ConfirmInfo: nil,
				CancelInfo: nil,
			}
		req := &cross_chain.CrossChainTryRequest{
			Version:        common.Version_V1_0_0,
			CrossChainId:   "$relayer",            // 来源跨链消息ID
			CrossChainName: "DID_MANAGER", // 跨链名称
			CrossChainFlag: "",            // 跨链标记
			CrossChainMsg:  crossChainMsg, // 跨链信息
			TxContent:      nil,
			From:           "",
		}
		res, err2 := request.RequestV1.CrossChainTry(req, 60, destGatewayInfo)
		if err2 != nil {
			msg := fmt.Sprintf("[sendCrossChain] distribute error: gatewayId:%s, res: %+v", gatewayId, res)
			c.log.Error(msg)
			continue
		}
		if res.Code != common.Code_GATEWAY_SUCCESS {
			msg := fmt.Sprintf("[sendCrossChain] distribute error: gatewayId:%s, res: %+v", gatewayId, res)
			c.log.Error(msg)
			continue
		}
	}
}

// crossChainResult 根据远程链执行结果获取整个跨链交易的结果
//
//	@receiver c
//	@param crossChainId
func (c *CrossChainTxManager) crossChainResult(crossChainId string) {
	crossChainInfo, err := c.GetCrossChainInfo(crossChainId)
	if err != nil {
		c.log.Errorf("[crossChainResult] can't get crossChainInfo: crossChainid: %s, error: %s",
			crossChainId, err.Error())
		return
	}
	crossChainInfoStr, _ := json.Marshal(crossChainInfo)
	c.log.Debugf(string(crossChainInfoStr))

	srcGatewayInfo, err := gateway.GatewayV1.GetGatewayInfo(crossChainInfo.From)
	if err != nil {
		msg := fmt.Sprintf("[crossChainResult] gateway not found: id: %s error: %s, crossChainId: %s",
			crossChainInfo.From, err.Error(), crossChainId)
		c.log.Error(msg)
		return
	}

	// 检查需不需要处理
	if srcGatewayInfo.RelayChainId != conf.Config.BaseConfig.GatewayID {
		return
	}

	hasFail := false
	hasSuccess := false

	for _, crossChainTxContent := range crossChainInfo.CrossChainTxContent {
		if crossChainTxContent.TxVerifyResult == common.TxVerifyRsult_VERIFY_INVALID {
			hasFail = true
			continue
		}
		hasSuccess = true
	}

	kv := make(map[string][]byte)
	kv[syscontract.UpdateCrossChainResult_CROSS_CHAIN_ID.String()] = []byte(crossChainId)
	kv[syscontract.UpdateCrossChainResult_CROSS_CHAIN_RESULT.String()] =
		[]byte(c.getCrossChainResult(hasSuccess, hasFail, crossChainInfo, srcGatewayInfo))
	kvJsonStr, _ := json.Marshal(kv)
	_, err = relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.UpdateCrossChainResult, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[crossChainResult]: [%s] invoke contract error: %s",
			utils.UpdateCrossChainResult, err.Error())
		return
	}

	c.log.Infof("[crossChainResult]end crossChainId: %s, result: %t",
		crossChainInfo.CrossChainId, crossChainInfo.CrossChainResult)
}

// crosschainConfirm 跨链第三阶段，远程链提交
//
//	@receiver c
//	@param crossChainId
func (c *CrossChainTxManager) crosschainConfirm(crossChainId string) {
	crossChainInfo, err := c.GetCrossChainInfo(crossChainId)
	if err != nil {
		c.log.Errorf("[crosschainConfirm] can't get crossChainInfo: crossChainid: %s, error: %s",
			crossChainId, err.Error())
		return
	}
	// 跨链查询不需要让目标网关回滚或者提交
	if crossChainInfo.CrossType == common.CrossType_QUERY {
		return
	}

	c.log.Infof("[crosschainConfirm]start crossChainId: %s", crossChainInfo.CrossChainId)
	wg := sync.WaitGroup{}
	wg.Add(len(crossChainInfo.CrossChainMsg))
	var crossChainConfirm []*common.CrossChainConfirmUpChain
	for i, crossChainMsg := range crossChainInfo.CrossChainMsg {
		go func(index int, crossChainMsg *common.CrossChainMsg) {
			defer wg.Done()
			// 如果有值，表明已经发送过了，不需要再发送一次
			if len(crossChainInfo.GatewayConfirmResult) > 0 &&
				crossChainInfo.GatewayConfirmResult[index] != nil &&
				*crossChainInfo.GatewayConfirmResult[index] != (common.CrossChainConfirm{}) {
				return
			}
			// 获取目标网关信息
			destGatewayInfo, err2 := gateway.GatewayV1.GetGatewayInfo(crossChainMsg.GatewayId)
			if err2 != nil {
				// 都走到这里了，不应该发生这个错误，为了保险，还是判断一下
				msg := fmt.Sprintf("[crosschainConfirm] before confirm: gateway not found: id:%s error:%s",
					crossChainMsg.GatewayId, err2.Error())
				c.log.Error(msg)
				crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
					Index: int32(index),
					CrossChainConfirm: &common.CrossChainConfirm{
						Code:    common.Code_GATEWAY_TIMEOUT,
						Message: msg,
					},
				})
				return
			}
			// 检查要不要当前网关处理
			if destGatewayInfo.RelayChainId != conf.Config.BaseConfig.GatewayID {
				return
			}
			c.log.Infof("[crosschainConfirm]start deal confirm: gateayId %s, crossChainId: %s",
				destGatewayInfo.GatewayId, crossChainInfo.CrossChainId)

			// 如果之前没有发送交易，那么现在也不发送交易
			if crossChainInfo.CrossChainTxContent[index].TxContent.TxResult ==
				common.TxResultValue_GATEWAY_PINGPONG_ERROR ||
				crossChainInfo.CrossChainTxContent[index].TxContent.TxResult == common.TxResultValue_CHAIN_PING_ERROR {
				msg := fmt.Sprintf("[crosschainConfirm] confirm : gateway try tx not send: id:%s",
					crossChainMsg.GatewayId)

				crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
					Index: int32(index),
					CrossChainConfirm: &common.CrossChainConfirm{
						Code:    common.Code_GATEWAY_TIMEOUT,
						Message: msg,
					},
				})
				c.log.Infof("[crosschainConfirm]confirm stop crossChainId: %s, failed",
					crossChainInfo.CrossChainId)
				return
			}
			// 如果之前是没有权限，现在也提示没有权限
			if crossChainInfo.CrossChainTxContent[index].TxContent.TxResult == common.TxResultValue_TX_NO_PERMISSIONS {
				msg := fmt.Sprintf("[crosschainConfirm] confirm : gateway try tx no permissions: id:%s",
					crossChainMsg.GatewayId)

				crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
					Index: int32(index),
					CrossChainConfirm: &common.CrossChainConfirm{
						Code:    common.Code_INTERNAL_ERROR,
						Message: msg,
					},
				})
				c.log.Infof("[crosschainConfirm]confirm stop crossChainId: %s, failed",
					crossChainInfo.CrossChainId)
				return
			}
			// 发送确认或者回滚
			if crossChainInfo.CrossChainResult {
				crossChainInfo.State = common.CrossChainStateValue_CONFIRM_END
				req := &cross_chain.CrossChainConfirmRequest{
					Version:        common.Version_V1_0_0,
					CrossChainId:   crossChainInfo.CrossChainId,
					CrossChainName: crossChainInfo.CrossChainName,
					CrossChainFlag: crossChainInfo.CrossChainFlag,
					ConfirmInfo:    crossChainMsg.ConfirmInfo,
					CrossType:      crossChainInfo.CrossType,
				}
				res, err3 := request.RequestV1.CrossChainConfirm(req, crossChainInfo.Timeout, destGatewayInfo)
				if err3 != nil {
					msg := fmt.Sprintf("[crosschainConfirm] confirm : gateway error: id:%s error:%s",
						crossChainMsg.GatewayId, err3.Error())
					c.log.Error(msg)
					crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
						Index: int32(index),
						CrossChainConfirm: &common.CrossChainConfirm{
							Code:    common.Code_GATEWAY_TIMEOUT,
							Message: msg,
						},
					})
					c.log.Infof("[crosschainConfirm]confirm stop crossChainId: %s, failed", crossChainInfo.CrossChainId)
					return
				}
				crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
					Index: int32(index),
					CrossChainConfirm: &common.CrossChainConfirm{
						Code:    res.Code,
						Message: res.Message,
					},
				})
				c.log.Infof("[crosschainConfirm]confirm stop crossChainId: %s, success", crossChainInfo.CrossChainId)
			} else {
				crossChainInfo.State = common.CrossChainStateValue_CANCEL_END
				if crossChainInfo.CrossChainTxContent[index].TxContent.TxResult ==
					common.TxResultValue_GATEWAY_PINGPONG_ERROR ||
					crossChainInfo.CrossChainTxContent[index].TxContent.TxResult == common.TxResultValue_CHAIN_PING_ERROR {
					msg := fmt.Sprintf("[crosschainConfirm] confirm : gateway try tx not send: id:%s",
						crossChainMsg.GatewayId)
					crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
						Index: int32(index),
						CrossChainConfirm: &common.CrossChainConfirm{
							Code:    common.Code_GATEWAY_TIMEOUT,
							Message: msg,
						},
					})
					c.log.Debugf("[crosschainConfirm]confirm stop crossChainId: %s, failed", crossChainInfo.CrossChainId)
					return
				}
				req := &cross_chain.CrossChainCancelRequest{
					Version:        common.Version_V1_0_0,
					CrossChainId:   crossChainInfo.CrossChainId,
					CrossChainName: crossChainInfo.CrossChainName,
					CrossChainFlag: crossChainInfo.CrossChainFlag,
					CancelInfo:     crossChainMsg.CancelInfo,
				}
				res, err3 := request.RequestV1.CrossChainCancel(req, crossChainInfo.Timeout, destGatewayInfo)
				if err3 != nil {
					msg := fmt.Sprintf("[crosschainConfirm] cancel : gateway error: id:%s error:%s",
						crossChainMsg.GatewayId, err3.Error())
					c.log.Error(msg)
					crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
						Index: int32(index),
						CrossChainConfirm: &common.CrossChainConfirm{
							Code:    common.Code_GATEWAY_TIMEOUT,
							Message: msg,
						},
					})
					c.log.Infof("[crosschainConfirm]cancel stop crossChainId: %s, failed", crossChainInfo.CrossChainId)
					return
				}
				crossChainConfirm = append(crossChainConfirm, &common.CrossChainConfirmUpChain{
					Index: int32(index),
					CrossChainConfirm: &common.CrossChainConfirm{
						Code:    res.Code,
						Message: res.Message,
					},
				})
				c.log.Infof("[crosschainConfirm]cancel stop crossChainId: %s, success", crossChainInfo.CrossChainId)
			}
		}(i, crossChainMsg)
	}
	wg.Wait()

	if len(crossChainConfirm) > 0 {
		kv := make(map[string][]byte)
		kv[syscontract.UpdateCrossChainResult_CROSS_CHAIN_ID.String()] = []byte(crossChainId)
		kv[syscontract.UpdateCrossChainResult_CROSS_CHAIN_RESULT.String()], _ = json.Marshal(crossChainConfirm)
		kvJsonStr, _ := json.Marshal(kv)
		_, err := relay_chain_chainmaker.RelayChainV1.InvokeContract(
			utils.CrossChainManager, utils.UpdateCrossChainConfirm, true, string(kvJsonStr), -1)
		if err != nil {
			c.log.Errorf("[crosschainConfirm]: [%s] invoke contract error: %s",
				utils.UpdateCrossChainConfirm, err.Error())
			return
		}
	}
}

// crosschainSrcGatewayConfirm 跨链第三阶段，源链提交
//
//	@receiver c
//	@param crossChainId
func (c *CrossChainTxManager) crosschainSrcGatewayConfirm(crossChainId string) {
	crossChainInfo, err := c.GetCrossChainInfo(crossChainId)
	if err != nil {
		c.log.Errorf("[crosschainSrcGatewayConfirm] can't get crossChainInfo: crossChainid: %s, error: %s",
			crossChainId, err.Error())
		return
	}

	// 如果有值，不需要再发送一次了
	if crossChainInfo.ConfirmResult != nil && *crossChainInfo.ConfirmResult != (common.CrossChainConfirm{}) {
		return
	}

	destGatewayInfo, err2 := gateway.GatewayV1.GetGatewayInfo(crossChainInfo.From)
	if err2 != nil {
		// 都走到这里了，不应该发生这个错误，为了保险，还是判断一下
		msg := fmt.Sprintf("[crosschainSrcGatewayConfirm] before confirm: gateway not found: id:%s error:%s",
			crossChainInfo.From, err2.Error())
		c.log.Error(msg)
		crossChainInfo.ConfirmResult = &common.CrossChainConfirm{
			Code:    common.Code_GATEWAY_TIMEOUT,
			Message: msg,
		}
		return
	}
	// 检查是否需要处理
	if destGatewayInfo.RelayChainId != conf.Config.BaseConfig.GatewayID {
		return
	}
	c.log.Infof("[crosschainSrcGatewayConfirm]start deal confirm: gateayId %s, crossChainId: %s",
		destGatewayInfo.GatewayId, crossChainInfo.CrossChainId)
	// 发送确认或者回滚
	if crossChainInfo.CrossChainResult {
		crossChainInfo.State = common.CrossChainStateValue_CONFIRM_END
		req := &cross_chain.CrossChainConfirmRequest{
			Version:        common.Version_V1_0_0,
			CrossChainId:   crossChainInfo.CrossChainId,
			CrossChainName: crossChainInfo.CrossChainName,
			CrossChainFlag: crossChainInfo.CrossChainFlag,
			ConfirmInfo:    crossChainInfo.ConfirmInfo,
			TryResult:      getTryResult(crossChainInfo),
			CrossType:      crossChainInfo.CrossType,
		}
		res, err3 := request.RequestV1.CrossChainConfirm(req, crossChainInfo.Timeout, destGatewayInfo)
		if err3 != nil {
			msg := fmt.Sprintf("[crosschainSrcGatewayConfirm] confirm : gateway error: id:%s error:%s",
				crossChainInfo.From, err3.Error())
			c.log.Error(msg)
			crossChainInfo.ConfirmResult = &common.CrossChainConfirm{
				Code:    common.Code_GATEWAY_TIMEOUT,
				Message: msg,
			}
			c.log.Infof("[crosschainSrcGatewayConfirm]confirm stop crossChainId: %s, failed",
				crossChainInfo.CrossChainId)
		} else {
			crossChainInfo.ConfirmResult = &common.CrossChainConfirm{
				Code:    res.Code,
				Message: res.Message,
			}
			c.log.Infof("[crosschainSrcGatewayConfirm]confirm stop crossChainId: %s, success",
				crossChainInfo.CrossChainId)
		}
	} else {
		crossChainInfo.State = common.CrossChainStateValue_CANCEL_END
		req := &cross_chain.CrossChainCancelRequest{
			Version:        common.Version_V1_0_0,
			CrossChainId:   crossChainInfo.CrossChainId,
			CrossChainName: crossChainInfo.CrossChainName,
			CrossChainFlag: crossChainInfo.CrossChainFlag,
			CancelInfo:     crossChainInfo.CancelInfo,
		}
		res, err3 := request.RequestV1.CrossChainCancel(req, crossChainInfo.Timeout, destGatewayInfo)
		if err3 != nil {
			msg := fmt.Sprintf("[crosschainSrcGatewayConfirm] cancel : gateway error: id:%s error:%s",
				crossChainInfo.From, err3.Error())
			c.log.Error(msg)
			crossChainInfo.ConfirmResult = &common.CrossChainConfirm{
				Code:    common.Code_GATEWAY_TIMEOUT,
				Message: msg,
			}
			c.log.Infof("[crosschainSrcGatewayConfirm]cancel stop crossChainId: %s, failed",
				crossChainInfo.CrossChainId)
		} else {
			crossChainInfo.ConfirmResult = &common.CrossChainConfirm{
				Code:    res.Code,
				Message: res.Message,
			}
			c.log.Infof("[crosschainSrcGatewayConfirm]cancel stop crossChainId: %s, success",
				crossChainInfo.CrossChainId)
		}
	}

	kv := make(map[string][]byte)
	kv[syscontract.UpdateSrcGatewayConfirm_CROSS_CHAIN_ID.String()] = []byte(crossChainId)
	kv[syscontract.UpdateSrcGatewayConfirm_CONFIRM_RESULT.String()], _ = proto.Marshal(crossChainInfo.ConfirmResult)
	kvJsonStr, _ := json.Marshal(kv)
	_, err = relay_chain_chainmaker.RelayChainV1.InvokeContract(
		utils.CrossChainManager, utils.UpdateSrcGatewayConfrim, true, string(kvJsonStr), -1)
	if err != nil {
		c.log.Errorf("[crosschainSrcGatewayConfirm]: [%s] invoke contract error: %s",
			utils.UpdateSrcGatewayConfrim, err.Error())
		return
	}
}

// getCrossChainResult 计算跨链交易结果
//
//	@receiver c
//	@param hasSuccess
//	@param hasFail
//	@param crossChainInfo
//	@param srcGatewayInfo
//	@return string
func (c *CrossChainTxManager) getCrossChainResult(
	hasSuccess,
	hasFail bool,
	crossChainInfo *common.CrossChainInfo,
	srcGatewayInfo *common.GatewayInfo) string {
	c.log.Debugf("[getCrossChainResult]start crossChainId: %s", crossChainInfo.CrossChainId)
	// 等到所有的网关返回或者超时，判断跨链结果或者询问源网关跨链结果
	if hasFail && !hasSuccess {
		// 存在失败的，但是不存在成功的，那么表明全部失败，不需要找源网关判断是否成功
		return crossChainResultFalse
	} else if hasSuccess && !hasFail {
		// 存在成功的，但是不存在失败的，那么表明全部成功，不需找源网关判断是否成功
		return crossChainResultTrue
	} else {
		// 其他情况需要找网关确认才能知道结果
		txRequest := &cross_chain.IsCrossChainSuccessRequest{
			Version:        common.Version_V1_0_0,
			CrossChainId:   crossChainInfo.CrossChainId,
			CrossChainName: crossChainInfo.CrossChainName,
			CrossChainFlag: crossChainInfo.CrossChainFlag,
			TxContent:      crossChainInfo.CrossChainTxContent,
		}
		// 如果出错或者超时，将被标记为跨链失败（暂时先这么搞，等沟通以后再说）
		res, err2 := request.RequestV1.IsCrossChainSuccess(txRequest, crossChainInfo.Timeout, srcGatewayInfo)
		if err2 != nil {
			msg := fmt.Sprintf("[startCrossChain] call IsCrossChainSuccess failed: gatewayId: %s, error: %s",
				srcGatewayInfo.GatewayId, err2.Error())
			c.log.Error(msg)
			return crossChainResultFalse
		}
		if res.CrossChainResult {
			return crossChainResultTrue
		}
		return crossChainResultFalse
	}
}

// buildTxContent 构建交易内容
//
//	@param txId
//	@param tx
//	@param txResultValue
//	@param gatewayId
//	@param chainRid
//	@param txProve
//	@return *common.TxContent
func buildTxContent(
	txId string, tx []byte, txResultValue common.TxResultValue,
	gatewayId, chainRid, txProve string) *common.TxContent {
	return &common.TxContent{
		TxId:      txId,
		Tx:        tx,
		TxResult:  txResultValue,
		GatewayId: gatewayId,
		ChainRid:  chainRid,
		TxProve:   txProve,
	}
}

// getTryResult 获取参数中的query结果
//
//	@param crossChainInfo
//	@return []string
func getTryResult(crossChainInfo *common.CrossChainInfo) []string {
	res := make([]string, 0)
	for _, v := range crossChainInfo.CrossChainTxContent {
		res = append(res, v.TryResult...)
	}
	return res
}
