var graphDataANZ = [];
var graphDataWF = [];
var graphDataStmtBalANZ = [];
var graphDataStmtBalWF = [];
var graphDataStmtBalWithoutBCANZ = [];
var graphDataStmtBalWithoutBCWF = [];
function init(){
      var message = document.getElementById("transactionSent");
      if(message)
        message.style.display = "none";
      installChaincode();
}
var pollInterval = 5000;
$(document).ready(function() {
    $(".various").fancybox({
        maxWidth  : 900,
        maxHeight : 300,
        fitToView : true,
        width   : '70%',
        height    : '70%',
        autoSize  : true,
        closeClick  : true,
        openEffect  : 'none',
        closeEffect : 'none'
    });

    var statementBalanceForWF;
    var statementBalanceForANZ;

    $("#confirmedPayments-tab").click(function () {
        clearInterval(unconfirmedForANZ);
        clearInterval(unconfirmedForWF);
        clearInterval(confirmForWF);
        clearInterval(confirmForANZ);
        institution = localStorage.bankName;
        if("WF" == institution){
            getStatementBalance("WF","ANZ");
        }else if("ANZ" == institution){
            getStatementBalance("ANZ","WF");
        }
        if("WF" == institution){
            clearInterval(statementBalanceForANZ);
            statementBalanceForWF = setInterval(function(){
                getStatementBalance("WF","ANZ");
            }, pollInterval);
        }else if("ANZ" == institution){
            clearInterval(statementBalanceForWF);
            statementBalanceForANZ = setInterval(function(){
                getStatementBalance("ANZ","WF");
            }, pollInterval);
        }
    });

    var unconfirmedForWF;
    var unconfirmedForANZ;
    $("#unconfirmedPayments-tab").click(function () {
        clearInterval(statementBalanceForANZ);
        clearInterval(statementBalanceForWF);
        clearInterval(confirmForWF);
        clearInterval(confirmForANZ);
        institution = localStorage.bankName;
        if("WF" == institution){
            getReceivedUnconfirmedOutgoingPayments("WF");
        }else if("ANZ" == institution){
            getReceivedUnconfirmedOutgoingPayments("ANZ");
        }
        if("WF" == institution){
            clearInterval(unconfirmedForANZ);
            unconfirmedForWF = setInterval(function(){
                getReceivedUnconfirmedOutgoingPayments("WF");
            }, pollInterval);
        }else if("ANZ" == institution){
            clearInterval(unconfirmedForWF);
            unconfirmedForANZ = setInterval(function(){
                getReceivedUnconfirmedOutgoingPayments("ANZ");
            }, pollInterval);
        }
    });

    var confirmForWF;
    var confirmForANZ;
    $("#confirmPayments-tab").click(function () {
        clearInterval(statementBalanceForANZ);
        clearInterval(statementBalanceForWF);
        clearInterval(unconfirmedForANZ);
        clearInterval(unconfirmedForWF);
        institution = localStorage.bankName;
        if("WF" == institution){
            getReceivedUnconfirmedPayments("WF");
            getReceivedConfirmedStatementBalance("ANZ","WF");
        }else if("ANZ" == institution){
            getReceivedUnconfirmedPayments("ANZ");
            getReceivedConfirmedStatementBalance("WF","ANZ");
        }
        if("WF" == institution){
            clearInterval(confirmForANZ);
            confirmForWF = setInterval(function(){
                getReceivedUnconfirmedPayments("WF");
                getReceivedConfirmedStatementBalance("ANZ","WF");
            }, pollInterval);
        }else if("ANZ" == institution){
            clearInterval(confirmForWF);
            confirmForANZ = setInterval(function(){
                getReceivedUnconfirmedPayments("ANZ");
                getReceivedConfirmedStatementBalance("WF","ANZ");
            }, pollInterval);
        }
    });
    $("#dashboard-tab").click(function () {
        clearInterval(statementBalanceForANZ);
        clearInterval(statementBalanceForWF);
        clearInterval(unconfirmedForANZ);
        clearInterval(unconfirmedForWF);
        clearInterval(confirmForWF);
        clearInterval(confirmForANZ);
        loadDashboard();
    });
    $("#thirdpartydashboard-tab").click(function () {
        loadThirdPartyDashboard();
    });

    $("#createPayments-tab, #fundingops-tab, #search-tab").click(function () {
        clearInterval(statementBalanceForANZ);
        clearInterval(statementBalanceForWF);
        clearInterval(unconfirmedForANZ);
        clearInterval(unconfirmedForWF);
        clearInterval(confirmForWF);
        clearInterval(confirmForANZ);
    });
});


function callChaincode(chainCodeMessage, task, institution, callbackFn) {
          var xhr = new XMLHttpRequest();
          xhr.onreadystatechange = function() {

            //Process Statment Account Transactions
            if (xhr.readyState == XMLHttpRequest.DONE) {
              if(task == 'getReceivedConfirmedStatementBalance'){
                var statementTableConfirmed = 'confirmedStatus'+institution;
                var chainCodeResponse = JSON.parse(xhr.responseText);
                if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                  deleteTable(statementTableConfirmed,"complete");
                  return;
                }
                deleteTable(statementTableConfirmed);
                var msgString = chainCodeResponse.result.message;
                msgString = JSON.parse(msgString);
                for (var key in msgString){
                  var msg = msgString[key];
                  populateTable(msg,statementTableConfirmed, task);
                }

              }
              if(task == 'getStatementBalance'){
                var statementTable = 'blockchainStatus'+institution;
                var statementTableConfirmed =  'confirmedStatus'+institution;
                var chainCodeResponse = JSON.parse(xhr.responseText);

                if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                  deleteTable(statementTable,"complete");
                  deleteTable(statementTableConfirmed,"complete");
                  return;
                }

                deleteTable(statementTable);
                deleteTable(statementTableConfirmed);
                var msgString = chainCodeResponse.result.message;
                msgString = JSON.parse(msgString);
                  for (var key in msgString){
                        var msg = msgString[key];
                        populateTable(msg,statementTable, task);
                        //getGrandTotals(statementTable);
                  }
                  for (var key in msgString){
                      var msg = msgString[key];
                      populateTable(msg,statementTableConfirmed, task);
                      //getGrandTotals(statementTable);
                }
              }
                if(task == 'getStatementBalanceForGraph'){
                    var chainCodeResponse = JSON.parse(xhr.responseText);
                    switch (institution){
                        case 'ANZ':
                            graphDataStmtBalANZ = [];
                            break;
                        case 'WF':
                            graphDataStmtBalWF = [];
                      }
                    if(chainCodeResponse.result && chainCodeResponse.result.message && chainCodeResponse.result.message != "null") {
                        var msgString = chainCodeResponse.result.message;
                        msgString = JSON.parse(msgString);
                        msgString.sort(function (x, y) {
                            if(new Date(x.TimestampCreated) > new Date(y.TimestampCreated))
                                return 1;
                            else if((new Date(x.TimestampCreated) == new Date(y.TimestampCreated)) && (x.StatementAccountBalance > y.StatementAccountBalance))
                                return 1;
                            return -1;
                        });

                        for (var key in msgString) {
                            var msg = msgString[key];
                            if(msg.Type == 'CONFIRMATION-RECORD' || msg.Type =='FUNDING-RECORD' || msg.Type =='DIRECT-CREDIT-RECORD'){
                                setGraphDataStmtBalance(msg, institution);
                            }
                        }
                    }
                    if (typeof(callbackFn) == 'function') {
                        callbackFn();
                    }
                }
              // getUnconfirmed Balance History for plotting the Graph
              if(task == 'getUnconfirmedBalanceHistory'){
                var chainCodeResponse = JSON.parse(xhr.responseText);
                switch (institution){
                    case 'ANZ':
                        graphDataANZ = [];
                        break;
                    case 'WF':
                        graphDataWF = [];
                  }
                  if(chainCodeResponse.result && chainCodeResponse.result.message && chainCodeResponse.result.message != "null") {
                      var msgString = chainCodeResponse.result.message;
                      msgString = JSON.parse(msgString);
                      msgString.sort(function (x, y) {
                          if(new Date(x.TimestampCreated) > new Date(y.TimestampCreated))
                              return 1;
                          else if((new Date(x.TimestampCreated) == new Date(y.TimestampCreated)) && (x.IndicativeBalance > y.IndicativeBalance))
                              return 1;
                          return -1;
                      });

                      for (var key in msgString) {
                          var msg = msgString[key];
                          if(msg.Type == 'REQUEST-RECORD'  || msg.Type =='FUNDING-RECORD' || msg.Type =='DIRECT-CREDIT-RECORD'){
                              setGraphData(msg, institution);
                          }
                      }
                  }
                  if (typeof(callbackFn) == 'function') {
                      callbackFn();
                  }
              }

              //Process Inward Payment Instructions
              if(task == 'getReceivedUnconfirmedPayments'){
                var receivedUnconfirmedPaymentsTable = 'unconfirmedReceivedPayments'+institution;
                var chainCodeResponse = JSON.parse(xhr.responseText);

                if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                  deleteTable(receivedUnconfirmedPaymentsTable,"complete");
                  return;
                }

                deleteTable(receivedUnconfirmedPaymentsTable);
                var msgString = chainCodeResponse.result.message;

                msgString = JSON.parse(msgString);
                msgString.sort(function(a, b){
                    return new Date(a.TimestampCreated).getTime() - new Date(b.TimestampCreated).getTime();
                });
                for (var key in msgString){
                      var msg = msgString[key];
                      populateTable(msg,receivedUnconfirmedPaymentsTable, task);
                    }

              }//Process Inward Payment Instructions
             if(task == 'getSentUnconfirmedPayments'){
               var receivedUnconfirmedPaymentsTable = 'unconfirmedOutgoingReceivedPayments'+institution;
               var chainCodeResponse = JSON.parse(xhr.responseText);

               if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                 deleteTable(receivedUnconfirmedPaymentsTable,"complete");
                 return;
               }
               deleteTable(receivedUnconfirmedPaymentsTable);


               var msgString = chainCodeResponse.result.message;
               msgString = JSON.parse(msgString);
               msgString.sort(function(a, b){
                   return new Date(a.TimestampCreated).getTime() - new Date(b.TimestampCreated).getTime();
               });

               for (var key in msgString){
                   var msg = msgString[key];
                   populateTable(msg,receivedUnconfirmedPaymentsTable, task);
               }
             }
              //Process Rejected Payment Instructions
              if(task == 'getRejectedPaymentInstructions'){
                var rejectedPaymentsTable = 'rejectedPayments'+institution;
                var chainCodeResponse = JSON.parse(xhr.responseText);

                if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                  deleteTable(rejectedPaymentsTable,"complete");
                  return;
                }

                deleteTable(rejectedPaymentsTable);
                var msgString = chainCodeResponse.result.message;
                msgString = JSON.parse(msgString);
                for (var key in msgString){
                      var msg = msgString[key];
                      populateTable(msg,rejectedPaymentsTable, task);
                    }
              }

              if(task == 'searchTransactions' || task == 'searchTransactionsThirdParty'){
                var searchTransactionsTable = 'searchTransactions';
                var chainCodeResponse = JSON.parse(xhr.responseText);

                if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
                  deleteTable(searchTransactionsTable,"complete");
                  return;
                }

                deleteTable(searchTransactionsTable);
                var msgString = chainCodeResponse.result.message;
                msgString = JSON.parse(msgString);
                for (var key in msgString){
                  var msg = msgString[key];
                  populateTable(msg,searchTransactionsTable, task);
                }
              }
                if(task == 'getTransactionDetails'){
                    var chainCodeResponse = JSON.parse(xhr.responseText);
                    callbackFn(chainCodeResponse);
                }
              //Process Create Transaction Requests
              if(task == 'sendMessage'){

                var message = document.getElementById("transactionSent");
                message.style.display = "table";

                var intervalID = setInterval(function(){
                    var message = document.getElementById("transactionSent");
                    message.style.display = "none";
                  },
                  pollInterval);
              }

              if(task == 'install'){
                initApp();

            }
          }
        }
      xhr.open('POST', ccURL, true);
      xhr.send(chainCodeMessage);
    }

    function installChaincode(){
      var chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "deploy",
                        "params": {
                          "type": 1,
                          "chaincodeID": {
                            "name": ccName
                          },
                          "ctorMsg": {
                            "function": "init",
                            "args": ["a","100"]
                          },
                          "secureContext": localStorage.user
                        },
                        "id": 2
                      };
      callChaincode(JSON.stringify(chainCodeMessage),'install');

    }

    function initApp(){
      var intervalID = setInterval(function(){
          //getStatementBalance("ANZ","WF");
          //getStatementBalance("WF","ANZ");
          //getReceivedUnconfirmedPayments("ANZ");
          //getReceivedUnconfirmedPayments("WF");

          getRejectedPaymentInstructions("ANZ");
          getRejectedPaymentInstructions("WF");

          //getReceivedUnconfirmedOutgoingPayments("ANZ");
          //getReceivedUnconfirmedOutgoingPayments("WF");
      },
      pollInterval);
    }

    function loadDashboard(){

        institution = localStorage.bankName
        if("WF" == institution){
            /*$.when(getStatementBalanceForGraph('WF', 'ANZ'), getUnconfirmedBalanceHistory("WF","ANZ", true)).done(function(a1, a2){
                // the code here will be executed when all four ajax requests resolve.
                // a1, a2, a3 and a4 are lists of length 3 containing the response text,
                // status, and jqXHR object for each of the four ajax calls respectively.
                populateCharts();
            });*/
            getStatementBalanceForGraph('WF', 'ANZ', function(){
                return getUnconfirmedBalanceHistory('WF', 'ANZ', function() {
                    return populateCharts();
                });
            });

        }else if("ANZ" == institution){

            getStatementBalanceForGraph('ANZ', 'WF', function(){
                return getUnconfirmedBalanceHistory('ANZ', 'WF', function() {
                    return populateCharts();
                });
            });
        }
    }

    function loadThirdPartyDashboard(){
        searchTransactionsThirdParty();
    }

    function getStatementBalance(accountOwner,accountHolder){
      var chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "query",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"getBalanceHistory",
                              "args":[accountOwner, accountHolder]
                              },
                              "secureContext": localStorage.user,
                              "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
        callChaincode(JSON.stringify(chainCodeMessage),'getStatementBalance', accountOwner);
    }

    function getStatementBalanceForGraph(accountOwner, accountHolder, callback) {
        var chainCodeMessage = {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID":{
                    "name":ccName
                },
                "ctorMsg": {
                    "function":"getBalanceHistory",
                    "args":[accountOwner, accountHolder]
                },
               "secureContext": localStorage.user,
               "attributes": ["enrolment"]
            },
            "id": Math.floor((Math.random() * 10) + 1)
        };
        callChaincode(JSON.stringify(chainCodeMessage),'getStatementBalanceForGraph', accountOwner, callback);
    }


    function getUnconfirmedBalanceHistory(accountOwner, accountHolder, callback){
        var chainCodeMessage = {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID":{
                    "name":ccName
                },
                "ctorMsg": {
                    "function":"getUnconfirmedBalanceHistory",
                    "args":[accountOwner, accountHolder]
                },
               "secureContext": localStorage.user,
               "attributes": ["enrolment"]
            },
            "id": Math.floor((Math.random() * 10) + 1)
        };
        callChaincode(JSON.stringify(chainCodeMessage),'getUnconfirmedBalanceHistory', accountOwner, callback);
    }
    function getReceivedConfirmedStatementBalance(accountOwner,accountHolder){
          var chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "query",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"getBalanceHistory",
                              "args":[accountOwner, accountHolder]
                              },
                             "secureContext": localStorage.user,
                             "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
        callChaincode(JSON.stringify(chainCodeMessage),'getReceivedConfirmedStatementBalance', accountOwner);
     }

    function getReceivedUnconfirmedPayments(accountOwner){
      var chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "query",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"getReceivedUnconfirmedPayments",
                              "args":[accountOwner]
                              },
                             "secureContext": localStorage.user,
                             "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };

      callChaincode(JSON.stringify(chainCodeMessage),'getReceivedUnconfirmedPayments',accountOwner);

    }

    function getReceivedUnconfirmedOutgoingPayments(accountOwner){
          var chainCodeMessage = {
                            "jsonrpc": "2.0",
                            "method": "query",
                            "params": {
                                  "type": 1,
                                  "chaincodeID":{
                                      "name":ccName
                                  },
                            "ctorMsg": {
                                  "function":"getSentUnconfirmedPayments",
                                  "args":[accountOwner]
                                  },
                             "secureContext": localStorage.user,
                             "attributes": ["enrolment"]
                            },
                            "id": Math.floor((Math.random() * 10) + 1)
                };

          callChaincode(JSON.stringify(chainCodeMessage),'getSentUnconfirmedPayments',accountOwner);

        }

    function getRejectedPaymentInstructions(accountOwner){
      var chainCodeMessage = {
                        "jsonrpc": "2.0",
                        "method": "query",
                        "params": {
                              "type": 1,
                              "chaincodeID":{
                                  "name":ccName
                              },
                        "ctorMsg": {
                              "function":"getRejectedPaymentInstructions",
                              "args":[accountOwner]
                              },
                         "secureContext": localStorage.user,
                         "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };

      callChaincode(JSON.stringify(chainCodeMessage),'getRejectedPaymentInstructions',accountOwner);

    }

    function confirmPaymentInstruction(txnKey){
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
                              "args":[""+txnKey+"",""+time+"", "WF", "2006-01-02T15:04:05.999999999Z07:00", "1", "1", "1"]
                              },
                         "secureContext": localStorage.user,
                         "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            callChaincode(JSON.stringify(chainCodeMessage),'sendMessage');
    }

    function rejectPaymentInstruction(txnKey){
      var time = new Date();
      var rationale = "Payment not honored";
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
                              "function":"rejectPaymentInstruction",
                              "args":[""+txnKey+"",""+rationale+"",""+time+""]
                              },
                         "secureContext": localStorage.user,
                         "attributes": ["enrolment"]
                        },
                        "id": Math.floor((Math.random() * 10) + 1)
            };
            callChaincode(JSON.stringify(chainCodeMessage),'sendMessage');
    }
    function searchTransactions(){
      var chainCodeMessage = {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                  "type": 1,
                  "chaincodeID":{
                      "name":ccName
                  },
            "ctorMsg": {
                  "function":"getAll",
                  "args":["true", "true", "true", "true","true","true"]
                  },
             "secureContext": localStorage.user,
             "attributes": ["enrolment"]
            },
            "id": Math.floor((Math.random() * 10) + 1)
      };

      callChaincode(JSON.stringify(chainCodeMessage),'searchTransactions', 'accountOwner');
    }
    function searchTransactionsThirdParty(){
          var chainCodeMessage = {
                "jsonrpc": "2.0",
                "method": "query",
                "params": {
                      "type": 1,
                      "chaincodeID":{
                          "name":ccName
                      },
                "ctorMsg": {
                      "function":"keys",
                      "args":[]
                      },
                 "secureContext": localStorage.user,
                 "attributes": ["enrolment"]
                },
                "id": Math.floor((Math.random() * 10) + 1)
          };

      callChaincode(JSON.stringify(chainCodeMessage),'searchTransactionsThirdParty', 'accountOwner');
    }
    function getTransactionDetails(key) {
        var chainCodeMessage = {
            "jsonrpc": "2.0",
            "method": "query",
            "params": {
                "type": 1,
                "chaincodeID":{
                    "name":ccName
                },
                "ctorMsg": {
                    "function":"get",
                    "args":["" + key]
                },
               "secureContext": localStorage.user,
               "attributes": ["enrolment"]
            },
            "id": Math.floor((Math.random() * 10) + 1)
        };

        var populateTransDetails = function(chainCodeResponse){
            populateTransactionDetails(key, chainCodeResponse);
        }
        callChaincode(JSON.stringify(chainCodeMessage),'getTransactionDetails', 'accountOwner', populateTransDetails);
    }

    function populateTransactionDetails(key, chainCodeResponse){
        if(!chainCodeResponse.result || !chainCodeResponse.result.message || chainCodeResponse.result.message == "null"){
            populateTable(null, null, task)
            return;
        }
        var counter = 0;
        var msgString = chainCodeResponse.result.message;
        var canParse = true;
        try {
          msgString = JSON.parse(msgString);
        } catch(err) {
          // Allow garbled data to show as key.
          canParse = false;
        }

        var transDetailsTable = '#accordion_' + key;
        var records = '<table style="width:100%">';
        $('' + transDetailsTable).html('<table>');
        record = '';
        if(canParse) {
          for (var key in msgString){
              var msg = msgString[key];
              record += ('<td class="transDetailsElement"><b>' + key  + ' : </b>' + msg + '</td>');
              counter ++;
              if(counter % 2 == 0){
                  records += '<tr>' + record +'</tr>';
                  record = '';
              }
          }
        } else {
          records += ('<tr><td colspan="8" class="transDetailsElement">' + msgString + '</td></tr>');
        }
        records += '</table>';
        $('' + transDetailsTable).html(records);
    }
    function populateTable(msgString,tableName, task){

      switch (task){
        case 'getStatementBalance':
            var table = "#"+tableName;
            var refKey = "";
            if(msgString.RefKey) {
              refKey="Reference Key";
            }
           $(table).append(
                            "<tr>"
                              +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                              +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.Key+" onclick=invokeKeySearch('" + msgString.Key + "');>"+msgString.Key+"</a></td>"
                              +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.RefKey+">"+refKey+"</a></td>"
                              +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                              +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                              +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampConfirmed).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampConfirmed+"</td>"
                              +"<td style='padding: 10px; color: green; font-weight: bold;' class='balance'>"+msgString.StatementAccountBalance+"</td>"
                            +"</tr>" );

          break;
          case 'getReceivedConfirmedStatementBalance':
              var table = "#"+tableName;
              var refKey = "";
              if(msgString.RefKey) {
                refKey="Reference Key";
              }
             $(table).append(
                              "<tr>"
                                +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                                +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.Key+" onclick=invokeKeySearch('" + msgString.Key + "');>"+msgString.Key+"</a></td>"
                                +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.RefKey+">"+refKey+"</a></td>"
                                +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                                +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                                +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                                +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                                +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampConfirmed).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampConfirmed+"</td>"
                                +"<td style='padding: 10px; color: green; font-weight: bold;' class='balance'>"+msgString.StatementAccountBalance+"</td>"
                              +"</tr>" );

            break;

        case 'getReceivedUnconfirmedPayments':
          var table = "#"+tableName;
            $(table).append(
                            "<tr>"
                              +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                              +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.Key+" onclick=invokeKeySearch('" + msgString.Key + "');>"+msgString.Key+"</a></td>"
                              +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                              +"<td style='padding: 10px'>"+msgString.FeeType+"</td>"
                              +"<td style='padding: 10px'>"+msgString.SendersCharge+"</td>"
                              +"<td style='padding: 10px'>"+msgString.BenePays+"</td>"
                              +"<td style='padding: 10px;'>"+msgString.ValueDate+"</td>"
                              +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                              +"<td style='padding: 10px;'>"+"<button class='btn btn-success' onclick=confirmPaymentInstruction('"+msgString.Key+"')>Confirm</button> &nbsp; <button class='btn btn-danger' onclick=rejectPaymentInstruction('"+msgString.Key+"')>Reject</button>"+"</td>"
                            +"</tr>" );


          break;
          case 'getSentUnconfirmedPayments':
            var table = "#"+tableName;
              $(table).append(
                              "<tr>"
                                +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                                +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.Key+" onclick=invokeKeySearch('" + msgString.Key + "');>"+msgString.Key+"</a></td>"
                                +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                                +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                                +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                                +"<td style='padding: 10px'>"+msgString.FeeType+"</td>"
                                +"<td style='padding: 10px'>"+msgString.SendersCharge+"</td>"
                                +"<td style='padding: 10px'>"+msgString.BenePays+"</td>"
                                +"<td style='padding: 10px;'>"+msgString.ValueDate+"</td>"
                                +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                                +"<td style='padding: 10px;'>"+msgString.IndicativeBalance+"</td>"
                              +"</tr>" );


            break;
        case 'getRejectedPaymentInstructions':
          var table = "#"+tableName;
            $(table).append(
                            "<tr>"
                              +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                              +"<td style='padding: 10px'><a href='#'' data-toggle='tooltip' data-placement='top' title="+msgString.Key+" onclick=invokeKeySearch('" + msgString.Key + "');>"+msgString.Key+"</a></td>"
                              +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                              +"<td style='padding: 10px'>"+msgString.FeeType+"</td>"
                              +"<td style='padding: 10px'>"+msgString.SendersCharge+"</td>"
                              +"<td style='padding: 10px'>"+msgString.BenePays+"</td>"
                              +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                            +"</tr>" );


          break;
          case 'searchTransactionsThirdParty':
          case 'searchTransactions':
            var table = "#"+tableName;
            var refKey = "";
            if(msgString.RefKey) {
              refKey="Reference Key";
            }
            if(msgString.Type == 'REQUEST-RECORD'){
                msgString.TimestampConfirmed = null;
                msgString.StatementAccountBalance = 0;
            }
            if(msgString.Type == 'CONFIRMATION-RECORD'){
                msgString.IndicativeBalance = 0;
            }
            var trnKey = (task == 'searchTransactionsThirdParty' ? msgString : msgString.Key);
            if(trnKey != 'a')
              $(table).append(
                            "<tr>"
                              +"<td style='padding: 10px'>"+msgString.Type+"</td>"
                              +"<td style='padding: 10px'><a data-toggle='collapse' data-target='#accordion_" + trnKey +"'>"+ trnKey + "</a></td>"
                              +"<td style='padding: 10px'><a href='#' data-toggle='tooltip' data-placement='top' title="+msgString.RefKey+">"+refKey+"</a></td>"
                              +"<td style='padding: 10px'>"+msgString.OrderingInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.AccountWithInstitution+"</td>"
                              +"<td style='padding: 10px'>"+msgString.PaymentAmount+"</td>"
                              +"<td style='padding: 10px;display: none'>"+msgString.TimestampCreated+"</td>"
                              +"<td style='padding: 10px;display: none'>"+msgString.TimestampConfirmed+"</td>"
                              +"<td style='padding: 10px;'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                              +"<td style='padding: 10px;'>"+confirmedTimeStampFormatter(msgString.TimestampConfirmed, true)+"<br />"+confirmedTimeStampFormatter(msgString.TimestampConfirmed)+"</td>"
                              +"<td style='padding: 10px;display: none'>"+trnKey+"</td>"
                            +"</tr>"

                            +"<tr>"
                            +"<td style='display: none'>"+msgString.Type+"</td>"
                            +"<td style='display: none'><a data-toggle='collapse' data-target='#accordion_" + trnKey +"'>" + trnKey + "</a></td>"
                            +"<td style='display: none'><a href='#' data-toggle='tooltip' data-placement='top' title="+msgString.RefKey+">"+refKey+"</a></td>"
                            +"<td style='display: none'>"+msgString.OrderingInstitution+"</td>"
                            +"<td style='display: none'>"+msgString.AccountWithInstitution+"</td>"
                            +"<td style='display: none'>"+msgString.PaymentAmount+"</td>"
                            +"<td style='padding: 10px;display: none'>"+msgString.TimestampCreated+"</td>"
                            +"<td style='padding: 10px;display: none'>"+msgString.TimestampConfirmed+"</td>"
                            +"<td style='display: none'>Local:&nbsp;"+new Date(msgString.TimestampCreated).toLocaleString()+"<br />UTC:&nbsp;"+msgString.TimestampCreated+"</td>"
                            +"<td style='display: none'>"+confirmedTimeStampFormatter(msgString.TimestampConfirmed, true)+"<br />"+confirmedTimeStampFormatter(msgString.TimestampConfirmed)+"</td>"
                            +"<td style='padding: 10px;display: none'>"+trnKey+"</td>"
                            +"<td colspan='10'>"
                            +"<div id='accordion_" + trnKey + "' class='collapse' data-key='" + trnKey +"'>"
                            +"</div>"
                            +"</td>"
                            +"</tr>");

          autoSearchKey = $('#filters option:selected').val() === 'KEY'
          if(autoSearchKey){
            filterByKey();
          }
          $('#accordion_' + trnKey).on('show.bs.collapse', function (e) {
              var key = $('#accordion_' + trnKey).attr('data-key');
              getTransactionDetails(key);
          });

          break;
          case 'getTransactionDetails':
              var tableRow = "#"+tableName;
              $(tableRow).append(
                  "<p>"
                  +"First Key is:"
                  +"</p>" );
              break;

          default:
          break;

      }
    }

    function confirmedTimeStampFormatter(timestampConfirmed, toLocal){
        if(timestampConfirmed){
            if(toLocal){
                return "Local:&nbsp;"+ new Date(timestampConfirmed).toLocaleString();
            }else{
                return "UTC:&nbsp;"+timestampConfirmed;
            }
        }else{
            return "&nbsp;";
        }
    }

    function setGraphData(msgString, institution){
        switch (institution){
          case 'ANZ':
            graphDataANZ.push([msgString.IndicativeBalance,msgString.TimestampCreated]);
            break;
          case 'WF':
            graphDataWF.push([msgString.IndicativeBalance,msgString.TimestampCreated]);
            break;
        }
    }

    function setGraphDataStmtBalance(msgString, institution){
        switch (institution){
            case 'ANZ':
                graphDataStmtBalANZ.push([msgString.StatementAccountBalance,msgString.TimestampCreated]);
                graphDataStmtBalWithoutBCANZ.push([msgString.StatementAccountBalance,msgString.TimestampCreated,msgString.TransactionReferenceNumber, msgString.PaymentAmount]);
                break;
            case 'WF':
                graphDataStmtBalWF.push([msgString.StatementAccountBalance,msgString.TimestampCreated]);
                graphDataStmtBalWithoutBCWF.push([msgString.StatementAccountBalance,msgString.TimestampCreated,msgString.TransactionReferenceNumber, msgString.PaymentAmount]);
                break;
        }
    }

    function getGrandTotals(tableName){
        var table = "#"+tableName; var r=$(table).length;
        var data = []; var grandTotalValue = 0;
        var isNegative = false; var negColor = "green";
        $(table+' tr.grandTotalRow').remove();
        //gets all the values from the last columns of the given table adds to an array
        $(table+' tr').each(function(i)
        {
            var thisValue = $(this).find('td:last').html();
            if (typeof thisValue !== undefined){
                data.push(thisValue);
            }
        });
        //sum totals from values in the array
        for (var i = 0; i < data.length; i++) {
            grandTotalValue += data[i] << 0;
        }
        if (grandTotalValue < 0 ){ isNegative = true; negColor = "red" };
        //$(table+ " > tbody:last").eq(r-1).after("<tr>"
        $(table+ " tbody").append("<tr class='grandTotalRow'>"
                         +"<td style='padding: 10px'><b>Grand Total</b></td>"
                         +"<td style='padding: 10px'></td>"
                         +"<td style='padding: 10px'></td>"
                         +"<td style='padding: 10px'></td>"
                         +"<td style='padding: 10px'></td>"
                         +"<td style='padding: 10px'></td>"
                         +"<td style='padding: 10px;'></td>"
                         +"<td style='padding: 10px;'></td>"
                         +"<td style='padding: 10px; color: "+negColor+"; font-weight: bolder;' class='grandTotalValue'>"+grandTotalValue+"</td>"
                         +"</tr>");

    }

    function invokeKeySearch(key){
        $('#search-tab').trigger('click');
        $('#filters option[value=KEY]').attr('selected','selected');
        $('#SearchKeyFilter').removeClass("displayNone").addClass("displayInline");
        $('#Status, #StartAmountRange, #EndAmountRange, #MessageType,#StartDatePicker, #EndDatePicker').removeClass("displayInline").addClass("displayNone");
        $('#SearchKey').val(key)
        $("#searchTrn").trigger('click');
    }