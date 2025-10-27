/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package grpcrequest

import (
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"chainmaker.org/chainmaker/tcip-go/v2/api"
	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

func testInit() *GrpcRequest {
	logger.InitLogConfig(log)
	conf.Config.BaseConfig = &conf.BaseConfig{
		DefaultTimeout: 10,
	}
	conf.Config.Relay = &conf.Relay{
		Address:    "https://127.0.0.1:19999",
		ServerName: "chainmaker.org",
		Tlsca:      "../../../config/cert/client/ca.crt",
		ClientKey:  "../../../config/cert/client/client.key",
		ClientCert: "../../../config/cert/client/client.crt",
	}
	return NewGrpcRequest(logger.GetLogger(logger.ModuleRequest))
}

func TestGrpcRequest_BeginCrossChain(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		req *relay_chain.BeginCrossChainRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.BeginCrossChainResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			args: args{
				req: &relay_chain.BeginCrossChainRequest{
					Version: common.Version_V1_0_0,
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcRequest{
				log: tt.fields.log,
			}
			got, err := g.BeginCrossChain(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("BeginCrossChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BeginCrossChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRequest_InitSpvContract(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		req *relay_chain.InitContractRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.InitContractResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			args: args{
				req: &relay_chain.InitContractRequest{
					Version: common.Version_V1_0_0,
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcRequest{
				log: tt.fields.log,
			}
			got, err := g.InitSpvContract(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitSpvContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitSpvContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRequest_SyncBlockHeader(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		req *relay_chain.SyncBlockHeaderRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.SyncBlockHeaderResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			args: args{
				req: &relay_chain.SyncBlockHeaderRequest{
					Version: common.Version_V1_0_0,
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcRequest{
				log: tt.fields.log,
			}
			got, err := g.SyncBlockHeader(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncBlockHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SyncBlockHeader() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRequest_UpdateSpvContract(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		req *relay_chain.UpdateContractRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.UpdateContractResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			args: args{
				req: &relay_chain.UpdateContractRequest{
					Version: common.Version_V1_0_0,
				},
			},
			wantErr: true,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcRequest{
				log: tt.fields.log,
			}
			got, err := g.UpdateSpvContract(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateSpvContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateSpvContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrpcRequest_getConnection(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		fields  fields
		want    api.RpcRelayChainClient
		want1   *grpc.ClientConn
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			wantErr: false,
			//want:    nil,
			//want1:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcRequest{
				log: tt.fields.log,
			}
			_, _, err := g.getConnection()
			if (err != nil) != tt.wantErr {
				t.Errorf("getConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("getConnection() got = %v, want %v", got, tt.want)
			//}
			//if !reflect.DeepEqual(got1, tt.want1) {
			//	t.Errorf("getConnection() got1 = %v, want %v", got1, tt.want1)
			//}
		})
	}
}

func TestNewGrpcRequest(t *testing.T) {
	testInit()
	type args struct {
		log *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want *GrpcRequest
	}{
		{
			name: "1",
			args: args{
				log: logger.GetLogger(logger.ModuleRequest),
			},
			want: &GrpcRequest{
				log: logger.GetLogger(logger.ModuleRequest),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGrpcRequest(tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGrpcRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
