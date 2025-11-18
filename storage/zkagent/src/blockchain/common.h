#pragma once
#include <future>
#include <iostream>
#include <map>
#include <memory>
#include <vector>
#include <boost/algorithm/string.hpp>
#include <boost/filesystem.hpp>
#include <boost/format.hpp>
#include <rapidjson/document.h>
#include <rapidjson/prettywriter.h>
#include <rapidjson/stringbuffer.h>
#include "cppsdk/cppsdk.h"
#include "cppsdk/crypto.h"
#include <sstream>


#if defined(WIN32)
#define SOLC_PLATFORM "win"
#define SOLC_PROGRAM "solc.exe"
#else
#define SOLC_PLATFORM "linux"
#define SOLC_PROGRAM "solc"
#endif
#undef GetObject

using namespace lscsdk;

const static std::map<std::string, std::string> status_map = {
    {"0x0", "SUCCESS"},
    {"0x1", "UNKNOWN"},
    {"0x2", "BAD_RLP"},
    {"0x3", "DATA_INVALID"},
    {"0x4", "OUT_OF_GAS_INTRINSIC"},
    {"0x5", "SIGNATURE_INVALID"},
    {"0x7", "BALANCE_NOT_ENOUGH"},
    {"0x8", "OUT_OF_GAS_BASE"},
    {"0x9", "CONTRACT_VALIDATION_FAILURE"},
    {"0xa", "BAD_INSTRUCTION"},
    {"0xb", "BAD_JUMP_DESTINATION"},
    {"0xc", "OUT_OF_GAS"},
    {"0xd", "STACK_OVERFLOW"},
    {"0xe", "STACK_UNDERFLOW"},
    {"0xf", "NONCE_ERROR"},
    {"0x10", "BLOCK_LIMIT_EXCEEDED"},
    {"0x15", "PRECOMPILED_ERROR"},
    {"0x16", "REVERT"},
    {"0x18", "ADDRESS_ALREADY_USED"},
    {"0x19", "NO_LEDGER_PERMISSION"},
    {"0x1a", "ADDRESS_NOT_EXIST"},
    {"0x1b", "GAS_OVERFLOW"},
    {"0x1c", "TXPOOL_FULL"},
    {"0x1d", "TRANSACTION_REFUSED"},
    {"0x1e", "CONTRACT_FROZEN"},
    {"0x1f", "ACCOUNT_FROZEN"},
    {"0x20", "PERMISSION_DENIED"},
    {"0x28", "WASM_TRAP"},
    {"0x29", "WASM_UNREACHABLE"},
    {"0x2a", "DOCKER_TIMEOUT"},
    {"0x2b", "DOCKER_CONNECTION_EXCEPTION"},
    {"0x2c", "DOCKER_UNSUPPORTED"},
    {"0x2d", "BAD_ARGUMENT"},
    {"0x2e", "MOVE_RESOURCE_ALREADY_EXISTS"},
    {"0x2f", "MOVE_RESOURCE_NOT_EXIST"},
    {"0x2710", "TRANSACTION_DUPLICATED"},
    {"0x2711", "TRANSACTION_DROPPED"},
    {"0x2712", "NODE_ID_INVALID"},
    {"0x2713", "LEDGER_ID_INVALID"},
    {"0x2714", "REQUEST_TO_WRONG_LEDGER"},
    {"0x2715", "BAD_TRANSACTION"},
    {"0x2716", "LEDGER_MEMORY_LIMIT_EXCEEDED"},
};

inline std::string GetSolcPath(bool sm = false) {
    std::string path =
        (boost::format("./solcs/solc-%1%/%2%/" SOLC_PLATFORM "/solc/" SOLC_PROGRAM) % "0.8.11" % (sm ? "sm" : "ecdsa"))
            .str();
#ifdef WIN32
    boost::algorithm::replace_all(path, "/", "\\");
#elif defined(__linux__)
    int ret = system(std::string("chmod +x ").append(path).c_str());
    if (ret != 0) {
        std::cerr << "Warning: chmod failed with return code " << ret << std::endl;
    }
#endif

    std::cout << "solc path is : " << path << std::endl;
    return path;
}

inline void compileSolidity(const std::string& contract_name, const std::string& file_name) {
    std::string contract_path = "./contracts/solidity/" + file_name + ".sol";
    // 生成国密版本字节码
    std::string command = (boost::format("%1% --bin --overwrite -o ./contracts/compiled/1/%2% %3%") %
                           GetSolcPath(true) % contract_name % contract_path)
                              .str();
    if (system(command.c_str()) != 0) {
        throw std::runtime_error("create abi and bin file had failed");
    }
    // 重命名字节码文件名
    std::string bin_prefix = "./contracts/compiled/1/" + contract_name + "/" + contract_name;
    boost::filesystem::rename(bin_prefix + ".bin", bin_prefix + ".sm.bin");
    // 生成国际版本字节码
    command = (boost::format("%1% --bin --abi --overwrite -o ./contracts/compiled/1/%2% %3%") % GetSolcPath(false) %
               contract_name % contract_path)
                  .str();
    if (system(command.c_str()) != 0) {
        throw std::runtime_error("create abi and bin file had failed");
    }
}

int64_t GetSysTimestamp();
std::string FormatJson(const std::string& json_str);
std::string GetAbi(const std::string& contract_type, const std::string& contract_name);
void ParseReceipt(std::shared_ptr<Sdk> sdk, const std::string& receipt, const std::string& contract_type,
    const std::string& contract_name, const std::string& contract_method, std::string& contract_addr,
    std::string& decoded_output, const std::string& abi);
std::string GenerateExtraData(bool is_create, const std::string& contract_type);
void SendTx(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const devcrypto::KeyPair& key_pair,
    const std::string& byte_code, std::string& contract_addr, std::string& decoded_output, bool is_create = false,
    const std::string& contract_type = "solidity", const std::string& contract_name = "",
    const std::string& contract_method = "", const std::string& abi = "");
void CallTx(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const std::string& from, const std::string& to,
    const std::string& data, const std::string& vm_type, const std::string& contract_type,
    const std::string& contract_name, const std::string& contract_method, std::string& decoded_output,
    const std::string& abi = "");

inline bool IsAddress(const std::string& addr) {
    // 必须以 "0x" 开头
    if (addr.size() != 42 || addr[0] != '0' || addr[1] != 'x') {
        return false;
    }

    // 检查后续 40 个字符是否都是 hex
    for (size_t i = 2; i < addr.size(); i++) {
        char c = addr[i];
        if (!std::isxdigit(static_cast<unsigned char>(c))) {
            return false;
        }
    }

    return true;
}

template <typename T>
typename std::enable_if<std::is_same<T, int32_t>::value, T>::type GetJsonValue(const rapidjson::Value& value) {
    return value.GetInt();
}

template <typename T>
typename std::enable_if<std::is_same<T, std::string>::value, T>::type GetJsonValue(const rapidjson::Value& value) {
    return value.GetString();
}

template <class T>
std::pair<std::string, std::set<T>> formatJsonAndGetResultArray(const std::string& json_str) {
    rapidjson::Document doc;
    doc.Parse(json_str.c_str());
    if (doc.HasParseError()) {
        throw std::runtime_error("format json failed, json: " + json_str);
    }
    std::set<T> res;
    if (doc.HasMember("result") && doc["result"].IsArray()) {
        auto result = doc["result"].GetArray();
        for (const auto& item : result) {
            res.insert(GetJsonValue<T>(item));
        }
    } else {
        std::cout << "receipt don't have result field" << std::endl;
        return std::make_pair("", std::set<T>());
    }
    rapidjson::StringBuffer buffer;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(buffer);

    doc["result"].Accept(writer);
    std::string formatted_json = buffer.GetString();
    return std::make_pair(formatted_json, res);
}

inline std::string GenerateJsonParam(const std::vector<std::string>& params) {
    if (params.size() == 1) {
        rapidjson::Document single;
        single.Parse(params[0].c_str());
        if (!single.HasParseError() && (single.IsObject() || single.IsArray())) {
            // 唯一参数且为 JSON 对象或数组，直接返回
            return params[0];
        }
    }

    // 多参数或不是合法 JSON，包装成数组
    rapidjson::StringBuffer buf;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(buf);
    writer.StartArray();
    for (const auto& p : params) {
        rapidjson::Document d;
        d.Parse(p.c_str());
        if (!d.HasParseError() && d.IsArray()) {
            d.Accept(writer);  // 嵌套 JSON 数组
        } else {
            writer.String(p.c_str());  // 普通字符串
        }
    }
    writer.EndArray();
    return buf.GetString();
}


inline uint64_t HexToUint64(const std::string& hexStr) {
    uint64_t value = 0;
    std::stringstream ss;
    if (hexStr.rfind("0x", 0) == 0) {
        ss << std::hex << hexStr.substr(2);  // 去掉 0x
    } else {
        ss << std::hex << hexStr;
    }
    ss >> value;
    return value;
}