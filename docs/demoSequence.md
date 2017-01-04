1) Deploy chaincode using `(I)nstall Chaincode`.
In the CLI:
```
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
i
You selected INSTALL CHAINCODE

Chaincode only needs to be deployed once per VM. Are you sure? [y/N] y
./obc-peer chaincode deploy -n anz -c '{"Function":"init", "Args": ["a","100"]}'
04:07:13.052 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
04:07:13.056 [logging] LoggingInit -> DEBU 002 Setting default logging level to DEBUG for command 'chaincode'
04:07:13.092 [main] chaincodeDeploy -> INFO 003 Deploy result: type:GOLANG chaincodeID:<name:"anz" > ctorMsg:<function:"init" args:"a" args:"100" >
anz
What else would you like to do?
```


You should see the following output on TERMINAL 1 (PEER):
```
03:03:54.348 [devops] Deploy -> DEBU 01e Creating deployment transaction (anz)
03:03:54.349 [devops] Deploy -> DEBU 01f Sending deploy transaction (anz) to validator
03:03:54.354 [peer] sendTransactionsToThisPeer -> DEBU 020 Marshalling transaction CHAINCODE_NEW to send to self
03:03:54.354 [peer] sendTransactionsToThisPeer -> DEBU 021 Sending message CHAIN_TRANSACTION with timestamp seconds:1465268634 nanos:354288658  to self
03:03:54.358 [peer] handleChat -> DEBU 022 Current context deadline = 0001-01-01 00:00:00 +0000 UTC, ok = false
03:03:54.364 [consensus/noops] newNoops -> DEBU 023 Creating a NOOPS object
03:03:54.374 [consensus/noops] newNoops -> INFO 024 NOOPS consensus type = *noops.Noops
03:03:54.378 [consensus/noops] newNoops -> INFO 025 NOOPS block size = 500
03:03:54.378 [consensus/noops] newNoops -> INFO 026 NOOPS block timeout = 1s
03:03:54.382 [consensus/handler] SendMessage -> DEBU 027 Sending to stream a message of type: RESPONSE
03:03:54.385 [peer] SendMessage -> DEBU 028 Sending message to stream of type: RESPONSE
03:03:54.388 [peer] func1 -> DEBU 029 Received RESPONSE message as expected, will wait for EOF
03:03:54.389 [consensus/noops] RecvMsg -> DEBU 02a Handling OpenchainMessage of type: CHAIN_TRANSACTION
03:03:54.391 [consensus/noops] broadcastConsensusMsg -> DEBU 02b Broadcasting CONSENSUS
03:03:54.393 [consensus/noops] RecvMsg -> DEBU 02c Sending to channel tx uuid: %!(EXTRA string=anz)
03:03:54.393 [peer] handleChat -> DEBU 02d Received EOF, ending Chat
03:03:54.395 [peer] func1 -> DEBU 02e Received EOF
03:03:55.394 [consensus/noops] handleChannels -> DEBU 02f Process block due to time
03:03:55.394 [consensus/noops] processTransactions -> DEBU 030 Starting TX batch with timestamp: seconds:1465268635 nanos:394371098
03:03:55.398 [consensus/noops] processTransactions -> DEBU 031 Executing batch of 1 transactions with timestamp seconds:1465268635 nanos:394371098
03:03:55.398 [chaincode] DeployChaincode -> DEBU 032 user runs chaincode, not deploying chaincode
03:03:55.400 [state] TxBegin -> DEBU 033 txBegin() for txUuid [anz]
03:03:55.401 [chaincode] LaunchChaincode -> DEBU 034 Container not in READY state(established)...send init/ready
03:03:55.403 [chaincode] initOrReady -> DEBU 035 sending INIT
03:03:55.408 [chaincode] setChaincodeSecurityContext -> DEBU 036 setting chaincode security context...
03:03:55.409 [chaincode] setChaincodeSecurityContext -> DEBU 037 setting chaincode security context. Transaction different from nil
03:03:55.409 [chaincode] setChaincodeSecurityContext -> DEBU 038 setting chaincode security context. Metadata []
03:03:55.410 [chaincode] processStream -> DEBU 039 [anz]Move state message INIT
03:03:55.410 [chaincode] HandleMessage -> DEBU 03a [anz]Handling ChaincodeMessage of type: INIT in state established
03:03:55.411 [chaincode] beforeInitState -> DEBU 03b Before state established.. notifying waiter that we are up
03:03:55.411 [chaincode] notifyDuringStartup -> DEBU 03c nothing to notify (dev mode ?)
03:03:55.411 [chaincode] enterInitState -> DEBU 03d [anz]Entered state init
03:03:55.415 [chaincode] processStream -> DEBU 03e [anz]Received message PUT_STATE from shim
03:03:55.415 [chaincode] HandleMessage -> DEBU 03f [anz]Handling ChaincodeMessage of type: PUT_STATE in state init
03:03:55.415 [chaincode] afterPutState -> DEBU 040 Received PUT_STATE in state busyinit, invoking put state to ledger
03:03:55.420 [chaincode] func1 -> DEBU 041 [anz]state is busyinit
03:03:55.424 [state] Set -> DEBU 042 set() chaincodeID=[anz], key=[a], value=[[]byte{0x31, 0x30, 0x30}]
03:03:55.425 [buckettree] newDataKey -> DEBU 043 Enter - newDataKey. chaincodeID=[anz], key=[a]
03:03:55.426 [buckettree] newDataKey -> DEBU 044 Exit - newDataKey=[bucketKey=[level=[5], bucketNumber=[9177]], compositeKey=[anza]]
03:03:55.429 [chaincode] func1 -> DEBU 045 [anz]Completed PUT_STATE. Sending RESPONSE
03:03:55.429 [chaincode] 1 -> DEBU 046 [anz]enterBusyState trigger event RESPONSE
03:03:55.429 [chaincode] processStream -> DEBU 047 [anz]Move state message RESPONSE
03:03:55.429 [chaincode] HandleMessage -> DEBU 048 [anz]Handling ChaincodeMessage of type: RESPONSE in state busyinit
03:03:55.429 [chaincode] enterInitState -> DEBU 049 [anz]Entered state init
03:03:55.429 [chaincode] processStream -> DEBU 04a [anz]sending state message RESPONSE
03:03:55.433 [chaincode] processStream -> DEBU 04b [anz]Received message COMPLETED from shim
03:03:55.433 [chaincode] HandleMessage -> DEBU 04c [anz]Handling ChaincodeMessage of type: COMPLETED in state init
03:03:55.433 [chaincode] beforeCompletedEvent -> DEBU 04d [anz]beforeCompleted - not in ready state will notify when in readystate
03:03:55.434 [chaincode] enterReadyState -> DEBU 04e [anz]Entered state ready
03:03:55.434 [chaincode] notify -> DEBU 04f notifying Uuid:anz
03:03:55.434 [chaincode] LaunchChaincode -> DEBU 050 sending init completed
03:03:55.434 [chaincode] LaunchChaincode -> DEBU 051 LaunchChaincode complete
03:03:55.434 [state] TxFinish -> DEBU 052 txFinish() for txUuid [anz], txSuccessful=[true]
03:03:55.434 [state] TxFinish -> DEBU 053 txFinish() for txUuid [anz] merging state changes
03:03:55.435 [statemgmt] getSortedKeys -> DEBU 054 Sorted keys = []string{"a"}
03:03:55.439 [statemgmt] ComputeCryptoHash -> DEBU 055 computing hash on []byte{0x61, 0x6e, 0x7a, 0x61, 0x31, 0x30, 0x30}
03:03:55.440 [state] GetHash -> DEBU 056 Enter - GetHash()
03:03:55.440 [state] GetHash -> DEBU 057 updating stateImpl with working-set
03:03:55.440 [buckettree] PrepareWorkingSet -> DEBU 058 Enter - PrepareWorkingSet()
03:03:55.441 [buckettree] newDataKey -> DEBU 059 Enter - newDataKey. chaincodeID=[anz], key=[a]
03:03:55.441 [buckettree] newDataKey -> DEBU 05a Exit - newDataKey=[bucketKey=[level=[5], bucketNumber=[9177]], compositeKey=[anza]]
03:03:55.441 [buckettree] add -> DEBU 05b Adding dataNode=[dataKey=[bucketKey=[level=[5], bucketNumber=[9177]], compositeKey=[anza]], value=[100]] against bucketKey=[level=[5], bucketNumber=[9177]]
03:03:55.443 [buckettree] ComputeCryptoHash -> DEBU 05c Enter - ComputeCryptoHash()
03:03:55.443 [buckettree] ComputeCryptoHash -> DEBU 05d Recomputing crypto-hash...

...

Crypto calcs (removed for simplicity)

...

03:03:55.476 [buckettree] computeCryptoHash -> DEBU 08d Appending crypto-hash for child bucket = [level=[1], bucketNumber=[1]]
03:03:55.476 [buckettree] computeCryptoHash -> DEBU 08e Propagating crypto-hash of single child node for bucket = [level=[0], bucketNumber=[1]]
03:03:55.476 [state] GetHash -> DEBU 08f Exit - GetHash()
03:03:55.477 [consensus/noops] processTransactions -> DEBU 090 Committing TX batch with timestamp: seconds:1465268635 nanos:394371098
03:03:55.480 [state] GetHash -> DEBU 091 Enter - GetHash()
03:03:55.480 [buckettree] ComputeCryptoHash -> DEBU 092 Enter - ComputeCryptoHash()
03:03:55.480 [buckettree] ComputeCryptoHash -> DEBU 093 Returing existing crypto-hash as recomputation not required
03:03:55.480 [state] GetHash -> DEBU 094 Exit - GetHash()
03:03:55.482 [indexes] addIndexDataForPersistence -> DEBU 095 Indexing block number [32] by hash = [1280e118c41ec6ccdbde156b22297f188222e2f1a246e215d90acca0fff7552ca7207eda997067916c0a2336c87fd1f9f95aa645ced7a919354892d8090ce03e]
03:03:55.487 [state] AddChangesForPersistence -> DEBU 096 state.addChangesForPersistence()...start
03:03:55.487 [buckettree] getAffectedBuckets -> DEBU 097 Adding changed bucket [level=[5], bucketNumber=[9177]]
03:03:55.488 [buckettree] getAffectedBuckets -> DEBU 098 Changed buckets are = [[level=[5], bucketNumber=[9177]]]
03:03:55.490 [state] AddChangesForPersistence -> DEBU 099 Adding state-delta corresponding to block number[32]
03:03:55.491 [state] AddChangesForPersistence -> DEBU 09a Not deleting previous state-delta. Block number [32] is smaller than historyStateDeltaSize [500]
03:03:55.491 [state] AddChangesForPersistence -> DEBU 09b state.addChangesForPersistence()...finished
03:03:55.492 [ledger] resetForNextTxGroup -> DEBU 09c resetting ledger state for next transaction batch
03:03:55.493 [buckettree] ClearWorkingSet -> DEBU 09d Enter - ClearWorkingSet()
03:03:55.493 [consensus/noops] getBlockData -> DEBU 09e Preparing to broadcast with block number 33
03:03:55.494 [consensus/noops] getBlockData -> DEBU 09f Got the delta state of block number 33
03:03:55.496 [consensus/noops] notifyBlockAdded -> DEBU 0a0 Broadcasting OpenchainMessage_SYNC_BLOCK_ADDED to non-validators
```

And on TERMINAL 2 (CHAINCODE):
```
2016/06/07 03:03:55 [anz]Received message INIT from shim
2016/06/07 03:03:55 [anz]Handling ChaincodeMessage of type: INIT(state:established)
2016/06/07 03:03:55 Entered state init
2016/06/07 03:03:55 [anz]Received INIT, initializing chaincode
2016/06/07 03:03:55 [anz]Inside putstate, isTransaction = true
2016/06/07 03:03:55 [anz]Sending PUT_STATE
2016/06/07 03:03:55 [anz]Received message RESPONSE from shim
2016/06/07 03:03:55 [anz]Handling ChaincodeMessage of type: RESPONSE(state:init)
2016/06/07 03:03:55 [anz]before send
2016/06/07 03:03:55 [anz]after send
2016/06/07 03:03:55 [anz]Received RESPONSE, communicated (state:init)
2016/06/07 03:03:55 [anz]Received RESPONSE. Successfully updated state
2016/06/07 03:03:55 [anz]Init succeeded. Sending COMPLETED
2016/06/07 03:03:55 [anz]Move state message COMPLETED
2016/06/07 03:03:55 [anz]Handling ChaincodeMessage of type: COMPLETED(state:init)
2016/06/07 03:03:55 [anz]send state message COMPLETED
```

2) Clear ledger of all entries.
    - Run (-) Delete Range of Records -> start key = 0, end key = zzzzz


********************************************************************************
Creating and Confirming a Payment
********************************************************************************

3) Create a statement account by creating an initial funding message.
    - Run (F) Add Funding Message ->
        Account Owner:    ANZ
        Account Holder:   Wells Fargo
        Funding Amount:   10000

4) View message via:

    - (p) PrettyPrint Range of Records -> start key = 0, end key = z
          Provides raw print of all ledger records.
    - (G) Get Balance History ->
            Account Owner: ANZ
            Account Holder: Wells Fargo
          Provides a list of balance records (i.e. confirmations and funding),
          ending with the current balance.
    - (B) Query Bank Statement Accounts -> Bank Name = ANZ
          Provides ANZ's accounts held at other banks, and accounts of other
          banks held by ANZ.
    - (P) Query Bank Transactions -> Bank Name = ANZ
          Provides a printout of the following:
            - Payment Instructions Created by ANZ
                - Confirmed
                - Unconfirmed
            - Payment Instructions Received by ANZ
                - Confirmed
                - Unconfirmed
            - Funding Messages for ANZ's Accounts
                - Account at Bank X (i.e. Wells Fargo)
                - Account at Bank Y
                - Account at Bank n

5) Request Wells Fargo transfer $1000 from ANZ's statement account to Bob the
   beneficiary by adding a Payment Instruction.
   ```
      - Run (+) Add Payment Instruction ->
          Amount:             1000
          Payer's Name:       Alice
          Payer's Bank:       ANZ
          Beneficiary's Bank: Wells Fargo
          Beneficiary's Name: Bob
          Enter Fee Type:     tran
```

6) View "Pending" Payment Instruction via:
```
    - (p) PrettyPrint Range, 0, z
      EXPECTED: Should see a pending Payment Instruction signified by a record
      with an empty ConfirmationID.
    - (P) Query Bank Transactions, ANZ
      EXPECTED: An unconfirmed entry under "PAYMENT INSTRUCTIONS CREATED BY ANZ"

    - NOTE: Balance History (G) and Statement Accounts (B) will not change
            until payments are confirmed, or funding messages are sent.
```

7) As Wells Fargo, confirm the "Pending" Payment Instruction.
```
      - Run (C) Add Payment Confirmation.
```

8) View the confirmed Payment Instruction via:
```
    - (p) PrettyPrint Range, 0, z
      EXPECTED: Should see a new Confirmation Record with:
          - RequestID =                 Record Key of corresponding Payment Instruction
          - Amount =                    Amount of Payment Instruction (i.e. 1000)
          - Statement Account Balance:  Reduced by amount of payment
                                        (i.e. 10000 - 1000 = 9000)

                                  AND

          The original Payment Instruction record, updated with the
          Confirmation ID of the confirmation.
```

    - (G) Get Balance History, ANZ, Wells Fargo
```
      EXPECTED: Two entries showing the statement account decreasing from 10000
      to 9000 via a confirmation record.
```

    - (B) ANZ's account at Wells Fargo updated to 9000.

    - (P) Query Bank Transactions, ANZ
      EXPECTED: The unconfirmed entry under "PAYMENT INSTRUCTIONS CREATED BY ANZ"
      moved to confirmed.

********************************************************************************
Viewing Multiple Statement Accounts
********************************************************************************

9) NOTE: Run (B) Query Bank Statement Accounts -> OTHER BANK's ACCOUNTS = null

10) Add a Statement Account owned by BBVA and held at ANZ.
    Run (F) Add Funding Message, BBVA, ANZ, 20000

11) View Statement Accounts relating to ANZ.
    Run (B) Query Bank Statement Accounts, ANZ
    EXPECTED: ANZ's account at Wells Fargo (9000)
              BBVA's account at ANZ (20000)
