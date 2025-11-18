#include "ipfs_helper.h"
#include <fstream>
#include <iostream>
#include <string>
#include <curl/curl.h>
#include <rapidjson/document.h>
// std::set<std::string> listPinnedCIDs() {
//     std::set<std::string> result;
//     std::string cmd = "ipfs pin ls --type=recursive -q";
//     std::string output = runShellCommand(cmd);
//     std::istringstream ss(output);
//     std::string line;
//     while (std::getline(ss, line)) {
//         if (!line.empty()) result.insert(line);
//     }
//     return result;
// }

// 回调函数：把响应写入 string
size_t WriteToFile(void* ptr, size_t size, size_t nmemb, void* userdata) {
    std::ofstream* ofs = static_cast<std::ofstream*>(userdata);
    ofs->write(static_cast<char*>(ptr), size * nmemb);
    return size * nmemb;
}

bool downloadFromIPFS(const std::string& endpoint, const std::string& cid, const std::string& outPath) {
    std::string url = "http://" + endpoint + "/api/v0/cat?arg=" + cid;

    std::cout << "url : " << url << std::endl;

    CURL* curl = curl_easy_init();
    if (!curl) return false;

    std::ofstream ofs(outPath, std::ios::binary);
    if (!ofs.is_open()) {
        std::cerr << "Failed to open output file: " << outPath << std::endl;
        curl_easy_cleanup(curl);
        return false;
    }

    curl_easy_setopt(curl, CURLOPT_POST, 1L);

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteToFile);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &ofs);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);  // 允许重定向

    CURLcode res = curl_easy_perform(curl);
    curl_easy_cleanup(curl);
    ofs.close();

    return res == CURLE_OK;
}

// 回调函数：把返回的数据写入 string（不要在这里解析 JSON！）
size_t WriteCallback(void* contents, size_t size, size_t nmemb, void* userp) {
    ((std::string*)userp)->append((char*)contents, size * nmemb);
    return size * nmemb;
}

std::string uploadFileToIpfs(const std::string& endpoint, const std::string& file_path) {
    std::string cid_res = "";
    CURL* curl;
    CURLcode res;
    std::string response;

    curl_global_init(CURL_GLOBAL_DEFAULT);
    curl = curl_easy_init();
    if (curl) {
        struct curl_httppost* formpost = nullptr;
        struct curl_httppost* lastptr = nullptr;

        // 替换成你要上传的文件绝对路径
        // const char* file_path = "/home/drw/storage/cpp-sdk/build/bin/config.yaml";
        // std::cout << "file_path: " << file_path << std::endl;

        // std::string full_file_path = "/home/drw/storage/dstorage/build/bin/" + file_path;

        curl_formadd(&formpost, &lastptr, CURLFORM_COPYNAME, "file", CURLFORM_FILE, file_path.c_str(),
            CURLFORM_CONTENTTYPE, "application/octet-stream", CURLFORM_END);

        std::string real_url = endpoint + "/api/v0/add";
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
            // std::cout << "Response: " << response << std::endl;
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