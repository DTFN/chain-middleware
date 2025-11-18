#pragma once

#include <string>
#include <vector>
#include <mutex>

class IpfsClient {
public:
    // 构造函数，指定IPFS的API URL
    IpfsClient(const std::string& api_url);

    // 上传文件到IPFS，返回文件的CID
    std::string UploadFile(const std::string& file_path);

    // 下载文件，给定CID，返回文件内容
    void DownloadFile(const std::string& cid, const std::string& file_name);

    bool IsExist(const std::string& cid);

    // 获取文件的CID，给定文件路径
    std::string GetFileCID(const std::string& file_path);

private:
    std::string api_url_;  // IPFS API URL
    // 内部实现，上传文件至IPFS
    std::string sendFileToIpfs(const std::string& file_path);

    // 内部实现，下载文件
    void fetchFileFromIpfs(const std::string& cid, const std::string& outPath);

private:
    static bool is_curl_init_;
    static std::mutex mutex_;
};
