#!/bin/bash

./cmc client contract user create \
--contract-name=crosschain1 \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/crosschain1.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config2.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"

./cmc client contract user upgrade \
--contract-name=crosschain1 \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/crosschain1.7z \
--version=1.1 \
--sdk-conf-path=./testdata/sdk_config2.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"

./cmc client contract user create \
--contract-name=crosschain2 \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/crosschain2.7z \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config2.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"

./cmc client contract user upgrade \
--contract-name=crosschain2 \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/crosschain2.7z \
--version=1.1 \
--sdk-conf-path=./testdata/sdk_config2.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"

./cmc client contract user upgrade \
--contract-name=cross_chain_manager \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/cross_chain_manager.7z \
--version=1.1 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"

./cmc client contract user invoke \
--contract-name=crosschain1 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config2.yml \
--params="{\"method\":\"cross_chain_transfer\",\"cross_chain_name\":\"first_event\",\"cross_chain_flag\":\"first_event\",\"cross_chain_msgs\":\"[{\\\"gateway_id\\\": \\\"0\\\",\\\"chain_id\\\":\\\"chain2\\\",\\\"contract_name\\\":\\\"crosschain2\\\",\\\"method\\\":\\\"invoke_contract\\\",\\\"parameter\\\":\\\"{\\\\\\\"method\\\\\\\":\\\\\\\"cross_chain_try\\\\\\\"}\\\",\\\"extra_data\\\":\\\"按需写，目标网关能解析就行\\\"}]\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=crosschain2 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config2.yml \
--params="{\"method\":\"query\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=spv0chain1 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config1.yml \
--params="{\"method\":\"get_block_header\",\"block_height\":\"30\"}" \
--sync-result=true