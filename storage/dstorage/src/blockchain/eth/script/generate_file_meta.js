// encode.js
import { ethers } from "ethers";
import fs from "fs";

const args = process.argv.slice(2);

const [abiPath, file_id, pk] = args;

// 加载 ABI
const abi = JSON.parse(fs.readFileSync(abiPath, "utf8"));

// 1. 创建接口对象
const iface = new ethers.Interface(abi);

// 2. 要调用的函数和参数
const data = iface.encodeFunctionData("getFileMeta", [file_id, pk]);

process.stdout.write(data);
