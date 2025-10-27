fabric_busi_chain_code=BUSI_FABRIC_6769


export FABRIC_CFG_PATH=${project_path}/test/test_fabric_didmanager/config-peer1
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_CLIENTAUTHREQUIRED=true
export CORE_PEER_TLS_CLIENTCERT_FILE=${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.crt
export CORE_PEER_TLS_CLIENTKEY_FILE=${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/server.key


# ed25519 成功
cd ${project_path}/contract/Busi-fabric
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com --tls \
  --cafile ${project_path}/hyperledger-fabric/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n ${fabric_busi_chain_code} \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"VerifyVcs","Args":["{\"origin\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"BUSI_FABRIC_1460_11278\",\"contract_func\":\"Erc20GetBalanceVcs\",\"contract_type\":\"7\",\"func_params\":{\"account\":\"did:lsid:1:1b5d5be8c5abad92145726a82474f52f46a8d9c2e6cb43e91826d481a0c0ee6e\"},\"gateway_id\":\"2\",\"param_name\":\"vcs\",\"resource_name\":\"chain_fabric_01:BUSI_FABRIC_1460_11278\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-09-18\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:51:05.460520501+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"42137b00493908fe99e6232a123dfb0c98322cabf5c9b83ab31418054f6482dc2b1a9afa7e2d925fbaa30a605dc04357991d1f9f6f57aaaa8a779dc01a298b01\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"contract_type\":\"1\",\"func_params\":{\"account\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:51:05.476427760+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"23903c55893d23ca41a76da78feb3d62a12a2f3b4014cdc8f5f5c268c6f6d3e9e1bcfc4ad3ecec803647a68e37c4c4bee8394d0211c7b92f47a5b69a28b82307\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"]}'


# ed25519 失败
cd ${project_path}/contract/Busi-fabric
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com --tls \
  --cafile ${project_path}/hyperledger-fabric/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n ${fabric_busi_chain_code} \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"VerifyVcs","Args":["{\"origin\":{\"@context\":[\"https://www.w3.org/2019/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"BUSI_FABRIC_1460_11278\",\"contract_func\":\"Erc20GetBalanceVcs\",\"contract_type\":\"7\",\"func_params\":{\"account\":\"did:lsid:1:1b5d5be8c5abad92145726a82474f52f46a8d9c2e6cb43e91826d481a0c0ee6e\"},\"gateway_id\":\"2\",\"param_name\":\"vcs\",\"resource_name\":\"chain_fabric_01:BUSI_FABRIC_1460_11278\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-09-18\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:51:05.460520501+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"42137b00493908fe99e6232a123dfb0c98322cabf5c9b83ab31418054f6482dc2b1a9afa7e2d925fbaa30a605dc04357991d1f9f6f57aaaa8a779dc01a298b01\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"contract_type\":\"1\",\"func_params\":{\"account\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:51:05.476427760+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"23903c55893d23ca41a76da78feb3d62a12a2f3b4014cdc8f5f5c268c6f6d3e9e1bcfc4ad3ecec803647a68e37c4c4bee8394d0211c7b92f47a5b69a28b82307\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"]}'

# eth 成功
cd ${project_path}/contract/Busi-fabric
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com --tls \
  --cafile ${project_path}/hyperledger-fabric/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n ${fabric_busi_chain_code} \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"VerifyVcs","Args":["{\"origin\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"BUSI_FABRIC_1460_11278\",\"contract_func\":\"Erc20GetBalanceVcs\",\"contract_type\":\"7\",\"func_params\":{\"account\":\"did:lsid:1:1b5d5be8c5abad92145726a82474f52f46a8d9c2e6cb43e91826d481a0c0ee6e\"},\"gateway_id\":\"2\",\"param_name\":\"vcs\",\"resource_name\":\"chain_fabric_01:BUSI_FABRIC_1460_11278\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-09-18\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\",\"contentHash\":\"d65fac3ed51f907cf0e68248155494c872af38467763cd39f0d02ae75ae502ee\",\"r\":\"769d1f8cde79266041b0f191487e2c65b7801ec6586ea118b0a089fec80bdf1b\",\"s\":\"2ac0fe88053f89caf58d0eb0e8758f8d204319039587b501c5d9cc5fec162923\",\"v\":28,\"verificationMethod\":\"did:1#eth\"},\"type\":[\"VerifiableCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"contract_type\":\"1\",\"func_params\":{\"account\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:53:02.905232658+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"23903c55893d23ca41a76da78feb3d62a12a2f3b4014cdc8f5f5c268c6f6d3e9e1bcfc4ad3ecec803647a68e37c4c4bee8394d0211c7b92f47a5b69a28b82307\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"]}'

# eth 失败
cd ${project_path}/contract/Busi-fabric
peer chaincode invoke \
  -o localhost:7050 \
  --ordererTLSHostnameOverride orderer.example.com --tls \
  --cafile ${project_path}/hyperledger-fabric/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
  -C mychannel \
  -n ${fabric_busi_chain_code} \
  --peerAddresses localhost:7051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
  --peerAddresses localhost:9051 \
  --tlsRootCertFiles ${project_path}/hyperledger-fabric/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt \
  -c '{"function":"VerifyVcs","Args":["{\"origin\":{\"@context\":[\"https://www.w3.org/2019/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain_fabric_01\",\"contract_address\":\"BUSI_FABRIC_1460_11278\",\"contract_func\":\"Erc20GetBalanceVcs\",\"contract_type\":\"7\",\"func_params\":{\"account\":\"did:lsid:1:1b5d5be8c5abad92145726a82474f52f46a8d9c2e6cb43e91826d481a0c0ee6e\"},\"gateway_id\":\"2\",\"param_name\":\"vcs\",\"resource_name\":\"chain_fabric_01:BUSI_FABRIC_1460_11278\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-09-18\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"address\":\"000000000000000000000000bd5d9f6b1b70cfd36c90aa11a8830a43062ab4ea\",\"contentHash\":\"d65fac3ed51f907cf0e68248155494c872af38467763cd39f0d02ae75ae502ee\",\"r\":\"769d1f8cde79266041b0f191487e2c65b7801ec6586ea118b0a089fec80bdf1b\",\"s\":\"2ac0fe88053f89caf58d0eb0e8758f8d204319039587b501c5d9cc5fec162923\",\"v\":28,\"verificationMethod\":\"did:1#eth\"},\"type\":[\"VerifiableCredential\"]},\"target\":{\"@context\":[\"https://www.w3.org/2018/credentials/v1\"],\"credentialSubject\":{\"chain_rid\":\"chain1\",\"contract_address\":\"contract01\",\"contract_func\":\"stock\",\"contract_type\":\"1\",\"func_params\":{\"account\":\"did:1\"},\"gateway_id\":\"0\",\"param_name\":\"content\",\"resource_name\":\"chain1:contract01\"},\"id\":\"https://uni.example.com/credentials/12345\",\"issuanceDate\":\"2025-05-20\",\"issuer\":\"https://uni.example.com/organization/1\",\"proof\":{\"created\":\"2025-09-23T14:53:02.905232658+08:00\",\"proofPurpose\":\"assertionMethod\",\"publicKey\":\"702fa6cd4f9133217448450067710d3f9ed0ec740edc8337e30245ba19b07d4a\",\"signature\":\"23903c55893d23ca41a76da78feb3d62a12a2f3b4014cdc8f5f5c268c6f6d3e9e1bcfc4ad3ecec803647a68e37c4c4bee8394d0211c7b92f47a5b69a28b82307\",\"type\":\"JsonWebSignature2020\",\"verificationMethod\":\"did:1#ed25519\"},\"type\":[\"VerifiableCredential\",\"UniversityDegreeCredential\"],\"z\":\"1\"}}"]}'
