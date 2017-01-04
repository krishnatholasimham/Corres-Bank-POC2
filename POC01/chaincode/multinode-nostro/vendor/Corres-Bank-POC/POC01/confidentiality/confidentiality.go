/*
Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)
833 Collins Street, Docklands 3008, ABN 11 005 357 522.
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by Heshan Peiris <heshan.peiris@anz.com> August 2016
*/

package confidentiality

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/crypto/attr"
	"github.com/hyperledger/fabric/core/crypto"
	"github.com/spf13/viper"
	"fmt"
	"github.com/op/go-logging"
	"encoding/asn1"
	"crypto/x509"
	"strings"
	"errors"
	"encoding/json"
)

var nostroLogger = logging.MustGetLogger("nostro")

var (
	vp0Admin  crypto.Client
	vp1Admin  crypto.Client
	vp2Admin  crypto.Client
	vp3Admin  crypto.Client
	peerID string
)

type SecurityMetaData struct {
	TxnSigma []byte
	TxnCert  []byte
	TxnBinding []byte
	InkSigma []byte
}

type EntrySecurityInfo struct {
	Signature	[]byte
}

//Evaluates whether "Privacy" is turned on in core.yaml
func IsConfidential() (bool){
	viper.SetConfigName("core")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/opt/gopath/src/github.com/hyperledger/fabric/peer")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file [%s] \n", err))
	}


	return viper.GetBool("security.privacy")
}

func setViperVariables(){
	viper.Set("peer.fileSystemPath","/var/hyperledger/production")
	viper.Set("peer.pki.eca.paddr","172.17.0.1:50051")
	viper.Set("peer.pki.tca.paddr","172.17.0.1:50051")
	viper.Set("peer.pki.tlsca.paddr","172.17.0.1:50051")
}

func Setup() {

	peerID = viper.GetString("peer.id")

	setViperVariables()

	// Logging
	/*var formatter = logging.MustStringFormatter(
		`%{color}[%{module}] %{shortfunc} [%{shortfile}] -> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	logging.SetFormatter(formatter)

	logging.SetLevel(logging.DEBUG, "peer")
	logging.SetLevel(logging.DEBUG, "chaincode")
	logging.SetLevel(logging.DEBUG, "cryptoain")*/

	// Init the crypto layer
	if err := crypto.Init(); err != nil {
		panic(fmt.Errorf("Failed initializing the crypto layer [%s]", err))
	}
}

func InitClients() error{
	var err error

	if(strings.Compare(peerID,"jdoe")==0 || strings.Compare(peerID,"vp0")==0){
		//Initialize vp0Admin
		if err := crypto.RegisterClient("vp0Admin", nil, "vp0Admin", "vp0admin_secret"); err != nil {
			return err
		}
		vp0Admin, err = crypto.InitClient("vp0Admin", nil)
		if err != nil {
			return err
		}

		//Initialize vp1Admin
		if err := crypto.RegisterClient("vp1Admin", nil, "vp1Admin", "vp1admin_secret"); err != nil {
			return err
		}
		vp1Admin, err = crypto.InitClient("vp1Admin", nil)
		if err != nil {
			return err
		}

		//Initialize vp2Admin
		if err := crypto.RegisterClient("vp2Admin", nil, "vp2Admin", "vp2admin_secret"); err != nil {
			return err
		}
		vp2Admin, err = crypto.InitClient("vp2Admin", nil)
		if err != nil {
			return err
		}

		//Initialize vp3Admin
		if err := crypto.RegisterClient("vp3Admin", nil, "vp3Admin", "vp3admin_secret"); err != nil {
			return err
		}
		vp3Admin, err = crypto.InitClient("vp3Admin", nil)
		if err != nil {
			return err
		}
	}
	viper.Reset()
	setViperVariables()

	return nil
}

//This function is called by all methods within the chaincode which creates a ledger entry.
// The function performs the following
/*
1. Authorises the invoking user based on the affiliation
2. Retrieves the corresponding Admin Crypto Client
3. Generates the transaction signature
 */
func GetInvokeTransactionSignature(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	if(isRegisteredInvoker(stub)){
		//Retrieve the invoker initiating the transaction
		invoker,error := getInvoker(stub)
		if error != nil {
			return nil, errors.New("Error retrieving the Invoker")
		}

		//Get invoker's Enrollment Certificate
		eCert, error := getInvokerEcert(invoker)
		if error != nil {
			return nil, errors.New("Error retrieving the Invoker Enrollment Certificate")
		}

		//Check if the Invoker has the correct affiliation
		affiliation, error := getInvokerAffiliation(eCert)
		if error != nil {
			return nil, errors.New("Error retrieving the Invoker Enrollment Certificate")
		}

		if ((strings.Compare(affiliation,strings.ToLower(args[0]))) != 0){
			return nil, errors.New("Invoker does not have the necessary affiliation")
		}
	}

	recipientAdministrator := getAdminUser(args[1])
	invokerAdministrator := getAdminUser(args[0])

	tSigma, tCertificate, tBinding, _ := generateTransactionSigma(recipientAdministrator, args[2])
	invokerAdminSigma,_ := generateInvokerSigma(invokerAdministrator, args[2])

	secMetadata := SecurityMetaData{tSigma, tCertificate, tBinding, invokerAdminSigma}
	secMetadataRaw, err := asn1.Marshal(secMetadata)
	if err != nil {
		return nil, err
	}
	return secMetadataRaw, nil
}

func IsAuthorisedToQuery(stub *shim.ChaincodeStub, args []string) (bool, error) {
	entry := &EntrySecurityInfo{}
	err := json.Unmarshal([]byte(args[1]), entry)
	if err != nil {
		return false, errors.New("Failed to unmarshal value into struct.")
	}

	secMetadataRaw := new(SecurityMetaData)
	_, err = asn1.Unmarshal(entry.Signature, secMetadataRaw)
	if err != nil {
		return false, errors.New("Failed unmarshalling metadata")
	}

	//Retrieve the invoker initiating the transaction
	invoker,error := getInvoker(stub)
	if error != nil {
		return false, errors.New("Error retrieving the Invoker")
	}

	//Get invoker's Enrollment Certificate
	eCert, error := getInvokerEcert(invoker)
	if error != nil {
		return false, errors.New("Error retrieving the Invoker Enrollment Certificate")
	}

	//Get the Invoker's affiliation
	affiliation, error := getInvokerAffiliation(eCert)
	if error != nil {
		return false, errors.New("Error retrieving the Invoker Enrollment Certificate")
	}

	//If the Invoker's affiliation Admin can match with the InvokerSigma, the Invoker is authorised to view the Transaction
	invokerAdministrator := getAdminUser(affiliation)
	nostroLogger.Debug("[Confidentiality] Params Collected for Invoker Signature Match.....")
	if(verifyInvokerSigma(invokerAdministrator,args[0],secMetadataRaw)){
		nostroLogger.Debug("[Confidentiality] Invoker Signature Matched !!!.....")
		return true,nil
	}
	nostroLogger.Debug("[Confidentiality] Invoker Signature Failed.....")

	//If the above fails, the system will attempt a transaction sigma match to verify invoker permissions

	recipientAdministrator := getAdminUser(affiliation)

	nostroLogger.Debug("[Confidentiality] Params Collected for Transaction Signature Match.....")
	if(verifyTransactionSigma(recipientAdministrator, args[0], secMetadataRaw)){
		nostroLogger.Debug("[Confidentiality] Transaction Signature Matched !!.....")
		return true,nil
	}
	nostroLogger.Debug("[Confidentiality] Transaction Signature Failed.....")
	return false,nil
}

func isRegisteredInvoker(stub *shim.ChaincodeStub) (bool){
	nostroLogger.Debug("[Confidentiality] Attempting to init the invoker....")
	invokerName,error := getInvoker(stub)

	nostroLogger.Debug("[Confidentiality] Invoker Name : "+invokerName)

	if error != nil {
		nostroLogger.Debug("[Confidentiality] Invoker not registered....")
		return false
	}
	_, err := crypto.InitClient(invokerName, nil)
	if err != nil {
		nostroLogger.Debug("[Confidentiality] Invoker not registered....")
		return false
	}
	nostroLogger.Debug("[Confidentiality] Invoker is registered....")
	return true
}

func getInvoker(stub *shim.ChaincodeStub) (string, error) {

	bytes, err := stub.GetCallerCertificate()
	if err != nil {
		nostroLogger.Debug("[Confidentiality] Error Reading Certificate : "+err.Error())
		return "",err
	}

	enrolment, err := attr.GetValueFrom("enrolment", bytes)
	if err != nil {
		nostroLogger.Debug("[Confidentiality] Error Reading Attribute : "+err.Error())
		return "",err
	}

	return string(enrolment), nil
}

func getInvokerEcert(invokerName string) ([]byte, error) {

	invoker, err := crypto.InitClient(invokerName, nil)
	if err != nil {
		return nil, err
	}
	eCertHandler, err := invoker.GetEnrollmentCertificateHandler()
	if err != nil {
		return nil, err
	}
	eCert := eCertHandler.GetCertificate()
	return eCert,nil
}


func getInvokerAffiliation(eCert []byte) (string, error) {

	x509Cert, err := x509.ParseCertificate(eCert);
	if err != nil {
		return "", errors.New("Couldn't parse certificate")
	}
	cn := x509Cert.Subject.CommonName
	res := strings.Split(cn,"\\")
	return res[1], nil
}



func getAdminUser (financialInstitution string) (crypto.Client){
	var (
		administrator  crypto.Client
		err error
	)
	switch strings.ToLower(financialInstitution) {
	case "anz":
		administrator, err = crypto.InitClient("vp0Admin", nil)
		if err != nil {
			return nil
		}
	case "wf":
		administrator, err = crypto.InitClient("vp1Admin", nil)
		if err != nil {
			return nil
		}
	case "ba":
		administrator, err = crypto.InitClient("vp2Admin", nil)
		if err != nil {
			return nil
		}
	case "lb":
		administrator, err = crypto.InitClient("vp3Admin", nil)
		if err != nil {
			return nil
		}
	}
	nostroLogger.Debug("[Confidentiality] Admin Returned : "+administrator.GetName())
	return administrator
}

func generateInvokerSigma(invokerAdmin crypto.Client, transactionKey string) ([]byte, error) {
	ecertHandler,err := invokerAdmin.GetEnrollmentCertificateHandler()
	if err != nil {
		return nil, err
	}
	invokerSigma, err := ecertHandler.Sign([]byte(transactionKey))
	if err != nil {
		return nil, err
	}
	return invokerSigma,nil
}

func verifyInvokerSigma(invokerAdmin crypto.Client, transactionKey string, secMetaData *SecurityMetaData) (bool) {
	ecertHandler,err := invokerAdmin.GetEnrollmentCertificateHandler()
	if err != nil {
		return false
	}
	errV := ecertHandler.Verify(secMetaData.InkSigma,[]byte(transactionKey))
	if(errV == nil){
		return true
	}
	return false
}

func generateTransactionSigma (recipientAdmin crypto.Client, transactionKey string)([]byte, []byte, []byte, error){
	tCertHandler, err := recipientAdmin.GetTCertificateHandlerNext()
	if err != nil {
		return nil, nil, nil, err
	}
	txnHandler, err := tCertHandler.GetTransactionHandler()
	if err != nil {
		return nil, nil, nil, err
	}
	binding, err := txnHandler.GetBinding()

	if err != nil {
		return nil, nil, nil, err
	}

	if err != nil {
		return nil, nil, nil, err
	}
	tCert := tCertHandler.GetCertificate()
	eCertHandler, err := recipientAdmin.GetEnrollmentCertificateHandler()
	if err != nil {
		return nil, nil, nil, err
	}
	eCert := eCertHandler.GetCertificate()

	sigma, err := tCertHandler.Sign(append(tCert, append(append(eCert, []byte(transactionKey)...), binding...)...))

	if err != nil {
		return nil, nil, nil, err
	}
	return sigma,tCert, binding, nil
}

func verifyTransactionSigma(recipientAdmin crypto.Client, transactionKey string, secMetaData *SecurityMetaData) (bool) {

	passedTxnSigma := secMetaData.TxnSigma
	passedTxnCert := secMetaData.TxnCert
	passedTxnBinding := secMetaData.TxnBinding

	txnHandler,err := recipientAdmin.GetTCertificateHandlerFromDER(passedTxnCert)
	if err != nil {
		return false
	}

	eCertHandler, err := recipientAdmin.GetEnrollmentCertificateHandler()
	if err != nil {
		return false
	}
	eCert := eCertHandler.GetCertificate()

	errV:=txnHandler.Verify(passedTxnSigma, append(passedTxnCert, append(append(eCert, []byte(transactionKey)...), passedTxnBinding...)...))

	if(errV == nil){
		return true
	}
	nostroLogger.Debug("[Confidentiality] Error Message: "+errV.Error())
	return false
}