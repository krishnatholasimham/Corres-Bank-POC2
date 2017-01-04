#!/bin/bash

# Argument = -n name


# Save the current directory and navigate to the running directory. Enables
# startup scripts to be stored centrally.
startDIR=$PWD
runDIR=$GOPATH/src/github.com/openblockchain/obc-peer/openchain/container


# Help message for script.
usage() {
cat << EOF
usage: $0 options

This script will create a validating peer with the name specified in the -n flag.

OPTIONS:
  -h    Shows this message.
  -f    Starts validating peer using hyperledger fabric.
  -n    Validating Peer Name. Must end in an numerical value of increasing order.
        E.g. vp0, vp1, vp2, vp3, etc.
  -a    The IP address of the root node that this peer will connect to. Not required for first peer.
  -l    Specify log mode, case-insensitive (CRITICAL | ERROR | WARNING | NOTICE | INFO | DEBUG)
  -m    Membership mode enabled.
  -b    Sets this peer as byzantine.

EOF
}

# Initiate flag variables
name=""
address=""
debug="info"
makePeerImage="make peer-image"
endpoint="OPENCHAIN_VM_ENDPOINT=http://172.17.0.1:4243"
peerID="OPENCHAIN_PEER_ID"
autodetect="OPENCHAIN_PEER_ADDRESSAUTODETECT=true"
discovery="OPENCHAIN_PEER_DISCOVERY_ROOTNODE"
executable="openchain-peer obc-peer peer"
membership=false
fabric=false
multinodeSecurity=""
security=""
byzantine=""

# Process flags
while getopts "fhn:a:l:mb" OPTION
do # (colon) denotes a flag that requires a value.
  case $OPTION in
    f)  runDIR=$GOPATH/src/github.com/hyperledger/fabric;
        endpoint="CORE_VM_ENDPOINT=http://172.17.0.1:2375"
        peerID="CORE_PEER_ID";
        autodetect="CORE_PEER_ADDRESSAUTODETECT=true"
        discovery="CORE_PEER_DISCOVERY_ROOTNODE"
        executable="hyperledger/fabric-peer peer node start"
        fabric=true
        ;;
    h)  usage
        echo "1"
        exit 1;;
    n)  name=$OPTARG;;
    a)  address=$OPTARG;;
    l)  debug=$OPTARG;;
    m)  membership=true;;
    b)  byzantine=" -e CORE_PBFT_GENERAL_BYZANTINE=true";;
    ?)  usage
        echo "2"
        exit;;
  esac
done

# If membership is enabled, then determine whether OBC or Fabric is being used.
if [ $membership == true ]; then
  if [ $fabric == true ]; then
    security="-e CORE_SECURITY_ENABLED=$membership -e CORE_SECURITY_PRIVACY=$membership -e CORE_PEER_PKI_ECA_PADDR=172.17.0.1:50051 -e CORE_PEER_PKI_TCA_PADDR=172.17.0.1:50051 -e CORE_PEER_PKI_TLSCA_PADDR=172.17.0.1:50051 -e CORE_SECURITY_ENROLLID=$name -e CORE_SECURITY_ENROLLSECRET=${name}_secret"
  else
    security="-e OPENCHAIN_SECURITY_ENABLED=$membership -e OPENCHAIN_SECURITY_PRIVACY=$membership -e OPENCHAIN_PEER_PKI_ECA_PADDR=172.17.0.1:50051 -e OPENCHAIN_PEER_PKI_TCA_PADDR=172.17.0.1:50051 -e OPENCHAIN_PEER_PKI_TLSCA_PADDR=172.17.0.1:50051 -e OPENCHAIN_SECURITY_ENROLLID=$name -e OPENCHAIN_SECURITY_ENROLLSECRET=${name}_secret"
  fi
else
  security=""
fi

cd $runDIR
echo "NAME=$name. ADDRESS=$address. membership=$membership"
# Ensure both name and address flags are provided.
if [ -z "$name" ]; then # if string is null
  echo "MUST PROVIDE NAME FOR FIRST VALIDATING PEER OR NAME AND ENDPOINT FOR SUBSEQUENT PEERS"
  usage
  exit 1
elif [ -z "$address" ]; then
  echo "NAME PROVIDED. STARTING FIRST VP."
# TODO: parameterise peer image build based on fabric or OBC.
  # echo "Making peer image"
  # cd $GOPATH/src/github.com/hyperledger/fabric
  # make peer-image
  echo "[EXEC]-> docker run --rm -it -p 5000:5000 -e $endpoint -e $peerID=$name -e $autodetect $security$byzantine $executable --logging-level=$debug"
  eval "docker run --rm -it -p 5000:5000 -e $endpoint -e $peerID=$name -e $autodetect $security$byzantine $executable --logging-level=$debug"
else
  VPID=$(echo $name| cut -d'p' -f 2)
  echo "NAME AND ADDRESS PROVIDED. STARTING SUBSEQUENT PEER LINKED TO $address."
  echo "[EXEC]-> docker run --rm -it -p 500$VPID:5000 -e $endpoint -e $peerID=$name -e $autodetect -e $discovery="$address:30303" $security$byzantine $executable --logging-level=$debug"
  eval "docker run --rm -it -p 500$VPID:5000 -e $endpoint -e $peerID=$name -e $autodetect -e $discovery="$address:30303" $security$byzantine $executable --logging-level=$debug"
fi


# Check if all flags that require a value have provided a value.
#if [[ -n $name ]]
#then
#  usage
#  echo "3"
#  exit 1
#fi

# Run command to create validating peer, with name $name.
#eval "docker run --rm -it -e OPENCHAIN_VM_ENDPOINT=http://172.17.0.1:4243 -e OPENCHAIN_PEER_ID=$name -e OPENCHAIN_PEER_ADDRESSAUTODETECT=true openchain-peer obc-peer peer"

cd $startDIR
