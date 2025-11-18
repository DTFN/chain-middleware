#include "isal.h"
#include <cassert>
#include <cmath>
#include <cstring>
#include <set>

// 读取文件内容到内存
bool read_file(const std::string& file_name, std::vector<uint8_t>& buffer) {
    std::ifstream file(file_name, std::ios::binary | std::ios::ate);
    if (!file.is_open()) {
        std::cerr << "无法打开文件: " << file_name << std::endl;
        return false;
    }
    size_t file_size = file.tellg();
    file.seekg(0, std::ios::beg);
    buffer.resize(file_size);
    file.read(reinterpret_cast<char*>(buffer.data()), file_size);
    return true;
}

// 将数据写入文件
bool write_file(const std::string& file_name, const std::vector<uint8_t>& buffer) {
    std::ofstream file(file_name, std::ios::binary);
    if (!file.is_open()) {
        std::cerr << "无法写入文件: " << file_name << std::endl;
        return false;
    }
    file.write(reinterpret_cast<const char*>(buffer.data()), buffer.size());
    return true;
}

void recover_data(int k, int m, int block_size,
    std::vector<uint8_t*>& shards,     // 所有分片（data + parity）
    std::vector<int>& shard_present,   // 标记每个分片是否存在
    std::vector<uint8_t*>& recovered)  // 输出恢复后的数据（指针数组）
{
    // Step 1: 生成编码矩阵
    std::vector<uint8_t> encode_matrix((k + m) * k);
    gf_gen_rs_matrix(encode_matrix.data(), k + m, k);

    // Step 2: 选取有效数据块，需 k 个
    std::vector<int> valid_indices;
    std::vector<uint8_t*> available_ptrs;

    for (int i = 0; i < k + m && valid_indices.size() < static_cast<size_t>(k); ++i) {
        if (shard_present[i]) {
            valid_indices.push_back(i);
            available_ptrs.push_back(shards[i]);
        }
    }

    assert(valid_indices.size() == static_cast<size_t>(k));  // 至少 k 个可用分片

    // Step 3: 构造解码矩阵
    std::vector<uint8_t> decode_matrix(k * k);
    for (int i = 0; i < k; ++i) {
        memcpy(&decode_matrix[i * k], &encode_matrix[valid_indices[i] * k], k);
    }

    // Step 4: 解码矩阵求逆
    std::vector<uint8_t> invert_matrix(k * k);
    std::vector<int> tmp(k);
    gf_invert_matrix(decode_matrix.data(), invert_matrix.data(), k);

    // Step 5: 构造恢复矩阵
    std::vector<int> recovery_indices;
    for (int i = 0; i < k; ++i) {
        if (!shard_present[i]) {
            recovery_indices.push_back(i);
        }
    }

    std::vector<uint8_t> recovery_matrix(recovery_indices.size() * k);
    for (size_t r = 0; r < recovery_indices.size(); ++r) {
        uint8_t* row = &encode_matrix[recovery_indices[r] * k];
        for (int j = 0; j < k; ++j) {
            uint8_t val = 0;
            for (int l = 0; l < k; ++l) {
                val ^= gf_mul(row[l], invert_matrix[l * k + j]);
            }
            recovery_matrix[r * k + j] = val;
        }
    }

    // Step 6: 构造 Galois Field 编码表
    std::vector<uint8_t> gftbls(32 * recovery_indices.size() * k);
    ec_init_tables(k, recovery_indices.size(), recovery_matrix.data(), gftbls.data());

    // Step 7: 构造输出指针
    std::vector<uint8_t*> recovery_ptrs;
    for (int idx : recovery_indices) {
        recovery_ptrs.push_back(recovered[idx]);  // recovered[i] 预分配好 block_size 大小
    }

    // Step 8: 执行恢复
    ec_encode_data(block_size, k, recovery_indices.size(), gftbls.data(), available_ptrs.data(), recovery_ptrs.data());

    // Step 9: 将恢复结果写回 shards[]
    for (size_t i = 0; i < recovery_indices.size(); ++i) {
        shards[recovery_indices[i]] = recovered[recovery_indices[i]];
        shard_present[recovery_indices[i]] = 1;  // 标记已恢复
    }
}

// 生成文件分片的函数
std::vector<std::string> IsalManager::GenerateErasureCodeShards(
    const std::string& input_file, int k, int m, uint64_t& file_size, bool write_to_disk) {
    // 文件读取
    std::cout << "isal enter erasure code shards" << std::endl;
    std::vector<uint8_t> input_data;
    if (!read_file(input_file, input_data)) {
        throw std::runtime_error("文件读取失败");
    }
    file_size = input_data.size();

    std::cout << "read file completed, file_size : " << file_size << std::endl;

    const int DATA_SIZE = ceil(file_size / k);  // 每个分片的大小
    const int NUM_DATA_BLOCKS = k;              // 数据块数量
    const int NUM_PARITY_BLOCKS = m;            // 校验块数量
    input_data.resize(DATA_SIZE * NUM_DATA_BLOCKS, 0);

    const int total_size = static_cast<int>(input_data.size());

    // 准备数据分片空间，每个分片大小 DATA_SIZE
    std::vector<std::vector<unsigned char>> data_blocks(NUM_DATA_BLOCKS, std::vector<unsigned char>(DATA_SIZE, 0));
    std::vector<std::vector<unsigned char>> parity_blocks(NUM_PARITY_BLOCKS, std::vector<unsigned char>(DATA_SIZE, 0));

    // 将输入数据分割填充到 data_blocks
    for (int i = 0; i < NUM_DATA_BLOCKS; ++i) {
        int offset = i * DATA_SIZE;
        int copy_len = std::min(DATA_SIZE, total_size - offset);
        if (copy_len > 0) {
            memcpy(data_blocks[i].data(), input_data.data() + offset, copy_len);
        }
    }

    // 构造指针数组，传入 ec_encode_data
    std::vector<unsigned char*> data_ptrs(NUM_DATA_BLOCKS);
    std::vector<unsigned char*> parity_ptrs(NUM_PARITY_BLOCKS);

    for (int i = 0; i < NUM_DATA_BLOCKS; ++i) {
        data_ptrs[i] = data_blocks[i].data();
    }
    for (int i = 0; i < NUM_PARITY_BLOCKS; ++i) {
        parity_ptrs[i] = parity_blocks[i].data();
    }

    // 分配编码表内存（每个编码器需要 32 bytes × k）
    std::vector<uint8_t> gftbls(32 * k * m);
    // 创建编码矩阵
    std::vector<uint8_t> encode_matrix(k * (k + m));  // 全部矩阵
    gf_gen_rs_matrix(encode_matrix.data(), k + m, k);
    // 取出校验块对应的矩阵部分
    std::vector<uint8_t> encode_rows(m * k);
    memcpy(encode_rows.data(), encode_matrix.data() + k * k, m * k);  // 取冗余部分
    // 初始化编码表
    ec_init_tables(k, m, encode_rows.data(), gftbls.data());

    // 调用纠删码库函数
    ec_encode_data(DATA_SIZE, NUM_DATA_BLOCKS, NUM_PARITY_BLOCKS, gftbls.data(), data_ptrs.data(), parity_ptrs.data());

    // 保存分片的文件名集合
    std::vector<std::string> file_names;

    // 保存数据分片
    for (int i = 0; i < NUM_DATA_BLOCKS; ++i) {
        std::vector<uint8_t> raw_block(data_blocks[i].begin(), data_blocks[i].end());
        std::string raw_filename = std::string("./isal/") + "raw-" + std::to_string(i);
        if (write_to_disk) {
            if (!write_file(raw_filename, raw_block)) {
                throw std::runtime_error("保存文件 " + raw_filename + " 失败!");
            }
        }
        file_names.push_back(raw_filename);
    }

    // 保存校验分片
    for (int i = 0; i < NUM_PARITY_BLOCKS; ++i) {
        std::vector<uint8_t> check_block(parity_blocks[i].begin(), parity_blocks[i].end());
        std::string check_filename = std::string("./isal/") + "check-" + std::to_string(i);
        if (write_to_disk) {
            if (!write_file(check_filename, check_block)) {
                throw std::runtime_error("保存文件 " + check_filename + " 失败!");
            }
        }
        file_names.push_back(check_filename);
    }

    std::cout << "isal deal completed" << std::endl;

    return file_names;
}

// 恢复函数：从数据分片和校验分片中恢复原始文件
std::vector<uint8_t> IsalManager::RecoverFile(const std::vector<std::string>& data_shards_file,
    const std::vector<std::string>& parity_shards_file, int k, int m, int data_size, const std::string& out_put_file) {
    int block_size = ceil(data_size / k);  // 每个块的大小（字节）

    // 从文件读取的分片数据
    std::vector<std::vector<uint8_t>> raw_blocks(k, std::vector<uint8_t>(block_size));
    std::vector<std::vector<uint8_t>> check_blocks(m, std::vector<uint8_t>(block_size));

    std::set<int> raw_lost_info;
    std::set<int> parity_lost_info;
    // 填充 raw_blocks[i] 和 check_blocks[i]，实际中是读文件
    for (int i = 0; i < k; i++) {
        if (!read_file(data_shards_file[i], raw_blocks[i])) {
            raw_lost_info.insert(i);
        }
    }

    for (int i = 0; i < m; i++) {
        if (!read_file(parity_shards_file[i], check_blocks[i])) {
            parity_lost_info.insert(i);
        }
    }

    // 构造 shards：data + parity 指针拼接
    std::vector<uint8_t*> shards(k + m);
    std::vector<int> shard_present(k + m);  // 标记每个分片是否存在
    for (int i = 0; i < k; ++i) {
        if (raw_lost_info.find(i) != raw_lost_info.end()) {
            shards[i] = nullptr;
            shard_present[i] = 0;
            std::cout << "raw file " << i << " is lost" << std::endl;
            continue;
        }
        shards[i] = raw_blocks[i].data();  // 如果某块缺失，也可以设为 nullptr
        shard_present[i] = 1;
    }
    for (int i = 0; i < m; ++i) {
        if (parity_lost_info.find(i) != parity_lost_info.end()) {
            shards[k + i] = nullptr;
            shard_present[k + i] = 0;
            std::cout << "check file " << i << " is lost" << std::endl;
            continue;
        }
        shards[k + i] = check_blocks[i].data();  // 同上
        shard_present[k + i] = 1;
    }

    // recovered 用来接收恢复的块，注意大小必须和 block_size 相同
    std::vector<std::vector<uint8_t>> recovered_blocks(k, std::vector<uint8_t>(block_size));
    std::vector<uint8_t*> recovered(k);  // 输出恢复后的数据（指针数组）
    for (int i = 0; i < k; ++i) {
        recovered[i] = recovered_blocks[i].data();
        memset(recovered[i], 0, block_size);
    }
    recover_data(k, m, block_size, shards, shard_present, recovered);

    std::vector<uint8_t> full_data(k * block_size);  // 用于拼接结果

    for (int i = 0; i < k; ++i) {
        std::memcpy(&full_data[i * block_size], shards[i], block_size);
    }
    // 截断为原始文件长度
    std::cout << "data_size : " << data_size << std::endl;
    full_data.resize(data_size);
    write_file(out_put_file, full_data);
    std::cout << "recover file successful" << " file_path : " << out_put_file << std::endl;

    return full_data;
}
