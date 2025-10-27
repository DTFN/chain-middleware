/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/utils"

	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/crosschaintx"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/prove"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"google.golang.org/grpc/peer"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"

	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
)

// Handler handler结构体
type Handler struct {
	log *zap.SugaredLogger
}

// NewHandler 新建handler
//  @return *Handler
func NewHandler() *Handler {
	handler := &Handler{
		log: logger.GetLogger(logger.ModuleHandler),
	}
	return handler
}

// SyncBlockHeader 同步区块头
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.SyncBlockHeaderResponse
//  @return error
func (h *Handler) SyncBlockHeader(
	ctx context.Context, req *relay_chain.SyncBlockHeaderRequest) (*relay_chain.SyncBlockHeaderResponse, error) {
	h.printRequest(ctx, "SyncBlockHeader", fmt.Sprintf("%+v", req))

	beginTime := time.Now().Unix()
	switch req.Version {
	case common.Version_V1_0_0:
		// 开始之前，需要先判断gateway是否已经注册，防止恶意gateway的连接
		g, err := gateway.GatewayV1.GetGatewayInfo(req.GatewayId)
		if err != nil || g == nil || g.GatewayId == "" {
			msg := fmt.Sprintf("no such gateway: %s", req.GatewayId)
			h.log.Error(msg)
			return &relay_chain.SyncBlockHeaderResponse{
				Code:    common.Code_INVALID_PARAMETER,
				Message: msg,
			}, nil
		}
		_, err = relay_chain_chainmaker.RelayChainV1.InvokeContract(
			utils.GetSpvContractName(req.GatewayId, req.ChainRid),
			utils.SyncBlockHeaderMethod,
			true,
			utils.GetSyncBlockHeaderParameter(req.BlockHeight, req.BlockHeader),
			-1)
		if err != nil {
			msg := fmt.Sprintf(
				"[SyncBlockHeader] call contract error: gatewayId: %s, chainRid: %s error: %s",
				req.GatewayId, req.ChainRid, err.Error())
			h.log.Error(msg)
			return &relay_chain.SyncBlockHeaderResponse{Code: common.Code_CONTRACT_FAIL, Message: msg}, nil
		}
		h.log.Debugf("[SyncBlockHeader]Finish, block height: %d, time used: %d",
			req.BlockHeight, time.Now().Unix()-beginTime)
		return &relay_chain.SyncBlockHeaderResponse{
			Code:    common.Code_GATEWAY_SUCCESS,
			Message: common.Code_GATEWAY_SUCCESS.String(),
		}, nil
	default:
		return &relay_chain.SyncBlockHeaderResponse{
			Code:    common.Code_INVALID_PARAMETER,
			Message: utils.UnsupportVersion(req.Version),
		}, nil
	}
}

// BeginCrossChain 接收跨链请求
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func (h *Handler) BeginCrossChain(
	ctx context.Context, req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error) {
	h.printRequest(ctx, "BeginCrossChain", fmt.Sprintf("%+v", req))

	beginTime := time.Now().Unix()
	switch req.Version {
	case common.Version_V1_0_0:
		// 开始之前，需要先判断gateway是否已经注册，防止恶意gateway的连接
		g, err := gateway.GatewayV1.GetGatewayInfo(req.From)
		if err != nil || g == nil || g.GatewayId == "" {
			msg := fmt.Sprintf("no such gateway: %s", req.From)
			h.log.Error(msg)
			return getBeginCrossChainReturn(common.Code_INVALID_PARAMETER, "", msg, h.log, beginTime)
		}
		// 第一步，生成跨链消息的结构
		crossChaintxInfo := crosschaintx.CrossChainTxV1.BuildCrossChainInfoFromBeginCrossChainRequest(req)
		h.log.Debugf("[BeginCrossChain] new cross chain info success: %s, time used: %d",
			crossChaintxInfo.CrossChainId, time.Now().Unix()-beginTime)
		// 第二步，验证交易是否有效
		txVerify, err := prove.ProveV1.ProveTx(
			crossChaintxInfo.FirstTxContent, crossChaintxInfo.From, req.TxContent.ChainRid)
		if err != nil {
			msg := fmt.Sprintf("[BeginCrossChain] tx verify error, cross name: %s, cross flag: %s",
				crossChaintxInfo.CrossChainName, crossChaintxInfo.CrossChainFlag)
			h.log.Error(msg)
			return getBeginCrossChainReturn(common.Code_TX_PROVE_ERROR, "", msg, h.log, beginTime)
		}
		if !txVerify {
			msg := fmt.Sprintf("[BeginCrossChain] tx verify failed, verify result: %d",
				crossChaintxInfo.FirstTxContent.TxVerifyResult)
			h.log.Error(msg)
			return getBeginCrossChainReturn(common.Code_TX_PROVE_ERROR, "", msg, h.log, beginTime)
		}
		h.log.Debugf("[BeginCrossChain] tx prove success: %s, time used: %d, bock height: %d",
			crossChaintxInfo.CrossChainId, time.Now().Unix()-beginTime,
			crossChaintxInfo.FirstTxContent.TxContent.BlockHeight)
		// 第三步，跨链信息上链
		crossChainId, err := crosschaintx.CrossChainTxV1.NewCrossChainInfo(crossChaintxInfo)
		if err != nil {
			msg := fmt.Sprintf("[BeginCrossChain] save cross chain error: cross name: %s, cross flag: %s",
				crossChaintxInfo.CrossChainName, crossChaintxInfo.CrossChainFlag)
			h.log.Error(msg)
			return getBeginCrossChainReturn(common.Code_CONTRACT_FAIL, "", msg, h.log, beginTime)
		}
		h.log.Debugf("[BeginCrossChain] save cross chain info to chain success: %s, time used: %d",
			crossChaintxInfo.CrossChainId, time.Now().Unix()-beginTime)
		return getBeginCrossChainReturn(common.Code_GATEWAY_SUCCESS, crossChainId,
			common.Code_GATEWAY_SUCCESS.String(), h.log, beginTime)
	default:
		return getBeginCrossChainReturn(common.Code_INVALID_PARAMETER,
			"", utils.UnsupportVersion(req.Version), h.log, beginTime)
	}
}

// GatewayRegister 网关注册
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.GatewayRegisterResponse
//  @return error
func (h *Handler) GatewayRegister(
	ctx context.Context, req *relay_chain.GatewayRegisterRequest) (*relay_chain.GatewayRegisterResponse, error) {
	h.printRequest(ctx, "GatewayRegister", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		req.GatewayInfo.RelayChainId = conf.BaseConf.GatewayID
		gatewayId, err := gateway.GatewayV1.GatewayRegister(req.GatewayInfo)
		if err != nil {
			return getGatewayRegisterReturn(common.Code_INTERNAL_ERROR, "", err.Error())
		}
		return getGatewayRegisterReturn(common.Code_GATEWAY_SUCCESS, gatewayId, common.Code_GATEWAY_SUCCESS.String())
	default:
		return getGatewayRegisterReturn(
			common.Code_INVALID_PARAMETER, "", utils.UnsupportVersion(req.Version))
	}
}

// GatewayUpdate 网关信息更新
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.GatewayUpdateResponse
//  @return error
func (h *Handler) GatewayUpdate(
	ctx context.Context, req *relay_chain.GatewayUpdateRequest) (*relay_chain.GatewayUpdateResponse, error) {
	h.printRequest(ctx, "GatewayRegister", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		req.GatewayInfo.RelayChainId = conf.BaseConf.GatewayID
		gatewayId, err := gateway.GatewayV1.GatewayUpdate(req.GatewayInfo)
		if err != nil {
			return getGatewayUpdateReturn(common.Code_CONTRACT_FAIL, "", err.Error())
		}
		return getGatewayUpdateReturn(common.Code_GATEWAY_SUCCESS, gatewayId, common.Code_GATEWAY_SUCCESS.String())
	default:
		return getGatewayUpdateReturn(
			common.Code_INVALID_PARAMETER, "", utils.UnsupportVersion(req.Version))
	}
}

// InitContract 安装区块头同步和spv验证合约
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.InitContractResponse
//  @return error
func (h *Handler) InitContract(
	ctx context.Context, req *relay_chain.InitContractRequest) (*relay_chain.InitContractResponse, error) {
	h.printRequest(ctx, "InitContract", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		// 开始之前，需要先判断gateway是否已经注册，防止恶意gateway的连接
		g, err := gateway.GatewayV1.GetGatewayInfo(req.GatewayId)
		if err != nil || g == nil || g.GatewayId == "" {
			msg := fmt.Sprintf("no such gateway: %s", req.GatewayId)
			h.log.Error(msg)
			return &relay_chain.InitContractResponse{
				Code: common.Code_INVALID_PARAMETER, Message: msg}, nil
		}
		jsonStr, _ := json.Marshal(req.KeyValuePairs)
		err = relay_chain_chainmaker.RelayChainV1.InitContract(
			utils.GetSpvContractName(req.GatewayId, req.ChainRid),
			req.ContractVersion,
			string(req.ByteCode),
			string(jsonStr),
			true,
			-1,
			req.RuntimeType)
		if err != nil {
			msg := fmt.Sprintf("[InitContract] init contract failed: error: %s", err.Error())
			h.log.Error(msg)
			return &relay_chain.InitContractResponse{
				Code: common.Code_INVALID_PARAMETER, Message: msg}, nil
		}
		return &relay_chain.InitContractResponse{
			Code: common.Code_GATEWAY_SUCCESS, Message: common.Code_GATEWAY_SUCCESS.String()}, nil
	default:
		return &relay_chain.InitContractResponse{
			Code: common.Code_INVALID_PARAMETER, Message: utils.UnsupportVersion(req.Version)}, nil
	}
}

// UpdateContract 升级区块头同步和spv验证合约
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.UpdateContractResponse
//  @return error
func (h *Handler) UpdateContract(
	ctx context.Context, req *relay_chain.UpdateContractRequest) (*relay_chain.UpdateContractResponse, error) {
	h.printRequest(ctx, "UpdateContract", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		// 开始之前，需要先判断gateway是否已经注册，防止恶意gateway的连接
		g, err := gateway.GatewayV1.GetGatewayInfo(req.GatewayId)
		if err != nil || g == nil || g.GatewayId == "" {
			msg := fmt.Sprintf("no such gateway: %s", req.GatewayId)
			h.log.Error(msg)
			return &relay_chain.UpdateContractResponse{
				Code: common.Code_INVALID_PARAMETER, Message: msg}, nil
		}
		jsonStr, _ := json.Marshal(req.KeyValuePairs)
		err = relay_chain_chainmaker.RelayChainV1.UpdateContract(
			utils.GetSpvContractName(req.GatewayId, req.ChainRid),
			req.ContractVersion,
			string(req.ByteCode),
			string(jsonStr),
			true,
			-1,
			req.RuntimeType)
		if err != nil {
			msg := fmt.Sprintf("[UpdateContract] update contract failed: error: %s", err.Error())
			h.log.Error(msg)
			return nil, errors.New(msg)
		}
		return &relay_chain.UpdateContractResponse{
			Code: common.Code_GATEWAY_SUCCESS, Message: common.Code_GATEWAY_SUCCESS.String()}, nil
	default:
		return &relay_chain.UpdateContractResponse{
			Code: common.Code_INVALID_PARAMETER, Message: utils.UnsupportVersion(req.Version)}, nil
	}
}

// QueryGateway 网关信息查询
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.QueryGatewayResponse
//  @return error
func (h *Handler) QueryGateway(
	ctx context.Context, req *relay_chain.QueryGatewayRequest) (*relay_chain.QueryGatewayResponse, error) {
	h.printRequest(ctx, "QuerGateway", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		if req.GatewayId != "" {
			gatewayInfo, err := gateway.GatewayV1.GetGatewayInfo(req.GatewayId)
			if err != nil {
				return getQueryGatewayReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			return getQueryGatewayReturn(
				common.Code_GATEWAY_SUCCESS,
				[]*common.GatewayInfo{gatewayInfo},
				common.Code_GATEWAY_SUCCESS.String(),
				nil)
		} else if req.PageSize != 0 && req.PageNumber != 0 {
			gatewayNumber, err := gateway.GatewayV1.GetGatewayNum()
			if err != nil {
				return getQueryGatewayReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			if (req.PageNumber-1)*req.PageSize > gatewayNumber {
				return getQueryGatewayReturn(common.Code_INVALID_PARAMETER, nil,
					fmt.Sprintf("out of range, PageSize: %d, PageNumber: %d", req.PageSize, req.PageNumber), nil)
			}

			if gatewayNumber == 0 {
				return getQueryGatewayReturn(common.Code_GATEWAY_SUCCESS, nil, "no gateway", nil)
			}
			pageInfo := &common.PageInfo{
				PageSize:   req.PageSize,
				PageNumber: req.PageNumber,
				TotalCount: gatewayNumber,
				// 向上取整
				Limit: (gatewayNumber + req.PageSize - 1) / req.PageSize,
			}
			startGatewayId := fmt.Sprintf("%d", (req.PageNumber-1)*req.PageSize)
			stopGatewayId := fmt.Sprintf("%d", req.PageNumber*req.PageSize)
			gatewayList, err := gateway.GatewayV1.GetGatewayInfoByRange(startGatewayId, stopGatewayId)
			if err != nil {
				return getQueryGatewayReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			return getQueryGatewayReturn(
				common.Code_GATEWAY_SUCCESS, gatewayList, common.Code_GATEWAY_SUCCESS.String(), pageInfo)
		} else {
			return getQueryGatewayReturn(
				common.Code_INVALID_PARAMETER, nil, common.Code_INVALID_PARAMETER.String(), nil)
		}
	default:
		return getQueryGatewayReturn(
			common.Code_INVALID_PARAMETER, nil, utils.UnsupportVersion(req.Version), nil)
	}
}

// QueryCrossChain 跨链查询
//  @receiver h
//  @param ctx
//  @param req
//  @return *relay_chain.QueryCrossChainResponse
//  @return error
func (h *Handler) QueryCrossChain(
	ctx context.Context, req *relay_chain.QueryCrossChainRequest) (*relay_chain.QueryCrossChainResponse, error) {
	h.printRequest(ctx, "QueryCrossChain", fmt.Sprintf("%+v", req))

	switch req.Version {
	case common.Version_V1_0_0:
		if req.CrossChainId != "" {
			crossChainInfo, err := crosschaintx.CrossChainTxV1.GetCrossChainInfo(req.CrossChainId)
			if err != nil {
				return getQueryCrossChainReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			return getQueryCrossChainReturn(
				common.Code_GATEWAY_SUCCESS,
				[]*common.CrossChainInfo{crossChainInfo},
				common.Code_GATEWAY_SUCCESS.String(),
				nil)
		} else if req.PageSize != 0 && req.PageNumber != 0 {
			crossChainNumber, err := crosschaintx.CrossChainTxV1.GetCrossChainNum()
			if err != nil {
				return getQueryCrossChainReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			if (req.PageNumber-1)*req.PageSize > crossChainNumber {
				return getQueryCrossChainReturn(common.Code_INVALID_PARAMETER, nil,
					fmt.Sprintf("out of range， PageSize: %d, PageNumber: %d", req.PageSize, req.PageNumber), nil)
			}

			if crossChainNumber == 0 {
				return getQueryCrossChainReturn(common.Code_GATEWAY_SUCCESS, nil, "no cross chain info", nil)
			}
			pageInfo := &common.PageInfo{
				PageSize:   req.PageSize,
				PageNumber: req.PageNumber,
				TotalCount: crossChainNumber,
				// 向上取整
				Limit: (crossChainNumber + req.PageSize - 1) / req.PageSize,
			}
			startCrossChainId := fmt.Sprintf("%d", (req.PageNumber-1)*req.PageSize)
			stopCrossChainId := fmt.Sprintf("%d", req.PageNumber*req.PageSize)
			h.log.Debugf("[QueryCrossChain] startCrossChainId: %s, stopCrossChainId: %s", startCrossChainId, stopCrossChainId)
			crossChainInfoList, err := crosschaintx.CrossChainTxV1.GetCrossChainInfoByRange(startCrossChainId, stopCrossChainId)
			if err != nil {
				return getQueryCrossChainReturn(common.Code_CONTRACT_FAIL, nil, err.Error(), nil)
			}
			return getQueryCrossChainReturn(
				common.Code_GATEWAY_SUCCESS, crossChainInfoList, common.Code_GATEWAY_SUCCESS.String(), pageInfo)
		} else {
			return getQueryCrossChainReturn(
				common.Code_INVALID_PARAMETER, nil, common.Code_INVALID_PARAMETER.String(), nil)
		}
	default:
		return getQueryCrossChainReturn(
			common.Code_INVALID_PARAMETER, nil, utils.UnsupportVersion(req.Version), nil)
	}
}

// printRequest 打印请求内容
//  @receiver h
//  @param ctx
//  @param method
//  @param request
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

// getGatewayRegisterReturn 网关注册返回值构建
//  @param code
//  @param gatewayId
//  @param msg
//  @return *relay_chain.GatewayRegisterResponse
//  @return error
func getGatewayRegisterReturn(code common.Code, gatewayId, msg string) (*relay_chain.GatewayRegisterResponse, error) {
	return &relay_chain.GatewayRegisterResponse{
		Code:      code,
		GatewayId: gatewayId,
		Message:   msg,
	}, nil
}

// getGatewayUpdateReturn 网关更新返回值构建
//  @param code
//  @param gatewayId
//  @param msg
//  @return *relay_chain.GatewayUpdateResponse
//  @return error
func getGatewayUpdateReturn(code common.Code, gatewayId, msg string) (*relay_chain.GatewayUpdateResponse, error) {
	return &relay_chain.GatewayUpdateResponse{
		Code:      code,
		GatewayId: gatewayId,
		Message:   msg,
	}, nil
}

// getQueryGatewayReturn 网关查询返回值构建
//  @param code
//  @param gatewayInfo
//  @param msg
//  @param pageInfo
//  @return *relay_chain.QueryGatewayResponse
//  @return error
func getQueryGatewayReturn(
	code common.Code,
	gatewayInfo []*common.GatewayInfo,
	msg string,
	pageInfo *common.PageInfo) (*relay_chain.QueryGatewayResponse, error) {
	return &relay_chain.QueryGatewayResponse{
		Code:        code,
		GatewayInfo: gatewayInfo,
		Message:     msg,
		PageInfo:    pageInfo,
	}, nil
}

// getQueryCrossChainReturn 跨链交易查询返回值构建
//  @param code
//  @param crossChainInfo
//  @param msg
//  @param pageInfo
//  @return *relay_chain.QueryCrossChainResponse
//  @return error
func getQueryCrossChainReturn(
	code common.Code,
	crossChainInfo []*common.CrossChainInfo,
	msg string,
	pageInfo *common.PageInfo) (*relay_chain.QueryCrossChainResponse, error) {
	return &relay_chain.QueryCrossChainResponse{
		Code:           code,
		CrossChainInfo: crossChainInfo,
		Message:        msg,
		PageInfo:       pageInfo,
	}, nil
}

// getBeginCrossChainReturn 发送跨链交易返回值构建
//  @param code
//  @param crossChainId
//  @param msg
//  @param log
//  @param beginTime
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func getBeginCrossChainReturn(
	code common.Code, crossChainId string, msg string,
	log *zap.SugaredLogger, beginTime int64) (*relay_chain.BeginCrossChainResponse, error) {
	log.Debugf("[BeginCrossChain] return %s, code: %s, time used: %d",
		crossChainId, common.Code_name[int32(code)], time.Now().Unix()-beginTime)
	return &relay_chain.BeginCrossChainResponse{Code: code, CrossChainId: crossChainId, Message: msg}, nil
}
