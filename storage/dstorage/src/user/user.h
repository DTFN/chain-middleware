#pragma once

#include <map>
#include <string>
#include <mutex>

class User {
public:
    static User& Instance() {
        static User u;
        return u;
    }

    bool HasUserInfo(const std::string& user_id);

    bool AddUserInfo(const std::string& user_id, const std::string& pwd, const std::string& pk, const std::string& sk);
    std::string GetPk(const std::string& user_id);
    std::string GetSk(const std::string& user_id, const std::string& pwd);

    std::pair<std::string, std::string> GetKeyPair(const std::string& user_id, const std::string& pwd);
private:
    // user_id, user_pwd
    std::map<std::string, std::string> user_map_;
    // user_id, pk
    std::map<std::string, std::string> user_pk_;
    // user_id|pwd, sk
    std::map<std::string, std::string> user_sk_;

    // pk, user_id
    std::map<std::string, std::string> user_pk_id_;

    std::mutex mutex_;
};
