#!/bin/bash

################################################################################
# Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)          #
# 100 Queen Street, Melbourne 3000, ABN 11 005 357 522.                        #
# Unauthorized copying of this file, via any medium is strictly prohibited     #
# Proprietary and confidential                                                 #
# Written by Chris T'en <chris.ten@anz.com> March 2016                         #
################################################################################

# Variables
loop=1
startDir=$PWD
runDir=$GOPATH/src/github.com/openblockchain/obc-peer
COLOUR='\033[0;34m'
NC='\033[0m' # No Color
TZ=""
sandbox=false
fabric=false
membership=false
peerID="0"

# Help message for script.
usage() {
cat << EOF
usage: $0 options

Description:
Command Line Interface for ANZ / Wells Nostro Chaincode.

Usage:
$(basename "$0") [-hfsm] [-t n]

Flags:

    -h  show this help text.
    -f  set commands to suit hyperledger fabric. Default is IBM\'s Openblockchain (obc).
    -t  set the timezone to be recorded against transactions, where n is the timezone in Country/City format. UTC if blank.
    -s  configures commands to operate in sandbox environment.
    -m  set the chaincode to work in membership enabled mode.

EOF
}

# cmd parameterisation - default is obc framework in multinode setup.
executable="./obc-peer"
# chaincodePath="-p Corres-Bank-POC/POC01/chaincode/multinode-nostro"
chaincodePath="-p github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro"
peerAddress="OPENCHAIN_PEER_ADDRESS=172.17.0."
peerID="2"
peerPort=":30303"
language=""
username=""
usernameFlag=""
attributeFlag=""
attributes=""

# flags
while getopts "fhsmt:" flag; do
  case "${flag}" in
    t)  TZ=$OPTARG;;
    f)  fabric=true;
        executable="peer";
        peerAddress="CORE_PEER_ADDRESS=172.17.0.";
        ;;
    s)  sandbox=true;
        chaincodePath="-n anz";
        peerAddress="";
        peerID=""
        peerPort=""
        language="-l golang";
        ;;
    m)  membership=true;
        usernameFlag="-u";
        attributeFlag="-a";
        attributes='["enrolment"]';
        ;;

    h)  echo -e $usage
        exit;;
  esac
done

clear
showMenu() {
cat<<EOF
===========================================
ANZ Nostro Reconciliations Blockchain POC01
-------------------------------------------
What would you like to do?

(I)nstall Chaincode
(T) Synchronize Admin Crypto Material
(S)pecify Chaincode Name
(M)enu
(E)xit

RECORD MANAGEMENT
-------------------------------------------
(F) Add Funding Message
(+) Add Payment Instruction
(C)onfirm Payment
(R)eject Payment
(D)elete Record
(-) Delete Range of Records

USER MANAGEMENT
-------------------------------------------
(L) Enroll User
(U) Specify User

REPORTS
-------------------------------------------
(Q)uery Record
(K) Get All Keys
(A) Get All Entries
(G) Get Balance History
(H) Get Unconfirmed Balance History
(B) Show All Statement Accounts for a Specified Bank
(P) Show Transaction Summary
    (1) Get Received Unconfirmed Payment Requests
    (2) Get Sent Unconfirmed Payment Requests
(X) Show All Rejected Payments for a Specified Bank

-------------------------------------------
EOF
}
showMenu

if [ "$fabric" = true ] ; then
  runDir=$GOPATH/src/github.com/hyperledger/fabric/peer
fi

# cd to obc-peer / fabric directory
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
                cmd="$peerAddress$peerID$peerPort $executable chaincode deploy $usernameFlag $username $chaincodePath -c '{\"Function\":\"init\", \"Args\": [\"a\",\"100\"]}'"
                echo $cmd
                NAME="$(eval $cmd)"
                chaincodePath="-n $NAME"
                echo
                echo "CHAINCODE NAME: $NAME"
                echo;;
                *)
                    ;;
            esac
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Ss] ) echo "You selected SPECIFY CHAINCODE NAME"
            echo
            read -p "Please provide Chaincode name: " NAME
            echo
            echo "Chaincode name set to $NAME"
            chaincodePath="-n $NAME"
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Tt] ) echo "You selected SYNCHRONIZE ADMIN CRYPTO MATERIAL"
            echo;
            if [ "$membership" = true ] && [ "$sandbox" = false ]; then
                cmd="rm -rf /opt/gopath/src/github.com/crypto && mkdir /opt/gopath/src/github.com/crypto"
                echo $cmd
                eval $cmd
                # read -p "Please provide VP0 Docker Container ID: " VP0

                # Get container IDs for each vp.
                container0=$(eval docker ps | grep '0.0.0.0:5000->5000' | awk '{print$1}')
                container1=$(eval docker ps | grep '0.0.0.0:5001->5000' | awk '{print$1}')
                container2=$(eval docker ps | grep '0.0.0.0:5002->5000' | awk '{print$1}')
                container3=$(eval docker ps | grep '0.0.0.0:5003->5000' | awk '{print$1}')

                # Determine how many peers in the network.
                read -p "How many peers are running (including VP0)? " peers
                COUNTER=0
                TMP="container$COUNTER"
                cmd="docker cp ${!TMP}:/var/hyperledger/production/crypto/client/vp0Admin /opt/gopath/src/github.com/crypto"
                echo $cmd && eval $cmd
                cmd="docker cp ${!TMP}:/var/hyperledger/production/crypto/client/vp1Admin /opt/gopath/src/github.com/crypto"
                echo $cmd && eval $cmd
                cmd="docker cp ${!TMP}:/var/hyperledger/production/crypto/client/vp2Admin /opt/gopath/src/github.com/crypto"
                echo $cmd && eval $cmd
                cmd="docker cp ${!TMP}:/var/hyperledger/production/crypto/client/vp3Admin /opt/gopath/src/github.com/crypto"
                echo $cmd && eval $cmd
                echo "Crypto material copied to temp location."
                echo;
                let COUNTER=COUNTER+1
                while [ $COUNTER -lt $peers ]; do
                  TMP="container$COUNTER"
                  cmd="docker cp /opt/gopath/src/github.com/crypto/vp0Admin/ ${!TMP}:/var/hyperledger/production/crypto/client"
                  echo $cmd && eval $cmd
                  cmd="docker cp /opt/gopath/src/github.com/crypto/vp0Admin/ ${!TMP}:/var/hyperledger/production/crypto/client"
                  echo $cmd && eval $cmd
                  cmd="docker cp /opt/gopath/src/github.com/crypto/vp1Admin/ ${!TMP}:/var/hyperledger/production/crypto/client"
                  echo $cmd && eval $cmd
                  cmd="docker cp /opt/gopath/src/github.com/crypto/vp2Admin/ ${!TMP}:/var/hyperledger/production/crypto/client"
                  echo $cmd && eval $cmd
                  cmd="docker cp /opt/gopath/src/github.com/crypto/vp3Admin/ ${!TMP}:/var/hyperledger/production/crypto/client"
                  echo $cmd && eval $cmd
                  let COUNTER=COUNTER+1
                done
                cmd="rm -rf /opt/gopath/src/github.com/crypto && mkdir /opt/gopath/src/github.com/crypto"
                eval $cmd
                echo "Crypto material copied to all validating peers."
                echo;
            else echo "CRYPTO MATERIAL SYNCHRONIZING IS NOT REQUIRED FOR A SECURITY DISABLED/SANDBOX CHAINCODE.";
            fi
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Uu] ) echo "You selected SPECIFY USERNAME"
            echo
            read -p "Please provide the username: " username
            echo
            if [ "$sandbox" = false ] ; then
                read -p "Please provide vp number user is logged in to: " peerID
                echo "Current username set to $username, logged in to vp$peerID"
                let "peerID += 2"
                echo
            else
                echo "Current username set to $username"
            fi

            printf "${COLOUR}What else would you like to do? ${NC}";;

    [+] )   echo "You selected ADD PAYMENT INSTRUCTION";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            echo "Enter Amount: ";
            read amount;
            echo "Enter Payer's Bank: "
            read payerBank
            echo "Enter Beneficiary's Bank: ";
            read beneficiaryBank;
            echo "Enter TRN: ";
            read trn;
            echo "Enter Currency: ";
            echo
            echo "  1) AUD"
            echo "  2) USD"
            read response
            case $response in
              1) currency="AUD";;
              2) currency="USD";;
            esac
            echo "Select Fee Type: ";
            echo
            echo "  1) OUR - Beneficiary receives all of payment."
            echo "           Receiving FI will settle its fees with the Sending FI in bulk at EOM."
            echo "           For DL, accrual of Receiving FI's fee can be viewed, but will"
            echo "           not impact Sending FI's statement account until settled at EOM."
            echo
            echo "  2) SHA - Beneficiary receives amount, less Receiving FI's fees."
            echo "           Receiving FI may also charge a fee to the Sending FI, which is settled"
            echo "           in bulk at EOM."
            echo
            echo "  3) BEN - Same as SHA."
            echo
            read response
            case $response in
              1) feeType="OUR";;
              2) feeType="SHA";;
              3) feeType="BEN";;
            esac
            echo
            echo "Select Payment Type: ";
            echo
            echo "  1) BOOK         - Beneficiary has an account at the Receiving FI."
            echo
            echo "  2) INTERMEDIARY - Beneficiary has an account at an FI other than the Receiving FI."
            echo
            read response
            case $response in
              1) bookInt="BOOK";;
              2) bookInt="INTERMEDIARY";;
            esac
            date="$(TZ=$TZ date '+%FT%T.%N%:z')";
            localDate="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
            DAY=$(date -d "$D" '+%d')
            MONTH=$(date -d "$D" '+%m')
            YEAR=$(date -d "$D" '+%Y')
            valueDate=${YEAR:2}$MONTH$DAY
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"$amount\",\"$payerBank\",\"$beneficiaryBank\",\"$date\",\"$feeType\",\"$valueDate\",\"$trn\",\"$currency\",\"$localDate\",\"123\",\"103\",\"$bookInt\"]}' $attributeFlag '$attributes'"
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Cc] )  echo "You selected ADD PAYMENT CONFIRMATION";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            echo "Enter Payment Request ID: ";
            read key;
            echo "Enter Name Bank Performing the Confirmation: ";
            read bankName;
            date="$(TZ=$TZ date '+%FT%T.%N%:z')";
            localTime="$(TZ='Australia/Melbourne' date '+%FT%T.%N%:z')";
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"addPaymentConfirmation\", \"Args\": [\"$key\",\"$date\",\"$bankName\",\"$localTime\",\"1\",\"103\",\"123\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Rr] )  echo "You selected REJECT PAYMENT";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Enter Payment Request ID: " key
            read -p "Enter Rationale for Rejection: " rationale
            date="$(TZ=$TZ date '+%FT%T.%N%:z')"
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"rejectPaymentInstruction\", \"Args\": [\"$key\",\"$rationale\",\"$date\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Dd] )  echo "You selected DELETE RECORD";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter record key: " key
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"remove\", \"Args\": [\"$key\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            printf "Transaction ID: ";
            eval $cmd;
            echo "Record [$key] deleted."
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [-] )   echo "You selected DELETE RANGE OF RECORDS";
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter start key: " startKey
            read -p "Please enter end key: " endKey
            # if [ "$membership" = true ] ; then
            #   read -p "Enter username: " username
            # fi
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"deleteRange\", \"Args\": [\"$startKey\",\"$endKey\"]}' $attributeFlag '$attributes'"
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

    [Ll] )  echo "You selected ENROLL USER";
            echo
            if [ "$membership" = true ] ; then
              if [ "$sandbox" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
                read -p "Please enter username: " username
                cmd="$peerAddress$peerID$peerPort $executable network login $username"
                echo $cmd;
                eval $cmd;
                echo
            else echo "USER ENROLLMENT IS NOT REQUIRED FOR A SECURITY DISABLED CHAINCODE.";
            fi
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Qq] )  echo "You selected QUERY";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter search key: " key
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $chaincodePath -c '{\"Function\":\"get\", \"Args\": [\"$key\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            printf "Value: ";
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;



    [Gg] )  echo "You selected GET BALANCE HISTORY";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getBalanceHistory\", \"Args\": [\"$owner\",\"$holder\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;



    [Hh] )  echo "You selected GET UNCONFIRMED (INDICATIVE) BALANCE HISTORY";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"$owner\",\"$holder\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;



    [Ff] )  echo "You selected ADD FUNDING MESSAGE";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            echo "Enter Funding Amount: ";
            read fundingAmount;
            date="$(TZ=$TZ date '+%FT%T.%N%:z')";
            cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"addLedgerEntryFunding\", \"Args\": [\"$owner\",\"$holder\",\"$fundingAmount\",\"$date\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Bb] )  echo "You selected SHOW ALL STATEMENT ACCOUNTS FOR A SPECIFIED BANK";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter bank name: " bankName
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getStatementAccounts\", \"Args\": [\"$bankName\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Pp] )  echo "You selected SHOW TRANSACTION SUMMARY";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter bank name: " bankName
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getTransactionSummary\", \"Args\": [\"$bankName\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [1] )   echo "You selected GET RECEIVED UNCONFIRMED REQUESTS";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter bank name: " bankName
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getReceivedUnconfirmedPayments\", \"Args\": [\"$bankName\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [2] )   echo "You selected GET RECEIVED UNCONFIRMED REQUESTS";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter bank name: " bankName
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getSentUnconfirmedPayments\", \"Args\": [\"$bankName\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Kk] )  echo "You selected GET ALL KEYS";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $chaincodePath -c '{\"Function\":\"keys\",\"Args\":[\"\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Aa] )  echo "You selected GET ALL ENTRIES";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $chaincodePath -c '{\"Function\":\"getAll\",\"Args\":[\"true\",\"true\",\"true\",\"true\",\"true\",\"true\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Xx] )  echo "You selected SHOW ALL REJECTED PAYMENTS FOR A SPECIFIED BANK";
            echo
            if [ "$sandbox" = false ] ; then
              if [ "$membership" = false ] ; then
                read -p "Please enter ID number of validating peer: " peerID
                echo
                let "peerID += 2"
              fi
            fi
            read -p "Please enter bank name: " bankName
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getRejectedPaymentInstructions\", \"Args\": [\"$bankName\"]}' $attributeFlag '$attributes'"
            echo $cmd;
            echo
            eval $cmd;
            echo
            printf "${COLOUR}What else would you like to do? ${NC}";;

    # load )  echo "You selected UPLOAD MT103s FROM FILE";
            # echo
            # # read -p "Please enter name of file: " INPUT
            # # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/anzToWells103AUD.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/anzToWells103AUDrequest.txt
            #
            # while IFS=, read MSG MSG_TYPE SENDER_BIC RECEIVER_BIC F20 UPDTIME F32AVD F32ACCY F32AAMT F33BCCY F33BAMT F57IND F571 F572 F573 F71A ValueDate LocalTime UTCTime UTCTimeUnformatted SenderFI ReceiverFI
            # do
            #     # echo -e "$line \n"
            #     echo $F32AAMT
            #     echo $SenderFI
            #     echo $ReceiverFI
            #     echo $UTCTime
            #     echo $F71A
            #     echo $ValueDate
            #     echo $F20
            #     echo $F32ACCY
            #     echo $LocalTime
            #     echo $MSG
            #     echo $MSG_TYPE
            #     echo
            #     # echo "Combined: $a1$a2$a3$a4"
            #     cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"$F32AAMT\",\"$SenderFI\",\"$ReceiverFI\",\"$UTCTime\",\"$F71A\",\"$ValueDate\", \"$F20\",\"$F32ACCY\",\"$LocalTime\",\"$MSG\",\"$MSG_TYPE\",\"$BookInt\"]}' $attributeFlag '$attributes'"
            #     echo "UPLOADING #$MSG: $cmd";
            #     eval $cmd;
            #     # sleep 0.5s
            #     echo
            # done <$INPUT
            # printf "${COLOUR}What else would you like to do? ${NC}";;


    # confirm )  echo "You selected CONFIRM MT103s FROM FILE";
    #         echo
    #         # read -p "Please enter name of file: " INPUT
    #         # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/anzToWells103AUDconfirm.txt
    #         INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/anzToWells103AUDconfirm.txt
    #
    #         while IFS=, read MSG MSG_TYPE SENDER_BIC RECEIVER_BIC F20 UPDTIME F32AVD F32ACCY F32AAMT F33BCCY F33BAMT F57IND F571 F572 F573 F71A ValueDate LocalTime UTCTime UTCTimeUnformatted SenderFI ReceiverFI
    #         do
    #             # echo -e "$line \n"
    #             echo $F32AAMT
    #             echo $SenderFI
    #             echo $ReceiverFI
    #             echo $UTCTime
    #             echo $F71A
    #             echo $ValueDate
    #             echo $F20
    #             echo $F32ACCY
    #             echo $LocalTime
    #             echo $MSG
    #             echo $MSG_TYPE
    #             echo
    #             # echo "Combined: $a1$a2$a3$a4"
    #             cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"matchUnconfirmedTransactions\", \"Args\": [\"$F32AAMT\",\"$SenderFI\",\"$ReceiverFI\",\"$UTCTime\",\"$F71A\",\"$ValueDate\", \"$F20\",\"$F32ACCY\",\"$LocalTime\",\"$MSG\",\"$MSG_TYPE\"]}' $attributeFlag '$attributes'"
    #             echo "CONFIRMING #$MSG: $cmd";
    #             eval $cmd;
    #             # sleep 0.5s
    #             echo
    #         done <$INPUT
    #         printf "${COLOUR}What else would you like to do? ${NC}";;

    upload )  echo "You selected UPLOAD COMBINED MT103s REQUESTS AND CONFIRMATIONS FROM FILE";
            echo
            # read -p "Please enter name of file: " INPUT
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv2.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv3.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv3shortTest.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv3RULES_TEST.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv4RULES_TEST.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv4.txt
            INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv4-dup-timestamp-test-big.txt
            # INPUT=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/data/upload-Masterv4-dup-2.txt

            # cat $INPUT |
            while IFS=, read -r MSG MSG_TYPE SENDER_BIC RECEIVER_BIC F20 UPDTIME F32AVD F32ACCY F32AAMT F33BCCY F33BAMT F57IND F571 F572 F573 F574 F575 F71A ValueDate LocalTime UTCTime UTCTimeUnformatted SenderFI ReceiverFI TYPE BookInt
            do
                echo "MSG#:      $MSG"
                echo "MSG_TYPE:  $MSG_TYPE"
                echo "TYPE:      $TYPE"
                echo "Sender:    $SenderFI"
                echo "Receiver:  $ReceiverFI"
                echo "F20:       $F20"
                echo "Currency:  $F32ACCY"
                echo "Amount:    $F32AAMT"
                echo "ValueDate: $ValueDate"
                echo "UTC Time:  $UTCTime"
                echo "Fee Type:  $F71A"
                echo "LocalTime: $LocalTime"
                echo "BookInt:   $BookInt"
                echo
                # echo "Combined: $a1$a2$a3$a4"
                if [ $TYPE = "REQUEST" ]; then
                  cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"addPaymentInstruction\", \"Args\": [\"$F32AAMT\",\"$SenderFI\",\"$ReceiverFI\",\"$UTCTime\",\"$F71A\",\"$ValueDate\",\"$F20\",\"$F32ACCY\",\"$LocalTime\",\"$MSG\",\"$MSG_TYPE\",\"$BookInt\",\"$F33BCCY\",\"$F33BAMT\"]}' $attributeFlag '$attributes'"
                  echo "UPLOADING #$MSG: $cmd";
                  eval $cmd;
                else
                  cmd="$peerAddress$peerID$peerPort $executable chaincode invoke $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"matchUnconfirmedTransactions\", \"Args\": [\"$F32AAMT\",\"$SenderFI\",\"$ReceiverFI\",\"$UTCTime\",\"$F71A\",\"$ValueDate\", \"$F20\",\"$F32ACCY\",\"$LocalTime\",\"$MSG\",\"$MSG_TYPE\",\"$BookInt\",\"$F33BCCY\",\"$F33BAMT\"]}' $attributeFlag '$attributes'"
                  echo "CONFIRMING #$MSG: $cmd";
                  eval $cmd;
                fi
                # sleep 0.5s
                echo
            # done
            done <$INPUT
            printf "${COLOUR}What else would you like to do? ${NC}";;


    print )
            # Get Unconfirmed Statement Balance & Pipe into balance.txt
            echo "GET UNCONFIRMED STATEMENT BALANCE"
            echo
            echo "Enter Account Owner: ";
            read owner;
            echo "Enter Account Holder: ";
            read holder;
            echo
            # cmd="peer chaincode query -l golang -n anz -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"WFB\",\"ANZ\"]}' $attributeFlag '$attributes'"
            cmd="$peerAddress$peerID$peerPort $executable chaincode query $usernameFlag $username $language $chaincodePath -c '{\"Function\":\"getUnconfirmedBalanceHistory\", \"Args\": [\"$owner\",\"$holder\"]}' $attributeFlag '$attributes'"
            echo $cmd && eval $cmd > $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/balance.txt
            printf "${COLOUR}What else would you like to do? ${NC}";;

    [Mm] )  echo "Menu selected"; clear; showMenu;;
    [Ee] )  echo "Exiting program. Thanks!"; echo;exit;;
    * )     echo "Error: Please select from the menu. Enter \"M\" to return to the menu.";;
  esac
done

cd $startDir
