#pragma once

#include <iostream>
#include "../../config/config.h"
#include "../blockchain.h"
#include "cppsdk/cppsdk.h"

using namespace lscsdk;

class ChainmakerChain : public BlockChain {
public:
    ChainmakerChain(std::string go_sdk_path) : sdk_path_(go_sdk_path) {}
    ~ChainmakerChain() {}

    virtual bool Deploy(const std::string& contract_type, const std::string& contract_name,
        const std::vector<std::string>& params) override;

    virtual bool Call(const std::string contract_type, const std::string contract_name,
        const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
        std::string& output) override;

    virtual std::string GetAddress(const std::string& contract_name) override;

private:
    std::string sdk_path_;
    std::map<std::string, std::string> addresses_;
};

inline std::shared_ptr<ChainmakerChain> CreateChainmakerChain() {
    std::string sdk_path = Config::Instance().sdk_path();

    auto block_chain = std::make_shared<ChainmakerChain>(sdk_path);
    return block_chain;
}