#include "eth.h"
#include <abi/abi_codec.h>
#include <abi/hash.h>
#include <curl/curl.h>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../common.h"
#include "eth_tool.h"
#include "libchainmaker_sdk.h"
#include "rapidjson/prettywriter.h"
#include "tools/tools.h"

using namespace rapidjson;

// libcurl 写回调
static size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp) {
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

// 发送 RPC 请求
std::string rpcRequest(const std::string& rpc_url, const std::string& data) {
    CURL* curl = curl_easy_init();
    std::string readBuffer;
    if (curl) {
        struct curl_slist* headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_URL, rpc_url.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &readBuffer);
        curl_easy_perform(curl);
        curl_easy_cleanup(curl);
    }
    return readBuffer;
}

bool EthChain::Deploy(
    const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) {
    std::cout << "Deploying " << contract_name << "..." << std::endl;
    (void)params;
    std::string rpc_url = Config::Instance().eth_url();
    std::string from = Config::Instance().eth_account();

    std::string abi = "eth/solidity/" + contract_name + ".abi";
    std::string bin = "eth/solidity/" + contract_name + ".bin";

    std::string bytecode = "";
    if (contract_name == "Challenge") {
        if (params.size() != 1) {
            return false;
        }

        // 创建一个 JSON Document，作为对象
        rapidjson::Document doc;
        doc.SetObject();

        // 分配器
        rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();

        // 添加 string 字段: bytecodeFile
        doc.AddMember("bytecodeFile", rapidjson::Value(bin.c_str(), allocator), allocator);

        // 添加 array: types
        rapidjson::Value types(rapidjson::kArrayType);
        types.PushBack("address", allocator);
        doc.AddMember("types", types, allocator);

        // 添加 array: values
        rapidjson::Value values(rapidjson::kArrayType);
        values.PushBack(rapidjson::Value(params[0].c_str(), allocator), allocator);  // string
        // values.PushBack("0xcf7ed3acca5a467e9e704c703e8d87f634fb0fc9", allocator); // string
        doc.AddMember("values", values, allocator);

        // 转换成字符串
        rapidjson::StringBuffer buffer;
        rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
        doc.Accept(writer);

        std::string exec_json = buffer.GetString();
        std::cout << "exec_json: " << exec_json << std::endl;
        bytecode = ExecNode("eth/script/generate_abi.js", exec_json);
        if (bytecode.substr(0, 2) != "0x") {
            bytecode = "0x" + bytecode;
        }
        std::cout << "bytecode: " << bytecode << std::endl;
    } else {
        bytecode = ReadFile(bin);  // 注意添加 0x 前缀
        if (bytecode.substr(0, 2) != "0x") {
            bytecode = "0x" + bytecode;
        }
    }

    // 读取 ABI 和 bytecode
    std::string abiStr = ReadFile(abi);

    // 使用 RapidJSON 构造部署请求
    Document deployDoc;
    deployDoc.SetObject();
    Document::AllocatorType& allocator = deployDoc.GetAllocator();
    deployDoc.AddMember("jsonrpc", "2.0", allocator);
    deployDoc.AddMember("method", "eth_sendTransaction", allocator);
    Value json_params(kArrayType);
    Value tx(kObjectType);
    tx.AddMember("from", Value().SetString(from.c_str(), allocator), allocator);

    tx.AddMember("data", Value().SetString(bytecode.c_str(), allocator), allocator);
    // tx.AddMember("data", Value().SetString(bytecode.c_str(), allocator), allocator);
    json_params.PushBack(tx, allocator);
    deployDoc.AddMember("params", json_params, allocator);
    deployDoc.AddMember("id", 1, allocator);
    // 转成字符串发送
    StringBuffer buffer;
    Writer<StringBuffer> writer(buffer);
    deployDoc.Accept(writer);
    std::string deployResp = rpcRequest(rpc_url, buffer.GetString());
    std::cout << "Deploy TX response: " << deployResp << std::endl;
    if (deployResp == "") {
        return false;
    }

    Document dTxRes;
    dTxRes.Parse(deployResp.c_str());

    if (dTxRes.HasMember("result") && dTxRes["result"].IsString()) {
        std::string txHash = dTxRes["result"].GetString();
        std::cout << "tx hash: " << txHash << std::endl;

        // 轮询等交易上链（简单做法：while 循环查询）
        std::string contractAddress;
        while (true) {
            std::string receiptReq = R"({
                "jsonrpc":"2.0",
                "method":"eth_getTransactionReceipt",
                "params":[")" + txHash +
                                     R"("],
                "id":1
            })";

            std::string receiptResp = rpcRequest(rpc_url, receiptReq);

            // 如果还没打包，result 会是 null
            if (receiptResp.find("\"result\":null") != std::string::npos) {
                std::this_thread::sleep_for(std::chrono::seconds(1));
                continue;
            }

            // 解析 JSON，拿到 contractAddress
            Document d;
            d.Parse(receiptResp.c_str());
            if (d.HasMember("result")) {
                const Value& result = d["result"];
                if (result.HasMember("contractAddress")) {
                    contractAddress = result["contractAddress"].GetString();
                    break;
                }
            }
        }

        std::cout << "Contract deployed at: " << contractAddress << std::endl;
        addresses_[contract_name] = contractAddress;

    } else {
        std::cout << "no result field or not string" << std::endl;
        return false;
    }

    return true;
}

bool EthChain::Call(const std::string contract_type, const std::string contract_name,
    const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
    std::string& output) {
    std::cout << "Calling " << contract_name << " " << contract_method << " ..." << std::endl;
    std::string rpc_url = Config::Instance().eth_url();
    std::string from = Config::Instance().eth_account();
    // 使用 eth_getTransactionReceipt 获取 contractAddress
    std::string contractAddress = contract_address;  // 从 receipt 里获取

    std::string abi = "eth/solidity/" + contract_name + ".abi";
    std::string bin = "eth/solidity/" + contract_name + ".bin";

    // 调用 storeFile
    std::string dataStoreFile = "";
    if (contract_method == "storeFile") {
        dataStoreFile = ExecNode("eth/script/generate_call.js " + abi, params[0]);
    } else if (contract_method == "createChallenge") {
        Document d;
        d.SetObject();
        Document::AllocatorType& allocator = d.GetAllocator();

        // 按顺序放到 JSON
        d.AddMember("fileHash", Value(params[0].c_str(), allocator), allocator);
        d.AddMember("nonce", Value(params[1].c_str(), allocator), allocator);
        d.AddMember("cid", Value(params[2].c_str(), allocator), allocator);
        d.AddMember("prover", Value(params[3].c_str(), allocator), allocator);

        // 转换成字符串
        StringBuffer buffer;
        Writer<StringBuffer> writer(buffer);
        d.Accept(writer);

        std::string json_str = buffer.GetString();
        dataStoreFile = ExecNode("eth/script/generate_call_create_challenge_abi.js " + abi, json_str);
    } else if (contract_method == "getChallengeInfo") {
        dataStoreFile = ExecNode("eth/script/generate_view_get_info.js " + abi, params[0]);
    } else if (contract_method == "getChallengesByUser") {
        dataStoreFile = ExecNode("eth/script/generate_view_get_by_user.js", abi);
    } else if (contract_method == "getFileMeta") {
        dataStoreFile = ExecNode("eth/script/generate_file_meta.js " + abi, params[0], params[1]);
    } else if (contract_method == "hasFile") {
        dataStoreFile = ExecNode(
            "eth/script/generate_has_file.js " + abi, params[0], params[1], params[2], params[3]);
    }
    std::cout << "dataStoreFile: " << dataStoreFile << std::endl;

    Document callDoc;
    callDoc.SetObject();
    Document::AllocatorType& alloc2 = callDoc.GetAllocator();
    callDoc.AddMember("jsonrpc", "2.0", alloc2);
    if (contract_method == "createChallenge" || contract_method == "storeFile" || contract_method == "submitProof") {
        callDoc.AddMember("method", "eth_sendTransaction", alloc2);
    } else {
        callDoc.AddMember("method", "eth_call", alloc2);
    }

    Value callParams(kArrayType);
    Value callTx(kObjectType);
    callTx.AddMember("from", Value().SetString(from.c_str(), alloc2), alloc2);
    callTx.AddMember("to", Value().SetString(contract_address.c_str(), alloc2), alloc2);
    callTx.AddMember("data", Value().SetString(dataStoreFile.c_str(), alloc2), alloc2);
    callParams.PushBack(callTx, alloc2);
    callDoc.AddMember("params", callParams, alloc2);
    callDoc.AddMember("id", 2, alloc2);

    StringBuffer buffer2;
    Writer<StringBuffer> writer2(buffer2);
    callDoc.Accept(writer2);

    std::string callResp = rpcRequest(rpc_url, buffer2.GetString());
    // std::cout << "storeFile TX response: " << callResp << std::endl;

    if (callResp.empty()) {
        return true;
    }

    Document dTxRes;
    dTxRes.Parse(callResp.c_str());
    if (dTxRes.HasMember("result") && dTxRes["result"].IsString()) {
        if (contract_method == "getChallengeInfo") {
            // 0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb922660000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000000000000000000000000000000000000068ad937400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000165734f8847cff904c6dce849292330090fecffb000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000
            std::string result = dTxRes["result"].GetString();
            output = getChallengeInfo(result);
            return true;
        } else if (contract_method == "getChallengesByUser") {
            std::string result = dTxRes["result"].GetString();
            output = getChallengeByUser(result);
            return true;
        } else if (contract_method == "submitProof") {
            return true;
        } else if (contract_method == "getFileMeta") {
            std::string result = dTxRes["result"].GetString();
            output = ExecNode("eth/script/decode_file_meta_data.js " + abi, result);
            return true;
        } else if (contract_method == "hasFile") {
            std::string result = dTxRes["result"].GetString();
            output = ExecNode("eth/script/decode_file_has_data.js " + abi, result);
            return true;
            // output = parseBoolResult(result);
            return true;
        }
        std::string txHash = dTxRes["result"].GetString();
        std::cout << "tx hash: " << txHash << std::endl;

        // 轮询等交易上链（简单做法：while 循环查询）
        std::string contractAddress;
        while (true) {
            std::string receiptReq = R"({
                "jsonrpc":"2.0",
                "method":"eth_getTransactionReceipt",
                "params":[")" + txHash +
                                     R"("],
                "id":2
            })";

            std::string receiptResp = rpcRequest(rpc_url, receiptReq);
            // 如果还没打包，result 会是 null
            if (receiptResp.find("\"result\":null") != std::string::npos) {
                std::this_thread::sleep_for(std::chrono::seconds(1));
                continue;
            }

            std::cout << "receipt response: " << receiptResp << std::endl;
            Document receipt_doc;
            receipt_doc.Parse(receiptResp.c_str());

            if (contract_method == "createChallenge") {
                if (!receipt_doc.HasParseError() && receipt_doc.HasMember("result")) {
                    const Value& result = receipt_doc["result"];
                    if (result.HasMember("logs") && result["logs"].IsArray()) {
                        const Value& logs = result["logs"];
                        if (!logs.Empty()) {
                            const Value& firstLog = logs[0];
                            if (firstLog.HasMember("topics") && firstLog["topics"].IsArray()) {
                                const Value& topics = firstLog["topics"];
                                if (topics.Size() > 1 && topics[1].IsString()) {
                                    std::string hexChallengeId = topics[1].GetString();
                                    uint64_t challenge_id = HexToUint64(hexChallengeId);
                                    std::cout << "challengeId = " << challenge_id << std::endl;
                                    // 创建 JSON 文档
                                    Document challenge_doc;
                                    challenge_doc.SetArray();  // 将 doc 设置为数组类型

                                    Document::AllocatorType& allocator = challenge_doc.GetAllocator();
                                    challenge_doc.PushBack(challenge_id, allocator);  // 将 uint64_t 添加到数组中

                                    // 写成字符串
                                    StringBuffer buffer;
                                    Writer<StringBuffer> writer(buffer);
                                    challenge_doc.Accept(writer);

                                    output = buffer.GetString();
                                }
                            }
                        }
                    }
                }
            }
            break;
        }
    } else {
        std::cout << "no result field or not string" << std::endl;
        return false;
    }
    return true;
}

std::string EthChain::GetAddress(const std::string& contract_name) { return addresses_[contract_name]; }