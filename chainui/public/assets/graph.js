
/*function loadDashboard() {
    var ANZBal =  getUnconfirmedBalanceHistory("ANZ","WF");
    var WFBal =  getUnconfirmedBalanceHistory("WF","ANZ");
}*/
dateDropdownData = [];dateDropdownDatawoBlockchain = [];
var oJS;
var chartwithBlockchain, chartwoBlockchain;
function populateCharts(task){
    dateDropdownData = [];

    if(localStorage.bankName == 'WF'){
        sendReqForCurrentView(graphDataWF, graphDataStmtBalWF);
    } else {
        sendReqForCurrentView(graphDataANZ, graphDataStmtBalANZ);
    }
}



function plotGraphwithoutBlockChain(graphData, filterDate){
    dataPointswoBlockchain = [];
    if(oJS){
        var count = 0;
        var sent = 0;
        var dataFromFile = [];
        var graphDataWithoutBC = [];
        if(localStorage.bankName == 'WF'){
            graphDataWithoutBC = graphDataStmtBalWithoutBCWF;
        }else{
            graphDataWithoutBC = graphDataStmtBalWithoutBCANZ;
        }
        var famt = 0;
        $.each(oJS, function(index, oJSActualData) {
                var time = new Date(oJSActualData.UPDTIME);
                var newMT950File = {
                    "UpdateTime": time,
                    "Amount": oJSActualData.F32AAMT,
                    "Currency":oJSActualData.F33BCCY,
                    "SenderFI":oJSActualData.SenderFI,
                    "ReceiverFI":oJSActualData.ReceiverFI,
                    "F20": oJSActualData.F20,
                    "F60M": oJSActualData.F60M,
                    "F60F" : oJSActualData.F60F
                };
                dataFromFile.push(newMT950File);

        });


        //Sort the values by date
        dataFromFile.sort(function(x,y){
            if(new Date(x.UpdateTime) > new Date(y.UpdateTime))
                return 1;
            else if(new Date(x.UpdateTime) == new Date(y.UpdateTime))
                return 1;
            return -1;
        });
        var orignalAmount = 0;
        var f60mAmt = 0;
        var f60fAmt = 0;
        var f60mCount = 0;
        $.each(dataFromFile, function(index, sortedData) {
            tempData = [];
            var time = sortedData.UpdateTime;
            tempData.push(new Date(time.getFullYear(), time.getMonth(), time.getDate(), time.getHours(), time.getMinutes(), time.getSeconds()));
            if(sortedData.F60M){
                f60mCount = f60mCount+1;
                f60mAmt = f60mAmt + Number(sortedData.Amount);
                orignalAmount = Number(orignalAmount) + Number(sortedData.Amount);
            }else{
                orignalAmount = orignalAmount - Number(sortedData.Amount);
            }
            console.log("F60M : "+f60mAmt);
            tempData.push(null);
            tempData.push(orignalAmount);
            dataPointswoBlockchain.push(tempData);
        });
        console.log("f60mCount " + f60mCount);
    }

    $.each(graphData, function(index, graphDatapoint) {
        var time = new Date(graphDatapoint[1]);
        var dateOnly = graphDatapoint[1].split('T')[0];
        var found = false;
        if(filterDate == null || filterDate == "" || filterDate == dateOnly){
            tempData = [];
            tempData.push(new Date(time.getFullYear(), time.getMonth(), time.getDate(), time.getHours(), time.getMinutes(), time.getSeconds()));
            tempData.push(graphDatapoint[0]);
            tempData.push(null);
            alreadythere = false;
            $.each(dateDropdownDatawoBlockchain, function(i, dateStr){
                if(dateStr == dateOnly){
                    alreadythere = true;
                }
            });
            if(!alreadythere){
                dateDropdownDatawoBlockchain.push(dateOnly);
            }
            dataPointswoBlockchain.push(tempData);
        }
    });

    if(filterDate == null){
        $('#DateDropdownwoBlockchain').html('');
        $('#DateDropdownwoBlockchain').append('<option value="">Select Date</option>');
        $.each(dateDropdownDatawoBlockchain, function(i, dateStr){
            $('#DateDropdownwoBlockchain').append('<option value="' + dateStr + '">' + dateStr + '</option>')
        });
    }

    // Plot the graph
    if(dataPointswoBlockchain.length > 0){
        google.charts.setOnLoadCallback(function() {drawBackgroundColor(dataPointswoBlockchain, 'chartwoblockchain_div', chartwoBlockchain, 'Actual Balance','950 Statement Balance'); });
    }
}

function sendReqForCurrentView(graphData, graphStatementData, filterDate){
    dataPoints = [];
    $.each(graphData, function(index, graphDatapoint) {
        var time = new Date(graphDatapoint[1]);
        var dateOnly = graphDatapoint[1].split('T')[0];
        var found = false;
        if(filterDate == null || filterDate == "" || filterDate == dateOnly){
            tempData = [];
            tempData.push(new Date(time.getFullYear(), time.getMonth(), time.getDate(), time.getHours(), time.getMinutes(), time.getSeconds()));
            tempData.push(null);
            tempData.push(graphDatapoint[0]);
            alreadythere = false;
            $.each(dateDropdownData, function(i, dateStr){
                if(dateStr == dateOnly){
                    alreadythere = true;
                }
            });
            if(!alreadythere){
                dateDropdownData.push(dateOnly);
            }
            dataPoints.push(tempData);
        }
    });

    $.each(graphStatementData, function(j, graphDatapoint) {
        found = false;
        var chaincodetime = new Date(graphDatapoint[1]);
        var time = new Date(graphDatapoint[1]);
        var dateOnly = graphDatapoint[1].split('T')[0];
        if(filterDate == null || filterDate == "" || filterDate == dateOnly){
            $.each(dataPoints, function(index, datapoint) {
                var dataPointTime = datapoint[0];
                if(dataPointTime == (time.getHours()+ ":" + time.getMinutes() + ":" + time.getSeconds())){
                    found = true;
                }
            });
            if( !found ) {
                tempData = [];
                tempData.push(new Date(time.getFullYear(), time.getMonth(), time.getDate(), time.getHours(), time.getMinutes(), time.getSeconds()));
                tempData.push(graphDatapoint[0]);
                tempData.push(null);
                alreadythere = false;
                $.each(dateDropdownData, function(i, dateStr){
                    if(dateStr == dateOnly){
                        alreadythere = true;
                    }
                });
                if(!alreadythere){
                    dateDropdownData.push(dateOnly);
                }
                dataPoints.push(tempData);
            }
        }

    });

    if(filterDate == null){
        $('#DateDropdown').html('');
        $('#DateDropdown').append('<option value="">Select Date</option>');
        $.each(dateDropdownData, function(i, dateStr){
            $('#DateDropdown').append('<option value="' + dateStr + '">' + dateStr + '</option>')
        });
    }

    // Plot the graph
    if(dataPoints.length > 0){
        google.charts.setOnLoadCallback(function() {drawBackgroundColor(dataPoints, 'chart_div',chartwithBlockchain,'Actual Balance','Indicative Balance'); });
    }

}

function timeToSeconds(time) {
    time = time.split(/:/);
    //return time[0] * 3600 + time[1] * 60;
    return parseInt(time[0] + "" + time[1]);
}

var chart;
function drawBackgroundColor(dataPoints, chartElementId, chartObj, balanceDesc1, balanceDesc2) {
      var data = new google.visualization.DataTable();
      //data.addColumn('string', 'Time of Day');
      data.addColumn('datetime', 'Time of Day');
      data.addColumn('number', balanceDesc1);
      data.addColumn('number', balanceDesc2);

      data.addRows(dataPoints);
      //data.setColumns([0,1,2]);

      var options = {
        height : 500,
        width: 1200,
        lineWidth: 2,
        series: {
            1 : { lineDashStyle: [5, 1] }
        },
        interpolateNulls: true,
        hAxis: {
          title: 'Time'
        },
        vAxis: {
          title: 'Statement Account Balance',
          gridlines: {
              color: 'transparent'
          }
        },
        view : {'columns': [0, 1]},
        backgroundColor: '#FFFFFF',
          color: '#FF0000'
      };

      chartObj = new google.visualization.LineChart(document.getElementById(chartElementId));
      //chart.bind(programmaticSlider);
      chartObj.draw(data, options);
    }


    $(document).ready(function() {
      $('#DateDropdown').change(function() {
        theVal = $('#DateDropdown').val();
        if(localStorage.bankName == 'WF'){
            sendReqForCurrentView(graphDataWF, graphDataStmtBalWF, theVal);
        } else {
            sendReqForCurrentView(graphDataANZ, graphDataStmtBalANZ, theVal);
        }
      });
      $('#DateDropdownwoBlockchain').change(function() {
          theVal = $('#DateDropdownwoBlockchain').val();
          if(localStorage.bankName == 'WF'){
              plotGraphwithoutBlockChain(graphDataWF, theVal);
          } else {
              plotGraphwithoutBlockChain(graphDataANZ, theVal);
          }
      });
    });

function handleGraph(e) {
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
            oJS = XLS.utils.sheet_to_row_object_array(workbook.Sheets[y]);
            if(localStorage.bankName == 'WF'){
                plotGraphwithoutBlockChain(graphDataStmtBalWF);
            } else {
                plotGraphwithoutBlockChain(graphDataStmtBalANZ);
            }
        });

    };
    reader.readAsBinaryString(f);
  }
}