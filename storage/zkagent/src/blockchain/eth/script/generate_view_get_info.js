// encode.js
import { ethers } from "ethers";
import fs from "fs";

const abiPath = process.argv[2];

// 加载 ABI
const abi = JSON.parse(fs.readFileSync(abiPath, "utf8"));

// 1. 创建接口对象
const iface = new ethers.Interface(abi);

// 2. 要调用的函数和参数
const challengeId = process.argv[3];
const data = iface.encodeFunctionData("getChallengeInfo", [challengeId]);

process.stdout.write(data);
