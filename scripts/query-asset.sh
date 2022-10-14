#!/bin/bash

. ./color.sh
export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")


print_yellow "********Query Bob Assets, You Will See Bob's Assets Has Increased By 20 Tokens    *******"
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ReadAsset","Args":["a1b52dca768d4bdd2c00e04b25fd1efc4625efe100e07d4de6b768a04751ffda"]}'


print_yellow "********Query Alice Assets, You Will See Alice's Assets Has Decreased By 20 Tokens    *******"
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ReadAsset","Args":["d03edccbb769916b4ad17295e5536de6c9913209aa77f30c0d6a918127e67d5f"]}'

