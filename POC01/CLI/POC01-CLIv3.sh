#!/bin/bash

################################################################################
# Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)          #
# 100 Queen Street, Melbourne 3000, ABN 11 005 357 522.                        #
# Unauthorized copying of this file, via any medium is strictly prohibited     #
# Proprietary and confidential                                                 #
# Written by Chris T'en <chris.ten@anz.com> March 2016                         #
################################################################################

clear

showMenu() {
cat<<EOF
====================================================
ANZ Nostro Reconciliations Blockchain POC01
----------------------------------------------------
What would you like to do?

(I)nstall Chaincode
(A)dd Record
(D)elete Record
(Q)uery
(M)enu
(E)xit

RECORD MANAGEMENT
----------------------------------------------------
(F) Add Funding Message
(+) Add Payment Instruction
(C) Add Payment Confirmation
(-) Delete Range of Records

REPORTS
----------------------------------------------------
(p) PrettyPrint Range of Records
(G) Get Balance History
(B) Show All Statement Accounts for a Specified Bank
(P) Show Transaction Summary

----------------------------------------------------
EOF
}
showMenu

# Variables
timesDeployed=0 # Ensures deploy function can only be run once.
                # Need a better way to check existence of chaincode deployment.
loop=1
startDir=$PWD
runDir=$GOPATH/src/github.com/openblockchain/obc-peer
COLOUR='\033[0;34m'
NC='\033[0m' # No Color


# cd to obc-peer directory
cd $runDir

while [ $loop -eq 1 ]
do
#  select option in "Initiate" "Query" "Exit"; do
  read option
  case $option in
    [Ii] )  echo "You selected INSTALL CHAINCODE";
            echo;
            read -r -p "Chaincode only needs to be deployed once per VM. Are you sure? [y/N] " response
            case $response in
              [yY][eE][sS]|[yY])
#                echo "Enter chaincode name to be created: ";
#                read ccName;
#                cmd="./obc-peer chaincode deploy -n $ccName -c '{\"Function\":\"init\", \"Args\": [\"a\",\"100\"]}'"
                cmd="./obc-peer chaincode deploy -n anz -c '{\"Function\":\"init\", \"Args\": [\"a\",\"100\"]}'"
                echo $cmd;
                eval $cmd;;
                *)
                    ;;
            esac
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Aa] )  echo "You selected ADD RECORD";
            echo
            echo "Enter key: ";
            read key;
            echo "Enter value: ";
            read value;
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"init\", \"Args\": [\"$key\",\"$value\"]}'"
            echo $cmd;
            printf "Transaction ID: ";
            eval $cmd;
            echo "Record [$key] added."
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [+] )  echo "You selected ADD PAYMENT INSTRUCTION";
            echo
            echo "Enter Amount: ";
            read amount;
            echo "Enter Payer's Name: ";
            read payer;
            echo "Enter Payer's Bank: "
            read payerBank
            echo "Enter Beneficiary's Bank: ";
            read beneficiaryBank;
            echo "Enter Beneficiary's Name: ";
            read beneficiary;
            echo "Enter Fee Type: ";
            read feeType;
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"$amount\",\"$payer\",\"$payerBank\",\"$beneficiaryBank\",\"$beneficiary\",\"$feeType\"]}'"
            echo $cmd;
            eval $cmd;
            echo "Record [Pay \$$amount from $payer at $payerBank to $beneficiary at $beneficiaryBank] added."
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Cc] )  echo "You selected ADD PAYMENT CONFIRMATION";
            echo
            echo "Choose confirmation option: ";
            echo "   1) via filename.";
            echo "   2) via request ID.";
            read option
            case $option in
              1)  echo "Enter MT103 filename: ";
                  read filename;;
              2)  echo "Enter Payment Request ID: ";
                  read key;;
            esac
#            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"mt103\"]}'"
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"$filename\",\"$key\"]}'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;



    [Dd] )  echo "You selected DELETE RECORD";
            read -p "Please enter record key: " key
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"delete\", \"Args\": [\"$key\"]}'"
            echo $cmd;
            printf "Transaction ID: ";
            eval $cmd;
            echo "Record [$key] deleted."
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;


    [Qq] )  echo "You selected QUERY";
#            read -p "Please enter chaincode name: " ccName
            read -p "Please enter search key: " key
            cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"query\", \"Args\": [\"$key\"]}'"
            echo "********** SEARCHING **********";
            echo $cmd;
            printf "Value: ";
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    # [p] )  echo "You selected PRINT RANGE OF RECORDS";
    #         read -p "Please enter start key: " startKey
    #         read -p "Please enter end key: " endKey
    #         cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"queryRange\", \"Args\": [\"$startKey\",\"$endKey\"]}'"
    #         echo "********** SEARCHING **********";
    #         echo $cmd;
    #         eval $cmd;
    #         echo
    #         printf "${COLOUR}What else would you like to do? ${NC}";;

    [p] )  echo "You selected PRETTYPRINT RANGE OF RECORDS";
            read -p "Please enter start key: " startKey
            read -p "Please enter end key: " endKey
            cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"queryRangePretty\", \"Args\": [\"$startKey\",\"$endKey\"]}'"
            echo "********** SEARCHING **********";
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;


    [-] )   echo "You selected DELETE RANGE OF RECORDS";
            read -p "Please enter start key: " startKey
            read -p "Please enter end key: " endKey
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"deleteRange\", \"Args\": [\"$startKey\",\"$endKey\"]}'"
            echo
            read -r -p "Are you sure? This process cannot be reversed. Delete range? [y/N] " response
            case $response in
              [yY][eE][sS]|[yY])
                echo $cmd;
                eval $cmd;
                echo;;
                *) echo;
                    ;;
            esac

            printf "${COLOUR}What else would you like to do? ${NC}";;

    # [Gg] )  echo "You selected GET BALANCE HISTORY";
    #         echo
    #         echo "Enter Account Owner: ";
    #         read owner;
    #         echo "Enter Account Holder: ";
    #         read holder;
    #         cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"getCurrentBalance\", \"Args\": [\"$owner\",\"$holder\"]}'"
    #         echo $cmd;
    #         eval $cmd;
    #         echo
    #         printf "${COLOUR}What else would you like to do? ${NC}";;
    #
    [Gg] )  echo "You selected GET BALANCE HISTORY";
            echo
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"getBalanceHistory\", \"Args\": [\"$owner\",\"$holder\"]}'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;


    [Ff] )  echo "You selected ADD FUNDING MESSAGE";
            echo
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            echo "Enter Funding Amount: ";
            read fundingAmount;
            cmd="./obc-peer chaincode invoke -l golang -n anz -c '{\"Function\":\"addLedgerEntryFunding\", \"Args\": [\"$owner\",\"$holder\",\"$fundingAmount\"]}'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Bb] )  echo "You selected SHOW ALL STATEMENT ACCOUNTS FOR A SPECIFIED BANK";
            read -p "Please enter bank name: " bankName
            cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"queryRangeBank\", \"Args\": [\"$bankName\"]}'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [P] )   echo "You selected SHOW TRANSACTION SUMMARY";
            read -p "Please enter bank name: " bankName
            cmd="./obc-peer chaincode query -l golang -n anz -c '{\"Function\":\"queryBankTransactions\", \"Args\": [\"$bankName\"]}'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Mm] )  echo "Menu selected"; clear; showMenu;;
    [Ee] )  echo "Exiting program. Thanks!"; echo;exit;;
    * )     echo "Error: Please select from the menu. Enter \"M\" to return to the menu.";;
  esac
done

cd $startDir
