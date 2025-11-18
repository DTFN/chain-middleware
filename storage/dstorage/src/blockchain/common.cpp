#include "common.h"
#include <chrono>
#include <abi/abi_codec.h>
#include <abi/hash.h>
#include <boost/filesystem.hpp>
#include <boost/format.hpp>
#include <cppsdk/transaction_factory.h>
#include "cppsdk/error.h"
#include "tools/tools.h"

namespace fs = boost::filesystem;
using namespace abicodec;
int64_t GetSysTimestamp() {
    auto now = std::chrono::steady_clock::now();
    return std::chrono::duration_cast<std::chrono::milliseconds>(now.time_since_epoch()).count();
}

std::string FormatJson(const std::string& json_str) {
    rapidjson::Document doc;
    doc.Parse(json_str.c_str());
    if (doc.HasParseError()) {
        throw std::runtime_error("format json failed, json: " + json_str);
    }
    rapidjson::StringBuffer buffer;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(buffer);
    if (doc.HasMember("result")) {
        doc["result"].Accept(writer);
    } else if (doc.HasMember("error")) {
        doc["error"].Accept(writer);
    } else {
        doc.Accept(writer);
    }

    std::string formatted_json = buffer.GetString();
    return formatted_json;
}
std::string GetAbi(const std::string& contract_type, const std::string& contract_name) {
    // fs::path root_path = "./contracts";
    // fs::path contract_path = root_path / contract_type / (contract_name + ".sol");
    // if (fs::exists(contract_path) && fs::is_regular_file(contract_path)) {
    //     std::string command = (boost::format("%1% --abi --overwrite -o ./ %2%") % GetSolcPath() %
    //     contract_path).str(); if (system(command.c_str()) != 0) {
    //         throw std::runtime_error("create abi file had failed");
    //     }
    //     std::string abi = ReadFile("./" + contract_name + ".abi");
    //     fs::remove("./" + contract_name + ".abi");
    //     return abi;
    // } else {
    //     throw std::runtime_error("can not find the contract");
    // }
    std::string abi = ReadFile("./contracts/compiled/1/" + contract_name + "/" + contract_name + ".abi");
    return abi;
}

void ParseReceipt(std::shared_ptr<Sdk> sdk, const std::string& receipt, const std::string& contract_type,
    const std::string& contract_name, const std::string& contract_method, std::string& contract_addr,
    std::string& decoded_output, const std::string& abi) {
#ifdef USE_LSSDK
    // std::cout << "receipt:" << receipt << std::endl;
    // std::cout << "----------------------------------------------------------" << std::endl;
    rapidjson::Document doc;
    doc.Parse(std::string(receipt.begin(), receipt.end()).c_str());
    if (doc.HasParseError()) {
        std::cout << "parse receipt failed\n";
        return;
    }
    std::string status;
    if (doc.HasMember("result") && doc["result"].IsObject()) {
        auto result = doc["result"].GetObject();

        if (result.HasMember("contractAddress") && result["contractAddress"].IsString()) {
            auto contract_address = result["contractAddress"].GetString();
            std::cout << "contract address: " << contract_address << std::endl;
        }

        if (result.HasMember("status") && result["status"].IsString()) {
            status = result["status"].GetString();
            auto it = status_map.find(status);
            // std::cout << "transaction status: " << status << " ("
            //           << (it != status_map.end() ? it->second : "unknown status") << ")" << std::endl;
        }

        if (result.HasMember("transactionHash") && result["transactionHash"].IsString()) {
            auto transaction_hash = result["transactionHash"].GetString();
            std::cout << "transaction hash: " << transaction_hash << std::endl;
        }

        if (result.HasMember("from") && result["from"].IsString()) {
            auto from = result["from"].GetString();
            std::cout << "current account: " << from << std::endl;
        }

        if (result.HasMember("contractAddress") && result["contractAddress"].IsString()) {
            contract_addr = result["contractAddress"].GetString();
        }

        if (((contract_type.empty() || contract_name.empty()) && abi.empty()) || contract_method.empty()) {
            return;
        }
        std::string tmp_abi = abi;
        if (tmp_abi.empty()) {
            tmp_abi = GetAbi(contract_type, contract_name);
        }
        abicodec::AbiFactory factory;
        auto contract_abi = factory.CreateAbi(tmp_abi);
        auto hash_type = sdk->config()->crypto.ssl_type == config::kSslTypeSsl ? abicodec::hash::HashType::kKeccak256
                                                                               : abicodec::hash::HashType::kSM3;
        abicodec::AbiCodec codec(hash_type);
        if (result.HasMember("output") && result["output"].IsString()) {
            auto output = std::make_shared<bytes>(FromHex(result["output"].GetString()));
            // std::cout << "----------------------------------------------------------" << std::endl;
            if (status == "0x0") {
                auto method = contract_abi->GetMethodByMethodName(contract_method);
                std::string ret = "(";
                for (const auto& output : method->outputs()) {
                    ret += output->GetTypeAsString();
                    ret += ",";
                }
                if (!ret.empty() && ret.back() == ',') {
                    ret.pop_back();
                }
                ret.push_back(')');
                // std::cout << "return value size: " << method->outputs().size() << std::endl;
                // std::cout << "return type: " << ret << std::endl;
                decoded_output = codec.DecodeMethodOutput(method, output);
                // std::cout << "return values: " << decoded_output << std::endl;
            } else if (status == "0x16") {
                if (output->empty()) {
                    decoded_output = "[]";
                } else {
                    decoded_output = codec.DecodeMethodInputByMethodSig("Error(string)", output);
                }
                std::cout << "revert reason: " << decoded_output << std::endl;
            } else {
                if (output->empty()) {
                    decoded_output = "[]";
                } else {
                    decoded_output = codec.DecodeMethodInputByMethodSig("Error(string)", output);
                }
                // std::cout << "return value size: 1" << std::endl;
                // std::cout << "return type: (string)" << std::endl;
                // std::cout << "return values: " << decoded_output << std::endl;
            }
        }

        if (result.HasMember("logs") && result["logs"].IsArray()) {
            auto logs = result["logs"].GetArray();
            for (const auto& log : logs) {
                if (log.IsObject()) {
                    auto log_entry = log.GetObject();
                    if (log_entry.HasMember("topics") && log_entry["topics"].IsArray() && log_entry.HasMember("data") &&
                        log_entry["data"].IsString()) {
                        auto topics = log_entry["topics"].GetArray();
                        std::string data = log_entry["data"].GetString();
                        if (topics.Empty()) {
                            throw std::runtime_error("topics is null");
                        }
                        if (data.empty()) {
                            throw std::runtime_error("data is null");
                        }
                        if (topics[0].IsString()) {
                            std::string event_topic = topics[0].GetString();
                            std::shared_ptr<abicodec::ContractMethod> e = nullptr;
                            auto events = contract_abi->name_to_events();
                            for (auto it = events.begin(); it != events.end() && e == nullptr; ++it) {
                                for (const auto& event : it->second) {
                                    if (event->GetEventTopic(hash_type) == ToHex(FromHex(event_topic), "")) {
                                        e = event;
                                        break;
                                    }
                                }
                            }
                            if (e == nullptr) {
                                return;
                            }

                            auto decoded_data = codec.DecodeEvent(e, std::make_shared<bytes>(FromHex(data)));
                            std::cout << "----------------------------------------------------------" << std::endl;
                            std::cout << "Logs topics:\n";
                            for (const auto& item : topics) {
                                std::cout << "    " << item.GetString() << "\n";
                            }
                            std::cout << "event: "
                                      << (boost::format(R"({"%1%" : %2%})") % e->name() % decoded_data).str()
                                      << std::endl;
                        }
                    }
                }
            }
        }
    } else if (doc.HasMember("error") && doc["error"].IsObject()) {
        const auto& error = doc["error"];
        std::cout << "RPC request failed\n"
                  << "----------------------------------------------------------\n"
                  << "error code: " << error["code"].GetInt() << "\nerror message: " << error["message"].GetString()
                  << std::endl;
    }
#endif
}

std::string GenerateExtraData(bool is_create, const std::string& contract_type) {
    (void)contract_type;
    std::string prefix = "2";
    if (is_create) {
        prefix = "1";
    }
    std::string suffix = "0";
    return prefix + "#" + suffix;
}

void SendTx(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const devcrypto::KeyPair& key_pair,
    const std::string& byte_code, std::string& contract_addr, std::string& decoded_output, bool is_create,
    const std::string& contract_type, const std::string& contract_name, const std::string& contract_method,
    const std::string& abi) {
#ifdef USE_LSSDK
    auto factory = std::make_shared<tx::TransactionFactory>();
    std::string block_limit = "1000";
    std::string extra_data = GenerateExtraData(is_create, contract_type);
    auto ret = factory->CreateSignedTransaction(
        key_pair, "", "", block_limit, contract_addr, byte_code, "1", std::to_string(ledger_id), extra_data);
    auto tx_str = ret.second;
    int64_t send_time = GetSysTimestamp();
    std::promise<int> promise;
    std::future<int> future = promise.get_future();
    sdk->json_rpc()->SendTxAndGetReceipt(
        ledger_id, tx_str, false,
        [&](Error err, std::vector<uint8_t> data) {
            if (err.errcode_ != 0) {
                std::cout << err.errmsg_ << "\n";
            } else {
                std::string res = std::string(data.begin(), data.end());
                // std::cout << FormatJson(res) << "\n";
                try {
                    ParseReceipt(
                        sdk, res, contract_type, contract_name, contract_method, contract_addr, decoded_output, abi);
                    std::cout << "consume time: " << GetSysTimestamp() - send_time << "ms\n";
                } catch (const std::exception& e) {
                    std::cout << e.what() << std::endl;
                }
            }
            promise.set_value(0);
        },
        10000);
    future.get();
#else
    return;
#endif
}

void CallTx(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const std::string& from, const std::string& to,
    const std::string& data, const std::string& vm_type, const std::string& contract_type,
    const std::string& contract_name, const std::string& contract_method, std::string& decoded_output,
    const std::string& abi) {
    rapidjson::StringBuffer buf;
    rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(buf);
    writer.StartObject();
    writer.Key("from");
    writer.String(from.c_str());

    writer.Key("to");
    writer.String(to.c_str());

    writer.Key("vmType");
    writer.String(vm_type.c_str());

    writer.Key("data");
    writer.String(data.c_str());
    writer.EndObject();

    auto json = buf.GetString();

    std::promise<int> promise;
    std::future<int> future = promise.get_future();
    sdk->json_rpc()->Call(ledger_id, json, [&](Error err, bytes data) {
        if (err.errcode_ != 0) {
            std::cout << err.errmsg_ << "\n";
        } else {
            auto res = std::string(data.begin(), data.end());
            // std::cout << FormatJson(res) << "\n";
            try {
                std::string tmp;
                ParseReceipt(sdk, res, contract_type, contract_name, contract_method, tmp, decoded_output, abi);
            } catch (const std::exception& e) {
                std::cout << e.what() << std::endl;
            }
        }
        promise.set_value(0);
    });
    future.get();
}