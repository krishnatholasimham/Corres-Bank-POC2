# REST ENDPOINTS

- **Chaincode**
	   
      - SETUP
      	  - [Enrol User](#enrol-user)
      	  - [Deploy Chaincode](#deploy-chaincode)
      	  
      - RECORD MANAGEMENT
	      - [Add Funding Message](#add-funding-message)
	      - [Add Payment Instruction](#add-payment-instruction)
	      - [Confirm Payment](#confirm-payment-instruction)
	      - [Reject Payment](#reject-payment)
	      - [Delete Record](#delete-record)
	      - [Delete Range of Records](#delete-range-of-records)
	      
      - REPORTING
          - [Query Record](#query-record)
          - [Get All Keys](#get-all-keys)
          - [Get Balance History](#get-balance-history)
          - [Get Unconfirmed Balance History](#get-unconfirmed-balance-history)
          - [Show All Statement Accounts for a Specified Bank](#show-all-statement-accounts-for-a-specified-bank)
          - [Show Transaction Summary](#show-transaction-summary)
            - [Show Received Unconfirmed Payment Requests Awaiting Action](#show-received-unconfirmed-payment-requests-awaiting-action)
            - [Show Sent Unconfirmed Payment Requests](#show-sent-unconfirmed-payment-requests)
          - [Show All Rejected Payments for a Specified Bank](#show-all-rejected-payments-for-a-specified-bank)
          - [Get All Records](#get-all-records)
          - [matchUnconfirmedTransactions]

##SETUP
###Enrol User
In a security enabled blockchain environment, an enrolled username must be provided to deploy, invoke or query the chaincode.

```json
POST host:port/registrar

{
  "enrollId": "username",
  "enrollSecret": "password"
}
```

###Deploy Chaincode
Prior to invoking or querying the chaincode must be deployed to all the peers within the network.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
    "type": 1,
    "chaincodeID":{
        "path":"github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro"
    },
    "ctorMsg": {
        "function":"init",
        "args":["a", "100"]
    },
    "secureContext": "username",
  },
  "id": 1
}
```
 

## RECORD MANAGEMENT
### Add Funding Message
Adds a funding message to the DL to top up the balance of a specified statement account.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"addLedgerEntryFunding",
        "args":["owner", "holder", "fundingAmount", "date"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```


### Add Payment Instruction
Adds a payment request to the DL, which represents an MT103 in the current SWIFT process.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"addPaymentInstruction",
        "args":["amount", "payerBank", "beneficiaryBank", "UTCtime", "feeType", "valueDate in YYMMDD", "trn", "currency", "LocalTime", "messageNumber", "messageType"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```


### Confirm Payment Instruction
Takes a request key (hash of an MT103 or agreed subset) and creates a confirmation message in the DL, as well as updating the original payment request entry with the confirmation key.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"addPaymentConfirmation",
        "args":["key", "UTD time of confirmation", "name of confirming bank", "local time of confirmation", "messageNumber", "messageType"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Match Unconfirmed Transactions
Reads the settled transactions and invokes the addPaymentConfirmation with the hashstrings of keys of pending payments.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"matchUnconfirmedTransactions",
        "args":["amount", "SendingFI", "ReceivingFI", "UTC time", "fee type", "value date", "Sender's Reference", "currency", "local Time", "messageNumber", "messageType"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
**Note:** The arguments here need to be refined to match the data file that will be provided by the business. Currently only a subset of these arguments are used to generate the relevant lookup key.

### Reject Payment
Takes a request key (hash of an MT103 or agreed subset) and updates its status to rejected, along with a rationale.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"rejectPaymentInstruction",
        "args":["key", "rationale", "date"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
### Delete Record
Deletes a record from the DL as specified by the lookup key provided.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"remove",
        "args":["key"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
### Delete Range of Records
Deletes all records from the DL whose lookup keys fall between the start and end key range (inclusive).

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "invoke",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"deleteRange",
        "args":["startKey", "endKey"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
## REPORTING
### Query Record
Retrieves the value of a record stored at the specified lookup key.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"get",
        "args":["key"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Get All Keys
Returns a list of all lookup keys stored in the DL.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"keys",
        "args":[""]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Get Balance History
Returns a filtered list of records from the DL that relate to a specific statement account. Records are sorted in ascending chronological order to ensure the current statement account balance is displayed as the last item in the list.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getBalanceHistory",
        "args":["owner", "holder"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```


### Get Unconfirmed Balance History
Returns a filtered list of records from the DL that relate to a specific statement account. Records are sorted in ascending chronological order to ensure the current **indicative** statement account balance is displayed as the last item in the list.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getUnconfirmedBalanceHistory",
        "args":["owner", "holder"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```


### Show All Statement Accounts for a Specified Bank
Returns the record relating to the latest statement account balance for every account owned by the specified bank, held at other banks.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getStatementAccounts",
        "args":["owner"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
### Show Transaction Summary
Returns a list of records for a specified bank, sorted into:

- **OUTWARD PAYMENT REQUESTS:** Payment instructions created by the specified bank.
  - CONFIRMED
  - UNCONFIRMED
- **INWARD PAYMENT REQUESTS:** Payment instructions received by the specified bank from other banks.
  - CONFIRMED
  - UNCONFIRMED
  - REJECTED
- **FUNDING MESSAGES:** Funding messages relating to accounts owned by the specified bank, held at other banks.
  - ACCOUNT AT BANK01
  - ACCOUNT AT BANK02
  - ACCOUNT AT BANK03
  - ACCOUNT AT BANK0n
- **FEE MESSAGES:** Fee entries relating to statement accounts owned by the specified bank.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getTransactionSummary",
        "args":["owner"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Show Received Unconfirmed Payment Requests Awaiting Action
Returns a list of all **inward** payment requests received by the specified bank that are awaiting confirmation (i.e. unconfirmed and not rejected).
A sub-list of the [Show Transaction Summary](#show-transaction-summary) report.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getReceivedUnconfirmedPayments",
        "args":["holder"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Show Sent Unconfirmed Payment Requests
Returns a list of all **outward** payment requests sent by the specified bank that are awaiting confirmation (i.e. unconfirmed and not rejected).
A sub-list of the [Show Transaction Summary](#show-transaction-summary) report.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getSentUnconfirmedPayments",
        "args":["holder"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
### Show All Rejected Payments for a Specified Bank
Returns a list of all outward payment requests from the specified bank that have been rejected by the receiving FI.
A sub-list of the [Show Transaction Summary](#show-transaction-summary) report.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getRejectedPaymentInstructions",
        "args":["owner"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```

### Get All Records
Returns all entries stored in the DL, but can be limited to specific entry types by specifying the relevant bools in the following order: request, confirmation, funding, fee.

**Note:** Omit "secureContext" and "attributes" if not using security.

```json
POST host:port/chaincode

{
  "jsonrpc": "2.0",
  "method": "query",
  "params": {
    "type": 1,
    "chaincodeID":{
        "name":"chaincode name"
    },
    "ctorMsg": {
        "function":"getAll",
        "args":["requestBool", "confirmationBool", "fundingBool", "feeBool"]
    },
    "secureContext": "username",
    "attributes": ["enrolment"]
  },
  "id": 1
}
```
