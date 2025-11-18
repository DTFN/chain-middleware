/**
 * @copyright (C) 2023 0shu.
 *
 * @file: tools.h
 * @brief: 工具类
 */

#pragma once

#include <array>
#include <memory>
#include <string>
#include <iostream>
#include <boost/uuid/uuid.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "rapidjson/prettywriter.h"
#include <fstream>

inline std::string uuid() {
    // 创建 UUID 生成器
    boost::uuids::random_generator gen;
    // 生成 UUID
    boost::uuids::uuid new_uuid = gen();

    return boost::uuids::to_string(new_uuid);
}

inline std::string uuid32() {
    // 创建 UUID 生成器 
    boost::uuids::random_generator gen;
    // 生成 UUID（36 字符，包括4个'-'）
    boost::uuids::uuid new_uuid = gen();

    // 转为字符串，并移除中划线
    std::string uuid_str = boost::uuids::to_string(new_uuid);
    uuid_str.erase(std::remove(uuid_str.begin(), uuid_str.end(), '-'), uuid_str.end());

    return uuid_str + uuid_str;  // 返回32位十六进制字符串
}

inline std::string ExecNode(const std::string js_file, const std::string& input_file) {
    std::string cmd = "node " + js_file + " " + input_file;
    std::cout << "cmd : " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }

    return result;
}

inline std::string GenerateJson(const std::string& field_id, const std::string& file_name, uint32_t file_size,
    uint32_t total_shards, uint32_t data_shards, const std::vector<std::string>& ipfs_urls,
    const std::vector<std::string>& shard_hashes, const std::vector<std::string>& poseidon_hashes) {
    using namespace rapidjson;

    // 初始化 JSON Document
    Document d;
    d.SetObject();
    Document::AllocatorType& allocator = d.GetAllocator();

    Value input(kObjectType);
    // 添加字符串字段
    input.AddMember("fileId", Value(field_id.c_str(), allocator), allocator);
    input.AddMember("fileName", Value(file_name.c_str(), allocator), allocator);

    // 添加数值字段
    input.AddMember("fileSize", file_size, allocator);
    input.AddMember("totalShards", total_shards, allocator);
    input.AddMember("dataShards", data_shards, allocator);

    // 添加数组字段
    Value ipfsUrls(kArrayType);
    for (auto ipfs_url : ipfs_urls) {
        ipfsUrls.PushBack(Value(ipfs_url.c_str(), allocator).Move(), allocator);
    }
    input.AddMember("ipfsUrls", ipfsUrls, allocator);

    Value shardHashes(kArrayType);
    for (auto shard_hash : shard_hashes) {
        shardHashes.PushBack(Value(shard_hash.c_str(), allocator).Move(), allocator);
    }
    input.AddMember("shardHashes", shardHashes, allocator);

    Value poseidonHashes(kArrayType);
    for (auto poseidon_hash : poseidon_hashes) {
        poseidonHashes.PushBack(Value(poseidon_hash.c_str(), allocator).Move(), allocator);
    }
    input.AddMember("poseidonHashes", poseidonHashes, allocator);

    // 把 input 加入到主对象
    d.AddMember("input", input, allocator);
    // 输出 JSON 字符串
    StringBuffer buffer;
    Writer<StringBuffer> writer(buffer);
    d.Accept(writer);

    std::string result = buffer.GetString();

    std::cout << "input json : " << result << std::endl;

    return result;
}

// 从文件读取整个内容
inline std::string ReadSolFile(const std::string& path) {
    std::ifstream in(path);
    if (!in) throw std::runtime_error("Cannot open file: " + path);
    return std::string((std::istreambuf_iterator<char>(in)), std::istreambuf_iterator<char>());
}