#!/bin/bash

CHANNEL_NAME=mychannel
. ./color.sh

export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")

print_yellow "******Print Alice's Asset Infomation,You Can See that Money Is 50 Now*******"
peer chaincode query -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"ReadAsset","Args":["d03edccbb769916b4ad17295e5536de6c9913209aa77f30c0d6a918127e67d5f"]}'


print_purple "******Create Htlc With Alice's Address*******"
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"CreateHash","Args":["01","20","123","d03edccbb769916b4ad17295e5536de6c9913209aa77f30c0d6a918127e67d5f"]}'

