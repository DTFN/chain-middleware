/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	chain_config "chainmaker.org/chainmaker/tcip-fabric/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/utils"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"

	"github.com/hyperledger/fabric-protos-go/peer"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/request"

	"go.uber.org/zap"

	tcipEvent "chainmaker.org/chainmaker/tcip-fabric/v2/module/event"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
)

// ChainClientItfc 链客户端接口
type ChainClientItfc interface {
	InvokeContract(
		chainRid, contractName, method string,
		args [][]byte, needTx bool) (string, string, *peer.ProcessedTransaction, error)
	GetTxProve(tx *peer.ProcessedTransaction, txId, chainRid string) string
	TxProve(txProve string) bool
	CheckChain() bool
}

// ChainClient 链交互对象结构体
type ChainClient struct {
	channelClient map[string]*channel.Client
	ledgerClinet  map[string]*ledger.Client
	eventClient   map[string]*event.Client
	listenCtx     map[string]fab.Registration
	listen        map[string]bool
	peers         map[string][]string
	//contractChannel map[string]*channel.Client
	log *zap.SugaredLogger
}

// ChainClientV1 链交互对象
var ChainClientV1 ChainClientItfc

const success = "success"

// InitChainClient 初始化链
//  @return error
func InitChainClient() error {
	log := logger.GetLogger(logger.ModuleChainClient)
	log.Debug("[InitChainClient] init")

	fabricClient := &ChainClient{
		channelClient: make(map[string]*channel.Client),
		ledgerClinet:  make(map[string]*ledger.Client),
		eventClient:   make(map[string]*event.Client),
		listen:        make(map[string]bool),
		listenCtx:     make(map[string]fab.Registration),
		peers:         make(map[string][]string),
		//contractChannel: make(map[string]*channel.Client),
		log: logger.GetLogger(logger.ModuleChainClient),
	}
	chainConfigList, err := chain_config.ChainConfigManager.Get("")
	if err != nil {
		panic(fmt.Sprintf("get chain config error: %v", err))
	}
	for _, chainConfig := range chainConfigList {
		log.Debugf("[InitChainClient] chainConfig: %+v", chainConfig)
		state := true
		msg := success
		channelClient, ledgerClinet, eventClient, peers, err := createSdk(chainConfig, log)
		if err != nil {
			state = false
			msg = err.Error()
			log.Errorf("[InitChainClient] newClient error: %s %v", chainConfig.ChainRid, err)
		}
		err1 := chain_config.ChainConfigManager.SetState(chainConfig, state, msg)
		if err1 != nil {
			log.Errorf("[chain config watch] SetState %s", err1.Error())
			continue
		}
		if !state {
			continue
		}

		log.Infof("[InitChainClient] SetState %t, chain rid %s",
			state, chainConfig.ChainRid)

		fabricClient.channelClient[chainConfig.ChainRid] = channelClient
		fabricClient.ledgerClinet[chainConfig.ChainRid] = ledgerClinet
		fabricClient.eventClient[chainConfig.ChainRid] = eventClient
		fabricClient.peers[chainConfig.ChainRid] = peers
		log.Debugf("[InitChainClient] create chain [%s] client success", chainConfig.ChainRid)
	}
	eventList, _ := tcipEvent.EventManagerV1.GetEvent("")
	for _, crossChainevent := range eventList {
		listenKey := getListenKey(crossChainevent.ChainRid, crossChainevent.ContractName)
		if _, ok := fabricClient.listen[listenKey]; ok {
			continue
		}
		fabricClient.listen[listenKey] = true
		go fabricClient.listenEvent(crossChainevent.CrossChainEventId, crossChainevent.ChainRid, crossChainevent.ContractName)
		time.Sleep(time.Millisecond * 100)
	}
	ChainClientV1 = fabricClient
	go fabricClient.eventStart()
	go fabricClient.chainConfigStart()
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
					c.log.Warn("[event watch]error event")
				}
				c.log.Infof("[event watch] receive chain config update: %v", eventInfo)
				listenKey := getListenKey(eventInfo.ChainRid, eventInfo.ContractName)
				_, has := c.listen[listenKey]
				if eventInfo.Operate == tcipcommon.Operate_SAVE && !has {
					go c.listenEvent(eventInfo.CrossChainEventID, eventInfo.ChainRid, eventInfo.ContractName)
				}
				if eventInfo.Operate == tcipcommon.Operate_SAVE && has {
					_ = tcipEvent.EventManagerV1.SetEventState(eventInfo.CrossChainEventID, true, success)
				}
				if eventInfo.Operate == tcipcommon.Operate_DELETE && has {
					client, err := c.getChainClient(eventInfo.ChainRid)
					if err != nil {
						c.log.Errorf("[event watch] %s", err.Error())
					} else {
						client.UnregisterChaincodeEvent(c.listenCtx[listenKey])
					}
					delete(c.listenCtx, listenKey)
					delete(c.listen, listenKey)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("[event watch]have close event happened")
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
					c.log.Warn("[chain config watch]error event")
				}
				c.log.Infof("[chain config watch] receive chain config update: %v", chainConfigOperate)
				_, has := c.channelClient[chainConfigOperate.ChainRid]
				chainConfig, err := chain_config.ChainConfigManager.Get(chainConfigOperate.ChainRid)
				c.log.Debugf("chainConfig: %+v", chainConfig)
				if err != nil {
					continue
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_SAVE {
					if has {
						delete(c.channelClient, chainConfigOperate.ChainRid)
						delete(c.ledgerClinet, chainConfigOperate.ChainRid)
						delete(c.eventClient, chainConfigOperate.ChainRid)
						delete(c.peers, chainConfigOperate.ChainRid)
					}
					channelClient, ledgerClinet, eventClient, peers, err := createSdk(chainConfig[0], c.log)
					if err != nil {
						c.log.Errorf("[chain config watch] %s", err.Error())
						err1 := chain_config.ChainConfigManager.SetState(chainConfig[0], false, err.Error())
						if err1 != nil {
							c.log.Errorf("[chain config watch] SetState %s", err1.Error())
							continue
						}
					} else {
						err1 := chain_config.ChainConfigManager.SetState(chainConfig[0], true, "success")
						if err1 != nil {
							c.log.Errorf("[chain config watch] SetState %s", err1.Error())
							continue
						}
						c.channelClient[chainConfig[0].ChainRid] = channelClient
						c.ledgerClinet[chainConfig[0].ChainRid] = ledgerClinet
						c.eventClient[chainConfig[0].ChainRid] = eventClient
						c.peers[chainConfig[0].ChainRid] = peers
					}
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_DELETE && has {
					delete(c.channelClient, chainConfigOperate.ChainRid)
					delete(c.ledgerClinet, chainConfigOperate.ChainRid)
					delete(c.eventClient, chainConfigOperate.ChainRid)
					delete(c.peers, chainConfigOperate.ChainRid)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("[chain config watch]have close event happened")
				close(utils.EventChan)
				return
			}
		}
	}()
}

// InvokeContract 调用合约
//  @receiver c
//  @param chainRid
//  @param contractName
//  @param method
//  @param args
//  @param needTx
//  @return string
//  @return string
//  @return *peer.ProcessedTransaction
//  @return error
func (c *ChainClient) InvokeContract(
	chainRid, contractName, method string,
	args [][]byte, needTx bool) (string, string, *peer.ProcessedTransaction, error) {
	chainClient, err := c.getChainClient(chainRid)
	if err != nil {
		c.log.Errorf("[InvokeContract] %s", err.Error())
		return "", "", nil, err
	}
	peers, err := c.getPeers(chainRid)
	if err != nil {
		c.log.Errorf("[InvokeContract] %s", err.Error())
		return "", "", nil, err
	}

	c.log.Debugf("[InvokeContract] chainRid: %s, peers: %+v", chainRid, peers)

	// peers 指定peer节点
	// peers = []string{"peer0.org1.example.com", "peer0.org2.example.com"}

	response, err := chainClient.Execute(channel.Request{
		ChaincodeID: contractName,
		Fcn:         method,
		Args:        args,
	}, channel.WithTargetEndpoints(peers...))

	if err != nil {
		errMsg := fmt.Sprintf(
			"[InvokeContract]invoke contract failed, res:%+v\n", err.Error())
		c.log.Error(errMsg)
		return "", "", nil, errors.New(errMsg)
	}
	c.log.Debugf("[InvokeContract] resp: %v", response)
	if response.TxValidationCode != 0 {
		if len(response.Payload) > 0 {
			errMsg := fmt.Sprintf(
				"[InvokeContract] %s", response.Payload)
			c.log.Errorf(errMsg)
			return "", "", nil, errors.New(errMsg)
		}
		return "", "", nil, fmt.Errorf("[InvokeContract] invoke fabric contract error, code = %v", response.TxValidationCode)
	}
	if needTx {
		ledgerClient, err := c.getLedgerClient(chainRid)
		if err != nil {
			c.log.Errorf("[InvokeContract] %s", err.Error())
			return "", "", nil, err
		}
		tx, err := ledgerClient.QueryTransaction(response.TransactionID, ledger.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com"))
		if err != nil {
			errMsg := fmt.Sprintf("[InvokeContract] can't get tx: txId %s, error %s",
				response.TransactionID, err.Error())
			c.log.Errorf(errMsg)
			return "", "", nil, errors.New(errMsg)
		}
		return string(response.Payload), string(response.TransactionID), tx, nil
	}
	return string(response.Payload), string(response.TransactionID), nil, nil
}

// listenEvent 监听合约事件
//  @receiver c
//  @param chainRid
//  @param contractName
func (c *ChainClient) listenEvent(crossChainEventId, chainRid, contractName string) {
	c.log.Debugf("[listenEvent]  listen start %s %s", chainRid, contractName)
	eventClient, err := c.getEventClient(chainRid)
	if err != nil {
		msg := fmt.Sprintf("[listenEvent] %s", err.Error())
		_ = tcipEvent.EventManagerV1.SetEventState(crossChainEventId, false, msg)
		c.log.Error(msg)
		return
	}
	chainClient, err := c.getChainClient(chainRid)
	if err != nil {
		msg := fmt.Sprintf("[listenEvent] %s", err.Error())
		_ = tcipEvent.EventManagerV1.SetEventState(crossChainEventId, false, msg)
		return
	}
	ledgerClient, err := c.getLedgerClient(chainRid)
	if err != nil {
		msg := fmt.Sprintf("[listenEvent] %s", err.Error())
		_ = tcipEvent.EventManagerV1.SetEventState(crossChainEventId, false, msg)
		return
	}
	// 监听所有的topic,这样就不需要每次都改配置了
	re := regexp.MustCompile(`.*`)
	registration, notifier, err := eventClient.RegisterChaincodeEvent(contractName, re.String())
	if err != nil {
		msg := fmt.Sprintf("[listenEvent]failed to register chaincode event: %s, %s", contractName, err.Error())
		c.log.Error(msg)
		_ = tcipEvent.EventManagerV1.SetEventState(crossChainEventId, false, msg)
		return
	}
	defer chainClient.UnregisterChaincodeEvent(registration)
	c.listenCtx[getListenKey(chainRid, contractName)] = registration
	c.listen[getListenKey(chainRid, contractName)] = true
	_ = tcipEvent.EventManagerV1.SetEventState(crossChainEventId, true, success)
	//for {
	//	select {
	//	case ccEvent, ok := <-notifier:
	//		if !ok {
	//			go c.listenEvent(chainRid, contractName)
	//		}
	//		tx, err := ledgerClient.QueryTransaction(fab.TransactionID(ccEvent.TxID))
	//		if err != nil {
	//			c.log.Errorf("[listenEvent] can't get tx: txId %s topic %s event %+v",
	//				ccEvent.TxID, ccEvent.EventName, ccEvent)
	//		}
	//		txProve := c.GetTxProve(tx, ccEvent.TxID, chainRid)
	//		eventInfo := &tcipEvent.EventInfo{
	//			Topic:        ccEvent.EventName,
	//			chainRid:     chainRid,
	//			ContractName: ccEvent.ChaincodeID,
	//			TxProve:      txProve,
	//			Data:         []string{string(ccEvent.Payload)},
	//			Tx:           tx.TransactionEnvelope.Payload,
	//			TxId:         ccEvent.TxID,
	//			BlockHeight:  int64(ccEvent.BlockNumber),
	//		}
	//
	//		c.log.Infof("[listenEvent] eventInfo: %v\n", eventInfo)
	//		go request.RequestV1.BeginCrossChain(eventInfo)
	//	}
	//}

	for ccEvent := range notifier {
		tx, err := ledgerClient.QueryTransaction(fab.TransactionID(ccEvent.TxID), ledger.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com"))
		if err != nil {
			c.log.Errorf("[listenEvent] can't get tx: txId %s, topic %s, err %+v, event %+v",
				ccEvent.TxID, ccEvent.EventName, err, ccEvent)
			c.log.Debugf("[listenEvent] ledgerClient %+v", ledgerClient)
		} else {
			txProve := c.GetTxProve(tx, ccEvent.TxID, chainRid)
			eventInfo := &tcipEvent.EventInfo{
				Topic:        ccEvent.EventName,
				ChainRid:     chainRid,
				ContractName: ccEvent.ChaincodeID,
				TxProve:      txProve,
				Data:         []string{string(ccEvent.Payload)},
				Tx:           tx.TransactionEnvelope.Payload,
				TxId:         ccEvent.TxID,
				BlockHeight:  int64(ccEvent.BlockNumber),
			}

			c.log.Infof("[listenEvent] eventInfo: %v\n", eventInfo.ToString())
			go request.RequestV1.BeginCrossChain(eventInfo)
		}
	}
}

// GetTxProve 获取交易认证信息
//  @receiver c
//  @param tx
//  @param txId
//  @param chainRid
//  @return string
func (c *ChainClient) GetTxProve(tx *peer.ProcessedTransaction, txId, chainRid string) string {
	txBase64 := base64.StdEncoding.EncodeToString(tx.TransactionEnvelope.Payload)
	txProveStr := fmt.Sprintf("{\"tx_id\": \"%s\", \"tx_byte\":\"%s\", \"chain_rid\":\"%s\"}",
		txId, txBase64, chainRid)
	txProve := cross_chain.TxVerifyRequest{
		Version: common.Version_V1_0_0,
		TxProve: txProveStr,
	}
	res, err := json.Marshal(txProve)
	if err != nil {
		return ""
	}
	return string(res)
}

// TxProve 交易认证 (fabric的方式)
//  @receiver c
//  @param txProve
//  @return bool
func (c *ChainClient) TxProve(txProve string) bool {
	c.log.Debugf("[txProve]: %s\n", txProve)
	txProveMap := make(map[string]string)
	err := json.Unmarshal([]byte(txProve), &txProveMap)
	if err != nil {
		c.log.Debugf("[txProve]: false. err != nil", txProve)
		return false
	}
	chainRid, ok := txProveMap["chain_rid"]
	if !ok {
		c.log.Debugf("[txProve]: false. not have chainRid", txProve)
		return false
	}
	txId, ok := txProveMap["tx_id"]
	if !ok {
		c.log.Debugf("[txProve]: false. not have tx_id", txProve)
		return false
	}
	txByteString, ok := txProveMap["tx_byte"]
	if !ok {
		c.log.Debugf("[txProve]: false. not have txByteString", txProve)
		return false
	}

	txByte, err := base64.StdEncoding.DecodeString(txByteString)
	if err != nil {
		c.log.Debugf("[txProve]: false. base64.StdEncoding.DecodeString(txByteString) fail", txProve)
		return false
	}
	ledgerClient, err := c.getLedgerClient(chainRid)
	if err != nil {
		c.log.Errorf("[TxProve] %s", err.Error())
		return false
	}
	tx, err := ledgerClient.QueryTransaction(fab.TransactionID(txId), ledger.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com"))
	if err != nil {
		c.log.Debugf("[txProve]: false. ledgerClient.QueryTransaction fail", txProve)
		return false
	}
	if string(txByte) != "" && tx.TransactionEnvelope.Payload != nil &&
		string(txByte) == string(tx.TransactionEnvelope.Payload) {
		return true
	}

	c.log.Debugf("[txProve]: false", txProve)
	return false
}

// CheckChain 检查链的连通性
//  @receiver c
//  @return bool
func (c *ChainClient) CheckChain() bool {
	for _, ledgerClient := range c.ledgerClinet {
		if _, err := ledgerClient.QueryConfig(ledger.WithTargetEndpoints("peer0.org1.example.com", "peer0.org2.example.com")); err != nil {
			return false
		}
	}
	return true
}

// getChainClient 获取链连接
//  @receiver c
//  @param chainRid
//  @return *channel.Client
//  @return error
func (c *ChainClient) getChainClient(chainRid string) (*channel.Client, error) {
	if _, ok := c.channelClient[chainRid]; !ok {
		msg := fmt.Sprintf("[getChainClient] no chain client: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return nil, fmt.Errorf(msg)
	}
	return c.channelClient[chainRid], nil
}

// getLedgerClient 获取链连接
//  @receiver c
//  @param chainRid
//  @return *ledger.Client
//  @return error
func (c *ChainClient) getLedgerClient(chainRid string) (*ledger.Client, error) {
	if _, ok := c.ledgerClinet[chainRid]; !ok {
		msg := fmt.Sprintf("[getLedgerClient] no ledger client: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return nil, fmt.Errorf(msg)
	}
	return c.ledgerClinet[chainRid], nil
}

// getEventClient 获取链连接
//  @receiver c
//  @param chainRid
//  @return *event.Client
//  @return error
func (c *ChainClient) getEventClient(chainRid string) (*event.Client, error) {
	if _, ok := c.eventClient[chainRid]; !ok {
		msg := fmt.Sprintf("[getEventClient] no event client: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return nil, fmt.Errorf(msg)
	}
	return c.eventClient[chainRid], nil
}

// getPeers 获取链连接
//  @receiver c
//  @param chainRid
//  @return []string
//  @return error
func (c *ChainClient) getPeers(chainRid string) ([]string, error) {
	if _, ok := c.peers[chainRid]; !ok {
		msg := fmt.Sprintf("[getPeers] no peer: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return nil, fmt.Errorf(msg)
	}
	return c.peers[chainRid], nil
}

// getListenKey 构建监听缓存的key
//  @param chainRid
//  @param contracrName
//  @return string
func getListenKey(chainRid, contracrName string) string {
	return fmt.Sprintf("%s#%s", chainRid, contracrName)
}
