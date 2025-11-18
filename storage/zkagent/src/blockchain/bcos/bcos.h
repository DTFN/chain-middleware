#pragma once

#include <condition_variable>
#include <iostream>
#include <bcos-c-sdk/bcos_sdk_c.h>
#include <bcos-c-sdk/bcos_sdk_c_uti_tx.h>
#include "../../config/config.h"
#include "../blockchain.h"
#include "bcos-c-sdk/bcos_sdk_c_error.h"
#include "bcos-c-sdk/bcos_sdk_c_rpc.h"
#include "bcos-c-sdk/bcos_sdk_c_uti_abi.h"
#include "bcos-c-sdk/bcos_sdk_c_uti_keypair.h"

class BcosChain : public BlockChain {
public:
    BcosChain(const std::string& config_file_path, int32_t ledger_id, const std::string& group_id);
    ~BcosChain();

    virtual bool Deploy(const std::string& contract_type, const std::string& contract_name,
        const std::vector<std::string>& params) override;

    virtual bool Call(const std::string contract_type, const std::string contract_name,
        const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
        std::string& output) override;

    virtual std::string GetAddress(const std::string& contract_name) override;

    virtual std::string GetAccount() override;
    void SetTempContractName(const std::string& contract_name) { temp_contract_name_ = contract_name; }
    std::string GetTempContractName() { return temp_contract_name_; }
    void BindContractAddress(const std::string& contract_name, const std::string& contract_address) {
        addresses_[contract_name] = contract_address;
    }

    void ReceiveRespData(const std::string& data) { receive_resp_data_ = data; }
    std::string GetReceiveRespData() {
        auto temp = receive_resp_data_;
        receive_resp_data_ = "";
        return temp;
    }

    void Notify() {
        {
            std::lock_guard<std::mutex> lock(mtx_);
            ready_ = true;
        }
        cv_.notify_one();  // 通知等待的线程
    }
    void Wait() {
        std::unique_lock<std::mutex> lock(mtx_);
        cv_.wait(lock, [&] { return ready_; });  // 条件满足才会返回
        ready_ = false;                          // 下次还能阻塞
    }

private:
    std::map<std::string, std::string> addresses_;
    std::string config_path_;
    int32_t ledger_id_;
    std::string group_id_;
    void* sdk_;
    void* key_pair_;
    std::string account_address_;
    int sdk_crypto_type_ = 0;
    std::string temp_contract_name_;

    std::condition_variable cv_;
    std::mutex mtx_;
    bool ready_;

    std::string receive_resp_data_;
};

inline std::shared_ptr<BcosChain> CreateBcosChain() {
    auto block_chain = std::make_shared<BcosChain>(
        Config::Instance().sdk_path(), Config::Instance().ledger_id(), Config::Instance().group_id());
    return block_chain;
}