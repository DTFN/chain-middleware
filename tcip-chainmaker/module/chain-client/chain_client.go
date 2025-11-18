/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/db"

	chain_config "chainmaker.org/chainmaker/tcip-chainmaker/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"

	"github.com/gogo/protobuf/proto"

	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/request"

	"chainmaker.org/chainmaker/pb-go/v2/common"

	"go.uber.org/zap"

	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-chainmaker/v2/module/logger"
	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	chainmakerutils "chainmaker.org/chainmaker/utils/v2"
)

// ChainClientItfc 链客户端接口
type ChainClientItfc interface {
	InvokeContract(chainRid, contractName,
		method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, *common.TransactionInfo, error)
	GetTxProve(blockHeight uint64, chainRid string, tx *common.TransactionInfo) string
	CheckChain() bool
}

// ChainClient 链客户端结构体
type ChainClient struct {
	client               map[string]*sdk.ChainClient
	listenCtx            map[string]context.CancelFunc
	listenBlockHeaderCtx map[string]context.CancelFunc
	listen               map[string]bool
	log                  *zap.SugaredLogger
}

// ChainClientV1 连交互模块对象
var ChainClientV1 ChainClientItfc

const (
	emptyJson = "{}"
	success   = "success"
)

// InitChainClient 初始化链客户端
//  @return error
func InitChainClient() error {
	log := logger.GetLogger(logger.ModuleChainmakerClient)
	log.Debug("[InitChainClient] init")
	chainmakerClient := &ChainClient{
		client:               make(map[string]*sdk.ChainClient),
		listen:               make(map[string]bool),
		listenCtx:            make(map[string]context.CancelFunc),
		listenBlockHeaderCtx: make(map[string]context.CancelFunc),
		log:                  logger.GetLogger(logger.ModuleChainClient),
	}
	eventList, _ := event.EventManagerV1.GetEvent("")
	chainConfigList, err := chain_config.ChainConfigManager.Get("")
	if err != nil {
		panic(fmt.Sprintf("get chain config error: %v", err))
	}
	for _, chainConfig := range chainConfigList {
		cc, err := newClient(chainConfig, chainmakerClient.log)
		state := true
		msg := success
		if err != nil {
			state = false
			msg = err.Error()
			log.Errorf("[InitChainClient] newClient error: %s %v", chainConfig.ChainRid, err)
		}
		log.Infof("[InitChainClient] SetState %t, chain rid %s", state, chainConfig.ChainRid)
		err1 := chain_config.ChainConfigManager.SetState(chainConfig, state, msg)
		if err1 != nil {
			log.Errorf("[InitChainClient] SetState %s", err1.Error())
			continue
		}
		if !state {
			continue
		}
		log.Debugf("[InitChainClient] create chain rid [%s] client success", chainConfig.ChainRid)

		// Enable certificate compression
		if cc.GetAuthType() == sdk.PermissionedWithCert {
			log.Debug("[InitChainClient] enable cert hash")
			err2 := cc.EnableCertHash()
			if err2 != nil {
				_ = chain_config.ChainConfigManager.SetState(chainConfig, false, err2.Error())
				log.Errorf(err2.Error())
				continue
			}
		}
		chainmakerClient.client[chainConfig.ChainRid] = cc

		if conf.Config.BaseConfig.TxVerifyType == conf.SpvTxVerify {
			go chainmakerClient.listenBlockHeader(chainConfig.ChainRid)
		}
	}
	for _, crossChainevent := range eventList {
		listenKey := getListenKey(crossChainevent.ChainRid, crossChainevent.ContractName)
		if _, ok := chainmakerClient.listen[listenKey]; ok {
			continue
		}
		lastBlock, err := chainmakerClient.client[crossChainevent.ChainRid].GetLastBlock(false)
		if err != nil {
			msg := fmt.Sprintf("[InitChainClient] get last block error: %s\n", err.Error())
			log.Error(msg)
			_ = event.EventManagerV1.SetEventState(crossChainevent.CrossChainEventId, false, msg)
			continue
		}
		//chainmakerClient.listen[listenKey] = true
		go chainmakerClient.listenEvent(
			crossChainevent.CrossChainEventId, crossChainevent.ContractName,
			crossChainevent.ChainRid, int64(lastBlock.Block.Header.BlockHeight))
		time.Sleep(time.Millisecond * 100)
	}
	ChainClientV1 = chainmakerClient
	//go chainmakerClient.keepClientAlive()
	go chainmakerClient.eventStart()
	go chainmakerClient.chainConfigStart()
	return nil
}

// 监听事件更新
//  @receiver c
func (c *ChainClient) eventStart() {
	go func() {
		utils.EventChan = make(chan *utils.EventOperate, 10)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case eventInfo, ok := <-utils.EventChan:
				if !ok {
					c.log.Warn("error event")
				}
				c.log.Infof("[event watch] receive chain config update: %v", eventInfo)
				listenKey := getListenKey(eventInfo.ChainRid, eventInfo.ContractName)
				_, has := c.listen[listenKey]
				if eventInfo.Operate == tcipcommon.Operate_SAVE && !has {
					client, err := c.getChainClient(eventInfo.ChainRid)
					if err != nil {
						msg := fmt.Sprintf("[event watch] %s", err.Error())
						c.log.Errorf(msg)
						_ = event.EventManagerV1.SetEventState(eventInfo.CrossChainEventID, false, msg)
					} else {
						lastBlock, err := client.GetLastBlock(false)
						if err != nil {
							msg := fmt.Sprintf("[event watch] get last block error: %s\n", err.Error())
							c.log.Error(msg)
							_ = event.EventManagerV1.SetEventState(eventInfo.CrossChainEventID, false, msg)
							continue
						}
						go c.listenEvent(
							eventInfo.CrossChainEventID, eventInfo.ContractName,
							eventInfo.ChainRid, int64(lastBlock.Block.Header.BlockHeight))
					}
				}
				if eventInfo.Operate == tcipcommon.Operate_SAVE && has {
					_ = event.EventManagerV1.SetEventState(eventInfo.CrossChainEventID, true, success)
				}
				if eventInfo.Operate == tcipcommon.Operate_DELETE && has {
					c.listenCtx[listenKey]()
					delete(c.listenCtx, listenKey)
					delete(c.listen, listenKey)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("have close event happened")
				close(utils.EventChan)
				return
			}
		}
	}()
}

// 监听chainconfig更新
//  @receiver c
func (c *ChainClient) chainConfigStart() {
	go func() {
		utils.UpdateChainConfigChan = make(chan *utils.ChainConfigOperate, 10)
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case chainConfigOperate, ok := <-utils.UpdateChainConfigChan:
				if !ok {
					c.log.Warn("error event")
					continue
				}
				c.log.Infof("[chain config watch] receive chain config update: %v", chainConfigOperate)
				_, has := c.client[chainConfigOperate.ChainRid]
				chainConfig, err := chain_config.ChainConfigManager.Get(chainConfigOperate.ChainRid)
				if err != nil {
					continue
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_SAVE {
					if has {
						_ = c.client[chainConfigOperate.ChainRid].Stop()
						delete(c.client, chainConfigOperate.ChainRid)
					}
					client, err := newClient(chainConfig[0], c.log)
					if err != nil {
						c.log.Errorf("[chain config watch] %s", err.Error())
						err1 := chain_config.ChainConfigManager.SetState(chainConfig[0], false, err.Error())
						if err1 != nil {
							c.log.Errorf("[chain config watch] SetState %s", err1.Error())
							continue
						}
					} else {
						err1 := chain_config.ChainConfigManager.SetState(chainConfig[0], true, success)
						if err1 != nil {
							c.log.Errorf("[chain config watch] SetState %s", err1.Error())
							continue
						}
						c.client[chainConfig[0].ChainRid] = client
						if conf.Config.BaseConfig.TxVerifyType == conf.SpvTxVerify {
							_ = c.resetLaseBlockHeaderHeight(chainConfig[0].ChainRid)
							go c.listenBlockHeader(chainConfig[0].ChainRid)
						}
					}
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_DELETE && has {
					if _, ok := c.listenBlockHeaderCtx[chainConfigOperate.ChainRid]; ok {
						c.listenBlockHeaderCtx[chainConfigOperate.ChainRid]()
						delete(c.listenBlockHeaderCtx, chainConfigOperate.ChainRid)
					}
					delete(c.client, chainConfigOperate.ChainRid)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("have close event happened")
				close(utils.EventChan)
				return
			}
		}
	}()
}

// InvokeContract invoke合约
//  @receiver c
//  @param chainRid
//  @param contractName
//  @param method
//  @param withSyncResult
//  @param kvJsonStr
//  @param timeout
//  @return []byte
//  @return *common.TransactionInfo
//  @return error
func (c *ChainClient) InvokeContract(chainRid, contractName,
	method string, withSyncResult bool, kvJsonStr string, timeout int64) ([]byte, *common.TransactionInfo, error) {
	c.log.Debugf("InvokeContract chainRid: %s, contractName: %s, method: %s, withSyncResult: %t, kvJsonStr: %s, timeout: %d", chainRid, contractName, method, withSyncResult, kvJsonStr, timeout)
	kvs, err := c.getKvsFromKvJsonStr(kvJsonStr)
	if err != nil {
		return nil, nil, c.errorFormat("InvokeContract", err)
	}
	c.log.Debugf("InvokeContract kvs: %+v", kvs);

	client, err := c.getChainClient(chainRid)
	if err != nil {
		c.log.Errorf("[InvokeContract] %s", err.Error())
		return nil, nil, err
	}
	resp, err := client.InvokeContract(contractName, method, "", kvs, timeout, withSyncResult)
	if err != nil {
		return nil, nil, c.errorFormat("InvokeContract", err)
	}
	c.log.Infof("[InvokeContract] resp: %v\n", resp)

	if resp.Code != common.TxStatusCode_SUCCESS {
		errMsg := fmt.Sprintf(
			"[InvokeContract]invoke contract failed, %v\n", resp)
		c.log.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}

	if withSyncResult {
		c.log.Debugf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n",
			resp.Code, resp.Message, resp.TxId)
		tx, err := client.GetTxByTxId(resp.TxId)
		if err != nil {
			c.log.Errorf("get tx error, txId: %s", resp.TxId)
			return resp.ContractResult.Result, nil, err
		}
		return resp.ContractResult.Result, tx, nil
	}
	c.log.Debug("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n",
		resp.Code, resp.Message, resp.ContractResult)
	return nil, nil, nil
}

// listenEvent 监听事件
//  @receiver c
//  @param contractName
//  @param chainRid
//  @param startBlock
func (c *ChainClient) listenEvent(crossChainEvnetId, contractName, chainRid string, startBlock int64) {
	c.log.Debugf("listenEvent crossChainEvnetId: %s, contractName: %s, chainRid: %s, startBlock: %d", crossChainEvnetId, contractName, chainRid, startBlock)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := c.getChainClient(chainRid)
	if err != nil {
		msg := fmt.Sprintf("[listenEvent] %s", err.Error())
		c.log.Error(msg)
		_ = event.EventManagerV1.SetEventState(crossChainEvnetId, false, msg)
		return
	}
	// 先检查一下合约存不存在
	contractInfo, err := client.GetContractInfo(contractName)
	if err != nil {
		msg := fmt.Sprintf("[listenEvent]GetContractInfo error: %s, contract name: %s", err.Error(),
			contractName)
		c.log.Errorf(msg)
		_ = event.EventManagerV1.SetEventState(crossChainEvnetId, false, msg)
		return
	}
	if contractInfo.Name != contractName || contractInfo.Status != common.ContractStatus_NORMAL {
		msg := fmt.Sprintf("[listenEvent]Contract name: %s not existed", contractName)
		c.log.Errorf(msg)
		_ = event.EventManagerV1.SetEventState(crossChainEvnetId, false, msg)
		return
	}
	// 监听所有的topic,这样就不需要每次都改配置了
	eventChan, err := client.SubscribeContractEvent(ctx, startBlock, -1, contractName, "")
	if err != nil {
		c.log.Errorf("[listenEvent] listen cross chain event failed, contract name: %s", contractName)
		return
	}
	c.log.Infof("[listenEvent]success listen contract event: chainRid: %s contract: %s start block: %d",
		chainRid, contractName, startBlock)
	c.listenCtx[getListenKey(chainRid, contractName)] = cancel
	c.listen[getListenKey(chainRid, contractName)] = true
	_ = event.EventManagerV1.SetEventState(crossChainEvnetId, true, success)
	for {
		select {
		case contractEvent, ok := <-eventChan:
			if !ok {
				c.log.Infof("[listenEvent]chan is close! contract name: %s", contractName)
				time.Sleep(time.Second * 5)
				go c.listenEvent(crossChainEvnetId, contractName, chainRid, startBlock)
				continue
			}
			if contractEvent == nil {
				c.log.Infof("[listenEvent]require not nil, contract name:: %s", contractName)
				continue
			}
			contractEventInfo, ok := contractEvent.(*common.ContractEventInfo)
			if !ok {
				c.log.Info("[listenEvent]require true")
				continue
			}
			c.log.Debugf("[listenEvent]recv contract event [%d] => %+v\n",
				contractEventInfo.BlockHeight, contractEventInfo)
			// 这里是为了防止已经产生的区块再次触发跨链事件
			if contractEventInfo.BlockHeight == uint64(startBlock) {
				continue
			}
			tx, err := client.GetTxByTxId(contractEventInfo.TxId)
			if err != nil {
				c.log.Errorf("[listenEvent] can't get tx: txId %s topic %s event %+v",
					contractEventInfo.TxId, contractEventInfo.Topic, contractEventInfo)
				continue
			}
			txProve := ""
			if conf.Config.BaseConfig.TxVerifyType == conf.SpvTxVerify {
				txProve = c.GetTxProve(contractEventInfo.BlockHeight, chainRid, tx)
			}
			txByte, err := proto.Marshal(tx.Transaction)
			if err != nil {
				c.log.Errorf("[listenEvent] Marshal tx error: txId %s topic %s event %+v err %s",
					contractEventInfo.TxId, contractEventInfo.Topic, contractEventInfo, err.Error())
				continue
			}
			eventInfo := &event.EventInfo{
				Topic:        contractEventInfo.Topic,
				ChainRid:     chainRid,
				ContractName: contractEventInfo.ContractName,
				TxProve:      txProve,
				Data:         contractEventInfo.EventData,
				Tx:           txByte,
				TxId:         contractEventInfo.TxId,
				BlockHeight:  int64(contractEventInfo.BlockHeight),
			}
			c.log.Infof("[listenEvent] eventInfo: %v\n", eventInfo.ToString())
			go request.RequestV1.BeginCrossChain(eventInfo)
		case <-ctx.Done():
			go c.listenEvent(crossChainEvnetId, contractName, chainRid, startBlock)
			return
		}
	}
}

// listenBlockHeader 监听区块头
//  @receiver c
//  @param chainRid
//  @param startBlock
func (c *ChainClient) listenBlockHeader(chainRid string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := c.getChainClient(chainRid)
	if err != nil {
		c.log.Errorf("[listenBlockHeader] %s", err.Error())
		return
	}
	startBlock := c.getLaseBlockHeaderHeight(chainRid)
	blockHeaderChan, err := client.SubscribeBlock(ctx, startBlock, -1,
		false, true)
	if err != nil {
		c.log.Errorf("[listenBlockHeader] listen block header failed, chainRid: %s", chainRid)
		return
	}
	c.log.Infof("[listenBlockHeader]success listen block header: chainRid: %s start block: %d",
		chainRid, startBlock)
	c.listenBlockHeaderCtx[chainRid] = cancel
	for {
		select {
		case block, ok := <-blockHeaderChan:
			if !ok {
				c.log.Infof("[listenBlockHeader]chan is close! chainRid: %s", chainRid)
				go c.listenBlockHeader(chainRid)
				return
			}
			if block == nil {
				c.log.Infof("[listenBlockHeader]require not nil, chainRid: %s", chainRid)
			}
			blockHeader, ok := block.(*common.BlockHeader)
			if !ok {
				c.log.Info("[listenBlockHeader]require true")
			}
			c.log.Debugf("[listenBlockHeader]recv block header [%d] => %+v\n",
				blockHeader.BlockHeight, blockHeader)
			go request.RequestV1.SyncBlockHeader(blockHeader, chainRid)
		case <-ctx.Done():
			go c.listenBlockHeader(chainRid)
			return
		}
	}
}

// GetTxProve 获取交易证明
//  @receiver c
//  @param blockHeight
//  @param chainRid
//  @param tx
//  @return string
func (c *ChainClient) GetTxProve(blockHeight uint64, chainRid string, tx *common.TransactionInfo) string {
	client, err := c.getChainClient(chainRid)
	if err != nil {
		c.log.Errorf("[listenBlockHeader] %s", err.Error())
		return ""
	}
	block, err := client.GetBlockByHeight(blockHeight, true)
	if err != nil {
		return emptyJson
	}
	txHashes := make([][]byte, 0)
	chainConfig, err := client.GetChainConfig()
	if err != nil {
		return emptyJson
	}
	for _, tx1 := range block.Block.Txs {
		txHash, err2 := chainmakerutils.CalcTxHashWithVersion(
			chainConfig.Crypto.Hash, tx1, int(block.Block.Header.BlockVersion))
		if err2 != nil {
			c.log.Errorf("calc tx hash failed, block height:%d, txId:%x, err2:%v",
				block.Block.Header.BlockHeight, tx1.Payload.TxId, err2)
			return emptyJson
		}
		txHashes = append(txHashes, txHash)
	}
	proveHash, err := getMerkleProve(chainConfig.Crypto.Hash, txHashes, tx.TxIndex)
	if err != nil {
		c.log.Errorf("[GetTxProve] getMerkleProve error: %s", err.Error())
		return emptyJson
	}
	txProve := make(map[string][]byte)
	if txProve["hash_array"], err = json.Marshal(proveHash); err != nil {
		return emptyJson
	}
	if block.Block.Header.BlockVersion > 2300 {
		if tx.Transaction.Result != nil && tx.Transaction.Result.ContractResult != nil {
			tx.Transaction.Result.ContractResult.GasUsed = 0
		}
	}
	if txProve["tx_byte"], err = proto.Marshal(tx); err != nil {
		return emptyJson
	}
	txProve["hash_type"] = []byte(chainConfig.Crypto.Hash)

	txProveByte, err := json.Marshal(txProve)
	if err != nil {
		return emptyJson
	}
	c.log.Debugf("[GetTxProve] %s", string(txProveByte))
	return string(txProveByte)
}

// errorFormat 错误格式化
//  @receiver c
//  @param methodName
//  @param err
//  @return error
func (c *ChainClient) errorFormat(methodName string, err error) error {
	msg := fmt.Sprintf("[%s]:%s", methodName, err.Error())
	c.log.Errorf(msg)
	return errors.New(msg)
}

// getKvsFromKvJsonStr 从参数获取合约参数
//  @receiver c
//  @param kvJsonStr
//  @return []*common.KeyValuePair
//  @return error
func (c *ChainClient) getKvsFromKvJsonStr(kvJsonStr string) ([]*common.KeyValuePair, error) {
	var kvMap map[string]string
	err := json.Unmarshal([]byte(kvJsonStr), &kvMap)
	if err != nil {
		errStr := fmt.Sprintf("[getKvsFromKvJsonStr] kvJsonStr must be json string: %s -> %s",
			kvJsonStr, err.Error())
		c.log.Error(errStr)
		return nil, errors.New(errStr)
	}
	kvs := []*common.KeyValuePair{}
	for k, v := range kvMap {
		kv := &common.KeyValuePair{
			Key:   k,
			Value: []byte(v),
		}
		kvs = append(kvs, kv)
	}
	return kvs, nil
}

// CheckChain 检查链的连通性
//  @receiver c
//  @return bool
func (c *ChainClient) CheckChain() bool {
	chainConfigList, _ := chain_config.ChainConfigManager.Get("")
	for _, chainConfig := range chainConfigList {
		if _, ok := c.client[chainConfig.ChainRid]; !ok {
			c.log.Errorf("[CheckChain] %s client not found", chainConfig.ChainRid)
			return false
		}
		if _, err := c.client[chainConfig.ChainRid].GetChainMakerServerVersion(); err != nil {
			c.log.Errorf("[CheckChain] %s", err.Error())
			return false
		}
	}
	return true
}

// keepClientAlive 连接性检查
//  @receiver c
//func (c *ChainClient) keepClientAlive() {
//	// 在实际使用过程中有连接失败的问题，所以暂时加了一个携程来保证连接不会断掉
//	for {
//		time.Sleep(time.Second)
//		chainConfigList, _ := chain_config.ChainConfigManager.Get("")
//		for _, chainConfig := range chainConfigList {
//			if _, err := c.client[chainConfig.ChainRid].GetChainMakerServerVersion(); err != nil {
//				// 这里是为了防止重连失败导致chainMakerClient为nil,报错，所以要反复重连
//				for {
//					// 先建立连接
//					client, err := newClient(chainConfig, c.log)
//					if err != nil {
//						c.log.Error("[InitChainmakerClient] To recreate the chain rid [%s] client Error: %s",
//							chainConfig.ChainRid, err.Error())
//						time.Sleep(time.Second)
//						continue
//					}
//					// 后释放旧连接
//					_ = c.client[chainConfig.ChainRid].Stop()
//					c.client[chainConfig.ChainRid] = client
//					c.log.Debugf("[InitChainmakerClient] To recreate the chain rid [%s] client success",
//						chainConfig.ChainRid)
//
//					if conf.Config.BaseConfig.TxVerifyType == conf.SpvTxVerify {
//						go c.listenBlockHeader(chainConfig.ChainRid)
//					}
//					break
//				}
//				for listenKey := range c.listen {
//					chainRid, contractName := getChainIdAndContractName(listenKey)
//					if chainRid == "" {
//						continue
//					}
//					lastBlock, err := c.client[chainRid].GetLastBlock(false)
//					if err != nil {
//						c.log.Errorf("[InitChainClient] get last block error: %s\n", err.Error())
//					}
//					go c.listenEvent(
//						contractName,
//						chainRid, int64(lastBlock.Block.Header.BlockHeight))
//				}
//			}
//		}
//	}
//}

// getChainClient 获取链的客户端对象
//  @receiver c
//  @param chainRid
//  @return *sdk.ChainClient
//  @return error
func (c *ChainClient) getChainClient(chainRid string) (*sdk.ChainClient, error) {
	if _, ok := c.client[chainRid]; !ok {
		msg := fmt.Sprintf("[getChainClient] no chain client: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return nil, fmt.Errorf(msg)
	}
	return c.client[chainRid], nil
}

func (c *ChainClient) getLaseBlockHeaderHeight(chainRid string) int64 {
	heightByte, err := db.Db.Get(utils.ParseHeaderKey(chainRid))
	if err != nil {
		return 0
	}
	if len(heightByte) == 0 {
		return 0
	}
	height, err := strconv.Atoi(string(heightByte))
	if err != nil {
		return 0
	}
	return int64(height)
}

func (c *ChainClient) resetLaseBlockHeaderHeight(chainRid string) error {
	return db.Db.Put(utils.ParseHeaderKey(chainRid), []byte("0"))
}
