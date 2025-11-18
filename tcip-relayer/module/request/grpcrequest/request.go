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

	"chainmaker.org/chainmaker/tcip-relayer/v2/module/conf"

	"google.golang.org/protobuf/types/known/emptypb"

	"chainmaker.org/chainmaker/tcip-go/v2/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"chainmaker.org/chainmaker/tcip-go/v2/common/cross_chain"

	"go.uber.org/zap"

	"chainmaker.org/chainmaker/tcip-go/v2/common"
)

// GrpcRequest grpc请求结构体
type GrpcRequest struct {
	// 后续可以实现一个connection pool
	//conn map[string]api.RpcCrossChainClient
	log *zap.SugaredLogger
}

//NewGrpcRequest 初始化grpc请求
//  @param log
//  @return *GrpcRequest
func NewGrpcRequest(log *zap.SugaredLogger) *GrpcRequest {
	return &GrpcRequest{
		log: log,
	}
}

// CrossChainTry 跨链执行
//  @receiver g
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainTryResponse
//  @return error
func (g *GrpcRequest) CrossChainTry(
	txRequest *cross_chain.CrossChainTryRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainTryResponse, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}
	client, conn, err := g.getConnection(destGatewayInfo)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := client.CrossChainTry(ctx, txRequest)
	// @todo 后续实现连接池以后可以复用，现在先按照一个请求，一个连接，请求完之后关闭
	_ = conn.Close()
	return response, err
}

// CrossChainConfirm 跨链提交
//  @receiver g
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainConfirmResponse
//  @return error
func (g *GrpcRequest) CrossChainConfirm(
	txRequest *cross_chain.CrossChainConfirmRequest,
	timeout int64, destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainConfirmResponse, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}
	client, conn, err := g.getConnection(destGatewayInfo)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := client.CrossChainConfirm(ctx, txRequest)
	// @todo 后续实现连接池以后可以复用，现在先按照一个请求，一个连接，请求完之后关闭
	_ = conn.Close()
	return response, err
}

// CrossChainCancel 跨链回滚
//  @receiver g
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.CrossChainCancelResponse
//  @return error
func (g *GrpcRequest) CrossChainCancel(
	txRequest *cross_chain.CrossChainCancelRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.CrossChainCancelResponse, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}
	client, conn, err := g.getConnection(destGatewayInfo)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := client.CrossChainCancel(ctx, txRequest)
	// @todo 后续实现连接池以后可以复用，现在先按照一个请求，一个连接，请求完之后关闭
	_ = conn.Close()
	return response, err
}

// IsCrossChainSuccess 询问跨链是否成功
//  @receiver g
//  @param txRequest
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.IsCrossChainSuccessResponse
//  @return error
func (g *GrpcRequest) IsCrossChainSuccess(
	txRequest *cross_chain.IsCrossChainSuccessRequest,
	timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.IsCrossChainSuccessResponse, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}
	client, conn, err := g.getConnection(destGatewayInfo)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, err := client.IsCrossChainSuccess(ctx, txRequest)
	// @todo 后续实现连接池以后可以复用，现在先按照一个请求，一个连接，请求完之后关闭
	_ = conn.Close()
	return response, err
}

// PingPong 心跳
//  @receiver g
//  @param timeout
//  @param destGatewayInfo
//  @return *cross_chain.PingPongResponse
//  @return error
func (g *GrpcRequest) PingPong(timeout int64,
	destGatewayInfo *common.GatewayInfo) (*cross_chain.PingPongResponse, error) {
	if timeout < 0 {
		timeout = conf.Config.BaseConfig.DefaultTimeout
	}
	client, conn, err := g.getConnection(destGatewayInfo)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	res, err := client.PingPong(ctx, &emptypb.Empty{})
	// @todo 后续实现连接池以后可以复用，现在先按照一个请求，一个连接，请求完之后关闭
	_ = conn.Close()
	// 如果这个错误是空的，那表明节点存活
	return res, err
}

// VerifyTx 交易验证
//  @receiver g
//  @param txVerifyInterface
//  @param txProve
//  @return []byte
//  @return error
func (g *GrpcRequest) VerifyTx(txVerifyInterface *common.TxVerifyInterface, txProve string) ([]byte, error) {
	return nil, errors.New("implement this method")
}

// getConnection 后续可以实现一个connection pool
//  @receiver g
//  @param destGatewayInfo
//  @return api.RpcCrossChainClient
//  @return *grpc.ClientConn
//  @return error
func (g *GrpcRequest) getConnection(
	destGatewayInfo *common.GatewayInfo) (api.RpcCrossChainClient, *grpc.ClientConn, error) {
	//conn, ok := g.conn[address]
	//if ok {
	//	return conn, nil
	//}
	var (
		conn *grpc.ClientConn
		err1 error
	)
	if destGatewayInfo.Tlsca != "" {
		clientKey, err := ioutil.ReadFile(conf.Config.RpcConfig.ClientKey)
		if err != nil {
			msg := fmt.Sprintf("[getConnection] Client key not exist, address: %s", conf.Config.RpcConfig.ClientKey)
			g.log.Errorf(msg)
			return nil, nil, errors.New(msg)
		}
		clientCert, err := ioutil.ReadFile(conf.Config.RpcConfig.ClientCert)
		if err != nil {
			msg := fmt.Sprintf("[getConnection] Client key not exist, address: %s", conf.Config.RpcConfig.ClientKey)
			g.log.Errorf(msg)
			return nil, nil, errors.New(msg)
		}
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			msg := "[getConnection] key and cert is not match"
			g.log.Errorf(msg)
			return nil, nil, err
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM([]byte(destGatewayInfo.Tlsca))

		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{cert},
			ServerName:   destGatewayInfo.ServerName,
			RootCAs:      certPool,
			MinVersion:   tls.VersionTLS12,
		})
		conn, err1 = grpc.Dial(destGatewayInfo.Address, grpc.WithTransportCredentials(creds))
	} else {
		conn, err1 = grpc.Dial(destGatewayInfo.Address, grpc.WithInsecure())
	}
	if err1 != nil {
		msg := fmt.Sprintf("[getConnection] connect to gateway error, address: %s, serverName: %s, error: %s",
			destGatewayInfo.Address, destGatewayInfo.ServerName, err1.Error())
		g.log.Errorf(msg)
		return nil, nil, errors.New(msg)
	}
	//g.conn[address] = conn
	return api.NewRpcCrossChainClient(conn), conn, nil
}
