#include "chainmaker.h"
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../common.h"
#include "libchainmaker_sdk.h"
#include "rapidjson/prettywriter.h"

// extern GoString deployContractEVM(GoString sdkPath, GoString contractName, GoString byteCodePath);
// extern GoString invokeStoreFile(GoString sdkPath, GoString abiPath, GoString contractName, GoString methodName,
// GoString fileId, GoString fileName, long long int fileSize, long long int totalShards, long long int dataShards,
// GoSlice ipfsUrls, GoSlice poseidonHashes);
// extern GoString invokeGetUser(GoString sdkPath, GoString abiPath, GoString
// contractName, GoString methodName); extern GoUint8 invokeChallengeInit(GoString sdkPath, GoString abiPath, GoString
// contractName, GoString methodName, GoString address);
// extern GoString invokeChallengeGetChallengesByUser(GoString
// sdkPath, GoString abiPath, GoString contractName, GoString methodName);
// extern GoString
// invokeChallengeGetChallengeInfo(GoString sdkPath, GoString abiPath, GoString contractName, GoString methodName,
// GoUint64 challengeId);
// extern GoUint8 invokeChallengeSubmitProof(GoString sdkPath, GoString abiPath, GoString
// contractName, GoString methodName, GoUint64 challengeId, GoString* pA, GoString* pB0, GoString* pB1, GoString* pC,
// GoString* pubSignals); extern GoInt invokeChallengeCreateChallenge(GoString sdkPath, GoString abiPath, GoString
// contractName, GoString methodName, GoString fileHash, GoString nonce, GoString cid, GoString prover);

bool ChainmakerChain::Deploy(
    const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) {
    std::string abi = "chainmaker/testdata/solidity/" + contract_name + ".abi";
    std::string bin = "chainmaker/testdata/solidity/" + contract_name + ".bin";

    GoString sdk_path = {sdk_path_.c_str(), (long)sdk_path_.size()};
    std::string real_contract_name = contract_name + "my_dstorage";
    GoString go_contract_name = {real_contract_name.c_str(), (long)real_contract_name.size()};
    GoString byte_code_path = {bin.c_str(), (long)bin.size()};
    char* address = deployContractEVM(sdk_path, go_contract_name, byte_code_path);
    std::cout << "contract_name : " << contract_name << " address : " << address << std::endl;
    if (contract_name == "Challenge") {
        GoString abi_path = {abi.c_str(), (long)abi.size()};
        std::string method_name = "init";
        GoString go_method_name = {method_name.c_str(), (long)method_name.size()};
        GoString address = {params[0].c_str(), (long)params[0].size()};
        std::cout << "address : " << params[0] << std::endl;
        invokeChallengeInit(sdk_path, abi_path, go_contract_name, go_method_name, address);
    }
    addresses_[contract_name] = std::string(address);
    return true;
}

std::vector<std::string> ConvertJsonToSlice(const std::string& json_str) {
    std::vector<std::string> result;
    rapidjson::Document doc;
    doc.Parse(json_str.c_str());

    if (!doc.IsArray()) {
        std::cerr << "Invalid JSON array" << std::endl;
        return result;
    }

    std::vector<std::string> cppStrings;
    for (auto& v : doc.GetArray()) {
        if (v.IsString()) {
            result.emplace_back(v.GetString());
        }
    }

    return result;
}

bool ChainmakerChain::Call(const std::string contract_type, const std::string contract_name,
    const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
    std::string& output) {
    (void)contract_address;
    GoString sdk_path = {sdk_path_.c_str(), (long)sdk_path_.size()};
    std::string abi = "chainmaker/testdata/solidity/" + contract_name + ".abi";
    GoString abi_path = {abi.c_str(), (long)abi.size()};
    std::string real_contract_name = contract_name + "my_dstorage";
    GoString go_contract_name = {real_contract_name.c_str(), (long)real_contract_name.size()};
    GoString go_method_name = {contract_method.c_str(), (long)contract_method.size()};
    if (contract_method == "createChallenge") {
        GoString file_hash = {params[0].c_str(), (long)params[0].size()};
        GoString nonce = {params[1].c_str(), (long)params[1].size()};
        GoString cid = {params[2].c_str(), (long)params[2].size()};
        GoString prover = {params[3].c_str(), (long)params[3].size()};
        char* result = invokeChallengeCreateChallenge(
            sdk_path, abi_path, go_contract_name, go_method_name, file_hash, nonce, cid, prover);
        output = std::string(result);
    } else if (contract_method == "getChallengeInfo") {
        GoUint64 challengeId = std::stoull(params[0]);
        char* result =
            invokeChallengeGetChallengeInfo(sdk_path, abi_path, go_contract_name, go_method_name, challengeId);
        output = std::string(result);
    } else if (contract_method == "getChallengesByUser") {
        char* result = invokeChallengeGetChallengesByUser(sdk_path, abi_path, go_contract_name, go_method_name);
        output = std::string(result);
    } else if (contract_method == "storeFile") {
        std::string json_str = params[0];
        rapidjson::Document doc;
        doc.Parse(json_str.c_str());

        if (!doc.IsObject() || !doc.HasMember("input") || !doc["input"].IsObject()) {
            std::cerr << "Invalid JSON format" << std::endl;
            return 1;
        }

        const rapidjson::Value& input = doc["input"];

        std::string json_raw_file_id = input["rawFileId"].GetString();
        std::string json_pk = input["pk"].GetString();
        std::string json_author = input["author"].GetString();
        std::string json_file_id = input["fileId"].GetString();
        std::string json_file_name = input["fileName"].GetString();
        int json_file_size = input["fileSize"].GetInt();
        int json_total_shards = input["totalShards"].GetInt();
        int json_data_shards = input["dataShards"].GetInt();

        std::vector<std::string> json_ipfs_urls;
        for (const auto& v : input["ipfsUrls"].GetArray()) {
            json_ipfs_urls.emplace_back(v.GetString());
        }

        std::vector<std::string> json_poseidon_hashes;
        for (const auto& v : input["shardHashes"].GetArray()) {
            json_poseidon_hashes.emplace_back(v.GetString());
        }

        GoString raw_file_id = {json_raw_file_id.c_str(), (long)json_raw_file_id.size()};
        GoString author = {json_author.c_str(), (long)json_author.size()};
        GoString pk = {json_pk.c_str(), (long)json_pk.size()};
        GoString file_id = {json_file_id.c_str(), (long)json_file_id.size()};
        GoString file_name = {json_file_name.c_str(), (long)json_file_name.size()};
        long long int file_size = json_file_size;
        long long int total_shards = json_total_shards;
        long long int data_shards = json_data_shards;

        std::vector<std::string> ipfs_urls = json_ipfs_urls;
        GoString* go_ipfs_urls = new GoString[total_shards];
        for (size_t i = 0; i < total_shards; ++i) {
            go_ipfs_urls[i].p = ipfs_urls[i].c_str();  // 指针直接指向 std::string 的 buffer
            go_ipfs_urls[i].n = ipfs_urls[i].size();   // 长度
        }
        GoSlice slice_ipfs_urls;
        slice_ipfs_urls.data = go_ipfs_urls;
        slice_ipfs_urls.len = ipfs_urls.size();
        slice_ipfs_urls.cap = ipfs_urls.size();

        std::vector<std::string> poseidon_hashes = json_poseidon_hashes;
        GoString* go_poseidon_hashes = new GoString[total_shards];
        for (size_t i = 0; i < total_shards; ++i) {
            go_poseidon_hashes[i].p = poseidon_hashes[i].c_str();  // 指针直接指向 std::string 的 buffer
            go_poseidon_hashes[i].n = poseidon_hashes[i].size();   // 长度
        }

        GoSlice slice_poseidon_hashes;
        slice_poseidon_hashes.data = go_poseidon_hashes;
        slice_poseidon_hashes.len = poseidon_hashes.size();
        slice_poseidon_hashes.cap = poseidon_hashes.size();

        char* result = invokeStoreFile(sdk_path, abi_path, go_contract_name, go_method_name, raw_file_id, pk, file_id,
            file_name, author, file_size);

        std::string contract_name_extent = contract_method + "Extent";
        GoString go_method_name_extent = {contract_name_extent.c_str(), (long)contract_name_extent.size()};
        invokeStoreFileExtent(sdk_path, abi_path, go_contract_name, go_method_name_extent, file_id, total_shards,
            data_shards, slice_ipfs_urls, slice_poseidon_hashes);
        output = std::string(result);
    } else if (contract_method == "getFileMeta") {
        GoString file_id = {params[0].c_str(), (long)params[0].size()};
        GoString pk = {params[1].c_str(), (long)params[1].size()};
        // std::string real_contract_name = "getFileMetaBasic";
        // GoString go_method_name = {real_contract_name.c_str(), (long)real_contract_name.size()};
        char* result = invokeFileStoreGetMetaData(sdk_path, abi_path, go_contract_name, go_method_name, file_id, pk);
        std::string ipfs_info = "getFileMetaIpfsUrls";
        std::string cid_info = "getFileMetaCids";
        GoString ipfs_method = {ipfs_info.c_str(), (long)ipfs_info.size()};
        GoString cids_method = {cid_info.c_str(), (long)cid_info.size()};
        char* result_ipfs = invokeFileStoreGetArray(sdk_path, abi_path, go_contract_name, ipfs_method, file_id, pk);
        std::cout << "ipfs" << result_ipfs << std::endl;
        char* result_cids = invokeFileStoreGetArray(sdk_path, abi_path, go_contract_name, cids_method, file_id, pk);

        rapidjson::Document doc;
        doc.SetArray();
        rapidjson::Document::AllocatorType& allocator = doc.GetAllocator();

        rapidjson::Document receiveResult;
        receiveResult.Parse(result);

        if (!receiveResult.IsArray()) {
            std::cout << "error" << std::endl;
        } else {
            for (auto i = 0; i < receiveResult.Size() - 2; i++) {
                doc.PushBack(receiveResult[i], doc.GetAllocator());
            }
        }

        rapidjson::Document receiveIpfs;
        receiveIpfs.Parse(result_ipfs);

        if (!receiveIpfs.IsArray()) {
            std::cout << "error" << std::endl;
        } else {
            doc.PushBack(receiveIpfs[0], doc.GetAllocator());
        }

        rapidjson::Document receiveCids;
        receiveCids.Parse(result_cids);

        if (!receiveCids.IsArray()) {
            std::cout << "error" << std::endl;
        } else {
            doc.PushBack(receiveCids[0], doc.GetAllocator());
            doc.PushBack(receiveCids[0], doc.GetAllocator());
        }

        if (!receiveIpfs.IsArray()) {
            std::cout << "error" << std::endl;
        } else {
            doc.PushBack(receiveResult[receiveResult.Size() - 2], doc.GetAllocator());
            doc.PushBack(receiveResult[receiveResult.Size() - 1], doc.GetAllocator());
        }

        // 输出 JSON
        rapidjson::StringBuffer buffer;
        rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
        doc.Accept(writer);
        std::cout << "output" << buffer.GetString() << std::endl;
        output = buffer.GetString();
        std::cout << "getFileMetaBasic output " << output << std::endl;
    } else if (contract_method == "hasFile") {
        GoString file_id = {params[0].c_str(), (long)params[0].size()};
        GoString file_name = {params[1].c_str(), (long)params[1].size()};
        GoString file_author = {params[2].c_str(), (long)params[2].size()};
        GoString pk = {params[3].c_str(), (long)params[3].size()};
        char* result =
            invokeHasFile(sdk_path, abi_path, go_contract_name, go_method_name, file_id, file_name, file_author, pk);
        output = std::string(result);
    }
    return true;
}

std::string ChainmakerChain::GetAddress(const std::string& contract_name) { return addresses_[contract_name]; }