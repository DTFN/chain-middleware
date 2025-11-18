#pragma once

#include <string>

bool submitProofOnChain(const std::string& cid, const std::string& proofFile,
                        const std::string& publicFile);

bool generateProof();