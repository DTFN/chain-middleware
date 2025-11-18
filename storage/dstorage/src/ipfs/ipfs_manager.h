#pragma once

#include <map>
#include <string>
#include "boost/noncopyable.hpp"

class IpfsManager : public boost::noncopyable {
public:
    static IpfsManager& Instance() {
        static IpfsManager instance;
        return instance;
    }

    void Add(const std::string& cid, const std::string& ipds_address);
    std::string Get(const std::string& cid);

private:
    std::map<std::string, std::string> ipfs_map_;
};