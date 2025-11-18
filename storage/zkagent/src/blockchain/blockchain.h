#pragma once

#include <memory>
#include <vector>
#include <string>

enum class BlockChainType {
    ls,
    chainmaker,
    bcos,
    fabric,
    eth
};

class BlockChain {
public:
    BlockChain() {}
    virtual ~BlockChain() {}

    virtual bool Deploy(
        const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) = 0;

    virtual bool Call(const std::string contract_type, const std::string contract_name,
        const std::string& contract_address, const std::string& contract_method,
        const std::vector<std::string>& params, std::string& output) = 0;

    virtual std::string GetAddress(const std::string& contract_name) = 0;

    virtual std::string GetAccount() = 0;

    virtual bool isSupportSolidity() { return true; }
};

