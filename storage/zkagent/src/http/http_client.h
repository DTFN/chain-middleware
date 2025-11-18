#pragma once

#include <curl/curl.h>
#include <string>

class HttpClient {
public:
    HttpClient(const std::string& ip, uint16_t port);
    ~HttpClient() {}
    std::string Post(const std::string& path, const std::string& body);
    std::string Get(const std::string& path);

private:
    std::string ip_;
    uint16_t port_;
};