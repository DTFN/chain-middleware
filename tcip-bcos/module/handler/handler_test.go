/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"context"
	"encoding/json"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/db"
	"github.com/gogo/protobuf/proto"

	bcostypes "github.com/FISCO-BCOS/go-sdk/core/types"

	chain_client "chainmaker.org/chainmaker/tcip-bcos/v2/module/chain-client"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/request"
	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	log = []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			FilePath:     path.Join(os.TempDir(), time.Now().String()),
			LogInConsole: true,
		},
	}
	tx = bcostypes.TransactionDetail{
		Hash:        "234567890",
		BlockNumber: "10",
	}
)

func testInit() {
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
	conf.Config.DbPath = path.Join(os.TempDir(), time.Now().String())
	chainConfigByte, _ := proto.Marshal(&common.BcosConfig{
		ChainRid: "chain1",
	})
	logger.InitLogConfig(log)
	db.NewDbHandle()
	_ = db.Db.Put([]byte("chain#config#chain1"), chainConfigByte)
	_ = request.InitRequestManagerMock()
	_ = chain_client.InitChainClientMock()
}

func TestHandler_CrossChainCancel(t *testing.T) {
	testInit()
	txByte, _ := json.Marshal(tx)
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *cross_chain.CrossChainCancelRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.CrossChainCancelResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainCancelRequest{
					Version: common.Version(10),
				},
			},
			want: &cross_chain.CrossChainCancelResponse{
				Code:    common.Code_INVALID_PARAMETER,
				Message: "Unsupported version: 10",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainCancelRequest{
					Version: common.Version_V1_0_0,
					CancelInfo: &common.CancelInfo{
						ChainRid:     "chain1",
						ContractName: "aaa",
						Method:       "bbb",
						Parameter:    "{\"a\":\"a\"}",
					},
				},
			},
			want: &cross_chain.CrossChainCancelResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				TxContent: &common.TxContent{
					TxId:      tx.Hash,
					Tx:        txByte,
					TxResult:  common.TxResultValue_TX_SUCCESS,
					GatewayId: conf.Config.BaseConfig.GatewayID,
					ChainRid:  "chain1",
					// 这里不验证不需要填
					TxProve:     "",
					BlockHeight: 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.CrossChainCancel(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossChainCancel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrossChainCancel() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CrossChainConfirm(t *testing.T) {
	testInit()
	txByte, _ := json.Marshal(tx)
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *cross_chain.CrossChainConfirmRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.CrossChainConfirmResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainConfirmRequest{
					Version: common.Version(10),
				},
			},
			want: &cross_chain.CrossChainConfirmResponse{
				Code:    common.Code_INVALID_PARAMETER,
				Message: "Unsupported version: 10",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainConfirmRequest{
					Version: common.Version_V1_0_0,
					ConfirmInfo: &common.ConfirmInfo{
						ChainRid:     "chain1",
						ContractName: "aaa",
						Method:       "bbb",
						Parameter:    "{\"a\":\"a\"}",
					},
				},
			},
			want: &cross_chain.CrossChainConfirmResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				TxContent: &common.TxContent{
					TxId:      tx.Hash,
					Tx:        txByte,
					TxResult:  common.TxResultValue_TX_SUCCESS,
					GatewayId: conf.Config.BaseConfig.GatewayID,
					ChainRid:  "chain1",
					// 这里不验证不需要填
					TxProve:     "",
					BlockHeight: 10,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.CrossChainConfirm(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossChainConfirm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrossChainConfirm() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_CrossChainTry(t *testing.T) {
	testInit()
	txByte, _ := json.Marshal(tx)
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *cross_chain.CrossChainTryRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.CrossChainTryResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainTryRequest{
					Version:        common.Version(10),
					CrossChainId:   "0",
					CrossChainFlag: "test",
					CrossChainName: "test",
				},
			},
			want: &cross_chain.CrossChainTryResponse{
				Code:           common.Code_INVALID_PARAMETER,
				Message:        "Unsupported version: 10",
				CrossChainId:   "0",
				CrossChainFlag: "test",
				CrossChainName: "test",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.CrossChainTryRequest{
					Version:        common.Version_V1_0_0,
					CrossChainId:   "0",
					CrossChainFlag: "test",
					CrossChainName: "test",
					CrossChainMsg: &common.CrossChainMsg{
						GatewayId:    "0",
						ChainRid:     "chain1",
						ContractName: "aaa",
						Method:       "bbb",
						Parameter:    "{}",
					},
				},
			},
			want: &cross_chain.CrossChainTryResponse{
				Code:           common.Code_GATEWAY_SUCCESS,
				Message:        common.Code_GATEWAY_SUCCESS.String(),
				CrossChainId:   "0",
				CrossChainFlag: "test",
				CrossChainName: "test",
				TxContent: &common.TxContent{
					TxId:        tx.Hash,
					Tx:          txByte,
					TxResult:    common.TxResultValue_TX_SUCCESS,
					GatewayId:   conf.Config.BaseConfig.GatewayID,
					ChainRid:    "chain1",
					TxProve:     "{}",
					BlockHeight: 16,
				},
				TryResult: []string{"123"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.CrossChainTry(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CrossChainTry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrossChainTry() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_IsCrossChainSuccess(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *cross_chain.IsCrossChainSuccessRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.IsCrossChainSuccessResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &cross_chain.IsCrossChainSuccessRequest{
					Version: common.Version_V1_0_0,
				},
			},
			want: &cross_chain.IsCrossChainSuccessResponse{
				CrossChainResult: false,
				Code:             common.Code_GATEWAY_SUCCESS,
				Message:          common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.IsCrossChainSuccess(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsCrossChainSuccess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IsCrossChainSuccess() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_PingPong(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *emptypb.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.PingPongResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &emptypb.Empty{},
			},
			want: &cross_chain.PingPongResponse{
				ChainOk: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.PingPong(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("PingPong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PingPong() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_TxVerify(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		in  *cross_chain.TxVerifyRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *cross_chain.TxVerifyResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				in: &cross_chain.TxVerifyRequest{
					Version: common.Version(0),
				},
			},
			want: &cross_chain.TxVerifyResponse{
				TxVerifyResult: true,
				Code:           common.Code_GATEWAY_SUCCESS,
				Message:        common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.TxVerify(tt.args.ctx, tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("TxVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TxVerify() got = %v, want %v", got, tt.want)
			}
		})
	}
}
