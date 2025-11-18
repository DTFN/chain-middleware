#pragma once

#include <string>
#include <boost/date_time/posix_time/posix_time.hpp>
#include <boost/filesystem.hpp>
#include <boost/format.hpp>
#include <boost/make_unique.hpp>
#include <openssl/pem.h>
#include <openssl/rsa.h>
#include "cppsdk/cppsdk.h"
#include "cppsdk/crypto.h"
#include "cppsdk/error.h"

namespace fs = boost::filesystem;

using namespace lscsdk;

struct Account {
    std::string store_dir = "account";
    std::string file_format = "pem";
    std::string address;
    int ssl_type = config::kSslTypeSsl;
    std::shared_ptr<devcrypto::KeyPair> key_pair;
};

class AccountManager {
public:
    static AccountManager& Instance() {
        static AccountManager account;
        return account;
    }
    ~AccountManager() {}

    void InitAccount(std::shared_ptr<Sdk> sdk, const fs::path& root_path);

    std::shared_ptr<Account> current_account() { return current_account_; }

private:
    AccountManager() {}
    bool createAccount(const fs::path& root_path);
    std::string getPrivateKey(EVP_PKEY* pkey);
    bool writekey(const std::string& filename, EVP_PKEY* pkey);
    std::string readPrivateKey(const std::string& file);

    std::shared_ptr<Account> current_account_ = nullptr;
    std::map<std::string, std::shared_ptr<Account>> accounts_;
    std::shared_ptr<Sdk> sdk_;
};
