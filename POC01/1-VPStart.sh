#!/bin/bash

[ "$GOPATH" ] || { echo Please define GOPATH && exit 1; }

# Open blockchain defaults
workDIR=$GOPATH/src/github.com/openblockchain/obc-peer
compileCmd="go build"
runCmd="./obc-peer peer"
compile=false
run=true
logLevel="info"
membership=""

function usage() {
  cat<<EOF

Description:

  Peer compilation and startup for sandbox environment.

Usage:

  $(basename "$0") [-crnfvhm]

Flags:
  -c -r    compile (kept -r for backwards compatibility)
  -n       no run - default to run.
  -f       fabric compile and/or run.  defaults to open blockchain when lacking -f for now - will update default to fabric soon.
  -l       specify log mode, case-insensitive (critical | error | warning | notice | info | debug)
  -v       set logging to debug
  -m       enable security and privacy needed for membership
  -h       print this usage help message

EOF
}

while getopts "crnfvhml:" flag; do
  case "$flag" in
    c|r)  compile=true;;
    n  )  run=false;;
    f  )  workDIR=$GOPATH/src/github.com/hyperledger/fabric
          compileCmd="make peer"
          runCmd="peer node start"
          ;;
    m  )  membership="CORE_SECURITY_ENABLED=true CORE_SECURITY_PRIVACY=true";;
    l  )  logLevel=$OPTARG;;
    v  )  logLevel="debug";;
    *|h)  usage && exit;;
  esac
done

fullRunCmd="$membership $runCmd --peer-chaincodedev --logging-level=$logLevel"

$compile && {
  ( cd $workDIR && echo compile: $compileCmd && time $compileCmd; ) || { echo compilation failed, will not run && exit 1; }
}
${run} && {
  ( cd $workDIR && echo run: $fullRunCmd && eval "time $fullRunCmd"; ) || { echo error running peer && exit 1; }
}

