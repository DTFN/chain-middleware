const fs = require("fs");
const circomlibjs = require("circomlibjs");


const FIELD_SIZE = 31;              // 每个 Poseidon field 最大安全字节数
const MAX_FILE_BYTES = 4096;        // 最大文件大小
const MAX_FIELD_COUNT = 133

function uint8ArrayToBigInt(arr) {
    let hex = Buffer.from(arr).toString("hex");
    return BigInt("0x" + hex);
  }

// 将 buffer 分块转为 BigInt field 数组
function bufferToBigIntArray(buffer) {

    const padded = Buffer.alloc(FIELD_SIZE * MAX_FIELD_COUNT); // 预先分配空间并补0
    buffer.copy(padded);  // 复制内容，自动零填充剩余部分
    
    const chunkSize = 31;
    const result = [];
    for (let i = 0; i < padded.length; i += chunkSize) {
        const chunk = padded.slice(i, i + chunkSize);
        result.push(BigInt("0x" + chunk.toString("hex")));
    }
    return result;
}

// 两两组合 Poseidon 递归哈希
function reducePoseidonArray(arr, poseidon) {
    
    while (arr.length > 1) {
        const next = [];
        for (let i = 0; i < arr.length; i += 2) {
            if (i + 1 < arr.length) {
                const p = poseidon([arr[i], arr[i + 1]]);
                // console.log(p);
                next.push(p);
            } else {
                next.push(arr[i]);
            }
        }
        arr = next;
    }
    return arr[0]; // 单个 field 元素
}

async function main() {
    const filePath = process.argv[2];

    if (!filePath) {
        console.error("用法: node generate_input.js <file_path>");
        process.exit(1);
    }

    const fileBuffer = fs.readFileSync(filePath);
    const limitedBuffer = fileBuffer.slice(0, MAX_FILE_BYTES);

    const poseidon = await circomlibjs.buildPoseidon();
    const F = poseidon.F;

    const fileFields = bufferToBigIntArray(limitedBuffer);

    const expected_chunk_hash = reducePoseidonArray(fileFields, poseidon);

    process.stdout.write(F.toObject(expected_chunk_hash).toString());
}

main();
