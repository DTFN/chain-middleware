/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"testing"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

func TestGetBlockHeaderParam(t *testing.T) {
	type args struct {
		blockHeight int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				blockHeight: 10,
			},
			want: "{\"block_height\":\"MTA=\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBlockHeaderParam(tt.args.blockHeight); got != tt.want {
				t.Errorf("GetBlockHeaderParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSpvContractName(t *testing.T) {
	type args struct {
		gatewayId string
		chainId   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				gatewayId: "0",
				chainId:   "chain1",
			},
			want: "spv0chain1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSpvContractName(tt.args.gatewayId, tt.args.chainId); got != tt.want {
				t.Errorf("GetSpvContractName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSyncBlockHeaderParameter(t *testing.T) {
	type args struct {
		blockHeight uint64
		blockHeader []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				blockHeight: 10,
				blockHeader: []byte("10"),
			},
			want: "{\"block_header\":\"MTA=\",\"block_height\":\"MTA=\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSyncBlockHeaderParameter(tt.args.blockHeight, tt.args.blockHeader); got != tt.want {
				t.Errorf("GetSyncBlockHeaderParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnsupportVersion(t *testing.T) {
	type args struct {
		version common.Version
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				version: 10,
			},
			want: "Unsupported version: 10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnsupportVersion(tt.args.version); got != tt.want {
				t.Errorf("GetSyncBlockHeaderParameter() = %v, want %v", got, tt.want)
			}
		})
	}
}
