#include <chrono>
#include <iostream>
#include <set>
#include <string>
#include <thread>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "blockchain/blockchain.h"
#include "blockchain/blockchain_factory.h"
#include "config/config.h"
#include "http/http_client.h"
#include "ipfs/ipfs_helper.h"
#include "poseidon/poseidon.h"
#include "proof/zk_proof.h"
#include "tools/tools.h"
#include "tools/utils.h"

// 读取文件内容为 string
std::string ReadFileToString(const std::string& filename) {
    std::ifstream ifs(filename);
    if (!ifs.is_open()) {
        throw std::runtime_error("Cannot open file: " + filename);
    }
    std::stringstream buffer;
    buffer << ifs.rdbuf();
    return buffer.str();
}

uint64_t parseUint64(const std::string& str) {
    try {
        return std::stoull(str);
    } catch (const std::exception& e) {
        throw std::runtime_error(std::string("Exception: ") + e.what() + " Invalid uint64 string: " + str);
    }
}

std::string GetJsonString(rapidjson::Document& doc) {
    // 输出 JSON 字符串
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    doc.Accept(writer);

    return buffer.GetString();
}

bool generate_proof_params(const std::string& challenge_id, const std::string& json_file,
    const std::string& public_json, std::vector<std::string>& params) {
    try {
        params.push_back(challenge_id);

        std::string jsonContent = ReadFileToString(json_file);
        rapidjson::Document doc;
        auto& alloc = doc.GetAllocator();

        rapidjson::Document proof_doc;
        if (proof_doc.Parse(jsonContent.c_str()).HasParseError()) {
            std::cerr << "JSON parse error!" << std::endl;
            return false;
        }

        rapidjson::Document a_arr(rapidjson::kArrayType);
        if (proof_doc.HasMember("pi_a") && proof_doc["pi_a"].IsArray()) {
            a_arr.PushBack(proof_doc["pi_a"][0], alloc);
            a_arr.PushBack(proof_doc["pi_a"][1], alloc);
            params.push_back(GetJsonString(a_arr));
        } else {
            throw std::runtime_error("Invalid JSON format.");
        }

        rapidjson::Document b_arr0(rapidjson::kArrayType);
        rapidjson::Document b_arr1(rapidjson::kArrayType);
        if (proof_doc.HasMember("pi_b") && proof_doc["pi_b"].IsArray()) {
            b_arr0.PushBack(proof_doc["pi_b"][0][1], alloc);
            b_arr0.PushBack(proof_doc["pi_b"][0][0], alloc);
            b_arr1.PushBack(proof_doc["pi_b"][1][1], alloc);
            b_arr1.PushBack(proof_doc["pi_b"][1][0], alloc);
            params.push_back(GetJsonString(b_arr0));
            params.push_back(GetJsonString(b_arr1));
        } else {
            throw std::runtime_error("Invalid JSON format.");
        }

        rapidjson::Document c_arr(rapidjson::kArrayType);
        if (proof_doc.HasMember("pi_c") && proof_doc["pi_c"].IsArray()) {
            c_arr.PushBack(proof_doc["pi_c"][0], alloc);
            c_arr.PushBack(proof_doc["pi_c"][1], alloc);
            params.push_back(GetJsonString(c_arr));
        } else {
            throw std::runtime_error("Invalid JSON format.");
        }

        std::string public_json_content = ReadFileToString(public_json);

        rapidjson::Document public_doc;
        if (public_doc.Parse(public_json_content.c_str()).HasParseError()) {
            std::cerr << "JSON parse error!" << std::endl;
            return false;
        }

        rapidjson::Document public_arr(rapidjson::kArrayType);
        if (public_doc.IsArray()) {
            public_arr.PushBack(public_doc[0], alloc);
            public_arr.PushBack(public_doc[1], alloc);
            public_arr.PushBack(public_doc[2], alloc);
            params.push_back(GetJsonString(public_arr));
        } else {
            throw std::runtime_error("Invalid JSON format.");
        }

        for (auto params_item : params) {
            std::cout << "params_item: " << params_item << std::endl;
        }
    } catch (const std::exception& ex) {
        std::cerr << "Exception: " << ex.what() << std::endl;
        return false;
    }

    return true;
}

void generate_input_json(const std::string& js_file, const std::string& input_file_path, const std::string& nonce,
    const std::string& out_put_file_path) {
    std::cout << "generate_input_json" << std::endl;
    // 生成input文件
    std::string cmd = "node " + js_file + " " + input_file_path + " " + nonce + " " + out_put_file_path;
    std::cout << "cmd: " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    // 检查input.json文件是否存在
}

void generate_witness(const std::string& js_file, const std::string& wasm_file_path, const std::string& input_json_file,
    const std::string& out_put_file_path) {
    std::cout << "generate_witness" << std::endl;
    // 生成witness.wtns
    // node ProveChunkWithRand_js/generate_witness.js ProveChunkWithRand_js/ProveChunkWithRand.wasm input.json
    // witness.wtns
    std::string cmd = "node " + js_file + " " + wasm_file_path + " " + input_json_file + " " + out_put_file_path;
    std::cout << "cmd: " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    // 检查witness.wtns文件
}

void generate_proof_json(const std::string& circuit_final_file, const std::string& witness_file_path,
    const std::string& proof_json_file, const std::string& public_json_file) {
    std::cout << "generate_proof_json" << std::endl;
    // 生成proof.json和public.json
    // snarkjs groth16 prove circuit_final.zkey witness.wtns proof.json public.json
    std::string cmd = "snarkjs groth16 prove " + circuit_final_file + " " + witness_file_path + " " + proof_json_file +
                      " " + public_json_file;
    std::cout << "cmd: " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    // 检查proof.json和public.json文件是否存在
}

bool verifyProof(
    const std::string& verification_key_file, const std::string& public_json_file, const std::string& proof_json_file) {
    // 验证证明 输出[INFO]  snarkJS: OK!
    // snarkjs groth16 verify verification_key.json public.json proof.json
    std::cout << "verify proof..." << std::endl;
    std::string cmd =
        "snarkjs groth16 verify " + verification_key_file + " " + proof_json_file + " " + public_json_file;
    std::cout << "cmd: " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }

    std::cout << "result: " << result << std::endl;
    return result.find("OK") != std::string::npos;
}

void generate_proof(const std::string& generate_input_js_path, const std::string& input_path, const std::string& nonce,
    const std::string& output_path, const std::string generate_witness_js_path, const std::string& wasm_file_path,
    const std::string& circuit_final_file, const std::string& verification_key_file) {
    std::string output_file_path = output_path + "/input.json";
    std::string witness_file_path = output_path + "/witness.wtns";
    std::string proof_json_file = output_path + "/proof.json";
    std::string public_json_file = output_path + "/public.json";
    generate_input_json(generate_input_js_path, input_path, nonce, output_file_path);

    generate_witness(generate_witness_js_path, wasm_file_path, output_file_path, witness_file_path);

    generate_proof_json(circuit_final_file, witness_file_path, proof_json_file, public_json_file);

    int ret = verifyProof(verification_key_file, proof_json_file, public_json_file);
    if (!ret) {
        throw std::runtime_error("verify proof failed");
    }

    std::cout << "verify proof success" << std::endl;
}

void download_file(const std::string& ipfs_url, const std::string& cid, const std::string& out_path) {
    std::cout << "download file..." << std::endl;
    if (!downloadFromIPFS(Config::Instance().ipfs_address(), cid, out_path)) {
        throw std::runtime_error("download file failed");
    }
}

std::shared_ptr<BlockChain> block_chain = nullptr;
int main() {
    // 初始化配置文件
    // 初始化日志
    // 初始化区块链sdk
    // 向超级矿工注册account
    Config::Instance().Init("./config.yml");

    try {
        block_chain = CreateBlockChain(GetBlockChainType(Config::Instance().chain_type()));

        HttpClient client(Config::Instance().http_ip(), Config::Instance().http_port());

        client.Post("/register", "{\"account\":\"" + block_chain->GetAccount() + "\", \"ipfs_address\":\"" +
                                     Config::Instance().ipfs_address() + "\"}");
        // std::this_thread::sleep_for(std::chrono::seconds(30));

        while (true) {
            // 获取待证明任务，然后生成证明提交到链上
            std::vector<std::string> params;
            std::string output;
            std::string contract_address = Config::Instance().challenge_address();
            std::cout << "contract_address: " << contract_address << std::endl;
            block_chain->Call("solidity", "Challenge", contract_address, "getChallengesByUser", params, output);
            std::cout << "output: " << output << std::endl;
            if (output == "") {
                std::this_thread::sleep_for(std::chrono::seconds(60));
                continue;
            }

            rapidjson::Document doc;
            if (doc.Parse(output.c_str()).HasParseError()) {
                std::cerr << "JSON parse error\n";
                return 1;
            }

            if (!doc.IsArray()) {
                std::cerr << "JSON is not an array\n";
                return 1;
            }

            if (doc.Empty()) {
                std::cout << "JSON array is empty.\n";
                std::this_thread::sleep_for(std::chrono::seconds(20));
                continue;
            } else {
                std::cout << "JSON array size: " << doc.Size() << std::endl;
            }

            // 逐个取出元素（你可以按类型转换）
            std::string challenge_id = doc[0].GetString();
            if (challenge_id == "0") {
                std::cout << "empty challenge" << std::endl;
                std::this_thread::sleep_for(std::chrono::seconds(20));
                continue;
            }
            std::string challenger = doc[1].GetString();
            std::string fileHash = doc[2].GetString();
            std::string nonce = doc[3].GetString();
            std::string cid = doc[4].GetString();
            std::string timestamp = doc[5].GetString();  // 注意是字符串
            bool proof_submitted = false;
            if (doc[6].IsBool()) {
                proof_submitted = doc[6].GetBool();
            } else {
                proof_submitted = (std::string(doc[6].GetString()) == "true");
            }
            bool proof_valid = false;
            if (doc[7].IsBool()) {
                proof_valid = doc[7].GetBool();
            } else {
                proof_valid = (std::string(doc[7].GetString()) == "true");
            }
            std::string prover = doc[8].GetString();

            std::cout << "challenger: " << challenger << "\n";
            std::cout << "fileHash: " << fileHash << "\n";
            std::cout << "nonce: " << nonce << "\n";
            std::cout << "cid: " << cid << "\n";
            std::cout << "timestamp: " << timestamp << "\n";
            std::cout << "proof_submitted: " << (proof_submitted ? "true" : "false") << "\n";
            std::cout << "proof_valid: " << (proof_valid ? "true" : "false") << "\n";
            std::cout << "prover: " << prover << "\n";

            if ((!block_chain->isSupportSolidity()) && (!proof_submitted)) {
                std::string download_file_path = "temp/input.bin";
                download_file(Config::Instance().ipfs_address(), cid, download_file_path);
                generate_proof("script/generate_input.js", download_file_path, nonce, "temp",
                    "temp/ProveChunkWithRand_js/generate_witness.js",
                    "temp/ProveChunkWithRand_js/ProveChunkWithRand.wasm", "temp/circuit_final.zkey",
                    "temp/verification_key.json");

                std::string proof_cid = uploadFileToIpfs(Config::Instance().ipfs_address(), "temp/proof.json");
                std::string public_cid = uploadFileToIpfs(Config::Instance().ipfs_address(), "temp/public.json");
                std::vector<std::string> params;
                params.push_back(challenge_id);
                params.push_back(proof_cid);
                params.push_back(Config::Instance().ipfs_address());
                params.push_back(proof_cid);

                params.push_back(public_cid);
                params.push_back(Config::Instance().ipfs_address());
                params.push_back(public_cid);

                block_chain->Call(
                    "solidity", "Challenge", Config::Instance().challenge_address(), "submitProof", params, output);
            } else if (!proof_submitted || !proof_valid) {
                std::string download_file_path = "temp/input.bin";
                download_file(Config::Instance().ipfs_address(), cid, download_file_path);
                generate_proof("script/generate_input.js", download_file_path, nonce, "temp",
                    "temp/ProveChunkWithRand_js/generate_witness.js",
                    "temp/ProveChunkWithRand_js/ProveChunkWithRand.wasm", "temp/circuit_final.zkey",
                    "temp/verification_key.json");
                // 提交证明
                std::vector<std::string> params;
                generate_proof_params(challenge_id, "temp/proof.json", "temp/public.json", params);
                std::cout << "params size: " << params.size() << std::endl;
                std::string output;
                block_chain->Call(
                    "solidity", "Challenge", Config::Instance().challenge_address(), "submitProof", params, output);
                std::cout << "output: " << output << std::endl;

                std::vector<std::string> params1;
                params1.push_back(challenge_id);
                std::string output1;
                block_chain->Call("solidity", "Challenge", Config::Instance().challenge_address(), "getChallengeInfo",
                    params1, output1);
                std::cout << "output1: " << output1 << std::endl;
            } else {
                throw std::runtime_error("proof already submitted");
            }

            std::this_thread::sleep_for(std::chrono::seconds(5));
        }

    } catch (std::exception& e) {
        std::cout << e.what() << std::endl;
    }
    return 0;
}
