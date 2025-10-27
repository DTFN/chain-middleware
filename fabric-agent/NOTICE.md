1. configtx.yaml里的组织名和mspId不建议有.也不能有_
2. shell脚本里以下cp方式会报文件找不到：
```
cp "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${peer_name}.${org_domain}/tls/tlscacerts/*" .
```
需要改为：
```
cp "${chain_path}/organizations/peerOrganizations/${org_domain}/peers/${peer_name}.${org_domain}/tls/tlscacerts/"* .
```
3. 启动ca-server后需要sleep10-15s
4. docker-compose低版本不支持语法 'networks.dev'，需要从1.25.0到1.29.2，或者yaml里version指定为3.7
5. 通道ID只能包含小写字母、数字和连字符（减号），且不能以数字打头