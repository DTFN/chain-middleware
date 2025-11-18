#!/usr/bin/env bash
# 拉取fabric-sample和工具并部署示例跨链合约
cp -r crosschain1 crosschain2 $GOPATH/src/
if [ ! -d "./fabric-samples" ];then
  curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.0 1.4.7
else
  echo "fabric-samples文件夹已经存在"
fi
#curl -sSL https://bit.ly/2ysbOFE | bash -s -- 2.2.0 1.4.7
cd fabric-samples/test-network
./network.sh up
sleep 3
./network.sh createChannel
sleep 3
./network.sh deployCC -ccn crosschain1 -ccp ../../crosschain1/ -ccl go
./network.sh deployCC -ccn crosschain2 -ccp ../../crosschain2/ -ccl go
rm -rf $GOPATH/src/crosschain1 $GOPATH/src/crosschain2
#cp -rf organizations ../../../config/crypto-config/
cd ../..
docker ps -a | grep "peer\|orderer"