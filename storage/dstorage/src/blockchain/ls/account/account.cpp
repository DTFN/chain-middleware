#include "account.h"
#include "../../common.h"

void AccountManager::InitAccount(std::shared_ptr<Sdk> sdk, const fs::path& root_path) {
#ifdef USE_LSSDK
    sdk_ = sdk;
    accounts_.clear();
    auto account_dir = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? "ecdsa" : "sm";
    auto crypto_type = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? devcrypto::CryptoType::kCryptoTypeIsca
                                                                              : devcrypto::CryptoType::kCryptoTypeSm;
    auto accounts_path = root_path / account_dir;
    std::cout << "accounts path: " << accounts_path << std::endl;
    if (fs::exists(accounts_path) && fs::is_directory(accounts_path)) {
        for (const auto& entry : fs::directory_iterator(accounts_path)) {
            const auto& current_path = entry.path();
            auto file_name = current_path.filename().string();
            auto addr = file_name.substr(0, file_name.find("."));
            if (accounts_.find(addr) != accounts_.end() ||
                (current_account_ != nullptr && addr == current_account_->address)) {
                continue;
            }
            auto tmp_account = std::make_shared<Account>();
            tmp_account->address = addr;
            tmp_account->ssl_type = sdk_->config()->crypto.ssl_type;
            auto private_key_path = current_path.parent_path() / (addr + ".pem");
            if (fs::exists(private_key_path) && fs::is_regular_file(private_key_path)) {
                tmp_account->key_pair =
                    std::move(devcrypto::CreateKeyPair(crypto_type, readPrivateKey(private_key_path.string())));
            } else {
                throw std::runtime_error("get private key file failed, path: " + private_key_path.string());
            }
            if (current_account_ == nullptr) {
                current_account_ = tmp_account;
            } else {
                accounts_.emplace(tmp_account->address, tmp_account);
            }
        }
    }

    if (accounts_.size() == 0 && current_account_ == nullptr) {
        std::cout << "***************************no account found, create one" << std::endl;
        createAccount(root_path);
    }
#endif
}

std::string AccountManager::getPrivateKey(EVP_PKEY* pkey) {
#ifdef USE_LSSDK
    auto ec_key = EVP_PKEY_get1_EC_KEY(pkey);
    if (ec_key == nullptr) {
        return {};
    }
    std::vector<uint8_t> buffer(32);
    auto private_key = EC_KEY_get0_private_key(ec_key);
    if (private_key != nullptr) {
        auto length = BN_bn2bin(private_key, buffer.data());
        buffer.resize(length);
    }
    EC_KEY_free(ec_key);
    return ToHex(buffer, "");
#else
    return "";
#endif
}

bool AccountManager::writekey(const std::string& filename, EVP_PKEY* pkey) {
#ifdef USE_LSSDK
    BIO* out = BIO_new_file(filename.c_str(), "w");
    if (out == nullptr) {
        return false;
    }
    auto ret = PEM_write_bio_PrivateKey(out, pkey, nullptr, nullptr, 0, nullptr, nullptr);
    BIO_free(out);
    if (ret <= 0) {
        return false;
    }
    out = BIO_new_file((filename + ".pub").c_str(), "w");
    if (out == nullptr) {
        return false;
    }
    ret = PEM_write_bio_PUBKEY(out, pkey);
    BIO_free(out);
    if (ret <= 0) {
        return false;
    }
    return true;
#else
    return false;
#endif
}

bool AccountManager::createAccount(const fs::path& root_path) {
#ifdef USE_LSSDK
    auto kEcParamSecp256k1 = R"(-----BEGIN EC PARAMETERS-----
BgUrgQQACg==
-----END EC PARAMETERS-----)";
    auto kEcParamSm2p256v1 = R"(-----BEGIN EC PARAMETERS-----
BggqgRzPVQGCLQ==
-----END EC PARAMETERS-----)";
    auto crypto_type = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? devcrypto::CryptoType::kCryptoTypeIsca
                                                                              : devcrypto::CryptoType::kCryptoTypeSm;
    auto account_dir = sdk_->config()->crypto.ssl_type == config::kSslTypeSsl ? "ecdsa" : "sm";

    BIO* bio = BIO_new_mem_buf(
        crypto_type == devcrypto::CryptoType::kCryptoTypeIsca ? kEcParamSecp256k1 : kEcParamSm2p256v1, -1);
    if (bio == nullptr) {
        std::cerr << "BIO_new_mem_buf error" << std::endl;
        return false;
    }
    EVP_PKEY* pkey = PEM_read_bio_Parameters(bio, nullptr);
    BIO_free(bio);
    bio = nullptr;
    if (pkey == nullptr) {
        std::cerr << "PEM_read_bio_Parameters error" << std::endl;
        return false;
    }
    EVP_PKEY_CTX* ctx = EVP_PKEY_CTX_new(pkey, nullptr);
    EVP_PKEY_free(pkey);
    pkey = nullptr;
    if (ctx == nullptr) {
        std::cerr << "EVP_PKEY_CTX_new error" << std::endl;
        return false;
    }
    if (EVP_PKEY_keygen_init(ctx) <= 0) {
        EVP_PKEY_CTX_free(ctx);
        std::cerr << "EVP_PKEY_keygen_init error" << std::endl;
        return false;
    }
    if (EVP_PKEY_keygen(ctx, &pkey) <= 0) {
        EVP_PKEY_CTX_free(ctx);
        std::cerr << "EVP_PKEY_keygen error" << std::endl;
        return false;
    }

    auto keypair = devcrypto::CreateKeyPair(crypto_type, getPrivateKey(pkey));
    if (keypair == nullptr) {
        EVP_PKEY_free(pkey);
        EVP_PKEY_CTX_free(ctx);
        return false;
    }
    auto address = "0x" + devcrypto::GetAddress(*keypair);
    fs::create_directories(root_path / account_dir);
    auto prikey_path = root_path / account_dir / (address + ".pem");
    writekey(prikey_path.string(), pkey);
    std::cout << address << std::endl;

    EVP_PKEY_free(pkey);
    EVP_PKEY_CTX_free(ctx);

    auto tmp_account = std::make_shared<Account>();
    tmp_account->key_pair = std::move(keypair);
    tmp_account->ssl_type = sdk_->config()->crypto.ssl_type;
    tmp_account->address = address;
    if (current_account_ == nullptr) {
        current_account_ = tmp_account;
    }
    return true;
#else
    return false;
#endif
}

std::string AccountManager::readPrivateKey(const std::string& file) {
#ifdef USE_LSSDK
    BIO* bio = BIO_new_file(file.c_str(), "r");
    if (bio == nullptr) {
        std::cerr << "BIO_new_file error" << std::endl;
        return {};
    }

    EVP_PKEY* pkey = PEM_read_bio_PrivateKey(bio, NULL, NULL, NULL);
    BIO_free(bio);
    if (pkey == nullptr) {
        std::cerr << "PEM_read_bio_PrivateKey error" << std::endl;
        return {};
    }

    EC_KEY* ec_key = EVP_PKEY_get1_EC_KEY(pkey);
    if (ec_key == nullptr) {
        EVP_PKEY_free(pkey);
        std::cerr << "EVP_PKEY_get1_EC_KEY error" << std::endl;
        return {};
    }

    // 获取私钥
    const BIGNUM* private_key = EC_KEY_get0_private_key(ec_key);
    if (!private_key) {
        EC_KEY_free(ec_key);
        EVP_PKEY_free(pkey);
        std::cerr << "EC_KEY_get0_private_key error" << std::endl;
        return {};
    }

    // 将私钥转换为16进制字符串
    std::vector<uint8_t> buffer(32);
    size_t length = BN_bn2bin(private_key, buffer.data());
    buffer.resize(length);

    EC_KEY_free(ec_key);
    EVP_PKEY_free(pkey);

    return ToHex(buffer, "");
#else
    return "";
#endif
}