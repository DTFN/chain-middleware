#include "http_client.h"
#include <iostream>
#include <curl/curl.h>
#include "../config/config.h"
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>

// 回调函数：将数据写入 string
size_t write_callback(void* contents, size_t size, size_t nmemb, void* userp) {
    size_t totalSize = size * nmemb;
    std::string* response = static_cast<std::string*>(userp);
    response->append(static_cast<char*>(contents), totalSize);
    return totalSize;
}

HttpClient::HttpClient(const std::string& ip, uint16_t port)
    : ip_(ip)
    , port_(port) {}

std::string HttpClient::Post(const std::string& path, const std::string& body) {
    CURL* curl = curl_easy_init();
    if (curl) {
        std::string response;
        struct curl_slist* headers = nullptr;
        headers = curl_slist_append(headers, "Content-Type: application/json");

        std::string url = std::string("http://") + ip_ + ":" + std::to_string(port_) + path;
        std::cout << "url: " << url << " path: " << path << " body: " << body << std::endl;

        curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        // 设置写入回调
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

        CURLcode res = curl_easy_perform(curl);
        if (res != CURLE_OK) {
            std::cerr << "Request failed: " << curl_easy_strerror(res) << std::endl;
        }

        std::cout << "response: " << response << std::endl;
        // 这里解析response，获取挑战的合约地址

        // 创建 Document 并解析
        rapidjson::Document doc;
        if (doc.Parse(response.c_str()).HasParseError()) {
            std::cerr << "JSON parse error!" << std::endl;
            return "";
        }

        // 检查字段是否存在并获取
        if (doc.HasMember("challenge_address") && doc["challenge_address"].IsString()) {
            std::string challenge_address = doc["challenge_address"].GetString();
            std::cout << "Challenge address: " << challenge_address << std::endl;
            Config::Instance().set_challenge_address(challenge_address);
        } else {
            std::cerr << "Missing or invalid 'challenge_address'" << std::endl;
            throw std::runtime_error("Missing or invalid challenge_address");
        }

        curl_slist_free_all(headers);
        curl_easy_cleanup(curl);
    }
    return "";
}

std::string HttpClient::Get(const std::string& path) { return ""; }
