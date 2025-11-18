/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package request

import (
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"go.uber.org/zap"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

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
	logger.InitLogConfig(log)
	_ = InitRequestManagerMock()
}

func TestRequestManager_GatewayRegister(t *testing.T) {
	testInit()
	type fields struct {
		log     *zap.SugaredLogger
		request Request
	}
	type args struct {
		objectPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log:     logger.GetLogger(logger.ModuleRequest),
				request: RequestV1.request,
			},
			args: args{
				objectPath: path.Join(os.TempDir(), time.Now().String()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestManager{
				log:     tt.fields.log,
				request: tt.fields.request,
			}
			if err := r.GatewayRegister(tt.args.objectPath); (err != nil) != tt.wantErr {
				t.Errorf("GatewayRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestManager_GatewayUpdate(t *testing.T) {
	testInit()
	type fields struct {
		log     *zap.SugaredLogger
		request Request
	}
	type args struct {
		objectPath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log:     logger.GetLogger(logger.ModuleRequest),
				request: RequestV1.request,
			},
			args: args{
				objectPath: path.Join(os.TempDir(), time.Now().String()),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestManager{
				log:     tt.fields.log,
				request: tt.fields.request,
			}
			if err := r.GatewayUpdate(tt.args.objectPath); (err != nil) != tt.wantErr {
				t.Errorf("GatewayUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestManager_InitSpvContracta(t *testing.T) {
	testInit()
	type fields struct {
		log     *zap.SugaredLogger
		request Request
	}
	type args struct {
		version     string
		path        string
		runtimeType string
		kvJsonStr   string
		chainId     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log:     logger.GetLogger(logger.ModuleRequest),
				request: RequestV1.request,
			},
			args: args{
				version:     "1.0",
				path:        "../../contract_demo/contract.sol",
				runtimeType: "DOCKER_GO",
				kvJsonStr:   "{}",
				chainId:     "chain1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestManager{
				log:     tt.fields.log,
				request: tt.fields.request,
			}
			if err := r.InitSpvContracta(tt.args.version, tt.args.path, tt.args.runtimeType, tt.args.kvJsonStr, tt.args.chainId); (err != nil) != tt.wantErr {
				t.Errorf("InitSpvContracta() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRequestManager_UpdateSpvContract(t *testing.T) {
	testInit()
	type fields struct {
		log     *zap.SugaredLogger
		request Request
	}
	type args struct {
		version     string
		path        string
		runtimeType string
		kvJsonStr   string
		chainId     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log:     logger.GetLogger(logger.ModuleRequest),
				request: RequestV1.request,
			},
			args: args{
				version:     "1.0",
				path:        "../../contract_demo/contract.sol",
				runtimeType: "DOCKER_GO",
				kvJsonStr:   "{}",
				chainId:     "chain1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RequestManager{
				log:     tt.fields.log,
				request: tt.fields.request,
			}
			if err := r.UpdateSpvContract(tt.args.version, tt.args.path, tt.args.runtimeType, tt.args.kvJsonStr, tt.args.chainId); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSpvContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getCallType(t *testing.T) {
	testInit()
	type args struct {
		log *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want common.CallType
	}{
		{
			name: "1",
			args: args{
				log: RequestV1.log,
			},
			want: common.CallType_GRPC,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCallType(tt.args.log); got != tt.want {
				t.Errorf("getCallType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getKvsFromKvJsonStr(t *testing.T) {
	type args struct {
		kvJsonStr string
	}
	tests := []struct {
		name    string
		args    args
		want    []*common.ContractKeyValuePair
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				kvJsonStr: "{\"a\":\"b\"}",
			},
			want: []*common.ContractKeyValuePair{
				{
					Key:   "a",
					Value: []byte("b"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getKvsFromKvJsonStr(tt.args.kvJsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("getKvsFromKvJsonStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getKvsFromKvJsonStr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTxVerifyInterface(t *testing.T) {
	testInit()
	var txVerifyInterface *common.TxVerifyInterface
	type args struct {
		log *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want *common.TxVerifyInterface
	}{
		{
			name: "1",
			args: args{
				log: RequestV1.log,
			},
			want: txVerifyInterface,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTxVerifyInterface(tt.args.log); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTxVerifyInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTxVerifyType(t *testing.T) {
	testInit()
	type args struct {
		log *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
		want common.TxVerifyType
	}{
		{
			name: "1",
			args: args{
				RequestV1.log,
			},
			want: common.TxVerifyType_NOT_NEED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTxVerifyType(tt.args.log); got != tt.want {
				t.Errorf("getTxVerifyType() = %v, want %v", got, tt.want)
			}
		})
	}
}
