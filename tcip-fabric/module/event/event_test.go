/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"github.com/gogo/protobuf/proto"

	"github.com/stretchr/testify/assert"

	chain_config "chainmaker.org/chainmaker/tcip-fabric/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/db"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-fabric/v2/module/utils"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

var (
	log         []*logger.LogModuleConfig
	chainConfig *common.FabricConfig
	event       *common.CrossChainEvent
	eventInfo   *EventInfo
)

const (
	tmp = "%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s"
)

func initTest() {
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	chainConfig = &common.FabricConfig{
		ChainRid: "chain001",
		ChainId:  "mychannel",
		TlsKey:   "sdkKeyText",
		TlsCert:  "sdkCrtText",
		Org: []*common.Org{
			{
				OrgId:    "Org1",
				MspId:    "Org1Msp",
				SignCert: "sign cert",
				SignKey:  "sign key",
				Peers: []*common.Peer{
					{
						NodeAddr:    "grpcs://127.0.0.1:7051",
						TrustRoot:   []string{"peer ca"},
						TlsHostName: "peer0.org1.example.com",
					},
				},
			},
		},
		Orderers: []*common.Orderer{
			{
				NodeAddr:    "grpcs://127.0.0.1:7051",
				TrustRoot:   []string{"orderer ca"},
				TlsHostName: "orderer.example.com",
			},
		},
	}

	event = &common.CrossChainEvent{
		CrossChainEventId:        "00001",
		ChainRid:                 "chain001",
		ContractName:             "contract1",
		CrossChainName:           "test",
		CrossChainFlag:           "test",
		TriggerCondition:         1,
		EventName:                "test",
		Timeout:                  1000,
		CrossType:                1,
		FabricEventDataType:      common.EventDataType_ARRAY,
		FabricArrayEventDataType: []common.EventDataType{common.EventDataType_STRING, common.EventDataType_STRING, common.EventDataType_STRING},
		IsCrossChain:             "zero == \"10\"",
		ConfirmInfo: &common.ConfirmInfo{
			ChainRid:     "chain001",
			ContractName: "ADDRESS2",
			Method:       "transfer",
			Parameter:    "[\"%s\"]",
			ParamData:    []int32{0},
		},
		CancelInfo: &common.CancelInfo{
			ChainRid:     "chain001",
			ContractName: "ADDRESS2",
			Method:       "burn",
			Parameter:    "[\"%s\"]",
			ParamData:    []int32{0},
		},
		CrossChainCreate: []*common.CrossChainMsg{
			{
				GatewayId:    "1",
				ChainRid:     "chain002",
				ContractName: "ADDRESS2",
				Method:       "minter",
				Parameter:    "[\"%s\"]",
				ParamData:    []int32{0},
				ConfirmInfo: &common.ConfirmInfo{
					ChainRid:     "chain002",
					ContractName: "ADDRESS2",
					Method:       "transfer",
					Parameter:    "[\"%s\"]",
					ParamData:    []int32{0},
				},
				CancelInfo: &common.CancelInfo{
					ChainRid:     "chain002",
					ContractName: "ADDRESS2",
					Method:       "burn",
					Parameter:    "[\"%s\"]",
					ParamData:    []int32{0},
				},
			},
		},
	}
	conf.Config.BaseConfig = &conf.BaseConfig{
		GatewayID:    "0",
		GatewayName:  "test",
		Address:      "https://127.0.0.1:19999",
		ServerName:   "chainmaker.org",
		Tlsca:        "../../config/cert/client/ca.crt",
		ClientKey:    "../../config/cert/client/client.key",
		ClientCert:   "../../config/cert/client/client.crt",
		TxVerifyType: "notneed",
		CallType:     "grpc",
	}
	eventInfo = &EventInfo{
		Topic:        "test",
		ChainRid:     "chain001",
		ContractName: "contract1",
		TxProve:      "{}",
		Data:         []string{"[\"10\",\"11\",\"13\"]"},
		Tx:           []byte("qweqrt"),
		TxId:         "hfjdkhsalufksa",
		BlockHeight:  10,
	}
	conf.Config.DbPath = path.Join(os.TempDir(), time.Now().String())
	logger.InitLogConfig(log)
	db.NewDbHandle()
	utils.EventChan = make(chan *utils.EventOperate)
	utils.UpdateChainConfigChan = make(chan *utils.ChainConfigOperate)
	go listenEventChan(utils.EventChan)
	go listenConfigChan(utils.UpdateChainConfigChan)
	chain_config.NewChainConfig()
}

func listenEventChan(updateChan chan *utils.EventOperate) {
	for {
		<-updateChan
	}
}

func listenConfigChan(updateChan chan *utils.ChainConfigOperate) {
	for {
		<-updateChan
	}
}

func tetsCancel(t *testing.T) {
	event.CrossChainCreate[0].CancelInfo.Parameter = tmp
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].CancelInfo.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func testConfirm(t *testing.T) {
	event.CrossChainCreate[0].ConfirmInfo.Parameter = tmp
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ConfirmInfo.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func testCrossChainCreate(t *testing.T) {
	event.CrossChainCreate[0].Parameter = tmp
	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamDataType = make([]common.EventDataType, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamData = make([]int32, 20)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamDataType = make([]common.EventDataType, 15)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ParamData = make([]int32, 10)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].Method = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].ChainRid = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate[0].GatewayId = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func Test_SaveEvent(t *testing.T) {
	initTest()
	InitEventManager()

	err := EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.NotNil(t, err)

	tetsCancel(t)
	testConfirm(t)

	testCrossChainCreate(t)

	event.CancelInfo.ChainRid = "0099"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CancelInfo = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ConfirmInfo.ChainRid = "0099"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ConfirmInfo = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.IsCrossChain = "abcdqwetrtrewq=cdwer"
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainCreate = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.IsCrossChain = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.FabricArrayEventDataType = nil
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.FabricEventDataType = common.EventDataType(100)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.TriggerCondition = common.TriggerCondition_COMPLETELY_CONTRACT_EVENT
	err = EventManagerV1.SaveEvent(event, false)
	assert.Nil(t, err)

	event.TriggerCondition = common.TriggerCondition(10)
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.CrossChainEventId = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ContractName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.ChainRid = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)

	event.EventName = ""
	err = EventManagerV1.SaveEvent(event, false)
	assert.NotNil(t, err)
}

func Test_DeleteEvent(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	err := EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	err = EventManagerV1.DeleteEvent(event)
	assert.Nil(t, err)

	event.CrossChainEventId = ""
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	event.CrossChainEventId = "123"
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)

	db.Db.Close()
	err = EventManagerV1.DeleteEvent(event)
	assert.NotNil(t, err)
}

func Test_GetEvent(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	res, err := EventManagerV1.GetEvent("1234")
	assert.NotNil(t, err)
	assert.Nil(t, res)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	res, err = EventManagerV1.GetEvent(event.CrossChainEventId)
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	res, err = EventManagerV1.GetEvent("")
	assert.Nil(t, err)
	assert.Equal(t, len(res), 1)

	err = EventManagerV1.DeleteEvent(event)
	assert.Nil(t, err)

	res, err = EventManagerV1.GetEvent(event.CrossChainEventId)
	assert.NotNil(t, err)
	assert.Equal(t, len(res), 0)
}

func Test_BuildCrossChainMsg(t *testing.T) {
	initTest()
	InitEventManager()

	_ = chain_config.ChainConfigManager.Save(chainConfig, common.Operate_SAVE)

	req, err := EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	err = EventManagerV1.SaveEvent(event, true)
	assert.Nil(t, err)

	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.NotNil(t, req)
	assert.Nil(t, err)

	event.FabricArrayEventDataType[0] = common.EventDataType(100)
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.Nil(t, err)

	event.TriggerCondition = common.TriggerCondition_COMPLETELY_CONTRACT_EVENT
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.Nil(t, req)
	assert.NotNil(t, err)

	reqByte, _ := proto.Marshal(&relay_chain.BeginCrossChainRequest{
		Version: common.Version_V1_0_0,
	})
	eventInfo.Data[0] = string(reqByte)
	req, err = EventManagerV1.BuildCrossChainMsg(eventInfo)
	assert.NotNil(t, req)
	assert.Nil(t, err)
}
