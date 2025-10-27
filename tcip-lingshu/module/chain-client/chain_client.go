/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chain_client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	chain_config "chainmaker.org/chainmaker/tcip-lingshu/v2/module/chain-config"

	//"github.com/ethereum/go-ethereum/crypto"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/utils"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/event"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/request"

	"go.uber.org/zap"

	tcipcommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/logger"
	bcosabi "github.com/FISCO-BCOS/go-sdk/abi"

	bcostypes "github.com/FISCO-BCOS/go-sdk/core/types"
	bcoscommon "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
)

const (
	toBlock         = "latest"
	eventTopicSplit = "("
	success         = "success"
)

// ChainClientItfc 链客户端接口
type ChainClientItfc interface {
	// InvokeContract 调用合约
	InvokeContract(chainRid, contractName, method, abiStr string, args string,
		needTx bool, paramType []tcipcommon.EventDataType) ([]string, *bcostypes.TransactionDetail, error)
	// InvokeContract 调用通用
	DIDManagerUpdate(ledgerId int, contractAddress string, did string, didDoc string) (*bcostypes.Receipt, error)
	InvokeContractByVc(ledgerId int, contractAddress string, funcName string, vc string) ([]string, *bcostypes.TransactionDetail, error)
	// GetTxProve 获取交易凭证
	GetTxProve(tx *bcostypes.TransactionDetail, chainRid string) string
	// TxProve 交易验证
	TxProve(txProve string) bool
	// CheckChain 验证了链的连通性
	CheckChain() bool
}

// ChainClient 链客户端结构体
type ChainClient struct {
	// 缓存链的客户端对象
	client map[string]bool
	// 缓存已经在监听的事件
	listen map[string]bool
	// 日志对象
	log *zap.SugaredLogger
}

// ChainClientV1 连交互模块对象
var ChainClientV1 ChainClientItfc

const (
	emptyJson = "{}"
)

// InitChainClient 初始化链客户端
//
//	@return error
func InitChainClient() error {
	log := logger.GetLogger(logger.ModuleChainmakerClient)
	log.Debug("[InitChainClient] init")
	bcosClient := &ChainClient{
		client: make(map[string]bool),
		listen: make(map[string]bool),
		log:    logger.GetLogger(logger.ModuleChainClient),
	}
	utils.EventChan = make(chan *utils.EventOperate, 10)
	utils.UpdateChainConfigChan = make(chan *utils.ChainConfigOperate, 10)
	eventList, _ := event.EventManagerV1.GetEvent("")
	chainConfigList, err := chain_config.ChainConfigManager.Get("")
	if err != nil {
		panic(fmt.Sprintf("get chain config error: %v", err))
	}
	for _, chainConfig := range chainConfigList {
		state := true
		msg := success
		flag, err := createSDK(chainConfig, log)
		if err != nil {
			state = false
			msg = err.Error()
			log.Errorf("[InitChainClient] Create chain client error failed, err: %v %t", msg, state)
		}
		err1 := chain_config.ChainConfigManager.SetState(chainConfig, state, msg)
		if err1 != nil {
			log.Errorf("[chain config watch] SetState %s", err1.Error())
			continue
		}
		if !state {
			continue
		}
		log.Debugf("[InitChainClient] create chain [%s] client success", chainConfig.ChainRid)

		bcosClient.client[chainConfig.ChainRid] = flag
	}
	for _, crossChainevent := range eventList {
		listenKey := getListenKey(crossChainevent.ChainRid, crossChainevent.ContractName, crossChainevent.EventName)
		if _, ok := bcosClient.listen[listenKey]; ok {
			continue
		}
		err := bcosClient.listenEvent(crossChainevent)
		if err != nil {
			msg := fmt.Sprintf("[InitChainClient] listen event [%+v] failed, err: %v", crossChainevent, err)
			_ = event.EventManagerV1.SetEventState(crossChainevent.CrossChainEventId, false, msg)
			log.Error(msg)
			continue
		}
		_ = event.EventManagerV1.SetEventState(crossChainevent.CrossChainEventId, true, success)
	}
	ChainClientV1 = bcosClient
	go bcosClient.evenStart()
	go bcosClient.chainConfigStart()
	return nil
}

// listenEvent 监听合约事件
//
//	@receiver c
//	@param dbEvent
//	@return error
func (c *ChainClient) listenEvent(dbEvent *tcipcommon.CrossChainEvent) error {
	// 有就不跑了
	listenKey := getListenKey(dbEvent.ChainRid, dbEvent.ContractName, dbEvent.EventName)
	fmt.Println(listenKey)
	if c.listen[getListenKey(dbEvent.ChainRid, dbEvent.ContractName, dbEvent.EventName)] {
		return nil
	}
	//client, err := c.getChainClient(dbEvent.ChainRid)
	//if err != nil {
	//	msg := fmt.Sprintf("[listenEvent] chain client error: %s\n", err.Error())
	//	c.log.Error(msg)
	//	return errors.New(msg)
	//}
	//var ctx context.Context
	//startBlock, err := client.GetBlockNumber(ctx)
	//if err != nil {
	//	c.log.Errorf("[listenEvent] GetBlockNumber error: %s\n", err.Error())
	//	return err
	//}
	//eventLogParams := bcostypes.EventLogParams{
	//	FromBlock: fmt.Sprintf("%d", startBlock),
	//	ToBlock:   toBlock,
	//	GroupID:   fmt.Sprintf("%d", client.GetGroupID()),
	//	Topics:    []string{bcoscommon.BytesToHash(crypto.Keccak256([]byte(dbEvent.EventName))).Hex()},
	//	Addresses: []string{dbEvent.ContractName},
	//}
	//err = client.SubscribeEventLogs(eventLogParams, func(status int, logs []bcostypes.Log) {

	c.log.Infof("[listenEvent] listen ChainRid %s success: eventName %s address %s",
		dbEvent.ChainRid, dbEvent.EventName, dbEvent.ContractName)
	c.listen[getListenKey(dbEvent.ChainRid, dbEvent.ContractName, dbEvent.EventName)] = true
	go c.SubscribeEventLogs(func(status int, logs []bcostypes.Log) {
		// 官方没有停止的接口，只能自行判断
		if !c.listen[getListenKey(dbEvent.ChainRid, dbEvent.ContractName, dbEvent.EventName)] {
			return
		}
		logRes, err2 := json.MarshalIndent(logs, "", "  ")
		if err2 != nil {
			c.log.Warnf("[listenEvent] logs marshalIndent error: %v", err2)
		}
		c.log.Debugf("[listenEvent] received: %s\n", logRes)
		args := make(bcosabi.Arguments, 0)
		for _, abi := range dbEvent.Abi {
			arg := &bcosabi.Argument{}
			err2 = arg.UnmarshalJSON([]byte(abi))
			if err2 != nil {
				msg := fmt.Sprintf("[listenEvent] UnmarshalJSON abi [%s] error: %s", abi, err2.Error())
				c.log.Errorf(msg)
				return
			}

			args = append(args, *arg)
		}
		tmpArr := strings.Split(dbEvent.EventName, eventTopicSplit)
		if len(tmpArr) == 0 {
			msg := fmt.Sprintf("[listenEvent] split topic [%s] error: %s", dbEvent.EventName, err2.Error())
			c.log.Errorf(msg)
			return
		}
		eve := &bcosabi.Event{
			Name:      tmpArr[0],
			RawName:   tmpArr[0],
			Anonymous: false,
			SMCrypto:  false,
			Inputs:    args,
		}
		text, err2 := eve.Inputs.UnpackValues(logs[0].Data)
		if err2 != nil {
			msg := fmt.Sprintf("[listenEvent] UnpackValues event data [%s] error: %s",
				dbEvent.EventName, err2.Error())
			c.log.Errorf(msg)
			return
		}
		data := fmt.Sprintf("%s", text)
		// 前后各去掉一个字符
		data = data[1:]
		data = data[:len(data)-1]

		if len(data) == 0 {
			msg := fmt.Sprintf("[listenEvent] nil data [%s], %s", dbEvent.EventName, fmt.Sprintf("%s", text))
			c.log.Warn(msg)
			return
		}
		//var ctx1 context.Context
		//tx, err2 := client.GetTransactionByHash(ctx1, logs[0].TxHash)
		tx, err2 := c.GetTransactionByHash(logs[0].TxHash.String())
		if err2 != nil {
			msg := fmt.Sprintf("[listenEvent] get tx error [%s]", dbEvent.EventName)
			c.log.Warn(msg)
			return
		}
		txProve := c.GetTxProve(tx, dbEvent.ChainRid)
		var txByte []byte
		if txByte, err2 = json.Marshal(tx); err2 != nil {
			msg := fmt.Sprintf("[listenEvent] Marshal tx error [%s]", dbEvent.EventName)
			c.log.Warn(msg)
			return
		}
		eventData := make([]string, 0)
		for _, v := range logs[0].Topics {
			eventData = append(eventData, v.Hex())
		}
		eventData = append(eventData, strings.Split(data, " ")...)
		eventInfo := &event.EventInfo{
			Topic:        dbEvent.EventName,
			ChainRid:     dbEvent.ChainRid,
			ContractName: dbEvent.ContractName,
			TxProve:      txProve,
			Data:         eventData,
			Tx:           txByte,
			TxId:         logs[0].TxHash.Hex(),
			BlockHeight:  int64(logs[0].BlockNumber),
		}
		c.log.Infof("[listenEvent] eventInfo: %v\n", eventInfo.ToString())

		go request.RequestV1.BeginCrossChain(eventInfo)
	})
	//if err != nil {
	//	c.log.Errorf("[listenEvent] listen ChainRid %s error: %s", dbEvent.ChainRid, err.Error())
	//	return fmt.Errorf("[listenEvent] listen ChainRid %s error: %s", dbEvent.ChainRid, err.Error())
	//}
	return nil
}

// evenStart 监听事件更新
//
//	@receiver c
func (c *ChainClient) evenStart() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case eventInfo, ok := <-utils.EventChan:
				if !ok {
					c.log.Warn("[event watch] error event")
				}
				c.log.Infof("[event watch] receive event: %v", eventInfo)
				has := true
				dbEvent, err := event.EventManagerV1.GetEvent(eventInfo.CrossChainEventID)
				if len(dbEvent) == 0 {
					if err != nil {
						c.log.Errorf("[event watch] get last block error: %s\n", err.Error())
					} else {
						c.log.Errorf("[event watch] get last block error")
					}
					has = false
					continue
				}
				listenKey := getListenKey(eventInfo.ChainRid, dbEvent[0].ContractName, dbEvent[0].EventName)
				if eventInfo.Operate == tcipcommon.Operate_SAVE && has {
					err := c.listenEvent(dbEvent[0])
					if err != nil {
						c.log.Errorf("[event watch] listen event [%+v] failed, err: %v",
							dbEvent[0], err)
						continue
					}
					_ = event.EventManagerV1.SetEventState(eventInfo.CrossChainEventID, true, success)
				}
				if eventInfo.Operate == tcipcommon.Operate_DELETE && !has {
					c.listen[listenKey] = false
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("[event watch] have close event happened")
				close(utils.EventChan)
				return
			}
		}
	}()
}

// chainConfigStart 监听chainconfig更新
//
//	@receiver c
func (c *ChainClient) chainConfigStart() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		for {
			select {
			case chainConfigOperate, ok := <-utils.UpdateChainConfigChan:
				if !ok {
					c.log.Warn("[chain config watch] error event")
				}
				c.log.Infof("[chain config watch] receive event: %v", chainConfigOperate)
				_, has := c.client[chainConfigOperate.ChainRid]
				chainConfig, err := chain_config.ChainConfigManager.Get(chainConfigOperate.ChainRid)
				if err != nil {
					continue
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_SAVE {
					if has {
						delete(c.client, chainConfigOperate.ChainRid)
					}
					client, err := createSDK(chainConfig[0], c.log)
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
						c.client[chainConfig[0].ChainRid] = client
					}
				}
				if chainConfigOperate.Operate == tcipcommon.Operate_DELETE && has {
					delete(c.client, chainConfigOperate.ChainRid)
				}
			case <-ch:
				// 有关闭事件
				c.log.Info("[chain config watch] have close event happened")
				close(utils.EventChan)
				return
			}
		}
	}()
}

// InvokeContract 调用合约
//
//	@receiver c
//	@param chainRid 链资源id
//	@param contractName 合约名称
//	@param method 调用方法
//	@param abiStr abi
//	@param args 参数
//	@param needTx 是否需要交易
//	@param paramType 参数类型
//	@return []string 返回参数
//	@return *bcostypes.TransactionDetail 交易
//	@return error 错误信息
func (c *ChainClient) InvokeContract(chainRid, contractName, method, abiStr string, args string,
	needTx bool, paramType []tcipcommon.EventDataType) ([]string, *bcostypes.TransactionDetail, error) {
	argsArr, err := dealParam(args, paramType)
	if err != nil {
		msg := fmt.Sprintf("[InvokeContract] dealParam error: %s\n", err.Error())
		c.log.Error(msg)
		return nil, nil, errors.New(msg)
	}
	//client, err := c.getChainClient(chainRid)
	//if err != nil {
	//	msg := fmt.Sprintf("[InvokeContract] chain client error: %s\n", err.Error())
	//	c.log.Error(msg)
	//	return nil, nil, errors.New(msg)
	//}

	//address := bcoscommon.HexToAddress(contractName)
	parsed, err := bcosabi.JSON(strings.NewReader(abiStr))
	if err != nil {
		msg := fmt.Sprintf("[InvokeContract] abi [%s] read error: %s", abiStr, err.Error())
		c.log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	//_, receipt, err := bcosbind.NewBoundContract(address, parsed, client, client, client).
	//	Transact(client.GetTransactOpts(), method, argsArr...)
	receipt, err := c.sendContractInvoke(contractName, abiStr, method, argsArr)

	if err != nil {
		msg := fmt.Sprintf("[InvokeContract] invoke contract [%s %s %s] error: %s\n, abi: %s, args: %v",
			chainRid, contractName, method, err.Error(), abiStr, args)
		c.log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	c.log.Debugf("[InvokeContract] invoke contract [%s %s %s] resp: %v\n, abi: %s, args: %v",
		chainRid, contractName, method, receipt, abiStr, args)
	if receipt.Status != bcostypes.Success {
		//if receipt.Status != 26 {
		msg := fmt.Sprintf("[InvokeContract] invoke contract [%s %s %s] error: %s\n, abi: %s, args: %v",
			chainRid, contractName, method, "status error", abiStr, args)
		c.log.Error(msg)
		return nil, nil, errors.New(msg)
	}

	resArr := make([]string, 0)
	if _, ok := parsed.Methods[method]; ok && len(parsed.Methods[method].Outputs) != 0 {
		b, err := hex.DecodeString(receipt.Output[2:])
		if err != nil {
			msg := fmt.Sprintf("[InvokeContract] Decode output [%s %s %s] error: %s\n, abi: %s, args: %v",
				chainRid, contractName, method, err.Error(), abiStr, args)
			c.log.Error(msg)
			return nil, nil, errors.New(msg)
		}
		methodApi := bcosabi.Method{
			Name:     method,
			RawName:  method,
			SMCrypto: false,
			Outputs:  parsed.Methods[method].Outputs,
		}
		text, err := methodApi.Outputs.UnpackValues(b)
		if err != nil {
			msg := fmt.Sprintf("[InvokeContract] UnpackValues output [%s %s %s] error: %s\n, abi: %s, args: %v",
				chainRid, contractName, method, err.Error(), abiStr, args)
			c.log.Error(msg)
			return nil, nil, errors.New(msg)
		}
		res := fmt.Sprintf("%s", text)
		res = res[1:]
		res = res[:len(res)-1]
		resArr = strings.Split(res, " ")
	}

	if needTx {
		//var ctx context.Context
		//tx, err := client.GetTransactionByHash(ctx, bcoscommon.HexToHash(receipt.TransactionHash))
		tx, err := c.GetTransactionByHash(receipt.GetTransactionHash())
		if err != nil {
			c.log.Debugf("[InvokeContract] get tx error [%s %s %s] error: %s\n, abi: %s, args: %v",
				chainRid, contractName, method, err.Error(), abiStr, args)
			msg := fmt.Sprintf("[InvokeContract] get tx error [%s]", err.Error())
			c.log.Warn(msg)
			return resArr, tx, nil
		}
		return resArr, tx, nil
	}

	// 如果是查询需要获取这里的结果，那么这里的值需要在跨链合约中进行处理，这里无法获取
	return resArr, nil, nil
}

func (c *ChainClient) GetBlockNumber() (int64, error) {
	frontIpPort := os.Getenv("FRONT_IP_PORT")
	if frontIpPort == "" {
		frontIpPort = "192.168.1.144:5007"
	}
	c.log.Debugf("[getBlockNumber] frontIpPort: %s", frontIpPort)

	resp, err := http.Get(fmt.Sprintf("http://%s/chain-proxy/1/rpc/blockNumber", frontIpPort))
	if err != nil {
		c.log.Errorf("[getBlockNumber] failed to send HTTP request: %s", err.Error())
		return -1, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}
	blockNum, err := strconv.ParseInt(strings.TrimSpace(string(body)), 0, 64)
	if err != nil {
		return -1, fmt.Errorf("failed to parse block number: %v", err)
	}
	return blockNum, nil
}

func (c *ChainClient) InvokeContractByVc(ledgerId int, contractAddress string, funcName string, vc string) ([]string, *bcostypes.TransactionDetail, error) {
	frontIpPort := os.Getenv("FRONT_IP_PORT")
	if frontIpPort == "" {
		frontIpPort = "192.168.1.144:5007"
	}
	c.log.Debugf("[InvokeContract] frontIpPort: %s", frontIpPort)

	payload := map[string]interface{}{
		"ledgerId": ledgerId,
		"address":  contractAddress,
		"funcName": funcName,
		"vc":       vc,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to marshal payload: %s", err.Error())
		return nil, nil, err
	}
	c.log.Debugf("[InvokeContract] payload: %s", string(payloadBytes))
	resp, err := http.Post(fmt.Sprintf("http://%s/chain-proxy/trans/vc-cal", frontIpPort),
		"application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to send HTTP request: %s", err.Error())
		return nil, nil, err
	}
	defer resp.Body.Close()

	// 直接解析为 bcostypes.Receipt
	var receipt *bcostypes.Receipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		c.log.Errorf("[InvokeContract] failed to decode HTTP response as receipt: %s", err.Error())
		return nil, nil, err
	}

	tx, err := c.GetTransactionByHash(receipt.GetTransactionHash())
	if err != nil {
		c.log.Debugf("[InvokeContract] get tx error TransactionHash: %s", receipt.GetTransactionHash())
		return []string{}, tx, nil
	}
	return []string{}, tx, nil
}

func (c *ChainClient) DIDManagerUpdate(ledgerId int, contractAddress string, did string, didDoc string) (*bcostypes.Receipt, error) {
	frontIpPort := os.Getenv("FRONT_IP_PORT")
	if frontIpPort == "" {
		frontIpPort = "192.168.1.144:5007"
	}
	c.log.Debugf("[InvokeContract] frontIpPort: %s", frontIpPort)

	payload := map[string]interface{}{
		"ledgerId": ledgerId,
		"address":  contractAddress,
		"did":      did,
		"didDoc":   didDoc,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to marshal payload: %s", err.Error())
		return nil, err
	}
	c.log.Debugf("[InvokeContract] payload: %s", string(payloadBytes))
	resp, err := http.Post(fmt.Sprintf("http://%s/chain-proxy/trans/did-update", frontIpPort),
		"application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to send HTTP request: %s", err.Error())
		return nil, err
	}
	c.log.Errorf("[InvokeContract] resp: %+v", resp)
	defer resp.Body.Close()

	// 直接解析为 bcostypes.Receipt
	var receipt *bcostypes.Receipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		c.log.Errorf("[InvokeContract] failed to decode HTTP response as receipt: %s", err.Error())
		return nil, err
	}
	return receipt, nil
}

func (c *ChainClient) sendContractInvoke(address string, abi string, method string, args []interface{}) (*bcostypes.Receipt, error) {
	frontIpPort := os.Getenv("FRONT_IP_PORT")
	if frontIpPort == "" {
		frontIpPort = "192.168.1.144:5007"
	}
	c.log.Debugf("[InvokeContract] frontIpPort: %s", frontIpPort)

	payload := map[string]interface{}{
		"contractAddress": address,
		"abi":             abi,
		"method":          method,
		"args":            args,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to marshal payload: %s", err.Error())
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/chain-proxy/1/rpc/invoke", frontIpPort),
		"application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		c.log.Errorf("[InvokeContract] failed to send HTTP request: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// 直接解析为 bcostypes.Receipt
	var receipt *bcostypes.Receipt
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		c.log.Errorf("[InvokeContract] failed to decode HTTP response as receipt: %s", err.Error())
		return nil, err
	}
	return receipt, nil
}

func (c *ChainClient) GetTransactionByHash(txhash string) (*bcostypes.TransactionDetail, error) {
	frontIpPort := os.Getenv("FRONT_IP_PORT")
	if frontIpPort == "" {
		frontIpPort = "192.168.1.144:5007"
	}
	c.log.Debugf("[GetTxDetail] frontIpPort: %s", frontIpPort)

	resp, err := http.Get(fmt.Sprintf("http://%s/chain-proxy/1/rpc/transaction/%s", frontIpPort, txhash))
	if err != nil {
		c.log.Errorf("[GetTxDetail] failed to send HTTP request: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	// 直接解析为 bcostypes.Receipt
	var receipt *bcostypes.TransactionDetail
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		c.log.Errorf("[GetTxDetail] failed to decode HTTP response as receipt: %s", err.Error())
		return nil, err
	}
	return receipt, nil
}

// GetTxProve 获取交易证明
//
//	@receiver c
//	@param tx 交易
//	@param chainRid 链资源id
//	@return string 交易证明
func (c *ChainClient) GetTxProve(tx *bcostypes.TransactionDetail, chainRid string) string {
	var err error
	txProve := make(map[string][]byte)
	if txProve["tx_byte"], err = json.Marshal(tx); err != nil {
		return emptyJson
	}
	txProve["tx_hash"] = []byte(tx.Hash)
	txProve["chain_rid"] = []byte(chainRid)

	txProveByte, err := json.Marshal(txProve)
	if err != nil {
		return emptyJson
	}
	c.log.Debugf("[GetTxProve] %s", string(txProveByte))
	txProveReq := cross_chain.TxVerifyRequest{
		Version: tcipcommon.Version_V1_0_0,
		TxProve: string(txProveByte),
	}
	res, err := json.Marshal(txProveReq)
	if err != nil {
		return emptyJson
	}
	return string(res)
}

// CheckChain 检查链的连通性
//
//	@receiver c
//	@return bool
func (c *ChainClient) CheckChain() bool {
	count := 0
	//for _, client := range c.client {
	//	var ctx context.Context
	//	if _, err := client.GetBlockNumber(ctx); err != nil {
	//		return false
	//	}
	//	count++
	//}
	for _ = range c.client {
		if _, err := c.GetBlockNumber(); err != nil {
			return false
		}
		count++
	}
	return count != 0
}

// TxProve 交易认证
//
//	@receiver c
//	@param txProve
//	@return bool
func (c *ChainClient) TxProve(txProve string) bool {
	c.log.Debugf("txProve: %s\n", txProve)
	txProveMap := make(map[string][]byte)
	err := json.Unmarshal([]byte(txProve), &txProveMap)
	if err != nil {
		c.log.Errorf("[TxProve] Unmarshal error: %s", err.Error())
		return false
	}
	//chainRid, ok := txProveMap["chain_rid"]
	_, ok := txProveMap["chain_rid"]
	if !ok {
		c.log.Errorf("[TxProve] chain_id not found: %s", err.Error())
		return false
	}
	txHash, ok := txProveMap["tx_hash"]
	if !ok {
		c.log.Errorf("[TxProve] tx_hash not found: %s", err.Error())
		return false
	}
	txByteString, ok := txProveMap["tx_byte"]
	if !ok {
		c.log.Errorf("[TxProve] tx_byte not found: %s", err.Error())
		return false
	}

	//client, err := c.getChainClient(string(chainRid))
	//if err != nil {
	//	c.log.Errorf("[TxProve] get client error %s", err.Error())
	//	return false
	//}
	//var ctx context.Context
	//tx, err := client.GetTransactionByHash(ctx, bcoscommon.HexToHash(string(txHash)))
	tx, err := c.GetTransactionByHash(string(txHash))
	if err != nil {
		c.log.Errorf("[TxProve] get tx error %s", err.Error())
		return false
	}
	txChainByte, err := json.Marshal(tx)
	if err != nil {
		c.log.Errorf("[TxProve] Marshal tx error %s", err.Error())
		return false
	}
	if string(txByteString) != "" && string(txByteString) == string(txChainByte) {
		return true
	}
	c.log.Errorf("[TxProve] Compare tx error\n%s\n%s\n", string(txByteString), string(txChainByte))
	return false
}

// getChainClient 获取链客户端
//
//	@receiver c
//	@param chainRid
//	@return *sdk.Client
//	@return error
func (c *ChainClient) getChainClient(chainRid string) (bool, error) {
	if _, ok := c.client[chainRid]; !ok {
		msg := fmt.Sprintf("[getChainClient] no chain client: chainRid %s", chainRid)
		c.log.Warnf(msg)
		return false, fmt.Errorf(msg)
	}
	return c.client[chainRid], nil
}

// getListenKey 拼接监听缓存的key
//
//	@param chainRid
//	@param eventName
//	@return string
func getListenKey(chainRid, contractName, eventName string) string {
	return fmt.Sprintf("%s#%s#%s", chainRid, contractName, eventName)
}

// dealParam 处理合约调用参数
//
//	@param args
//	@param paramType
//	@return []interface{}
//	@return error
func dealParam(args string, paramType []tcipcommon.EventDataType) ([]interface{}, error) {
	argsArr := make([]interface{}, 0)
	err := json.Unmarshal([]byte(args), &argsArr)
	if err != nil {
		return nil, fmt.Errorf("Umarshal args error: %s", err.Error())
	}

	if len(argsArr) != len(paramType) {
		return nil, errors.New("len(argsArr) != len(paramType)")
	}

	for i, arg := range argsArr {
		switch paramType[i] {
		case tcipcommon.EventDataType_ADDRESS:
			argsArr[i] = bcoscommon.HexToAddress(strings.ToLower(arg.(string)))
		case tcipcommon.EventDataType_HASH:
			argsArr[i] = bcoscommon.HexToHash(arg.(string))
		case tcipcommon.EventDataType_INT:
			tmpNum, err2 := strconv.Atoi(arg.(string))
			if err2 != nil {
				return nil, fmt.Errorf("Int parse error: data %s error %s ",
					argsArr[i].(string), err2.Error())
			}
			argsArr[i] = big.NewInt(int64(tmpNum))
		case tcipcommon.EventDataType_FLOAT:
			tmpNum, err2 := strconv.ParseFloat(arg.(string), 64)
			if err2 != nil {
				return nil, fmt.Errorf("Float parse error: data %s error %s ",
					argsArr[i].(string), err2.Error())
			}
			argsArr[i] = tmpNum
		case tcipcommon.EventDataType_BOOL:
			tmpBool, err2 := strconv.ParseBool(arg.(string))
			if err2 != nil {
				return nil, fmt.Errorf("Bool parse error: data %s error %s ",
					argsArr[i].(string), err2.Error())
			}
			argsArr[i] = tmpBool
		case tcipcommon.EventDataType_MAP:
			data := make(map[string]interface{})
			err2 := json.Unmarshal([]byte(arg.(string)), &data)
			if err2 != nil {
				return nil, fmt.Errorf("Map parse error, only support "+
					"map[string]interface{}: data %s error %s ",
					argsArr[i].(string), err2.Error())
			}
			argsArr[i] = data
		case tcipcommon.EventDataType_STRING:
			argsArr[i], _ = arg.(string)
		case tcipcommon.EventDataType_ARRAY:
			data := make([]interface{}, 0)
			err2 := json.Unmarshal([]byte(arg.(string)), &data)
			if err2 != nil {
				return nil, fmt.Errorf("Array parse error: data %s error %s ",
					argsArr[i].(string), err2.Error())
			}
			argsArr[i] = data
		case tcipcommon.EventDataType_BYTE:
			argsArr[i] = []byte(arg.(string))
		}
	}
	return argsArr, nil
}

func (c *ChainClient) SubscribeEventLogs(handler func(int, []bcostypes.Log)) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:5007", Path: "/chain-proxy/event"}
	log.Printf("connect to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("fail to connect:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			log.Printf("receive msg: %s", message)
			var logs *[]bcostypes.Log
			var logs2 *[]TestLog
			// 另一种方式：如果 message 包含字符数据
			strMessage := string(message) // 这将把 []uint 转换为字符串
			//var items []map[string]interface{}
			//if err := json.Unmarshal([]byte(message), &items); err != nil {
			//	log.Fatalf("JSON解析失败: %v", err)
			//}
			//for _, item := range items {
			//	// 获取指定字段值
			//	logData := item["data"]
			//	hex.DecodeString(logData)
			//}
			//hex.DecodeString(testStr)
			if err := json.NewDecoder(strings.NewReader(strMessage)).Decode(&logs); err != nil {
				c.log.Warn("[InvokeContract] failed to decode HTTP response as receipt: %s", err.Error())
			}
			if err := json.NewDecoder(strings.NewReader(strMessage)).Decode(&logs2); err != nil {
				c.log.Errorf("[InvokeContract] failed to decode HTTP response as receipt: %s", err.Error())
			}

			dereferenced := *logs2
			logsPtr := *logs // 解引用logs指针

			for i := 0; i < len(dereferenced); i++ {
				log := dereferenced[i]
				// 假设Data是16进制字符串
				val, _ := hex.DecodeString(log.Data[2:])

				// 直接操作切片
				logsPtr[i].Data = val
			}

			// 更新logs指针指向的切片
			*logs = logsPtr
			handler(0, *logs)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return err
		case t := <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write error:", err)
				return err
			}
		case <-interrupt:
			log.Println("receive interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("close error:", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

type TestLog struct {
	// Consensus fields:
	// address of the contract that generated the event
	Address bcoscommon.Address `json:"address" gencodec:"required"`
	// list of topics provided by the contract.
	Topics []bcoscommon.Hash `json:"topics" gencodec:"required"`
	// supplied by the contract, usually ABI-encoded
	Data string `json:"data" gencodec:"required"`

	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	BlockNumber uint64 `json:"blockNumber"`
	// hash of the transaction
	TxHash bcoscommon.Hash `json:"transactionHash" gencodec:"required"`
	// index of the transaction in the block
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
	// hash of the block in which the transaction was included
	BlockHash bcoscommon.Hash `json:"blockHash"`
	// index of the log in the block
	Index uint `json:"logIndex" gencodec:"required"`

	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	Removed bool `json:"removed"`
}
