#pragma once

#include "blockchain.h"
#include "chainmaker/chainmaker.h"
#include "ls/ls_chain.h"
#include "eth/eth.h"
#include "fabric/fabric.h"
#include "bcos/bcos.h"

inline std::shared_ptr<BlockChain> CreateBlockChain(BlockChainType type) {
    std::shared_ptr<BlockChain> block_chain = nullptr;
    if (type == BlockChainType::ls) {
        block_chain = CreateLsChain();
    } else if (type == BlockChainType::chainmaker) {
        block_chain = CreateChainmakerChain();
    } else if (type == BlockChainType::eth) {
        block_chain = CreateEthChain();
    } else if (type == BlockChainType::fabric) {
        block_chain = CreateFabricChain();
    } else if (type == BlockChainType::bcos) {
        block_chain = CreateBcosChain();
    } else {
        throw std::runtime_error("block chain type error");
    }

    return block_chain;
}

inline BlockChainType GetBlockChainType(const std::string& str_type) {
    if (str_type == "ls") {
        return BlockChainType::ls;
    } else if (str_type == "chainmaker") {
        return BlockChainType::chainmaker;
    } else if (str_type == "bcos") {
        return BlockChainType::bcos;
    }  else if (str_type == "eth") {
        return BlockChainType::eth;
    }  else if (str_type == "fabric") {
        return BlockChainType::fabric;
    } else {
        throw std::runtime_error("block chain type error");
    }
}
