#!/bin/bash
workDIR=$GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro
compileCmd="go build multinode-nostro.go"
runCmd="OPENCHAIN_CHAINCODE_ID_NAME=anz OPENCHAIN_PEER_ADDRESS=0.0.0.0:30303 ./multinode-nostro $debug"
compile=false
run=true

function usage() {
  cat<<EOF

Description:

  Chaincode compilation and startup for sandbox environment.

Usage:

  $(basename "$0") [-crnfvhm]

Flags:
  -c -r    compile (kept -r for backwards compatibility)
  -n       no run - default to run.
  -f       fabric compile and/or run.  defaults to open blockchain when lacking -f for now - will update default to fabric soon.
  -h       print this usage help message

EOF
}

while getopts "crnfhl:" flag; do
  case "${flag}" in
    c|r)  compile=true;;
    n  )  run=false;;
    f  )  runCmd="CORE_CHAINCODE_ID_NAME=anz CORE_PEER_ADDRESS=0.0.0.0:30303 ./multinode-nostro";;
    l  )  debug="-logMode $OPTARG";;
    *|h)  usage && exit;;

  esac
done

fullRunCmd="$runCmd" # keeping here for consistency with other start file.

$compile && {
  ( cd $workDIR && echo compile: $compileCmd && time $compileCmd; ) || { echo compilation failed, will not run && exit 1; }
}
${run} && {
  ( cd $workDIR && echo run: $fullRunCmd && eval "time $fullRunCmd"; ) || { echo error running multinode-nostro && exit 1; }
}
