##!/usr/bin/env bash
#echo "============================================================================"
#echo "   mac下chainmaker1.4的prepare脚本生成的配置文件会有缩进问题，需要手动修改后手动启动"
#echo "============================================================================"
#sleep 2
## 拉取chainmaker并创建网络安装示例存证合约
#if [ ! -d "./chainmaker-go" ];then
#  git clone git@git.code.tencent.com:ChainMaker/chainmaker-go.git
#else
#  echo "chainmaker-go文件夹已经存在"
#fi
#if [ ! -d "./chainmaker-cryptogen" ];then
#  git clone git@git.code.tencent.com:ChainMaker/chainmaker-cryptogen.git
#else
#  echo "chainmaker-cryptogen"
#fi
##git clone git@git.code.tencent.com:ChainMaker/chainmaker-go.git
##git clone git@git.code.tencent.com:ChainMaker/chainmaker-cryptogen.git
#cd chainmaker-cryptogen
#make
#cd ..
#cd chainmaker-go/tools
#ln -s ../../chainmaker-cryptogen ./
#cd ..
#git checkout -b v2.2.0_alpha_qc origin/v2.2.0_alpha_qc
#cd tools/sdk
#git submodule update --init
#cd ../../
#git submodule update --init
#cd ./scripts/
#echo "设置要启动的网络类型"
#./prepare.sh 4 1
#./build_release.sh
#./cluster_quick_start.sh normal
#sleep 2
#ps -ef | grep chainmaker
#cd ../build/
#cp -r crypto-config ../../../config/crypto-config/
#cp -r crypto-config ../tools/sdk/testdata/
#cd ../tools/cmc
#git submodule update --init
#go build
#cp -r ../sdk/testdata ./
#./cmc client contract user create \
#  --contract-name=fact \
#  --runtime-type=WASMER \
#  --byte-code-path=../../test/wasm/rust-fact-1.2.0.wasm \
#  --version=1.0 \
#  --sdk-conf-path=./testdata/sdk_config.yml \
#  --admin-org-ids=wx-org1.chainmaker.org \
#  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key \
#  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt \
#  --org-id=wx-org1.chainmaker.org \
#  --sync-result=true \
#  --params="{}"
#./cmc client contract user create \
#  --contract-name=cross_chain_save \
#  --runtime-type=WASMER \
#  --byte-code-path=/Users/chengliang/workspaces/tbis-chainmaker-tinygo-contracts/chainmaker-contract-go.wasm \
#  --version=1.0 \
#  --sdk-conf-path=./testdata/sdk_config.yml \
#  --admin-org-ids=wx-org1.chainmaker.org \
#  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key \
#  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt \
#  --org-id=wx-org1.chainmaker.org \
#  --sync-result=true \
#  --params="{}"

#cd ../../..
