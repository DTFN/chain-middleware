#include "poseidon.h"

std::vector<uint8_t> readFile(const std::string& filePath) {
    std::ifstream file(filePath, std::ios::binary);  // 打开文件（二进制模式）
    if (!file) {
        throw std::runtime_error("Failed to open file: " + filePath);
    }

    // 获取文件大小
    file.seekg(0, std::ios::end);
    std::streamsize size = file.tellg();
    file.seekg(0, std::ios::beg);

    // 读取文件内容
    std::vector<uint8_t> buffer(size);
    if (!file.read(reinterpret_cast<char*>(buffer.data()), size)) {
        throw std::runtime_error("Failed to read file: " + filePath);
    }

    return buffer;
}

std::string computePoseidonHash(const std::vector<uint8_t>& data) {
    // 可替换为真正的 Poseidon
    return ""; 
}
