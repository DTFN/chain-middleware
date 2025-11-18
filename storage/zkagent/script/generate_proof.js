const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const CIRCUIT_NAME = "ProveChunkHashMatch";
const CIRCUIT_DIR = path.join(__dirname, "..", "circuits");
const PROOF_DIR = path.join(__dirname, "proof");

// 文件路径配置
const INPUT_JSON = path.join(__dirname, "input.json");
const WASM_FILE = path.join(CIRCUIT_DIR, `${CIRCUIT_NAME}.wasm`);
const ZKEY_FILE = path.join(CIRCUIT_DIR, `${CIRCUIT_NAME}.zkey`);
const WITNESS_FILE = path.join(PROOF_DIR, "witness.wtns");
const PROOF_FILE = path.join(PROOF_DIR, "proof.json");
const PUBLIC_FILE = path.join(PROOF_DIR, "public.json");

// 创建 proof 目录（如果不存在）
if (!fs.existsSync(PROOF_DIR)) {
    fs.mkdirSync(PROOF_DIR);
}

// 执行 shell 命令并打印输出
function run(cmd) {
    console.log(`> ${cmd}`);
    execSync(cmd, { stdio: "inherit" });
}

function main() {
    if (!fs.existsSync(INPUT_JSON)) {
        console.error("❌ input.json 不存在，请先运行 gen_input.js");
        process.exit(1);
    }

    console.log("📦 Step 1: 生成 witness.wtns");
    run(`snarkjs wtns calculate ${WASM_FILE} ${INPUT_JSON} ${WITNESS_FILE}`);

    console.log("🔐 Step 2: 生成 proof.json 和 public.json");
    run(`snarkjs groth16 prove ${ZKEY_FILE} ${WITNESS_FILE} ${PROOF_FILE} ${PUBLIC_FILE}`);

    console.log("✅ ZK 证明已生成：");
    console.log("  - proof.json");
    console.log("  - public.json");
}

main();
