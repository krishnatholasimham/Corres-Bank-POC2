<?php
header('Content-Type: text/event-stream');
header('Cache-Control: no-cache');

$data = file_get_contents("txnlog.txt");
// Trim the last new line character from the data
$trimmed = rtrim($data, "\n");
// Create an array in which each new line is an element
$data_array = explode("\n", $trimmed);


foreach($data_array as $item) {
    echo "data: $item\n\n";
}

echo "retry: 100\n";
ob_flush();
flush();
if(count($data_array) == 4){
file_put_contents("txnlog.txt", "");
}
?>