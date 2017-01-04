#!/bin/bash

echo
echo "*** REMOVING OLD CRYPTO FILES ***"
echo "rm -rf /var/hyperledger/production/ && mkdir /var/hyperledger/production/"
rm -rf /var/hyperledger/production/ && mkdir /var/hyperledger/production/
sleep 2s 
echo
echo "*** STARTING MEMBERSHIP SERVER ***"
cd $GOPATH/src/github.com/hyperledger/fabric && make membersrvc && membersrvc
