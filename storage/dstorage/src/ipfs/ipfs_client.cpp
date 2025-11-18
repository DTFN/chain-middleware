#include "ipfs_client.h"
#include <fstream>
#include <iostream>
#include <sstream>
#include <curl/curl.h>
#include <rapidjson/document.h>

bool IpfsClient::is_curl_init_ = false;
std::mutex IpfsClient::mutex_;

// 回调函数：把返回的数据写入 string（不要在这里解析 JSON！）
size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp) {
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

// 回调函数：把响应写入 string
size_t WriteToFile(void* ptr, size_t size, size_t nmemb, void* userdata) {
    std::ofstream* ofs = static_cast<std::ofstream*>(userdata);
    ofs->write(static_cast<char*>(ptr), size * nmemb);
    return size * nmemb;
}

IpfsClient::IpfsClient(const std::string& api_url)
    : api_url_(api_url) {}

std::string IpfsClient::UploadFile(const std::string& file_path) { return sendFileToIpfs(file_path); }

void IpfsClient::DownloadFile(const std::string& cid, const std::string& file_name) {
    fetchFileFromIpfs(cid, file_name);
}

size_t WriteToString(void* ptr, size_t size, size_t nmemb, void* userdata) {
    std::string* resp = reinterpret_cast<std::string*>(userdata);
    resp->append(static_cast<char*>(ptr), size * nmemb);
    return size * nmemb;
}

bool IpfsClient::IsExist(const std::string& cid) {
    std::string url = api_url_ + "/api/v0/block/stat?arg=" + cid;

    if (!is_curl_init_) {
        std::lock_guard<std::mutex> lock(mutex_);
        if (!is_curl_init_) {
            curl_global_init(CURL_GLOBAL_DEFAULT);
            is_curl_init_ = true;
        }
    }

    curl_global_init(CURL_GLOBAL_DEFAULT);
    CURL* curl = curl_easy_init();
    if (!curl) return false;

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_POST, 1L);  // GET 请求

    // 写入内存而不是文件
    std::string response;
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteToString);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

    CURLcode res = curl_easy_perform(curl);

    long http_code = 0;
    curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

    curl_easy_cleanup(curl);

    if (res == CURLE_OK && http_code == 200) {
        return true;
    } else {
        return false;
    }
}

std::string IpfsClient::GetFileCID(const std::string& file_path) { return sendFileToIpfs(file_path); }

std::string IpfsClient::sendFileToIpfs(const std::string& file_path) {
    std::string cid_res = "";
    CURL* curl;
    CURLcode res;
    std::string response;

    if (!is_curl_init_) {
        std::lock_guard<std::mutex> lock(mutex_);
        if (!is_curl_init_) {
            curl_global_init(CURL_GLOBAL_DEFAULT);
            is_curl_init_ = true;
        }
    }
    curl = curl_easy_init();
    if (curl) {
        struct curl_httppost* formpost = nullptr;
        struct curl_httppost* lastptr = nullptr;

        curl_formadd(&formpost, &lastptr, CURLFORM_COPYNAME, "file", CURLFORM_FILE, file_path.c_str(),
            CURLFORM_CONTENTTYPE, "application/octet-stream", CURLFORM_END);

        std::string real_url = api_url_ + "/api/v0/add";
        curl_easy_setopt(curl, CURLOPT_URL, real_url.c_str());
        curl_easy_setopt(curl, CURLOPT_HTTPPOST, formpost);

        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

        // 输出日志
        // curl_easy_setopt(curl, CURLOPT_VERBOSE, 1L);

        res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            std::cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << std::endl;

        } else {
            rapidjson::Document doc;
            if (doc.Parse(response.c_str()).HasParseError()) {
                std::cerr << "JSON parse error" << std::endl;
                return "";
            }

            if (doc.HasMember("Hash") && doc["Hash"].IsString()) {
                cid_res = doc["Hash"].GetString();
            } else {
                std::cerr << "No Hash field or not a string" << std::endl;
            }
        }

        curl_easy_cleanup(curl);
        curl_formfree(formpost);
    }
    curl_global_cleanup();
    return cid_res;
}

void IpfsClient::fetchFileFromIpfs(const std::string& cid, const std::string& outPath) {
    if (cid == "") {
        return;
    }
    std::string url = api_url_ + "/api/v0/cat?arg=" + cid;

    std::cout << "url : " << url << std::endl;

    if (!is_curl_init_) {
        std::lock_guard<std::mutex> lock(mutex_);
        if (!is_curl_init_) {
            curl_global_init(CURL_GLOBAL_DEFAULT);
            is_curl_init_ = true;
        }
    }
    CURL* curl = curl_easy_init();
    if (!curl) return;

    std::ofstream ofs(outPath, std::ios::binary);
    if (!ofs.is_open()) {
        std::cerr << "Failed to open output file: " << outPath << std::endl;
        curl_easy_cleanup(curl);
        return;
    }

    curl_easy_setopt(curl, CURLOPT_POST, 1L);

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteToFile);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &ofs);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);  // 允许重定向

    CURLcode res = curl_easy_perform(curl);
    curl_easy_cleanup(curl);
    ofs.close();

    return;
}
