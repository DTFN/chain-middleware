#include "fabric.h"
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../common.h"
#include "libfabric_sdk.h"
#include "rapidjson/prettywriter.h"

// extern void initConfig(GoString m, GoString c, GoString p, GoString g, GoString cn, GoString ccn);
// mspID        = m
// cryptoPath   = c
// peerEndpoint = p
// gatewayPeer  = g
// channelName  = cn
// chaincodeName = ccn
// extern char* invokeStoreFile(GoString contractName, GoString methodName, GoString fileId, GoString fileName, long
// long int fileSize, long long int totalShards, long long int dataShards, GoSlice ipfsUrls, GoSlice poseidonHashes);
// extern char* invokeGetUser(GoString contractName, GoString methodName);
// extern char* invokeChallengeGetChallengesByUser(GoString contractName, GoString methodName);
// extern char* invokeChallengeGetChallengeInfo(GoString contractName, GoString methodName, GoUint64 challengeId);
// extern char* invokeChallengeSubmitProof(GoString contractName, GoString methodName, GoUint64 challengeId, GoString
// proofCID, GoString proofURL, GoString proofHash, GoString publicCID, GoString publicURL, GoString publicHash); extern
// char* invokeChallengeGetProof(GoString contractName, GoString methodName, GoUint64 challengeId); extern char*
// invokeChallengeUpdateChallenge(GoString contractName, GoString methodName, GoUint64 challengeId, GoUint8
// ProofVerified); extern char* invokeChallengeCreateChallenge(GoString contractName, GoString methodName, GoString
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

inline std::string ParseMyChallenge(std::string json) {
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
    auto issuer = d["issuer"].GetString();
    auto cid = d["cid"].GetString();
    auto proof_submitted = d["proof_submitted"].GetBool();
    auto proof_verified = d["proof_verified"].GetBool();
    auto prover = d["prover"].GetString();

    rapidjson::Document output(rapidjson::kArrayType);
    rapidjson::Document::AllocatorType& allocator = output.GetAllocator();
    // 按顺序加入元素
    output.PushBack(rapidjson::Value(challenge_id.c_str(), allocator), allocator);
    output.PushBack(rapidjson::Value(issuer, allocator), allocator);
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
        output = std::string(result);
    } else if (contract_method == "submitProof") {
        std::string real_contract_method = "SubmitProof";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        GoUint64 challenge_id = std::stoull(params[0]);
        GoString proof_cid = {params[1].c_str(), (long)params[1].size()};
        GoString proof_url = {params[2].c_str(), (long)params[2].size()};
        GoString proof_hash = {params[3].c_str(), (long)params[3].size()};
        GoString public_cid = {params[4].c_str(), (long)params[4].size()};
        GoString public_url = {params[5].c_str(), (long)params[5].size()};
        GoString public_hash = {params[6].c_str(), (long)params[6].size()};
        invokeFabricChallengeSubmitProof(go_contract_name, go_method_name, challenge_id, proof_cid, proof_url,
            proof_hash, public_cid, public_url, public_hash);
    } else if (contract_method == "getChallengeInfo") {
        std::string real_contract_method = "GetChallengeInfo";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        GoUint64 challenge_id = std::stoull(params[0]);
        char* result = invokeFabricChallengeGetChallengeInfo(go_contract_name, go_method_name, challenge_id);
        output = std::string(result);
    } else if (contract_method == "getChallengesByUser") {
        std::string real_contract_method = "GetMyChallenges";
        GoString go_method_name = {real_contract_method.c_str(), (long)real_contract_method.size()};
        char* result = invokeFabricChallengeGetChallengesByUser(go_contract_name, go_method_name);
        output = ParseMyChallenge(std::string(result));
    }
    return true;
}

std::string FabricChain::GetAddress(const std::string& contract_name) { return addresses_[contract_name]; }

std::string FabricChain::GetAccount() {
    std::string contract_name = "FileStoreContract";
    std::string contract_method = "GetUser";
    GoString go_contract_name = {contract_name.c_str(), (long)contract_name.size()};
    GoString go_method_name = {contract_method.c_str(), (long)contract_method.size()};
    char* account_address = invokeFabricGetUser(go_contract_name, go_method_name);

    std::string account = std::string(account_address);
    return account;
    // return "";
}
