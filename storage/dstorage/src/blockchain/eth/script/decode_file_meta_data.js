import { ethers } from "ethers";
import fs from "fs";
import { decode } from "punycode";

const abiPath = process.argv[2];
// 1. 读取 ABI
const abi = JSON.parse(fs.readFileSync(abiPath, "utf-8"));

// 2. 创建 Interface
const iface = new ethers.Interface(abi);

// 3. 你拿到的 hex 返回值
const response = process.argv[3];

// 4. 解码 getFileMeta 返回值
const decoded = iface.decodeFunctionResult("getFileMeta", response);

// 5. 构造数组形式的 JSON
const jsonArray = [
      decoded.fileId,
      decoded.fileName,
      decoded.fileSize.toString(),
      decoded.totalShards.toString(),
      decoded.dataShards.toString(),
      decoded.ipfsUrls,
      decoded.shardHashes,
      decoded.poseidonHashes,
      decoded.timestamp.toString(),
      decoded.author
  ];

process.stdout.write(JSON.stringify(jsonArray, null, 2));
