# 去中心化存储系统（基于零知识证明）

## 🌐 项目简介
本系统是一套 **基于零知识证明机制的去中心化存储系统**，旨在实现可信数据的安全存储与完整性验证。  
系统整体架构由以下核心组成部分构成：
- **存储代理节点（Storage Agent Node）**
- **存储证明节点（Storage Proof Node）**
- **分布式 IPFS 存储集群**
- **区块链层（多链接入）**

系统支持包括：
- **零数链**
- **长安链（ChainMaker）**
- **FISCO BCOS**
- **Fabric**
- **以太坊（Ethereum）**
等五条主流区块链平台的 SDK 接入。

---

## 系统架构概述

```mermaid
flowchart LR
    subgraph User
        U[用户应用]
    end

    subgraph Proxy[存储代理节点]
        F1[文件分片纠删码编码]
        F2[分片上传至IPFS]
        F3[元数据上链]
        F4[挑战任务分配]
    end

    subgraph Proof[存储证明节点]
        P1[挑战监听]
        P2[分片下载]
        P3[零知识证明生成π]
        P4[验证结果上链]
    end

    subgraph IPFS[IPFS 分布式存储集群]
        I1[分片1]
        I2[分片2]
        I3[分片N]
    end

    subgraph Blockchain[多链适配层]
        B1[零数链 SDK]
        B2[长安链 SDK]
        B3[FISCO BCOS SDK]
        B4[Fabric SDK]
        B5[以太坊 SDK]
    end

    U --> Proxy
    Proxy --> IPFS
    Proxy --> Blockchain
    Blockchain --> Proof
    Proof --> Blockchain
