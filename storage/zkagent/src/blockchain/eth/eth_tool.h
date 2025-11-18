#include <iomanip>
#include <iostream>
#include <sstream>
#include <string>
#include <vector>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>

// 把 hex 转成 bytes
inline std::vector<uint8_t> hexToBytes(const std::string& hex) {
    std::vector<uint8_t> bytes;
    for (size_t i = 2; i < hex.length(); i += 2) {  // 跳过 "0x"
        std::string byteString = hex.substr(i, 2);
        uint8_t byte = (uint8_t)strtol(byteString.c_str(), nullptr, 16);
        bytes.push_back(byte);
    }
    return bytes;
}

// 从 bytes 取 32 字节 big-endian 转 uint64_t
inline uint64_t getUint64(const std::vector<uint8_t>& data, size_t offset) {
    uint64_t value = 0;
    for (int i = 24; i < 32; i++) {
        value = (value << 8) | data[offset + i];
    }
    return value;
}

// 取 32字节里的地址 (最后 20 字节)
inline std::string getAddress(const std::vector<uint8_t>& data, size_t offset) {
    std::ostringstream oss;
    oss << "0x";
    for (int i = 12; i < 32; i++) {
        oss << std::hex << std::setw(2) << std::setfill('0') << (int)data[offset + i];
    }
    return oss.str();
}

// 解析动态 string
inline std::string getString(const std::vector<uint8_t>& data, size_t base, size_t offset) {
    std::cout << "getString " << std::endl;
    // 偏移位置
    size_t dynamicOffset = getUint64(data, offset);
    size_t realOffset = base + dynamicOffset;

    std::cout << "getString realOffset" << realOffset << std::endl;
    size_t length = getUint64(data, realOffset);
    std::cout << "getString length" << length << std::endl;
    std::string str;
    for (size_t i = 0; i < length; i++) {
        str.push_back((char)data[realOffset + 32 + i]);
    }
    return str;
}

inline std::string getChallengeInfo(const std::string& result) {
    if (result == "0x") {
        std::cout << "the challenge is null" << std::endl;
        return "";
    }
    //   std::string result =
    //       "0x000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000"
    //       "000000000000000000000000000000000000000000000000000000010000000000000000"
    //       "000000000000000000000000000000000000000000000001600000000000000000000000"
    //       "0000000000000000000000000000000000000001c0000000000000000000000000000000"
    //       "0000000000000000000000000068ad937400000000000000000000000000000000000000"
    //       "000000000000000000000000000000000000000000000000000000000000000000000000"
    //       "000000000000000000000000000000000000000000165734f8847cff904c6dce84929233"
    //       "0090fecffb00000000000000000000000000000000000000000000000000000000000000"
    //       "283136353733346638383437636646393034633664634538343932393233333030393046"
    //       "654366666200000000000000000000000000000000000000000000000000000000000000"
    //       "000000000000000000000000000000000000000000000000283136353733346638383437"
    //       "636646393034633664634538343932393233333030393046654366666200000000000000"
    //       "000000000000000000000000000000000000000000000000000000000000000000000000"
    //       "000000000000000000000000283136353733346638383437636646393034633664634538"
    //       "343932393233333030393046654366666200000000000000000000000000000000000000"
    //       "0000000000";

    auto data = hexToBytes(result);

    // 固定部分解析
    std::string challenger = getAddress(data, 0x00);
    size_t fileHashOffset = 0x20;
    size_t nonceOffset = 0x40;
    size_t cidOffset = 0x60;
    uint64_t timestamp = getUint64(data, 0x80);
    bool proofSubmitted = getUint64(data, 0xa0);
    bool proofValid = getUint64(data, 0xc0);
    std::string prover = getAddress(data, 0xe0);

    // 动态字符串
    std::string fileHash = getString(data, 0x00, fileHashOffset);
    std::string nonce = getString(data, 0x00, nonceOffset);
    std::string cid = getString(data, 0x00, cidOffset);

    // 打印
    std::cout << "challenger: " << challenger << "\n";
    std::cout << "fileHash:   " << fileHash << "\n";
    std::cout << "nonce:      " << nonce << "\n";
    std::cout << "cid:        " << cid << "\n";
    std::cout << "timestamp:  " << timestamp << "\n";
    std::cout << "proofSubmitted: " << proofSubmitted << "\n";
    std::cout << "proofValid:     " << proofValid << "\n";
    std::cout << "prover:     " << prover << "\n";

    // RapidJSON 构建 JSON
    rapidjson::Document doc;
    doc.SetArray();
    rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();

    // 顺序添加元素
    doc.PushBack(rapidjson::Value(challenger.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(fileHash.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(nonce.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(cid.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(std::to_string(timestamp).c_str(), allocator), allocator);
    doc.PushBack(proofSubmitted, allocator);
    doc.PushBack(proofValid, allocator);
    doc.PushBack(rapidjson::Value(prover.c_str(), allocator), allocator);

    // 输出 JSON 字符串
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    doc.Accept(writer);

    std::string json_str = buffer.GetString();
    std::cout << json_str << std::endl;

    return json_str;
}

inline std::string getChallengeByUser(const std::string& result) {
    if (result == "0x") {
        std::cout << "the challenge is null" << std::endl;
        return "";
    }
    //   std::string result =
    //       "0x0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb922660000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001e00000000000000000000000000000000000000000000000000000000068aebbac00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002831363537333466383834376366463930346336646345383439323932333330303930466543666662000000000000000000000000000000000000000000000000";

    auto data = hexToBytes(result);

    // 固定部分
    uint64_t challengeId = getUint64(data, 0x00);
    std::string challenger = getAddress(data, 0x20);
    size_t fileHashOffset = 0x40;
    size_t nonceOffset = 0x60;
    size_t cidOffset = 0x80;

    uint64_t timestamp = getUint64(data, 0xa0);
    bool proofSubmitted = (getUint64(data, 0xc0) != 0);
    bool proofValid = (getUint64(data, 0xe0) != 0);
    std::string prover = getAddress(data, 0x100);

    // 动态字符串
    std::string fileHash = getString(data, 0x00, fileHashOffset);
    std::string nonce = getString(data, 0x00, nonceOffset);
    std::string cid = getString(data, 0x00, cidOffset);

    // 打印
    std::cout << "challengeId: " << challengeId << "\n";
    std::cout << "challenger: " << challenger << "\n";
    std::cout << "fileHash:   " << fileHash << "\n";
    std::cout << "nonce:      " << nonce << "\n";
    std::cout << "cid:        " << cid << "\n";
    std::cout << "timestamp:  " << timestamp << "\n";
    std::cout << "proofSubmitted: " << proofSubmitted << "\n";
    std::cout << "proofValid:     " << proofValid << "\n";
    std::cout << "prover:     " << prover << "\n";

    // RapidJSON 构建 JSON
    rapidjson::Document doc;
    doc.SetArray();
    rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();

    // 顺序添加元素
    doc.PushBack(rapidjson::Value(std::to_string(challengeId).c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(challenger.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(fileHash.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(nonce.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(cid.c_str(), allocator), allocator);
    doc.PushBack(rapidjson::Value(std::to_string(timestamp).c_str(), allocator), allocator);
    doc.PushBack(proofSubmitted, allocator);
    doc.PushBack(proofValid, allocator);
    doc.PushBack(rapidjson::Value(prover.c_str(), allocator), allocator);

    // 输出 JSON 字符串
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    doc.Accept(writer);

    std::string json_str = buffer.GetString();
    std::cout << json_str << std::endl;

    return json_str;
}
