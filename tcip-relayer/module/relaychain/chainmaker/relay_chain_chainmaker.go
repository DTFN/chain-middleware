/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package relay_chain_chainmaker

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	"encoding/binary"
	"encoding/hex"
	"encoding/json"

	"chainmaker.org/chainmaker/sdk-go/v2/examples"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	sdkutils "chainmaker.org/chainmaker/sdk-go/v2/utils"
	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	relayChainManager "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain"
)

// RelayChainV1 中继链交互全局变量
var RelayChainV1 relayChainManager.RelayChainManager

// RelayChainChainmaker 中继链交互结构体
type RelayChainChainmaker struct {
	chainMakerClient *sdk.ChainClient
	config           *conf.RelayChain
	log              *zap.SugaredLogger
}

// InitRelayChain 初始化中继链
//  @param relayChainConfig
//  @return error
func InitRelayChain(relayChainConfig *conf.RelayChain) error {
	log := logger.GetLogger(logger.ModuleRelayChain)
	log.Debug("[InitRelayChain] init")
	cc, err := sdk.NewChainClient(
		sdk.WithConfPath(relayChainConfig.ChainmakerSdkConfigPath),
	)
	if err != nil {
		return err
	}
	log.Debug("[InitRelayChain] create chain client success")

	// Enable certificate compression
	if cc.GetAuthType() == sdk.PermissionedWithCert {
		log.Debug("[InitRelayChain] enable cert hash")
		err2 := cc.EnableCertHash()
		if err2 != nil {
			return err2
		}
	}
	log.Debug("[InitRelayChain] enable cert success")
	relayChainChainmaker := &RelayChainChainmaker{
		chainMakerClient: cc,
		config:           relayChainConfig,
		log:              log,
	}

	RelayChainV1 = relayChainChainmaker
	go relayChainChainmaker.keepClientAlive()
	utils.CrossChainTryChan = make(chan string)
	utils.CrossChainResultChan = make(chan string)
	utils.CrossChainConfirmChan = make(chan string)
	utils.CrossChainSrcGatewayConfirmChan = make(chan string)
	utils.DIDManagerUpdateDIDChan = make(chan string)
	lastBlock, err := cc.GetLastBlock(false)
	if err != nil {
		log.Errorf("[InitRelayChain] get last block error: %s\n", err.Error())
		return err
	}
	go relayChainChainmaker.listenEvent(lastBlock.Block.Header.BlockHeight)
	go relayChainChainmaker.listenDIDManagerEvent(lastBlock.Block.Header.BlockHeight)
	return nil
}

// InitContract 初始化合约
//  @receiver r
//  @param contractName
//  @param version
//  @param byteCodeBase64
//  @param kvJsonStr
//  @param withSyncResult
//  @param timeout
//  @param runtime
//  @return error
func (r *RelayChainChainmaker) InitContract(
	contractName, version, byteCodeBase64 string,
	kvJsonStr string,
	withSyncResult bool,
	timeout int64,
	runtime tcipcommon.ChainmakerRuntimeType) error {
	kvs, err := r.getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}
	payload, err := r.chainMakerClient.CreateContractCreatePayload(
		contractName, version, byteCodeBase64, common.RuntimeType(runtime), kvs)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}

	endorsers, err := r.getEndorsersWithAuthType(
		r.chainMakerClient.GetHashType(), r.chainMakerClient.GetAuthType(), payload, r.config.Users)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}

	resp, err := r.chainMakerClient.SendContractManageRequest(payload, endorsers, timeout, withSyncResult)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		r.log.Errorf("InitContract err resp %v", resp)
		return r.errorFormat("InitContract",
			fmt.Errorf("InitContract err resp %v", resp))
	}
	r.log.Debugf("CREATE %s contract success, resp: %+v\n", contractName, resp)
	return nil
}

// UpdateContract 更新合约
//  @receiver r
//  @param contractName
//  @param version
//  @param byteCodeBase64
//  @param kvJsonStr
//  @param withSyncResult
//  @param timeout
//  @param runtime
//  @return error
func (r *RelayChainChainmaker) UpdateContract(
	contractName, version, byteCodeBase64 string,
	kvJsonStr string,
	withSyncResult bool,
	timeout int64,
	runtime tcipcommon.ChainmakerRuntimeType) error {
	kvs, err := r.getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return r.errorFormat("UpdateContract", err)
	}
	payload, err := r.chainMakerClient.CreateContractUpgradePayload(
		contractName, version, byteCodeBase64, common.RuntimeType(runtime), kvs)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}
	endorsers, err := r.getEndorsersWithAuthType(
		r.chainMakerClient.GetHashType(), r.chainMakerClient.GetAuthType(), payload, r.config.Users)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}

	resp, err := r.chainMakerClient.SendContractManageRequest(payload, endorsers, timeout, withSyncResult)
	if err != nil {
		return r.errorFormat("InitContract", err)
	}

	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		r.log.Errorf("UpdateContract err resp %v", resp)
		return r.errorFormat("UpdateContract",
			fmt.Errorf("UpdateContract err resp %v", resp))
	}
	r.log.Debugf("UPDATE %s contract success, resp: %+v\n", contractName, resp)
	return nil
}

// InvokeContract invoke合约
//  @receiver r
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return error
func (r *RelayChainChainmaker) InvokeContract(
	contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error) {
	start := time.Now().Unix()
	kvs, err := r.getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return nil, r.errorFormat("InvokeContract", err)
	}

	// 经常会报找不到交易的问题，所以把延时设置的大一点
	if timeout == -1 {
		timeout = 1000
	}

	// 经常出现节点断掉的情况，这里进行一个重试，每五秒重试一次，重试10次
	resp, err := r.chainMakerClient.InvokeContract(contractName, method, "", kvs, timeout, withSyncResult)
	if err != nil {
		return nil, r.errorFormat("InvokeContract", err)
	}

	if resp.Code != common.TxStatusCode_SUCCESS {
		r.log.Errorf("InvokeContract err resp %v", resp)
		errMsg := fmt.Sprintf(
			"[InvokeContract]invoke contract failed, [code:%d]/[msg:%s]/[contractName:%s]/[txId:%s]/[kvs:%+v]\n",
			resp.Code, resp.Message, contractName, resp.TxId, kvs)
		r.log.Error(errMsg)
		return nil, errors.New(errMsg)
	}

	if withSyncResult {
		r.log.Debugf("[InvokeContract]invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n",
			resp.Code, resp.Message, resp.TxId)
		r.log.Debugf("[InvokeContract] time used: %d", time.Now().Unix()-start)
		return resp.ContractResult.Result, nil
	}
	r.log.Debug("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n",
		resp.Code, resp.Message, resp.ContractResult)
	return nil, nil
}

// QueryContract 查询合约
//  @receiver r
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return error
func (r *RelayChainChainmaker) QueryContract(
	contractName, method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, error) {
	kvs, err := r.getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return nil, r.errorFormat("QueryContract", err)
	}
	resp, err := r.chainMakerClient.QueryContract(contractName, method, kvs, -1)
	if err != nil {
		r.log.Debugf("QueryContract %v", resp)
		return nil, r.errorFormat("QueryContract", err)
	}
	err = examples.CheckProposalRequestResp(resp, true)
	if err != nil {
		return nil, r.errorFormat("QueryContract", err)
	}
	return resp.ContractResult.Result, nil
}

// getEndorsersWithAuthType 获取背书内容
//  @receiver r
//  @param hashType
//  @param authType
//  @param payload
//  @param users
//  @return []*common.EndorsementEntry
//  @return error
func (r *RelayChainChainmaker) getEndorsersWithAuthType(
	hashType crypto.HashType,
	authType sdk.AuthType,
	payload *common.Payload,
	users []*conf.User) ([]*common.EndorsementEntry, error) {
	var endorsers []*common.EndorsementEntry

	for _, u := range users {
		var entry *common.EndorsementEntry
		var err error
		switch authType {
		case sdk.PermissionedWithCert:
			entry, err = sdkutils.MakeEndorserWithPath(u.SignKeyPath, u.SignCrtPath, payload)
			if err != nil {
				return nil, err
			}

		case sdk.PermissionedWithKey:
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, u.OrgId, payload)
			if err != nil {
				return nil, err
			}

		case sdk.Public:
			entry, err = sdkutils.MakePkEndorserWithPath(u.SignKeyPath, hashType, "", payload)
			if err != nil {
				return nil, err
			}

		default:
			return nil, errors.New("invalid authType")
		}
		endorsers = append(endorsers, entry)
	}

	return endorsers, nil
}


func (r *RelayChainChainmaker) listenDIDManagerEvent(startBlock uint64) {
	r.log.Debugf("[didManagerEvent] listenDIDManagerEvent")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventChan, err := r.chainMakerClient.SubscribeContractEvent(ctx, int64(startBlock), -1,
		utils.DIDManager, "")
	if err != nil {
		r.log.Errorf("[didManagerEvent] listen cross chain manager event failed, contract_demo name: %s",
			utils.CrossChainManager)
		return
	}
	r.log.Infof("[didManagerEvent]success listen contract manager event: contract: %s start block: %d",
		utils.DIDManager, startBlock)

	var readedBlock int64 = 0
	for {
		select {
		case event, ok := <-eventChan:
			r.log.Debugf("event: %+v", event)
			if !ok {
				r.log.Infof("[didManagerEvent]chan is close! contract_demo name: %s", utils.CrossChainManager)
				return
			}
			if event == nil {
				r.log.Infof("[didManagerEvent]require not nil, contract_demo name:: %s", utils.CrossChainManager)
			}
			contractEventInfo, _ := event.(*common.ContractEventInfo)
			r.log.Debugf("[didManagerEvent]recv contract_demo event [%d] => %+v\n",
				contractEventInfo.BlockHeight, contractEventInfo)

			// 为了防止已经产生的区块再次触发事件
			r.log.Debugf("[didManagerEvent] blockHieght: %s, readedBlock: %s", contractEventInfo.BlockHeight, readedBlock)
			if int64(contractEventInfo.BlockHeight) <= readedBlock {
				r.log.Debugf("[didManagerEvent] blockHieght: %s, readedBlock: %s", contractEventInfo.BlockHeight, readedBlock)
				continue
			}
			readedBlock = int64(contractEventInfo.BlockHeight)

			switch contractEventInfo.Topic {
			case "bc4b9177d949f0995aac6fc1e139b1e35a50170b6e9ca10360f73c3ac103d905":
				r.log.Debugf("[didManagerEvent]catch event")
				// todo 解析eventData为json字符串
				didUpdateMsg, error := r.parse(contractEventInfo.EventData[0])
				r.log.Debugf("[didManagerEvent]parse event, %+v", didUpdateMsg)
				if error != nil {
					r.log.Debugf("[didManagerEvent]parse event data fail")
				} else {
					r.log.Debugf("[didManagerEvent]parse event data success")
					utils.DIDManagerUpdateDIDChan <- didUpdateMsg
				}
			default:
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

type DIDUpdateMsg struct {
	DID       string      `json:"did"`
	DIDDocument       string      `json:"didDocument"`
}

func (r *RelayChainChainmaker) parse(rawHexStr string) (string, error) {
	// 移除可能存在的0x前缀
	hexStr := strings.TrimPrefix(rawHexStr, "0x")
	
	// 转换为字节数组
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Printf("parse error: %v\n", err)
		return "", errors.New("decode fail")
	}

	r.log.Debugf("[didManagerEvent]1 %+v", data)
	didIdentifierOffset := bytesToUint256(data[0:32])
	r.log.Debugf("[didManagerEvent]2 %s %+v", didIdentifierOffset, data[0:32])
	didDocumentOffset := bytesToUint256(data[32:64])
	r.log.Debugf("[didManagerEvent]3 %s %+v", didDocumentOffset, data[32:64])

	didIdentifierLen := bytesToUint256(data[didIdentifierOffset:didIdentifierOffset + 32])
	r.log.Debugf("[didManagerEvent]4")
	didIdentifier := string(data[didIdentifierOffset + 32 : didIdentifierOffset + 32 + didIdentifierLen])

	didDocumentLen := bytesToUint256(data[didDocumentOffset:didDocumentOffset + 32])
	r.log.Debugf("[didManagerEvent]5")
	didDocument := string(data[didDocumentOffset + 32 : didDocumentOffset + 32 + didDocumentLen])
	
	// 生成json
	didUpdateMsg := DIDUpdateMsg {
		DID: didIdentifier,
		DIDDocument: didDocument,
	}
	jsonStr, err := json.Marshal(didUpdateMsg)
	return string(jsonStr), err
}

// 将32字节转换为uint64整数,只取后8位
func bytesToUint256(b []byte) uint64 {
	bs := b[24: 32]
	var result uint64 = binary.BigEndian.Uint64(bs)
	return result
}

// listenEvent 监听事件
//  @receiver r
//  @param startBlock
func (r *RelayChainChainmaker) listenEvent(startBlock uint64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventChan, err := r.chainMakerClient.SubscribeContractEvent(ctx, int64(startBlock), -1,
		utils.CrossChainManager, "")
	if err != nil {
		r.log.Errorf("[listenEvent] listen cross chain manager event failed, contract_demo name: %s",
			utils.CrossChainManager)
		return
	}
	r.log.Infof("[listenEvent]success listen contract manager event: contract: %s start block: %d",
		utils.CrossChainManager, startBlock)

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				r.log.Infof("[listenEvent]chan is close! contract_demo name: %s", utils.CrossChainManager)
				return
			}
			if event == nil {
				r.log.Infof("[listenEvent]require not nil, contract_demo name:: %s", utils.CrossChainManager)
			}
			contractEventInfo, _ := event.(*common.ContractEventInfo)
			r.log.Debugf("[listenEvent]recv contract_demo event [%d] => %+v\n",
				contractEventInfo.BlockHeight, contractEventInfo)
			// 这里是为了防止已经产生的区块再次触发事件
			if contractEventInfo.BlockHeight == startBlock {
				continue
			}
			switch contractEventInfo.Topic {
			case tcipcommon.EventName_NEW_CROSS_CHAIN.String():
				utils.CrossChainTryChan <- contractEventInfo.EventData[0]
			case tcipcommon.EventName_CROSS_CHAIN_TRY_END.String():
				utils.CrossChainResultChan <- contractEventInfo.EventData[0]
			case tcipcommon.EventName_UPADATE_RESULT_END.String():
				utils.CrossChainConfirmChan <- contractEventInfo.EventData[0]
			case tcipcommon.EventName_GATEWAY_CONFIRM_END.String():
				utils.CrossChainSrcGatewayConfirmChan <- contractEventInfo.EventData[0]
			case tcipcommon.EventName_SRC_GATEWAY_CONFIRM_END.String():
				r.log.Infof("[listenEvent]cross chain tx finish %s\n",
					contractEventInfo.EventData[0])
			default:
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}

// getKvsFromKvJsonStr 根据参数构建合约参数
//  @receiver r
//  @param kvJsonStr
//  @return []*common.KeyValuePair
//  @return error
func (r *RelayChainChainmaker) getKvsFromKvJsonStr(kvJsonStr string) ([]*common.KeyValuePair, error) {
	var kvMap map[string][]byte
	err := json.Unmarshal([]byte(kvJsonStr), &kvMap)
	if err != nil {
		errStr := fmt.Sprintf("[getKvsFromKvJsonStr] kvJsonStr must be json string: %s -> %s",
			kvJsonStr, err.Error())
		r.log.Error(errStr)
		return nil, errors.New(errStr)
	}
	kvs := []*common.KeyValuePair{}
	for k, v := range kvMap {
		kv := &common.KeyValuePair{
			Key:   k,
			Value: v,
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
}

// errorFormat 错误打印
//  @receiver r
//  @param methodName
//  @param err
//  @return error
func (r *RelayChainChainmaker) errorFormat(methodName string, err error) error {
	msg := fmt.Sprintf("[%s]:%s", methodName, err.Error())
	r.log.Errorf(msg)
	return errors.New(msg)
}

// keepClientAlive 链的存活性检查，虽然sdk有重连机制，但是他
//					的重连不是被重复尝试，只有发送了请求才会重连，而且不会一直重连，不符合我们的要求，所以手动执行重连
//  @receiver r
func (r *RelayChainChainmaker) keepClientAlive() {
	// 在实际使用过程中有连接失败的问题，所以暂时加了一个携程来保证连接不会断掉
	for {
		time.Sleep(10 * time.Second)
		r.log.Debugf("keepClientAlive")
		if _, err := r.chainMakerClient.GetChainMakerServerVersion(); err != nil {
			r.log.Debugf("reconnect")
			// 这里是为了防止重连失败导致chainMakerClient为nil，所以要反复重连
			for {
				// 先建立连接
				client, err := sdk.NewChainClient(
					sdk.WithConfPath(r.config.ChainmakerSdkConfigPath),
				)
				if err != nil {
					r.log.Errorf("[keepClientAlive]Reconnection failure: %s", err.Error())
					time.Sleep(time.Second)
					continue
				}
				// 后释放旧连接
				_ = r.chainMakerClient.Stop()
				r.chainMakerClient = client
				lastBlock, err := r.chainMakerClient.GetLastBlock(false)
				if err != nil {
					r.log.Errorf("[keepClientAlive] get last block error: %s, listen event error\n", err.Error())
					break
				}
				go r.listenEvent(lastBlock.Block.Header.BlockHeight)
				break
			}

			// 启动DID合约监听
			r.log.Debugf("reconnect didmanager")
			lastBlock, err := r.chainMakerClient.GetLastBlock(false)
			if err != nil {
				r.log.Errorf("[keepClientAlive] get last block error: %s, listen event error\n", err.Error())
				break
			}
			go r.listenDIDManagerEvent(lastBlock.Block.Header.BlockHeight)
		} else {
			r.log.Debugf("connection err: %+v", err)
		}
	}
}
