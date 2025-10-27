#!/bin/bash

./tcip-lingshu register -c config/tcip_chainmaker.yml

./tcip-lingshu update -c config/tcip_chainmaker.yml

./tcip-lingshu spv -c config/tcip_chainmaker.yml \
-v 1.0 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O install

./tcip-lingshu spv -c config/tcip_chainmaker.yml \
-v 1.1 \
-p ./contract_demo/spv0chain2.7z \
-r DOCKER_GO \
-P "{}" \
-C chain2 \
-O update