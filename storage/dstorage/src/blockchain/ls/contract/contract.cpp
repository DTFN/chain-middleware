#include "contract.h"
#include <boost/filesystem.hpp>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "rapidjson/prettywriter.h"

#include "../account/account.h"

#include <string>
#include <abi/abi_codec.h>
#include <abi/hash.h>
#include <boost/date_time/posix_time/posix_time.hpp>

#include <boost/format.hpp>
#include <boost/make_unique.hpp>

#include <cppsdk/transaction_factory.h>
#include <openssl/pem.h>
#include <openssl/rsa.h>
#include "cppsdk/cppsdk.h"
#include "cppsdk/error.h"

#include "../../common.h"
#include "tools/tools.h"

void LsContract::Call(const std::string contract_name, const std::string& address, const std::string& contract_method,
    const std::vector<std::string>& params, std::string& output) {
#ifdef USE_LSSDK
    std::string contract_address = address; //Address(address, Address::kFromHex, Address::kAlignRight).HexPrefixed();
    if (!IsAddress(contract_address)) {
        throw std::invalid_argument(
            "contract address should be '0x' prefixed hex string whose digital length must be within 40 (20 "
            "bytes)");
    }

    std::string abi;
    {
        //处理solidity合约
        // compileSolidity(contract_name, contract_name);
        // std::string dst_path = "./contracts/compiled/1/" + contract_name + "/" + contract_address + "/";
        // std::string src_path = "./contracts/compiled/1/" + contract_name + "/";
        // boost::filesystem::create_directories(dst_path);
        // if (boost::filesystem::exists(src_path) && boost::filesystem::is_directory(src_path)) {
        //     for (const auto& entry : boost::filesystem::directory_iterator(src_path)) {
        //         const auto& current_path = entry.path();
        //         if (boost::filesystem::is_regular_file(current_path)) {
        //             boost::filesystem::rename(current_path, dst_path + current_path.filename().string());
        //         }
        //     }
        // }
        abi = ReadFile(
            "./contracts/compiled/1/" + contract_name + "/" + contract_name + ".abi");
    }

    auto hash_type = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? abicodec::hash::HashType::kKeccak256
                                                                            : abicodec::hash::HashType::kSM3;
    abicodec::AbiCodec codec(hash_type);

    auto json_param = GenerateJsonParam(params);  // buf.GetString();

    // std::cout << "json_param: " << json_param << std::endl;

    std::shared_ptr<abicodec::ContractMethod> method;
    std::shared_ptr<std::vector<uint8_t>> data;
    abicodec::AbiFactory abi_factory;
    auto contract_abi = abi_factory.CreateAbi(abi);
    method = contract_abi->GetMethodByMethodName(contract_method);
    if (method == nullptr) {
        throw std::runtime_error("abicodec GetMethodByMethodName failed");
    }
    data = codec.EncodeMethodInput(method, json_param, abicodec::VmType::kEvm);

    auto key_pair = AccountManager::Instance().current_account()->key_pair;
    if (method->state_mutability() == "view") {
        std::string vm_type = std::to_string((int)abicodec::VmType::kEvm);
        CallTx(sdk_, ledger_id_, AccountManager::Instance().current_account()->address, contract_address, ToHex(*data, ""),
            vm_type, "solidity", contract_name, contract_method, output);
    } else {
        SendTx(sdk_, ledger_id_, *key_pair, ToHex(*data, ""), contract_address, output, false, "solidity", contract_name,
            contract_method);
    }

// std::cout << "call contract output : " << output << std::endl;
#endif
}

std::string LsContract::Deploy(const std::string& contract_name, const std::vector<std::string>& params) {
#ifdef USE_LSSDK
    std::cout << "Deploy contract, contract_name is " << contract_name << " params size is " << params.size()
              << std::endl;
    std::string abi;
    std::string bin;
    {
        compileSolidity(contract_name, contract_name);
        abi = ReadFile("./contracts/compiled/1/" + contract_name + "/" + contract_name + ".abi");
        bin = ReadFile("./contracts/compiled/1/" + contract_name + "/" + contract_name + ".bin");
        if (sdk_->config()->crypto.ssl_type == config::kSslTypeSmSsl) {
            bin = ReadFile("./contracts/compiled/1/" + contract_name + "/" + contract_name + ".sm.bin");
        }
    }

    auto hash_type = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? abicodec::hash::HashType::kKeccak256
                                                                            : abicodec::hash::HashType::kSM3;
    abicodec::AbiCodec codec(hash_type);

    rapidjson::StringBuffer buf;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(buf);
    writer.StartArray();
    for (const auto& p : params) {
        writer.String(p.c_str());
    }
    writer.EndArray();
    auto json_param = buf.GetString();

    std::shared_ptr<std::vector<uint8_t>> data;

    data = codec.EncodeConstructor(abi, bin, json_param, abicodec::VmType::kEvm);

    std::string contract_addr;
    auto key_pair = AccountManager::Instance().current_account()->key_pair;
    std::string output;
    SendTx(sdk_, ledger_id_, *key_pair, ToHex(*data, ""), contract_addr, output, true, "solidity");

    if (contract_addr.empty()) {
        throw std::runtime_error("deploy contract failed, cant not get contract_addr");
    }

    std::cout << "contract address is " << contract_addr << std::endl;
    return contract_addr;
#else
    return "";
#endif
}
