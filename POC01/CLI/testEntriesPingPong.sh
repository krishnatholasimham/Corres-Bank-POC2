#!/bin/bash

# Clear DB
echo "CLEAR DB"
peer chaincode invoke -l golang -n anz -c '{"Function":"deleteRange", "Args": ["0","zzzz"]}'
sleep 1s

# Add 10000 Funding
echo
echo "ADDING 10000 FUNDING"
date="$(TZ=$TZ date '+%FT%T.%N%:z')"
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addLedgerEntryFunding\", \"Args\": [\"ANZ\",\"WF\",\"10000\",\"$date\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Add 1000 Request
echo
echo "ADD 1000 REQUEST"
date="$(TZ=$TZ date '+%FT%T.%N%:z')"
DAY=$(date -d "$date" '+%d')
MONTH=$(date -d "$date" '+%m')
YEAR=$(date -d "$date" '+%Y')
vDate="${YEAR:2}$MONTH$DAY"
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"1000\",\"ANZ\",\"WF\",\"$date\",\"OUR\",\"$vDate\",\"12345\",\"USD\",\"$localDate\",\"001\",\"103\",\"WF\"]}'"
    #  peer chaincode invoke -l golang -n anz -c '{"Function":"addPaymentInstruction", "Args": ["1000","Alice","ANZ","WF","Bob","2016-07-19T13:33:12.435929075+00:00","BEN"]}'
echo $cmd && eval $cmd
sleep 1s

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt

# Transaction Summary
echo
echo "PRINT TXN SUMMARY"
peer chaincode query -l golang -n anz -c '{"Function":"getTransactionSummary", "Args": ["ANZ"]}'

# Confirm Request
echo
echo "CONFIRM REQUEST"
# read -p "Please provide Request ID: " ID
date="$(TZ=$TZ date '+%FT%T.%N%:z')";
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"42624451a0028a178146ac5358595d38c9778a77dae6e1615a97f9c32e0d10ab\",\"$date\",\"WF\",\"$localDate\",\"001c\",\"103\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Add 2000 request
echo
echo "ADD 2000 REQUEST"
date="$(TZ=$TZ date '+%FT%T.%N%:z')";
DAY=$(date -d "$date" '+%d')
MONTH=$(date -d "$date" '+%m')
YEAR=$(date -d "$date" '+%Y')
vDate="${YEAR:2}$MONTH$DAY"
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
# echo $vDate
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"2000\",\"ANZ\",\"WF\",\"$date\",\"SHA\",\"$vDate\",\"123456\",\"USD\",\"$localDate\",\"002\",\"103\",\"WF\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt

# Transaction Summary
echo
echo "PRINT TXN SUMMARY"
peer chaincode query -l golang -n anz -c '{"Function":"getTransactionSummary", "Args": ["ANZ"]}'

# Confirm Request
echo
echo "CONFIRM REQUEST"
# read -p "Please provide Request ID: " ID
date="$(TZ=$TZ date '+%FT%T.%N%:z')";
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"76b833827a84ba3f8e85991debbbbd692537fd63acd68ebe0a8eb686834012e6\",\"$date\",\"WF\",\"$localDate\",\"002c\",\"103\"]}'"
echo $cmd && eval $cmd
sleep 1s

# ADD 3000 REQUEST
echo
echo "ADDING 3000 REQUEST"
date="$(TZ=$TZ date '+%FT%T.%N%:z')"
DAY=$(date -d "$date" '+%d')
MONTH=$(date -d "$date" '+%m')
YEAR=$(date -d "$date" '+%Y')
vDate="$YEAR$MONTH$DAY"
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"3000\",\"ANZ\",\"WF\",\"$date\",\"BEN\",\"$vDate\",\"123457\",\"USD\",\"$localDate\",\"003\",\"103\",\"WF\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt

# Transaction Summary
echo
echo "PRINT TXN SUMMARY"
peer chaincode query -l golang -n anz -c '{"Function":"getTransactionSummary", "Args": ["ANZ"]}'

# Confirm Request
echo
echo "CONFIRM REQUEST"
# read -p "Please provide Request ID: " ID
date="$(TZ=$TZ date '+%FT%T.%N%:z')";
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"f1fae93836c5cd3dec0c5890a5648822292ef011fbd1f3eb4c8238d4ef02d064\",\"$date\",\"WF\",\"$localDate\",\"003c\",\"103\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt

# Transaction Summary
echo
echo "PRINT TXN SUMMARY"
peer chaincode query -l golang -n anz -c '{"Function":"getTransactionSummary", "Args": ["ANZ"]}'

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt

# Confirm Request
echo
echo "CONFIRM REQUEST"
# read -p "Please provide Request ID: " ID
date="$(TZ=$TZ date '+%FT%T.%N%:z')";
localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
cmd="peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"9e863680fc6493872770bf332e904d31179317f1a1ad8258d120569e30017b88\",\"$date\",\"WF\",\"$localDate\",\"1\",\"103\"]}'"
echo $cmd && eval $cmd
sleep 1s

# Transaction Summary
echo
echo "PRINT TXN SUMMARY"
peer chaincode query -l golang -n anz -c '{"Function":"getTransactionSummary", "Args": ["ANZ"]}'

# Get Unconfirmed Statement Balance & Pipe into balance.txt
echo
echo "GET UNCONFIRMED STATEMENT BALANCE"
cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
echo $cmd && eval $cmd > balance.txt


# # Get Statement Balance
# echo "GET STATEMENT BALANCE"
# cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getBalanceHistory\", \"Args\": [\"ANZ\",\"WF\"]}'"
# echo $cmd && eval $cmd
