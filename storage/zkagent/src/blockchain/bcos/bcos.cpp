#include "bcos.h"
#include "../common.h"
#include "../eth/eth_tool.h"

#ifdef USE_BCOS
// callback for rpc interfaces
void on_deploy_resp_callback(struct bcos_sdk_c_struct_response* resp) {
    if (resp->error != BCOS_SDK_C_SUCCESS) {
        printf("\t deploy contract failed, error: %d, message: %s\n", resp->error, resp->desc);
        exit(-1);
    }

    const char* cflag = "\"contractAddress\":\"";

    // find the "contractAddress": "0xxxxx"
    char* p0 = strstr((char*)resp->data, cflag);
    if (p0 == NULL) {
        printf("\t cannot find the \"contractAddress\" filed, resp: %s\n", (char*)resp->data);
        exit(-1);
    }

    char* p1 = (char*)p0 + strlen(cflag);
    char* p2 = strstr(p1, "\"");
    if (p2 == NULL) {
        printf("\t cannot find the \"contractAddress\" filed, resp: %s\n", (char*)resp->data);
        exit(-1);
    }

    char* address = (char*)malloc(p2 - p1 + 1);

    memcpy(address, p1, p2 - p1);
    address[p2 - p1] = '\0';

    std::string contract_address = std::string(address);
    std::string contract_name = ((BcosChain*)resp->context)->GetTempContractName();
    ((BcosChain*)resp->context)->BindContractAddress(contract_name, contract_address);
    printf("contractAddress ===>>>>: %s\n", contract_address.c_str());

    // printf("transaction receipt ===>>>>: %s\n", (char*)resp->data);
    if (address != NULL) {
        free(address);
    }
    ((BcosChain*)resp->context)->Notify();
}

void on_send_tx_resp_callback(struct bcos_sdk_c_struct_response* resp) {
    if (resp->error != BCOS_SDK_C_SUCCESS) {
        printf("\t send tx failed, error: %d, message: %s\n", resp->error, resp->desc);
        exit(-1);
    }

    // printf(" ===>> send tx resp: %s\n", (char*)resp->data);
    std::string resp_data = std::string((char*)resp->data);
    ((BcosChain*)resp->context)->ReceiveRespData(resp_data);
    ((BcosChain*)resp->context)->Notify();
}

void on_call_resp_callback(struct bcos_sdk_c_struct_response* resp) {
    if (resp->error != BCOS_SDK_C_SUCCESS) {
        printf("\t call failed, error: %d, message: %s\n", resp->error, resp->desc);
        exit(-1);
    }

    printf(" ===>> call resp: %s\n", (char*)resp->data);
}

// callback for rpc interfaces
void on_recv_resp_callback(struct bcos_sdk_c_struct_response* resp) {
    if (resp->error != 0) {
        printf("\t something is wrong, error: %d, errorMessage: %s\n", resp->error, resp->desc);
        exit(-1);
    } else {
        printf(" \t recv rpc resp from server ===>>>> resp: %s\n", (char*)resp->data);
    }
}
#endif
BcosChain::BcosChain(const std::string& config_path, int32_t ledger_id, const std::string& group_id)
    : config_path_(config_path)
    , ledger_id_(ledger_id)
    , group_id_(group_id) {
#ifdef USE_BCOS
    std::cout << "construct config path: " << config_path_ << std::endl;
    sdk_ = bcos_sdk_create_by_config_file(config_path_.c_str());
    // check success or not
    if (!bcos_sdk_is_last_opr_success()) {
        printf("bcos_sdk_create_by_config_file failed, error: %s\n", bcos_sdk_get_last_error_msg());
        // exit(-1);
    }
    std::cout << "start sdk ... " << std::endl;

    bcos_sdk_start(sdk_);
    // chcek
    if (!bcos_sdk_is_last_opr_success()) {
        printf("bcos_sdk_start failed, error: %s\n", bcos_sdk_get_last_error_msg());
        // exit(-1);
    }

    key_pair_ = bcos_sdk_create_keypair(sdk_crypto_type_);
    if (!key_pair_) {
        printf("create keypair failed, error: %s\n", bcos_sdk_get_last_error_msg());
        exit(-1);
    }

    const char* address = bcos_sdk_get_keypair_address(key_pair_);
    account_address_ = std::string(address);
    bcos_sdk_c_free((void*)address);
    printf("new account, address: %s\n", account_address_.c_str());
    std::cout << "construct sdk success" << std::endl;
#endif
}

BcosChain::~BcosChain() {
#ifdef USE_BCOS
    // stop sdk
    bcos_sdk_stop(sdk_);
    // release sdk
    bcos_sdk_destroy(sdk_);
    // release keypair
    bcos_sdk_destroy_keypair(key_pair_);
#endif
}

bool BcosChain::Deploy(
    const std::string& contract_type, const std::string& contract_name, const std::vector<std::string>& params) {
#ifdef USE_BCOS
    (void)contract_type;
    if (sdk_ == nullptr) {
        return false;
    }
    std::string abi_path = "bcos/solidity/" + contract_name + ".abi";
    std::string bin_path = "bcos/solidity/" + contract_name + ".bin";
    std::string abi = ReadSolFile(abi_path);
    std::string bin = ReadSolFile(bin_path);

    // 4. get chain_id of the group_id
    const char* chain_id = bcos_sdk_get_group_chain_id(sdk_, group_id_.c_str());
    if (!bcos_sdk_is_last_opr_success()) {
        printf("bcos_sdk_get_group_chain_id failed, error: %s\n", bcos_sdk_get_last_error_msg());
        exit(-1);
    }

    // 5. get blocklimit of the group_id.c_str()
    int64_t block_limit = bcos_rpc_get_block_limit(sdk_, group_id_.c_str());
    if (block_limit < 0) {
        printf("group not exist, group: %s\n", group_id_.c_str());
        exit(-1);
    }

    char* tx_hash = nullptr;
    char* signed_tx = nullptr;
    const char* extra_data = "ExtraData";

    printf("extra_data: %s\n", extra_data);
    // 8. deploy HelloWorld contract
    // 8.1 create signed transaction
    bcos_sdk_create_signed_transaction_ver_extra_data(
        key_pair_, group_id_.c_str(), chain_id, "", bin.c_str(), "", block_limit, 0, extra_data, &tx_hash, &signed_tx);

    printf(
        "create deploy contract transaction success, tx_hash: "
        "%s\n",
        tx_hash);
    SetTempContractName(contract_name);
    bcos_rpc_send_transaction(sdk_, group_id_.c_str(), "", signed_tx, 0, on_deploy_resp_callback, this);
    Wait();
    bcos_sdk_c_free((void*)chain_id);
    bcos_sdk_c_free((void*)signed_tx);
    bcos_sdk_c_free((void*)tx_hash);
    return true;
#else
    return false;
#endif
}
bool BcosChain::Call(const std::string contract_type, const std::string contract_name,
    const std::string& contract_address, const std::string& contract_method, const std::vector<std::string>& params,
    std::string& output) {
#ifdef USE_BCOS
    (void)contract_type;

    std::string abi_path = "bcos/solidity/" + contract_name + ".abi";
    std::string bin_path = "bcos/solidity/" + contract_name + ".bin";
    std::string abi = ReadSolFile(abi_path);
    std::string bin = ReadSolFile(bin_path);

    std::string param_str = GenerateJsonParam(params);
    const char* set_data =
        bcos_sdk_abi_encode_method(abi.c_str(), contract_method.c_str(), param_str.c_str(), sdk_crypto_type_);
    const char* chain_id = bcos_sdk_get_group_chain_id(sdk_, group_id_.c_str());
    if (!bcos_sdk_is_last_opr_success()) {
        printf("bcos_sdk_get_group_chain_id failed, error: %s\n", bcos_sdk_get_last_error_msg());
        exit(-1);
    }
    // get blocklimit of the group_id.c_str()
    int64_t block_limit = bcos_rpc_get_block_limit(sdk_, group_id_.c_str());
    if (block_limit < 0) {
        printf("group not exist, group: %s\n", group_id_.c_str());
        exit(-1);
    }
    const char* extra_data = "ExtraData";
    {
        void* transaction_data = bcos_sdk_create_transaction_data(
            group_id_.c_str(), chain_id, contract_address.c_str(), set_data, abi.c_str(), block_limit);
        printf("tx data: %s\n", (char*)transaction_data);
        const char* transaction_data_hash = bcos_sdk_calc_transaction_data_hash(sdk_crypto_type_, transaction_data);
        printf("set tx hash: %s\n", transaction_data_hash);
        const char* signed_hash = bcos_sdk_sign_transaction_data_hash(key_pair_, transaction_data_hash);
        const char* signed_tx = bcos_sdk_create_signed_transaction_with_signed_data_ver_extra_data(
            transaction_data, signed_hash, transaction_data_hash, 0, extra_data);

        bcos_rpc_send_transaction(sdk_, group_id_.c_str(), "", signed_tx, 0, on_send_tx_resp_callback, this);

        Wait();

        bcos_sdk_destroy_transaction_data(transaction_data);
        bcos_sdk_c_free((void*)transaction_data_hash);
        bcos_sdk_c_free((void*)signed_hash);
        bcos_sdk_c_free((void*)signed_tx);
    }
    auto resp_data = GetReceiveRespData();
    // 找到最后一个 '}'
    size_t pos = resp_data.rfind('}');
    if (pos != std::string::npos) {
        resp_data = resp_data.substr(0, pos + 1);  // 保留到最后一个 '}'
    }
    rapidjson::Document d;
    d.Parse(resp_data.c_str());

    if (d.HasParseError()) {
        std::cerr << "JSON parse error!" << std::endl;
        return -1;
    }

    if (d.HasMember("result") && d["result"].HasMember("output")) {
        std::string result_out = d["result"]["output"].GetString();
        if (contract_method == "getChallengeInfo") {
            output = bcos_sdk_abi_decode_method_output(
                abi.c_str(), contract_method.c_str(), result_out.c_str(), sdk_crypto_type_);
        } else if (contract_method == "getChallengesByUser") {
            output = bcos_sdk_abi_decode_method_output(
                abi.c_str(), contract_method.c_str(), result_out.c_str(), sdk_crypto_type_);
            return true;
        } else if (contract_method == "createChallenge") {
            uint64_t challenge_id = HexToUint64(result_out);

            // 创建 JSON 文档
            rapidjson::Document challenge_doc;
            challenge_doc.SetArray();  // 将 doc 设置为数组类型

            rapidjson::Document::AllocatorType& allocator = challenge_doc.GetAllocator();
            challenge_doc.PushBack(challenge_id, allocator);  // 将 uint64_t 添加到数组中

            // 写成字符串
            rapidjson::StringBuffer buffer;
            rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
            challenge_doc.Accept(writer);

            output = buffer.GetString();
        }
    } else {
        std::cerr << "No output found in JSON" << std::endl;
    }
    return true;
#else
    return false;
#endif
}

std::string BcosChain::GetAddress(const std::string& contract_name) {
    auto it = addresses_.find(contract_name);
    if (it == addresses_.end()) {
        return "";
    } else {
        return it->second;
    }
}

std::string BcosChain::GetAccount() { return account_address_; }
