#include "user.h"
#include <iostream>

bool User::HasUserInfo(const std::string& user_id) {
    std::lock_guard<std::mutex> guard(mutex_);
    return user_map_.find(user_id) != user_map_.end();
}

bool User::AddUserInfo(
    const std::string& user_id, const std::string& pwd, const std::string& pk, const std::string& sk) {
    std::lock_guard<std::mutex> guard(mutex_);
    auto it = user_map_.find(user_id);
    if (it == user_map_.end()) {
        user_map_[user_id] = pwd;
        user_pk_[user_id] = pk;
        user_sk_[user_id + "|" + pwd] = sk;
        user_pk_id_[pk] = user_id;
        return true;
    } else {
        std::cout << "add user info failed" << std::endl;
        return false;
    }
}

std::string User::GetPk(const std::string& user_id) {
    std::lock_guard<std::mutex> guard(mutex_);
    return user_pk_[user_id];
}

std::string User::GetSk(const std::string& user_id, const std::string& pwd) {
    std::lock_guard<std::mutex> guard(mutex_);
    return user_sk_[user_id + "|" + pwd];
}

std::pair<std::string, std::string> User::GetKeyPair(const std::string& user_id, const std::string& pwd) {
    std::lock_guard<std::mutex> guard(mutex_);
    if (user_map_.find(user_id) == user_map_.end()) {
        return std::make_pair("", "");
    }
    std::string pk = user_pk_[user_id];
    std::string sk = user_sk_[user_id + "|" + pwd];
    return std::make_pair(pk, sk);
}
