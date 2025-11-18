#include "config.h"
#include <filesystem>
#include <iostream>

const int MIN_LISTEN_PORT = 1024;
const int MAX_LISTEN_PORT = 65535;

bool IsValidPort(int port) {
    if (port <= MIN_LISTEN_PORT || port > MAX_LISTEN_PORT) {
        return false;
    }
    return true;
}

void Config::Init(const std::string& path) {
    if ("" == path) {
        std::cout << "the current config path is empty!" << std::endl;
        BOOST_THROW_EXCEPTION(std::runtime_error("Invalid config file path"));
    }

    if (!std::filesystem::exists(path)) {
        std::cout << "the current config path is not exist!" << std::endl;
        BOOST_THROW_EXCEPTION(std::runtime_error("can not find the file, path = " + path));
    }

    try {
        YAML::Node doc = YAML::LoadFile(path);
        initHttp(doc);
        initChain(doc);
        initIpfs(doc);
        initIsal(doc);
    } catch (const std::exception& e) {
        std::string error_info = std::string("the config file initialize failed! error_info is ") + e.what();
        std::cout << error_info << std::endl;
        BOOST_THROW_EXCEPTION(std::runtime_error(error_info));
    }

    std::cout << "config load complete. the current config path is " + path << std::endl;
}

void Config::initHttp(const YAML::Node& doc) {
    const auto& http_node = doc["http"];
    if (http_node) {
        http_doc_root_ = http_node["doc_root"].as<std::string>("./tmp");
        http_listen_ip_ = http_node["listen_ip"].as<std::string>("0.0.0.0");
        http_listen_port_ = http_node["listen_port"].as<int>(56656);
        if (!IsValidPort(http_listen_port_)) {
            std::string error_info =
                std::string("tquic listen port is invalid! port = ") + std::to_string(http_listen_port_);
            BOOST_THROW_EXCEPTION(std::runtime_error(error_info));
        }
    }

    std::cout << "http config: "
              << "listen_ip: " << http_listen_ip_ << ", "
              << "listen_port: " << http_listen_port_ << std::endl;
}

void Config::initChain(const YAML::Node& doc) {
    const auto& chain_node = doc["chain"];
    if (chain_node) {
        chain_type_ = chain_node["type"].as<std::string>("ls");
        sdk_path_ = chain_node["sdk_path"].as<std::string>("./config.yml");
        account_path_ = chain_node["account_path"].as<std::string>("./account");
        ledger_id_ = chain_node["ledger_id"].as<int>(8080);
        challenge_address_ = chain_node["challange_address"].as<std::string>("");
        msp_id_ = chain_node["msp_id"].as<std::string>("");
        crypto_path_ = chain_node["crypto_path"].as<std::string>("");
        cert_path_ = chain_node["cert_path"].as<std::string>("");
        key_path_ = chain_node["key_path"].as<std::string>("");
        tls_cert_path_ = chain_node["tls_cert_path"].as<std::string>("");
        peer_endpoint_ = chain_node["peer_endpoint"].as<std::string>("");
        gateway_peer_ = chain_node["gateway_peer"].as<std::string>("");
        channel_name_ = chain_node["channel_name"].as<std::string>("");
        chaincode_name_ = chain_node["chaincode_name"].as<std::string>("");
        group_id_ = chain_node["group_id"].as<std::string>("");
        eth_account_ = chain_node["eth_account"].as<std::string>("");
        eth_url_ = chain_node["eth_url"].as<std::string>("");

    }

    std::cout << "chain config: "
              << "chain_type: " << chain_type_ << ", "
              << "sdk_path: " << sdk_path_ << ", "
              << "account_path: " << account_path_ << ", "
              << "ledger_id: " << ledger_id_ << ", "
              << "challenge_address_: " << challenge_address_ << ", "
              << "msp_id" << msp_id_ << ", "
              << "crypto_path: " << crypto_path_ << ", "
              << "peer_endpoint: " << peer_endpoint_ << ", "
              << "gateway_peer: " << gateway_peer_ << ", "
              << "channel_name: " << channel_name_ << ", "
              << "chaincode_name: " << chaincode_name_ << ", "
              << "group_id: " << group_id_ << std::endl;
}

void Config::initIpfs(const YAML::Node& doc) {
    const auto& ipfs_node = doc["ipfs"];
    if (ipfs_node) {
        if (ipfs_node["address"]) {
            for (const auto& addrNode : ipfs_node["address"]) {
                ipfs_address_.push_back(addrNode.as<std::string>());
            }
        }
    }
}

void Config::initIsal(const YAML::Node& doc) {
    const auto& isal_node = doc["isal"];
    if (isal_node) {
        isal_k_ = isal_node["isal_k"].as<uint64_t>(4);
        isal_m_ = isal_node["isal_m"].as<uint64_t>(2);
        isal_threshold_ = isal_node["isal_threshold"].as<uint64_t>(5);
    }
}
