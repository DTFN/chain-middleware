import { ethers } from "ethers";
import fs from "fs";

// 用法: node query.js abi.json <contractAddress> <fileId>

const abiPath = process.argv[2];
const contractAddress = process.argv[3];
const fileId = process.argv[4];

// 读取 ABI
const abiJson = JSON.parse(fs.readFileSync(abiPath, "utf-8"));
const abi = abiJson.abi ?? abiJson;

// RPC Provider（这里用本地 Hardhat，换成你自己的 RPC 地址）
const provider = new ethers.JsonRpcProvider("http://127.0.0.1:8545");

// 创建合约实例
const contract = new ethers.Contract(contractAddress, abi, provider);

// 调用合约
const fileMeta = await contract.files(fileId);

console.log("📂 FileMeta:", fileMeta);

// 如果想输出更整齐
console.log(JSON.stringify(fileMeta, (key, value) =>
    typeof value === "bigint" ? value.toString() : value,
    2
  ));
