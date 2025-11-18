#!/bin/bash
docker_engine_port=`grep "docker_engine_port" /chain/node/config.yaml | grep -o "[0-9]\+"`
port_occupied=$(netstat -alt | grep :$docker_engine_port | wc -l)
if [ $port_occupied -eq "1" ];then
    echo "port:$docker_engine_port had occupied." >> /chain/node/console.log
    exit 1
fi
export LSVM_ENGINE_PORT=$docker_engine_port
cd /vm-engine
nohup ./bin/vm-engine config/config.yaml >> nohup.out 2>&1 &
sleep 2
exist=$(ps -ef | grep "vm-engine" | grep -v "grep" | wc -l)
if [ $exist -ne "1" ];then
    echo "$exist, vm-engine start failed." >> /chain/node/console.log
    exit 1
fi

cd /chain/node/
nohup /chain/bin/LingShuChain -c config.yaml >>nohup.out 2>&1 &
sleep 2
exist=$(ps -ef | grep "LingShuChain" | grep -v "grep" | wc -l)
if [ $exist -ne "1" ];then
    echo "$exist, node start failed." >> /chain/node/console.log
    exit 1
fi

echo "start success." >> /chain/node/console.log
exit 0