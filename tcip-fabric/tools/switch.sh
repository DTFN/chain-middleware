#!/bin/bash

pid=$(cat ./pid | sed -n "1p")

if [ $1 == up ]
then
	nohup ./tcip-fabric start -c ./config/tcip_fabric.yml > panic.log 2>&1 & echo $! > pid
	echo "tcip-fabric start"
	cat ./pid
	exit
fi
if [ $1 == down ]
then
	kill $pid
	echo "tcip-fabric stop: $pid"
	exit
fi
echo "error parameter"