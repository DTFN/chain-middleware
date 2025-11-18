#pragma once


#include <vector>
#include <fstream>
#include <string>

std::vector<uint8_t> readFile(const std::string& filePath);

std::string computePoseidonHash(const std::vector<uint8_t>& data);