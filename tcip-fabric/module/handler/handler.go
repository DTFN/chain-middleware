/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	// "os"
	"strings"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/utils"

	//tbis_event "chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"

	chain_config "chainmaker.org/chainmaker/tcip-fabric/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/event"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"

	chain_client "chainmaker.org/chainmaker/tcip-fabric/v2/module/chain-client"

	"google.golang.org/grpc/peer"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/logger"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler 结构体
type Handler struct {
	log *zap.SugaredLogger
}

// NewHandler 初始化handler模块
//  @return *Handler
func NewHandler() *Handler {
	return &Handler{
		log: logger.GetLogger(logger.ModuleHandler),
	}
}

type DIDUpdateMsg struct {
	DID       string      `json:"did"`
	DIDDocument       string      `json:"didDocument"`
}

// CrossChainTry 接收跨链请求的接口
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.CrossChainTryResponse
//  @return error
func (h *Handler) CrossChainTry(ctx context.Context,
	req *cross_chain.CrossChainTryRequest) (*cross_chain.CrossChainTryResponse, error) {
	h.printRequest(ctx, "CrossChainTry", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:

		// // 构造参数
		// paramt := []string{"{\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"DID_MANAGER\",\"contract_func\":\"Erc20MintVcs\",\"contract_type\":\"2\",\"func_params\":{\"amount\":\"10\",\"to\":\"did:1\"},\"gateway_id\":\"2\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":null,\"contentHash\":null,\"created\":\"2025-09-14T16:06:44.694686803+00:00\",\"jws\":null,\"proofPurpose\":\"assertionMethod\",\"r\":null,\"s\":null,\"signature\":\"b9f05613a1d5a5c100d734ff353a022f08261011200ba14fbb6ed50a8e6bdd5421cffcca0e16997bf99036fab2b0fd9e74712e77488c6fdbbe225384664f6802\",\"type\":\"JsonWebSignature2020\",\"v\":null,\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"]}}"}
		// paramBytet := make([][]byte, len(paramt))
		// for i, v := range paramt {
		// 	paramBytet[i] = []byte(v)
		// }

		// // todo 删除
		// tryReultt, txIdt, txt, _ := chain_client.ChainClientV1.InvokeContract("chain_fabric_01", "DID_MANAGER", "Erc20MintVcs", paramBytet, true)
		// h.log.Debugf("[CrossChainTry] test call vc contract, txId: %s, tryReult: %+v, tx: %+v", txIdt, tryReultt, txt)

		if req.CrossChainId == "$relayer" {
			// 发起方为中继链
			h.log.Debugf("[CrossChainTry] from relayer chain")
			// 遍历所有链，调用指定合约
			chainRids := allChainRids();
			h.log.Debugf("[CrossChainTry] all chainRids: %+v, parameter: %s", chainRids, req.CrossChainMsg.Parameter)
			for _, chainRid := range chainRids {
				// 合约参数序列化为json字符串
				var didUpdateMsg DIDUpdateMsg
				err := json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &didUpdateMsg)
				if err != nil {
					h.log.Errorf("[CrossChainTry] Unmarshal fail: %s", req.CrossChainMsg.Parameter)
					continue
				}

				// parameter {data:$VC}
				var dataVc map[string]interface{}
				json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &dataVc)

				// 构造参数
				param := []string{didUpdateMsg.DID, didUpdateMsg.DIDDocument}
				paramByte := make([][]byte, len(param))
				for i, v := range param {
					paramByte[i] = []byte(v)
				}

				// 发起交易
				tryReult, txId, tx, err := chain_client.ChainClientV1.InvokeContract(chainRid, "DID_MANAGER", "UpdateDID", paramByte, true)
				h.log.Debugf("[CrossChainTry] call did contract, txId: %s, tryReult: %+v, tx: %+v", txId, tryReult, tx)
			}

			// 调用成功,返回执行失败
			return getCrossChainTryReturn(common.Code_CONTRACT_FAIL,
				req.CrossChainId, req.CrossChainName,
				req.CrossChainFlag, utils.UnsupportVersion(req.Version), nil, nil)
		} else if (req.CrossChainMsg.ExtraData == "$VC_CROSS") {
			// VC跨链交易的第二阶段
			h.log.Debugf("[CrossChainTry] VC_CROSS")

			// parameter {data:$VC}
			var dataVc map[string]interface{}
			json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &dataVc)

			// 构造参数
			param := []string{dataVc["data"].(string)}
			paramByte := make([][]byte, len(param))
			for i, v := range param {
				paramByte[i] = []byte(v)
			}

			// 发起交易
			h.log.Debugf("[CrossChainTry] call vc params, ChainRid: %s, ContractName: %s, Method: %s, param: %+v", req.CrossChainMsg.ChainRid, req.CrossChainMsg.ContractName, req.CrossChainMsg.Method, param)
			tryReult, txId, tx, err := chain_client.ChainClientV1.InvokeContract(req.CrossChainMsg.ChainRid, req.CrossChainMsg.ContractName, req.CrossChainMsg.Method, paramByte, true)
			// tryReult, txId, tx, err := chain_client.ChainClientV1.InvokeContract(req.CrossChainMsg.ChainRid, "mychannel", "Echo", paramByte, true)
			h.log.Debugf("[CrossChainTry] call vc contract, txId: %s, tryReult: %+v, tx: %+v", txId, tryReult, tx)
			if (err != nil) {
				h.log.Errorf("[CrossChainTry] call vc contract fail: %+v", err)
			}

			return getCrossChainTryReturn(common.Code_GATEWAY_SUCCESS,
				req.CrossChainId, req.CrossChainName, req.CrossChainFlag,
				common.Code_GATEWAY_SUCCESS.String(), &common.TxContent{
					TxId:      txId,
					Tx:        tx.TransactionEnvelope.Payload,
					TxResult:  common.TxResultValue_TX_SUCCESS,
					GatewayId: conf.Config.BaseConfig.GatewayID,
					ChainRid:  req.CrossChainMsg.ChainRid,
					TxProve:   chain_client.ChainClientV1.GetTxProve(tx, txId, req.CrossChainMsg.ChainRid),
					// 不走spv不需要该字段
					//BlockHeight: ,
				}, 
				[]string{tryReult},
			)
		}

		chainExist := checkChain(req.CrossChainMsg.ChainRid)
		if !chainExist {
			return getCrossChainTryReturn(common.Code_INVALID_PARAMETER,
				req.CrossChainId, req.CrossChainName,
				req.CrossChainFlag, fmt.Sprintf("%s not exist", req.CrossChainMsg.ChainRid), nil, nil)
		}
		param := make([]string, 0)
		err := json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &param)
		if err != nil {
			h.log.Errorf("[CrossChainTry] Failed to Unmarshal Parameter: cross chain id: %s, error: %s",
				req.CrossChainId, err.Error())
			return getCrossChainTryReturn(common.Code_INTERNAL_ERROR,
				req.CrossChainId, req.CrossChainName,
				req.CrossChainFlag, common.Code_INTERNAL_ERROR.String(), nil, nil)
		}
		paramByte := make([][]byte, len(param))
		for i, v := range param {
			paramByte[i] = []byte(v)
		}
		tryReult, txId, tx, err := chain_client.ChainClientV1.InvokeContract(
			req.CrossChainMsg.ChainRid, req.CrossChainMsg.ContractName,
			req.CrossChainMsg.Method, paramByte, true)
		if err != nil {
			h.log.Errorf("[CrossChainTry] Failed to execute cross-chain transaction: cross chain id: %s",
				req.CrossChainId)
			return getCrossChainTryReturn(common.Code_INTERNAL_ERROR,
				req.CrossChainId, req.CrossChainName,
				req.CrossChainFlag, err.Error(), nil, nil)
		}
		return getCrossChainTryReturn(common.Code_GATEWAY_SUCCESS,
			req.CrossChainId, req.CrossChainName, req.CrossChainFlag,
			common.Code_GATEWAY_SUCCESS.String(), &common.TxContent{
				TxId:      txId,
				Tx:        tx.TransactionEnvelope.Payload,
				TxResult:  common.TxResultValue_TX_SUCCESS,
				GatewayId: conf.Config.BaseConfig.GatewayID,
				ChainRid:  req.CrossChainMsg.ChainRid,
				TxProve:   chain_client.ChainClientV1.GetTxProve(tx, txId, req.CrossChainMsg.ChainRid),
				// 不走spv不需要该字段
				//BlockHeight: ,
			}, []string{tryReult})
	default:
		return getCrossChainTryReturn(common.Code_INVALID_PARAMETER,
			req.CrossChainId, req.CrossChainName,
			req.CrossChainFlag, utils.UnsupportVersion(req.Version), nil, nil)
	}
}

// allChain 获取所有chainRid
//  @return bool
func allChainRids() []string {
	chainRids := chain_config.ChainConfigManager.AllChainRids()
	return chainRids
}

// CrossChainConfirm 跨链结果确认
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.CrossChainConfirmResponse
//  @return error
func (h *Handler) CrossChainConfirm(ctx context.Context,
	req *cross_chain.CrossChainConfirmRequest) (*cross_chain.CrossChainConfirmResponse, error) {
	h.printRequest(ctx, "CrossChainConfirm", fmt.Sprintf("%+v", req))
	// 根据业务做一些处理，这里模拟调用一个confirm方法，这样就把业务的逻辑放在了合约中，减少跨链网关的定制化开发
	switch req.Version {
	case common.Version_V1_0_0:
		if req.ConfirmInfo == nil || req.ConfirmInfo.ChainRid == "" {
			return &cross_chain.CrossChainConfirmResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		}
		var (
			param [][]byte
			err   error
			err1  error
		)
		if req.CrossChainFlag == "tbis_event" {
			//param, err1 = fillTbisResult(req.ConfirmInfo.Parameter, req.ConfirmInfo.ChainRid,
			//	tbis_event.SubSuccess, tbis_event.SubSuccess, req.TryResult[0])
			//if err1 != nil {
			//	h.log.Errorf("[CrossChainConfirm] %s", err1.Error())
			//	return &cross_chain.CrossChainConfirmResponse{
			//		Code:    common.Code_INTERNAL_ERROR,
			//		Message: err1.Error(),
			//	}, nil
			//}
			h.log.Errorf("not support")
		} else {
			param, err1 = fillTryResult(req.ConfirmInfo.Parameter, req.TryResult, req.CrossType)
			if err1 != nil {
				h.log.Errorf("[CrossChainConfirm] %s", err1.Error())
				return &cross_chain.CrossChainConfirmResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err1.Error(),
				}, nil
			}
		}
		_, txId, tx, err := chain_client.ChainClientV1.InvokeContract(req.ConfirmInfo.ChainRid,
			req.ConfirmInfo.ContractName, req.ConfirmInfo.Method, param, true)
		if err != nil {
			h.log.Errorf("[CrossChainConfirm] Failed to execute cross-chain transaction: cross chain id: %s", req.CrossChainId)
			return &cross_chain.CrossChainConfirmResponse{
				Code:    common.Code_INTERNAL_ERROR,
				Message: err.Error(),
			}, nil
		}
		return &cross_chain.CrossChainConfirmResponse{
			Code:    common.Code_GATEWAY_SUCCESS,
			Message: common.Code_GATEWAY_SUCCESS.String(),
			TxContent: &common.TxContent{
				TxId:      txId,
				Tx:        tx.TransactionEnvelope.Payload,
				TxResult:  common.TxResultValue_TX_SUCCESS,
				GatewayId: conf.Config.BaseConfig.GatewayID,
				ChainRid:  req.ConfirmInfo.ChainRid,
				// 这里不验证不需要填
				TxProve: "",
				// 不走spv不需要该字段
				//BlockHeight: ,
			},
		}, nil
	default:
		return &cross_chain.CrossChainConfirmResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// CrossChainCancel 跨链结果确认
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.CrossChainCancelResponse
//  @return error
func (h *Handler) CrossChainCancel(ctx context.Context,
	req *cross_chain.CrossChainCancelRequest) (*cross_chain.CrossChainCancelResponse, error) {
	h.printRequest(ctx, "CrossChainCancel", fmt.Sprintf("%+v", req))
	// 根据业务做一些处理,这里模拟调用一个合约中的cancel方法，这样就把业务的逻辑放在了合约中，减少跨链网关的定制化开发
	switch req.Version {
	case common.Version_V1_0_0:
		if req.CancelInfo == nil || req.CancelInfo.ChainRid == "" {
			return &cross_chain.CrossChainCancelResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		}
		paramByte := make([][]byte, 0)
		if req.CrossChainFlag == "tbis_event.TbisFlag" {
			//var err1 error
			//paramByte, err1 = fillTbisResult(req.CancelInfo.Parameter, req.CancelInfo.ChainRid,
			//	tbis_event.SubSuccess, tbis_event.SubSuccess, "failed")
			//if err1 != nil {
			//	h.log.Errorf("[CrossChainConfirm] %s", err1.Error())
			//	return &cross_chain.CrossChainCancelResponse{
			//		Code:    common.Code_INTERNAL_ERROR,
			//		Message: err1.Error(),
			//	}, nil
			//}
			h.log.Errorf("not support")
		} else {
			param := make([]string, 0)
			if req.CancelInfo.Parameter != "" {
				err1 := json.Unmarshal([]byte(req.CancelInfo.Parameter), &param)
				if err1 != nil {
					return &cross_chain.CrossChainCancelResponse{
						Code:    common.Code_INTERNAL_ERROR,
						Message: "unmarshal param error: " + err1.Error(),
					}, nil
				}
			}
			for _, v := range param {
				paramByte = append(paramByte, []byte(v))
			}
		}
		_, txId, tx, err := chain_client.ChainClientV1.InvokeContract(req.CancelInfo.ChainRid,
			req.CancelInfo.ContractName, req.CancelInfo.Method, paramByte, true)
		if err != nil {
			h.log.Errorf("[CrossChainCancel] Failed to execute cross-chain transaction: cross chain id: %s", req.CrossChainId)
			return &cross_chain.CrossChainCancelResponse{
				Code:    common.Code_INTERNAL_ERROR,
				Message: err.Error(),
			}, nil
		}
		return &cross_chain.CrossChainCancelResponse{
			Code:    common.Code_GATEWAY_SUCCESS,
			Message: common.Code_GATEWAY_SUCCESS.String(),
			TxContent: &common.TxContent{
				TxId:      txId,
				Tx:        tx.TransactionEnvelope.Payload,
				TxResult:  common.TxResultValue_TX_SUCCESS,
				GatewayId: conf.Config.BaseConfig.GatewayID,
				ChainRid:  req.CancelInfo.ChainRid,
				// 这里不验证不需要填
				TxProve: "",
				// 不走spv不需要该字段
				//BlockHeight: ,
			},
		}, nil
	default:
		return &cross_chain.CrossChainCancelResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// IsCrossChainSuccess 判断跨链结果
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.IsCrossChainSuccessResponse
//  @return error
func (h *Handler) IsCrossChainSuccess(ctx context.Context,
	req *cross_chain.IsCrossChainSuccessRequest) (*cross_chain.IsCrossChainSuccessResponse, error) {
	h.printRequest(ctx, "IsCrossChainSuccess", fmt.Sprintf("%+v", req))
	// 根据业务做一些处理，这里一律让他失败
	return &cross_chain.IsCrossChainSuccessResponse{
		CrossChainResult: false,
		Code:             common.Code_GATEWAY_SUCCESS,
		Message:          common.Code_GATEWAY_SUCCESS.String(),
	}, nil
}

// CrossChainEvent 跨链触发事件管理
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.CrossChainEventResponse
//  @return error
func (h *Handler) CrossChainEvent(ctx context.Context,
	req *cross_chain.CrossChainEventRequest) (*cross_chain.CrossChainEventResponse, error) {
	h.printRequest(ctx, "CrossChainEvent", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		switch req.Operate {
		case common.Operate_GET:
			crossChainEvent, err := event.EventManagerV1.GetEvent(string(event.EventKey(req.CrossChainEvent.EventName,
				req.CrossChainEvent.ContractName, req.CrossChainEvent.ChainRid)))
			if err != nil {
				return &cross_chain.CrossChainEventResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.CrossChainEventResponse{
				Code:                common.Code_GATEWAY_SUCCESS,
				Message:             common.Code_GATEWAY_SUCCESS.String(),
				CrossChainEventList: crossChainEvent,
			}, nil
		case common.Operate_DELETE:
			err := event.EventManagerV1.DeleteEvent(req.CrossChainEvent)
			if err != nil {
				return &cross_chain.CrossChainEventResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.CrossChainEventResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		case common.Operate_SAVE:
			err := event.EventManagerV1.SaveEvent(req.CrossChainEvent, true)
			if err != nil {
				return &cross_chain.CrossChainEventResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.CrossChainEventResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		case common.Operate_UPDATE:
			err := event.EventManagerV1.SaveEvent(req.CrossChainEvent, false)
			if err != nil {
				return &cross_chain.CrossChainEventResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.CrossChainEventResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		default:
			return &cross_chain.CrossChainEventResponse{
				Code:    common.Code_INVALID_PARAMETER,
				Message: "unsupported operate",
			}, nil
		}
	default:
		return &cross_chain.CrossChainEventResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// TxVerify rpc交易验证，不是非要在当前服务中实现
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.TxVerifyResponse
//  @return error
func (h *Handler) TxVerify(ctx context.Context,
	req *cross_chain.TxVerifyRequest) (*cross_chain.TxVerifyResponse, error) {
	h.printRequest(ctx, "TxVerify", fmt.Sprintf("%+v", req))
	switch req.Version {
	case common.Version_V1_0_0:
		res := chain_client.ChainClientV1.TxProve(req.TxProve)
		return &cross_chain.TxVerifyResponse{
			TxVerifyResult: res,
			Code:           common.Code_GATEWAY_SUCCESS,
			Message:        common.Code_GATEWAY_SUCCESS.String(),
		}, nil
	default:
		return &cross_chain.TxVerifyResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// ChainIdentity 链配置
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.ChainIdentityResponse
//  @return error
func (h *Handler) ChainIdentity(ctx context.Context,
	req *cross_chain.ChainIdentityRequest) (*cross_chain.ChainIdentityResponse, error) {
	h.printRequest(ctx, "ChainIdentity", fmt.Sprintf("%+v", req))

	if req.FabricConfig == nil {
		return &cross_chain.ChainIdentityResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: "FabricConfig is required",
		}, nil
	}
	switch req.Version {
	case common.Version_V1_0_0:
		switch req.Operate {
		case common.Operate_GET:
			fabricConfig, err := chain_config.ChainConfigManager.Get(req.FabricConfig.ChainRid)
			if err != nil {
				return &cross_chain.ChainIdentityResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.ChainIdentityResponse{
				Code:         common.Code_GATEWAY_SUCCESS,
				Message:      common.Code_GATEWAY_SUCCESS.String(),
				FabricConfig: fabricConfig,
			}, nil
		case common.Operate_DELETE:
			err := chain_config.ChainConfigManager.Delete(req.FabricConfig.ChainRid)
			if err != nil {
				return &cross_chain.ChainIdentityResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.ChainIdentityResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		case common.Operate_SAVE:
			err := chain_config.ChainConfigManager.Save(req.FabricConfig, common.Operate_SAVE)
			if err != nil {
				return &cross_chain.ChainIdentityResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.ChainIdentityResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		case common.Operate_UPDATE:
			err := chain_config.ChainConfigManager.Save(req.FabricConfig, common.Operate_UPDATE)
			if err != nil {
				return &cross_chain.ChainIdentityResponse{
					Code:    common.Code_INTERNAL_ERROR,
					Message: err.Error(),
				}, nil
			}
			return &cross_chain.ChainIdentityResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			}, nil
		default:
			return &cross_chain.ChainIdentityResponse{
				Code:    common.Code_INVALID_PARAMETER,
				Message: "unsupported operate",
			}, nil
		}
	default:
		return &cross_chain.ChainIdentityResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// PingPong 心跳
//  @receiver h
//  @param ctx
//  @param req
//  @return *cross_chain.PingPongResponse
//  @return error
func (h *Handler) PingPong(ctx context.Context, req *emptypb.Empty) (*cross_chain.PingPongResponse, error) {
	//h.printRequest(ctx, "PingPong", fmt.Sprintf("%+v", req))
	PingPong := chain_client.ChainClientV1.CheckChain()
	h.log.Debugf("PingPong %+v", PingPong)
	return &cross_chain.PingPongResponse{
		ChainOk: PingPong,
	}, nil
}

func (h *Handler) printRequest(ctx context.Context, method, request string) {
	pr, ok := peer.FromContext(ctx)
	var addr string
	if !ok || pr.Addr == net.Addr(nil) {
		h.log.Errorf("getClientAddr FromContext failed")
		addr = "unknown"
	} else {
		addr = pr.Addr.String()
	}

	h.log.Infof("[%s]: |%s|%s", method, addr, request)
}

func getCrossChainTryReturn(
	code common.Code, crossChainId, crossChainName, crossChainFlag,
	msg string, txContent *common.TxContent, tryResult []string) (*cross_chain.CrossChainTryResponse, error) {
	return &cross_chain.CrossChainTryResponse{
		CrossChainId:   crossChainId,
		CrossChainName: crossChainName,
		CrossChainFlag: crossChainFlag,
		TxContent:      txContent,
		TryResult:      tryResult,
		Code:           code,
		Message:        msg,
	}, nil
}

func checkChain(chainRid string) bool {
	res, err := chain_config.ChainConfigManager.Get(chainRid)
	if err != nil || len(res) == 0 {
		return false
	}
	return true
}

func fillTryResult(param string, tryResult []string, crossType common.CrossType) ([][]byte, error) {
	paramArr := make([]string, 0)
	if param == "" {
		return parseToByte(paramArr), nil
	}
	err := json.Unmarshal([]byte(param), &paramArr)
	if err != nil {
		return nil, fmt.Errorf("unmarshal param error: %s", err.Error())
	}
	if crossType == common.CrossType_INVOKE {
		return parseToByte(paramArr), nil
	}
	tryResultCount := strings.Count(param, common.TryResult_TRY_RESULT.String())
	if len(tryResult) == 0 || tryResultCount == 0 {
		return parseToByte(paramArr), nil
	}
	if len(tryResult) != tryResultCount {
		return nil, fmt.Errorf("\"%s\" count != len(TryResult), please update event config",
			common.TryResult_TRY_RESULT.String())
	}
	param = strings.Replace(param, common.TryResult_TRY_RESULT.String(), "%s", -1)
	paramData := make([]interface{}, len(tryResult))
	for j, v := range tryResult {
		paramData[j] = v
	}
	param = fmt.Sprintf(param, paramData...)
	err = json.Unmarshal([]byte(param), &paramArr)
	if err != nil {
		return nil, fmt.Errorf("unmarshal param error(TryResult): %s", err.Error())
	}
	return parseToByte(paramArr), nil
}

func parseToByte(params []string) [][]byte {
	paramsByte := make([][]byte, len(params))
	for i, v := range params {
		paramsByte[i] = []byte(v)
	}
	return paramsByte
}

// fillTbisResult 填充tbis执行结果
//  @param kvJsonStr
//  @param chainRid
//  @param proveStatus
//  @param contractStatus
//  @param contractResult
//  @return []byte
//  @return error
//func fillTbisResult(param, chainRid string,
//	proveStatus, contractStatus int, contractResult string) ([][]byte, error) {
//	res := tbis_event.GetCommitParam(chainRid, proveStatus, contractStatus, contractResult)
//	paramArr := make([]string, 0)
//	if param == "" {
//		return [][]byte{[]byte(res)}, nil
//	}
//	err := json.Unmarshal([]byte(param), &paramArr)
//	if err != nil {
//		return nil, fmt.Errorf("unmarshal param error: %s", err.Error())
//	}
//	paramArr[0] = res
//	return parseToByte(paramArr), nil
//}
