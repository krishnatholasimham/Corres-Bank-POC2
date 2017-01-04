// Local envs
var ccURL = "http://0.0.0.0:5000/chaincode";
var ccName = "anz";

remoteEnv = false;
if(remoteEnv) {
  ccURL   = "https://636189dc-fdf0-4066-b45e-ac10c59ff84d_vp1-api.zone.blockchain.ibm.com:443/chaincode";  // dev-1
  ccURL   = "https://636189dc-fdf0-4066-b45e-ac10c59ff84d_vp3-api.zone.blockchain.ibm.com:443/chaincode";
  ccName1 = "b402ec094f7f079f324657da6765321fb4838255d0553d93e562c3ad321a9f2e085e3a16a5524566fe7168b348c206c9dc5755b730481342fd1ff653f744c0f7";
  ccName2 = "1ea3ee57400343cdce4665e457e06773a59d4a4a07b360d6b41e5b4fdbdcc8b52745bcccef2b98ddc7bcbd1726f23335bfea6de81e5d3f2b25c34d101108ccc7";
  ccName3 = "526e6e7d7677ef39b777454d835b5ddee9e68b2268a3d83bb8ca0bac4e734a78e300024bb63305a6dde169ea5108c023bb6519860584002d70dc74742ced63fb";
  ccName  = ccName3;
}
var secureUserWF = "user_type1_ebafb7be7a";
var secureUserANZ = "user_type1_ebafb7be7a";
var secureUserBOA = "BA";

var BANK_NAME_WF = "WF";
var BANK_NAME_ANZ = "ANZ";
var BANK_NAME_BA = "BA";

secureLocalEnv = false;  // Should be used together with membersvcs.secure-sample.yaml - can read this flag from config later.
if(secureLocalEnv) {
  secureUserWF = "wf1";
  secureUserANZ = "anz1";
  secureUserBOA = "boa1";
}

function getTxnTypes() {
      return [
        {
          "name": "Funding Message",
          "value":"funding_message"
        },
        {
          "name":"Payment Instruction",
          "value":"payment_instruction"
        },
        {
          "name":"Confirm Payment Instruction",
          "value":"confirm_payment_instruction"
        },
      ];
    };

    function getPaticipatingBanks() {
      return [
        {
          "name": "ANZ",
          "value":"ANZ"
        },
        {
          "name":"Wells Fargo",
          "value":"WF"
        },
      ];
    };

    function getCurrencies() {
          return [
            {
              "name": "Australian Dollar",
              "value":"AUD"
            },
            {
              "name":"US Dollar",
              "value":"USD"
            },
          ];
        };

    function getConfirmationMethod() {
      return [
        {
          "name": "via File Name",
          "value":"filename"
        },
        {
          "name":"via Request ID",
          "value":"requestID"
        },
      ];
    };

    function getFeeTypes() {
      return [
        {
          "name": "BEN - Beneficiary bears the fee amount",
          "value":"BEN"
        },
        {
          "name":"OUR - Payer bears the fee amount",
          "value":"OUR"
        },
        {
          "name":"SHA - Fee amount is shared between the two parties",
          "value":"SHA"
        },
      ];
    };

    function getMsgTypes() {
      return [
        {
          "name": "MT103",
          "value":"103"
        },
        {
          "name":"MT202",
          "value":"202"
        },
      ];
    };

    function ISODateString(d){
        function pad(n){return n<10 ? '0'+n : n}
        return d.getUTCFullYear()+'-'
            + pad(d.getUTCMonth()+1)+'-'
            + pad(d.getUTCDate())+'T'
            + pad(d.getUTCHours())+':'
            + pad(d.getUTCMinutes())+':'
            + pad(d.getUTCSeconds())+'Z'
    }

    function deleteTable(tableName, complete){
      if(complete){
        var table = document.getElementById(tableName);
        table.style.display = "none";
        var messageName = tableName+"Message";
        var message = document.getElementById(messageName);
        message.style.display = "inline";
      }
      else{
        var messageName = tableName+"Message";
        var message = document.getElementById(messageName);
        message.style.display = "none";

        var tableHeaderRowCount = 1;
          var table = document.getElementById(tableName);
          table.style.display = "table";
          var rowCount = table.rows.length;
          for (var i = tableHeaderRowCount; i < rowCount; i++) {
                table.deleteRow(tableHeaderRowCount);
          }
      }
    }

    function generateChaincodeMessage(vmModel){
        var chainCodeMessage="";
        switch (vmModel.txnType){
          case 'funding_message':
          var time = new Date();
          time = ISODateString(time);
          var fundingDate = vmModel.fundingDate;
          if(fundingDate){
            var fundDate = new Date();
            var date = Date.parse(fundDate);
            var newDate = new Date(date);
            fundingDate = ISODateString(new Date(newDate));
          }else{
              fundingDate = time;
          }
            chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "invoke",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"addLedgerEntryFunding",
                              "args": [""+vmModel.accountOwner+"", ""+vmModel.accountHolder+"", ""+vmModel.fundingAmount+"", fundingDate]
                              },
                        "secureContext":localStorage.user,
                        "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            break;

            case 'payment_instruction':
            var time = new Date();
            var upDateTime;
            if(vmModel.updateTime){
                upDateTime = vmModel.updateTime;
            }else{
                upDateTime = ISODateString(time);
            }
            var valueDate = new Date();
            var date = Date.parse(vmModel.valueDate);
            var newDate = new Date(date);
            valueDate = ISODateString(new Date(newDate));
            var paymentAmount = vmModel.paymentAmount;

            chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "invoke",
                        "params": {
                              "type": 1,
                              "chaincodeID": {
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"addPaymentInstruction",
                              "args":[""+paymentAmount+"", ""+vmModel.payerBank+"", ""+vmModel.beneBank+"", ""+upDateTime+"",""+vmModel.feeType+"",""+valueDate+"",""+vmModel.senderReference+"",""+vmModel.currency+"",""+ISODateString(time)+"",""+vmModel.msgNum+"",""+vmModel.msgType+"",""+vmModel.BookInt+"",""+vmModel.F33BCCY+"",""+vmModel.F33BAMT+""]
                              },
                        "secureContext":localStorage.user,
                        "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            break;

            case 'confirm_payment_instruction':
            var time = new Date();
            time = ISODateString(time);

            chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "invoke",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"addPaymentConfirmation",
                              "args":[""+vmModel.txnKey+"",""+time+"", ""+vmModel.beneBank+"",""+time+"",""+vmModel.msgNum+"",""+vmModel.msgType+"","1234567"]
                              },
                        "secureContext":localStorage.user,
                        "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            break;
            case 'reject_payment_instruction':
            var time = new Date();
            time = ISODateString(time);
            var rationale = "Payment not honored";

            chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "invoke",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"rejectPaymentInstruction",
                              "args":[""+vmModel.txnKey+"",""+rationale+"",""+time+""]
                              },
                        "secureContext":localStorage.user,
                        "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            break;

            case 'match_confirm_payment_instruction':
            var time = new Date();
            time = ISODateString(time);
            var valueDate = new Date();
            var date = Date.parse(vmModel.valueDate);
            var newDate = new Date(date);
            valueDate = ISODateString(new Date(newDate));
            var bankName = localStorage.bankName;
            var paymentAmount = vmModel.paymentAmount;

            chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "invoke",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"matchUnconfirmedTransactions",
                              "args":[""+paymentAmount+"", ""+vmModel.payerBank+"", ""+vmModel.beneBank+"", ""+vmModel.updateTime+"",""+vmModel.feeType+"",""+valueDate+"",""+vmModel.senderReference+"",""+vmModel.currency+"",""+time+"",""+vmModel.msgNum+"",""+vmModel.msgType+"",""+vmModel.BookInt+"",""+vmModel.F33BCCY+"",""+vmModel.F33BAMT+""]
                              },
                          "secureContext":localStorage.user,
                          "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            break;
        }
        return chainCodeMessage;
    }