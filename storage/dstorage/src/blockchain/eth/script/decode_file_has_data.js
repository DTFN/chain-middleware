import { ethers } from "ethers";
import fs from "fs";

// 命令行参数：ABI文件路径 + 返回值hex
// 用法：node decode_hasFile.js ./artifacts/FileStore.json 0x1234abcd...
const [abiPath, responseHex] = process.argv.slice(2);

if (!abiPath || !responseHex) {
  console.error("Usage: node decode_hasFile.js <abiPath> <responseHex>");
  process.exit(1);
}

// 1. 读取 ABI 文件
const abiJson = JSON.parse(fs.readFileSync(abiPath, "utf-8"));
const iface = new ethers.Interface(abiJson);

// 2. 解码函数返回值
const decoded = iface.decodeFunctionResult("hasFile", responseHex);

// 4. 输出 JSON
process.stdout.write(JSON.stringify(decoded, null, 2));
