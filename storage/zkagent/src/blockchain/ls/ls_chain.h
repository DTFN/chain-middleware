#pragma once

#include <iostream>
#include "../../config/config.h"
#include "../blockchain.h"
#include "cppsdk/cppsdk.h"

using namespace lscsdk;

class LsChain : public BlockChain {
public:
    LsChain(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const std::string& root_path);
    ~LsChain();

    virtual bool Deploy(const std::string& contract_type, const std::string& contract_name,
        const std::vector<std::string>& params) override;

    virtual bool Call(const std::string contract_type, const std::string contract_name,
        const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
        std::string& output) override;

    virtual std::string GetAddress(const std::string& contract_name) override;

    virtual std::string GetAccount() override;

private:
    std::map<std::string, std::string> addresses_;
};

inline std::shared_ptr<LsChain> CreateLsChain() {
#ifdef USE_LSSDK
    // contract sdk,ledgerid,rootpath
    auto sdk = std::make_shared<Sdk>();
    std::cout << "Config::Instance().sdk_path() : " << Config::Instance().sdk_path() << std::endl;
    if (!sdk->InitSdk(Config::Instance().sdk_path()) || !sdk->StartSdk()) {
        sdk->CloseSdk();
        throw std::runtime_error("start sdk failed");
    }

    auto block_chain =
        std::make_shared<LsChain>(sdk, Config::Instance().ledger_id(), Config::Instance().account_path());
    return block_chain;
#else
    return nullptr;
#endif
}