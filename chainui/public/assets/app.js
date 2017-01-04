
console.clear(); // <-- keep the console clean on refresh

/*var intervalID = setInterval(function(){
  getStatementBalance("ANZ","WF");
  getStatementBalance("WF","ANZ");
  getReceivedUnconfirmedPayments("ANZ");
  getReceivedUnconfirmedPayments("WF");
}, 
5000);*/

/* global angular */
(function() {
  'use strict';

  var app = angular.module('formlyExample', ['formly', 'formlyBootstrap', '720kb.datepicker'], function config(formlyConfigProvider) {
    // set templates here
    formlyConfigProvider.setType({
      name: 'custom',
      templateUrl: 'custom.html'
    });
  });
  

  app.controller('MainCtrl', function MainCtrl(formlyVersion) {
    var vm = this;
    // funcation assignment
    vm.onSubmit = onSubmit;
    
    vm.fields = [
      {
      key: "txnType",
      type: "select",
      templateOptions: {
          label: "Select Transaction Type",
          options: getTxnTypes()
        },
        watcher: {
          listener: function(field, newValue, oldValue, formScope, stopWatching) {
            if (newValue) {
              vm.model = {"txnType": newValue};
            }
          }
        }
      },
      //Funding Message Parameters
      {
      key: "accountOwner",
      type: "select",
      templateOptions: {
          label: "Select Account Owner",
          options: getPaticipatingBanks()
        },
        hideExpression: 'model.txnType != "funding_message"'
      },
      {
      key: "accountHolder",
      type: "select",
      templateOptions: {
          label: "Select Account Holder",
          options: getPaticipatingBanks()
        },
        hideExpression: 'model.txnType != "funding_message"'
      },
      {
        key: 'fundingAmount',
        type: 'input',
        templateOptions: {
          label: 'Amount',
          placeholder: '40000'
        },
        hideExpression: 'model.txnType != "funding_message"'
      },
      {
          key: 'fundingDate',
          type: 'input',
          templateOptions: {
            label: 'Funding Date (DD-MMM-YYYY)',
            placeholder: '27-AUG-2016'
          },
          hideExpression: 'model.txnType != "funding_message"'
        },
    

    //Payment Instruction Message Parameters
    {
      key: 'senderReference',
      type: 'input',
      templateOptions: {
          label: 'Sender\'s Reference',
      },
      hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: "payerBank",
      type: "select",
      templateOptions: {
          label: "Select Payer Bank",
          options: getPaticipatingBanks()
        },
        hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: "beneBank",
      type: "select",
      templateOptions: {
          label: "Select Beneficiary Bank",
          options: getPaticipatingBanks()
        },
        hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: "currency",
      type: "select",
      templateOptions: {
          label: "Select Currency",
          options: getCurrencies()
        },
        hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: 'valueDate',
      type: 'input',
      templateOptions: {
          label: 'Value Date (DD-MMM-YYYY)',
      },
      hideExpression: 'model.txnType != "payment_instruction"'
    },

    {
        key: 'paymentAmount',
        type: 'input',
        templateOptions: {
          label: 'Amount',
          placeholder: '1000'
        },
        hideExpression: 'model.txnType != "payment_instruction"'
      },
    {
      key: "feeType",
      type: "select",
      templateOptions: {
          label: "Select Fee Type",
          options: getFeeTypes()
        },
        hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: 'msgNum',
      type: 'input',
      templateOptions: {
          label: 'MSG#'
      },
      hideExpression: 'model.txnType != "payment_instruction"'
    },
    {
      key: 'msgType',
      type: 'select',
      templateOptions: {
          label: 'MSG Type',
          options: getMsgTypes()
      },
      hideExpression: 'model.txnType != "payment_instruction"'
    },


    //Payment Confirmation Message Parameters
    {
      key: 'txnKey',
      type: 'input',
      templateOptions: {
          label: 'Transaction Key',
      },
      hideExpression: 'model.txnType != "confirm_payment_instruction"'
    },
    {
      key: "beneBank",
      type: "select",
      templateOptions: {
          label: "Select Beneficiary Bank",
          options: getPaticipatingBanks()
        },
        hideExpression: 'model.txnType != "confirm_payment_instruction"'
    },
    {
      key: 'msgNum',
      type: 'input',
      templateOptions: {
          label: 'MSG#'
      },
      hideExpression: 'model.txnType != "confirm_payment_instruction"'
    },
    {
      key: 'msgType',
      type: 'select',
      templateOptions: {
          label: 'MSG Type',
          options: getMsgTypes()
      },
      hideExpression: 'model.txnType != "confirm_payment_instruction"'
    }
];

    // function definition
    function onSubmit() {

      var vmModel = JSON.parse(JSON.stringify(vm.model), null, 2);
      var chainCodeMessage = generateChaincodeMessage(vmModel);

      callChaincode(JSON.stringify(chainCodeMessage),'sendMessage');

    }
    });

  app.directive('exampleDirective', function() {
    return {
      templateUrl: 'example-directive.html'
    };
  });
})();