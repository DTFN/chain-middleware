/**
 * @copyright (C) 2023 0shu.
 *
 * @file: tools.h
 * @brief: 工具类
 */

#pragma once

#include <array>
#include <fstream>
#include <iomanip>
#include <iostream>
#include <memory>
#include <sstream>
#include <string>
#include <vector>
#include <sodium.h>
#include <boost/uuid/uuid.hpp>
#include <boost/uuid/uuid_generators.hpp>
#include <boost/uuid/uuid_io.hpp>
#include <openssl/evp.h>
#include <openssl/md5.h>
#include <openssl/sha.h>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../ipfs/ipfs_client.h"
#include "rapidjson/prettywriter.h"

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
    std::string cmd = "node " + js_file + " '" + input_file + "'";
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

// 使用模板 + 折叠表达式
template <typename... Args>
inline std::string ExecNode(const std::string& js_file, Args&&... args) {
    std::ostringstream cmd;
    cmd << "node " << js_file;

    // 逐个追加参数并加上引号
    ((cmd << " '" << args << "'"), ...);

    std::cout << "cmd : " << cmd.str() << std::endl;

    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.str().c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }

    return result;
}

inline std::string GenerateJson(const std::string& user, const std::string& raw_file_id, const std::string& pk,
    const std::string& field_id, const std::string& file_name, uint32_t file_size, uint32_t total_shards,
    uint32_t data_shards, const std::vector<std::string>& ipfs_urls, const std::vector<std::string>& shard_hashes,
    const std::vector<std::string>& poseidon_hashes) {
    using namespace rapidjson;

    // 初始化 JSON Document
    Document d;
    d.SetObject();
    Document::AllocatorType& allocator = d.GetAllocator();

    Value input(kObjectType);
    // 添加字符串字段
    input.AddMember("rawFileId", Value(raw_file_id.c_str(), allocator), allocator);
    input.AddMember("pk", Value(pk.c_str(), allocator), allocator);
    input.AddMember("author", Value(user.c_str(), allocator), allocator);
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

inline bool verifyProof(
    const std::string& verification_key_file, const std::string& public_json_file, const std::string& proof_json_file) {
    // 验证证明 输出[INFO]  snarkJS: OK!
    // snarkjs groth16 verify verification_key.json public.json proof.json
    std::cout << "verify proof..." << std::endl;
    std::string cmd =
        "snarkjs groth16 verify " + verification_key_file + " " + proof_json_file + " " + public_json_file;
    std::cout << "cmd: " << cmd << std::endl;
    std::array<char, 128> buffer;
    std::string result;

    std::unique_ptr<FILE, decltype(&pclose)> pipe(popen(cmd.c_str(), "r"), pclose);
    if (!pipe) {
        throw std::runtime_error("popen() failed!");
    }

    while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
        result += buffer.data();
    }

    std::cout << "result: " << result << std::endl;
    return result.find("OK") != std::string::npos;
}
inline bool ParseChallengeJsonAndVerify(std::string json) {
    std::vector<std::string> result;
    rapidjson::Document d;
    rapidjson::ParseResult ok = d.Parse(json.c_str());
    if (!ok) {
        std::string err = std::string("JSON parse error");
        std::cout << err << std::endl;
        return false;
    }
    if (!d.IsObject()) {
        std::string err = std::string("Root must be an object");
        std::cout << err << std::endl;
        return false;
    }

    // 必填：challenge_id
    // if (!d.HasMember("challenge_id") || !d["challenge_id"].IsUint64()) {
    //     std::string err = std::string("Field 'challenge_id' missing or not uint64");
    //     std::cout << err << std::endl;
    //     return result;
    // }
    auto challenge_id = d["challenge_id"].GetUint64();
    auto challenge_id_str = std::to_string(challenge_id);
    auto proof_cid = d["proof_cid"].GetString();
    auto proof_url = d["proof_url"].GetString();
    auto proof_hash = d["proof_hash"].GetString();
    auto public_cid = d["public_cid"].GetString();
    auto public_url = d["public_url"].GetString();
    auto public_hash = d["public_hash"].GetString();
    result.push_back(challenge_id_str);
    result.push_back(proof_cid);
    result.push_back(proof_url);
    result.push_back(proof_hash);
    result.push_back(public_cid);
    result.push_back(public_url);
    result.push_back(public_hash);

    std::string proof_file_name = "./tmp/proof.json";
    IpfsClient proof_ipfs_client(proof_url);
    proof_ipfs_client.DownloadFile(proof_cid, proof_file_name);
    std::cout << "download proof.json" << std::endl;

    std::string public_file_name = "./tmp/public.json";
    IpfsClient public_pfs_client(public_url);
    public_pfs_client.DownloadFile(public_cid, public_file_name);
    std::cout << "download public.json" << std::endl;

    return verifyProof("script/verification_key.json", proof_file_name, public_file_name);
}

// 从文件读取整个内容
inline std::string ReadFile(const std::string& path) {
    std::ifstream in(path);
    if (!in) throw std::runtime_error("Cannot open file: " + path);
    return std::string((std::istreambuf_iterator<char>(in)), std::istreambuf_iterator<char>());
}

inline std::string CharToHex(const unsigned char* data, std::size_t len) {
    std::ostringstream oss;
    oss << std::hex << std::setfill('0');
    for (std::size_t i = 0; i < len; ++i) {
        oss << std::setw(2) << static_cast<int>(data[i]);
    }
    return oss.str();
}

inline bool HashFile(const std::string& path, std::string& out_hex) {
    std::cout << "cal hash path : " << path << std::endl;
    std::ifstream ifs(path, std::ios::binary);
    if (!ifs) return false;

    EVP_MD_CTX* ctx = EVP_MD_CTX_new();
    if (!ctx) return false;

    if (EVP_DigestInit_ex(ctx, EVP_sha256(), nullptr) != 1) {
        EVP_MD_CTX_free(ctx);
        return false;
    }

    const std::size_t BUF_SIZE = 8192;
    std::vector<char> buffer(BUF_SIZE);
    while (ifs.good()) {
        ifs.read(buffer.data(), static_cast<std::streamsize>(BUF_SIZE));
        std::streamsize readBytes = ifs.gcount();
        if (readBytes > 0) {
            if (EVP_DigestUpdate(ctx, buffer.data(), static_cast<size_t>(readBytes)) != 1) {
                EVP_MD_CTX_free(ctx);
                return false;
            }
        }
    }

    unsigned int out_len = EVP_MD_size(EVP_sha256());
    std::vector<unsigned char> digest(out_len);
    if (EVP_DigestFinal_ex(ctx, digest.data(), &out_len) != 1) {
        EVP_MD_CTX_free(ctx);
        return false;
    }

    out_hex = CharToHex(digest.data(), out_len);
    // EVP_MD_CTX_free(ctx);
    return true;
}

inline bool HashData(const std::vector<unsigned char>& data, std::string& out_hex) {
    if (data.empty()) return false;

    EVP_MD_CTX* ctx = EVP_MD_CTX_new();
    if (!ctx) return false;

    if (EVP_DigestInit_ex(ctx, EVP_sha256(), nullptr) != 1) {
        EVP_MD_CTX_free(ctx);
        return false;
    }

    if (EVP_DigestUpdate(ctx, data.data(), data.size()) != 1) {
        EVP_MD_CTX_free(ctx);
        return false;
    }

    unsigned int out_len = EVP_MD_size(EVP_sha256());
    std::vector<unsigned char> digest(out_len);

    if (EVP_DigestFinal_ex(ctx, digest.data(), &out_len) != 1) {
        EVP_MD_CTX_free(ctx);
        return false;
    }

    // EVP_MD_CTX_free(ctx);

    out_hex = CharToHex(digest.data(), out_len);
    return true;
}

inline bool DeriveSeed(const std::string& username, const std::string& password, const std::string& salt_str,
    unsigned char seed[crypto_sign_SEEDBYTES]) {
    // Combine username and password as input to KDF (you can adjust)
    std::string input = username + ":" + password;

    // Prepare salt: must be crypto_pwhash_SALTBYTES = 16
    unsigned char salt[crypto_pwhash_SALTBYTES];
    // If provided salt_str is shorter, pad/derive deterministically; here we
    // simple fill.
    memset(salt, 0, sizeof(salt));
    if (!salt_str.empty()) {
        size_t copy = std::min(sizeof(salt), salt_str.size());
        memcpy(salt, salt_str.data(), copy);
    } else {
        // Use fixed system-wide salt is acceptable but not ideal; recommended to
        // use per-user random salt stored off-chain
        const char* default_salt = "default_system_salt";
        size_t copy = std::min(sizeof(salt), strlen(default_salt));
        memcpy(salt, default_salt, copy);
    }

    // Use crypto_pwhash to generate 32-byte seed (seed size for Ed25519)
    // params: opslimit & memlimit choose moderate values. For production tune
    // these.
    if (crypto_pwhash(seed, crypto_sign_SEEDBYTES, input.c_str(), input.size(), salt, crypto_pwhash_OPSLIMIT_MODERATE,
            crypto_pwhash_MEMLIMIT_MODERATE, crypto_pwhash_ALG_ARGON2ID13) != 0) {
        // out of memory
        return false;
    }
    return true;
}

inline std::vector<unsigned char> HexToBytes(const std::string& hex) {
    std::vector<unsigned char> bytes;
    bytes.reserve(hex.size() / 2);
    for (size_t i = 0; i < hex.size(); i += 2) {
        std::string byteString = hex.substr(i, 2);
        unsigned char byte = static_cast<unsigned char>(strtol(byteString.c_str(), nullptr, 16));
        bytes.push_back(byte);
    }
    return bytes;
}

inline bool Encrypt(std::vector<unsigned char>& ciphertext, std::vector<unsigned char>& source, const std::string& pk) {
    std::cout << "encrypt file"
              << " pk: " << pk << std::endl;
    auto pk_bytes = HexToBytes(pk);
    if (crypto_box_seal(ciphertext.data(), source.data(), source.size(), pk_bytes.data()) != 0) {
        std::cerr << "crypto_box_seal failed\n";
        return false;
    }
    std::cout << "encrypt file end" << std::endl;

    return true;
}

inline bool Decrypt(std::vector<unsigned char>& target, std::vector<unsigned char>& ciphertext, const std::string& pk,
    const std::string& sk) {
    std::cout << "decrypt file"
              << " pk: " << pk << std::endl;
    std::cout << "decrypt file"
              << " sk: " << sk << std::endl;
    auto pk_bytes = HexToBytes(pk);
    auto sk_bytes = HexToBytes(sk);

    if (crypto_box_seal_open(target.data(), ciphertext.data(), ciphertext.size(), pk_bytes.data(), sk_bytes.data()) !=
        0) {
        std::cerr << "crypto_box_seal_open failed\n";
        return false;
    }

    return true;
}