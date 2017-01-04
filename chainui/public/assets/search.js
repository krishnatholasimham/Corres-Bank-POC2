$('#startDate, #endDate').on('propertychange change keyup input paste', function() {
    filterDate();
});
$('#statusFilter').on('change', function() {
    addToFilter('STATUS', $(this).find('option:selected')[0].value);
});
$('#messageTypeFilter').on('change', function() {
    addToFilter('MESSAGETYPE', $(this).find('option:selected')[0].value);
});
$('#StartAmountRange, #EndAmountRange').on('propertychange change keyup input paste', function() {
    filterAmount();
});
$('#SearchKey').on('propertychange change keyup input paste', function() {
    filterByKey();
});
$('#filters').on('change', function() {
    console.log($(this).val()); // the selected options’s value

    // if you want to do stuff based on the OPTION element:
    var opt = $(this).find('option:selected')[0];
    switch (opt.value){
        case 'KEY':
            showKeyFilter();
            break;
        case 'DATERANGE':
            showDateRange();
            break;
        case 'AMOUNTRANGE':
            showAmountRange();
            //addToFilter('CONFIRMATION', 'CONFIRMATION');
            break;
        case 'MESSAGETYPE':
            showMessageTypeFilter();
            //addToFilter('FUNDING', 'FUNDING');
            break;
        case 'STATUS':
            showStatusFilter();
            //addToFilter('STATUS', 'FEE');
            break;
        default:
            break;
    }
    // use switch or if/else etc.
});
function showKeyFilter(){
    $('#SearchKeyFilter').removeClass("displayNone").addClass("displayInline");
    $('#Status, #StartAmountRange, #EndAmountRange, #MessageType,#StartDatePicker, #EndDatePicker').removeClass("displayInline").addClass("displayNone");
}
function showDateRange(){
    $('#StartDatePicker, #EndDatePicker').removeClass("displayNone").addClass("displayInline");
    $('#Status, #StartAmountRange, #EndAmountRange, #MessageType, #SearchKeyFilter').removeClass("displayInline").addClass("displayNone");
}
function showAmountRange(){
    $('#StartAmountRange, #EndAmountRange').removeClass("displayNone").addClass("displayInline");
    $('#Status, #StartDatePicker, #EndDatePicker, #MessageType, #SearchKeyFilter').removeClass("displayInline").addClass("displayNone");
}
function showMessageTypeFilter(){
    $('#MessageType').removeClass("displayNone").addClass("displayInline");
    $('#Status, #StartAmountRange, #EndAmountRange, #StartDatePicker, #EndDatePicker, #SearchKeyFilter').removeClass("displayInline").addClass("displayNone");
}
function showStatusFilter(){
    $('#Status').removeClass("displayNone").addClass("displayInline");
    $('#StartDatePicker, #EndDatePicker, #StartAmountRange, #EndAmountRange, #MessageType, #SearchKeyFilter').removeClass("displayInline").addClass("displayNone");
}
function addToFilter(filterType, value){
    $('#filterField_' + filterType).remove();
    $('#filterField').clone().find('[id]').andSelf().each(function () { this.id = this.id + '_' + filterType }).appendTo('#filterFields');
    $('#filterField_' + filterType).removeClass("displayNone").addClass("displayInline").html(value + '  <a id="close_' + filterType +'" style="margin-left:10px;">x</a>').attr('data-filterType', filterType).attr('data-filterVal', value);
    $('#close_' + filterType).click(function() {
        $(this).parent().remove();
        filterTrans();
    });
    filterTrans();
}
function filterDate(){
    var stDate = $('#startDate').val();
    var edDate = $('#endDate').val();
    if(stDate && stDate != '' && edDate && edDate != ''){
        addToFilter('DATERANGE', stDate + ' - ' + edDate);
    }
}
function filterAmount(){
    var stAmt = $('#StartAmount').val();
    var edAmt = $('#EndAmount').val();
    if(stAmt && stAmt != '' && $.isNumeric(stAmt) && edAmt && edAmt != '' && $.isNumeric(edAmt)) {
        addToFilter('AMOUNTRANGE', stAmt + ' - ' + edAmt);
    }
}
function filterByKey(){
    var keyVal = $('#SearchKey').val();
    if(keyVal && keyVal != '') {
        addToFilter('KEY', keyVal);
    }
}
function filterTrans(){
    var jo = $("#searchTransactions").find("tr");
    filtered = false;
    if($('#filterFields .filterElement').length > 1){
        jo.hide();
        $('#filterFields .filterElement').each(function(index, value) {
            var fType = $(this).attr('data-filterType');
            if(fType === 'KEY'){
                filtered = true;
                aVal = $(this).attr('data-filterVal');
                jo.filter(function (i, v) {
                    var $t = $(this);
                    if($(this).find('td').length > 0){
                        theVal = $($(this).find('td')[10]).text() ;
                        if(theVal == aVal){
                            return true;
                        } else {
                            return false;
                        }
                    } else {
                        return true;
                    }
                }).show();
            }  else if(fType == 'AMOUNTRANGE'){
                filtered = true;
                aVal = $(this).attr('data-filterVal');
                amtVal = aVal.split(' - ');
                jo.filter(function (i, v) {
                    var $t = $(this);
                    if($(this).find('td').length > 0){
                        theVal = parseFloat($($(this).find('td')[5]).text()) ;
                        if((parseFloat(amtVal[0]) <= theVal) && (theVal <= parseFloat(amtVal[1]))){
                            return true;
                        } else {
                            return false;
                        }
                    } else {
                        return true;
                    }
                }).show();
            } else if(fType === 'DATERANGE'){
                filtered = true;
                dVal = $(this).attr('data-filterVal');
                dateVal = dVal.split(' - ');
                jo.filter(function (i, v) {
                    var $t = $(this);
                    if($(this).find('td').length > 0){
                        theCreatedVal = $($(this).find('td')[6]).text() ;
                        theConfirmedVal = $($(this).find('td')[7]).text() ;
                        if(((new Date(dateVal[0]) <= new Date(theCreatedVal)) && (new Date(theCreatedVal) <= new Date(dateVal[1])))
                            || ((new Date(dateVal[0]) <= new Date(theConfirmedVal)) && (new Date(theConfirmedVal) <= new Date(dateVal[1]))) ){
                            return true;
                        } else {
                            return false;
                        }
                    } else {
                        return true;
                    }
                }).show();
            } else if(fType === 'MESSAGETYPE'){
                filtered = true;
                messageType = $('#messageTypeFilter').find('option:selected')[0].value;
                jo.filter(function (i, v) {
                    var $t = $(this);
                    if($(this).find('td').length > 0){
                        theVal = $($(this).find('td')[0]).text() ;
                        if((theVal == 'REQUEST-RECORD' && messageType == 'MT103')
                            || (theVal == 'CONFIRMATION-RECORD' && messageType == 'MT940')){
                            return true;
                        } else {
                            return false;
                        }
                    } else {
                        return true;
                    }
                }).show();
            } else if(fType === 'STATUS'){
                filtered = true;
                status = $('#statusFilter').find('option:selected')[0].value
                jo.filter(function (i, v) {
                    var $t = $(this);
                    if($(this).find('td').length > 0){
                        theVal = $($(this).find('td')[2]).text() ;
                        if((theVal == '' && status == 'UnConfirmed')
                            || (theVal != '' && status == 'Confirmed')){
                            return true;
                        } else {
                            return false;
                        }
                    } else {
                        return true;
                    }
                }).show();
            }
        });
        if(!filtered){
            jo.show();
        }
    } else {
        jo.show();
    }




}
$("#searchTrn").click(function(){
    searchTransactions();
});
$("#searchKey").click(function(){
    searchKey();
});

