#!/bin/bash

if [ ! -d "./release" ];then
  mkdir release
else
  echo "release文件夹已经存在"
fi

make

mv ./tcip-bcos ./release/
cp -rf config ./release/
cp tools/switch.sh ./release/