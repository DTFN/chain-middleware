import { Interface } from "ethers";
import fs from "fs";

// 用法: node encode.js abi.json config_json
const abiPath = process.argv[2];
const config_json = process.argv[3];

// 读取 ABI
const abiJson = JSON.parse(fs.readFileSync(abiPath, "utf-8"));
const abi = abiJson.abi ?? abiJson;
const iface = new Interface(abi);

// 读取参数
const config = JSON.parse(config_json);

// 编码函数调用
const data = iface.encodeFunctionData("createChallenge", [
  config.fileHash,
  config.nonce,
  config.cid,
  config.prover
]);

process.stdout.write(data);
