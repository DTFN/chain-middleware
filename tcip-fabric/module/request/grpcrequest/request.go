/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package grpcrequest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"chainmaker.org/chainmaker/tcip-fabric/v2/module/conf"

	"chainmaker.org/chainmaker/tcip-go/v2/common/relay_chain"

	"chainmaker.org/chainmaker/tcip-go/v2/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GrpcRequest grpc请求结构体
type GrpcRequest struct {
	// 后续可以实现一个connection pool
	//conn map[string]api.RpcCrossChainClient
	log *zap.SugaredLogger
}

// NewGrpcRequest 初始化grpc请求
//  @param log
//  @return *GrpcRequest
func NewGrpcRequest(log *zap.SugaredLogger) *GrpcRequest {
	return &GrpcRequest{
		log: log,
	}
}

// BeginCrossChain 调用跨链接口
//  @receiver g
//  @param req
//  @return *relay_chain.BeginCrossChainResponse
//  @return error
func (g *GrpcRequest) BeginCrossChain(
	req *relay_chain.BeginCrossChainRequest) (*relay_chain.BeginCrossChainResponse, error) {
	timeout := conf.Config.BaseConfig.DefaultTimeout
	client, conn, err := g.getConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := client.BeginCrossChain(ctx, req)
	_ = conn.Close()
	return response, err
}

// getConnection 后续可以实现一个connection pool
//  @receiver g
//  @return api.RpcRelayChainClient
//  @return *grpc.ClientConn
//  @return error
func (g *GrpcRequest) getConnection() (api.RpcRelayChainClient, *grpc.ClientConn, error) {
	var (
		conn *grpc.ClientConn
	)
	cert, err := tls.LoadX509KeyPair(conf.Config.Relay.ClientCert, conf.Config.Relay.ClientKey)
	if err != nil {
		msg := "[getConnection] key and cert is not match"
		g.log.Errorf(msg)
		return nil, nil, err
	}
	certPool := x509.NewCertPool()

	caCrt, err := ioutil.ReadFile(conf.Config.Relay.Tlsca)
	if err != nil {
		msg := "[getConnection] read tlaca failed"
		g.log.Errorf(msg)
		return nil, nil, err
	}
	certPool.AppendCertsFromPEM(caCrt)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   conf.Config.Relay.ServerName,
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
	})
	conn, err = grpc.Dial(conf.Config.Relay.Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		msg := fmt.Sprintf("[getConnection] connect to gateway error, address: %s, serverName: %s, error: %s",
			conf.Config.Relay.Address, conf.Config.Relay.ServerName, err.Error())
		g.log.Errorf(msg)
		return nil, nil, errors.New(msg)
	}
	return api.NewRpcRelayChainClient(conn), conn, nil
}
