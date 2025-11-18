#include "zk_proof.h"
#include <iostream>


bool submitProofOnChain(const std::string& cid, const std::string& proofFile, const std::string& publicFile) {
    // 可先写入 log 或发送给本地 relay
    std::cout << "Simulated submit proof for " << cid << "\n";
    return true;
}


bool generateProof() {
    int ret = system("node generate_witness.js proof.wasm input.json witness.wtns");
    if (ret != 0) return false;

    ret = system("snarkjs groth16 prove circuit.zkey witness.wtns proof.json public.json");
    return ret == 0;
}
