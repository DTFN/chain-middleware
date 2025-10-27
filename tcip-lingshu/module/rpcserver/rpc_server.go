/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package rpcserver

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tmc/grpc-websocket-proxy/wsproxy"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/handler"
	"google.golang.org/grpc/credentials"

	tcipApi "chainmaker.org/chainmaker/tcip-go/v2/api"
	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/logger"

	"chainmaker.org/chainmaker/tcip-lingshu/v2/module/conf"

	"go.uber.org/zap"

	"github.com/cloudflare/cfssl/log"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

var (
	rpcLog *zap.SugaredLogger
)

// RPCServer rpc服务结构体
type RPCServer struct {
	grpcServer *grpc.Server
	config     *conf.RpcConfig
	log        *zap.SugaredLogger
	ctx        context.Context
	cancel     context.CancelFunc
	isShutdown bool
	mixServer  *http.Server
}

// NewRpcServer 新建rpc服务
//
//	@return *RPCServer
//	@return error
func NewRpcServer() (*RPCServer, error) {
	rpcLog = logger.GetLogger(logger.ModuleRpcServer)
	rpcConfig := conf.Config.RpcConfig
	server, err := newGrpc(rpcConfig)
	if err != nil {
		return nil, fmt.Errorf("new grpc server failed, %s", err.Error())
	}

	mixServer, err := newMixServer(server, rpcConfig)
	if err != nil {
		return nil, fmt.Errorf("new http grpc server failed, %s", err.Error())
	}

	return &RPCServer{
		grpcServer: server,
		mixServer:  mixServer,
		config:     rpcConfig,
		log:        rpcLog,
	}, nil
}

// Start start RPCServer
//
//	@receiver s
//	@return error
func (s *RPCServer) Start() error {
	var (
		err error
	)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.isShutdown = false

	if err = s.RegisterHandler(); err != nil {
		return fmt.Errorf("register handler failed, %s", err.Error())
	}

	endPoint := fmt.Sprintf(":%d", s.config.Port)
	if err != nil {
		return fmt.Errorf("TCP listen failed, %s", err.Error())
	}

	go func() {
		err = s.mixServer.ListenAndServeTLS(s.config.TLSConfig.CertFile, s.config.TLSConfig.KeyFile)
		if err != nil {
			s.log.Errorf("grpc Serve failed, %s", err.Error())
		}
	}()

	s.log.Infof("gRPC server listen on %s", endPoint)

	return nil
}

// RegisterHandler register apiservice handler to rpcserver
//
//	@receiver s
//	@return error
func (s *RPCServer) RegisterHandler() error {
	apiHandler := handler.NewHandler()
	tcipApi.RegisterRpcCrossChainServer(s.grpcServer, apiHandler)
	return nil
}

// Stop stop RPCServer
//
//	@receiver s
func (s *RPCServer) Stop() {
	s.isShutdown = true
	s.cancel()
	s.grpcServer.GracefulStop()
	s.log.Info("RPCServer is stopped!")
}

func newGrpc(rpcConfig *conf.RpcConfig) (*grpc.Server, error) {
	var opts []grpc.ServerOption

	opts = []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
			BlackListInterceptor(),
		),
	}

	checkClientAuth := tls.RequireAndVerifyClientCert
	rpcLog.Infof("need check client auth")

	certs, err := tls.LoadX509KeyPair(rpcConfig.TLSConfig.CertFile, rpcConfig.TLSConfig.KeyFile)
	if err != nil {
		log.Errorf("load X509 key pair failed, %s", err.Error())
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(rpcConfig.TLSConfig.CaFile)
	if err != nil {
		log.Errorf("read ca file failed, %s", err.Error())
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certs},
		ClientAuth:   checkClientAuth,
		ClientCAs:    certPool,
		MinVersion:   tls.VersionTLS12,
	})

	opts = append(opts, grpc.Creds(creds))

	opts = append(opts, grpc.MaxSendMsgSize(rpcConfig.MaxSendMsgSize*1024*1024))
	opts = append(opts, grpc.MaxRecvMsgSize(rpcConfig.MaxRecvMsgSize*1024*1024))

	server := grpc.NewServer(opts...)

	return server, nil
}

func newMixServer(grpcServer *grpc.Server, rpcConfig *conf.RpcConfig) (*http.Server, error) {
	restMux := runtime.NewServeMux()
	ctx := context.Background()
	dopts := []grpc.DialOption{}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(rpcConfig.TLSConfig.CaFile)
	if err != nil {
		log.Errorf("read ca file failed, %s", err.Error())
		return nil, err
	}
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewClientTLSFromCert(certPool, rpcConfig.TLSConfig.ServerName)

	dopts = append(dopts, grpc.WithTransportCredentials(creds))
	endPoint := fmt.Sprintf(":%d", rpcConfig.Port)
	if err := tcipApi.RegisterRpcCrossChainHandlerFromEndpoint(ctx, restMux, "localhost"+endPoint, dopts); err != nil {
		log.Errorf("new restful server failed, RegisterRpcRelayChainHandlerFromEndpoint err: %v", err)
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/", restMux)
	handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
	}), &http2.Server{})
	httpServer := &http.Server{Handler: wsproxy.WebsocketProxy(handler, wsproxy.WithMaxRespBodyBufferSize(
		rpcConfig.RestfulConfig.MaxRespBodySize*1024*1024)), Addr: endPoint}

	return httpServer, nil
}
