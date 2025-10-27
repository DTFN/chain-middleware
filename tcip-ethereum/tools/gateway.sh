#!/bin/bash

./tcip-ethereum register -c config/tcip_chainmaker.yml

./tcip-ethereum update -c config/tcip_chainmaker.yml

./tcip-ethereum spv -c config/tcip_chainmaker.yml \
-v 1.0 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O install

./tcip-ethereum spv -c config/tcip_chainmaker.yml \
-v 1.1 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O update