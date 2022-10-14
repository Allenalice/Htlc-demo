#!/bin/bash
JQ_EXEC=`which jq`

FILE_PATH=${PWD}/test-trans.json

len=$(cat $FILE_PATH | ${JQ_EXEC} length | sed 's/\"//g')

export TARGET_TLS_OPTIONS=(-o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt")

for i in $(seq 0 2)
  do
    name=$(cat $FILE_PATH | ${JQ_EXEC} keys[${i}] | sed 's/\"//g')
    #echo $name
    price=$(cat $FILE_PATH | ${JQ_EXEC} .$name[0] | sed 's/\"//g')
    #echo $price
    quantity=$(cat $FILE_PATH | ${JQ_EXEC} .$name[1] | sed 's/\"//g')
    time=$(cat $FILE_PATH | ${JQ_EXEC} .$name[2] | sed 's/\"//g')
    #echo $quantity
   # echo $time
    #echo $money
    #dep=$(echo $price | awk '{printf("%.18f\n"),$0}')
    #money=`awk -v x=${price} -v y=${quantity} 'BEGIN{printf "%.18f\n",x*y}'`
    
   peer chaincode invoke "${TARGET_TLS_OPTIONS[@]}" -C mychannel -n token_erc20 -c '{"function":"CreateAsset","Args":["'${name}'"]}'
    
  done

