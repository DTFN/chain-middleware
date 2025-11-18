#include "ipfs_manager.h"
#include <stdexcept>
#include <iostream>

void IpfsManager::Add(const std::string& cid, const std::string& ipfs_address) {
    std::cout << "add ipfs ipfs_address: " << ipfs_address << " cid: " << cid << std::endl;
    auto old_ipfs_address = Get(cid);
    if (old_ipfs_address.empty()) {
        ipfs_map_[cid] = ipfs_address;
    } else if (old_ipfs_address != ipfs_address) {
        throw std::runtime_error("ipfs cid already exists");
    }
}

std::string IpfsManager::Get(const std::string& cid) {
    auto it = ipfs_map_.find(cid);
    if (it != ipfs_map_.end()) {
        return it->second;
    }

    return "";
}
