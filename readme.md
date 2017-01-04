# ANZ Nostro Reconciliations Proof of Concept

**TABLE OF CONTENTS**
<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [1. Business Scenario](#1-business-scenario)
- [2. Structures](#2-structures)
	- [2.1. Ledger Entry Structs](#21-ledger-entry-structs)
		- [2.1.1. A Payment Instruction](#211-a-payment-instruction)
		- [2.1.2. A Payment Confirmation](#212-a-payment-confirmation)
		- [2.1.3. A Funding Message](#213-a-funding-message)
- [3. Example Chaincode Flow](#3-example-chaincode-flow)
- [4. Chaincode Setup Instructions](docs/setup.md) (On Fabric.  ~~Former Open Blockchain [setup](docs/setupInstructions.md)~~)
- [5. Membership Setup Instructions](docs/membershipService.md)
- [6. Fee Calculation Logic](docs/feeCalculation.md) 
- [7. Chaincode Endpoints](docs/chaincodeAPI.md)

<!-- /TOC -->
## 1. Business Scenario
For the purposes of this explanation, the following scenario will be used.
- Alice is a customer of Bank A, who wishes to transfer $100 to Bob, who is a customer of Bank B.
- Alice is known as the "payer".
- Bob is known as the "beneficiary".
- Bank A and Bank B reside in different geographies and have a bilateral correspondent banking relationship in place.
- Bank A **owns** a **statement account** that is **held** by Bank B. Note re terminology:
  - From Bank A's perspective, this **statement account** can also be referred to as a **nostro account** ('our' account).
  - From Bank B's perspective, this **statement account** can also be referred to as a **vostro account** ('your' account).
- In order for Bank A to track the activity on its **statement account** at Bank B, it replicates the transactions in a **nostro mirror account**. This account holds no legal value. It simply mirrors the activity in Bank A's statement account.
- The following is a high-level sequence of steps involved in satisfying Alice's request:
  1. Alice instructs Bank A to transfer $100 to Bob at Bank B via a channel (e.g. internet banking).
  2. Bank A deducts $100 from Alice's account and credits Bank A's **nostro mirror account** for its **statement account** at Bank B.
  3. Bank A sends a payment instruction to Bank B via the Swift network in the form of an MT103.
  4. When Bank B receives the MT103, it debits $100 from Bank A's **statement account** and credits it to Bob's account.
  5. At Bank B's end-of-day, it sends an MT940 back to Bank A. The MT940 is an EOD statement, which contains every transaction that was executed against Bank A's **statement account**. Each transactions is stored in a "transaction block".
  6. Bank A receives the MT940 and matches each transaction against a transaction recorded in the **nostro mirror account**.
- Note: for simplicity, the involvement of suspense accounts have been omitted. These are typically used to hold a payment until the relevant validation steps are complete (e.g. AML / Sanctions checking). This will not be modelled as part of this POC.

## 2. Structures
To model the business process, this chaincode uses two types of structs:

1. **Ledger Entry structs:** define the information stored in the world state in order to track and maintain the payment messages and account balances involved in correspondent banking.

2. **Reporting structs:** used to collate relevant information from the various ledger entries in order to generate certain reports.

### 2.1. Ledger Entry Structs
Four constructs are used to model an interbank payment:

1. A Payment Instruction
2. A Payment Confirmation
3. A Funding Message
4. A Fee Message

#### 2.1.1. A Payment Instruction
**When an MT103 is created and sent by Bank A to Bank B, a ```LedgerEntryRequest``` is also created in the world state. The ```PaymentInstruction``` within this entry contains the contents of the MT103.**

A Payment Instruction represents a request from Bank A to transfer a certain amount from its **statement account** to one of Bank B's customers.
In the chaincode, a Payment Instruction is represented by the ```LedgerEntryRequest``` struct.
```Go
type LedgerEntryRequest struct {
	ConfirmationID string
	PaymentInstruction
	Timestamp string
}
```
It contains the following elements:
- ```ConfirmationID```: The key in the world state against which a corresponding Payment Confirmation is stored. If empty, this indicates an unconfirmed payment instruction.
- ```PaymentInstruction```: A struct which contains an element for each field in an MT103 Swift message.
- ```Timestamp```: The time that this entry was created.

The MT103 is represented by the ```PaymentInstruction``` struct.
```Go
type PaymentInstruction struct {
	InwardSequenceNumber              string // Unique 8 digit sequence number assigned by Swift.
	TransactionReferenceNumber        string // [20] 16 alphanumeric
	BankOperationCode                 string // [23B] 4 digit CRED code
	ValueDateCurrencyInterbankSettled        // [32A] YYMMDD, 3-digit, 12345,98
	CurrencyOriginalAmount            string // [33B] USD12345,98 15-digits
	Payer                                    // [50A] Ordering Customer, Address, Account Number 123456789012345
	OrderingInstitution               string // [52A] Payer's Bank
	AccountWithInstitution            string // [57A] Beneficiary's Bank / Wells Fargo PNBPUS3NNYC
	Beneficiary                              // [59] Recipient name address and account number
	RemittanceInfo                    int    // [70] Payment Reference number
	FeeType                           string // [71A] Details of Charges (BEN-Beneficiary will pay / OUR-Ordering customer will pay / SHA-Shared between ordering customer and beneficiary)
	Comments                          string // [72]	Sender to Receiver Information
	//	RegulatoryReporting                         string // [77B]	Regulatory Reporting
}
```
Its contents are described in the comments contained in-line. Note that the ```RegulatoryReporting``` field is not used in this POC.

#### 2.1.2. A Payment Confirmation
**When an MT103 is received by Bank B and executed, a ```LedgerEntryConfirmation``` is created in the world state. The ```PaymentConfirmation``` within this entry contains the contents of a transaction block, which will be recorded in the EOD statement (i.e. MT940) sent by Bank B to Bank.**

A Payment Confirmation represents a confirmation by Bank B that a requested transfer has been executed (i.e. money has been transferred from Bank A's **statement account** into a beneficiary's account.
In the chaincode, a Payment Confirmation is represented by the ```LedgerEntryConfirmation``` struct.
```Go
type LedgerEntryConfirmation struct {
	RequestID 		string
	PaymentConfirmation
	StatementAccountBalance int
	OrderingInstitution     string
	AccountWithInstitution  string
	Timestamp               string
}
```
It contains the following elements:
- ```RequestID```: The key in the world state against which the request being executed is stored (i.e. the ```LedgerEntryRequest```). Every confirmation must have a corresponding request. As such, this element should never be blank.
- ```PaymentConfirmation```: A struct which contains an element for each field in a transaction block of an MT940 Swift message.
- ```StatementAccountBalance```: The balance of Bank A's **statement account** as a result of this request being confirmed.
- ```OrderingInstitution```: The name of the bank that generated the payment request (i.e. Bank A).
- ```AccountWithInstitution```: The name of the bank that holds the relevant **statement account** (i.e. Bank B).
- ```Timestamp```: The time that this entry was created.

The transaction block in an MT940 is represented by the ```PaymentConfirmation``` struct.
```Go
type PaymentConfirmation struct {
	StatementLine      string // Optional. Details of each transaction.
	Date               string // Mandatory. Date as YYMMDD.
	EntryDate          string // Optional. Entry date as MMDD.
	FundsCode          string // Mandatory. 1-2 digit code. C = credit, D = debit, RC = Reversal credit, RD = Reversal debit.
	Amount             int    // Mandatory. Amount with comma as decimal separator.
	SwiftCode          string // Mandatory. F and 3 signs of Swift Code(?)
	RefAccountOwner    string // Mandatory. Client's information (first line).
	RefOriginatingBank string // Optional. The originating bank's own reference. The last 16 digits of the Transaction Reference Number.
	TransactionDesc    string // Optional. Description according to txn code (e.g. Card transaction, calc of default interest).
	InfoToAccountOwner string // Optional. Additional information passed to account owner.
}
```
Its contents are described in the comments contained in-line.

#### 2.1.3. A Funding Message
**When Bank A's additional funds (i.e. a "top up") are sent to Bank A's statement account, a ```LedgerEntryFunding``` message is created in the world state.**

A Funding Message represents a transfer by Bank A of additional funds to its **statement account** at Bank B.

In the chaincode, a Funding Message is represented by the ```LedgerEntryFunding``` struct.
```Go
type LedgerEntryFunding struct {
	OrderingInstitution     string
	AccountWithInstitution  string
	FundingAmount           int
	StatementAccountBalance int
	Timestamp               string
}
```
It contains the following elements:
- ```OrderingInstitution```: As above.
- ```AccountWithInstitution```: As above.
- ```FundingAmount```: The amount by which Bank A's **statement account** balance should be increased.
- ```StatementAccountBalance```: The balance of Bank A's **statement account** as a result of this funding message being completed.
- ```Timestamp```: The time that this entry was created.

Note: This message will only ever **increase** a **statement account** balance.

## 3. Example Chaincode Flow

<table>
	<tr>
		<td>#</td>
		<td>Current Process</td>
		<td>Shared Ledger</td>
	</tr>
	<tr>
		<td>1.</td>
		<td>Bank A funds statement account held by Bank B with $10000.</td>
		<td>```LedgerEntryFunding``` created.<br>
			<ul>
				<li>```OrderingInstitution```=Bank A.</li>
				<li>```AccountWithInstitution```=Bank B.</li>
				<li>```FundingAmount```=10000.</li>
				<li>```StatementAccountBalance```=10000</li>
			</ul>
		</td>
	</tr>
	<tr>
		<td>2.</td>
		<TD>Bank A sends an **MT103** to Bank B, requesting $2000 be transferred from its **statement account** to a beneficiary.
		<BR><BR>Bank A also records this transaction in a **nostro mirror account**, mirroring Bank A's **statement account** at Bank B.</TD>
		<TD>`LedgerEntryRequest` created.<BR>
			<UL>
				<LI>```ConfirmationID```=""</LI>
				<LI>```PaymentInstruction```=Use fields of MT103</LI></TD>
	</tr>
	<TR>
		<TD>3. </TD>
		<TD>Bank B receives an **MT103**, and executes the payment instruction by disbursing $2000 from Bank A's **statement account** to the beneficiary's account.</TD>
		<TD>`LedgerEntryConfirmation` created.
			<UL>
				<LI>`RequestID`=Key of `LedgerEntryRequest`<BR>**Note:** this allows Bank B to identify the payment request that is being confirmed by this `LedgerEntryConfirmation`.</LI>
				<LI>`StatementAccountBalance`=10000-2000=8000</LI>
				<LI>`PaymentConfirmation`=Information that will eventually be used to construct a transaction block in an **MT940**.</LI>
				<LI>`OrderingInstitution`=Bank A</LI>
				<LI>`AccountWithInstitution`=Bank B</LI>
			</UL>
			`LedgerEntryRequest` updated.
			<UL>
				<LI>`ConfirmationID`=Key of `LedgerEntryConfirmation`<BR>**Note:** this allows Bank A to identify that a payment request has been confirmed, and to locate the corresponding `LedgerEntryConfirmation`.</LI>
			</UL>
		</TD>
	</TR>
	<TR>
		<TD>4.</TD>
		<TD>At Bank B's EOD, Bank B sends Bank A an **MT940**, which lists all the transactions Bank B has processed against Bank A's **statement account**.</TD>
		<TD>**REPORT:** `QueryBankTransactions` function can be called to generate a report of confirmed and unconfirmed `LedgerEntryRequest`s for a given bank name.<BR>
		<BR>
		I.e.
		<UL>
			<LI>Payment Instructions Created By Bank X</LI>
				<UL>
					<LI>Confirmed</LI>
					<LI>Unconfirmed</LI>
				</UL>
			<LI>Payment Instructions Received by Bank X</LI>
			<UL>
				<LI>Confirmed</LI>
				<LI>Unconfirmed</LI>
			</UL>
		</UL>
		</TD>
	</TR>
	<TR>
		<TD>5.</TD>
		<TD>Bank A receives the **MT940** from Bank B, and reconciles each transaction block recorded in the Swift message against the transactions recorded against Bank A's **nostro mirror account**.</TD>
		<TD>
		**REPORT:** The status of every payment instruction can be checked at any point in time by running the `QueryBankTransactions` function.
		</TD>
	</TR>

</table>
