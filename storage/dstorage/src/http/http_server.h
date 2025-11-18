#pragma once

#include <cstdint>
#include <map>
#include <string>
#include "boost/noncopyable.hpp"

class HttpServer : public boost::noncopyable {
public:
    static HttpServer& Instance() {
        static HttpServer http_server;
        return http_server;
    }
    ~HttpServer() {}

    void Init(const std::string& listen_ip, uint16_t listen_port, const std::string& doc_root);

    std::string GetAccount(const std::string& ipfs_address);

    void AddAccount(const std::string& ipfs_address, const std::string& account);

private:
    HttpServer() {}

    std::map<std::string, std::string> account_map_;

    std::map<std::string, std::string> file_name_map_;
};
