## 用于存放区块链抽象中间件，跨链、链上资源解析、分布式存储，分开单独文件夹存放。

# 项目结构
| 模块名      |     功能描述 |       说明       |
|:---------:|:---------:|:--------------|
| contract | 合约项目 | 代币合约，DID管理合约，链上资源域名管理合约。 |
| tcip-relayer | 跨链中继服务 | 跨链网关的注册，管理跨链消息的转发和回滚。 |
| cross-chain-mgr | 跨链管理服务 | 服务同时访问业务链和中继链，通过HTTP接口提供DID注册、资源域名管理和业务合约访问的功能。 |
| fabric-agent | Fabric 代理服务 | 提供HTTP接口发起fabric合约交易。 |
| tcip-chainmaker | 长安链跨链网关 | 连接长安链和中继网关，提供区块链管理，跨链事件管理，跨链交易动态解析的功能。 |
| tcip-fabric | Fabric 跨链网关 | 连接Hyperledger Fabric链和中继网关，提供区块链管理，跨链事件管理，跨链交易动态解析的功能。 |
| tcip-bcos | BCOS 跨链网关 | 连接BCOS链和中继网关，提供区块链管理，跨链事件管理，跨链交易动态解析的功能。 |
| tcip-ethereum | Ethereum 跨链网关 | 连接Ethereum和中继网关，提供区块链管理，跨链事件管理，跨链交易动态解析的功能。 |
| tcip-lingshu | 零数链跨链网关 | 连接零数链和中继网关，提供区块链管理，跨链事件管理，跨链交易动态解析的功能。 |
| front | 零数链和Ethereum的代理服务 | 通过HTTP接口提供零数链和Ethereum的访问功能。 |

# 环境要求：
* 操作系统: x86_64 GNU/Linux
* JDK: Oracle JDK8+或Open JDK8 +
* Docker: docker version 20+
