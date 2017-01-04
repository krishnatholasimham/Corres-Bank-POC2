#!/bin/bash

# List of valid cURL commands that call APIs from outside the VM.

# GET /transactions/{UUID}
# Use the Transaction API to retrieve an individual transaction matching the UUID from the blockchain. The returned transaction message is defined in section 3.1.2.1.
# Transaction Retrieval Request:
# GET host:port/transactions/f5978e82-6d8c-47d1-adec-f18b794f570e

# curl -I --> prints headers
# Executing curl with no parameters on a URL (resource) will execute a GET.



host="http://127.0.0.1"
port="3000"
UUID="c3ed7a52-c854-4e99-8679-ff2906a787ec"
block="10"
COLOUR='\033[0;34m'
NC='\033[0m' # No Color

clear
printf "${COLOUR}======================================================================================================================${NC}\n"
printf "${COLOUR}IBM Blockchain API${NC}\n"
printf "${COLOUR}======================================================================================================================${NC}\n"
echo

# GET /chain/blocks/{Block}
# Returns information about a specific block within the Blockchain
# Transaction Retrieval Request:
# GET host:port/transactions/f5978e82-6d8c-47d1-adec-f18b794f570e
printf "${COLOUR}GET /chain/blocks/{Block}${NC}\n"
echo "  Returns information about a specific block within the Blockchain"
echo "  block={$block}"
echo
eval "curl -s $host:$port/chain/blocks/$block | python -m json.tool"
echo


# GET /transactions/{UUID}
# Use the Transaction API to retrieve an individual transaction matching the UUID from the blockchain. The returned transaction message is defined in section 3.1.2.1.
# Transaction Retrieval Request:
# GET host:port/transactions/f5978e82-6d8c-47d1-adec-f18b794f570e
printf "${COLOUR}GET /transactions/{UUID}${NC}\n"
echo "  Retrieves an individual transaction matching the UUID from the blockchain."
echo "  UUID={$UUID}"
echo
#eval "curl -I $host:$port/transactions/$UUID | sed 's/^/  /'"
eval "curl -s $host:$port/transactions/$UUID | python -m json.tool"
echo


printf "TRANSACTIONS\n"
printf "* Chaincode ID: A hash of the chaincode source, path to the source code, constructor function, and parameters.\n"
printf "* Payload Hash: As the payload can be large, only the payload hash is included directly in the transaction message.\n"
printf "* UUID:         A unique ID for the transaction.\n"
printf "${COLOUR}======================================================================================================================${NC}\n"
#printf "*************************************************** TRANSACTIONS *****************************************************\n"
printf "${COLOUR}| Chaincode ID |                        Payload Hash                        |                  UUID                  |${NC}\n"
printf "${COLOUR}======================================================================================================================${NC}\n"

txnccID=`curl -s http://127.0.0.1:3000/transactions/{c3ed7a52-c854-4e99-8679-ff2906a787ec} | jsawk -n 'out(this.chaincodeID)'`

txnpayload=`curl -s http://127.0.0.1:3000/transactions/{c3ed7a52-c854-4e99-8679-ff2906a787ec} | jsawk -n 'out(this.payload)'`

txnuuid=`curl -s http://127.0.0.1:3000/transactions/{c3ed7a52-c854-4e99-8679-ff2906a787ec} | jsawk -n 'out(this.uuid)'`

printf "|%11s   |%58s  |%38s  |\n\n" $txnccID $txnpayload $txnuuid


# COLOURS
# Black        0;30     Dark Gray     1;30
# COLOUR          0;31     Light COLOUR     1;31
# Green        0;32     Light Green   1;32
# Brown/Orange 0;33     Yellow        1;33
# Blue         0;34     Light Blue    1;34
# Purple       0;35     Light Purple  1;35
# Cyan         0;36     Light Cyan    1;36
# Light Gray   0;37     White         1;37
