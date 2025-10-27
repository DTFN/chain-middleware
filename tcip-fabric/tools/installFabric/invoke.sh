#!/usr/bin/env bash

cd fabric-samples/test-network
#在fabric-samples/test-network下调用合约的方法:
export FABRIC_CFG_PATH=$PWD/../config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051
../bin/peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
-C mychannel \
-n crosschain1 \
--peerAddresses localhost:7051 \
--tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c \
'{"function":"CrossChainTransfer","Args":["test","test","{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain1\",\"method\":\"CrossChainConfirm\"}","{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain1\",\"method\":\"Cross_ChainCancel\"}","[{\"gateway_id\": \"0\",\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"CrossChainTry\",\"parameter\":\"[\\\"method\\\",\\\"CrossChainTry\\\"]\",\"extra_data\":\"按需写，目标网关能解析就行\",\"confirm_info\":{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"CrossChainConfirm\"},\"cancel_info\":{\"chain_id\":\"mychannel\",\"contract_name\":\"crosschain2\",\"method\":\"CrossChainCancel\"}}]"]}'
../bin/peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
-C mychannel \
-n crosschain1 \
--peerAddresses localhost:7051 \
--tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c \
'{"function":"Query","Args":["1","1"]}'
../bin/peer chaincode invoke \
-o localhost:7050 \
--ordererTLSHostnameOverride orderer.example.com --tls --cafile ${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem \
-C mychannel \
-n crosschain2 \
--peerAddresses localhost:7051 \
--tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt \
--peerAddresses localhost:9051 --tlsRootCertFiles ${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c \
'{"function":"Query","Args":["1","1"]}'