#!/bin/bash
CHANNEL_NAME=mychannel
. ./color.sh
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")

print_yellow "*****Query Htlc Information By Transaction ID(01) ******"
peer chaincode query -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"QueryTransId","Args":["01"]}'


print_purple "******Bob Uses Perimages To get Alice 20 Token *******"
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"AcrossTransfer","Args":["123","01","a1b52dca768d4bdd2c00e04b25fd1efc4625efe100e07d4de6b768a04751ffda"]}'

