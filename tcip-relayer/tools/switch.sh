#!/bin/bash

pid=$(cat ./pid | sed -n "1p")

if [ $1 == up ]
then
	nohup ./tcip-relayer start -c ./config/tcip_relayer.yml > panic.log 2>&1 & echo $! > pid
	echo "tcip-relayer start"
	cat ./pid
	exit
fi
if [ $1 == down ]
then
	kill $pid
	echo "tcip-relayer stop: $pid"
	exit
fi
echo "error parameter"