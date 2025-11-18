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
	"os"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/utils"

	//tbis_event "chainmaker.org/chainmaker/tcip-chainmaker/v2/module/event"

	chain_config "chainmaker.org/chainmaker/tcip-ethereum/v2/module/chain-config"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/conf"

	chain_client "chainmaker.org/chainmaker/tcip-ethereum/v2/module/chain-client"

	"google.golang.org/grpc/peer"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/event"
	"chainmaker.org/chainmaker/tcip-ethereum/v2/module/logger"
	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Handler handler结构体
type Handler struct {
	log *zap.SugaredLogger
}

const nilParam = "{}"

// NewHandler 初始化handler模块
//
//	@Description:
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
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.CrossChainTryResponse
//	@return error
func (h *Handler) CrossChainTry(ctx context.Context,
	req *cross_chain.CrossChainTryRequest) (*cross_chain.CrossChainTryResponse, error) {
	h.printRequest(ctx, "CrossChainTry", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		chainExist := checkChain(req.CrossChainMsg.ChainRid)
		if !chainExist {
			h.log.Debugf("[CrossChainTry] ChainRid[%s] not exist", req.CrossChainMsg.ChainRid)
			if req.CrossChainId == "$relayer" {
				h.log.Debugf("[CrossChainTry] from relayer chain")
				// 遍历所有链，调用指定合约
				chainRids := allChainRids();
				h.log.Debugf("[CrossChainTry] all chainRids: %+v, parameter: %s", chainRids, req.CrossChainMsg.Parameter)

				// 合约参数序列化为json字符串
				var param DIDUpdateMsg
				err := json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &param)
				if err != nil {
					h.log.Errorf("[CrossChainTry] Unmarshal fail: %s", req.CrossChainMsg.Parameter)
				}
				// 调用零数链
				didManagerAddress := os.Getenv("DID_MANAGER_ADDRESS")
				var didManagerUpdateAbi map[string]interface{}
				err2 := json.Unmarshal([]byte(DID_MANAGER_ABI_UPDATE_DID), &didManagerUpdateAbi)
				if err2 != nil {
					h.log.Errorf("[CrossChainTry] Unmarshal fail: %s", req.CrossChainMsg.Parameter)
				}
				tryResult, err := chain_client.ChainClientV1.DIDManagerUpdate(1, didManagerAddress, param.DID, param.DIDDocument)
				h.log.Debugf("[CrossChainTry] didManager contract result: %+v, err: %+v", tryResult, err)

				// 调用成功,返回执行失败
				return getCrossChainTryReturn(common.Code_CONTRACT_FAIL,
					req.CrossChainId, req.CrossChainName,
					req.CrossChainFlag, utils.UnsupportVersion(req.Version), nil, nil)
			} else {
				return getCrossChainTryReturn(common.Code_INVALID_PARAMETER,
				req.CrossChainId, req.CrossChainName, req.CrossChainFlag,
				fmt.Sprintf("%s not exist", req.CrossChainMsg.ChainRid), nil, nil)
			}
		}

		// 跨链交易的第二阶段
		if (req.CrossChainMsg.ExtraData == "$VC_CROSS") {
			// {data:[VC]}
			var parameter map[string]interface{}
			json.Unmarshal([]byte(req.CrossChainMsg.Parameter), &parameter)

			tryResult, tx, err := chain_client.ChainClientV1.InvokeContractByVc(1, req.CrossChainMsg.ContractName, req.CrossChainMsg.Method, parameter["data"].(string))
			if (err != nil) {
				h.log.Errorf("[CrossChainTry] call vc contract fail: %+v", err)
			}
			h.log.Debugf("[CrossChainTry] call vc contract tx: %+v", tx)

			txByte, _ := json.Marshal(tx)
			blockNumber := strings.Replace(tx.BlockNumber, "0x", "", -1)
			h, _ := strconv.ParseUint(blockNumber, 16, 64)
			return getCrossChainTryReturn(common.Code_GATEWAY_SUCCESS,
				req.CrossChainId, req.CrossChainName,
				req.CrossChainFlag, common.Code_GATEWAY_SUCCESS.String(), &common.TxContent{
					TxId:        tx.Hash,
					Tx:          txByte,
					TxResult:    common.TxResultValue_TX_SUCCESS,
					GatewayId:   conf.Config.BaseConfig.GatewayID,
					ChainRid:    req.CrossChainMsg.ChainRid,
					TxProve:     chain_client.ChainClientV1.GetTxProve(tx, req.CrossChainMsg.ChainRid),
					BlockHeight: int64(h),
				}, tryResult)
		}

		tryResult, tx, err := chain_client.ChainClientV1.InvokeContract(req.CrossChainMsg.ChainRid,
			req.CrossChainMsg.ContractName, req.CrossChainMsg.Method,
			req.CrossChainMsg.Abi, req.CrossChainMsg.Parameter, true, req.CrossChainMsg.ParamDataType)
		if err != nil {
			h.log.Errorf("[CrossChainTry] Failed to execute cross-chain transaction: cross chain id: %s",
				req.CrossChainId)
			return getCrossChainTryReturn(common.Code_INTERNAL_ERROR,
				req.CrossChainId, req.CrossChainName, req.CrossChainFlag,
				err.Error(), nil, nil)
		}
		txByte, _ := json.Marshal(tx)
		blockNumber := strings.Replace(tx.BlockNumber, "0x", "", -1)
		h, _ := strconv.ParseUint(blockNumber, 16, 64)
		return getCrossChainTryReturn(common.Code_GATEWAY_SUCCESS,
			req.CrossChainId, req.CrossChainName,
			req.CrossChainFlag, common.Code_GATEWAY_SUCCESS.String(), &common.TxContent{
				TxId:        tx.Hash,
				Tx:          txByte,
				TxResult:    common.TxResultValue_TX_SUCCESS,
				GatewayId:   conf.Config.BaseConfig.GatewayID,
				ChainRid:    req.CrossChainMsg.ChainRid,
				TxProve:     chain_client.ChainClientV1.GetTxProve(tx, req.CrossChainMsg.ChainRid),
				BlockHeight: int64(h),
			}, tryResult)
	default:
		return getCrossChainTryReturn(common.Code_INVALID_PARAMETER,
			req.CrossChainId, req.CrossChainName,
			req.CrossChainFlag, utils.UnsupportVersion(req.Version), nil, nil)
	}
}

const DID_MANAGER_ABI_UPDATE_DID = "{\"constant\":false,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"},{\"name\":\"newDidDocument\",\"type\":\"string\"}],\"name\":\"updateDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}"

const DID_MANAGER_ABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"},{\"name\":\"didDocument\",\"type\":\"string\"}],\"name\":\"createDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"getDIDDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"},{\"name\":\"newDidDocument\",\"type\":\"string\"}],\"name\":\"updateDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"doesDIDExist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"did\",\"type\":\"string\"}],\"name\":\"echo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"did\",\"type\":\"string\"},{\"indexed\":false,\"name\":\"didDocument\",\"type\":\"string\"}],\"name\":\"DIDModify\",\"type\":\"event\"}]"

// CrossChainConfirm 跨链结果确认
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.CrossChainConfirmResponse
//	@return error
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
			param string
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
					Message: err.Error(),
				}, nil
			}
		}
		_, tx, err := chain_client.ChainClientV1.InvokeContract(req.ConfirmInfo.ChainRid,
			req.ConfirmInfo.ContractName, req.ConfirmInfo.Method, req.ConfirmInfo.Abi,
			param, true, req.ConfirmInfo.ParamDataType)
		if err != nil {
			h.log.Errorf("[CrossChainTry] Failed to execute cross-chain transaction: cross chain id: %s",
				req.CrossChainId)
			return &cross_chain.CrossChainConfirmResponse{
				Code:    common.Code_INTERNAL_ERROR,
				Message: err.Error(),
			}, nil
		}
		txByte, _ := json.Marshal(tx)
		blockHeight, _ := strconv.Atoi(tx.BlockNumber)
		return &cross_chain.CrossChainConfirmResponse{
			Code:    common.Code_GATEWAY_SUCCESS,
			Message: common.Code_GATEWAY_SUCCESS.String(),
			TxContent: &common.TxContent{
				TxId:      tx.Hash,
				Tx:        txByte,
				TxResult:  common.TxResultValue_TX_SUCCESS,
				GatewayId: conf.Config.BaseConfig.GatewayID,
				ChainRid:  req.ConfirmInfo.ChainRid,
				// 这里不验证不需要填
				TxProve:     "",
				BlockHeight: int64(blockHeight),
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
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.CrossChainCancelResponse
//	@return error
func (h *Handler) CrossChainCancel(
	ctx context.Context, req *cross_chain.CrossChainCancelRequest) (*cross_chain.CrossChainCancelResponse, error) {
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
		param := nilParam
		if req.CancelInfo.Parameter != "" {
			param = req.CancelInfo.Parameter
		}
		//var err1 error
		if req.CrossChainFlag == "tbis_event.TbisFlag" {
			//param, err1 = fillTbisResult(req.CancelInfo.Parameter, req.CancelInfo.ChainRid,
			//	tbis_event.SubFailed, tbis_event.SubFailed, "failed")
			//if err1 != nil {
			//	h.log.Errorf("[CrossChainConfirm] %s", err1.Error())
			//	return &cross_chain.CrossChainCancelResponse{
			//		Code:    common.Code_INTERNAL_ERROR,
			//		Message: err1.Error(),
			//	}, nil
			//}
			h.log.Errorf("not support")
		}
		_, tx, err := chain_client.ChainClientV1.InvokeContract(req.CancelInfo.ChainRid,
			req.CancelInfo.ContractName, req.CancelInfo.Method,
			req.CancelInfo.Abi, param, true, req.CancelInfo.ParamDataType)
		if err != nil {
			h.log.Errorf("[CrossChainTry] Failed to execute cross-chain transaction: cross chain id: %s",
				req.CrossChainId)
			return &cross_chain.CrossChainCancelResponse{
				Code:    common.Code_INTERNAL_ERROR,
				Message: err.Error(),
			}, nil
		}
		txByte, _ := json.Marshal(tx)
		blockHeight, _ := strconv.Atoi(tx.BlockNumber)
		return &cross_chain.CrossChainCancelResponse{
			Code:    common.Code_GATEWAY_SUCCESS,
			Message: common.Code_GATEWAY_SUCCESS.String(),
			TxContent: &common.TxContent{
				TxId:      tx.Hash,
				Tx:        txByte,
				TxResult:  common.TxResultValue_TX_SUCCESS,
				GatewayId: conf.Config.BaseConfig.GatewayID,
				ChainRid:  req.CancelInfo.ChainRid,
				// 这里不验证不需要填
				TxProve:     "",
				BlockHeight: int64(blockHeight),
			},
		}, nil
	default:
		return &cross_chain.CrossChainCancelResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// CrossChainEvent 跨链触发事件管理
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.CrossChainEventResponse
//	@return error
func (h *Handler) CrossChainEvent(ctx context.Context,
	req *cross_chain.CrossChainEventRequest) (*cross_chain.CrossChainEventResponse, error) {
	h.printRequest(ctx, "CrossChainEvent", fmt.Sprintf("req: %+v", req))
	h.printRequest(ctx, "CrossChainEvent", fmt.Sprintf("req.Version: %+v", req.Version))
	h.printRequest(ctx, "CrossChainEvent", fmt.Sprintf("req.Operate: %+v", req.Operate))

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

// TxVerify rpc交易验证，不是非要在当前服务中实现,本项目不支持rpc验证，不需要实现
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.TxVerifyResponse
//	@return error
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

// IsCrossChainSuccess 判断跨链结果
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.IsCrossChainSuccessResponse
//	@return error
func (h *Handler) IsCrossChainSuccess(
	ctx context.Context,
	req *cross_chain.IsCrossChainSuccessRequest) (*cross_chain.IsCrossChainSuccessResponse, error) {
	h.printRequest(ctx, "IsCrossChainSuccess", fmt.Sprintf("%+v", req))
	// 根据业务做一些处理，这里一律让他失败
	return &cross_chain.IsCrossChainSuccessResponse{
		CrossChainResult: false,
		Code:             common.Code_GATEWAY_SUCCESS,
		Message:          common.Code_GATEWAY_SUCCESS.String(),
	}, nil
}

// ChainIdentity 链配置
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.ChainIdentityResponse
//	@return error
func (h *Handler) ChainIdentity(ctx context.Context,
	req *cross_chain.ChainIdentityRequest) (*cross_chain.ChainIdentityResponse, error) {
	h.printRequest(ctx, "ChainIdentity", fmt.Sprintf("%+v", req))

	if req.BcosConfig == nil {
		return &cross_chain.ChainIdentityResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: "EthConfig is required",
		}, nil
	}
	switch req.Version {
	case common.Version_V1_0_0:
		switch req.Operate {
		case common.Operate_GET:
			_, err := chain_config.ChainConfigManager.Get(req.BcosConfig.ChainRid)
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
		case common.Operate_DELETE:
			err := chain_config.ChainConfigManager.Delete(req.BcosConfig.ChainRid)
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
			ethConfig := conf.EthConfig{
				ChainRid: req.BcosConfig.ChainRid,
				ChainId:  req.BcosConfig.ChainId,
			}
			err := chain_config.ChainConfigManager.Save(&ethConfig, common.Operate_SAVE)
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
			ethConfig := conf.EthConfig{
				ChainRid: req.BcosConfig.ChainRid,
				ChainId:  req.BcosConfig.ChainId,
			}
			err := chain_config.ChainConfigManager.Save(&ethConfig, common.Operate_UPDATE)
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
//
//	@receiver h
//	@param ctx
//	@param req
//	@return *cross_chain.PingPongResponse
//	@return error
func (h *Handler) PingPong(ctx context.Context, req *emptypb.Empty) (*cross_chain.PingPongResponse, error) {
	//h.printRequest(ctx, "PingPong", fmt.Sprintf("%+v", req))

	return &cross_chain.PingPongResponse{
		ChainOk: chain_client.ChainClientV1.CheckChain(),
	}, nil
}

// printRequest 打印请求信息
//
//	@receiver h
//	@param ctx
//	@param method
//	@param request
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

// getCrossChainTryReturn 创建crosschaintry返回值
//
//	@param code
//	@param crossChainId
//	@param crossChainName
//	@param crossChainFlag
//	@param msg
//	@param txContent
//	@param tryResult
//	@return *cross_chain.CrossChainTryResponse
//	@return error
func getCrossChainTryReturn(
	code common.Code, crossChainId,
	crossChainName, crossChainFlag,
	msg string, txContent *common.TxContent,
	tryResult []string) (*cross_chain.CrossChainTryResponse, error) {
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

// checkChain 检查链是否存在
//
//	@param chainRid
//	@return bool
func checkChain(chainRid string) bool {
	res, err := chain_config.ChainConfigManager.Get(chainRid)
	if err != nil || len(res) == 0 {
		return false
	}
	return true
}

// allChain 获取所有chainRid
//  @return bool
func allChainRids() []string {
	chainRids := chain_config.ChainConfigManager.AllChainRids()
	return chainRids
}

// fillTryResult 填充跨链查询内容
//
//	@param param
//	@param tryResult
//	@param crossType
//	@return string
//	@return error
func fillTryResult(param string, tryResult []string, crossType common.CrossType) (string, error) {
	if param == "" {
		return nilParam, nil
	}
	if crossType == common.CrossType_INVOKE {
		return param, nil
	}
	tryResultCount := strings.Count(param, common.TryResult_TRY_RESULT.String())
	if len(tryResult) == 0 || tryResultCount == 0 {
		return param, nil
	}
	if len(tryResult) != tryResultCount {
		return "", fmt.Errorf("\"%s\" count != len(TryResult), please update event config",
			common.TryResult_TRY_RESULT.String())
	}
	param = strings.Replace(param, common.TryResult_TRY_RESULT.String(), "%s", -1)
	paramData := make([]interface{}, len(tryResult))
	for j, v := range tryResult {
		paramData[j] = v
	}
	param = fmt.Sprintf(param, paramData...)
	return param, nil
}

// fillTbisResult 填充tbis执行结果
//  @param kvJsonStr
//  @param chainRid
//  @param proveStatus
//  @param contractStatus
//  @param contractResult
//  @return string
//  @return error
//func fillTbisResult(param, chainRid string,
//	proveStatus, contractStatus int, contractResult string) (string, error) {
//	res := tbis_event.GetCommitParam(chainRid, proveStatus, contractStatus, contractResult)
//	if param == "" {
//		return fmt.Sprintf("[\"%s\"]", res), nil
//	}
//	paramArr := make([]interface{}, 0)
//	err := json.Unmarshal([]byte(param), &paramArr)
//	if err != nil {
//		return "", fmt.Errorf("unmarshal param error: %s", err.Error())
//	}
//	paramArr[0] = res
//	resStr, _ := json.Marshal(paramArr)
//	return string(resStr), nil
//}
