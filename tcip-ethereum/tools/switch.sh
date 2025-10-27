#!/bin/bash

pid=$(cat ./pid | sed -n "1p")

if [ $1 == up ]
then
	nohup ./tcip-ethereum start -c ./config/tcip_ethereum.yml > panic.log 2>&1 & echo $! > pid
	echo "tcip-ethereum start"
	cat ./pid
	exit
fi
if [ $1 == down ]
then
	kill $pid
	echo "tcip-ethereum stop: $pid"
	exit
fi
echo "error parameter"