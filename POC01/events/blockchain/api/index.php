<?php
if (isset($_POST['peer']) 
        && isset($_POST['txn_uuid']) 
        && isset($_POST['seconds']) 
        && isset($_POST['nanos'])
        && isset($_POST['height'])
        && isset($_POST['currentBlockHash'])
        && isset($_POST['previousBlockHash'])
        ){
    
    $timestamp = getDatetime($_POST['seconds'], $_POST['nanos']);
    
    $handle = fopen("txnlog.txt", "a");
    fwrite($handle, 
            $_POST['peer']."|".
            $_POST['txn_uuid']."|".
            $timestamp."|".
            $_POST['height']."|".
            $_POST['currentBlockHash']."|".
            $_POST['previousBlockHash']."|".
            "\n");
    fclose($handle);
    exit();
}

function getDatetime($second, $nanos){
    $ts = $second;
$date = new DateTime("@$ts");
return $date->format('Y-m-d H:i:s').".".$nanos;
}

?>