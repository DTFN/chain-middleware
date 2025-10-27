/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package prove

import (
	"os"
	"path"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"

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
	conf.Config.RelayChain = &conf.RelayChain{}
	logger.InitLogConfig(log)
	_ = relay_chain_chainmaker.InitRelayChainMock(nil)
	_ = request.InitRequestManagerMock()
	_ = gateway.InitGatewayManagerMock()
}

func TestInitProveManager(t *testing.T) {
	testInit()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitProveManager(); (err != nil) != tt.wantErr {
				t.Errorf("InitProveManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProveManager_ProveTx(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		tx        *common.TxContentWithVerify
		gatewayId string
		chainId   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				logger.GetLogger(logger.ModuleProve),
			},
			args: args{
				tx: &common.TxContentWithVerify{
					TxContent: &common.TxContent{
						BlockHeight: 10,
						TxProve:     "{}",
					},
				},
				gatewayId: "0",
				chainId:   "chain1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				logger.GetLogger(logger.ModuleProve),
			},
			args: args{
				tx: &common.TxContentWithVerify{
					TxContent: &common.TxContent{
						BlockHeight: 10,
						//TxProve:     "{}",
					},
				},
				gatewayId: "0",
				chainId:   "chain1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "3",
			fields: fields{
				logger.GetLogger(logger.ModuleProve),
			},
			args: args{
				tx: &common.TxContentWithVerify{
					TxContent: &common.TxContent{
						BlockHeight: 10,
						//TxProve:     "{}",
					},
				},
				gatewayId: "1",
				chainId:   "chain1",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "4",
			fields: fields{
				logger.GetLogger(logger.ModuleProve),
			},
			args: args{
				tx: &common.TxContentWithVerify{
					TxContent: &common.TxContent{
						BlockHeight: 10,
						TxProve:     "{}",
					},
				},
				gatewayId: "1",
				chainId:   "chain1",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "5",
			fields: fields{
				logger.GetLogger(logger.ModuleProve),
			},
			args: args{
				tx: &common.TxContentWithVerify{
					TxContent: &common.TxContent{
						BlockHeight: 10,
						TxProve:     "{}",
					},
				},
				gatewayId: "2",
				chainId:   "chain1",
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProveManager{
				log: tt.fields.log,
			}
			got, err := p.ProveTx(tt.args.tx, tt.args.gatewayId, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProveTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ProveTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}
