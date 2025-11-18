#include "http_server.h"
#include <cstring>
#include <fstream>
#include <iostream>
#include <sodium.h>
#include <rapidjson/document.h>
#include <rapidjson/stringbuffer.h>
#include <rapidjson/writer.h>
#include "../blockchain/blockchain.h"
#include "../config/config.h"
#include "../event/task_queue.h"
#include "../ipfs/ipfs_client.h"
#include "../isal/isal.h"
#include "../user/user.h"
#include "CivetServer.h"
#include "civetweb.h"
#include "rapidjson/prettywriter.h"
#include "tools/tools.h"

extern std::shared_ptr<BlockChain> block_chain;
const std::string kSalt = "per-user-salt-or-store-offchain";  // ideally a random per-user salt
void NotifyTaskQueue(const std::string &user, const std::string &raw_file_id, const std::string &pk,
    const std::string &file_id, const std::string &file_name, const std::string &full_file_path) {
    std::cout << "[+] Notifying task queue: " << file_name << "\n";
    FileStruct file_info;
    file_info.user = user;
    file_info.raw_file_id = raw_file_id;
    file_info.pk = pk;
    file_info.file_id = file_id;
    file_info.file_name = file_name;
    file_info.full_fule_path = full_file_path;
    TaskQueue::Instance().Push(file_info);
}

void GrantAccessFile(const std::string &user, std::string &file_id) {}

bool CanAccessFile(const std::string &user, const std::string &file_id) {
    std::string contract_address = block_chain->GetAddress("FileAccess");
    std::vector<std::string> params;

    return true;
}

std::vector<unsigned char> DownFile(
    const std::string &raw_file_id, const std::string &pk, int size, std::string &filename) {
    std::vector<std::string> params;
    params.push_back(raw_file_id);
    params.push_back(pk);
    std::string output;
    std::string contract_address = block_chain->GetAddress("FileStore");
    bool ret = block_chain->Call("solidity", "FileStore", contract_address, "getFileMeta", params, output);
    if (!ret) {
        std::cout << "download failed" << std::endl;
        return std::vector<unsigned char>();
    }
    std::cout << "output: " << output << std::endl;

    rapidjson::Document d;
    d.Parse(output.c_str());

    if (!d.IsArray()) {
        std::cerr << "Not an array" << std::endl;
        return std::vector<unsigned char>();
    }

    std::vector<std::string> datas;
    std::vector<std::string> replica;
    filename = d[1].GetString();
    if (filename == "") {
        std::cout << "filename is null, download failed" << std::endl;
        return std::vector<unsigned char>();
    }
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
    const rapidjson::Value &urls = d[5];
    const rapidjson::Value &cids = d[6];
    for (rapidjson::SizeType i = 0; i < urls.Size(); i++) {
        std::string cid = cids[i].GetString();
        std::string url = urls[i].GetString();
        std::cout << "cid: " << cid << " url: " << url << std::endl;
        std::string file_name = "./tmp/" + filename + std::to_string(i) + ".bin";
        IpfsClient ipfs_client(url);
        ipfs_client.DownloadFile(cid, file_name);
        if (i < replica_count) {
            datas.push_back(file_name);
        } else {
            replica.push_back(file_name);
        }
    }

    return IsalManager::Instance().RecoverFile(
        datas, replica, replica_count, shard_count - replica_count, size, "./tmp/" + filename);
}

class UploadHandler : public CivetHandler {
public:
    bool handlePost(CivetServer *server, struct mg_connection *conn) override {
        std::cout << "handle post" << std::endl;
        char buf[1024];
        int n;

        const mg_context *ctx = server->getContext();
        const char *doc_root = mg_get_option(ctx, "document_root");
        if (!doc_root) doc_root = ".";  // fallback

        const struct mg_request_info *req_info = mg_get_request_info(conn);
        std::string query = req_info->query_string ? req_info->query_string : "";

        std::string file_name = "uploaded_file";
        size_t pos = query.find("file_name=");
        if (pos != std::string::npos) {
            size_t start = pos + 10;
            size_t end = query.find("&", start);
            file_name = query.substr(start, end - start);
        }

        std::string user = "";
        pos = query.find("user=");
        if (pos != std::string::npos) {
            size_t start = pos + 5;
            size_t end = query.find("&", start);
            user = query.substr(start, end - start);
        }

        if (user == "") {
            std::cout << "has no user info" << std::endl;
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "has no user info\n");
            return true;
        }
        std::string pk = User::Instance().GetPk(user);

        std::vector<unsigned char> data;
        while ((n = mg_read(conn, buf, sizeof(buf))) > 0) {
            data.insert(data.end(), buf, buf + n);
        }
        std::string raw_file_id = "";
        HashData(data, raw_file_id);

        std::vector<unsigned char> ciphertext(data.size() + crypto_box_SEALBYTES);
        bool ret = Encrypt(ciphertext, data, pk);
        if (!ret) {
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "Failed encrypt file faild\n");
            return true;
        }
        // std::vector<unsigned char> target(ciphertext.size() - crypto_box_SEALBYTES);
        // Decrypt(target, ciphertext, pk, sk);
        std::string file_id = "";
        ret = HashData(ciphertext, file_id);
        std::cout << "file " << file_name << " hash is " << file_id << std::endl;
        if (!ret) {
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "Failed cal hash faild\n");
            return true;
        }
        std::string file_full_path = std::string(doc_root) + "/" + file_id;
        std::ofstream out(file_full_path, std::ios::binary);
        if (!out.is_open()) {
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "Failed to open output file\n");
            return true;
        }
        out.write(reinterpret_cast<const char *>(ciphertext.data()), ciphertext.size());
        out.close();
        std::cout << "content write to file successful, data size is " << ciphertext.size() << std::endl;

        NotifyTaskQueue(user, raw_file_id, pk, file_id, file_name, file_full_path);

        mg_printf(conn,
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: text/plain\r\n"
            "Connection: close\r\n\r\n"
            "{\"file_id\": \"%s\"}\n",
            raw_file_id.c_str());

        return true;
    }
};

class JsonHandler : public CivetHandler {
public:
    bool handlePost(CivetServer *server, struct mg_connection *conn) override {
        const char *contentType = mg_get_header(conn, "Content-Type");
        if (!contentType || std::string(contentType).find("application/json") == std::string::npos) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nUnsupported Content-Type");
            return true;
        }

        char buffer[4096];
        int ret = mg_read(conn, buffer, sizeof(buffer) - 1);
        buffer[ret] = '\0';
        std::string body(buffer);

        try {
            rapidjson::Document doc;
            if (doc.Parse(body.c_str()).HasParseError()) {
                mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
                return true;
            }

            rapidjson::StringBuffer strBuf;
            rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuf);
            doc.Accept(writer);
            std::string account = doc["account"].GetString();
            std::string ipfs_address = doc["ipfs_address"].GetString();
            std::cout << "Received JSON:\n" << strBuf.GetString() << std::endl;
            std::cout << "account: " << account << " ipfs_address:  " << ipfs_address << std::endl;
            HttpServer::Instance().AddAccount(ipfs_address, account);
        } catch (...) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
            return true;
        }

        mg_printf(conn, "HTTP/1.1 200 OK\r\nContent-Type: application/json\r\n\r\n{\"challenge_address\": \"%s\"}",
            Config::Instance().challenge_address().c_str());
        return true;
    }
};

class UserRegisterHandler : public CivetHandler {
public:
    bool handlePost(CivetServer *server, struct mg_connection *conn) override {
        const char *contentType = mg_get_header(conn, "Content-Type");
        if (!contentType || std::string(contentType).find("application/json") == std::string::npos) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nUnsupported Content-Type");
            return true;
        }

        char buffer[4096];
        int ret = mg_read(conn, buffer, sizeof(buffer) - 1);
        buffer[ret] = '\0';
        std::string body(buffer);

        try {
            rapidjson::Document doc;
            if (doc.Parse(body.c_str()).HasParseError()) {
                mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
                return true;
            }

            rapidjson::StringBuffer strBuf;
            rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuf);
            doc.Accept(writer);
            std::string user = doc["user"].GetString();
            std::string pwd = doc["pwd"].GetString();
            std::cout << "Received JSON:\n" << strBuf.GetString() << std::endl;
            std::cout << "user: " << user << " pwd:  " << pwd << std::endl;
            if (User::Instance().HasUserInfo(user)) {
                mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nUser has already exists");
                return true;
            }

            // 1. Derive seed
            unsigned char seed[crypto_sign_SEEDBYTES];
            if (!DeriveSeed(user, pwd, kSalt, seed)) {
                mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nderive_seed failed");
                std::cout << "derive_seed failed" << std::endl;
                return true;
            }
            std::cout << "Derived seed (hex): " << CharToHex(seed, sizeof(seed)) << std::endl;

            // 2. Generate Ed25519 keypair
            unsigned char pk[crypto_box_PUBLICKEYBYTES];
            unsigned char sk[crypto_box_SECRETKEYBYTES];
            crypto_box_seed_keypair(pk, sk, seed);
            std::string pk_str = CharToHex(pk, sizeof(pk));
            std::string sk_str = CharToHex(sk, sizeof(sk));
            std::cout << "public key: " << pk_str << std::endl;
            std::cout << "secret key: " << sk_str << " (keep secret!)" << std::endl;

            std::string response = "{\"public_key\":\"" + pk_str + "\", \"secret_key\":\"" + sk_str + "\"}";

            User::Instance().AddUserInfo(user, pwd, pk_str, sk_str);

            mg_printf(conn,
                "HTTP/1.1 200 OK\r\n"
                "Content-Type: application/json\r\n"
                "Content-Length: %lu\r\n"
                "\r\n%s",
                response.length(), response.c_str());
        } catch (...) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
            return true;
        }
        return true;
    }
};

class ShareToOtherHandler : public CivetHandler {
public:
    bool handlePost(CivetServer *server, struct mg_connection *conn) override {
        const char *contentType = mg_get_header(conn, "Content-Type");
        if (!contentType || std::string(contentType).find("application/json") == std::string::npos) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nUnsupported Content-Type");
            return true;
        }

        char buffer[4096];
        int ret = mg_read(conn, buffer, sizeof(buffer) - 1);
        buffer[ret] = '\0';
        std::string body(buffer);

        try {
            rapidjson::Document doc;
            if (doc.Parse(body.c_str()).HasParseError()) {
                mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
                return true;
            }

            rapidjson::StringBuffer strBuf;
            rapidjson::PrettyWriter<rapidjson::StringBuffer> writer(strBuf);
            doc.Accept(writer);
            std::string user = doc["user"].GetString();
            std::string pwd = doc["pwd"].GetString();

            std::string to_user = doc["to_user"].GetString();
            std::string to_pk = User::Instance().GetPk(to_user);
            std::string file_id = doc["file_id"].GetString();
            std::cout << "Received JSON:\n" << strBuf.GetString() << std::endl;
            std::cout << "account: " << user << " pwd:  " << pwd << std::endl;

            std::string user_sk = User::Instance().GetSk(user, pwd);
            int size = 0;
            std::string filename = "";
            auto user_pk = User::Instance().GetPk(user);
            auto result_datas = DownFile(file_id, user_pk, size, filename);
            if (result_datas.size() == 0) {
                mg_printf(conn,
                    "HTTP/1.1 500 Internal Server Error\r\n"
                    "Content-Type: text/plain\r\n"
                    "Connection: close\r\n\r\n"
                    "the file is not exist!\n");
                return true;
            }

            std::vector<unsigned char> target(result_datas.size() - crypto_box_SEALBYTES);
            Decrypt(target, result_datas, user_pk, user_sk);
            std::vector<unsigned char> ciphertext(target.size() + crypto_box_SEALBYTES);
            Encrypt(ciphertext, target, to_pk);

            std::string new_file_id = "";
            ret = HashData(ciphertext, new_file_id);
            if (!ret) {
                mg_printf(conn,
                    "HTTP/1.1 500 Internal Server Error\r\n"
                    "Content-Type: text/plain\r\n"
                    "Connection: close\r\n\r\n"
                    "Failed cal hash faild\n");
                return true;
            }
            std::cout << "file " << filename << " hash is " << new_file_id << std::endl;
            const mg_context *ctx = server->getContext();
            const char *doc_root = mg_get_option(ctx, "document_root");
            if (!doc_root) doc_root = ".";  // fallback
            std::string file_full_path = std::string(doc_root) + "/" + new_file_id;
            std::ofstream out(file_full_path, std::ios::binary);
            if (!out.is_open()) {
                mg_printf(conn,
                    "HTTP/1.1 500 Internal Server Error\r\n"
                    "Content-Type: text/plain\r\n"
                    "Connection: close\r\n\r\n"
                    "Failed to open output file\n");
                return true;
            }

            out.write(reinterpret_cast<const char *>(ciphertext.data()), ciphertext.size());
            out.close();

            NotifyTaskQueue(user, file_id, to_pk, new_file_id, filename, file_full_path);

            std::string response = "{\"file_id\":\"" + file_id + "\"}";

            mg_printf(conn,
                "HTTP/1.1 200 OK\r\n"
                "Content-Type: application/json\r\n"
                "Content-Length: %lu\r\n"
                "\r\n%s",
                response.length(), response.c_str());
        } catch (...) {
            mg_printf(conn, "HTTP/1.1 400 Bad Request\r\n\r\nInvalid JSON");
            return true;
        }
        return true;
    }
};

class ContractAddressHandler : public CivetHandler {
public:
    bool handleGet(CivetServer *server, struct mg_connection *conn) override {
        auto addr = block_chain->GetAddress("Challenge");
        std::string response = "{\"challenge_contract_address\":\"" + addr + "\"}";

        mg_printf(conn,
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: application/json\r\n"
            "Content-Length: %lu\r\n"
            "\r\n%s",
            response.length(), response.c_str());

        return true;
    }
};

class DownloadHandler : public CivetHandler {
public:
    bool handleGet(CivetServer *server, struct mg_connection *conn) override {
        // 获取上行参数的file_id
        char file_id[128] = {0};
        char user[128] = {0};
        char pwd[128] = {0};
        const struct mg_request_info *req_info = mg_get_request_info(conn);
        if (req_info->query_string) {
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "file_id", file_id, sizeof(file_id));
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "user", user, sizeof(user));
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "password", pwd, sizeof(pwd));
        }
        std::string file_id_str(file_id);
        std::string user_str(user);

        {
            std::string contract_address = block_chain->GetAddress("FileStore");
            std::vector<std::string> params;
            std::string pk = User::Instance().GetPk(user_str);

            params.push_back(file_id_str);
            params.push_back("");
            params.push_back("");
            params.push_back(pk);
            std::string output;
            bool ret = block_chain->Call("solidity", "FileStore", contract_address, "hasFile", params, output);
            if (!ret) {
                std::cout << "qryFile failed" << std::endl;
            }
            std::cout << "output: " << output << std::endl;

            rapidjson::Document d;
            d.Parse(output.c_str());

            if (!d.IsArray()) {
                std::cout << "output: " << output << std::endl;

                std::cerr << "Not an array" << std::endl;
                return false;
            }

            // 转成 id-true 对象数组
            rapidjson::Value result(rapidjson::kArrayType);
            for (rapidjson::SizeType i = 0; i < d[0].Size(); i++) {
                if (d[0][i].GetString() == file_id_str && !d[1][i].GetBool()) {
                    mg_printf(conn,
                        "HTTP/1.1 500 Internal Server Error\r\n"
                        "Content-Type: text/plain\r\n"
                        "Connection: close\r\n\r\n"
                        "the user no access to the file\n");
                    return true;
                }
            }
        }

        int size = 0;
        std::string file_name = "";
        auto key_pair = User::Instance().GetKeyPair(user_str, pwd);
        auto result_datas = DownFile(file_id_str, key_pair.first, size, file_name);
        if (result_datas.size() == 0) {
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "the file is not exist!\n");
            return true;
        }
        if (key_pair.first == "" || key_pair.second == "") {
            mg_printf(conn,
                "HTTP/1.1 500 Internal Server Error\r\n"
                "Content-Type: text/plain\r\n"
                "Connection: close\r\n\r\n"
                "can not found the user info\n");
            return true;
        }
        std::vector<unsigned char> target(result_datas.size() - crypto_box_SEALBYTES);
        Decrypt(target, result_datas, key_pair.first, key_pair.second);
        mg_printf(conn,
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: application/octet-stream\r\n"
            "Content-Length: %d\r\n"
            "Content-Disposition: attachment; filename=\"%s\"\r\n"
            "\r\n",
            (int)target.size(), file_id_str.c_str());

        // 一次性写内存里的数据
        if (!target.empty()) {
            mg_write(conn, target.data(), target.size());
        }

        return true;
    }
};

class QryFileHandler : public CivetHandler {
public:
    bool handleGet(CivetServer *server, struct mg_connection *conn) override {
        std::string contract_address = block_chain->GetAddress("FileStore");
        std::vector<std::string> params;
        // 获取上行参数的file_id
        char file_id[128] = {0};
        char user[128] = {0};
        const struct mg_request_info *req_info = mg_get_request_info(conn);
        if (req_info->query_string) {
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "file_id", file_id, sizeof(file_id));
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "user", user, sizeof(user));
        }
        std::string file_id_str(file_id);
        std::string pk = User::Instance().GetPk(user);

        params.push_back(file_id_str);
        params.push_back(pk);
        std::string output;
        bool ret = block_chain->Call("solidity", "FileStore", contract_address, "getFileMeta", params, output);
        if (!ret) {
            std::cout << "qryFile failed" << std::endl;
        }
        // std::cout << "output: " << output << std::endl;

        rapidjson::Document d;
        d.Parse(output.c_str());

        if (!d.IsArray()) {
            std::cout << "output: " << output << std::endl;

            std::cerr << "Not an array" << std::endl;
            return false;
        }

        std::vector<std::string> replica;
        std::string filename = d[1].GetString();
        std::string size = "";
        if (d[2].IsString()) {
            size = d[2].GetString();
        } else if (d[2].IsInt()) {
            size = std::to_string(d[2].GetInt());
        }
        std::string shardCount = "";
        if (d[3].IsString()) {
            size = d[3].GetString();
        } else if (d[3].IsInt()) {
            size = std::to_string(d[3].GetInt());
        }
        std::string replicaCount = "";
        if (d[4].IsString()) {
            size = d[4].GetString();
        } else if (d[4].IsInt()) {
            size = std::to_string(d[4].GetInt());
        }
        const rapidjson::Value &urls = d[5];
        const rapidjson::Value &cids = d[6];
        int exist_count = 0;
        for (rapidjson::SizeType i = 0; i < urls.Size(); i++) {
            std::string cid = cids[i].GetString();
            std::string url = urls[i].GetString();
            // std::cout << "cid: " << cid << " url: " << url << std::endl;
            std::string file_name = "./tmp/" + filename + std::to_string(i) + ".bin";
            IpfsClient ipfs_client(url);
            bool is_exist = ipfs_client.IsExist(cid);
            if (is_exist) {
                exist_count++;
                if (exist_count >= Config::Instance().isal_threshold()) {
                    break;
                }
            }
        }
        std::string response = "";
        if (exist_count >= Config::Instance().isal_threshold()) {
            response = "{\"result\":\"true\"}";
        } else {
            response = "{\"result\":\"false\"}";
        }

        mg_printf(conn,
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: application/json\r\n"
            "Content-Length: %lu\r\n"
            "\r\n%s",
            response.length(), response.c_str());

        return true;
    }
};

class QryFileExistsHandler : public CivetHandler {
public:
    bool handleGet(CivetServer *server, struct mg_connection *conn) override {
        std::cout << "QryFileExistsHandler" << std::endl;
        std::string contract_address = block_chain->GetAddress("FileStore");
        std::vector<std::string> params;
        // 获取上行参数的file_id
        char file_hash[128] = {0};
        char file_name[256] = {0};
        char author[128] = {0};
        char user[128] = {0};
        const struct mg_request_info *req_info = mg_get_request_info(conn);
        if (req_info->query_string) {
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "file_id", file_hash, sizeof(file_hash));
            mg_get_var(
                req_info->query_string, strlen(req_info->query_string), "file_name", file_name, sizeof(file_name));
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "author", author, sizeof(author));
            mg_get_var(req_info->query_string, strlen(req_info->query_string), "user", user, sizeof(user));
        }
        std::string file_hash_str(file_hash);
        std::string file_name_str(file_name);
        std::string author_str(author);
        std::string pk = User::Instance().GetPk(user);

        params.push_back(file_hash_str);
        params.push_back(file_name_str);
        params.push_back(author_str);
        params.push_back(pk);
        std::string output;
        bool ret = block_chain->Call("solidity", "FileStore", contract_address, "hasFile", params, output);
        if (!ret) {
            std::cout << "qryFile failed" << std::endl;
        }
        std::cout << "output: " << output << std::endl;

        rapidjson::Document d;
        d.Parse(output.c_str());

        if (!d.IsArray()) {
            std::cout << "output: " << output << std::endl;

            std::cerr << "Not an array" << std::endl;
            return false;
        }

        // 转成 id-true 对象数组
        rapidjson::Value result(rapidjson::kArrayType);
        for (rapidjson::SizeType i = 0; i < d[0].Size(); i++) {
            rapidjson::Value obj(rapidjson::kObjectType);
            obj.AddMember("id", d[0][i], d.GetAllocator());
            obj.AddMember("access", d[1][i], d.GetAllocator());
            result.PushBack(obj, d.GetAllocator());
        }

        // 输出 JSON
        rapidjson::StringBuffer buffer;
        rapidjson::Writer<rapidjson::StringBuffer> writer(buffer);
        result.Accept(writer);

        std::string response = buffer.GetString();
        std::cout << response << std::endl;

        mg_printf(conn,
            "HTTP/1.1 200 OK\r\n"
            "Content-Type: application/json\r\n"
            "Content-Length: %lu\r\n"
            "\r\n%s\r\n",
            response.length(), response.c_str());

        return true;
    }
};

void HttpServer::Init(const std::string &listen_ip, uint16_t listen_port, const std::string &doc_root) {
    std::string listen_host = listen_ip + ":" + std::to_string(listen_port);
    const char *options[] = {"document_root", doc_root.c_str(),  // 静态文件根目录
        "listening_ports", listen_host.c_str(),                  // 监听端口
        "num_threads", "64",                                    // 开启 8 个工作线程
        nullptr};

    static CivetServer server(options);
    static UploadHandler handler;
    static JsonHandler json_handler;
    static UserRegisterHandler user_register_handler;
    static ContractAddressHandler contract_addr_handler;
    static DownloadHandler download_handler;
    static QryFileHandler qry_handler;
    static QryFileExistsHandler qry_file_exists_handler;
    static ShareToOtherHandler share_handler;

    server.addHandler("/upload", handler);
    server.addHandler("/register", json_handler);
    server.addHandler("/user_register", user_register_handler);
    server.addHandler("/contract/address", contract_addr_handler);
    server.addHandler("/download", download_handler);
    server.addHandler("/qryfile", qry_handler);
    server.addHandler("/qryfileexists", qry_file_exists_handler);
    server.addHandler("/share", share_handler);
}

std::string HttpServer::GetAccount(const std::string &ipfs_address) {
    auto it = account_map_.find(ipfs_address);
    if (it != account_map_.end()) {
        return it->second;
    }

    return "";
}

void HttpServer::AddAccount(const std::string &ipfs_address, const std::string &account) {
    auto it = account_map_.find(ipfs_address);
    if (it != account_map_.end()) {
        it->second = account;
        return;
    }

    account_map_[ipfs_address] = account;
}
