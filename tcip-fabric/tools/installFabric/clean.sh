#!/usr/bin/env bash
isAll=$1
cd fabric-samples/test-network
./network.sh down
cd ../../

if [ "${isAll}" = "all" ]; then
    rm -rf ./fabric-samples
    echo "停止并清除fabric网络成功，清除文件成功"
else
    echo "停止并清除fabric网络成功"
fi