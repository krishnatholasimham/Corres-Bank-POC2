function convertLocalDateToUTCDate(date, toUTC) {
    date = new Date(date);
    //Local time converted to UTC
    var localOffset = 300 * 60000;
    var localTime = date.getTime();
    if (toUTC) {
        date = localTime + localOffset;
    } else {
        date = localTime - localOffset;
    }
    date = new Date(date);
    return date;
}

function handleFile(e) {
  var files = e.target.files;
  var i,f;
  for (i = 0, f = files[i]; i != files.length; ++i) {
    var reader = new FileReader();
    var name = f.name;
    reader.onload = function(e) {
      var data = e.target.result;


      var workbook = XLSX.read(data, {type: 'binary'});
        var sheet_name_list = workbook.SheetNames;
        sheet_name_list.forEach(function(y) {
            var worksheet = workbook.Sheets[y];
            var headers = {};
            var data = []; //var newData = {};
            for(z in worksheet) {
                if(z[0] === '!') continue;
                //parse out the column, row, and value
                var col = z.substring(0,1);
                var row = parseInt(z.substring(1));
                var value = worksheet[z].v;

                //store header names
                if(row == 1) {
                    headers[col] = value;
                    continue;
                }

                if(!data[row]) data[row]={};
                data[row][headers[col]] = value;
            }
            //drop those first two rows which are empty
            data.shift();
            data.shift();
            //console.log(data);

            var count = 0;
            var sent = 0;
            $.each(data, function (i, value) {

                if(value['F20'] == undefined && value['F32AAMT'] == undefined){
                    return;
                }
                //Value Date
                var vd = value['F32AVD'];
                var valueDate = new Date((vd - (25569))*86400*1000);
                if(valueDate == 'Invalid Date'){
                    valueDate = new Date((value['ValueDate'] - (25569))*86400*1000);
                }

                var updateTime = value['UTCTime'];
                var requestType = value['TYPE'];
                var newData;
                if(requestType == 'REQUEST'){
                    count = count +1;
                    newData = {
                        "txnType": "payment_instruction",
                        "senderReference": value['F20'],
                        "payerBank": value['SenderFI'],
                        "beneBank": value['ReceiverFI'],
                        "currency": value['F32ACCY'],
                        "paymentAmount": value['F32AAMT'],
                        "feeType": value['F71A'],
                        "updateTime": updateTime,
                        "valueDate": valueDate,
                        "debitCredit":value['CRDR'],
                        "msgNum": value['MSG#'],
                        "msgType":value['MSG_TYPE'],
                        "localTime":updateTime,
                        "BookInt":value['FEE'],
                        "F33BCCY": value['F33BCCY'],
                        "F33BAMT":value['F33BAMT']
                    };
                }else if(requestType == 'CONFIRMATION'){
                    count = count +1;
                   newData = {
                        "txnType": "match_confirm_payment_instruction",
                        "senderReference": value['F20'],
                        "payerBank": value['SenderFI'],
                        "beneBank": value['ReceiverFI'],
                        "currency": value['F32ACCY'],
                        "paymentAmount": value['F32AAMT'],
                        "feeType": value['F71A'],
                        "updateTime": updateTime,
                        "valueDate": valueDate,
                        "debitCredit":value['CRDR'],
                        "msgNum": value['MSG#'],
                        "msgType":value['MSG_TYPE'],
                        "BookInt":value['FEE'],
                        "F33BCCY": value['F33BCCY'],
                        "F33BAMT":value['F33BAMT']
                    };
                }
                $('#loading-indicator').show();
                $( "div[id='UploadStatus']" ).hide();
                setTimeout(function() {
                    if(newData){
                        sent = sent +1;
                        console.log("Total request + confirmations sent : "+sent);
                        var chainCodeMessage = generateChaincodeMessage(newData);
                        callChaincode(JSON.stringify(chainCodeMessage),'sendMessage');
                        if(sent == count){
                            $('#loading-indicator').hide();
                            $( "div[id='UploadStatus']" ).show();
                        }
                    }
                }, 1000 * i);
            });
            console.log("Total Number of request + confirmations : "+count);
        });

    };
    reader.readAsBinaryString(f);
  }
}

function handleConfirmationFile(e) {
  var files = e.target.files;
  var i,f;
  for (i = 0, f = files[i]; i != files.length; ++i) {
    var reader = new FileReader();
    var name = f.name;
    reader.onload = function(e) {
      var data = e.target.result;


      var workbook = XLSX.read(data, {type: 'binary'});
        var sheet_name_list = workbook.SheetNames;
        sheet_name_list.forEach(function(y) {
            var worksheet = workbook.Sheets[y];
            var headers = {};
            var data = []; //var newData = {};
            for(z in worksheet) {
                if(z[0] === '!') continue;
                //parse out the column, row, and value
                var col = z.substring(0,1);
                var row = parseInt(z.substring(1));
                var value = worksheet[z].v;

                //store header names
                if(row == 1) {
                    headers[col] = value;
                    continue;
                }

                if(!data[row]) data[row]={};
                data[row][headers[col]] = value;
            }
            //drop those first two rows which are empty
            data.shift();
            data.shift();
            //console.log(data);


            $.each(data, function (i, value) {

                //Value Date
                var vd = value['F32AVD'];
                var updateTime = value['UTCTime'];
                var valueDate = new Date((vd - (25569))*86400*1000);
                if(valueDate == 'Invalid Date'){
                    valueDate = new Date((value['ValueDate'] - (25569))*86400*1000);
                }
                if(value['F20'] == undefined && value['F32AAMT'] == undefined){
                    return;
                }
                newData = {
                    "txnType": "match_confirm_payment_instruction",
                    "senderReference": value['F20'],
                    "payerBank": value['SenderFI'],
                    "beneBank": value['ReceiverFI'],
                    "currency": value['F32ACCY'],
                    "paymentAmount": value['F32AAMT'],
                    "feeType": value['F71A'],
                    "updateTime": updateTime,
                    "valueDate": valueDate,
                    "debitCredit":value['CRDR'],
                    "msgNum": value['MSG#'],
                    "msgType":value['MSG_TYPE'],
                    "BookInt":value['FEE'],
                    "F33BCCY": value['F33BCCY'],
                    "F33BAMT":value['F33BAMT']
                };
                setTimeout(function() {
                    var chainCodeMessage = generateChaincodeMessage(newData);
                    callChaincode(JSON.stringify(chainCodeMessage),'sendMessage');
                }, 1000 * i);

            });

        });

    };
    reader.readAsBinaryString(f);
  }
}



  $(function() {

    $(document).on('change', ':file', function(e) {

    var input = $(this).parents('.input-group').find(':text');
    input.val(e.target.files[0].name);

     if(e.target.id == "inputFile"){
        handleFile(e);
        $("#inputFile").val("");
     }else  if(e.target.id == "inputConfirmFile"){
        handleConfirmationFile(e);
        $("#inputConfirmFile").val("");
     }else  if(e.target.id == "graphInputFile"){
          handleGraph(e);
          $("#graphInputFile").val("");
       }
     $("p.status-block").html("File uploaded successfully.");
     });
  });