import { ethers } from "ethers";
import fs from "fs";

// example
// {
//     "bytecodeFile": "./../eth/solidity/FileStore.bin",
//     "types": ["address", "uint256", "string"],
//     "values": ["0x1234567890abcdef1234567890abcdef12345678", 42, "hello world"]
// }
  

// 用法: node deploy.js deployConfig.json
const configFile = process.argv[2];
const config = JSON.parse(configFile);

// 读取合约 bytecode
let bytecode = fs.readFileSync(config.bytecodeFile, "utf-8").trim();

// 编码构造函数参数
const encodedArgs = ethers.AbiCoder.defaultAbiCoder().encode(
  config.types,
  config.values
);

// 拼接部署 data
const deployData = bytecode + encodedArgs.slice(2);

process.stdout.write(deployData);
