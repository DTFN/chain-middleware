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

	"chainmaker.org/chainmaker/tcip-go/v2/common"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/crosschaintx"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/gateway"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/logger"
	relay_chain_chainmaker "chainmaker.org/chainmaker/tcip-relayer/v2/module/relaychain/chainmaker"
	"chainmaker.org/chainmaker/tcip-relayer/v2/module/request"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"
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
	conf.Config.RelayChain = &conf.RelayChain{}
	conf.BaseConf = &conf.BaseConfig{
		GatewayID: "relay1",
	}
	logger.InitLogConfig(log)
	_ = relay_chain_chainmaker.InitRelayChainMock(nil)
	_ = request.InitRequestManagerMock()
	_ = gateway.InitGatewayManagerMock()
}

func TestHandler_BeginCrossChain(t *testing.T) {
	testInit()
	_ = crosschaintx.InitCrossChainManagerMock()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
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
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.BeginCrossChainRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.BeginCrossChainResponse{
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
				req: &relay_chain.BeginCrossChainRequest{
					Version:        common.Version_V1_0_0,
					CrossChainName: "test",
					CrossChainFlag: "test",
					CrossChainMsg:  []*common.CrossChainMsg{},
					From:           "0",
					TxContent: &common.TxContent{
						TxProve: "{}",
					},
				},
			},
			want: &relay_chain.BeginCrossChainResponse{
				CrossChainId: "0",
				Code:         common.Code_GATEWAY_SUCCESS,
				Message:      common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.BeginCrossChain(tt.args.ctx, tt.args.req)
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

func TestHandler_GatewayRegister(t *testing.T) {
	testInit()
	var gatewayInfo common.GatewayInfo
	_ = json.Unmarshal([]byte(gatewayInfoString), &gatewayInfo)
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *relay_chain.GatewayRegisterRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.GatewayRegisterResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.GatewayRegisterRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.GatewayRegisterResponse{
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
				req: &relay_chain.GatewayRegisterRequest{
					Version:     common.Version_V1_0_0,
					GatewayInfo: &gatewayInfo,
				},
			},
			want: &relay_chain.GatewayRegisterResponse{
				Code:      common.Code_GATEWAY_SUCCESS,
				Message:   common.Code_GATEWAY_SUCCESS.String(),
				GatewayId: "0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.GatewayRegister(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayRegister() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GatewayRegister() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_GatewayUpdate(t *testing.T) {
	testInit()
	var gatewayInfo common.GatewayInfo
	_ = json.Unmarshal([]byte(gatewayInfoString), &gatewayInfo)
	gatewayInfo.GatewayId = "0"
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *relay_chain.GatewayUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.GatewayUpdateResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.GatewayUpdateRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.GatewayUpdateResponse{
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
				req: &relay_chain.GatewayUpdateRequest{
					Version:     common.Version_V1_0_0,
					GatewayInfo: &gatewayInfo,
				},
			},
			want: &relay_chain.GatewayUpdateResponse{
				Code:      common.Code_GATEWAY_SUCCESS,
				Message:   common.Code_GATEWAY_SUCCESS.String(),
				GatewayId: "0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.GatewayUpdate(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GatewayUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GatewayUpdate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_InitContract(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
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
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.InitContractRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.InitContractResponse{
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
				req: &relay_chain.InitContractRequest{
					Version:         common.Version_V1_0_0,
					ContractName:    "test",
					ContractVersion: "1.0",
					ByteCode:        []byte("123"),
					RuntimeType:     common.ChainmakerRuntimeType_DOCKER_GO,
					KeyValuePairs:   []*common.ContractKeyValuePair{},
					GatewayId:       "0",
					ChainRid:        "chain1",
				},
			},
			want: &relay_chain.InitContractResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.InitContract(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_QueryCrossChain(t *testing.T) {
	testInit()
	_ = crosschaintx.InitCrossChainManagerMock()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *relay_chain.QueryCrossChainRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.QueryCrossChainResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.QueryCrossChainRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.QueryCrossChainResponse{
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
				req: &relay_chain.QueryCrossChainRequest{
					Version:      common.Version_V1_0_0,
					CrossChainId: "0",
				},
			},
			want: &relay_chain.QueryCrossChainResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				CrossChainInfo: []*common.CrossChainInfo{
					{
						CrossChainId:   "0",
						CrossChainName: "test",
						CrossChainFlag: "test",
						CrossChainMsg: []*common.CrossChainMsg{
							{
								GatewayId:    "0",
								ChainRid:     "chain1",
								ContractName: "test",
								Method:       "test",
								Parameter:    "{\"a\":\"MA==\"}",
								ConfirmInfo: &common.ConfirmInfo{
									ChainRid: "chain1",
								},
								CancelInfo: &common.CancelInfo{
									ChainRid: "chain1",
								},
							},
						},
						FirstTxContent: &common.TxContentWithVerify{},
						From:           "0",
						Timeout:        1000,
						ConfirmInfo: &common.ConfirmInfo{
							ChainRid: "chain1",
						},
						CancelInfo: &common.CancelInfo{
							ChainRid: "chain1",
						},
						CrossType: common.CrossType_INVOKE,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.QueryCrossChainRequest{
					Version:    common.Version_V1_0_0,
					PageSize:   10,
					PageNumber: 1,
				},
			},
			want: &relay_chain.QueryCrossChainResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				CrossChainInfo: []*common.CrossChainInfo{
					{
						CrossChainId:   "0",
						CrossChainName: "test",
						CrossChainFlag: "test",
						CrossChainMsg: []*common.CrossChainMsg{
							{
								GatewayId:    "0",
								ChainRid:     "chain1",
								ContractName: "test",
								Method:       "test",
								Parameter:    "{\"a\":\"MA==\"}",
								ConfirmInfo: &common.ConfirmInfo{
									ChainRid: "chain1",
								},
								CancelInfo: &common.CancelInfo{
									ChainRid: "chain1",
								},
							},
						},
						FirstTxContent: &common.TxContentWithVerify{},
						From:           "0",
						Timeout:        1000,
						ConfirmInfo: &common.ConfirmInfo{
							ChainRid: "chain1",
						},
						CancelInfo: &common.CancelInfo{
							ChainRid: "chain1",
						},
						CrossType: common.CrossType_INVOKE,
					},
				},
				PageInfo: &common.PageInfo{
					PageSize:   10,
					PageNumber: 1,
					TotalCount: 1,
					Limit:      1,
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
			got, err := h.QueryCrossChain(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryCrossChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryCrossChain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_QuerGateway(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
		req *relay_chain.QueryGatewayRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *relay_chain.QueryGatewayResponse
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.QueryGatewayRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.QueryGatewayResponse{
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
				req: &relay_chain.QueryGatewayRequest{
					Version:   common.Version_V1_0_0,
					GatewayId: "0",
				},
			},
			want: &relay_chain.QueryGatewayResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				GatewayInfo: []*common.GatewayInfo{
					{
						GatewayId:    "0",
						TxVerifyType: common.TxVerifyType_SPV,
						RelayChainId: "0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "3",
			fields: fields{
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.QueryGatewayRequest{
					Version:    common.Version_V1_0_0,
					PageSize:   10,
					PageNumber: 1,
				},
			},
			want: &relay_chain.QueryGatewayResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
				GatewayInfo: []*common.GatewayInfo{
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
				PageInfo: &common.PageInfo{
					PageSize:   10,
					PageNumber: 1,
					TotalCount: 5,
					Limit:      1,
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
			got, err := h.QueryGateway(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerGateway() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuerGateway() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandler_SyncBlockHeader(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
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
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.SyncBlockHeaderRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.SyncBlockHeaderResponse{
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
				req: &relay_chain.SyncBlockHeaderRequest{
					Version:     common.Version_V1_0_0,
					GatewayId:   "0",
					ChainRid:    "chain1",
					BlockHeight: 10,
					BlockHeader: []byte("abcd"),
				},
			},
			want: &relay_chain.SyncBlockHeaderResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.SyncBlockHeader(tt.args.ctx, tt.args.req)
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

func TestHandler_UpdateContract(t *testing.T) {
	testInit()
	type fields struct {
		log *zap.SugaredLogger
	}
	type args struct {
		ctx context.Context
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
				log: logger.GetLogger(logger.ModuleHandler),
			},
			args: args{
				ctx: context.Background(),
				req: &relay_chain.UpdateContractRequest{
					Version: common.Version(10),
				},
			},
			want: &relay_chain.UpdateContractResponse{
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
				req: &relay_chain.UpdateContractRequest{
					Version:         common.Version_V1_0_0,
					ContractName:    "test",
					ContractVersion: "1.0",
					ByteCode:        []byte("123"),
					RuntimeType:     common.ChainmakerRuntimeType_DOCKER_GO,
					KeyValuePairs:   []*common.ContractKeyValuePair{},
					GatewayId:       "0",
					ChainRid:        "chain1",
				},
			},
			want: &relay_chain.UpdateContractResponse{
				Code:    common.Code_GATEWAY_SUCCESS,
				Message: common.Code_GATEWAY_SUCCESS.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log: tt.fields.log,
			}
			got, err := h.UpdateContract(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateContract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateContract() got = %v, want %v", got, tt.want)
			}
		})
	}
}
