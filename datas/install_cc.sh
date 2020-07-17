#!/bin/bash

if [ ! -f "version" ]; then echo "0" > version;  fi

VERSION=`more version`
COMMAND=""

if [[ $VERSION > "0" ]]; then
	COMMAND="upgrade"
else
	COMMAND="instantiate"
fi

((VERSION++))

echo "$VERSION" > version
NAME_CC="datas"

docker exec -e "CORE_PEER_ADDRESS=peer.bank1.kz:7051" -e "CORE_PEER_LOCALMSPID=Bank1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank1.kz/users/Admin@bank1.kz/msp" cli peer chaincode install -n $NAME_CC -v $VERSION -l golang -p github.com/chaincode/$NAME_CC
docker exec -e "CORE_PEER_ADDRESS=peer.bank2.kz:7051" -e "CORE_PEER_LOCALMSPID=Bank2MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank2.kz/users/Admin@bank2.kz/msp" cli peer chaincode install -n $NAME_CC -v $VERSION -l golang -p github.com/chaincode/$NAME_CC
docker exec -e "CORE_PEER_ADDRESS=peer.bank3.kz:7051" -e "CORE_PEER_LOCALMSPID=Bank3MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank3.kz/users/Admin@bank3.kz/msp" cli peer chaincode install -n $NAME_CC -v $VERSION -l golang -p github.com/chaincode/$NAME_CC
docker exec -e "CORE_PEER_ADDRESS=peer.bank1.kz:7051" -e "CORE_PEER_LOCALMSPID=Bank1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/bank1.kz/users/Admin@bank1.kz/msp" cli peer chaincode $COMMAND -o orderer.kz:7050  -C channel1 -n $NAME_CC -l golang -v $VERSION -c '{"Args":["init"]}' 