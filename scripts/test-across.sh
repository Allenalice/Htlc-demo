#!/bin/bash

CHANNEL_NAME=mychannel

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'
YELLOW='\033[0;33m'
NC2='\033[4m'
PURPLE='\033[0;35m'

function print_blue() {
  printf "${BLUE}%s${NC}\n" "$1"
}

function print_green() {
  printf "${GREEN}%s${NC}\n" "$1"
}

function print_red() {
  printf "${RED}%s${NC}\n" "$1"
}

function print_yellow() {
  printf "${YELLOW}%s${NC}\n" "$1"
}

function print_yellow2() {
  printf "${YELLOW}%s${NC2}\n" "$1"
}

function print_purple() {
  printf "${PURPLE}%s${NC}\n" "$1"
}


export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")



print_yellow  "*****Print All User's Asset Infomation********"
peer chaincode query -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"GetAllAssets","Args":[]}'


print_purple  "*******Mint 50 Tokens For Alice******\n"
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"Mint","Args":["50","d03edccbb769916b4ad17295e5536de6c9913209aa77f30c0d6a918127e67d5f"]}'

<<comment
print_yellow "******Print Alice's Asset Infomation,You Can See that Money Is 50 Now*******"
peer chaincode query -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"ReadAsset","Args":["38f3e4d1a6327167a93086124ebc9d83cafd08f6c54394ef662d1fd6c5e1c0bd"]}'


print_purple "******Create Htlc With Alice's Address*******"
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"CreateHash","Args":["01","20","123","38f3e4d1a6327167a93086124ebc9d83cafd08f6c54394ef662d1fd6c5e1c0bd"]}'

print_yellow "*****Query Htlc Information By Transaction ID(01) ******"
peer chaincode query -C ${CHANNEL_NAME} -n token_erc20 -c '{"function":"QueryTransId","Args":["01"]}'


print_purple "******Bob Uses Peri<F4><F4>mages To get Alice 20 Token *******"
peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"AcrossTransfer","Args":["123","01","c5c08e75b064ae3fcb832403ff4068da009b432d306d6e9fccf400a20582ceab"]}'

print_yellow "********Query Bob Assets, You Will See Bob's Assets Has Increased By 20 Tokens    *******"
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ReadAsset","Args":["c5c08e75b064ae3fcb832403ff4068da009b432d306d6e9fccf400a20582ceab"]}'


print_yellow "********Query Alice Assets, You Will See Alice's Assets Has Decreased By 20 Tokens    *******"
peer chaincode query -C mychannel -n token_erc20 -c '{"function":"ReadAsset","Args":["38f3e4d1a6327167a93086124ebc9d83cafd08f6c54394ef662d1fd6c5e1c0bd"]}'




comment

