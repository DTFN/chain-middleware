#!/bin/bash

pid=$(cat ./pid | sed -n "1p")

if [ $1 == up ]
then
	nohup ./tcip-bcos start -c ./config/tcip_bcos.yml > panic.log 2>&1 & echo $! > pid
	echo "tcip-bcos start"
	cat ./pid
	exit
fi
if [ $1 == down ]
then
	kill $pid
	echo "tcip-bcos stop: $pid"
	exit
fi
echo "error parameter"