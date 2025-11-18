#!/bin/bash

mkdir -p fisco && cd fisco
curl -#LO https://github.com/FISCO-BCOS/FISCO-BCOS/releases/download/v2.9.0/build_chain.sh && chmod u+x build_chain.sh

bash build_chain.sh -l 127.0.0.1:4 -p 30300,20200,8545

bash nodes/127.0.0.1/start_all.sh
ps -ef | grep -v grep | grep fisco-bcos