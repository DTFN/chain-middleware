#!/bin/bash

./tcip-chainmaker register -c config/tcip_chainmaker.yml

./tcip-chainmaker update -c config/tcip_chainmaker.yml

./tcip-chainmaker spv -c config/tcip_chainmaker.yml \
-v 1.0 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O install

./tcip-chainmaker spv -c config/tcip_chainmaker.yml \
-v 1.1 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O update