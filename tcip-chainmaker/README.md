# TCIP中继网关

## 版本对照表
| chainmaker-go|tcip-relay|tcip-chainmaker|tcip-fabric|fabric-bcos|
| -------- | --------- | --------- | --------- | --------- |
| v2.3.1 | v2.3.1| v2.3.1| v2.3.1| v2.3.1|

## 编译安装包

```shell
make release
```

## 编译镜像

```shell
make docker_build
```

根据文档修改release下的配置文件

## 启动

```shell
cd release
./switch.sh up
```

## 关闭

```shell
cd release
./switch.sh down
```