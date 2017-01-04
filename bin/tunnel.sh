#!/usr/bin/env bash

function vsg() {
  cat $tmp | grep "$1" | tr -d '"' | awk '{print $2}'
}

function vssh() {
  tmp=$(mktemp /tmp/abc-script.XXXXXX)
  vagrant ssh-config >> $tmp
  ssh `vsg "User "`@`vsg HostName` -p `vsg Port` -o Compression=yes -o DSAAuthentication=yes -o LogLevel=`vsg LogLevel` -o StrictHostKeyChecking=`vsg StrictHostKeyChecking` -o UserKnownHostsFile=`vsg UserKnownHostsFile` -o IdentitiesOnly=`vsg IdentitiesOnly` -i `vsg IdentityFile` $*
  rm "$tmp"
}

function tunnel() {
  defaultTunnel="-L 5000:172.17.0.2:5000 cat -"
  [ $# = 0 ] && tunnelArgs="$defaultTunnel" || tunnelArgs="$*"
  echo tunneling with: "$tunnelArgs"
  echo "( Pass any other arguments to map another host port to an internal peer.  For instance '-L 5002:172.17.0.2:5000 cat -' will map host 5002 to the internal 172.17.0.2:5000"
  echo "  To test port forwarding on your laptop use: curl localhost:5000/chain"
  echo "  to update the port next to the POST command to something other than 5000: vi $GOPATH/src/github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/Blockchain\ POC\ UI/assets/api.js"
  echo "  If a port is in use find the pid and kill -9 <pid>, i.e. if your port is 5002: ps auxgww | egrep '5002|USER' | grep -v [g]rep )"

  vssh $tunnelArgs
}
