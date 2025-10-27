/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package rpcserver

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime/debug"
	"strings"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"

	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	//UNKNOWN unknown string
	UNKNOWN = "unknown"
)

// BlackListInterceptor set ip blacklist interceptor
//  @return grpc.UnaryServerInterceptor
func BlackListInterceptor() grpc.UnaryServerInterceptor {

	blackIps := conf.Config.RpcConfig.BlackList

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		interface{}, error) {

		ipAddr := getClientIp(ctx)
		for _, blackIp := range blackIps {
			if ipAddr == blackIp {
				errMsg := fmt.Sprintf("%s is rejected by black list [%s]", info.FullMethod, ipAddr)
				rpcLog.Warn(errMsg)
				return nil, status.Error(codes.ResourceExhausted, errMsg)
			}
		}

		return handler(ctx, req)
	}
}

// LoggingInterceptor set logging interceptor
//  @return unc
func LoggingInterceptor(ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	addr := GetClientAddr(ctx)

	rpcLog.Debugf("[%s] call gRPC method: %s", addr, info.FullMethod)
	str := fmt.Sprintf("req detail: %+v", req)
	if len(str) > 1024 {
		str = str[:1024] + " ......"
	}
	rpcLog.Info(str)
	resp, err := handler(ctx, req)
	rpcLog.Debugf("[%s] call gRPC method: %s, resp detail: %+v", addr, info.FullMethod, resp)
	return resp, err
}

// RecoveryInterceptor set recovery interceptor
//  @return unc
func RecoveryInterceptor(ctx context.Context, req interface{},
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			stack := debug.Stack()
			os.Stderr.Write(stack)
			rpcLog.Errorf("panic stack: %s", string(stack))
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}

func getClientIp(ctx context.Context) string {
	addr := GetClientAddr(ctx)
	return strings.Split(addr, ":")[0]
}

// GetClientAddr 获取客户端地址
//  @param ctx
//  @return string
func GetClientAddr(ctx context.Context) string {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		rpcLog.Errorf("getClientAddr FromContext failed")
		return UNKNOWN
	}

	if pr.Addr == net.Addr(nil) {
		rpcLog.Errorf("getClientAddr failed, peer.Addr is nil")
		return UNKNOWN
	}

	return pr.Addr.String()
}
