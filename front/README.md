# 节点前置服务

## 简介
front是和LingShuChain节点配合使用的一个子系统，需要和节点统计部署，可和节点进行通信，集成了java-sdk，对接口进行了封装和抽象。

# 项目结构
| 模块名      |     功能描述 |       说明       |
|:---------|---------:|:--------------:|
| docker   |     镜像脚本 | 依赖LingShuChain |
| lib      |  依赖的链SDK |       无        |
| script   |     运维脚本 |       无        |
| src/main | 前置节点业务代码 |       无        |
| src/test | 前置节点测试代码 |       无        |

# 主要源码目录
| 目录          | 功能描述                      |
|:------------|:--------------------------|
| abi         | 维护合约abi文件中所包含的的方法信息       |
| base        | 基础类,包含配置、枚举、异常等           |
| contract    | 合约接口，包含合约部署、合约编译等         |
| event       | 链上事件通知信息                  |
| gm          | 获取加密类型                    |
| link        | 对已运行的区块链纳入平台进行统一管控        |
| monitor     | 监控链上数据,检查节点进程连接           |
| performance | 前置节点性能,tps/平均耗时/异常比率      |
| builtin     | 内置API，包含权限管理，节点管理，NS管理等信息 |
| task        | 定时任务,事件同步任务               |
| tool        | 工具接口                      |
| transaction | 交易信息，包含交易查询、交易发送等         |
| util        | 工具类，包含错误码工具类、证书工具类等       |
| version     | 版本信息，查询前置服务版本与签名服务版本      |
| rpcapi      | rpc接口，与链节点数据交互的接口         |

# 环境要求：
* 操作系统: x86_64 GNU/Linux
* JDK: Oracle JDK8+或Open JDK8 +
* Docker: docker version 20+
* Aliyun Docker: 阿里云docker账号
* LingShuChainSDK: 2.2.6