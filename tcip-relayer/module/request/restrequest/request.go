/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package restrequest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"
	"go.uber.org/zap"
)

const (
	//DefaultSendTxTimeout = 10

	crossChainTry       = "/v1/crossChainTry"
	crossChainConfirm   = "/v1/crossChainConfirm"
	crossChainCancel    = "/v1/CrossChainCancel"
	isCrossChainSuccess = "/v1/IsCrossChainSuccess"
	pingPong            = "/v1/PingPong"
)

// RestRequest rest请求结构体
type RestRequest struct {
	log *zap.SugaredLogger
}

// NewRestRequest 新建rest请求
func NewRestRequest(log *zap.SugaredLogger) *RestRequest {
	return &RestRequest{
		log: log,
	}
}

// CrossChainTry 跨链试运行
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return res
//  @return err
func (r *RestRequest) CrossChainTry(
	txRequest *cross_chain.CrossChainTryRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (res *cross_chain.CrossChainTryResponse, err error) {
	reqByte, err := json.Marshal(txRequest)
	if err != nil {
		return nil, err
	}
	url := destGatewayInfo.Address + crossChainTry
	body, err := r.post(timeout, destGatewayInfo, string(reqByte), url)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(body, res)
	return res, nil
}

// CrossChainConfirm 跨链提交
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return res
//  @return err
func (r *RestRequest) CrossChainConfirm(
	txRequest *cross_chain.CrossChainConfirmRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (res *cross_chain.CrossChainConfirmResponse, err error) {
	reqByte, err := json.Marshal(txRequest)
	if err != nil {
		return nil, err
	}
	url := destGatewayInfo.Address + crossChainConfirm
	body, err := r.post(timeout, destGatewayInfo, string(reqByte), url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// CrossChainCancel 回滚跨链
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return res
//  @return err
func (r *RestRequest) CrossChainCancel(
	txRequest *cross_chain.CrossChainCancelRequest,
	timeout int64, destGatewayInfo *common.GatewayInfo) (res *cross_chain.CrossChainCancelResponse, err error) {
	reqByte, err := json.Marshal(txRequest)
	if err != nil {
		return nil, err
	}
	url := destGatewayInfo.Address + crossChainCancel
	body, err := r.post(timeout, destGatewayInfo, string(reqByte), url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// IsCrossChainSuccess 询问跨链是否成功
//  @receiver r
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return res
//  @return err
func (r *RestRequest) IsCrossChainSuccess(
	txRequest *cross_chain.IsCrossChainSuccessRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (res *cross_chain.IsCrossChainSuccessResponse, err error) {
	reqByte, err := json.Marshal(txRequest)
	if err != nil {
		return nil, err
	}
	url := destGatewayInfo.Address + isCrossChainSuccess
	body, err := r.post(timeout, destGatewayInfo, string(reqByte), url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PingPong 心跳
//  @receiver r
//  @param timeout
//  @param destGatewayInfo
//  @return res
//  @return err
func (r *RestRequest) PingPong(timeout int64,
	destGatewayInfo *common.GatewayInfo) (res *cross_chain.PingPongResponse, err error) {
	url := destGatewayInfo.Address + pingPong
	body, err := r.post(timeout, destGatewayInfo, "", url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// VerifyTx 交易验证
//  @receiver r
//  @param txVerifyInterface
//  @param txProve
//  @return []byte
//  @return error
func (r *RestRequest) VerifyTx(txVerifyInterface *common.TxVerifyInterface, txProve string) ([]byte, error) {
	timeout := conf.Config.BaseConfig.DefaultTimeout
	var (
		client *http.Client
		url    string
	)
	if txVerifyInterface.Tlsca != "" {
		clientKey, err := ioutil.ReadFile(conf.Config.RpcConfig.ClientKey)
		if err != nil {
			msg := fmt.Sprintf("[post] Client key not exist, address: %s", conf.Config.RpcConfig.ClientKey)
			r.log.Errorf(msg)
			return nil, errors.New(msg)
		}
		cert, err2 := tls.X509KeyPair([]byte(txVerifyInterface.ClientCert), []byte(clientKey))
		if err2 != nil {
			const msg = "[getConnection] key and cert is not match"
			r.log.Errorf(msg)
			return nil, err2
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(txVerifyInterface.Tlsca))

		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      certPool,
					Certificates: []tls.Certificate{cert},
					MinVersion:   tls.VersionTLS12,
					ServerName:   txVerifyInterface.HostName,
				},
			},
		}
		url = "https://" + txVerifyInterface.Address
	} else {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
		url = "http://" + txVerifyInterface.Address
	}
	//r.conn[address] = conn

	resp, err := client.Post(url, "application/json", strings.NewReader(txProve))
	if err != nil {
		return nil, err
	}
	client.CloseIdleConnections()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// post 后续可以实现一个connection pool
//  @receiver r
//  @param timeout
//  @param destGatewayInfo
//  @param reqStr
//  @param url
//  @return []byte
//  @return error
func (r *RestRequest) post(timeout int64, destGatewayInfo *common.GatewayInfo, reqStr, url string) ([]byte, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}

	//conn, ok := r.conn[address]
	//if ok {
	//	return conn, nil
	//}
	var (
		client *http.Client
	)
	if destGatewayInfo.Tlsca != "" {
		clientKey, err := ioutil.ReadFile(conf.Config.RpcConfig.ClientKey)
		if err != nil {
			msg := fmt.Sprintf("[post] Client key not exist, address: %s", conf.Config.RpcConfig.ClientKey)
			r.log.Errorf(msg)
			return nil, errors.New(msg)
		}
		clientCert, err := ioutil.ReadFile(conf.Config.RpcConfig.ClientCert)
		if err != nil {
			msg := fmt.Sprintf("[post] Client key not exist, address: %s", conf.Config.RpcConfig.ClientKey)
			r.log.Errorf(msg)
			return nil, errors.New(msg)
		}
		cert, err2 := tls.X509KeyPair(clientCert, clientKey)
		if err2 != nil {
			msg := "[post] key and cert is not match"
			r.log.Errorf(msg)
			return nil, err2
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(destGatewayInfo.Tlsca))

		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs:      certPool,
					Certificates: []tls.Certificate{cert},
					MinVersion:   tls.VersionTLS12,
				},
			},
		}
	} else {
		client = &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}
	}
	//r.conn[address] = conn

	resp, err := client.Post(url, "application/json", strings.NewReader(reqStr))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
