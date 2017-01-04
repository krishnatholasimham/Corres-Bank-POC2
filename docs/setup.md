
#### Steps to install

**TABLE OF CONTENTS**

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Run Sandbox / UI env](#run-sandbox--ui-env)
- [Run Multinode / CLI env](#run-multinode--cli-env)
- [Current Issues](#current-issues)

##### Overview

- We have migrated from IBM's Openblockchain to Hyperledger's fabric (i.e. "the fabric").
- The fabric needs a new VM.
- Note on setup time - the time it takes to bring the VM up which is only needed once can vary.  It can be anywhere from
  ~5 minutes on a fast network (100MBit/sec) to 30+ minutes, and up to 2+ hours on a mifi which runs a
  slow 1-3MBit/sec in office (and we know it can be ~40MBit/sec elsewhere).
- We plan to discontinue use of IBM's Openblockchain, but in the interim, you can set up dev environments for both frameworks on a single machine.
- Some tips for maintaining both VMs on one machine:
  - Only one VM can run at a time.
  - To switch between VMs, first suspend or halt your current VM via the `vagrant suspend` or `vagrant halt` command.
  - Navigate to the directory of the VM you wish to switch to, and start it up via the `vagrant up` command.
- If you wish to remove the Openblockchain framework from your machine:
  - Navigate to the obc-dev-env directory.
  - Run `vagrant halt` to shut down the VM.
  - Run `vagrant destroy` to release all resources used to provision the VM.
  - You can then safely remove the relevant Openblockchain directories. Namely obc-dev-env and github.com/openblockchain.

##### Prerequisites

- [Fabric prerequisites](https://github.com/hyperledger/fabric/blob/master/docs/dev-setup/devenv.md#prerequisites)
- [Go path](https://github.com/hyperledger/fabric/blob/master/docs/dev-setup/devenv.md#set-your-gopath)

##### Installation

1. Remove old Hyperledger Fabric VM (if prior environment existed).
  ```bash
echo ===========================================================
echo == NOTE: Be careful this will destroy the current VM !!  ==
echo == Please keep this exact location, and move your other  ==
echo == existing working dir out for now.  We can implement   ==
echo == something more convenient which does not need to      ==
echo == rename the dirs/symlinks later on.                    ==
echo ===========================================================
cd $GOPATH/src/github.com/hyperledger/fabric/devenv
vagrant destroy
```

2. Clone Fabric and set to specific revision.
  ```bash
cd $GOPATH/src
mkdir -p github.com/hyperledger && cd github.com/hyperledger
git clone https://github.com/ANZ-Blockchain-Lab/fabric.git && cd fabric && git checkout 5651cfcb11a5749ae479937705d479856045db6c
```

3. Clone chaincode project (Corres-Bank-POC) into the fabric directory as a submodule.
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric
git submodule add https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC
```

4. Add the Corres-Bank-POC folder to the fabric's git index (the index is used to generate the docker container's file structure).
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric
git add Corres-Bank-POC
```
  **Note:** You can check that the index has been successfully updated by running `git ls-files | less` and looking for `Corres-Bank-POC` in the list.

5. Provision VM
  - This step does take time so please be patient. Also being on a faster connection will reduce wait time.
  - See more details on times on the [overview](#overview)
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/devenv && time vagrant up
```

  **Notes for the following steps**
  - The example outputs below result from commands executed from the local machine, and from the `fabric/devenv` directory. If the above steps were followed, this should map to  `$GOPATH/src/github.com/hyperledger/fabric/devenv`.
  - Please check that your output for each command matches the one in the sample here for validation.

6. Enter VM and check docker images and containers
  ```bash
fabric/devenv $
vagrant ssh -c 'docker images; echo; docker ps -a'
REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
hyperledger/fabric-ccenv       latest              704500cbe8c3        2 minutes ago       1.745 GB
hyperledger/fabric-src         latest              113cd57901a6        2 minutes ago       1.728 GB
hyperledger/fabric-baseimage   latest              5844e1518133        2 minutes ago       1.687 GB
busybox                        latest              2b8fd9751c4c        4 weeks ago         1.093 MB
hyperledger/fabric-baseimage   x86_64-0.0.10       a2392cc113fd        8 weeks ago         1.076 GB

CONTAINER ID        IMAGE                                        COMMAND                  CREATED             STATUS                      PORTS               NAMES
02b47e860190        hyperledger/fabric-baseimage                 "go install github.co"   2 minutes ago       Exited (0) 2 minutes ago                        nostalgic_hopper
d170e95132fb        hyperledger/fabric-baseimage:x86_64-0.0.10   "/bin/bash -l -c prin"   13 minutes ago      Exited (0) 13 minutes ago                       stoic_leakey
Connection to 127.0.0.1 closed.
```

7. Build chaincode
  ```bash
fabric/devenv $
vagrant ssh -c 'cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro/ && echo && echo building chaincode: && go build multinode-nostro.go'

building chaincode:
Connection to 127.0.0.1 closed.
```

8. Build fabric project and the peer docker image
  ```bash
fabric/devenv $
vagrant ssh -c 'cd /opt/gopath/src/github.com/hyperledger/fabric/ && echo && echo making peer: && make peer && echo && echo making peer-image: && make peer-image && echo'

making peer:
make: Nothing to be done for 'peer'.

making peer-image:
Building build/docker/bin/examples/events/block-listener
Building build/docker/bin/peer
Building docker peer-image
cp build/docker/bin/peer build/image/peer/bin
docker build -t hyperledger/fabric-peer:latest build/image/peer
Sending build context to Docker daemon 26.73 MB
Step 1 : FROM hyperledger/fabric-src:latest
 ---> 4e27932832ec
Step 2 : RUN mkdir -p /var/hyperledger/db
 ---> Running in 6fd77221f346
 ---> 8a64a6a8a348
Removing intermediate container 6fd77221f346
Step 3 : COPY bin/* $GOPATH/bin/
 ---> 9ceeb9931557
Removing intermediate container 732e13e7adcd
Step 4 : WORKDIR $GOPATH/src/github.com/hyperledger/fabric
 ---> Running in cf6ee7464f1a
 ---> 60dcfb7a3cd3
Removing intermediate container cf6ee7464f1a
Successfully built 60dcfb7a3cd3
rm build/docker/bin/peer

Connection to 127.0.0.1 closed.
```
  **Note:** You can also force the peer image to be built using the following command:
  ```bash
  go test github.com/hyperledger/fabric/core/container -run=BuildImage_Peer
  ```

9. Check docker images
  ```bash
fabric/devenv $
vagrant ssh -c 'docker images'

REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
hyperledger/fabric-peer        latest              6a374aee6a73        6 minutes ago       1.755 GB
hyperledger/fabric-ccenv       latest              704500cbe8c3        12 minutes ago      1.745 GB
hyperledger/fabric-src         latest              113cd57901a6        12 minutes ago      1.728 GB
hyperledger/fabric-baseimage   latest              5844e1518133        13 minutes ago      1.687 GB
busybox                        latest              2b8fd9751c4c        4 weeks ago         1.093 MB
hyperledger/fabric-baseimage   x86_64-0.0.10       a2392cc113fd        8 weeks ago         1.076 GB

Connection to 127.0.0.1 closed.
```

  **Notes for most following steps**
  - Spawn multiple shells off the same directory (such as: `~/dev/lang/go/gopath/src/github.com/hyperledger/fabric/devenv`).
  - Vagrant ssh on each then cont below steps *inside the vagrant machine*.
  - Some commands for outside vagrant will be clearly marked with `user@laptop:path $` prefix.


##### Run Sandbox / UI env

1.  Run UI only in sandbox mode.

- Term 1:
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-VPStart.sh -rf
```

- Term 2:
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./2-CCStart.sh -rf
```

- Term 3: Run CLI in fabric / standalone mode to validate
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/ && ./MultinodeCLI.sh -fs
Select the option to (I) Install Chaincode
```

- Open the local index.html file in a browser (**outside vagrant, on your laptop**):
   ```bash
user@laptop:path $
open -a "Google Chrome" $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/chainui/public/index.html
echo This command will also work with "Firefox" or "Safari", yet possibly best to stick w/ chrome whcih ui is developed on.
```


##### Run Multinode / CLI env

1. Start peers and CLI for multi-node env:

  - Start first validating peer (vp0)
  ```bash
cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -f -n vp0 -l debug
```

  - In new terminal within the VM, Start second validating peer (vp1)
  ```bash
cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -f -n vp1 -a 172.17.0.2 -l debug
```

  - In new terminal within the VM, Start third validating peer (vp2)
  ```bash
cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -f -n vp2 -a 172.17.0.3 -l debug
```

  - In new terminal within the VM, Start fourth validating peer (vp3)
  ```bash
cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -f -n vp3 -a 172.17.0.4 -l debug
```

  - In new terminal within VM, Check that four containers are running, one per peer
  ```bash
vagrant@hyperledger-devenv:v0.0.10-a94f5a7:~$ docker ps
CONTAINER ID        IMAGE                     COMMAND                  CREATED             STATUS              PORTS                    NAMES
25f0bba90947        hyperledger/fabric-peer   "peer node start --lo"   18 seconds ago      Up 17 seconds       0.0.0.0:5003->5000/tcp   suspicious_ride
96eef6dcb3da        hyperledger/fabric-peer   "peer node start --lo"   27 seconds ago      Up 26 seconds       0.0.0.0:5002->5000/tcp   angry_allen
7a4722d230db        hyperledger/fabric-peer   "peer node start --lo"   2 minutes ago       Up 2 minutes        0.0.0.0:5001->5000/tcp   insane_yalow
7e5dd8c12c97        hyperledger/fabric-peer   "peer node start --lo"   2 minutes ago       Up 2 minutes        0.0.0.0:5000->5000/tcp   distracted_liskov
```

  - Start CLI
  ```bash
cd /opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/ && ./MultinodeCLI.sh -f
```

2. Install chaincode and test Get Keys query
  ```bash
===========================================
ANZ Nostro Reconciliations Blockchain POC01
-------------------------------------------
What would you like to do?

(I)nstall Chaincode
(S)pecify Chaincode Name
(M)enu
(E)xit

RECORD MANAGEMENT
-------------------------------------------
(F) Add Funding Message
(+) Add Payment Instruction
(C)onfirm Payment
(D)elete Record
(-) Delete Range of Records

REPORTS
-------------------------------------------
(Q)uery Record
(K) Get All Keys
(G) Get Balance History
(B) Show All Statement Accounts for a Specified Bank
(P) Show Transaction Summary

-------------------------------------------
i
You selected INSTALL CHAINCODE

Chaincode only needs to be deployed once per VM. Are you sure? [y/N] y
CORE_PEER_ADDRESS=172.17.0.2:30303 ./peer chaincode deploy -p Corres-Bank-POC/POC01/chaincode/multinode-nostro -c '{"Function":"init", "Args": ["a","100"]}'
06:05:09.441 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO

CHAINCODE NAME: e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5

What else would you like to do? k
You selected GET ALL KEYS

Please enter ID number of validating peer: 0

CORE_PEER_ADDRESS=172.17.0.2:30303 ./peer chaincode query -n e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5 -c '{"Function":"keys","Args":[""]}'

06:06:30.252 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
[
  "a"
]

What else would you like to do?
```

  **Note:** if you get an error as per the below when running the Get Keys query, it may mean that the peer network is still in the process of creating docker images of the chaincode for each peer as part of the deployment process. You have two options:
  - Wait a few seconds and retry until it works.
  - Exit and re-enter `MultinodeCLI.sh` and re-run the install command.
  ```bash
What else would you like to do? k
You selected GET ALL KEYS

Please enter ID number of validating peer: 0

CORE_PEER_ADDRESS=172.17.0.2:30303 ./peer chaincode query -n e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5 -c '{"Function":"keys","Args":[""]}'

06:06:00.115 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO

Usage:
  peer chaincode query [flags]

Flags:
  -x, --hex[=false]: If true, output the query value byte array in hexadecimal. Incompatible with --raw
  -r, --raw[=false]: If true, output the query value as raw bytes, otherwise format as a printable string


Global Flags:
  -c, --ctor="{}": Constructor message for the chaincode in JSON format
  -l, --lang="golang": Language the chaincode is written in
      --logging-level="": Default logging level and overrides, see core.yaml for full syntax
  -n, --name="": Name of the chaincode returned by the deploy transaction
  -p, --path="": Path to chaincode
  -u, --username="": Username for chaincode operations when security is enabled

Error: Error querying chaincode: rpc error: code = 2 desc = "Error:Failed to launch chaincode spec(Could not get deployment transaction for e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5 - LedgerError - ResourceNotFound: ledger: resource not found)"
```

3. Check docker images
  ```bash
vagrant@hyperledger-devenv:v0.0.10-a94f5a7:/opt/gopath/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI$ docker images
REPOSITORY                                                                                                                                 TAG                 IMAGE ID            CREATED             SIZE
dev-vp2-e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5   latest              8775f5569c6e        7 minutes ago       1.469 GB
dev-vp3-e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5   latest              85b2dd188b59        7 minutes ago       1.469 GB
dev-vp0-e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5   latest              2dfbf55d6778        7 minutes ago       1.469 GB
dev-vp1-e35a4831627a9bd7408d0916cb6af8a6893ba06e38ea567f4220f567301b8f24d97e8f86187051331575f64730ad857594f52255f7e42a366274d71c9c42c5b5   latest              79f064ec809d        7 minutes ago       1.469 GB
hyperledger-peer                                                                                                                           latest              8dee504471b2        19 minutes ago      1.526 GB
hyperledger/fabric-baseimage                                                                                                               latest              980f2b3c6da4        58 minutes ago      1.384 GB
busybox                                                                                                                                    latest              2b8fd9751c4c        6 hours ago         1.093 MB
hyperledger/fabric-baseimage                                                                                                               x86_64-0.0.10       a2392cc113fd        4 weeks ago         1.076 GB
```



#### Current Issues
