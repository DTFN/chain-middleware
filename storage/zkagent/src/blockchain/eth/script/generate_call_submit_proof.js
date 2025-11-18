#!/usr/bin/env node
import { ethers } from "ethers";
import fs from "fs";

const args = process.argv.slice(2);

if (args.length < 7) {
  console.error("用法: node generate_call_submit_proof.js abi.json challengeId pA pB0 pB1 pC pubSignals");
  console.error('示例: node generate_call_submit_proof.js abi.json 1 "[1,2]" "[3,4]" "[5,6]" "[7,8]" "[9,10,11]"');
  process.exit(1);
}

const [abiPath, challengeId, pA, pB0, pB1, pC, pubSignals] = args;

// 读取 ABI
const abi = JSON.parse(fs.readFileSync(abiPath, "utf8"));
const iface = new ethers.Interface(abi);

// 解析 JSON 数组参数
function parseArray(str) {
  try {
    return JSON.parse(str);
  } catch (err) {
    console.error("数组参数必须是 JSON 格式，比如 \"[1,2]\"");
    process.exit(1);
  }
}

function toBigIntArray(arr) {
  return arr.map((x,i) => {
    return typeof x === "string" ? BigInt(x) : x;
  });
}

const arrP_A = toBigIntArray(parseArray(pA));
const arrP_B0 = toBigIntArray(parseArray(pB0));
const arrP_B1 = toBigIntArray(parseArray(pB1));
const arrP_C = toBigIntArray(parseArray(pC));
const arrPubSignals = toBigIntArray(parseArray(pubSignals));

// 编码
const calldata = iface.encodeFunctionData("submitProof", [
  BigInt(challengeId), // 确保是 uint256
  arrP_A,
  arrP_B0,
  arrP_B1,
  arrP_C,
  arrPubSignals
]);

process.stdout.write(calldata);
