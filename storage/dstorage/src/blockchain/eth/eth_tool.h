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

inline uint64_t getUint256AsUint64(const std::vector<uint8_t>& data, size_t offset) {
    uint64_t value = 0;
    for (int i = 24; i < 32; ++i) {  // 取最后8字节作为uint64
        value = (value << 8) | data[offset + i];
    }
    return value;
}

// inline size_t getUint256(const std::vector<uint8_t>& data, size_t offset) {
//     if (offset + 32 > data.size()) throw std::out_of_range("getUint256: out of range");

//     size_t value = 0;
//     for (size_t i = 24; i < 32; i++) {  // 只取低 8 字节
//         value = (value << 8) | data[offset + i];
//     }
//     return value;
// }

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
inline std::string getString(
    const std::vector<uint8_t>& data, size_t base, size_t offset, bool need_cal_offset = true) {
    // std::cout << "offset" << offset << std::endl;
    // 偏移位置
    size_t length = 0;
    if (need_cal_offset) {
        size_t dynamicOffset = getUint64(data, offset);
        // std::cout << "dynamicOffset" << dynamicOffset << std::endl;
        size_t realOffset = base + dynamicOffset;
        offset = realOffset;

        length = getUint64(data, realOffset);
    } else {
        length = getUint64(data, offset);
        std::cout << std::endl;
        for (int i = 0; i < 32; i++) {
            std::cout << data[offset + i];
        }
        std::cout << std::endl;
    }
    // std::cout << "length: " << length << std::endl;
    std::string str;
    for (size_t i = 0; i < length; i++) {
        std::string s(data.begin() + offset + 32, data.begin() + offset + 32 + 96);
        str = s;
    }
    return str;
}

inline std::vector<std::string> getStringArray(const std::vector<uint8_t>& data,
    size_t base,   // 一般传 0
    size_t offset  // 数组的偏移量（相对于 base）
) {
    std::vector<std::string> result;

    size_t arrayStart = base + offset;                 // 数组实际开始位置
    size_t arrayLength = getUint64(data, arrayStart);  // 读取长度 N
    std::cout << "arrayStart: " << arrayStart << " arrayLength : " << arrayLength << std::endl;

    for (size_t i = 0; i < arrayLength; i++) {
        size_t elemOffsetRel = getUint64(data, arrayStart + 32 + i * 32);  // 元素偏移（相对于数组起点）
        size_t elemOffset = arrayStart + elemOffsetRel;                    // 元素实际偏移
        uint64_t strLen = getUint64(data, elemOffset);                     // 读取真正的字符串长度
        std::string s(data.begin() + elemOffset + 32, data.begin() + elemOffset + 32 + strLen);

        std::cout << "i: " << i << " elemOffsetRel " << elemOffsetRel << " eleOffset " << elemOffset
                  << " length: " << strLen << std::endl;
        // std::string s = getString(data, 0x00, elemOffset, false);          // 读取字符串
        std::cout << "i: " << i << " s: " << s << std::endl;
        result.push_back(s);
    }

    return result;
}

// inline std::string getString1(const std::vector<uint8_t>& data, size_t offset) {
//     if (offset + 32 > data.size()) throw std::out_of_range("getString: out of range");

//     size_t length = getUint64(data, offset);
//     std::cout << "length: " << length << std::endl;
//     if (offset + 32 + length > data.size()) throw std::out_of_range("getString: content out of range");

//     return std::string(data.begin() + offset + 32, data.begin() + offset + 32 + length);
// }

inline std::string getChallengeInfo(const std::string& result) {
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
    // std::cout << "challenger: " << challenger << "\n";
    // std::cout << "fileHash:   " << fileHash << "\n";
    // std::cout << "nonce:      " << nonce << "\n";
    // std::cout << "cid:        " << cid << "\n";
    // std::cout << "timestamp:  " << timestamp << "\n";
    // std::cout << "proofSubmitted: " << proofSubmitted << "\n";
    // std::cout << "proofValid:     " << proofValid << "\n";
    // std::cout << "prover:     " << prover << "\n";

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

inline std::string parseBoolResult(const std::string& hexResult) {
    // 安全检查
    if (hexResult.size() < 66 || hexResult.rfind("0x", 0) != 0) {
        std::cerr << "Invalid hex result: " << hexResult << std::endl;
        return "";
    }

    // 取最后两位（对应最后一个字节）
    std::string lastByte = hexResult.substr(hexResult.size() - 2);

    // 判断最后一个字节是否为 01
    bool ret = (lastByte == "01");

    // RapidJSON 构建 JSON
    rapidjson::Document doc;
    doc.SetArray();
    rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();

    doc.PushBack(ret, allocator);
    // 输出 JSON 字符串
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    doc.Accept(writer);

    std::string json_str = buffer.GetString();
    std::cout << json_str << std::endl;

    return json_str;
}

inline std::string getChallengeByUser(const std::string& result) {
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

inline std::string getFileMetaResult(const std::string& result) {
    // 1. 把 hex 转 bytes
    auto data = hexToBytes(result);

    // === 固定部分 (前 32 bytes 对应动态偏移) ===
    // uploader: address
    std::string uploader = getAddress(data, 0x00);

    // fileName: string (动态类型，slot=0x20 存 offset)
    size_t fileNameOffset = 0x20;

    // fileSize: uint256 (直接取)
    uint64_t fileSize = getUint64(data, 0x40);

    // totalShards
    uint64_t totalShards = getUint64(data, 0x60);

    // dataShards
    uint64_t dataShards = getUint64(data, 0x80);

    // ipfsUrls[] (动态数组，slot=0xa0)
    size_t ipfsUrlsOffset = getUint64(data, 0xa0);

    // shardHashes[] (slot=0xc0)
    size_t shardHashesOffset = getUint64(data, 0xc0);

    // poseidonHashes[] (slot=0xe0)
    size_t poseidonHashesOffset = getUint64(data, 0xe0);

    // timestamp (最后一个 uint256，直接取)
    uint64_t timestamp = getUint64(data, 0x100);

    // === 解析动态字段 ===
    std::string fileName = getString(data, 0x00, fileNameOffset);

    std::vector<std::string> ipfsUrls = getStringArray(data, 0x00, ipfsUrlsOffset);
    // std::vector<std::string> shardHashes = getStringArray(data, 0x00, shardHashesOffset);
    // std::vector<std::string> poseidonHashes = getStringArray(data, 0x00, poseidonHashesOffset);

    // === 打印调试 ===
    std::cout << "uploader: " << uploader << "\n";
    std::cout << "fileName: " << fileName << "\n";
    std::cout << "fileSize: " << fileSize << "\n";
    std::cout << "totalShards: " << totalShards << "\n";
    std::cout << "dataShards: " << dataShards << "\n";
    std::cout << "timestamp: " << timestamp << "\n";

    std::cout << "ipfsUrls: ";
    for (auto& u : ipfsUrls) std::cout << u << " ";
    std::cout << "\n";

    // std::cout << "shardHashes: ";
    // for (auto& h : shardHashes) std::cout << h << " ";
    // std::cout << "\n";

    // std::cout << "poseidonHashes: ";
    // for (auto& h : poseidonHashes) std::cout << h << " ";
    // std::cout << "\n";

    // === 构建 JSON ===
    rapidjson::Document doc;
    doc.SetObject();
    auto& allocator = doc.GetAllocator();

    doc.AddMember("uploader", rapidjson::Value(uploader.c_str(), allocator), allocator);
    doc.AddMember("fileName", rapidjson::Value(fileName.c_str(), allocator), allocator);
    doc.AddMember("fileSize", fileSize, allocator);
    doc.AddMember("totalShards", totalShards, allocator);
    doc.AddMember("dataShards", dataShards, allocator);

    // // ipfsUrls[]
    // {
    //     rapidjson::Value arr(rapidjson::kArrayType);
    //     for (auto& u : ipfsUrls) {
    //         arr.PushBack(rapidjson::Value(u.c_str(), allocator), allocator);
    //     }
    //     doc.AddMember("ipfsUrls", arr, allocator);
    // }

    // // shardHashes[]
    // {
    //     rapidjson::Value arr(rapidjson::kArrayType);
    //     for (auto& h : shardHashes) {
    //         arr.PushBack(rapidjson::Value(h.c_str(), allocator), allocator);
    //     }
    //     doc.AddMember("shardHashes", arr, allocator);
    // }

    // // poseidonHashes[]
    // {
    //     rapidjson::Value arr(rapidjson::kArrayType);
    //     for (auto& h : poseidonHashes) {
    //         arr.PushBack(rapidjson::Value(h.c_str(), allocator), allocator);
    //     }
    //     doc.AddMember("poseidonHashes", arr, allocator);
    // }

    doc.AddMember("timestamp", timestamp, allocator);

    // === 转 JSON 字符串 ===
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    doc.Accept(writer);

    std::string json_str = buffer.GetString();
    std::cout << json_str << std::endl;

    return json_str;
}
