#pragma once

#include <set>
#include <string>
#include <vector>

// std::set<std::string> listPinnedCIDs();

bool downloadFromIPFS(const std::string& endpoint, const std::string& cid, const std::string& outPath);

std::string uploadFileToIpfs(const std::string& endpoint, const std::string& file_path);