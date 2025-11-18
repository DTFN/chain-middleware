/**
 * @copyright (C) 2023 0shu.
 *
 * @file: config.h
 * @brief: 配置类
 */

#pragma once

#include <atomic>
#include <random>
#include <string>
#include <tuple>
#include <vector>
#include <boost/property_tree/ini_parser.hpp>
#include <boost/property_tree/ptree.hpp>
#include <yaml-cpp/yaml.h>
#include "boost/noncopyable.hpp"
#include "tools/tools.h"

using namespace boost::property_tree;

class Config : public boost::noncopyable {
public:
    static Config& Instance() {
        static Config ins;
        return ins;
    }

    ~Config() {}

    void Init(const std::string& path);

public:
    /// [http]
    const std::string& http_listen_ip() { return http_listen_ip_; }
    uint16_t http_listen_port() { return http_listen_port_; }

    const std::string& http_doc_root() { return http_doc_root_; }

    /// [chain]
    const std::string& chain_type() { return chain_type_; }
    const std::string& sdk_path() { return sdk_path_; }
    const std::string& account_path() { return account_path_; }
    uint32_t ledger_id() { return ledger_id_; }
    const std::string& group_id() { return group_id_; }
    void set_challenge_address(const std::string& address) { challenge_address_ = address; }
    const std::string& challenge_address() { return challenge_address_; }
    const std::string& msp_id() { return msp_id_; }
    const std::string& crypto_path() { return crypto_path_; }
    const std::string& cert_path() { return cert_path_; }
    const std::string& key_path() { return key_path_; }
    const std::string& tls_cert_path() { return tls_cert_path_; }
    const std::string& peer_endpoint() { return peer_endpoint_; }
    const std::string& gateway_peer() { return gateway_peer_; }
    const std::string& channel_name() { return channel_name_; }
    const std::string& chaincode_name() { return chaincode_name_; }

    const std::string& eth_account() { return eth_account_; }
    const std::string& eth_url() { return eth_url_; }

    uint64_t isal_all() { return isal_k_ + isal_m_; }
    uint64_t isal_k() { return isal_k_; }

    uint64_t isal_m() { return isal_m_; }

    uint64_t isal_threshold() { return isal_threshold_; }

    const std::vector<std::string>& ipfs_address() { return ipfs_address_; }
    std::string GetRandomAddress() {
        if (ipfs_address_.size() == 0) {
            return "";
        }

        // 轮询
        static std::vector<std::string>::iterator it = ipfs_address_.begin();
        if (it == ipfs_address_.end()) {
            it = ipfs_address_.begin();
        }

        std::string result = *it;
        it++;
        return result;
    }

private:
    Config() {}
    void initHttp(const YAML::Node& pt);
    void initChain(const YAML::Node& pt);
    void initIpfs(const YAML::Node& pt);
    void initIsal(const YAML::Node& pt);

private:
    /// [http]
    std::string http_listen_ip_ = "0.0.0.0";
    uint16_t http_listen_port_ = 1234;
    std::string http_doc_root_ = "./tmp";

    /// [chain]
    std::string chain_type_ = "ls";
    std::string sdk_path_ = "./ls/config.yml";
    std::string account_path_ = "./ls/account";
    uint32_t ledger_id_ = 1;
    std::string challenge_address_ = "";
    std::string msp_id_ = "";
    std::string crypto_path_ = "./organizations/peerOrganizations/org1.example.com";
    std::string cert_path_ = "";
    std::string key_path_ = "";
    std::string tls_cert_path_ = "";
    std::string peer_endpoint_ = "localhost:7051";
    std::string gateway_peer_ = "peer0.org1.example.com";
    std::string channel_name_ = "mychannel";
    std::string chaincode_name_ = "basic";
    std::string group_id_ = "group0";
    std::string eth_url_ = "http://127.0.0.1:8545";
    std::string eth_account_ = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266";

    /// [ipfs]
    std::vector<std::string> ipfs_address_;

    /// [isal]
    uint64_t isal_k_ = 4;  // 数据块
    uint64_t isal_m_ = 2;  // 校验块
    uint64_t isal_threshold_ = 5;
};
