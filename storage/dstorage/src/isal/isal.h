#pragma once

#include <fstream>
#include <iostream>
#include <string>
#include <vector>
#include "isa-l.h"  // 引入 ISA-L 库头文件
#include "boost/noncopyable.hpp"

class IsalManager : public boost::noncopyable {
public:
    static IsalManager& Instance() {
        static IsalManager isal_manager;
        return isal_manager;
    }
    ~IsalManager() {}

    std::vector<std::string> GenerateErasureCodeShards(
        const std::string& input_file, int k, int m, uint64_t& file_size, bool write_to_disk = false);

    std::vector<uint8_t> RecoverFile(const std::vector<std::string>& data_shards_file,
        const std::vector<std::string>& parity_shards_file, int k, int m, int data_size,
        const std::string& out_put_file);

private:
    IsalManager() {}
};