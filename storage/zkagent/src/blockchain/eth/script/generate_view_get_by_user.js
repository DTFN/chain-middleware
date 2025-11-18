// encode_getChallengesByUser.js
import { ethers } from "ethers";
import fs from "fs";

const abiPath = process.argv[2];
// 读 ABI
const abi = JSON.parse(fs.readFileSync(abiPath, "utf8"));
const iface = new ethers.Interface(abi);

// encode 调用（注意没有参数）
const data = iface.encodeFunctionData("getChallengesByUser", []);

console.log("calldata:", data);
