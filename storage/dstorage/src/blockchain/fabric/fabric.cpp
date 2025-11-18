#include "fabric.h"
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../common.h"
#include "libfabric_sdk.h"
#include "rapidjson/prettywriter.h"
#include "tools/tools.h"

// extern void initConfig(GoString m, GoString c, GoString p, GoString g, GoString cn, GoString ccn);
// mspID        = m
// cryptoPath   = c
// peerEndpoint = p
// gatewayPeer  = g
// channelName  = cn
// chaincodeName = ccn
// extern char* invokeFabricStoreFile(GoString contractName, GoString methodName, GoString fileId, GoString fileName,
// long long int fileSize, long long int totalShards, long long int dataShards, GoSlice ipfsUrls, GoSlice
// poseidonHashes); extern char* invokeFabricGetUser(GoString contractName, GoString methodName); extern char*
// invokeFabricChallengeGetChallengesByUser(GoString contractName, GoString methodName); extern char*
// invokeFabricChallengeGetChallengeInfo(GoString contractName, GoString methodName, GoUint64 challenge_id); extern
// char* invokeFabricChallengeSubmitProof(GoString contractName, GoString methodName, GoUint64 challenge_id, GoString
// proofCID, GoString proofURL, GoString proofHash, GoString publicCID, GoString publicURL, GoString publicHash);
// extern char* invokeFabricChallengeGetProof(GoString contractName, GoString methodName, GoUint64 challenge_id); extern
// char* invokeFabricChallengeUpdateChallenge(GoString contractName, GoString methodName, GoUint64 challenge_id, GoUint8
// ProofVerified);
// extern char* invokeFabricChallengeCreateChallenge(GoString contractName, GoString methodName, GoString
// fileHash, GoString nonce, GoString cid, GoString prover);

void FabricChain::Init() {
    const std::string& msp_id = Config::Instance().msp_id();
    const std::string& crypto_path = Config::Instance().crypto_path();
    const std::string& cert_path = Config::Instance().cert_path();
    const std::string& key_path = Config::Instance().key_path();
    const std::string& tls_cert_path = Config::Instance().tls_cert_path();
    const std::string& peer_endpoint = Config::Instance().peer_endpoint();
    const std::string& gateway_peer = Config::Instance().gateway_peer();
    const std::string& channel_name = Config::Instance().channel_name();
    const std::string& chaincode_name = Config::Instance().chaincode_name();
    GoString m = {msp_id.c_str(), (long)msp_id.size()};
    GoString c = {crypto_path.c_str(), (long)crypto_path.size()};
    GoString cp = {cert_path.c_str(), (long)cert_path.size()};
    GoString kp = {key_path.c_str(), (long)key_path.size()};
    GoString tlscp = {tls_cert_path.c_str(), (long)tls_cert_path.size()};
    GoString p = {peer_endpoint.c_str(), (long)peer_endpoint.size()};
    GoString g = {gateway_peer.c_str(), (long)gateway_peer.size()};
    GoString cn = {channel_name.c_str(), (long)channel_name.size()};
    GoString ccn = {chaincode_name.c_str(), (long)chaincode_name.size()};
    initFabricConfig(m, c, p, g, cn, ccn, cp, kp, tlscp);
}

bool FabricChain::Deploy(
    const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) {
    (void)contract_type;
    (void)contract_name;
    (void)params;
    return true;
}

inline std::string ParseFileMetaDataJson(std::string json) {
    std::string result = "";
    rapidjson::Document d;
    rapidjson::ParseResult ok = d.Parse(json.c_str());
    if (!ok) {
        std::string err = std::string("JSON parse error");
        std::cout << err << std::endl;
        return result;
    }
    if (!d.IsObject()) {
        std::string err = std::string("Root must be an object");
        std::cout << err << std::endl;
        return result;
    }

    auto uploader = d["uploader"].GetString();
    auto filename = d["fileName"].GetString();
    auto size = std::to_string(d["fileSize"].GetInt());
    auto shardCount = std::to_string(d["totalShards"].GetInt());
    auto replicaCount = std::to_string(d["dataShards"].GetInt());
    auto urls = d["ipfsUrls"].GetArray();
    auto cids = d["cids"].GetArray();

    rapidjson::Document output(rapidjson::kArrayType);
    rapidjson::Document::AllocatorType& allocator = output.GetAllocator();
    // 按顺序加入元素
    output.PushBack(rapidjson::Value(uploader, allocator), allocator);
    output.PushBack(rapidjson::Value(filename, allocator), allocator);
    output.PushBack(rapidjson::Value(size.c_str(), allocator), allocator);
    output.PushBack(rapidjson::Value(shardCount.c_str(), allocator), allocator);
    output.PushBack(rapidjson::Value(replicaCount.c_str(), allocator), allocator);
    output.PushBack(urls, allocator);
    output.PushBack(cids, allocator);

    // 输出结果
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    output.Accept(writer);
    result = buffer.GetString();

    return result;
}

inline std::string ParseChallengeInfoJson(std::string json) {
    std::string result = "";
    rapidjson::Document d;
    rapidjson::ParseResult ok = d.Parse(json.c_str());
    if (!ok) {
        std::string err = std::string("JSON parse error");
        std::cout << err << std::endl;
        return result;
    }
    if (!d.IsObject()) {
        std::string err = std::string("Root must be an object");
        std::cout << err << std::endl;
        return result;
    }

    auto challenge_id = std::to_string(d["challenge_id"].GetInt());
    auto file_hash = d["file_hash"].GetString();
    auto nonce = d["nonce"].GetString();
    auto cid = d["cid"].GetString();
    auto proof_submitted = d["proof_submitted"].GetBool();
    auto proof_verified = d["proof_verified"].GetBool();
    auto prover = d["prover"].GetString();

    rapidjson::Document output(rapidjson::kArrayType);
    rapidjson::Document::AllocatorType& allocator = output.GetAllocator();
    // 按顺序加入元素
    output.PushBack(rapidjson::Value(challenge_id.c_str(), allocator), allocator);
    output.PushBack(rapidjson::Value(file_hash, allocator), allocator);
    output.PushBack(rapidjson::Value(nonce, allocator), allocator);
    output.PushBack(rapidjson::Value(cid, allocator), allocator);
    output.PushBack("", allocator);
    output.PushBack(proof_submitted, allocator);
    output.PushBack(proof_verified, allocator);
    output.PushBack(rapidjson::Value(prover, allocator), allocator);

    // 输出结果
    rapidjson::StringBuffer buffer;
    rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
    output.Accept(writer);
    result = buffer.GetString();

    return result;
}

bool FabricChain::Call(const std::string contract_type, const std::string contract_name,
    const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
    std::string& output) {
    (void)contract_address;
    std::string real_contract_name = contract_name;
    if (contract_name == "FileStore") {
        real_contract_name = "FileStoreContract";
    } else if (contract_name == "Challenge") {
        real_contract_name = "ChallengeContract";
    }
    GoString go_contract_name = {real_contract_name.c_str(), (long)real_contract_name.size()};
    if (contract_method == "createChallenge") {
        std::string real_contract_method = "CreateChallenge";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        GoString file_hash = {params[0].c_str(), (long)params[0].size()};
        GoString nonce = {params[1].c_str(), (long)params[1].size()};
        GoString cid = {params[2].c_str(), (long)params[2].size()};
        GoString prover = {params[3].c_str(), (long)params[3].size()};
        char* result =
            invokeFabricChallengeCreateChallenge(go_contract_name, go_method_name, file_hash, nonce, cid, prover);

        rapidjson::Document output_json(rapidjson::kArrayType);
        rapidjson::Document::AllocatorType& allocator = output_json.GetAllocator();
        // 按顺序加入元素
        output_json.PushBack(rapidjson::Value(result, allocator), allocator);

        // 输出结果
        rapidjson::StringBuffer buffer;
        rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
        output_json.Accept(writer);
        output = buffer.GetString();
    } else if (contract_method == "getChallengeInfo") {
        std::string real_contract_method = "GetProof";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        // 这里还需要验证零知识证明是不是对的，需要下载challenge并本地验证，验证后结果上传
        // {"challenge_id":1,"proof_cid":"proofCID","proof_url":"proofURL","proof_hash":"proofHash","public_cid":"publicCID","public_url":"publicURL","public_hash":"publicHash","vk_cid":"","vk_url":"","vk_hash":""}
        GoUint64 challenge_id = std::stoull(params[0]);
        char* proof = invokeFabricChallengeGetProof(go_contract_name, go_method_name, challenge_id);
        std::cout << "challengeID:" << params[0] << " proof:" << proof  << std::endl;
        auto verify_res = ParseChallengeJsonAndVerify(std::string(proof));
        real_contract_method = "UpdateChallenge";
        GoString go_update_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        invokeFabricChallengeUpdateChallenge(go_contract_name, go_update_method_name, challenge_id, verify_res);

        real_contract_method = "GetChallenge";
        GoString go_get_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        char* result = invokeFabricChallengeGetChallengeInfo(go_contract_name, go_get_method_name, challenge_id);

        output = ParseChallengeInfoJson(std::string(result));
    } else if (contract_method == "getFileMeta") {
        std::string real_contract_method = "GetFileMeta";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        GoString go_file_id = {params[0].c_str(), (long)params[0].size()};
        GoString go_pk = {params[1].c_str(), (long)params[1].size()};
        char* result = invokeFabricFileStoreGetMetaData(go_contract_name, go_method_name, go_file_id, go_pk);
        output = ParseFileMetaDataJson(std::string(result));
    } else if (contract_method == "hasFile") {
        std::string real_contract_method = "HasFile";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        GoString go_file_id = {params[0].c_str(), (long)params[0].size()};
        GoString go_file_name = {params[1].c_str(), (long)params[1].size()};
        GoString author = {params[2].c_str(), (long)params[2].size()};
        GoString go_pk = {params[3].c_str(), (long)params[3].size()};
        char* result = invokeFabricHasFile(go_contract_name, go_method_name, go_file_id, go_file_name, author, go_pk);
        std::cout << "has file result : " << std::string(result) << std::endl;
        output = std::string(result);
    } else if (contract_method == "storeFile") {
        std::string real_contract_method = "StoreFile";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
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

        std::vector<std::string> json_cids;
        for (const auto& v : input["shardHashes"].GetArray()) {
            json_cids.emplace_back(v.GetString());
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

        std::vector<std::string> cids = json_cids;
        GoString* go_poseidon_hashes = new GoString[total_shards];
        for (size_t i = 0; i < total_shards; ++i) {
            go_poseidon_hashes[i].p = cids[i].c_str();  // 指针直接指向 std::string 的 buffer
            go_poseidon_hashes[i].n = cids[i].size();   // 长度
        }

        GoSlice slice_cids;
        slice_cids.data = go_poseidon_hashes;
        slice_cids.len = cids.size();
        slice_cids.cap = cids.size();

        std::cout << "invokeFabricStoreFile begin" << std::endl;
        char* result = invokeFabricStoreFile(go_contract_name, go_method_name, raw_file_id, pk, file_id, file_name, author, file_size,
            total_shards, data_shards, slice_ipfs_urls, slice_cids);
        std::cout << "invokeFabricStoreFile end" << std::endl;
        output = std::string(result);
    }
    return true;
}

std::string FabricChain::GetAddress(const std::string& contract_name) { return addresses_[contract_name]; }