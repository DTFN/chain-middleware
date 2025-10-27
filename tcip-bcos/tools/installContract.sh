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
--params="{\"method\":\"cross_chain_transfer\",\"cross_chain_name\":\"first_event\",\"cross_chain_flag\":\"first_event\",\"confirm_info\":\"{\\\"chain_id\\\":\\\"chain2\\\",\\\"contract_name\\\":\\\"crosschain1\\\",\\\"method\\\":\\\"cross_chain_confirm\\\"}\",\"cancel_info\":\"{\\\"chain_id\\\":\\\"chain2\\\",\\\"contract_name\\\":\\\"crosschain1\\\",\\\"method\\\":\\\"cross_chain_cancel\\\"}\",\"cross_chain_msgs\":\"[{\\\"gateway_id\\\": \\\"1\\\",\\\"chain_id\\\":\\\"chain3\\\",\\\"contract_name\\\":\\\"crosschain2\\\",\\\"method\\\":\\\"invoke_contract\\\",\\\"parameter\\\":\\\"{\\\\\\\"method\\\\\\\":\\\\\\\"cross_chain_try\\\\\\\"}\\\",\\\"confirm_info\\\":{\\\"chain_id\\\":\\\"chain3\\\",\\\"contract_name\\\":\\\"crosschain2\\\",\\\"method\\\":\\\"cross_chain_confirm\\\"},\\\"cancel_info\\\":{\\\"chain_id\\\":\\\"chain3\\\",\\\"contract_name\\\":\\\"crosschain2\\\",\\\"method\\\":\\\"cross_chain_cancel\\\"},\\\"extra_data\\\":\\\"按需写，目标网关能解析就行\\\"}]\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=crosschain1 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config2.yml \
--params="{\"method\":\"query\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=crosschain2 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config2.yml \
--params="{\"method\":\"query\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=spv0chain2 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"method\":\"get_block_header\",\"block_height\":\"3\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=spv0chain2 \
--method=invoke_contract \
--sdk-conf-path=./sdk_config_chain1.yml \
--params="{\"method\":\"get_block_header\",\"block_height\":\"3\"}" \
--sync-result=true

./cmc client contract user invoke \
--contract-name=spv0chain2 \
--method=invoke_contract \
--sdk-conf-path=./testdata/sdk_config.yml \
--params="{\"method\":\"verify_tx\",\"hash_array\":\"WyJOWFJ0cVgva1A4Y2NtRlUxZDlaMUloYTdQSTZvc3MwZkZtWDZQNjdjUGI4PSJd\",\"hash_type\":\"U0hBMjU2\",\"tx_byte\":\"CtEGCrADCgZjaGFpbjIaQDEwOWU1NDViYTU5NDRhMmM5MWMyNzNmNDZmNWRlYTBlMjhlZDQ4NGY0ZDI0NGQ5MzhlMGMzZmY2Y2YwNjRhYWUguL/BlAYyC2Nyb3NzY2hhaW4xOg9pbnZva2VfY29udHJhY3RCHgoGbWV0aG9kEhRjcm9zc19jaGFpbl90cmFuc2ZlckIfChBjcm9zc19jaGFpbl9uYW1lEgtmaXJzdF9ldmVudEIfChBjcm9zc19jaGFpbl9mbGFnEgtmaXJzdF9ldmVudELdAQoQY3Jvc3NfY2hhaW5fbXNncxLIAVt7ImdhdGV3YXlfaWQiOiAiMCIsImNoYWluX2lkIjoiY2hhaW4yIiwiY29udHJhY3RfbmFtZSI6ImNyb3NzY2hhaW4yIiwibWV0aG9kIjoiaW52b2tlX2NvbnRyYWN0IiwicGFyYW1ldGVyIjoie1wibWV0aG9kXCI6XCJjcm9zc19jaGFpbl90cnlcIn0iLCJleHRyYV9kYXRhIjoi5oyJ6ZyA5YaZ77yM55uu5qCH572R5YWz6IO96Kej5p6Q5bCx6KGMIn1dEnMKKQoWd3gtb3JnMS5jaGFpbm1ha2VyLm9yZxAEGg1teV9jZXJ0X2FsaWFzEkYwRAIgBEbwRIsi6S4S8nqHqXjv2U6p97oDapzbHSVRJbBJLKcCIAXw3Bkpf+2NprqRE00HTTji0OLpMn/uXBf3poPjpWp/IqYCEoECEgdzdWNjZXNzGgdTdWNjZXNzIPJgKukBCgR0ZXN0EkAxMDllNTQ1YmE1OTQ0YTJjOTFjMjczZjQ2ZjVkZWEwZTI4ZWQ0ODRmNGQyNDRkOTM4ZTBjM2ZmNmNmMDY0YWFlGgtjcm9zc2NoYWluMSIDMS4wKowBGgtmaXJzdF9ldmVudCILZmlyc3RfZXZlbnQqcAoBMBIGY2hhaW4yGgtjcm9zc2NoYWluMiIPaW52b2tlX2NvbnRyYWN0Mhx7Im1ldGhvZCI6ImNyb3NzX2NoYWluX3RyeSJ9OifmjInpnIDlhpnvvIznm67moIfnvZHlhbPog73op6PmnpDlsLHooYwaIHgypsZMO4fbcR/mScrrGDpaQXHRVFsOn9NlfGEueXB0EC4ouL/BlAY=\"}" \
--sync-result=true

./cmc query tx 16f73509be8792ebca4a363d01020ae0184aa79c55ae400483fcafc987537059 \
--chain-id=chain2 \
--sdk-conf-path=./sdk_config_chain2.yml


./cmc client contract user upgrade \
--contract-name=spv0chain2 \
--runtime-type=DOCKER_GO \
--byte-code-path=/root/chainmaker-contract-sdk-docker-go/spv0chain2.7z \
--version=1.1 \
--sdk-conf-path=./testdata/sdk_config.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.tls.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.tls.crt \
--sync-result=true \
--params="{}"