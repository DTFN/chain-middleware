#!/bin/bash

pid=$(cat ./pid | sed -n "1p")

if [ $1 == up ]
then
	nohup ./tcip-chainmaker start -c ./config/tcip_chainmaker.yml > panic.log 2>&1 & echo $! > pid
	echo "tcip-chainmaker start"
	cat ./pid
	exit
fi
if [ $1 == down ]
then
	kill $pid
	echo "tcip-chainmaker stop: $pid"
	exit
fi
echo "error parameter"