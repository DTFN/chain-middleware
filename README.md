## 用于存放区块链抽象中间件，跨链、链上资源解析、分布式存储，分开单独文件夹存放。

# 项目模块

## contract: 智能合约

跨链功能相关的智能合约，适配以太坊、Hyperledger Fabric和ChainMaker等区块链平台，提供业务处理、DID管理和资源域管理等核心功能。

主要包含以下核心合约：

- **busi** - 业务合约，包含BusiCenter和DIDManager的Java及Solidity实现，提供基础业务逻辑和去中心化身份管理功能
- **busi-fabric** - Hyperledger Fabric平台的Go语言实现，包含配置文件和测试脚本
- **did-manager-fabric** - 针对Fabric平台的DID管理合约，包含相关Go语言工具类
- **did-manager-wasm/chainmaker** - ChainMaker平台的DID管理WASM合约实现
- **resource-domain** - 资源域管理合约，包含ResourceDomain的双语言实现

## tcip-relayer: 跨链中继网关

基于长安链跨链的[TCIP中继网关](https://git.chainmaker.org.cn/chainmaker/tcip-relayer) 进行二次开发，主要用于处理不同区块链网络之间的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

主要包含以下核心功能：

- **配置管理**，管理应用配置和系统参数，支持YAML格式配置文件加载
- **中继链管理**，提供合约调用、查询等核心功能，实现跨链交易的中继处理
- **跨链交易管理**，提供跨链交易处理逻辑，实现跨链交易的发起、确认、取消等操作
- **网关管理**，管理不同区块链网络的网关信息，提供网关注册、查询等功能
- **交易证明**，提供跨链交易证明的生成和验证功能，保障交易完整性

## cross-chain-mgr: 跨链管理服务

跨链管理服务项目是区块链抽象中间件体系的核心组件之一。该项目主要负责提供访问业务链与中继链的功能，通过HTTP接口为用户提供DID注册、资源域名管理和业务合约访问等核心服务。

主要包含以下核心功能：

- **账户管理**，提供账户创建、备份、恢复等功能，实现账户管理
- **合约管理**，提供合约部署、调用、查询等功能，实现智能合约管理
- **跨链交易**，提供跨链交易发起、确认、取消等功能，实现跨链交易管理
- **DID管理**，提供DID注册、查询等功能，实现去中心化身份管理
- **资源域名管理**，提供资源域名注册、查询等功能，实现资源域名管理
- **业务合约访问**，提供业务合约的访问功能，实现业务合约的调用和查询

## fabric-agent: Fabric代理服务

基于Java的区块链中间件项目，主要用于与 Hyperledger Fabric 区块链网络进行交互和管理。

主要包含以下核心功能：
- **链码管理**：提供链码的部署、升级、查询、删除等功能，实现智能合约管理
- **部署管理**：提供区块链网络的部署和配置，实现区块链网络管理
- **节点信息管理**：提供区块链节点信息查询等功能，实现区块链网络管理

## lingshu-ethereum-agent: 零数与以太坊代理服务

基于Java的区块链中间件项目，通过封装LingShuChain SDK与Web3j库和自定义业务逻辑，为跨链操作提供了便捷的区块链操作接口。

主要包含以下核心功能：

- **事件处理**：提供零数链与以太坊链的事件处理功能，实现跨链交易动态解析功能
- **RPC接口**：封装零数链与以太坊链的链交互接口，为跨链操作提供了便捷的区块链操作接口
- **跨链交易管理**：提供零数链与以太坊链的跨链交易处理功能，实现跨链交易管理

## tcip-chainmaker: 长安链跨链网关
基于长安链跨链的[Fabric跨链网关](https://git.chainmaker.org.cn/chainmaker/tcip-chainmaker) 进行二次开发，主要用于处理长安链的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

## tcip-fabric: Fabric跨链网关
基于长安链跨链的[Fabric跨链网关](https://git.chainmaker.org.cn/chainmaker/tcip-fabric) 进行二次开发，主要用于Fabric的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

## tcip-bcos: BCOS跨链网关
基于长安链跨链的[Bcos跨链网关](https://git.chainmaker.org.cn/chainmaker/tcip-bcos) 进行二次开发，主要用于处理Bcos的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

## tcip-ethereum: Ethereum跨链网关
基于长安链跨链的[Bcos跨链网关](https://git.chainmaker.org.cn/chainmaker/tcip-bcos) 进行二次开发，移除Bcos相关模块，通过以太坊的代理服务实现与以太坊私有链的通信。主要用于处理以太坊的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

## tcip-lingshu: 零数链跨链网关
基于长安链跨链的[Bcos跨链网关](https://git.chainmaker.org.cn/chainmaker/tcip-bcos) 进行二次开发，移除Bcos相关模块，通过零数链的代理服务实现与零数链的通信。主要用于处理零数链的跨链交易和通信，提供跨链网关的注册，管理跨链消息的转发和回滚。

# 环境要求：
* 操作系统: x86_64 GNU/Linux
* JDK: Oracle JDK8+或Open JDK8 +
* Docker: docker version 20+
* Go: go1.17.7 linux/amd64
