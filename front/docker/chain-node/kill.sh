#!/bin/bash

container_name=$1

pid=`ps -ef | grep "LingShuChain" | grep -v "grep" | awk '{print \$2}'`

stop_flag="false"
for ((i=1; i<=10; i++))
do
    if [ -n "$pid" ]; then
        kill $pid
        sleep 1
    fi

    pid=`ps -ef | grep "LingShuChain" | grep -v "grep" | awk '{print \$2}'`
    if [ -n "$pid" ]; then
        continue
    else
        stop_flag="true"
        break
    fi
done

if [[ "${stop_flag}" == "false" ]]; then
    echo "stop node failed." >> /chain/node/console.log
    echo "failed"
    exit 1
fi

pid=`ps -ef | grep "vm-engine" | grep -v "grep" | awk '{print \$2}'`
if [ -n "$pid" ]; then
    kill $pid
fi
echo "stop success." >> /chain/node/console.log
echo "success"