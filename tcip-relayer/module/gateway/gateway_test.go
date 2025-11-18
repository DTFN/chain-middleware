/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package gateway

import (
	"encoding/json"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/emirpasic/gods/maps/treemap"
	"go.uber.org/zap"
)

var log = []*logger.LogModuleConfig{
	{
		ModuleName:   "default",
		FilePath:     path.Join(os.TempDir(), time.Now().String()),
		LogInConsole: true,
	},
}

const gatewayInfoString = "{\"relay_chain_id\":\"relay1\",\"gateway_name\":\"relay_gateway\",\"address\":\"127.0.0.1:19998\"," +
	"\"server_name\":\"www.example.com\",\"tlsca\":\"-----BEGIN CERTIFICATE-----\\nMIIDijCCAnICCQCpQsgRCbWAHzANB" +
	"gkqhkiG9w0BAQsFADCBhjEXMBUGA1UEAwwO\\nKi5leGFtcGxlcy5jb20xCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAw\\nDgY" +
	"DVQQHDAdCZWlqaW5nMQswCQYDVQQKDAJjYTELMAkGA1UECwwCY2ExIDAeBgkq\\nhkiG9w0BCQEWEWVtYWlsQGV4YW1wbGUuY29tMB4XDTIy" +
	"MDMxNzA3MjQwN1oXDTM1\\nMTEyNDA3MjQwN1owgYYxFzAVBgNVBAMMDiouZXhhbXBsZXMuY29tMQswCQYDVQQG\\nEwJDTjEQMA4GA1UECA" +
	"wHQmVpamluZzEQMA4GA1UEBwwHQmVpamluZzELMAkGA1UE\\nCgwCY2ExCzAJBgNVBAsMAmNhMSAwHgYJKoZIhvcNAQkBFhFlbWFpbEBleGF" +
	"tcGxl\\nLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANM+ssybOOYawup1\\n8+Xsd9hWeF9UVEZooeJV07xoXPTHLTMdC" +
	"Xdhm8rfrYdKICB0+6rBm7zxVNSYgZxp\\ndoiQ5bhrQtaKLi/ermFcIMqSbzo7pjOC/FoZqESWcdhVn3lVHllk6x87CSoJ9S0M\\nX95/LgQ" +
	"fRYCThk5TAwBKO21GY7HhLg55hKKWbZ7Zo++nsRnaaUO5KA8TOH29dgnP\\n2SorQHViI7JSFJ3PHVmQWunIY3cQd+LRCfVfIE7GO47uRnC" +
	"sxk+slQZJIVlz2RAG\\n6VYyDzsGT27tHjcJDjqLszMSzIfqWkgzjjlBR0IMyP+wlIXCY618Tb1SSf0r6iXA\\nscXpn8cCAwEAATANBgk" +
	"qhkiG9w0BAQsFAAOCAQEAibgzU7L++2FSVjfFE3EnTIP6\\nDEDyuoK7OEq4BjQtyNVWOaJF73H4imWlCkbHTdI3BjIfAgdC7eZPjBe74kj" +
	"Sxai1\\nyxjQ/wY/sKnalKqP2m5EUs/3ledpZq7rnVCHCKOzXMT7x2hGdTxC5A1ao+2RE9qW\\nylod0qAeS708SQO1ZGwWakxu/Pr1SJXqd" +
	"ksgxWwUbkK4GwE8naNH5u/KS1nDZeK2\\nbNYtrEyQc+mxdGGKgj6CnrDfq2jW4b5SUbGSO9VStxnQuu+vz0tqxy7wtWrqLm0v\\nFmdqAil" +
	"2fidCA1tgquIiLBnQ28iT7tLcDfCZC3TXZtl0RdXkxhYgX0S9Jg0LxQ==\\n-----END CERTIFICATE-----\\n\",\"client_cert\":" +
	"\"-----BEGIN CERTIFICATE-----\\nMIID8zCCAtugAwIBAgIJAKawhvugHF7cMA0GCSqGSIb3DQEBBQUAMIGGMRcwFQYD\\nVQ" +
	"QDDA4qLmV4YW1wbGVzLmNvbTELMAkGA1UEBhMCQ04xEDAOBgNVBAgMB0JlaWpp\\nbmcxEDAOBgNVBAcMB0JlaWppbmcxCzAJBgNVBA" +
	"oMAmNhMQswCQYDVQQLDAJjYTEg\\nMB4GCSqGSIb3DQEJARYRZW1haWxAZXhhbXBsZS5jb20wHhcNMjIwMzE3MDcyNDA4\\nWhcNMzUxMT" +
	"I0MDcyNDA4WjCBmDEWMBQGA1UEAwwNKi5leGFtcGxlLmNvbTELMAkG\\nA1UEBhMCQ04xEDAOBgNVBAgMB0JlaWppbmcxEDAOBgNVBAcMB" +
	"0JlaWppbmcxEzAR\\nBgNVBAoMCk15IENvbXBhbnkxFjAUBgNVBAsMDU15IERlcGFydG1lbnQxIDAeBgkq\\nhkiG9w0BCQEWEWVtYWlsQG" +
	"V4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC\\nAQ8AMIIBCgKCAQEAwoE0/maembY1eK3+DEv7iSub6ra0DG5DfQMs0fc91U+tWL" +
	"57\\nTNp5EvBEQ8IePA1R4ZZN0bjYbulWKnRN7N4CTlCOYcpcr" +
	"KqN37QC9IMitUhFCWNc\\nrsPCNR/ZOUiibEb9mQKUqpOAj933vUBqO3Cyr98gkSxqG/fqK6jSsd0MViqrR3xA\\njrN2sXayGHQiIG" +
	"547jAgcFdv5Fno0PjnqViwTj8POO1tqHfTrQNU7oMNg6qsa3Xs\\navRvZRGtXBQRqqShxW+wgq0W+IDIV6ei4uZX/qXPEDp4wGc4+CWpx" +
	"n1jshqKaJQf\\nLZGxW5qXtHJJiMQdNkRMu6VDGUIBWrkTqETQzwIDAQABo1AwTjAJBgNVHRMEAjAA\\nMAsGA1UdDwQEAwIF4DA0BgNVHR" +
	"EELTArgg93d3cuZXhhbXBsZS5jb22CDSouZXhh\\nbXBsZS5jb22CCWxvY2FsaG9zdDANBgkqhkiG9w0BAQUFAAOCAQEAYJ8yx1yxA89Q\\n" +
	"FnIh6zS95p67QCjh0XejuVOYUFZH+zpEU4jRc+rphUlYWp1/IB2CJr3fIhodwEU2\\nqhsovimMTxcbvdk6kQjS0jkghc25v3BGP4DbNnC" +
	"QpjrAOZuE0wbvaeuAmOUfuzwy\\n/U2Gb97Ka3/UhE/c8Xf4hVDMjFrim2VyqctnfnqFkvzNUaId6wcp43eoSYxOpm9/\\n3rR67KnO/f" +
	"/x21JfaWRQpdWrNcSFkTtYOlKXhooeWF8of64gvWbyZgJku2aKk+sK\\n8IEjA427mZEHFp1gSmRYXFl+e+Gx4t8GfaXQ20dF0ar0DJYp" +
	"ewZnYpUmlfr6QNTv\\nbJAR4PWu/A==\\n-----END CERTIFICATE-----\\n\",\"chain_list\":[{\"chain_id\":\"mychanne" +
	"l\",\"identity\":[\"Org1\",\"Org2\"]}],\"tx_verify_type\":1,\"TxVerifyInterface\":[{\"chain_id\":\"mychan" +
	"nel\",\"address\":\"127.0.0.1:19998/v1/TxVerify\",\"tls_enable\":true,\"tlsca\":\"-----BEGIN CERTIFICAT" +
	"E-----\\nMIIDijCCAnICCQCpQsgRCbWAHzANBgkqhkiG9w0BAQsFADCBhjEXMBUGA1UEAwwO\\nKi5le" +
	"GFtcGxlcy5jb20xCzAJBgNVBAYTAkNOMRAwDgYDVQQIDAdCZWlqaW5nMRAw\\nDgYDVQQHDAdCZWlqaW5nMQswCQYDVQQKDAJjYTELMA" +
	"kGA1UECwwCY2ExIDAeBgkq\\nhkiG9w0BCQEWEWVtYWlsQGV4YW1wbGUuY29tMB4XDTIyMDMxNzA3MjQwN1oXDTM1\\nMTEyNDA" +
	"3MjQwN1owgYYxFzAVBgNVBAMMDiouZXhhbXBsZXMuY29tMQswCQYDVQQG\\nEwJDTjEQMA4GA1UECAwHQmVpamluZzEQMA4GA1UEBwwHQ" +
	"mVpamluZzELMAkGA1UE\\nCgwCY2ExCzAJBgNVBAsMAmNhMSAwHgYJKoZIhvcNAQkBFhFlbWFpbEBleGFtcGxl\\nLmNvbTCCASIwDQY" +
	"JKoZIhvcNAQEBBQADggEPADCCAQoCggEBANM+ssybOOYawup1\\n8+Xsd9hWeF9UVEZooeJV07xoXPTHLTMdCXdhm8rfrYdKICB0+6rB" +
	"m7zxVNSYgZxp\\ndoiQ5bhrQtaKLi/ermFcIMqSbzo7pjOC/FoZqESWcdhVn3lVHllk6x87CSoJ9S0M\\nX95/LgQfRYCThk5TAwBKO21" +
	"GY7HhLg55hKKWbZ7Zo++nsRnaaUO5KA8TOH29dgnP\\n2SorQHViI7JSFJ3PHVmQWunIY3cQd+LRCfVfIE7GO47uRnCsxk+slQZJIVlz2R" +
	"AG\\n6VYyDzsGT27tHjcJDjqLszMSzIfqWkgzjjlBR0IMyP+wlIXCY618Tb1SSf0r6iXA\\nscXpn8cCAwEAATANBgkqhkiG9w0BAQsFAAO" +
	"CAQEAibgzU7L++2FSVjfFE3EnTIP6\\nDEDyuoK7OEq4BjQtyNVWOaJF73H4imWlCkbHTdI3BjIfAgdC7eZPjBe74kjSxai1\\nyxjQ/wY" +
	"/sKnalKqP2m5EUs/3ledpZq7rnVCHCKOzXMT7x2hGdTxC5A1ao+2RE9qW\\nylod0qAeS708SQO1ZGwWakxu/Pr1SJXqdksgxWwUbkK4GwE" +
	"8naNH5u/KS1nDZeK2\\nbNYtrEyQc+mxdGGKgj6CnrDfq2jW4b5SUbGSO9VStxnQuu+vz0tqxy7wtWrqLm0v\\nFmdqAil2fidCA1tgquI" +
	"iLBnQ28iT7tLcDfCZC3TXZtl0RdXkxhYgX0S9Jg" +
	"0LxQ==\\n-----END CERTIFICATE-----\\n\",\"client_cert\":\"-----BEGIN CERTIFICATE-----\\nMIID8zCCAtugAwIBA" +
	"gIJAKawhvugHF7cMA0GCSqGSIb3DQEBBQUAMIGGMRcwFQYD\\nVQQDDA4qLmV4YW1wbGVzLmNvbTELMAkGA1UEBhMCQ04xEDAOBgNVBAgMB" +
	"0JlaWpp\\nbmcxEDAOBgNVBAcMB0JlaWppbmcxCzAJBgNVBAoMAmNhMQswCQYDVQQLDAJjYTEg\\nMB4GCSqGSIb3DQEJARYRZW1haWxA" +
	"ZXhhbXBsZS5jb20wHhcNMjIwMzE3MDcyNDA4\\nWhcNMzUxMTI0MDcyNDA4WjCBmDEWMBQGA1UEAwwNKi5leGFtcGxlLmNvbTELMAkG\\" +
	"nA1UEBhMCQ04xEDAOBgNVBAgMB0JlaWppbmcxEDAOBgNVBAcMB0JlaWppbmcxEzAR\\nBgNVBAoMCk15IENvbXBhbnkxFjAUBgNVBAsM" +
	"DU15IERlcGFydG1lbnQxIDAeBgkq\\nhkiG9w0BCQEWEWVtYWlsQGV4YW1wbGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOC\\nAQ8AMII" +
	"BCgKCAQEAwoE0/maembY1eK3+DEv7iSub6ra0DG5DfQMs0fc91U+tWL57\\nTNp5EvBEQ8IePA1R4ZZN0bjYbulWKnRN7N4CTlCOYc" +
	"pcrKqN37QC9IMitUhFCWNc\\nrsPCNR/ZOUiibEb9mQKUqpOAj933vUBqO3Cyr98gkSxqG/fqK6jSsd0MViqrR3xA\\njrN2sXayGHQi" +
	"IG547jAgcFdv5Fno0PjnqViwTj8POO1tqHfTrQNU7oMNg6qsa3Xs\\navRvZRGtXBQRqqShxW+wgq0W+IDIV6ei4uZX/qXPEDp4wGc4+C" +
	"Wpxn1jshqKaJQf\\nLZGxW5qXtHJJiMQdNkRMu6VDGUIBWrkTqETQzwIDAQABo1AwTjAJBgNVHRMEAjAA\\nMAsGA1UdDwQEAwIF4DA0B" +
	"gNVHREELTArgg93d3cuZXhhbXBsZS5jb22CDSouZXhh\\nbXBsZS5jb22CCWxvY2FsaG9zdDANBgkqhkiG9w0BAQUFAAOCAQEAYJ8yx1y" +
	"xA89Q\\nFnIh6zS95p67QCjh0XejuVOYUFZH+zpEU4jRc+rphUlYWp1/IB2CJr3fIhodwEU2\\nqhsovimMTxcbvdk6kQjS0jkghc25v3B" +
	"GP4DbNnCQpjrAOZuE0wbvaeuAmOUfuzwy\\n/U2Gb97Ka3/UhE/c8Xf4hVDMjFrim2VyqctnfnqFkvzNUaId6wcp43eoSYxOpm9/\\n3rR" +
	"67KnO/f/x21JfaWRQpdWrNcSFkTtYOlKXhooeWF8of64gvWbyZgJku2aKk+sK\\n8IEjA427mZEHFp1gSmRYXFl+e+Gx4t8GfaXQ20dF0ar" +
	"0DJYpewZnYpUmlfr6QNTv\\nbJAR4PWu/A==\\n-----END CERTIFICATE-----\\n\",\"host_name\":\"localhost\"}]," +
	"\"call_type\":1}"

func testInit() {
	logger.InitLogConfig(log)
	conf.Config.RelayChain = &conf.RelayChain{}
	_ = relay_chain_chainmaker.InitRelayChainMock(nil)
	_ = request.InitRequestManagerMock()
	_ = InitGatewayManagerMock()
}

func TestGatewayManager_GatewayRegister(t *testing.T) {
	testInit()
	var gatewayInfo common.GatewayInfo
	_ = json.Unmarshal([]byte(gatewayInfoString), &gatewayInfo)
	type fields struct {
		gatewayInfo  *treemap.Map
		gateWayState *treemap.Map
		log          *zap.SugaredLogger
	}
	type args struct {
		gatewayInfo *common.GatewayInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				gatewayInfo: &common.GatewayInfo{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "2",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				gatewayInfo: &gatewayInfo,
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GatewayManager{
				gatewayInfo:  tt.fields.gatewayInfo,
				gateWayState: tt.fields.gateWayState,
				log:          tt.fields.log,
			}
			got, err := g.GatewayRegister(tt.args.gatewayInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GatewayRegister() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayManager_GatewayUpdate(t *testing.T) {
	testInit()
	var gatewayInfo common.GatewayInfo
	_ = json.Unmarshal([]byte(gatewayInfoString), &gatewayInfo)
	gatewayInfo.GatewayId = "0"
	type fields struct {
		gatewayInfo  *treemap.Map
		gateWayState *treemap.Map
		log          *zap.SugaredLogger
	}
	type args struct {
		gatewayInfo *common.GatewayInfo
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				gatewayInfo: &gatewayInfo,
			},
			want:    "0",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GatewayManager{
				gatewayInfo:  tt.fields.gatewayInfo,
				gateWayState: tt.fields.gateWayState,
				log:          tt.fields.log,
			}
			got, err := g.GatewayUpdate(tt.args.gatewayInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GatewayUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayManager_GetGatewayInfo(t *testing.T) {
	testInit()
	type fields struct {
		gatewayInfo  *treemap.Map
		gateWayState *treemap.Map
		log          *zap.SugaredLogger
	}
	type args struct {
		gatewayId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *common.GatewayInfo
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				gatewayId: "0",
			},
			want: &common.GatewayInfo{
				GatewayId:    "0",
				TxVerifyType: common.TxVerifyType_SPV,
				RelayChainId: "0",
			},
			wantErr: false,
		},
		{
			name: "2",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				gatewayId: "100",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GatewayManager{
				gatewayInfo:  tt.fields.gatewayInfo,
				gateWayState: tt.fields.gateWayState,
				log:          tt.fields.log,
			}
			got, err := g.GetGatewayInfo(tt.args.gatewayId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGatewayInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGatewayInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayManager_GetGatewayInfoByRange(t *testing.T) {
	testInit()
	type fields struct {
		gatewayInfo  *treemap.Map
		gateWayState *treemap.Map
		log          *zap.SugaredLogger
	}
	type args struct {
		startGatewayId string
		stopGatewayId  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*common.GatewayInfo
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			args: args{
				startGatewayId: "0",
				stopGatewayId:  "1",
			},
			want: []*common.GatewayInfo{
				{
					GatewayId:    "0",
					TxVerifyType: common.TxVerifyType_SPV,
					RelayChainId: "0",
				},
				{
					GatewayId:         "1",
					TxVerifyType:      common.TxVerifyType_RPC_INTERFACE,
					TxVerifyInterface: &common.TxVerifyInterface{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GatewayManager{
				gatewayInfo:  tt.fields.gatewayInfo,
				gateWayState: tt.fields.gateWayState,
				log:          tt.fields.log,
			}
			got, err := g.GetGatewayInfoByRange(tt.args.startGatewayId, tt.args.stopGatewayId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGatewayInfoByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGatewayInfoByRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayManager_GetGatewayNum(t *testing.T) {
	testInit()
	type fields struct {
		gatewayInfo  *treemap.Map
		gateWayState *treemap.Map
		log          *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				gatewayInfo:  treemap.NewWithStringComparator(),
				gateWayState: treemap.NewWithStringComparator(),
				log:          logger.GetLogger(logger.ModuleGateway),
			},
			want:    5,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GatewayManager{
				gatewayInfo:  tt.fields.gatewayInfo,
				gateWayState: tt.fields.gateWayState,
				log:          tt.fields.log,
			}
			got, err := g.GetGatewayNum()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGatewayNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGatewayNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGatewayManager_pingPong(t *testing.T) {

}

func Test_checkGateway(t *testing.T) {
	var gatewayInfo common.GatewayInfo
	_ = json.Unmarshal([]byte(gatewayInfoString), &gatewayInfo)
	type args struct {
		gatewayInfo *common.GatewayInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				gatewayInfo: &gatewayInfo,
			},
			wantErr: false,
		},
		{
			name: "2",
			args: args{
				gatewayInfo: &common.GatewayInfo{
					GatewayId:    "0",
					TxVerifyType: common.TxVerifyType_SPV,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkGateway(tt.args.gatewayInfo); (err != nil) != tt.wantErr {
				t.Errorf("checkGateway() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
