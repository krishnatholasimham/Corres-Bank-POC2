/*
Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)
100 Queen Street, Melbourne 3000, ABN 11 005 357 522.
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by Chris T'en <chris.ten@anz.com> March 2016
*/

package main

//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"sort"
//	"strconv"
//	"time"
//
//	"github.com/hyperledger/fabric/Corres-Bank-POC/POC01/utility"
//
//	//"math/big"
//
//	"Corres-Bank-POC/POC01/confidentiality"
//
//	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"github.com/op/go-logging"
//)

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"
	//"math/big"

	"Corres-Bank-POC/POC01/confidentiality"
	"Corres-Bank-POC/POC01/utility" // Not an ideal location to pick from the vendor dir, but works for now ...Bluemix specific change

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/op/go-logging"
)

// This chaincode implements a simple map that is stored in the state.
// The following operations are available.

// Invoke operations
// put - requires two arguments, a key and value
// remove - requires a key

// Query operations
// get - requires one argument, a key, and returns a value
// keys - requires no arguments, returns all keys

var nostroLogger = logging.MustGetLogger("nostro")
var confidentialityEnabled = false

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// LedgerEntry is a structure to store key / value pairs for iteration.
// Note: Not stored on the Shared Ledger.
type LedgerEntry struct {
	Key                    string
	Type                   string
	Timestamp              time.Time
	Balance                float64
	OrderingInstitution    string
	AccountWithInstitution string
}

// LedgerEntryTxn is a structure to store key / value pairs for iteration.
// Extends LedgerEntry, which was designed to retrieve balances. This one needs
// FromBank, ToBank, Amount, DatePosted, DateConfirmed
// Note: Not stored on the Shared Ledger.
type LedgerEntryTxn struct {
	Key                        string
	RefKey                     string
	Type                       string
	SendingFI                  string
	ReceivingFI                string
	Currency                   string
	OrderingInstitution        string
	AccountWithInstitution     string
	PaymentAmount              float64
	FeeType                    string
	SendersCharge              float64
	ReceiversCharge            float64
	BenePays                   float64
	Rebate                     float64
	TimestampCreated           time.Time
	Sequence                   int
	TimestampConfirmed         time.Time
	IsRejected                 bool //Next three elements will only be used if the payment instruction is rejected
	RejectRationale            string
	RejectTimestamp            time.Time
	StatementAccountBalance    float64
	IndicativeBalance          float64
	ValueDate                  string
	UnconfirmedSortTime        time.Time
	LocalTime                  time.Time
	MsgNum                     string
	MsgType                    string
	TransactionReferenceNumber string
	SecondKey                  string
	ThirdKey                   string
	FourthKey                  string
}

// LedgerEntryRequest represnts the fields stored in a key/value entry.
type LedgerEntryRequest struct {
	MsgNum         string
	MsgType        string
	EntryType      string
	ConfirmationID string
	SendingFI      string
	ReceivingFI    string
	PaymentInstruction
	StatementAccountBalance float64
	IndicativeBalance       float64
	IsRejected              bool //Next three elements will only be used if the payment instruction is rejected
	RejectRationale         string
	RejectTimestamp         time.Time
	Timestamp               time.Time
	Sequence                int
	TimestampLocal          time.Time
	Signature               []byte
	SecondKey               string
	ThirdKey                string
	FourthKey               string
}

// PaymentInstruction instructions (MT103) to Bank X to execute a payment.
type PaymentInstruction struct {
	InwardSequenceNumber              string  // Unique 8 digit sequence number assigned by Swift.
	TransactionReferenceNumber        string  // [20] 16 alphanumeric
	BankOperationCode                 string  // [23B] 4 digit CRED code
	ValueDateCurrencyInterbankSettled         // [32A] YYMMDD, 3-digit, 12345,98
	CurrencyOriginalAmount            string  // [33B] USD12345,98 15-digits
	Payer                                     // [50A] Ordering Customer, Address, Account Number 123456789012345
	OrderingInstitution               string  // [52A] Payer's Bank
	AccountWithInstitution            string  // [57A] Beneficiary's Bank / Wells Fargo PNBPUS3NNYC
	Beneficiary                               // [59] Recipient name address and account number
	RemittanceInfo                    int     // [70] Payment Reference number
	FeeType                           string  // [71A] Details of Charges (BEN-Beneficiary will pay / OUR-Ordering customer will pay / SHA-Shared between ordering customer and beneficiary)
	SendersCharge                     float64 // [71F] Sender's Charge
	ReceiversCharge                   float64 // [71G] Receiver's Charges
	BenePays                          float64 // The amount the beneficiary will pay out of the transfer proceeds.
	Rebate                            float64 // In BEN / SHA arrangements, the amount owed by the Receiving FI to the Sending FI from fee collected from the bene.
	Comments                          string  // [72]	Sender to Receiver Information
	//	RegulatoryReporting                         string // [77B]	Regulatory Reporting
}

// ValueDateCurrencyInterbankSettled contains details of payment.
type ValueDateCurrencyInterbankSettled struct {
	ValueDate string
	Currency  string
	Amount    float64
}

// Payer contains the sender's identification details.
type Payer struct {
	PayerName          string
	PayerAddress       string
	PayerAccountNumber int
}

// Beneficiary contains the recipient's identification details.
type Beneficiary struct {
	BenName          string
	BenAddress       string
	BenAccountNumber int
}

// LedgerEntryConfirmation represnts the fields stored in a key/value entry.
type LedgerEntryConfirmation struct {
	MsgNum    string
	MsgType   string
	EntryType string
	RequestID string
	PaymentConfirmation
	SendingFI                  string
	ReceivingFI                string
	Currency                   string
	StatementAccountBalance    float64
	IndicativeBalance          float64
	OrderingInstitution        string
	AccountWithInstitution     string
	Timestamp                  time.Time
	Sequence                   int
	LocalTimestamp             time.Time
	Signature                  []byte
	TransactionReferenceNumber string //TransactionReference on MT950 (Senders Reference). Used to match from file input to show the current view of on the Graph
}

// PaymentConfirmation contains a transaction block [61] from the MT940 that
// confirms execution of a PaymentInstruction.
type PaymentConfirmation struct {
	StatementLine      string  // Optional. Details of each transaction.
	Date               string  // Mandatory. Date as YYMMDD.
	EntryDate          string  // Optional. Entry date as MMDD.
	FundsCode          string  // Mandatory. 1-2 digit code. C = credit, D = debit, RC = Reversal credit, RD = Reversal debit.
	Amount             float64 // Mandatory. Amount with comma as decimal separator.
	SwiftCode          string  // Mandatory. F and 3 signs of Swift Code(?)
	RefAccountOwner    string  // Mandatory. Client's information (first line).
	RefOriginatingBank string  // Optional. The originating bank's own reference. The last 16 digits of the Transaction Reference Number.
	TransactionDesc    string  // Optional. Description according to txn code (e.g. Card transaction, calc of default interest).
	InfoToAccountOwner string  // Optional. Additional information passed to account owner.
	FeeType            string  // [71A] Details of Charges (BEN-Beneficiary will pay / OUR-Ordering customer will pay / SHA-Shared between ordering customer and beneficiary)
	SendersCharge      float64 // [71F] Sender's Charge
	ReceiversCharge    float64 // [71G] Receiver's Charges
	BenePays           float64 // The amount the beneficiary will pay out of the transfer proceeds.
	Rebate             float64 // In BEN / SHA arrangements, the amount owed by the Receiving FI to the Sending FI from fee collected from the bene.
}

// LedgerEntryFunding tops up the specified statement account with the specific
// amount. Current statement account balance is determined by looking at the
// most recent LedgerEntryConfirmation or LedgerEntryFunding
type LedgerEntryFunding struct {
	EntryType               string
	OrderingInstitution     string
	AccountWithInstitution  string
	FundingAmount           float64
	StatementAccountBalance float64
	IndicativeBalance       float64
	Timestamp               time.Time
	Sequence                int
	LocalTimestamp          time.Time
	Signature               []byte
}

// LedgerEntryFee records a fee attributable to a confirmed payment and the resulting statement account balance.
type LedgerEntryFee struct {
	EntryType               string
	RequestID               string
	ConfirmationID          string
	OrderingInstitution     string
	AccountWithInstitution  string
	PaymentAmount           float64
	FeeType                 string  // [71A] Details of Charges (BEN-Beneficiary will pay / OUR-Ordering customer will pay / SHA-Shared between ordering customer and beneficiary)
	SendersCharge           float64 // [71F] Sender's Charge - charges deducted by the Sender and by previous banks in the transaction chain.
	ReceiversCharge         float64 // [71G] Receiver's Charges - Amount due to Receiver FI
	BenePays                float64 // The amount the beneficiary will pay out of the transfer proceeds.
	Rebate                  float64 // In BEN / SHA arrangements, the amount owed by the Receiving FI to the Sending FI from fee collected from the bene.
	StatementAccountBalance float64
	IndicativeBalance       float64
	Timestamp               time.Time
	Sequence                int
	LocalTimestamp          time.Time
	Signature               []byte
}

// type PaymentInstructionKey struct {
// 	ValueDate       string
// 	senderReference string
// 	Amount          float64
// }
//
type PaymentInstructionKey struct {
	F32AAMT         float64
	F32ACCY         string
	F33BAMT         float64
	F33BCCY         string
	SendingFI       string
	ReceivingFI     string
	FeeType         string
	BookInt         string
	ValueDate       string
	SenderReference string
	Currency        string
}

// ByDate implements sort.Interface for []LedgerEntryTxn based on the TimestampConfirmed field.
type ByDate []LedgerEntryTxn

func (slice ByDate) Len() int {
	return len(slice)
}

func (slice ByDate) Less(i, j int) bool {
	return slice[i].TimestampConfirmed.Before(slice[j].TimestampConfirmed)
}

func (slice ByDate) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// ByDateUnconfirmed implements sort.Interface for []LedgerEntryTxn based on the TimestampConfirmed field.
type ByDateUnconfirmed []LedgerEntryTxn

func (slice ByDateUnconfirmed) Len() int {
	return len(slice)
}

func (slice ByDateUnconfirmed) Less(i, j int) bool {
	return slice[i].UnconfirmedSortTime.Before(slice[j].UnconfirmedSortTime)
}

func (slice ByDateUnconfirmed) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

//
// // A Change is a record of source code changes, recording user, language, and delta size.
// type Change struct {
// 	user     string
// 	language string
// 	lines    int
// }

type lessFunc func(p1, p2 *LedgerEntryTxn) bool

// multiSorter implements the Sort interface, sorting the LedgerEntryTxns within.
type multiSorter struct {
	ledgerEntryTxns []LedgerEntryTxn
	less            []lessFunc
}

// Sort sorts the argument slice according to the less functions passed to OrderedBy.
func (ms *multiSorter) Sort(ledgerEntryTxns []LedgerEntryTxn) {
	ms.ledgerEntryTxns = ledgerEntryTxns
	sort.Sort(ms)
}

// OrderedBy returns a Sorter that sorts using the less functions, in order.
// Call its Sort method to sort the data.
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

// Len is part of sort.Interface.
func (ms *multiSorter) Len() int {
	return len(ms.ledgerEntryTxns)
}

// Swap is part of sort.Interface.
func (ms *multiSorter) Swap(i, j int) {
	ms.ledgerEntryTxns[i], ms.ledgerEntryTxns[j] = ms.ledgerEntryTxns[j], ms.ledgerEntryTxns[i]
}

// Less is part of sort.Interface. It is implemented by looping along the
// less functions until it finds a comparison that is either Less or
// !Less. Note that it can call the less functions twice per call.
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.ledgerEntryTxns[i], &ms.ledgerEntryTxns[j]
	// Try all but the last comparison.
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			// p < q, so we have a decision.
			return true
		case less(q, p):
			// p > q, so we have a decision.
			return false
		}
		// p == q; try the next comparison.
	}
	// All comparisons to here said "equal", so just return whatever
	// the final comparison reports.
	return ms.less[k](p, q)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting chaincode: %s", err)
	}
}

// Run callback representing the invocation of a chaincode
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "addPaymentInstruction" {
		return t.addPaymentInstruction(stub, args)
	} else if function == "addPaymentConfirmation" {
		return t.addPaymentConfirmation(stub, args)
	} else if function == "rejectPaymentInstruction" {
		return t.rejectPaymentInstruction(stub, args)
	} else if function == "addLedgerEntryFunding" {
		return t.addLedgerEntryFunding(stub, args)
	} else if function == "deleteRange" {
		return t.deleteRange(stub, args)
	} else if function == "matchUnconfirmedTransactions" {
		return t.matchUnconfirmedTransactions(stub, args)
	}

	switch function {

	case "init":
		// Do nothing

	case "put":
		if len(args) < 2 {
			return nil, errors.New("put operation must include two arguments, a key and value")
		}
		key := args[0]
		value := args[1]

		err := stub.PutState(key, []byte(value))
		if err != nil {
			fmt.Printf("Error putting state %s", err)
			return nil, fmt.Errorf("put operation failed. Error updating state: %s", err)
		}
		return nil, nil

	case "remove":
		if len(args) < 1 {
			return nil, errors.New("remove operation must include one argument, a key")
		}
		key := args[0]

		err := stub.DelState(key)
		if err != nil {
			return nil, fmt.Errorf("remove operation failed. Error updating state: %s", err)
		}
		return nil, nil

	default:
		return nil, errors.New("Unsupported operation")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if function == "getBalanceHistory" {
		return t.getBalanceHistory(stub, args)
	} else if function == "getUnconfirmedBalanceHistory" {
		return t.getUnconfirmedBalanceHistory(stub, args)
	} else if function == "getStatementAccounts" {
		return t.getStatementAccounts(stub, args)
	} else if function == "getTransactionSummary" {
		return t.getTransactionSummary(stub, args)
	} else if function == "getReceivedUnconfirmedPayments" {
		return t.getReceivedUnconfirmedPayments(stub, args)
	} else if function == "getSentUnconfirmedPayments" {
		return t.getSentUnconfirmedPayments(stub, args)
	} else if function == "getRejectedPaymentInstructions" {
		return t.getRejectedPaymentInstructions(stub, args)
	} else if function == "getAll" {
		return t.getAll(stub, args)
	}

	switch function {

	case "get":
		if len(args) < 1 {
			return nil, errors.New("get operation must include one argument, a key")
		}
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return nil, fmt.Errorf("get operation failed. Error accessing state: %s", err)
		}
		if confidentialityEnabled {
			params := []string{key, string(value)}
			passed, _ := confidentiality.IsAuthorisedToQuery(stub, params)
			if passed {
				return value, nil
			} else {
				payload, _ := stub.GetCallerCertificate()
				return payload, nil
			}
		} else {
			return value, nil
		}
		return nil, nil

	case "keys":

		keysIter, err := stub.RangeQueryState("", "")
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		defer keysIter.Close()

		var keys []string
		for keysIter.HasNext() {
			key, _, err := keysIter.Next()
			if err != nil {
				return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
			}
			keys = append(keys, key)
		}

		jsonKeys, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error marshaling JSON: %s", err)
		}

		return jsonKeys, nil

	default:
		return nil, errors.New("Unsupported operation")
	}
}

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	t.initConfidentiality(stub, args)

	var A string     // Entities
	var Aval float64 // Asset holdings
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting key and value")
	}

	// Initialize the chaincode
	A = args[0]
	Aval, err = strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	// Write the state to the ledger
	err = stub.PutState(A, []byte(strconv.FormatFloat(Aval, 'f', 2, 64)))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) deleteRange(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var startKey, endKey string // from / to key range
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting start and end key to delete")
	}

	startKey = args[0]
	endKey = args[1]
	// Get a range of states from the ledger
	iterator, err := stub.RangeQueryState(startKey, endKey)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get iterator for " + startKey + "\"}"
		return nil, errors.New(jsonResp)
	}

	// Retrieve current key and value
	key, _, err := iterator.Next()
	if err != nil {
		return nil, errors.New("[CHRIS] Failed to retrieve key and value via iterator")
	}
	valuesRetrieved := 1
	// loggerCT.Trace.Printf("Item %v, Key = %v, Value = %s\n", valuesRetrieved, key, value)
	_, _ = t.delete(stub, []string{key})

	// Iterate and print subsequent keys and values
	for iterator.HasNext() {
		valuesRetrieved++
		key, _, err := iterator.Next()
		if err != nil {
			return nil, errors.New("[CHRIS] Failed to retrieve key and value via iterator")
		}
		// loggerCT.Trace.Printf("bItem %v, Key = %v, Value = %s\n", valuesRetrieved, key, value)
		_, _ = t.delete(stub, []string{key})
	}
	iterator.Close()
	return nil, nil
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

func (t *SimpleChaincode) initConfidentiality(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("[Confidentiality] Inside Confidentiality.....")
	if confidentiality.IsConfidential() {
		confidentialityEnabled = true
		nostroLogger.Debug("[Confidentiality] Initialiting.....")
		confidentiality.Setup()
		nostroLogger.Debug("[Confidentiality] Setup Complete......")

		nostroLogger.Debug("[Confidentiality] Initialiting Clients.....")
		if err := confidentiality.InitClients(); err != nil {
			panic(fmt.Errorf("Failed initializing Clients [%s]", err))
		}
		nostroLogger.Debug("[Confidentiality] Client Setup Complete......")
	}
	return nil, nil
}

// addPaymentInstruction adds a payment request to the DL, which represents an MT103 in the current SWIFT process.
func (t *SimpleChaincode) addPaymentInstruction(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("=================================================================  ENTERED ADDPAYMENTINSTRUCTION ==================================================")

	if len(args) != 14 {
		return nil, errors.New("Incorrect number of arguments. Expecting 14:\nAmount F32AAMT\nOrdering Institution\nAccount With Institution\nUTC Time\nFee Type\nValue Date\nSender's Reference\nCurrency F32ACCY\nLocal Time\nMSG#\nMSG Type\nBookInt\nF33BCCY\nF33BAMT\n")
	}

	// Parse arguments
	var F32AAMT float64
	var F33BAMT float64
	var entryType string
	var refKey string
	sendingFI := args[1]
	receivingFI := args[2]
	timeStamp, e := time.Parse(time.RFC3339Nano, args[3])
	if e != nil {
		return nil, errors.New("Could not parse the time...")
	}
	feeType := args[4]
	vDate := args[5]
	F20 := args[6]
	F32ACCY := args[7]
	localTime, e := time.Parse(time.RFC3339Nano, args[8])
	msgNum := args[9]
	msgType := args[10]
	bookInt := args[11]
	F33BCCY := args[12]
	// F33BAMT := args[13]
	// nostroLogger.Debug("F57 Received is >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> " + bookInt)

	currencySameOrDiff := "same"
	isn := "12345678"
	//Temp Constants
	beneName := "Bene"
	payerName := "Payer"

	senderFIOwesReceivingFI, benePays, rebate := calculateBankFees(feeType, sendingFI, receivingFI, bookInt, currencySameOrDiff)

	// Determine account to be updated and credit or debit.
	// ANZ->WF & USD = DR from ANZ's USD account held by WF
	// ANZ->WF & AUD =	CR to WF's AUD account held by ANZ *** Direct Credit. No confirmation required.
	// WF->ANZ & AUD =	DR from WF's AUD account held by ANZ
	// WF->ANZ & USD =	CR to ANZ's USD account held by WF *** Direct Credit. No confirmation required.
	orderingInstitution, accountWithInstitution, crdr := determineAccountCrDr(sendingFI, receivingFI, F32ACCY, msgType)

	// Parse F32AAMT arg as a debit or credit.
	switch crdr {
	case "DR":
		F32AAMT, _ = strconv.ParseFloat("-"+args[0], 64)
		F33BAMT, _ = strconv.ParseFloat("-"+args[13], 64)
		entryType = "REQUEST-RECORD"
	case "CR":
		F32AAMT, _ = strconv.ParseFloat(args[0], 64)
		F33BAMT, _ = strconv.ParseFloat(args[13], 64)
		entryType = "DIRECT-CREDIT-RECORD"
		refKey = "NA"
	case "Do Nothing":
		F32AAMT, _ = strconv.ParseFloat(args[0], 64)
		entryType = "ADVICE"
		refKey = "NA"
	}

	// 1. SET KEY MATCHING COMBINATIONS
	// SET FIRST KEY arguments: Senders Ref, F33BAMT, F33BCCY, and Value Date.
	firstKeyArgs := &PaymentInstructionKey{
		SenderReference: F20,
		F33BAMT:         F33BAMT,
		F33BCCY:         F33BCCY,
		ValueDate:       vDate,
	}

	// SET SECOND KEY arguments: F32AAMT, F32ACCY, F33BAMT, F33BCCY, value date, and BookInt.
	secondKeyArgs := &PaymentInstructionKey{
		F32AAMT:   F32AAMT,
		F32ACCY:   F32ACCY,
		F33BAMT:   F33BAMT,
		F33BCCY:   F33BCCY,
		ValueDate: vDate,
		BookInt:   bookInt,
	}

	// SET THIRD KEY arguments: F33BAMT, F33BCCY, value date, and BookInt.
	thirdKeyArgs := &PaymentInstructionKey{
		F33BAMT:   F33BAMT,
		F33BCCY:   F33BCCY,
		ValueDate: vDate,
		BookInt:   bookInt,
	}

	// SET FOURTH KEY arguments: F33BAMT, F33BCCY, value date, and BookInt.
	fourthKeyArgs := &PaymentInstructionKey{
		F32AAMT:   F32AAMT,
		F32ACCY:   F32ACCY,
		ValueDate: vDate,
	}

	// 2. GENERATE KEYS USING MATCH COMBINATIONS
	// First key gen
	out, err := json.MarshalIndent(firstKeyArgs, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal settlementData for addPaymentInstruction.")
	}
	key := utility.GenerateKey(string(out))
	keyString := fmt.Sprintf("%x", key[0:])
	nostroLogger.Debug("First Key ", keyString)

	// NOTE: SECONDARY CHECK will be F20 directly.
	// Second key gen
	out2, err := json.MarshalIndent(secondKeyArgs, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal settlementData for addPaymentInstruction.")
	}
	key2 := utility.GenerateKey(string(out2))
	keyString2 := fmt.Sprintf("%x", key2[0:])

	// Third key gen
	out3, err := json.MarshalIndent(thirdKeyArgs, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal settlementData for addPaymentInstruction.")
	}
	key3 := utility.GenerateKey(string(out3))
	keyString3 := fmt.Sprintf("%x", key3[0:])

	// Fourth key gen
	out4, err := json.MarshalIndent(fourthKeyArgs, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal settlementData for addPaymentInstruction.")
	}
	key4 := utility.GenerateKey(string(out4))
	keyString4 := fmt.Sprintf("%x", key4[0:])

	// 3. POPULATE PAYMENT INSTRUCTION (based on MT103).
	MT103 := &PaymentInstruction{
		InwardSequenceNumber:       isn,
		TransactionReferenceNumber: F20,
		CurrencyOriginalAmount:     F32ACCY + args[0],
		ValueDateCurrencyInterbankSettled: ValueDateCurrencyInterbankSettled{
			ValueDate: vDate,
			Currency:  F32ACCY,
			Amount:    F32AAMT,
		},
		Payer: Payer{
			PayerName:          payerName,
			PayerAddress:       "Address",
			PayerAccountNumber: 123456789012345,
		},
		OrderingInstitution:    orderingInstitution,
		AccountWithInstitution: accountWithInstitution,
		Beneficiary: Beneficiary{
			BenName:          beneName,
			BenAddress:       "Address",
			BenAccountNumber: 123456789012345,
		},
		FeeType:         feeType,
		SendersCharge:   0,
		ReceiversCharge: senderFIOwesReceivingFI,
		BenePays:        benePays,
		Rebate:          rebate,
	}

	// Get ledger entry with the latest confirmed balance.
	ledgerEntryTxn := getLatestBalance(stub, orderingInstitution, accountWithInstitution)

	// Get ledger entry with the latest unconfirmed balance.
	ledgerEntryTxnUnconfirmed := getLatestUnconfirmedBalance(stub, orderingInstitution, accountWithInstitution)
	nostroLogger.Debug("**UNCONFIRMED Balance = %e", ledgerEntryTxnUnconfirmed.IndicativeBalance)

	// Instantiate a new LedgerEntryRequest and populate. Note: ConfirmationID will be
	// populated by the receiving bank once the payment is confirmed.
	// If record is a direct credit, also add to statement account balance.
	var directCredit float64
	var adviceCorrection float64
	directCredit = 0
	adviceCorrection = 0

	sequence := 1

	nostroLogger.Debug("New Entry Timestamp ", timeStamp)
	nostroLogger.Debug("Latest MSG# ", ledgerEntryTxnUnconfirmed.MsgNum)
	nostroLogger.Debug("Latest Timestamp ", ledgerEntryTxnUnconfirmed.TimestampCreated)

	// If timestamp of new entry is the same as timestamp of latest entry, make new entry a higher sequence number.
	if timeStamp.Equal(ledgerEntryTxnUnconfirmed.TimestampCreated) {
		nostroLogger.Debug("Sequence start", sequence)
		sequence = ledgerEntryTxnUnconfirmed.Sequence + 1
		nostroLogger.Debug("Sequence updated", sequence)
	}

	if crdr == "CR" {
		directCredit = F32AAMT
	} else if crdr == "Do Nothing" {
		adviceCorrection = F32AAMT
	}
	value := &LedgerEntryRequest{
		EntryType:               entryType,
		ConfirmationID:          refKey,
		PaymentInstruction:      *MT103,
		StatementAccountBalance: ledgerEntryTxn.StatementAccountBalance + directCredit,
		IndicativeBalance:       ledgerEntryTxnUnconfirmed.IndicativeBalance + F32AAMT - adviceCorrection,
		IsRejected:              false,
		SendingFI:               sendingFI,
		ReceivingFI:             receivingFI,
		Timestamp:               timeStamp,
		Sequence:                sequence,
		TimestampLocal:          localTime,
		MsgNum:                  msgNum,
		MsgType:                 msgType,
		SecondKey:               keyString2,
		ThirdKey:                keyString3,
		FourthKey:               keyString4,
	}

	//Adds the Transaction Signature to the Ledger Entry if confidentiality is enabled
	if confidentialityEnabled {
		nostroLogger.Debug("[Confidentiality] Adding Signature.....")
		params := []string{sendingFI, receivingFI, keyString}
		securityMetaData, err := confidentiality.GetInvokeTransactionSignature(stub, params)
		if err != nil {
			return nil, err
		}
		nostroLogger.Debug("[Confidentiality] Adding Added.....")
		value.Signature = securityMetaData
	}

	// Convert LedgerEntryRequest struct to a json byte array.
	entry, err2 := json.MarshalIndent(value, "", "  ")
	if err2 != nil {
		return nil, errors.New("json.Marshal() FAILED")
	}

	// Write the state to the ledger
	err3 := stub.PutState(keyString, entry)
	if err3 != nil {
		return nil, err3
	}

	return key, nil
}

// addPaymentConfirmation takes a request ID (hash of an MT103) and creates a confirmation message in
// the shared ledger, as well as updating the original request entry with the confirmation ID.
func (t *SimpleChaincode) addPaymentConfirmation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED ADDPAYMENTCONFIRMATION ", args[4])

	if len(args) != 7 {
		return nil, errors.New("Expecting Payment Request ID, timestamp, receiving FI,localTime, messageNumber, messageType, F20")
	}

	// Parse arguments
	requestKeyString := args[0]
	timestamp, e := time.Parse(time.RFC3339Nano, args[1]) // RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	if e != nil {
		return nil, errors.New("Could not parse the time...")
	}
	date := fmt.Sprintf("%0s%02d%02d", strconv.Itoa(timestamp.Year())[2:], timestamp.Month(), timestamp.Day())
	entryDate := fmt.Sprintf("%02d%02d", timestamp.Month(), timestamp.Day())
	receivingFI := args[2]
	localTime, e := time.Parse(time.RFC3339Nano, args[3])
	msgNum := args[4]
	msgType := args[5]
	F20 := args[6]

	// Retrieve value stored at requestKeyString.
	queryBytes, _ := t.query(stub, []string{requestKeyString})

	// Unpack into LedgerEntryRequest struct
	requestEntry := &LedgerEntryRequest{}
	err := json.Unmarshal(queryBytes, requestEntry)
	if err != nil {
		nostroLogger.Debug("ERROR: Record not found for key ", args[0])
		return nil, errors.New("Could not map query value to Ledger Entry Request struct.")
	}

	nostroLogger.Debug("SUCCESS: Retrieved record using key ", args[0])

	// Prevent the user from confirming a rejected payment
	if requestEntry.IsRejected == true {
		return nil, errors.New("Cannot confirm a rejected payment.")
	}

	// Only allow the Receiving FI to confirm a payment request.
	if receivingFI != requestEntry.PaymentInstruction.AccountWithInstitution {
		return nil, errors.New("Unauthorized to confirm payments.")
	}

	// Populate confirmation struct.
	confirmation := &PaymentConfirmation{
		Date:               date,
		EntryDate:          entryDate,
		FundsCode:          "D",
		Amount:             requestEntry.PaymentInstruction.Amount,
		SwiftCode:          "TEST",
		RefAccountOwner:    requestEntry.PaymentInstruction.BenName,
		RefOriginatingBank: requestEntry.PaymentInstruction.TransactionReferenceNumber,
		FeeType:            requestEntry.PaymentInstruction.FeeType,
		SendersCharge:      requestEntry.PaymentInstruction.SendersCharge,
		ReceiversCharge:    requestEntry.PaymentInstruction.ReceiversCharge,
		BenePays:           requestEntry.PaymentInstruction.BenePays,
		Rebate:             requestEntry.PaymentInstruction.Rebate,
	}

	out, err := json.MarshalIndent(confirmation, "", "  ")
	if err != nil {
		return nil, errors.New("Could not marshal confirmation struct.")
	}
	// Generate database key based on MT103.
	key := utility.GenerateKey(string(out))

	// Convert key from []byte to string in hex format for readability.
	confirmationKeyString := fmt.Sprintf("%x", key[0:])

	// Instantiate a new LedgerEntryConfirmation and populate.
	owner := requestEntry.PaymentInstruction.OrderingInstitution
	holder := requestEntry.PaymentInstruction.AccountWithInstitution

	// Get latest StatementAccountBalance
	ledgerEntryTxn := getLatestBalance(stub, owner, holder)
	nostroLogger.Debug("Latest Balance = %e", ledgerEntryTxn.StatementAccountBalance)
	nostroLogger.Debug("Payment Amount = %e", requestEntry.PaymentInstruction.Amount)

	// Get the latest unconfirmed balance.
	ledgerEntryTxnUnconfirmed := getLatestUnconfirmedBalance(stub, owner, holder)
	nostroLogger.Debug("UNCONFIRMED Balance = %e", ledgerEntryTxnUnconfirmed.IndicativeBalance)

	nostroLogger.Debug("New Balance = %e", ledgerEntryTxn.StatementAccountBalance+requestEntry.PaymentInstruction.Amount)

	// If timestamp of new entry is the same as timestamp of latest entry, make new entry a higher sequence number.
	sequence := 1
	nostroLogger.Debug("New Entry Timestamp ", timestamp)
	nostroLogger.Debug("Latest Timestamp ", ledgerEntryTxnUnconfirmed.TimestampCreated)

	if timestamp.Equal(ledgerEntryTxnUnconfirmed.TimestampCreated) {
		nostroLogger.Debug("Sequence start", sequence)
		sequence = ledgerEntryTxnUnconfirmed.Sequence + 1
		nostroLogger.Debug("Sequence updated", sequence)
	}

	value := &LedgerEntryConfirmation{
		EntryType:                  "CONFIRMATION-RECORD",
		RequestID:                  requestKeyString,
		PaymentConfirmation:        *confirmation,
		Timestamp:                  timestamp,
		Sequence:                   sequence,
		OrderingInstitution:        owner,
		AccountWithInstitution:     holder,
		Currency:                   requestEntry.Currency,
		StatementAccountBalance:    ledgerEntryTxn.StatementAccountBalance + requestEntry.PaymentInstruction.Amount,
		IndicativeBalance:          ledgerEntryTxnUnconfirmed.IndicativeBalance,
		LocalTimestamp:             localTime,
		MsgNum:                     msgNum,
		MsgType:                    msgType,
		TransactionReferenceNumber: F20,
	}

	//Adds the Transaction Signature to the Ledger Entry if confidentiality is enabled
	if confidentialityEnabled {
		params := []string{holder, owner, confirmationKeyString}
		securityMetaData, err := confidentiality.GetInvokeTransactionSignature(stub, params)
		if err != nil {
			return nil, err
		}
		value.Signature = securityMetaData
	}

	// Convert LedgerEntryRequest struct to a []byte
	out, err = json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal LedgerEntryConfirmation during addPaymentConfirmation().")
	}

	// Write the new confirmation entry to the ledger
	err = stub.PutState(confirmationKeyString, out)
	if err != nil {
		return nil, errors.New("Could not putstate confirmation key and value.")
	}

	// Update LedgerEntryRequest with confirmationID.
	requestEntry.ConfirmationID = confirmationKeyString

	// Convert LedgerEntryRequest struct to a []byte
	out, err = json.MarshalIndent(requestEntry, "", "  ")
	if err != nil {
		return nil, errors.New("Could not marshal requestEntry struct.")
	}
	err = stub.PutState(requestKeyString, out)
	if err != nil {
		return nil, errors.New("Could not putstate updated request entry. Tried updating with confirmation ID.")
	}

	// Create fee entry
	// if requestEntry.FeeType == "OUR" {
	params := []string{requestKeyString, confirmationKeyString}
	t.addLedgerEntryFee(stub, params)
	// }

	return nil, nil
}

// rejectPaymentInstruction takes a request key (hash of an MT103 or agreed subset) and updates its status to rejected, along with a rationale.
func (t *SimpleChaincode) rejectPaymentInstruction(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Expecting Payment Request ID, Reject Rationale and time.")
	}

	time, e := time.Parse(time.RFC3339Nano, args[2]) // RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	if e != nil {
		return nil, errors.New("Could not parse the time...")
	}
	nostroLogger.Debug("Reject Payments")

	requestKeyString := args[0]

	// Retrieve value stored at requestKeyString.
	queryBytes, _ := t.query(stub, []string{requestKeyString})

	// Unpack into LedgerEntryRequest struct
	requestEntry := &LedgerEntryRequest{}
	err := json.Unmarshal(queryBytes, requestEntry)
	if err != nil {
		return nil, errors.New("Could not map query value to Ledger Entry Request struct.")
	}
	requestEntry.IsRejected = true
	requestEntry.RejectRationale = args[1]
	requestEntry.RejectTimestamp = time

	// Get the latest unconfirmed balance.
	ledgerEntryTxnUnconfirmed := getLatestUnconfirmedBalance(stub, requestEntry.OrderingInstitution, requestEntry.AccountWithInstitution)

	requestEntry.IndicativeBalance = ledgerEntryTxnUnconfirmed.IndicativeBalance + requestEntry.Amount

	// Convert LedgerEntryRequest struct to a []byte
	out, err := json.MarshalIndent(requestEntry, "", "  ")
	if err != nil {
		return nil, errors.New("Could not marshal requestEntry struct.")
	}
	err = stub.PutState(requestKeyString, out)
	if err != nil {
		return nil, errors.New("Could not putstate updated request entry. Tried updating with confirmation ID.")
	}

	return nil, nil
}

func (t *SimpleChaincode) addLedgerEntryFee(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED ADDLEDGERENTRYFEE")

	requestKeyString := args[0]
	confirmationKeyString := args[1]

	// Query the payment instruction
	queryBytes, _ := t.query(stub, []string{requestKeyString})
	// Unpack into LedgerEntryRequest struct
	queriedLedgerEntryRequest := &LedgerEntryRequest{}
	err := json.Unmarshal(queryBytes, queriedLedgerEntryRequest)
	if err != nil {
		return nil, errors.New("Failed to unmarshal query results into LedgerEntryRequest struct during addLedgerEntryFee().")
	}

	// Query the payment confirmation
	queryBytes, _ = t.query(stub, []string{confirmationKeyString})
	// Unpack into LedgerEntryConfirmation struct
	queriedLedgerEntryConfirmation := &LedgerEntryConfirmation{}
	err = json.Unmarshal(queryBytes, queriedLedgerEntryConfirmation)
	if err != nil {
		return nil, errors.New("Failed to unmarshal query results into LedgerEntryConfirmation struct during addLedgerEntryFee().")
	}

	// Get record containing the latest statement account balance.
	// latestLedgerEntryTxn := getLatestBalance(stub, queriedLedgerEntryConfirmation.OrderingInstitution, queriedLedgerEntryConfirmation.AccountWithInstitution)

	// Create new ledger entry for Fee record.
	ledgerEntryFee := &LedgerEntryFee{
		EntryType:               "FEE-RECORD",
		RequestID:               queriedLedgerEntryConfirmation.RequestID,
		ConfirmationID:          confirmationKeyString,
		OrderingInstitution:     queriedLedgerEntryConfirmation.OrderingInstitution,
		AccountWithInstitution:  queriedLedgerEntryConfirmation.AccountWithInstitution,
		FeeType:                 queriedLedgerEntryConfirmation.FeeType,
		PaymentAmount:           queriedLedgerEntryConfirmation.Amount,
		SendersCharge:           queriedLedgerEntryConfirmation.SendersCharge,
		ReceiversCharge:         queriedLedgerEntryConfirmation.ReceiversCharge,
		BenePays:                queriedLedgerEntryConfirmation.BenePays,
		Rebate:                  queriedLedgerEntryConfirmation.Rebate,
		StatementAccountBalance: 0,
		IndicativeBalance:       0,
		Timestamp:               queriedLedgerEntryConfirmation.Timestamp.Add(1),
	}

	// Convert Fee Message struct to JSON format as []byte.
	out, err := json.MarshalIndent(ledgerEntryFee, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal fee message for keygen.")
	}

	// Generate database key based on Funding Message.
	key := utility.GenerateKey(string(out))

	// Write the state to the ledger
	keyString := fmt.Sprintf("%x", key[0:])

	//Adds the Transaction Signature to the Funding Ledger Entry if confidentiality is enabled
	if confidentialityEnabled {
		params := []string{queriedLedgerEntryConfirmation.AccountWithInstitution, queriedLedgerEntryConfirmation.OrderingInstitution, keyString}
		securityMetaData, err := confidentiality.GetInvokeTransactionSignature(stub, params)
		if err != nil {
			return nil, err
		}
		ledgerEntryFee.Signature = securityMetaData
		out, err = json.MarshalIndent(ledgerEntryFee, "", "  ")
		if err != nil {
			return nil, errors.New("Failed to marshal funding message for keygen.")
		}
	}

	err = stub.PutState(keyString, out)
	if err != nil {
		return nil, errors.New("Error in putstate for fee entry.")
	}
	return out, nil
}

// addLedgerEntryFunding adds a funding message to the ledger to top up balance
// of a specified statement account.
func (t *SimpleChaincode) addLedgerEntryFunding(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED ADDLEDGERENTRYFUNDING")

	if len(args) != 4 {
		return nil, errors.New("Expecting OrderingInstitution, AccountWithInstitution, FundingAmount, and Time")
	}

	fundingAmount, _ := strconv.ParseFloat(args[2], 64)
	time, e := time.Parse(time.RFC3339Nano, args[3]) // RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	if e != nil {
		return nil, errors.New("Could not parse the time...")
	}

	latestLedgerEntryTxn := getLatestBalance(stub, args[0], args[1])
	latestLedgerEntryTxnUnconfirmed := getLatestUnconfirmedBalance(stub, args[0], args[1])

	// If timestamp of new entry is the same as timestamp of latest entry, make new entry a higher sequence number.
	sequence := 1
	if time.Equal(latestLedgerEntryTxnUnconfirmed.TimestampCreated) {
		sequence = latestLedgerEntryTxnUnconfirmed.Sequence + 1
	}

	ledgerEntryFunding := &LedgerEntryFunding{
		EntryType:               "FUNDING-RECORD",
		OrderingInstitution:     args[0],
		AccountWithInstitution:  args[1],
		FundingAmount:           fundingAmount,
		StatementAccountBalance: latestLedgerEntryTxn.StatementAccountBalance + fundingAmount,
		IndicativeBalance:       latestLedgerEntryTxnUnconfirmed.IndicativeBalance + fundingAmount,
		Sequence:                sequence,
		Timestamp:               time,
	}

	// Convert Funding Message struct to JSON format as []byte.
	out, err := json.MarshalIndent(ledgerEntryFunding, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal funding message for keygen.")
	}

	// Generate database key based on Funding Message.
	key := utility.GenerateKey(string(out))

	// Write the state to the ledger
	out, err = json.MarshalIndent(ledgerEntryFunding, "", "  ")
	keyString := fmt.Sprintf("%x", key[0:])

	//Adds the Transaction Signature to the Funding Ledger Entry if confidentiality is enabled
	if confidentialityEnabled {
		params := []string{args[0], args[1], keyString}
		securityMetaData, err := confidentiality.GetInvokeTransactionSignature(stub, params)
		if err != nil {
			return nil, err
		}
		ledgerEntryFunding.Signature = securityMetaData
		out, err = json.MarshalIndent(ledgerEntryFunding, "", "  ")
		if err != nil {
			return nil, errors.New("Failed to marshal funding message for keygen.")
		}
	}

	err = stub.PutState(keyString, out)
	if err != nil {
		return nil, errors.New("PutState failed for Funding Entry.")
	}
	return out, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var A string // Entities
	var err error
	// var prettyPrint []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}
	// prettyPrint = PrettyPrinter(Avalbytes)

	return Avalbytes, nil
}

// getStatementAccounts returns a list organised into accounts owned by Bank X and
// accounts held by Bank X.
func (t *SimpleChaincode) getStatementAccounts(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED GETSTATEMENTACCOUNTS")

	if len(args) != 1 {
		return nil, errors.New("Expecting Bank Name.")
	}

	bankName := args[0]

	// Get a list of transactions with statement balances (i.e. confirmations, funding, fees)
	list, _ := getLedgerEntries(stub, false, true, true, true, true, false)

	var heldList []LedgerEntryTxn
	var ownedList []LedgerEntryTxn

	var listOfHoldingBanks []string
	var listOfHoldingBanksUnique []string
	var ownedListLatestBalances []LedgerEntryTxn
	var latestBalance LedgerEntryTxn

	var listOfOwningBanks []string
	var listOfOwningBanksUnique []string
	var heldListLatestBalances []LedgerEntryTxn

	var result []byte

	width := 86

	// Sift through balance entries list and sort into accounts held by this bank
	// and accounts owned by this bank.
	if len(list) == 0 {
		return nil, errors.New("Shared Ledger is EMPTY")
	}
	for i := 0; i < len(list); i++ {
		if list[i].OrderingInstitution == bankName { // If it's owned by this bank
			ownedList = append(ownedList, list[i])
		} else if list[i].AccountWithInstitution == bankName { // If it's held by this bank
			heldList = append(heldList, list[i])
		}
	}
	// For the accounts we own, create a unique list of banks that hold our accounts.
	// Start by creating a string array of all holding banks, then make it unique.
	if len(ownedList) > 0 {
		for i := 0; i < len(ownedList); i++ {
			listOfHoldingBanks = append(listOfHoldingBanks, ownedList[i].AccountWithInstitution)
		}
		listOfHoldingBanksUnique = append(listOfHoldingBanksUnique, listOfHoldingBanks[0])
	}
	if len(listOfHoldingBanks) > 1 {
		for i := 1; i < len(listOfHoldingBanks); i++ {
			if contains(listOfHoldingBanksUnique, listOfHoldingBanks[i]) {
			} else {
				listOfHoldingBanksUnique = append(listOfHoldingBanksUnique, listOfHoldingBanks[i])
			}
		}
	}
	// For each of these banks, find the entry in ownedList that has the most recent timestamp. Then add to a list.
	if len(listOfHoldingBanksUnique) > 0 {
		for i := 0; i < len(listOfHoldingBanksUnique); i++ {
			latestBalance = getLatestBalance(stub, bankName, listOfHoldingBanksUnique[i])
			ownedListLatestBalances = append(ownedListLatestBalances, latestBalance)
		} // ownedListLatestBalances contains a list of this bank's statement accounts held at other banks.
	}
	// For the accounts held by us, create a unique list of banks that own this accounts.
	// Start by creating a string array of all owning banks, then make it unique.
	if len(heldList) > 0 {
		for i := 0; i < len(heldList); i++ {
			listOfOwningBanks = append(listOfOwningBanks, heldList[i].OrderingInstitution)
		}
		listOfOwningBanksUnique = append(listOfOwningBanksUnique, listOfOwningBanks[0])
	}
	if len(listOfOwningBanks) > 1 {
		for i := 1; i < len(listOfOwningBanks); i++ {
			if contains(listOfOwningBanksUnique, listOfOwningBanks[i]) {
			} else {
				listOfOwningBanksUnique = append(listOfOwningBanksUnique, listOfOwningBanks[i])
			}
		}
	}
	// For each of these banks, find the entry in heldList that has the most recent timestamp. Then add to a list.
	if len(listOfOwningBanksUnique) > 0 {
		for i := 0; i < len(listOfOwningBanksUnique); i++ {
			latestBalance = getLatestBalance(stub, listOfOwningBanksUnique[i], bankName)
			heldListLatestBalances = append(heldListLatestBalances, latestBalance)
		} // heldListLatestBalances contains a list of the statement accounts this bank holds for other banks.
	}
	// Combine the two lists (ownedListLatestBalances & heldListLatestBalances) along with a separator.

	titleO := fmt.Sprintf("%s'S STATEMENT ACCOUNTS HELD AT OTHER BANKS", bankName)
	textLen := len(titleO)
	a := padLeft(titleO, " ", (width/2)+(textLen/2))
	b := padRight(a, " ", width-2)
	border := ""
	border = padLeft(border, "#", width)
	titleOwnedList := fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleOwnedList...)
	tmp, _ := json.MarshalIndent(ownedListLatestBalances, "", "  ")
	result = append(result, tmp...)

	titleH := fmt.Sprintf("OTHER BANK'S STATEMENT ACCOUNTS HELD BY %s", bankName)
	textLen = len(titleH)
	a = padLeft(titleH, " ", (width/2)+(textLen/2))
	b = padRight(a, " ", width-2)
	titleHeldList := fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleHeldList...)
	tmp, _ = json.MarshalIndent(heldListLatestBalances, "", "  ")
	result = append(result, tmp...)

	return result, nil
}

func contains(array []string, item string) bool {
	for _, this := range array {
		if this == item {
			return true
		}
	}
	return false
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0 : length+1]
		}
	}
}

func padLeft(str, pad string, length int) string {
	for {
		str = pad + str
		if len(str) > length {
			return str[0 : length+1]
		}
	}
}

// getBalanceHistory takes an account owner and holder and returns the balance-related entries (i.e. confirmations and funding messages) in chronological order.
func (t *SimpleChaincode) getBalanceHistory(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED GETBALANCEHISTORY")

	if len(args) != 2 {
		return nil, errors.New("Expecting Owner, Holder.")
	}

	owner := args[0]
	holder := args[1]

	var ledgerList []LedgerEntryTxn

	ledgerList, err := getLedgerEntries(stub, false, true, true, false, true, true)
	if err != nil {
		return nil, errors.New("getLedgerEntries() failed during getBalanceHistory")
	}
	if len(ledgerList) == 0 {
		list, err := json.MarshalIndent(ledgerList, "", "  ")
		if err != nil {
			return nil, errors.New("json.MarshalIndent() FAILED")
		}
		return list, nil
		// return nil, errors.New("getLedgerEntries returned 0 entries")
	}
	// Filter down by owner and holder
	var filteredList []LedgerEntryTxn
	for i := 0; i < len(ledgerList); i++ {
		if ledgerList[i].OrderingInstitution == owner && ledgerList[i].AccountWithInstitution == holder {
			filteredList = append(filteredList, ledgerList[i])
		}
	}
	// // Sort list
	// sort.Sort(ByDate(filteredList))

	// Sort by TimestampCreated then Sequence
	// Closures that order the LedgerEntryTxn structure.
	timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.TimestampCreated.Before(c2.TimestampCreated)
	}
	sequence := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.Sequence < c2.Sequence
	}
	OrderedBy(timestampCreated, sequence).Sort(filteredList)

	list, err := json.MarshalIndent(filteredList, "", "  ")
	if err != nil {
		return nil, errors.New("json.MarshalIndent() FAILED")
	}
	return list, nil
}

// getUnconfirmedBalanceHistory takes an account owner and holder and returns the balance-related entries (i.e. confirmations and funding messages) in chronological order.
func (t *SimpleChaincode) getUnconfirmedBalanceHistory(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED GETUNCONFIRMEDBALANCEHISTORY")

	if len(args) != 2 {
		return nil, errors.New("Expecting Owner, Holder.")
	}

	owner := args[0]
	holder := args[1]

	var ledgerList []LedgerEntryTxn

	ledgerList, err := getLedgerEntries(stub, true, true, true, false, true, true)
	if err != nil {
		return nil, errors.New("getLedgerEntries() failed during getUnconfirmedBalanceHistory")
	}
	nostroLogger.Debug("ledgerListLENGTH = %d", len(ledgerList))
	if len(ledgerList) == 0 {
		list, err := json.MarshalIndent(ledgerList, "", "  ")
		if err != nil {
			return nil, errors.New("json.MarshalIndent() FAILED")
		}
		return list, nil
		// return nil, errors.New("getLedgerEntries returned 0 entries")
	}
	// Filter down by owner and holder
	var filteredList []LedgerEntryTxn
	for i := 0; i < len(ledgerList); i++ {
		if ledgerList[i].OrderingInstitution == owner && ledgerList[i].AccountWithInstitution == holder {
			filteredList = append(filteredList, ledgerList[i])
		}
	}
	// // Sort list
	// sort.Sort(ByDateUnconfirmed(filteredList))

	// Sort by TimestampCreated then Sequence
	// Closures that order the LedgerEntryTxn structure.
	timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.TimestampCreated.Before(c2.TimestampCreated)
	}
	sequence := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.Sequence < c2.Sequence
	}
	OrderedBy(timestampCreated, sequence).Sort(filteredList)

	list, err := json.MarshalIndent(filteredList, "", "  ")
	if err != nil {
		return nil, errors.New("json.MarshalIndent() FAILED")
	}
	return list, nil
}

// getAll returns all entries stored in the ledger, but can be limited to
// specific entry types by specifying the relevant bools in the following order:
// request, confirmation, funding, fee.
func (t *SimpleChaincode) getAll(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("ENTERED GETALL")

	if len(args) != 6 {
		return nil, errors.New("Insufficient arguments. Expecting 6 bools representing request, confirmation, funding, fee, direct, advice.")
	}

	request, err1 := strconv.ParseBool(args[0])
	if err1 != nil {
		return nil, errors.New("Error parsing bool.")
	}
	confirmation, err2 := strconv.ParseBool(args[1])
	if err2 != nil {
		return nil, errors.New("Error parsing bool.")
	}
	funding, err3 := strconv.ParseBool(args[2])
	if err3 != nil {
		return nil, errors.New("Error parsing bool.")
	}
	fee, err4 := strconv.ParseBool(args[3])
	if err4 != nil {
		return nil, errors.New("Error parsing bool.")
	}
	direct, err5 := strconv.ParseBool(args[4])
	if err5 != nil {
		return nil, errors.New("Error parsing bool.")
	}
	advice, err6 := strconv.ParseBool(args[5])
	if err6 != nil {
		return nil, errors.New("Error parsing bool.")
	}

	direct = true
	advice = true

	var ledgerList []LedgerEntryTxn

	ledgerList, err := getLedgerEntries(stub, request, confirmation, funding, fee, direct, advice)

	// // Sort list
	// sort.Sort(ByDate(ledgerList))

	// Sort by TimestampCreated then Sequence
	// Closures that order the LedgerEntryTxn structure.
	timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.TimestampCreated.Before(c2.TimestampCreated)
	}
	sequence := func(c1, c2 *LedgerEntryTxn) bool {
		return c1.Sequence < c2.Sequence
	}
	OrderedBy(timestampCreated, sequence).Sort(ledgerList)

	list, err := json.MarshalIndent(ledgerList, "", "  ")
	if err != nil {
		return nil, errors.New("json.MarshalIndent() FAILED")
	}

	return list, nil
}

// getLedgerEntries is a modified version of getBalanceEntries which takes additional bools for confirmation, request, funding and fee to determine which types of entries to return.
func getLedgerEntries(stub *shim.ChaincodeStub, request bool, confirmation bool, funding bool, fee bool, direct bool, advice bool) ([]LedgerEntryTxn, error) {
	// nostroLogger.Debug("ENTERED GETLEDGERENTRIES")
	// nostroLogger.Debug("direct bool = ", direct)

	var ledgerListTxn []LedgerEntryTxn

	// Get a range of states from the ledger
	iterator, err := stub.RangeQueryState("0", "zzzzz")
	if err != nil {
		return nil, errors.New("Failed to get iterator")
	}

	// Retrieve first key and value
	key, value, err := iterator.Next()
	if err != nil {
		return nil, errors.New("Failed to retrieve key and value via iterator")
	}

	if confidentialityEnabled {
		params := []string{key, string(value)}
		passed, _ := confidentiality.IsAuthorisedToQuery(stub, params)
		if passed {
			ledgerListTxn, err = selectiveAppend(stub, ledgerListTxn, key, value, request, confirmation, funding, fee, direct, advice)
			if err != nil {
				return nil, errors.New("selectiveAppend FAILED")
			}
		}
	} else {
		ledgerListTxn, err = selectiveAppend(stub, ledgerListTxn, key, value, request, confirmation, funding, fee, direct, advice)
		if err != nil {
			return nil, errors.New("selectiveAppend FAILED")
		}
	}

	// Iterate and store subsequent keys and values
	for iterator.HasNext() {
		key, value, err := iterator.Next()
		if err != nil {
			return nil, errors.New("Failed to retrieve key and value via iterator")
		}
		if confidentialityEnabled {
			params := []string{key, string(value)}
			passed, _ := confidentiality.IsAuthorisedToQuery(stub, params)
			if passed {
				ledgerListTxn, err = selectiveAppend(stub, ledgerListTxn, key, value, request, confirmation, funding, fee, direct, advice)
				if err != nil {
					return nil, errors.New("selectiveAppend FAILED")
				}
			}
		} else {
			ledgerListTxn, err = selectiveAppend(stub, ledgerListTxn, key, value, request, confirmation, funding, fee, direct, advice)
			if err != nil {
				return nil, errors.New("selectiveAppend FAILED")
			}
		}
	}
	iterator.Close()

	return ledgerListTxn, nil
}

// getLatestBalance same as getCurrentBalance but usable by other functions.
// getCurrentBalance is accessed via CLI.
func getLatestBalance(stub *shim.ChaincodeStub, owner string, holder string) LedgerEntryTxn {
	// nostroLogger.Debug("ENTERED GETLATESTBALANCE")
	list, _ := getLedgerEntries(stub, false, true, true, false, true, false)

	// Filter by owner and holder
	var filteredList []LedgerEntryTxn
	for i := 0; i < len(list); i++ {
		if list[i].OrderingInstitution == owner && list[i].AccountWithInstitution == holder {
			filteredList = append(filteredList, list[i])
		}
	}

	// Find entry with most recent timeStamp
	latestEntry := &LedgerEntryTxn{}
	for i := 0; i < len(filteredList); i++ {
		if filteredList[i].TimestampCreated.Equal(latestEntry.TimestampCreated) {
			if filteredList[i].Sequence > latestEntry.Sequence {
				latestEntry = &filteredList[i]
			}
		} else if filteredList[i].TimestampCreated.After(latestEntry.TimestampCreated) {
			latestEntry = &filteredList[i]
		}
	}
	return *latestEntry
}

// getLatestUnconfirmedBalance returns a LedgerEntryTxn containing an entry with
// the most recent unconfirmed balance.
func getLatestUnconfirmedBalance(stub *shim.ChaincodeStub, owner string, holder string) LedgerEntryTxn {
	// nostroLogger.Debug("ENTERED GETLATESTUNCONFIRMEDBALANCE")
	list, _ := getLedgerEntries(stub, true, true, true, false, true, false)

	// Filter by owner and holder
	var filteredList []LedgerEntryTxn
	for i := 0; i < len(list); i++ {
		if list[i].OrderingInstitution == owner && list[i].AccountWithInstitution == holder {
			filteredList = append(filteredList, list[i])
		}
	}

	// Find entry with most recent timeStamp
	latestEntry := &LedgerEntryTxn{}
	for i := 0; i < len(filteredList); i++ {
		if filteredList[i].TimestampCreated.Equal(latestEntry.TimestampCreated) {
			if filteredList[i].Sequence > latestEntry.Sequence {
				latestEntry = &filteredList[i]
			}
		} else if filteredList[i].TimestampCreated.After(latestEntry.TimestampCreated) {
			latestEntry = &filteredList[i]
		}
	}
	return *latestEntry
}

// getBalanceEntries returns a list of LedgerEntry that have balances. I.e.
// Confirmation messages and Funding messages.
// TODO: Possibly made redundant by getLedgerEntries()
func getBalanceEntries(stub *shim.ChaincodeStub) ([]LedgerEntry, error) {
	// nostroLogger.Debug("ENTERED GETBALANCEENTRIES")

	ledgerEntry := &LedgerEntry{}
	var ledgerList []LedgerEntry
	entryType := ""

	// Get a range of states from the ledger
	iterator, err := stub.RangeQueryState("0", "zzzzz")
	if err != nil {
		return nil, errors.New("Failed to get iterator")
	}
	// Retrieve first key and value
	key, value, err := iterator.Next()
	if err != nil {
		return nil, errors.New("Failed to retrieve key and value via iterator")
	}

	entryType = utility.IdentifyLedgerType(value)
	switch entryType {
	case "Confirmation":
		entry := &LedgerEntryConfirmation{}
		err = json.Unmarshal(value, entry)
		if err != nil {
			return nil, errors.New("Failed to unmarshal value into struct.")
		}
		ledgerEntry.Key = key
		ledgerEntry.Balance = entry.StatementAccountBalance
		ledgerEntry.Type = entryType
		ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
		ledgerEntry.OrderingInstitution = entry.OrderingInstitution
		// timeStamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", entry.Timestamp)
		// if err != nil {
		// 	return nil, errors.New("Could not parse time from ledger.")
		// }
		ledgerEntry.Timestamp = entry.Timestamp
		ledgerList = append(ledgerList, *ledgerEntry)

	case "Funding":
		entry := &LedgerEntryFunding{}
		err = json.Unmarshal(value, entry)
		if err != nil {
			return nil, errors.New("Failed to unmarshal value into struct.")
		}
		ledgerEntry.Key = key
		ledgerEntry.Balance = entry.StatementAccountBalance
		ledgerEntry.Type = entryType
		ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
		ledgerEntry.OrderingInstitution = entry.OrderingInstitution
		// timeStamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", entry.Timestamp)
		// if err != nil {
		// 	return nil, errors.New("Could not parse time from ledger.")
		// }
		ledgerEntry.Timestamp = entry.Timestamp
		ledgerList = append(ledgerList, *ledgerEntry)

	case "Fee":
		entry := &LedgerEntryFee{}
		err = json.Unmarshal(value, entry)
		if err != nil {
			return nil, errors.New("Failed to unmarshal value into struct.")
		}
		ledgerEntry.Key = key
		ledgerEntry.Balance = entry.StatementAccountBalance
		ledgerEntry.Type = entryType
		ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
		ledgerEntry.OrderingInstitution = entry.OrderingInstitution
		ledgerEntry.Timestamp = entry.Timestamp
		ledgerList = append(ledgerList, *ledgerEntry)
	}

	// Iterate and store subsequent keys and values
	for iterator.HasNext() {
		key, value, err := iterator.Next()
		if err != nil {
			return nil, errors.New("Failed to retrieve key and value via iterator")
		}
		entryType = utility.IdentifyLedgerType(value)
		switch entryType {
		case "Confirmation":
			entry := &LedgerEntryConfirmation{}
			err = json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntry.Key = key
			ledgerEntry.Balance = entry.StatementAccountBalance
			ledgerEntry.Type = entryType
			ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntry.OrderingInstitution = entry.OrderingInstitution
			ledgerEntry.Timestamp = entry.Timestamp
			ledgerList = append(ledgerList, *ledgerEntry)

		case "Funding":
			entry := &LedgerEntryFunding{}
			err = json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntry.Key = key
			ledgerEntry.Balance = entry.StatementAccountBalance
			ledgerEntry.Type = entryType
			ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntry.OrderingInstitution = entry.OrderingInstitution
			// timeStamp, err := time.Parse("2006-01-02 15:04:05 -0700 MST", entry.Timestamp)
			// if err != nil {
			// 	return nil, errors.New("Could not parse time from ledger.")
			// }
			ledgerEntry.Timestamp = entry.Timestamp
			ledgerList = append(ledgerList, *ledgerEntry)

		case "Fee":
			entry := &LedgerEntryFee{}
			err = json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntry.Key = key
			ledgerEntry.Balance = entry.StatementAccountBalance
			ledgerEntry.Type = entryType
			ledgerEntry.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntry.OrderingInstitution = entry.OrderingInstitution
			ledgerEntry.Timestamp = entry.Timestamp
			ledgerList = append(ledgerList, *ledgerEntry)
		}
	}
	iterator.Close()
	return ledgerList, nil
}

// selectiveAppend takes a value from the ledger and either ignores it or appends it to a []ledgerEntryTxn based on the bools provided.
func selectiveAppend(stub *shim.ChaincodeStub, ledgerList []LedgerEntryTxn, key string, value []byte, request bool, confirmation bool, funding bool, fee bool, direct bool, advice bool) ([]LedgerEntryTxn, error) {
	// nostroLogger.Debug("ENTERED SELECTIVEAPPEND")

	ledgerEntryTxn := &LedgerEntryTxn{}
	entryType := utility.IdentifyLedgerType(value)
	// nostroLogger.Debug("Type = ", entryType)

	switch entryType {
	case "REQUEST-RECORD":
		if request {
			// nostroLogger.Debug("*APPEND REQUEST")
			entry := &LedgerEntryRequest{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.RefKey = entry.ConfirmationID
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.SendingFI = entry.SendingFI
			ledgerEntryTxn.ReceivingFI = entry.ReceivingFI
			ledgerEntryTxn.Currency = entry.ValueDateCurrencyInterbankSettled.Currency
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.PaymentAmount = entry.Amount
			ledgerEntryTxn.FeeType = entry.FeeType
			ledgerEntryTxn.SendersCharge = entry.SendersCharge
			ledgerEntryTxn.ReceiversCharge = entry.ReceiversCharge
			ledgerEntryTxn.BenePays = entry.BenePays
			ledgerEntryTxn.Rebate = entry.Rebate
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.IsRejected = entry.IsRejected
			ledgerEntryTxn.RejectRationale = entry.RejectRationale
			ledgerEntryTxn.RejectTimestamp = entry.RejectTimestamp
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.ValueDate = entry.ValueDate
			ledgerEntryTxn.LocalTime = entry.TimestampLocal
			ledgerEntryTxn.MsgNum = entry.MsgNum
			ledgerEntryTxn.MsgType = entry.MsgType
			ledgerEntryTxn.SecondKey = entry.SecondKey
			ledgerEntryTxn.ThirdKey = entry.ThirdKey
			ledgerEntryTxn.FourthKey = entry.FourthKey
			ledgerEntryTxn.TransactionReferenceNumber = entry.TransactionReferenceNumber
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}

	case "CONFIRMATION-RECORD":
		if confirmation {
			// nostroLogger.Debug("*APPEND CONFIRMATION")
			entry := &LedgerEntryConfirmation{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.RefKey = entry.RequestID
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.SendingFI = entry.SendingFI
			ledgerEntryTxn.ReceivingFI = entry.ReceivingFI
			ledgerEntryTxn.Currency = entry.Currency
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.PaymentAmount = entry.Amount
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.FeeType = entry.FeeType
			ledgerEntryTxn.SendersCharge = entry.SendersCharge
			ledgerEntryTxn.ReceiversCharge = entry.ReceiversCharge
			ledgerEntryTxn.BenePays = entry.BenePays
			ledgerEntryTxn.Rebate = entry.Rebate
			ledgerEntryTxn.TimestampConfirmed = entry.Timestamp
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerEntryTxn.LocalTime = entry.LocalTimestamp
			ledgerEntryTxn.MsgNum = entry.MsgNum
			ledgerEntryTxn.MsgType = entry.MsgType
			vDate := fmt.Sprintf("%s%02d%02d", strconv.Itoa(entry.Timestamp.Year())[2:], entry.Timestamp.Month(), entry.Timestamp.Day())
			ledgerEntryTxn.ValueDate = vDate
			ledgerEntryTxn.TransactionReferenceNumber = entry.TransactionReferenceNumber
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}

	case "DIRECT-CREDIT-RECORD":
		if direct {
			// nostroLogger.Debug("*APPEND DIRECT CREDIT")
			entry := &LedgerEntryRequest{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.RefKey = entry.ConfirmationID
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.SendingFI = entry.SendingFI
			ledgerEntryTxn.ReceivingFI = entry.ReceivingFI
			ledgerEntryTxn.Currency = entry.ValueDateCurrencyInterbankSettled.Currency
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.PaymentAmount = entry.Amount
			ledgerEntryTxn.FeeType = entry.FeeType
			ledgerEntryTxn.SendersCharge = entry.SendersCharge
			ledgerEntryTxn.ReceiversCharge = entry.ReceiversCharge
			ledgerEntryTxn.BenePays = entry.BenePays
			ledgerEntryTxn.Rebate = entry.Rebate
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.TimestampConfirmed = entry.Timestamp
			ledgerEntryTxn.IsRejected = entry.IsRejected
			ledgerEntryTxn.RejectRationale = entry.RejectRationale
			ledgerEntryTxn.RejectTimestamp = entry.RejectTimestamp
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.ValueDate = entry.ValueDate
			ledgerEntryTxn.LocalTime = entry.TimestampLocal
			ledgerEntryTxn.MsgNum = entry.MsgNum
			ledgerEntryTxn.MsgType = entry.MsgType
			ledgerEntryTxn.SecondKey = entry.SecondKey
			ledgerEntryTxn.ThirdKey = entry.ThirdKey
			ledgerEntryTxn.FourthKey = entry.FourthKey
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}

	case "ADVICE":
		if advice {
			nostroLogger.Debug("*APPEND ADVICE")
			entry := &LedgerEntryRequest{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.RefKey = entry.ConfirmationID
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.SendingFI = entry.SendingFI
			ledgerEntryTxn.ReceivingFI = entry.ReceivingFI
			ledgerEntryTxn.Currency = entry.ValueDateCurrencyInterbankSettled.Currency
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.PaymentAmount = entry.Amount
			ledgerEntryTxn.FeeType = entry.FeeType
			ledgerEntryTxn.SendersCharge = entry.SendersCharge
			ledgerEntryTxn.ReceiversCharge = entry.ReceiversCharge
			ledgerEntryTxn.BenePays = entry.BenePays
			ledgerEntryTxn.Rebate = entry.Rebate
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.TimestampConfirmed = entry.Timestamp
			ledgerEntryTxn.IsRejected = entry.IsRejected
			ledgerEntryTxn.RejectRationale = entry.RejectRationale
			ledgerEntryTxn.RejectTimestamp = entry.RejectTimestamp
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.ValueDate = entry.ValueDate
			ledgerEntryTxn.LocalTime = entry.TimestampLocal
			ledgerEntryTxn.MsgNum = entry.MsgNum
			ledgerEntryTxn.MsgType = entry.MsgType
			ledgerEntryTxn.SecondKey = entry.SecondKey
			ledgerEntryTxn.ThirdKey = entry.ThirdKey
			ledgerEntryTxn.FourthKey = entry.FourthKey
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}

	case "FUNDING-RECORD":
		if funding {
			// nostroLogger.Debug("*APPEND FUNDING")
			entry := &LedgerEntryFunding{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.PaymentAmount = entry.FundingAmount
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.TimestampConfirmed = entry.Timestamp
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}

	case "FEE-RECORD":
		if fee {
			// nostroLogger.Debug("*APPEND FEE")
			entry := &LedgerEntryFee{}
			err := json.Unmarshal(value, entry)
			if err != nil {
				return nil, errors.New("Failed to unmarshal value into struct.")
			}
			ledgerEntryTxn.Key = key
			ledgerEntryTxn.RefKey = entry.ConfirmationID
			ledgerEntryTxn.StatementAccountBalance = entry.StatementAccountBalance
			ledgerEntryTxn.IndicativeBalance = entry.IndicativeBalance
			ledgerEntryTxn.Type = entryType
			ledgerEntryTxn.OrderingInstitution = entry.OrderingInstitution
			ledgerEntryTxn.AccountWithInstitution = entry.AccountWithInstitution
			ledgerEntryTxn.FeeType = entry.FeeType
			ledgerEntryTxn.SendersCharge = entry.SendersCharge
			ledgerEntryTxn.ReceiversCharge = entry.ReceiversCharge
			ledgerEntryTxn.BenePays = entry.BenePays
			ledgerEntryTxn.Rebate = entry.Rebate
			ledgerEntryTxn.PaymentAmount = entry.PaymentAmount
			ledgerEntryTxn.TimestampConfirmed = entry.Timestamp
			ledgerEntryTxn.Sequence = entry.Sequence
			ledgerEntryTxn.TimestampCreated = entry.Timestamp
			ledgerEntryTxn.UnconfirmedSortTime = entry.Timestamp
			ledgerList = append(ledgerList, *ledgerEntryTxn)
		}
	}
	return ledgerList, nil
}

func (t *SimpleChaincode) getReceivedUnconfirmedPayments(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Expecting Bank Name.")
	}
	bankName := args[0]

	var payReceivedUnconfirmed []LedgerEntryTxn
	payReceivedUnconNotEmpty := false

	var result []byte
	// (1) Retrieve the ledger entry list.
	paymentInstructions, err := getLedgerEntries(stub, true, false, false, false, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(paymentInstructions) != 0 {
		for i := 0; i < len(paymentInstructions); i++ {
			if paymentInstructions[i].AccountWithInstitution == bankName {
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
					payReceivedUnconfirmed = append(payReceivedUnconfirmed, paymentInstructions[i])
				}
			}
		}

		if len(payReceivedUnconfirmed) != 0 {
			payReceivedUnconNotEmpty = true
		}
	}

	if payReceivedUnconNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(payReceivedUnconfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(payReceivedUnconfirmed)

		tmp, _ := json.MarshalIndent(payReceivedUnconfirmed, "", "  ")
		result = append(result, tmp...)
	}
	return result, nil

}

func (t *SimpleChaincode) getRejectedPaymentInstructions(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Expecting Bank Name.")
	}
	bankName := args[0]

	var rejectedPayInstructions []LedgerEntryTxn
	rejectedPayInstructionsNotEmpty := false

	var result []byte
	// (1) Retrieve the ledger entry list.
	paymentInstructions, err := getLedgerEntries(stub, true, false, false, false, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(paymentInstructions) != 0 {
		for i := 0; i < len(paymentInstructions); i++ {
			if paymentInstructions[i].OrderingInstitution == bankName {
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == true {
					rejectedPayInstructions = append(rejectedPayInstructions, paymentInstructions[i])
				}
			}
		}

		if len(rejectedPayInstructions) != 0 {
			rejectedPayInstructionsNotEmpty = true
		}
	}

	if rejectedPayInstructionsNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(rejectedPayInstructions))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(rejectedPayInstructions)

		tmp, _ := json.MarshalIndent(rejectedPayInstructions, "", "  ")
		result = append(result, tmp...)
	}
	return result, nil

}

// getTransactionSummary takes a bank name and returns a list of payment instructions created by it, and addressed to it. Also includes funding messages relating to its statement accounts.
func (t *SimpleChaincode) getTransactionSummary(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	// nostroLogger.Debug("ENTERED GETTRANSACTIONSUMMARY")

	if len(args) != 1 {
		return nil, errors.New("Expecting Bank Name.")
	}
	bankName := args[0]
	var paySentUnconfirmed []LedgerEntryTxn
	var paySentConfirmed []LedgerEntryTxn
	var paySentRejected []LedgerEntryTxn
	var payReceivedUnconfirmed []LedgerEntryTxn
	var payReceivedConfirmed []LedgerEntryTxn
	var payReceivedRejected []LedgerEntryTxn
	var fundingList []LedgerEntryTxn
	var filteredFundingMessages [][]LedgerEntryTxn
	var feeList []LedgerEntryTxn
	var filteredFeeMessages [][]LedgerEntryTxn

	paySentUnconNotEmpty := false
	paySentConNotEmpty := false
	paySentRejNotEmpty := false
	payReceivedConNotEmpty := false
	payReceivedUnconNotEmpty := false
	payReceivedRejNotEmpty := false
	fundingNotEmpty := false
	feeNotEmpty := false

	var result []byte
	var holdingBanks []string
	var feeHoldingBanks []string

	width := 86 // Sets pretty printing format width for CLI.

	// (1) Retrieve the ledger entry list.
	paymentInstructions, err := getLedgerEntries(stub, true, false, false, false, true, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(paymentInstructions) != 0 {
		// (2) Create list of Payment Instructions created by Bank X (i.e. SendingFI = Bank X).
		// (3) Group into confirmed and unconfirmed. Sort by descending date. Include FromBank(always Bank X), ToBank, Amount, DatePosted, DateConfirmed.
		for i := 0; i < len(paymentInstructions); i++ {
			if paymentInstructions[i].SendingFI == bankName {
				if paymentInstructions[i].RefKey != "" {
					paySentConfirmed = append(paySentConfirmed, paymentInstructions[i])
				}
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
					paySentUnconfirmed = append(paySentUnconfirmed, paymentInstructions[i])
				}
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == true {
					paySentRejected = append(paySentRejected, paymentInstructions[i])
				}
			}
		}
		if len(paySentConfirmed) != 0 {
			paySentConNotEmpty = true
		}

		if len(paySentUnconfirmed) != 0 {
			paySentUnconNotEmpty = true
		}
		if len(paySentRejected) != 0 {
			paySentRejNotEmpty = true
		}

		// (4) Create list of Payment Instructions addressed to Bank X (i.e. ReceivingFI = Bank X).
		// (5) Group into confirmed and unconfired. Sort by descending date. Include FromBank, ToBank(always Bank X), Amount, DatePosted, DateConfirmed.
		for i := 0; i < len(paymentInstructions); i++ {
			if paymentInstructions[i].ReceivingFI == bankName {
				if paymentInstructions[i].RefKey != "" {
					payReceivedConfirmed = append(payReceivedConfirmed, paymentInstructions[i])
				}
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
					payReceivedUnconfirmed = append(payReceivedUnconfirmed, paymentInstructions[i])
				}
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == true {
					payReceivedRejected = append(payReceivedRejected, paymentInstructions[i])
				}
			}
		}
		if len(payReceivedConfirmed) != 0 {
			payReceivedConNotEmpty = true
		}
		if len(payReceivedUnconfirmed) != 0 {
			payReceivedUnconNotEmpty = true
		}
		if len(payReceivedRejected) != 0 {
			payReceivedRejNotEmpty = true
		}
	}

	// (6) Create list of Funding messages related to statement accounts owned by Bank X. Sort descending date.
	fundingMessages, err := getLedgerEntries(stub, false, false, true, false, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(fundingMessages) != 0 {
		fundingNotEmpty = true
		for i := 0; i < len(fundingMessages); i++ {
			if fundingMessages[i].OrderingInstitution == bankName {
				fundingList = append(fundingList, fundingMessages[i])
			}
		}
		// TODO: Try sorting here.
		// sort.Sort(fundingList)

		// Get unique list of holding banks.
		if len(fundingList) > 0 {
			holdingBanks = append(holdingBanks, fundingList[0].AccountWithInstitution)
			if len(fundingList) > 1 {
				for i := 1; i < len(fundingList); i++ {
					if contains(holdingBanks, fundingList[i].AccountWithInstitution) {
					} else {
						holdingBanks = append(holdingBanks, fundingList[i].AccountWithInstitution)
					}
				}
			}
		}

		filteredFundingMessages = make([][]LedgerEntryTxn, len(holdingBanks))
		var temp []LedgerEntryTxn

		for i := 0; i < len(holdingBanks); i++ {
			for j := 0; j < len(fundingList); j++ {
				if holdingBanks[i] == fundingList[j].AccountWithInstitution {
					temp = append(temp, fundingList[j])
				}
			}
			filteredFundingMessages[i] = temp
			temp = nil
		}
	}

	// (7) Create list of Fee messages related to statement accounts owned by Bank X. Sort descending date.
	feeMessages, err := getLedgerEntries(stub, false, false, false, true, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(feeMessages) != 0 {
		feeNotEmpty = true
		// loggerCT.Trace.Println("CHECK02")
		for i := 0; i < len(feeMessages); i++ {
			if feeMessages[i].OrderingInstitution == bankName {
				feeList = append(feeList, feeMessages[i])
			}
		}
		// Get unique list of holding banks.
		if len(feeList) > 0 {
			feeHoldingBanks = append(feeHoldingBanks, feeList[0].AccountWithInstitution)
			if len(feeList) > 1 {
				for i := 1; i < len(feeList); i++ {
					if contains(feeHoldingBanks, feeList[i].AccountWithInstitution) {
					} else {
						feeHoldingBanks = append(feeHoldingBanks, feeList[i].AccountWithInstitution)
					}
				}
			}
		}

		filteredFeeMessages = make([][]LedgerEntryTxn, len(feeHoldingBanks))
		var temp []LedgerEntryTxn

		for i := 0; i < len(feeHoldingBanks); i++ {
			for j := 0; j < len(feeList); j++ {
				if feeHoldingBanks[i] == feeList[j].AccountWithInstitution {
					temp = append(temp, feeList[j])
				}
			}
			filteredFeeMessages[i] = temp
			temp = nil
		}
	}

	// Format lists and combine into a single return variable of type []byte
	// PAYMENTS SENT
	titleC := fmt.Sprintf("PAYMENT INSTRUCTIONS CREATED BY %s", bankName)
	textLen := len(titleC)
	a := padLeft(titleC, " ", (width/2)+(textLen/2))
	b := padRight(a, " ", width-2)
	border := ""
	border = padLeft(border, "#", width)
	titleOwnedList := fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleOwnedList...)
	// Sub-heading --> Confirmed
	titleC = "CONFIRMED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleConfirmed := fmt.Sprintf("\n%s\n", b)
	result = append(result, titleConfirmed...)
	// Append confirmed list
	if paySentConNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(paySentConfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(paySentConfirmed)

		tmp, _ := json.MarshalIndent(paySentConfirmed, "", "  ")
		result = append(result, tmp...)
	}

	// Sub-heading --> Unconfirmed
	titleC = "UNCONFIRMED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleUnconfirmed := fmt.Sprintf("\n%s\n", b)
	result = append(result, titleUnconfirmed...)
	// Append confirmed list
	if paySentUnconNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(paySentUnconfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(paySentUnconfirmed)

		tmp, _ := json.MarshalIndent(paySentUnconfirmed, "", "  ")
		result = append(result, tmp...)
	}

	// Sub-heading --> Rejected
	titleC = "REJECTED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleRejected := fmt.Sprintf("\n%s\n", b)
	result = append(result, titleRejected...)
	// Append confirmed list
	if paySentRejNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(paySentRejected))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(paySentRejected)

		tmp, _ := json.MarshalIndent(paySentRejected, "", "  ")
		result = append(result, tmp...)
	}

	// PAYMENTS RECEIVED
	titleC = fmt.Sprintf("PAYMENT INSTRUCTIONS RECEIVED BY %s", bankName)
	textLen = len(titleC)
	a = padLeft(titleC, " ", (width/2)+(textLen/2))
	b = padRight(a, " ", width-2)
	titleOwnedList = fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleOwnedList...)
	// Sub-heading --> Confirmed
	titleC = "CONFIRMED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleConfirmed = fmt.Sprintf("\n%s\n", b)
	result = append(result, titleConfirmed...)
	// Append confirmed list
	if payReceivedConNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(payReceivedConfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(payReceivedConfirmed)

		tmp, _ := json.MarshalIndent(payReceivedConfirmed, "", "  ")
		result = append(result, tmp...)
	}

	// Sub-heading --> Unconfirmed
	titleC = "UNCONFIRMED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleUnconfirmed = fmt.Sprintf("\n%s\n", b)
	result = append(result, titleUnconfirmed...)
	// Append confirmed list
	if payReceivedUnconNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(payReceivedUnconfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(payReceivedUnconfirmed)

		tmp, _ := json.MarshalIndent(payReceivedUnconfirmed, "", "  ")
		result = append(result, tmp...)
	}

	// Sub-heading --> Rejected
	titleC = "REJECTED"
	textLen = len(titleC)
	a = padLeft(titleC, "*", (width/2)+(textLen/2))
	b = padRight(a, "*", width-2)
	titleRejected = fmt.Sprintf("\n%s\n", b)
	result = append(result, titleRejected...)
	// Append confirmed list
	if payReceivedRejNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(payReceivedRejected))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(payReceivedRejected)

		tmp, _ := json.MarshalIndent(payReceivedRejected, "", "  ")
		result = append(result, tmp...)
	}

	// FUNDING MESSAGES
	titleC = fmt.Sprintf("FUNDING MESSAGES FOR %s's ACCOUNTS", bankName)
	textLen = len(titleC)
	a = padLeft(titleC, " ", (width/2)+(textLen/2))
	b = padRight(a, " ", width-2)
	titleOwnedList = fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleOwnedList...)
	if fundingNotEmpty {
		for i := 0; i < len(filteredFundingMessages); i++ {
			if len(filteredFundingMessages[i]) > 0 {
				titleC = fmt.Sprintf("ACCOUNT AT %s", filteredFundingMessages[i][0].AccountWithInstitution)
				textLen = len(titleC)
				a = padLeft(titleC, "*", (width/2)+(textLen/2))
				b = padRight(a, "*", width-2)
				titleC = fmt.Sprintf("\n%s\n", b)
				result = append(result, titleC...)
				// // Sort messages
				// sort.Sort(ByDate(filteredFundingMessages[i]))

				// Sort by TimestampCreated then Sequence
				// Closures that order the LedgerEntryTxn structure.
				timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
					return c1.TimestampCreated.Before(c2.TimestampCreated)
				}
				sequence := func(c1, c2 *LedgerEntryTxn) bool {
					return c1.Sequence < c2.Sequence
				}
				OrderedBy(timestampCreated, sequence).Sort(filteredFundingMessages[i])

				// Append funding messages held at bank x
				tmp, _ := json.MarshalIndent(filteredFundingMessages[i], "", "  ")
				result = append(result, tmp...)
			}
		}
	}
	// FEE MESSAGES
	titleC = fmt.Sprintf("FEE MESSAGES FOR %s's ACCOUNTS", bankName)
	textLen = len(titleC)
	a = padLeft(titleC, " ", (width/2)+(textLen/2))
	b = padRight(a, " ", width-2)
	titleOwnedList = fmt.Sprintf("\n%s\n#%s#\n%s\n", border, b, border)
	result = append(result, titleOwnedList...)
	if feeNotEmpty {
		for i := 0; i < len(filteredFeeMessages); i++ {
			if len(filteredFeeMessages[i]) > 0 {
				titleC = fmt.Sprintf("FEE MESSAGES FOR ACCOUNT AT %s", filteredFeeMessages[i][0].AccountWithInstitution)
				a = padLeft(titleC, "*", (width/2)+(textLen/2))
				b = padRight(a, "*", width-2)
				titleC = fmt.Sprintf("\n%s\n", b)
				result = append(result, titleC...)
				// // Sort messages
				// sort.Sort(ByDate(filteredFeeMessages[i]))

				// Sort by TimestampCreated then Sequence
				// Closures that order the LedgerEntryTxn structure.
				timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
					return c1.TimestampCreated.Before(c2.TimestampCreated)
				}
				sequence := func(c1, c2 *LedgerEntryTxn) bool {
					return c1.Sequence < c2.Sequence
				}
				OrderedBy(timestampCreated, sequence).Sort(filteredFeeMessages[i])

				// Append funding messages held at bank x
				tmp, _ := json.MarshalIndent(filteredFeeMessages[i], "", "  ")
				result = append(result, tmp...)
			}
		}
	}
	return result, nil
}

func calculateBankFees(feeType string, sendingFI string, receivingFI string, bookInt string, currency string) (senderFIOwesReceivingFI float64, benePays float64, rebate float64) {
	// nostroLogger.Debug("ENTERED CALCULATEBANKFEES")

	senderFIOwesReceivingFI = 0
	benePays = 0
	rebate = 0

	switch feeType {

	case "OUR":
		switch sendingFI {
		case "ANZ":
			switch receivingFI {
			case "WF", "Wells Fargo":
				senderFIOwesReceivingFI = 6.55
			}
		case "WF", "Wells Fargo":
			switch receivingFI {
			case "ANZ":
				switch bookInt {
				case "BOOK":
					senderFIOwesReceivingFI = 15
				case "INTERMEDIARY":
					switch currency {
					case "same":
						senderFIOwesReceivingFI = 25
					case "diff":
						senderFIOwesReceivingFI = 35
					}
				}
			}
		}

	case "BEN", "SHA":
		switch sendingFI {
		case "ANZ":
			switch receivingFI {
			case "WF", "Wells Fargo":
				senderFIOwesReceivingFI = 0.55
				benePays = 25
			}
		case "WF", "Wells Fargo":
			switch receivingFI {
			case "ANZ":
				switch bookInt {
				case "BOOK":
					benePays = 15
				case "INTERMEDIARY":
					switch currency {
					case "same":
						benePays = 25
					case "diff":
						benePays = 35
					}
				}
			}
		}
	}

	return senderFIOwesReceivingFI, benePays, rebate
}

// matchUnconfirmedTransactions reads the settled transactions and invokes the addPaymentConfirmation with the hashstrings of keys of pending payments
func (t *SimpleChaincode) matchUnconfirmedTransactions(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if len(args) != 14 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11:\nAmount\nOrdering Institution\nAccount With Institution\nUTC Time\nFee Type\nValue Date\nSender's Reference\nCurrency\nLocal Time\nMSG#\nMSG Type\nBookInt\nF33BCCY\nF33BAMT\n")
	}
	nostroLogger.Debug(" ================================================  ENTERED MATCHUNCONFIRMEDTRANSACTIONS ============================================================")

	// Parse arguments
	var F32AAMT float64
	var F33BAMT float64
	var entryType string
	sendingFI := args[1]
	receivingFI := args[2]
	timestampString := args[3]
	// timeStamp, e := time.Parse(time.RFC3339Nano, args[3])
	// if e != nil {
	// 	return nil, errors.New("Could not parse the time...")
	// }
	// feeType := args[4]
	vDate := args[5]
	F20 := args[6]
	F32ACCY := args[7]
	localTime := args[8]
	msgNum := args[9]
	msgType := args[10]
	bookInt := args[11]
	F33BCCY := args[12]

	// No confirmation needed if type = ADVICE or DIRECT CREDIT
	_, _, crdr := determineAccountCrDr(sendingFI, receivingFI, F32ACCY, msgType)
	switch crdr {
	case "DR":
		F32AAMT, _ = strconv.ParseFloat("-"+args[0], 64)
		F33BAMT, _ = strconv.ParseFloat("-"+args[13], 64)
		entryType = "REQUEST-RECORD"
	case "CR":
		F32AAMT, _ = strconv.ParseFloat(args[0], 64)
		F33BAMT, _ = strconv.ParseFloat(args[13], 64)
		entryType = "DIRECT-CREDIT-RECORD"
	case "Do Nothing":
		nostroLogger.Debug("NO CONFIRMATION REQUIRED FOR ADVICE MESSAGES")
		entryType = "ADVICE"
		return nil, nil
	}

	// 1. SET KEY MATCHING COMBINATIONS
	// SET FIRST KEY arguments: Senders Ref, F33BAMT, F33BCCY, and Value Date.
	firstKeyArgs := &PaymentInstructionKey{
		SenderReference: F20,
		F33BAMT:         F33BAMT,
		F33BCCY:         F33BCCY,
		ValueDate:       vDate,
	}

	// SET SECOND KEY arguments: F32AAMT, F32ACCY, F33BAMT, F33BCCY, value date, and BookInt.
	secondKeyArgs := &PaymentInstructionKey{
		F32AAMT:   F32AAMT,
		F32ACCY:   F32ACCY,
		F33BAMT:   F33BAMT,
		F33BCCY:   F33BCCY,
		ValueDate: vDate,
		BookInt:   bookInt,
	}

	// SET THIRD KEY arguments: F33BAMT, F33BCCY, value date, and BookInt.
	thirdKeyArgs := &PaymentInstructionKey{
		F33BAMT:   F33BAMT,
		F33BCCY:   F33BCCY,
		ValueDate: vDate,
		BookInt:   bookInt,
	}

	// SET FOURTH KEY arguments: F33BAMT, F33BCCY, value date, and BookInt.
	fourthKeyArgs := &PaymentInstructionKey{
		F32AAMT:   F32AAMT,
		F32ACCY:   F32ACCY,
		ValueDate: vDate,
	}

	// 2. CREATE CONFIRMATION ARGS FOR ALL KEYS
	confirmationArgs1, err := createConfirmationArgs(firstKeyArgs, timestampString, receivingFI, localTime, msgNum, msgType, F20)
	if err != nil {
		return nil, err
	}
	confirmationArgs2, err := createConfirmationArgs(secondKeyArgs, timestampString, receivingFI, localTime, msgNum, msgType, F20)
	if err != nil {
		return nil, err
	}
	confirmationArgs3, err := createConfirmationArgs(thirdKeyArgs, timestampString, receivingFI, localTime, msgNum, msgType, F20)
	if err != nil {
		return nil, err
	}
	confirmationArgs4, err := createConfirmationArgs(fourthKeyArgs, timestampString, receivingFI, localTime, msgNum, msgType, F20)
	if err != nil {
		return nil, err
	}
	nostroLogger.Debug("")
	nostroLogger.Debug("FIRST KEY INFO ", firstKeyArgs)
	nostroLogger.Debug("FIRST KEY ", confirmationArgs1[0])
	nostroLogger.Debug("")
	nostroLogger.Debug("SECOND KEY INFO ", secondKeyArgs)
	nostroLogger.Debug("SECOND KEY ", confirmationArgs2[0])
	nostroLogger.Debug("")
	nostroLogger.Debug("THIRD KEY INFO ", thirdKeyArgs)
	nostroLogger.Debug("THIRD KEY ", confirmationArgs3[0])
	nostroLogger.Debug("")
	nostroLogger.Debug("FOURTH KEY INFO ", fourthKeyArgs)
	nostroLogger.Debug("FOURTH KEY ", confirmationArgs4[0])

	val, er := t.addPaymentConfirmation(stub, confirmationArgs1)
	if er == nil {
		return val, nil
	}

	// SECOND, THIRD & FOURTH KEY MATCHES.
	// Retrieve every request entry and sequentially compare second
	// Logic to handle Secondary Match - Rule 14
	nostroLogger.Debug("FIRST KEY NOT MATCHED. TRYING SECOND KEY")

	paymentInstructions, err := getLedgerEntries(stub, true, false, false, false, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(paymentInstructions) != 0 {
		for i := 0; i < len(paymentInstructions); i++ {
			// nostroLogger.Debug("receivingFI ", paymentInstructions[i].AccountWithInstitution, receivingFI)
			if paymentInstructions[i].AccountWithInstitution == receivingFI {
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
					if paymentInstructions[i].TransactionReferenceNumber == F20 {
						nostroLogger.Debug("SUCCESS: F20 MATCHED")
						confirmationArgs := []string{paymentInstructions[i].Key, timestampString, receivingFI, localTime, msgNum, msgType, F20}
						val, er := t.addPaymentConfirmation(stub, confirmationArgs)
						if er == nil {
							return val, nil
						}
					}
					if paymentInstructions[i].SecondKey == confirmationArgs2[0] {
						nostroLogger.Debug("SUCCESS: SECOND KEY MATCHED")
						confirmationArgs := []string{paymentInstructions[i].Key, timestampString, receivingFI, localTime, msgNum, msgType, F20}
						val, er := t.addPaymentConfirmation(stub, confirmationArgs)
						if er == nil {
							return val, nil
						}
					}
					if paymentInstructions[i].ThirdKey == confirmationArgs3[0] {
						nostroLogger.Debug("SUCCESS: THIRD KEY MATCHED")
						confirmationArgs := []string{paymentInstructions[i].Key, timestampString, receivingFI, localTime, msgNum, msgType, F20}
						val, er := t.addPaymentConfirmation(stub, confirmationArgs)
						if er == nil {
							return val, nil
						}
					}
					if paymentInstructions[i].FourthKey == confirmationArgs4[0] {
						nostroLogger.Debug("SUCCESS: FOURTH KEY MATCHED")
						confirmationArgs := []string{paymentInstructions[i].Key, timestampString, receivingFI, localTime, msgNum, msgType, F20}
						val, er := t.addPaymentConfirmation(stub, confirmationArgs)
						if er == nil {
							return val, nil
						}
					}
				}
			}
		}
	}

	// if len(paymentInstructions) != 0 {
	// 	//Check if payment amounts match.
	// 	for i := 0; i < len(paymentInstructions); i++ {
	// 		if paymentInstructions[i].AccountWithInstitution == receivingFI {
	// 			if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
	// 				if paymentInstructions[i].PaymentAmount == amount {
	// 					confirmationArgs := []string{paymentInstructions[i].Key, timestampString, receivingFI, localTime, msgNum, msgType, F20}
	// 					val, er := t.addPaymentConfirmation(stub, confirmationArgs)
	// 					if er == nil {
	// 						return val, nil
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	nostroLogger.Error("================ Following Transaction did not find a MATCH ================")
	nostroLogger.Error("MSG# ", msgNum)
	nostroLogger.Error("MSG_TYPE ", msgType)
	nostroLogger.Error("ENTRY_TYPE ", entryType)
	nostroLogger.Error("Value Date ", vDate)
	nostroLogger.Error("Amount ", F32AAMT)
	nostroLogger.Error("SendingFI ", sendingFI)
	nostroLogger.Error("ReceivingFI ", receivingFI)
	nostroLogger.Error("Senders Reference ", F20)
	nostroLogger.Error(" ============================================================================")
	return nil, errors.New("Failed to invoke addPaymentConfirmation.")
}

func createConfirmationArgs(settlementData *PaymentInstructionKey, time string, bankName string, localTime string, msgNum string, msgType string, F20 string) ([]string, error) {
	// Convert ConfirmationRequest struct to JSON format as []byte.
	// nostroLogger.Debug("ENTERED CREATECONFIRMATIONARGS")
	out, err := json.MarshalIndent(settlementData, "", "  ")
	if err != nil {
		return nil, errors.New("Failed to marshal settlementData for matchUnconfirmedTransactions.")
	}

	// Generate database key based on settlementData.
	key := utility.GenerateKey(string(out))

	// Convert key from []byte to string in hex format for readability.
	requestKeyString := fmt.Sprintf("%x", key[0:])

	confirmationArgs := []string{requestKeyString, time, bankName, localTime, msgNum, msgType, F20}

	return confirmationArgs, nil
}

// getSentUnconfirmedPayments returns all the unconfirmed payments sent by the specified bank.
func (t *SimpleChaincode) getSentUnconfirmedPayments(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	nostroLogger.Debug("Invoked function getSentUnconfirmedPayments")
	if len(args) != 1 {
		return nil, errors.New("Expecting Bank Name.")
	}
	bankName := args[0]

	var payReceivedUnconfirmed []LedgerEntryTxn
	payReceivedUnconNotEmpty := false

	var result []byte
	// (1) Retrieve the ledger entry list.
	paymentInstructions, err := getLedgerEntries(stub, true, false, false, false, false, false)
	if err != nil {
		return nil, errors.New("getLedgerEntries() FAILED")
	}
	if len(paymentInstructions) != 0 {
		for i := 0; i < len(paymentInstructions); i++ {
			if paymentInstructions[i].OrderingInstitution == bankName {
				if paymentInstructions[i].RefKey == "" && paymentInstructions[i].IsRejected == false {
					payReceivedUnconfirmed = append(payReceivedUnconfirmed, paymentInstructions[i])
				}
			}
		}

		if len(payReceivedUnconfirmed) != 0 {
			payReceivedUnconNotEmpty = true
		}
	}

	if payReceivedUnconNotEmpty {
		// // Sort messages
		// sort.Sort(ByDate(payReceivedUnconfirmed))

		// Sort by TimestampCreated then Sequence
		// Closures that order the LedgerEntryTxn structure.
		timestampCreated := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.TimestampCreated.Before(c2.TimestampCreated)
		}
		sequence := func(c1, c2 *LedgerEntryTxn) bool {
			return c1.Sequence < c2.Sequence
		}
		OrderedBy(timestampCreated, sequence).Sort(payReceivedUnconfirmed)

		tmp, _ := json.MarshalIndent(payReceivedUnconfirmed, "", "  ")
		result = append(result, tmp...)
	}
	return result, nil

}

// determineAccountCrDr determines which account to modify and whether it is a debit or a credit.
func determineAccountCrDr(sendingFI string, receivingFI string, currency string, msgType string) (orderingInstitution, accountWithInstitution, crdr string) {

	switch {
	case receivingFI == "ANZ" && currency == "AUD": // DR OI's AUD account at ANZ
		orderingInstitution = sendingFI
		accountWithInstitution = receivingFI
		crdr = "DR"
		return
	case receivingFI == "ANZ" && currency == "USD":
		orderingInstitution = receivingFI
		accountWithInstitution = sendingFI
		if msgType == "103" { // for MT103s, CR ANZ's USD account at WF
			crdr = "CR"
		} else { // for MT2xx, Do nothing
			crdr = "Do Nothing"
		}
		return
	case receivingFI == "WF" && currency == "USD": // DR OI's USD account at WF
		orderingInstitution = sendingFI
		accountWithInstitution = receivingFI
		crdr = "DR"
		return
	case receivingFI == "WF" && currency == "AUD":
		orderingInstitution = receivingFI
		accountWithInstitution = sendingFI
		if msgType == "103" { // for MT103s, CR WF's AUD account at WF
			crdr = "CR"
		} else { // for MT2xx, Do nothing
			crdr = "Do Nothing"
		}
		return
	}

	// switch {
	// case receivingFI == "ANZ" && currency == "AUD": // DR WF's AUD account at ANZ
	// 	orderingInstitution = sendingFI
	// 	accountWithInstitution = receivingFI
	// 	crdr = "DR"
	// 	return
	// case receivingFI == "ANZ" && currency == "USD": // CR ANZ's USD account at WF
	// 	orderingInstitution = receivingFI
	// 	accountWithInstitution = sendingFI
	// 	crdr = "CR"
	// 	return
	// case receivingFI == "WF" && currency == "USD": // DR ANZ's USD account at WF
	// 	orderingInstitution = sendingFI
	// 	accountWithInstitution = receivingFI
	// 	crdr = "DR"
	// 	return
	// case receivingFI == "WF" && currency == "AUD": // CR WF's AUD account at ANZ
	// 	orderingInstitution = receivingFI
	// 	accountWithInstitution = sendingFI
	// 	crdr = "CR"
	// 	return
	// }
	return
}
