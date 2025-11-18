#include <iostream>
#include <tuple>
#include <sodium.h>
#include <boost/filesystem.hpp>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "blockchain/blockchain.h"
#include "blockchain/blockchain_factory.h"
#include "config/config.h"
#include "event/task_queue.h"
#include "http/http_server.h"
#include "ipfs/ipfs_client.h"
#include "ipfs/ipfs_manager.h"
#include "isal/isal.h"
#include "tools/tools.h"
#include "user/user.h"

std::map<std::string, std::string> challenge_id_file_id_map;
std::map<std::string, std::string> cid_file_id_map;
// std::map<std::string, std::string> cid_file_id_map;

namespace fs = boost::filesystem;

std::mutex mutex;
std::map<std::string, time_t> cid_map;        // 存储了需要发起挑战的文件名和上次挑战的时间
std::map<std::string, std::string> hash_map;  // 存储了cid与poseidon的映射，需要放在挑战中

std::string GetFileName(const std::string file_path) {
    size_t pos = file_path.find_last_of("/\\");
    std::string file_name = (pos == std::string::npos) ? file_path : file_path.substr(pos + 1);
    return file_name;
}

void CreateDirIfNotExist(const std::string path) {
    fs::path dirPath(path);

    try {
        if (!fs::exists(dirPath)) {
            if (fs::create_directories(dirPath)) {
                std::cout << "目录已创建: " << dirPath << '\n';
            } else {
                std::cout << "目录未能创建，可能已存在。\n";
            }
        } else {
            std::cout << "目录已存在。\n";
        }
    } catch (const fs::filesystem_error& ex) {
        std::cerr << "错误: " << ex.what() << '\n';
    }
}

std::shared_ptr<BlockChain> block_chain = nullptr;

void BackgroundHandler(const std::string& user, const std::string& raw_file_id, const std::string& pk,
    const std::string& file_id, const std::string& file_name, const std::string& file_full_name) {
    try {
        std::cout << "[*] Processing: " << file_full_name << "\n";

        // 处理逻辑
        // 是纠删码进行分片
        uint64_t file_size = 0;
        auto isal_slice = IsalManager::Instance().GenerateErasureCodeShards(
            file_full_name, Config::Instance().isal_k(), Config::Instance().isal_m(), file_size, true);

        std::vector<std::string> data_shards_file;
        std::vector<std::string> parity_shards_file;
        for (int i = 0; i < Config::Instance().isal_k(); i++) {
            data_shards_file.push_back(isal_slice[i]);
        }

        for (int i = 0; i < Config::Instance().isal_m(); i++) {
            parity_shards_file.push_back(isal_slice[Config::Instance().isal_k() + i]);
        }

        // 将CID上链
        std::vector<std::string> ipfses;
        std::vector<std::string> cids;
        std::vector<std::string> poseidons;
        std::vector<std::pair<std::string, std::string>> cid_vec;

        // // 将纠删码分片上传到ipfs并记录CID
        for (auto slice : isal_slice) {
            std::cout << "slice is " << slice << std::endl;
            auto endpoint = Config::Instance().GetRandomAddress();
            std::string url = "http://" + endpoint;
            IpfsClient ipfs_client(url);
            auto cid = ipfs_client.UploadFile(slice);
            std::string poseidon_hash = ExecNode("script/generate_poseidon_hash.js", slice);
            cid_vec.push_back(std::make_pair(cid, poseidon_hash));
            ipfses.push_back(url);
            cids.push_back(cid);
            poseidons.push_back(poseidon_hash);

            IpfsManager::Instance().Add(cid, endpoint);
        }

        std::string contract_address = block_chain->GetAddress("FileStore");
        std::vector<std::string> params;
        std::string json = GenerateJson(user, raw_file_id, pk, file_id, file_name, file_size,
            Config::Instance().isal_all(), Config::Instance().isal_k(), ipfses, cids, poseidons);
        params.push_back(json);

        std::string output;
        bool ret = block_chain->Call("solidity", "FileStore", contract_address, "storeFile", params, output);
        if (!ret) {
            std::cout << "storeFile failed" << std::endl;
        }

        std::cout << "storeFile success" << std::endl;
        {
            std::lock_guard<std::mutex> lock(mutex);
            for (auto cid : cid_vec) {
                cid_map[cid.first] = 0;
                hash_map[cid.first] = cid.second;
                cid_file_id_map[cid.first] = raw_file_id + "|" + pk + "|" + file_id;
            }
        }
    } catch (...) {
        std::cout << "file store failed" << std::endl;
    }
}

std::string GetChallengeID(const std::string& json_str) {
    std::cout << "json_str : " << json_str << std::endl;
    // 解析 JSON 字符串
    rapidjson::Document doc;
    if (doc.Parse(json_str.c_str()).HasParseError()) {
        std::cerr << "JSON parse error" << std::endl;
        return "";
    }

    // 确保是数组且非空
    if (doc.IsArray() && !doc.Empty()) {
        const auto& val = doc[0];
        if (val.IsString()) {
            std::string id = val.GetString();
            return id;
        } else {
            std::string id = std::to_string(val.GetInt());
            return id;
        }
    } else {
        std::cerr << "Not an array or array is empty" << std::endl;
    }
    return "";
}

std::tuple<std::string, std::string, std::string, std::string, uint64_t, bool, bool, std::string> GetCheckChallengeInfo(
    const std::string& json_str) {
    rapidjson::Document doc;
    doc.Parse(json_str.c_str());

    if (!doc.IsArray() || doc.Size() != 8) {
        std::cerr << "Invalid JSON format" << std::endl;
        return std::make_tuple("", "", "", "", 0, false, false, "");
    }

    std::string challenger = doc[0].GetString();
    std::string fileHashStr = doc[1].GetString();  // 如果你想转 uint256 可用 boost::multiprecision
    std::string nonce = doc[2].GetString();
    std::string cid = doc[3].GetString();
    bool proofSubmitted = false;
    if (doc[5].IsBool()) {
        proofSubmitted = doc[5].GetBool();
    } else {
        proofSubmitted = (std::string(doc[5].GetString()) == "true");
    }
    bool proofValid = false;
    if (doc[6].IsBool()) {
        proofValid = doc[6].GetBool();
    } else {
        proofValid = (std::string(doc[6].GetString()) == "true");
    }
    std::string prover = doc[7].GetString();

    // 输出确认
    std::cout << "challenger:     " << challenger << std::endl;
    std::cout << "fileHash:       " << fileHashStr << std::endl;
    std::cout << "nonce:          " << nonce << std::endl;
    std::cout << "cid:            " << cid << std::endl;
    std::cout << "proofSubmitted: " << std::boolalpha << proofSubmitted << std::endl;
    std::cout << "proofValid:     " << std::boolalpha << proofValid << std::endl;
    std::cout << "prover:         " << prover << std::endl;

    return std::make_tuple(challenger, fileHashStr, nonce, cid, 0, proofSubmitted, proofValid, prover);
}

void GenerateChallenge(const std::vector<std::string>& cids, const std::vector<std::string>& hashes,
    std::vector<std::string>& challenge_ids) {
    // 构造文件哈希
    if (cids.size() != hashes.size()) {
        // error
        return;
    }

    for (auto i = 0; i < cids.size(); i++) {
        std::cout << "create the " << i << " challenge" << std::endl;
        std::string nonce = uuid32();
        std::vector<std::string> params;
        std::string cid = cids[i];
        std::string file_id = cid_file_id_map[cid];
        params.push_back(hashes[i]);
        params.push_back(nonce);
        params.push_back(cid);

        auto endpoint = IpfsManager::Instance().Get(cids[i]);
        auto account_address = HttpServer::Instance().GetAccount(endpoint);
        if (account_address.empty()) {
            std::cerr << "account_address is empty" << std::endl;
            // account_address = "165734f8847cfF904c6dcE849292330090FeCffb";
            continue;
        }
        params.push_back(account_address);  // 挑战者地址
        std::string contract_address = block_chain->GetAddress("Challenge");
        std::string output;
        bool ret = block_chain->Call("solidity", "Challenge", contract_address, "createChallenge", params, output);
        if (ret) {
            // challenge id
            auto id = GetChallengeID(output);
            if (id == "") {
                continue;
            }
            challenge_ids.push_back(id);
            challenge_id_file_id_map[id] = file_id;
        }
    }
}

std::vector<std::string> split(const std::string& s, char delimiter) {
    std::vector<std::string> tokens;
    std::string token;
    std::istringstream tokenStream(s);
    while (std::getline(tokenStream, token, delimiter)) {
        tokens.push_back(token);
    }
    return tokens;
}

bool reBuildAndUpload(const std::string& raw_file_id, const std::string& pk);
void CheckChallenge(std::queue<std::pair<std::vector<std::string>, time_t>>& challenge_queue) {
    static int checkcount = 0;
    checkcount++;
    if (checkcount <= 1) {
        // std::cout << "do nothing" << std::endl;
        return;
    }
    std::queue<std::pair<std::vector<std::string>, time_t>> new_queue;
    std::cout << "challenge_queue.size(): " << challenge_queue.size() << std::endl;
    while (challenge_queue.size() > 0) {
        std::pair<std::vector<std::string>, time_t> challenge = challenge_queue.front();
        std::cout << "challenge size is " << challenge.first.size() << std::endl;
        if (time(0) - challenge.second >= 1) {
            challenge_queue.pop();
        } else {
            std::cout << "time is not arrive" << std::endl;
            return;
        }

        bool need_enqueue = false;

        // check challenge
        auto submit_count = 0;
        auto valid_count = 0;
        std::string file_id = "";
        for (auto it : challenge.first) {
            std::vector<std::string> params;
            params.push_back(it);
            if (file_id == "") {
                file_id = challenge_id_file_id_map[it];
            }
            // std::cout << "the challenge " << it << " is match file " << challenge_id_file_id_map[it] << std::endl;
            std::string contract_address = block_chain->GetAddress("Challenge");
            std::string output;
            bool ret = block_chain->Call("solidity", "Challenge", contract_address, "getChallengeInfo", params, output);
            if (ret) {
                std::cout << "out_put: " << output << std::endl;
                // address challenger, 1
                // string memory fileHash, 2
                // string memory nonce, 3
                // string memory cid, 4
                // uint timestamp, 5
                // bool proofSubmitted, 6
                // bool proofValid, 7
                // address prover 8
                // out_put:
                // ["0x5ed87ab8c43120dc8553f04a31cf4be0e37b578e","11365769016959303865992615604602367532759447157249400316645834812718765714672","6b0f11bdbee74ce79992f582475df290","QmajjidtL73D6CXRbBJgciKZsQjAAY4VtF26nSWpDs8r3t","1750840616155",false,false,"0x0000000000000000000000000000000000000000"]
                // TODO 统计通过校验的分片个数，大于阈值，则认为成功，小于 阈值，则需要重新组合，分片上传逻辑
                auto [challenger, fileHash, nonce, cid, timestamp, proofSubmitted, proofValid, prover] =
                    GetCheckChallengeInfo(output);
                if (proofSubmitted) {
                    submit_count++;
                }
                if (proofValid) {
                    valid_count++;
                }
            }
        }

        if (submit_count > 0 && valid_count > 0) {
            if (valid_count >= Config::Instance().isal_threshold()) {
                // 大于阈值
                std::cout << "this check successful, valid count is " << valid_count << " this threshold is "
                          << Config::Instance().isal_threshold() << std::endl;
            } else {
                if (valid_count < Config::Instance().isal_k()) {
                    std::cout << "[error] the unvalid count is bigger than " << Config::Instance().isal_m()
                              << std::endl;
                    std::vector<std::string> file_info = split(file_id, '|');
                    reBuildAndUpload(file_info[0], file_info[1]);
                    // exit(0);
                    // error
                } else {
                    std::cout << "some files are not submitted" << std::endl;
                    // reBuild
                    std::vector<std::string> file_info = split(file_id, '|');
                    reBuildAndUpload(file_info[0], file_info[1]);
                }
            }
        } else {
            // 没有节点提交证明
            std::cout << "No proof submitted" << std::endl;
        }
    }
}

bool reBuildAndUpload(const std::string& raw_file_id, const std::string& pk) {
    std::cout << "receive the zkp is not enough, rebuild the file on ipfs and chain!" << std::endl;
    std::string contract_address = block_chain->GetAddress("FileStore");
    std::string output;
    std::vector<std::string> params;
    params.push_back(raw_file_id);
    params.push_back(pk);
    bool ret = block_chain->Call("solidity", "FileStore", contract_address, "getFileMeta", params, output);
    if (!ret) {
        std::cout << "storeFile failed" << std::endl;
        return false;
    }
    std::cout << "output: " << output << std::endl;

    rapidjson::Document d;
    d.Parse(output.c_str());

    if (!d.IsArray()) {
        std::cerr << "Not an array" << std::endl;
        return false;
    }

    std::vector<std::string> datas;
    std::vector<std::string> replica;
    std::string file_id = d[0].GetString();
    std::string filename = d[1].GetString();
    int size = 0;
    if (d[2].IsString()) {
        std::string size_str = d[2].GetString();
        size = std::stoi(size_str);
    } else if (d[2].IsInt()) {
        size = d[2].GetInt();
    }

    int shard_count = 0;
    if (d[3].IsString()) {
        std::string shard_count_str = d[3].GetString();
        shard_count = std::stoi(shard_count_str);
    } else if (d[3].IsInt()) {
        shard_count = d[3].GetInt();
    }

    int replica_count = 0;
    if (d[4].IsString()) {
        std::string replica_count_str = d[4].GetString();
        replica_count = std::stoi(replica_count_str);
    } else if (d[4].IsInt()) {
        replica_count = d[4].GetInt();
    }
    std::string replicaCount = d[4].GetString();
    const rapidjson::Value& urls = d[5];
    const rapidjson::Value& cids = d[6];
    for (rapidjson::SizeType i = 0; i < urls.Size(); i++) {
        std::string cid = cids[i].GetString();
        std::string url = urls[i].GetString();
        std::cout << "cid: " << cid << " url: " << url << std::endl;
        std::string tmp_ipfs_file_name = "./tmp/" + filename + std::to_string(i) + ".bin";
        IpfsClient ipfs_client(url);
        ipfs_client.DownloadFile(cid, tmp_ipfs_file_name);
        if (i < replica_count) {
            datas.push_back(tmp_ipfs_file_name);
        } else {
            replica.push_back(tmp_ipfs_file_name);
        }
    }
    std::string user = d[9].GetString();

    std::string tmp_file_name = "./tmp/" + filename;
    auto result_datas = IsalManager::Instance().RecoverFile(
        datas, replica, replica_count, shard_count - replica_count, size, tmp_file_name);
    BackgroundHandler(user, raw_file_id, pk, file_id, filename, tmp_file_name);
    return true;
}

int main() {
    try {
        if (sodium_init() < 0) {
            std::cerr << "libsodium init failed\n";
            return 1;
        }

        // 加载配置文件
        Config::Instance().Init("./config.yml");

        CreateDirIfNotExist("./isal");
        CreateDirIfNotExist("./tmp");

        block_chain = CreateBlockChain(GetBlockChainType(Config::Instance().chain_type()));
        std::vector<std::string> params;
        block_chain->Deploy("solidity", "FileStore", params);
        std::cout << "FileStore address: " << block_chain->GetAddress("FileStore") << std::endl;
        block_chain->Deploy("solidity", "Groth16Verifier", params);
        auto verifier_addr = block_chain->GetAddress("Groth16Verifier");
        std::cout << "Groth16Verifier address: " << block_chain->GetAddress("Groth16Verifier") << std::endl;
        params.push_back(verifier_addr);
        block_chain->Deploy("solidity", "Challenge", params);
        std::cout << "Challenge address: " << block_chain->GetAddress("Challenge") << std::endl;

        std::string challenge_address = block_chain->GetAddress("Challenge");
        Config::Instance().set_challenge_address(challenge_address);

        TaskQueue::Instance().RunWorker(BackgroundHandler);

        std::cout << "start http server" << std::endl;

        // 启动HttpServer监听数据提交请求
        HttpServer::Instance().Init(Config::Instance().http_listen_ip(), Config::Instance().http_listen_port(),
            Config::Instance().http_doc_root());

        std::cout << "start http server successful, listen_ip is " << Config::Instance().http_listen_ip()
                  << " listen_port is " << Config::Instance().http_listen_port() << " file_root_path is "
                  << Config::Instance().http_doc_root() << std::endl;

        std::queue<std::pair<std::vector<std::string>, time_t>> challenges;

        time_t last_create_challenge_time = 0;
        while (true) {
            std::this_thread::sleep_for(std::chrono::seconds(1));
            std::cout
                << "--------------------------------------------------------------------------------------------------"
                << std::endl;
            std::cout << "generate challenge:" << std::endl;
            if (time(0) - last_create_challenge_time > 10 || last_create_challenge_time == 0) {
                std::vector<std::string> cids;
                std::vector<std::string> hashes;
                {
                    std::lock_guard<std::mutex> lock(mutex);
                    for (auto& cid_it : cid_map) {
                        if (cid_it.second == 0 || time(0) - cid_it.second > 240) {
                            cids.push_back(cid_it.first);
                            hashes.push_back(hash_map[cid_it.first]);
                            cid_it.second = time(0);  // 更新挑战时间
                        }
                    }
                }

                if (cids.size() > 0) {
                    std::vector<std::string> challenge_ids;
                    GenerateChallenge(cids, hashes, challenge_ids);
                    std::cout << "challenge_ids: " << challenge_ids.size() << std::endl;
                    last_create_challenge_time = time(0);
                    if (challenge_ids.size() > 0) {
                        challenges.push(std::make_pair(challenge_ids, time(0)));
                    }
                }
            }
            std::cout
                << "--------------------------------------------------------------------------------------------------"
                << std::endl;

            std::cout << "CheckChallenge:" << std::endl;
            CheckChallenge(challenges);
            std::cout
                << "--------------------------------------------------------------------------------------------------"
                << std::endl;
            std::this_thread::sleep_for(std::chrono::seconds(60));
        }
    } catch (const std::exception& e) {
        std::cerr << "Exception: " << e.what() << std::endl;
        return -1;
    }

    std::cout << "exit successfully" << std::endl;
    return 0;
}