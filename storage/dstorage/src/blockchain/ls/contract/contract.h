#pragma once

#include "cppsdk/cppsdk.h"

using namespace lscsdk;

class LsContract {
public:
    static LsContract& Instance() {
        static LsContract contract;
        return contract;
    }
    ~LsContract() {}

    void Init(std::shared_ptr<Sdk> sdk, int32_t ledger_id) {
        sdk_ = sdk;
        ledger_id_ = ledger_id;
    }

    std::string Deploy(const std::string& contract_name, const std::vector<std::string>& params);
    void Call(const std::string contract_name, const std::string& contract_address, const std::string& contract_method,
        const std::vector<std::string>& params, std::string& output);

private:
    LsContract() {}

    std::shared_ptr<Sdk> sdk_;
    int32_t ledger_id_;
};