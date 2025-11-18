import { Interface } from "ethers";
import fs from "fs";

// 用法: node call.js abi.json callConfig.json
const abiPath = process.argv[2];
const config_json = process.argv[3];

// 读取 ABI
const abiJson = JSON.parse(fs.readFileSync(abiPath, "utf-8"));
const abi = abiJson.abi ?? abiJson; // 有的编译产物包了一层 { "abi": [...] }

// 读取调用参数配置
const config = JSON.parse(config_json);

// 创建接口
const iface = new Interface(abi);

// 编码函数调用
const data = iface.encodeFunctionData("storeFile", [config.input]);

// 输出，不带换行
process.stdout.write(data);
