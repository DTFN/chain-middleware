#include "ls_chain.h"
#include "account/account.h"
#include "contract/contract.h"
LsChain::LsChain(std::shared_ptr<Sdk> sdk, int32_t ledger_id, const std::string& root_path) {
    std::cout << "init ls chain" << std::endl;
    AccountManager::Instance().InitAccount(sdk, root_path);
    LsContract::Instance().Init(sdk, ledger_id);
}

LsChain::~LsChain() {}

bool LsChain::Deploy(
    const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) {
    (void)contract_type;
    try {
        std::string addr = LsContract::Instance().Deploy(contract_name, params);
        addresses_[contract_name] = addr;
    } catch (const std::exception& e) {
        std::cout << "lschain deploy failed, exception:" << e.what() << std::endl;
        return false;
    }
    return true;
}
bool LsChain::Call(const std::string contract_type, const std::string contract_name,
    const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
    std::string& output) {
    (void)contract_type;
    try {
        LsContract::Instance().Call(contract_name, contract_address, contract_method, params, output);
    } catch (const std::exception& e) {
        std::cout << "lschain call failed, exception:" << e.what() << std::endl;
        return false;
    }
    return true;
}

std::string LsChain::GetAddress(const std::string& contract_name) {
    auto it = addresses_.find(contract_name);
    if (it == addresses_.end()) {
        std::cout << "can not find the contract address" << std::endl;
        return "";
    } else {
        // std::cout << "the contract address is " << it->second << std::endl;
        return it->second;
    }
}
