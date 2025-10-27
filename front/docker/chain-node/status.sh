#!/bin/bash
# node_dir=$(cd $(dirname $0);pwd)
lingshu_chain=/chain/bin/LingShuChain
node_pid=$(pgrep -f ${lingshu_chain})
if [ -z "${node_pid}" ]; then
    echo "节点未启动."
    exit 1
fi

echo "节点运行中. PID=${node_pid}"
exit 0

