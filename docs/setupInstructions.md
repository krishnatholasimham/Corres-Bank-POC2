# Chaincode Setup Instructions


### Note - these are the pre fabric steps.  For Fabric, and latests setup steps please see [setup.md](setup.md)


**TABLE OF CONTENTS**
<!-- TOC depthFrom:2 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [1. Environment Setup](#1-environment-setup)
- [2. Sandbox vs Multinode Environments](#2-sandbox-vs-multinode-environments)
	- [2.1. Sandbox Environment](#21-sandbox-environment)
	- [2.2. Multinode Environment](#22-multinode-environment)
		- [2.2.1. Build obc-peer Project and Create Docker Image](#221-build-obc-peer-project-and-create-docker-image)
		- [2.2.2. TERMINAL 1 - First Validating Peer (vp0)](#222-terminal-1---first-validating-peer-vp0)
		- [2.2.3. TERMINAL 2 - Second Validating Peer (vp1)](#223-terminal-2---second-validating-peer-vp1)
		- [2.2.4. TERMINAL 3 - Open CLI](#224-terminal-3---open-cli)


<!-- /TOC -->

## 1. Environment Setup
### 1.1. Repo Cloning and VM Configuration
The current chaincode and related command line scripts are designed to be located outside of IBM's fabric directory (i.e. `obc-peer`). This is so that the `obc-peer` directory can be refreshed without the need to relocate chaincode files each time.

As a result, some VM configuration is required to ensure the chaincode directory (i.e. `Corres-Bank-POC`) is accessible from within the VM. Below are the end-to-end steps involved in setting up a local development environment with the nostro chaincode deployed.

1. Follow first 3 steps for [Setting up the development environment](https://github.com/ANZ-Blockchain-Lab/obc-docs/blob/master/dev-setup/devenv.md#setting-up-the-development-environment).
  1. [Set your GOPATH](https://github.com/ANZ-Blockchain-Lab/obc-docs/blob/master/dev-setup/devenv.md#set-your-gopath).
  2. [Cloning the Open Blockchain Peer project](https://github.com/ANZ-Blockchain-Lab/obc-docs/blob/master/dev-setup/devenv.md#cloning-the-open-blockchain-peer-project).
  3. [Cloning the Open Blockchain Development Environment project](https://github.com/ANZ-Blockchain-Lab/obc-docs/blob/master/dev-setup/devenv.md#cloning-the-open-blockchain-development-environment-project).

2. Before starting up the VM via Vagrant:
  1. Navigate to your `$GOPATH/src` directory and clone the `Corres-Bank-POC` repository.
  ```bash
cd $GOPATH/src
git clone https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC.git
```
  3. Navigate to the `obc-dev-env` directory.
  4. Open the `Vagrantfile` and add the following line to sync the local `Corres-Bank-POC` chaincode folder to the VM. This should be added around line 49.
  ```bash
config.vm.synced_folder "#{HOST_GOPATH}/src/Corres-Bank-POC", "/opt/gopath/src/Corres-Bank-POC"
```
3. Continue with obc instructions for [bootstrapping the VM using Vagrant](https://github.com/ANZ-Blockchain-Lab/obc-docs/blob/master/dev-setup/devenv.md#boostrapping-the-vm-using-vagrant)

Once you ssh into the VM, you should be able to access the chaincode and CLI folder using the following command:
```bash
cd $GOPATH/src/Corres-Bank-POC
```

## 2. Sandbox vs Multinode Environments
The IBM codebase provides two environments for development:
1. **Sandbox Environment:** A single node network that allows for rapid code, compile and test cycles. Ideal for quick chaincode development.
2. **Multinode Environment:** A network of multiple peer nodes, which allows more rigorous testing of blockchain capabilities (e.g. networking, consensus, membership, etc).

### 2.1. Sandbox Environment
~~While the current chaincode works in the sandbox environment, it is producing errors when run in a multinode development environment. An issue has been opened to resolve this.~~

A number of shell scripts have been written to run the chaincode.

Below are the steps required to run the current `nostro.go` chaincode.

#### 2.1.1. Build obc-peer Project
1. Navigate to the obc-peer project directory and build the binary:
```
cd $GOPATH/src/github.com/openblockchain/obc-peer
go build
```

2. Run the obc-peer binary to view the commands.
```bash
./obc-peer
```
You should see the following output:
```bash
Usage:
  obc-peer [command]

Available Commands:
  peer        Run openchain peer.
  status      Status of the openchain peer.
  stop        Stop openchain peer.
  login       Login user on CLI.
  vm          VM functionality of openchain.
  network     List of network peers.
  chaincode   chaincode specific commands.
  help        Help about any command

Flags:
  -h, --help[=false]: help for obc-peer
      --logging-level="": Default logging level and overrides, see openchain.yaml for full syntax


Use "obc-peer [command] --help" for more information about a command.

```

#### 2.1.2. TERMINAL 1 - Start obc-peer Project in Dev Mode.
In a new terminal, enter the Vagrant VM, navigate to the `POC01` directory and run the shell script `1-VPStart.sh` to start the validating peer. Note: you may need to edit the script permissions in the first instance.
```bash
vagrant ssh
cd $GOPATH/src/Corres-Bank-POC/POC01
chmod u+x ./1-VPStart.sh
./1-VPStart.sh
```
**Note:** The `-r` flag can be used to rebuild the `obc-peer` project before starting the validating peer.
* In essence, the script performs the following steps:
  - Navigates to the `$GOPATH/src/github.com/openblockchain/obc-peer` directory.
	- If the `-r` flag is set, builds the obc-peer project.
	- Starts the validating peer in sandbox mode with the logging flag set to debug using the following command: `./obc-peer peer --peer-chaincodedev --logging-level=debug`

You should get the following output:
```
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/Corres-Bank-POC/POC01$ ./1-VPStart.sh
01:31:42.569 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
01:31:42.592 [logging] LoggingInit -> DEBU 002 Setting default logging level to DEBUG for command 'peer'
01:31:42.596 [main] serve -> DEBU 003 Listen address not specified, using peer endpoint address
01:31:42.617 [eventhub_producer] AddEventType -> DEBU 004 registering block
01:31:42.619 [eventhub_producer] AddEventType -> DEBU 005 registering register
01:31:42.621 [main] serve -> INFO 007 Running in chaincode development mode
01:31:42.621 [main] serve -> INFO 008 Set consensus to NOOPS and user starts chaincode
01:31:42.621 [main] serve -> INFO 009 Disable loading validity system chaincode
01:31:42.623 [main] serve -> INFO 00a Security enabled status: false
01:31:42.623 [main] serve -> INFO 00b Privacy enabled status: false
01:31:42.619 [eventhub_producer] start -> INFO 006 event processor started
01:31:42.624 [main] serve -> DEBU 00c Running as validating peer - installing consensus noops
01:31:42.628 [db] createDBIfDBPathEmpty -> DEBU 00d Is db path [/var/openchain/production/db] empty [false]
01:31:42.699 [state] NewState -> INFO 00e Initializing state implementation [buckettree]
01:31:42.702 [buckettree] initConfig -> INFO 00f configs passed during initialization = map[string]interface {}{"numBuckets":10009, "maxGroupingAtEachLevel":10}
01:31:42.704 [buckettree] initConfig -> INFO 010 Initializing bucket tree state implemetation with configurations &{maxGroupingAtEachLevel:10 lowestLevel:5 levelToNumBucketsMap:map[0:1 5:10009 4:1001 3:101 2:11 1:2] hashFunc:0xa14f50}
01:31:42.708 [buckettree] computeCryptoHash -> DEBU 011 Appending crypto-hash for child bucket = [level=[1], bucketNumber=[1]]
01:31:42.709 [buckettree] computeCryptoHash -> DEBU 012 Propagating crypto-hash of single child node for bucket = [level=[0], bucketNumber=[1]]
01:31:42.711 [peer] chatWithPeer -> DEBU 013 Starting up the first peer
01:31:42.713 [chaincode] NewChaincodeSupport -> INFO 014 Chaincode support using peerAddress: 0.0.0.0:30303
01:31:42.719 [rest] StartOpenchainRESTServer -> INFO 016 Initializing the REST service...
01:31:42.719 [main] serve -> INFO 015 Starting peer with id=name:"jdoe" , network id=dev, address=0.0.0.0:30303, discovery.rootnode=, validator=true
```


#### 2.1.3. TERMINAL 2 - Start nostro.go Chaincode.
In a new terminal, enter the Vagrant VM, navigate to the `POC01` directory and run the shell script `2-CCStart.sh` to start and register the `nostro.go` chaincode with the validating peer in terminal 1. Note: again, you may need to edit the script permissions.
```bash
cd $GOPATH/src/Corres-Bank-POC/POC01
chmod u+x ./2-CCStart.sh
./2-CCStart.sh -r
```
**Note:** The `-r` flag is used to build the `nostro.go` chaincode before it is run.
* In essence, the script performs the following steps:
  - Navigates to the `$GOPATH/src/Corres-Bank-POC/POC01/chaincode/nostro` directory.
	- If the `-r` flag is set, builds the `nostro.go` chaincode.
	- Starts the chaincode using the following command: `OPENCHAIN_CHAINCODE_ID_NAME=anz OPENCHAIN_PEER_ADDRESS=0.0.0.0:30303 ./nostro`

On terminal 1 (PEER), you should get the following output:
```
01:35:24.426 [chaincode] HandleChaincodeStream -> DEBU 017 Current context deadline = 0001-01-01 00:00:00 +0000 UTC, ok = false
01:35:24.440 [chaincode] processStream -> DEBU 018 []Received message REGISTER from shim
01:35:24.442 [chaincode] HandleMessage -> DEBU 019 []Handling ChaincodeMessage of type: REGISTER in state created
01:35:24.444 [chaincode] beforeRegisterEvent -> DEBU 01a Received REGISTER in state created
01:35:24.445 [chaincode] registerHandler -> DEBU 01b registered handler complete for chaincode anz
01:35:24.446 [chaincode] beforeRegisterEvent -> DEBU 01c Got REGISTER for chaincodeID = name:"anz" , sending back REGISTERED
01:35:24.455 [chaincode] notifyDuringStartup -> DEBU 01d nothing to notify (dev mode ?)
```

On terminal 2 (CHAINCODE), you should get the following output:
```
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/Corres-Bank-POC/POC01$ ./2-CCStart.sh -r
2016/06/07 01:36:07 Peer address: 0.0.0.0:30303
2016/06/07 01:36:07 os.Args returns: [./nostro]
2016/06/07 01:36:07 Chaincode ID: anz
2016/06/07 01:36:07 Registering.. sending REGISTER
2016/06/07 01:36:07 []Received message REGISTERED from shim
2016/06/07 01:36:07 []Handling ChaincodeMessage of type: REGISTERED(state:created)
2016/06/07 01:36:07 Received REGISTERED, ready for invocations
```

#### 2.1.4. TERMINAL 3 - Open CLI
The `POC01-CLIv3.sh` CLI has been created to construct the commands passed to the validating peer. Details of the command elements are documented in the [obc-docs](https://github.com/openblockchain/obc-docs/blob/master/api/SandboxSetup.md#vagrant-terminal-3-cli-or-rest-api) repo.

In a new terminal, enter the Vagrant VM, navigate to the `CLI` directory, and run the `POC01-CLI.sh` script.
```bash
cd $GOPATH/src/Corres-Bank-POC/POC01/CLI
chmod u+x ./POC01-CLIv3.sh
./POC01-CLIv3.sh
```

A [demo sequence](https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/blob/master/POC01/demoSequence.md) is available on the github repo.

Current functionality includes:

##### RECORD MANAGEMENT
|Function | Description |
|---|---|
| **(F) Add Funding Message** | Add funds to a statement account, specified by the `account owner` and the `account holder`. |
| **(+) Add Payment Instruction** | Create a payment request entry in the distributed ledger. |
| **(C) Add Payment Confirmation** | Confirm an existing payment request entry in the distributed ledger. |
| **(-) Delete Range of Records** | Delete all entries on the distributed ledger within the specified key range. |


##### REPORTS
|Function | Description |
|---|---|
| **(p) PrettyPrint Range of Records** | Returns all distributed ledger entries within the specified key range. |
| **(G) Get Balance History** | Returns all balance-related entries (i.e. confirmations and funding messages) in chronological order for a statement account, specified by the `account owner` and the `account holder`. |
| **(B) Show All Statement Accounts for a Specified Bank** | Returns a list of all statement accounts owned by the specified bank, and all statement accounts held at the specified bank, along with the current balances. |
| **(P) Show Transaction Summary** | For a specified bank, returns a list of all confirmed and unconfirmed payments created or received by that bank. Also returns a list of all funding messages sent to that bank's statement accounts held at other banks. |


### 2.2. Multinode Environment
A multinode network is needed to test more blockchain-specific traits such as consensus and confidentiality. However, it is less suitable for rapid testing of chaincode, as each update requires the peer docker image to be rebuilt. 

At ~1gb in size, rebuilding the docker image can take ~40-50 seconds. An image is then created each time a peer node is started. These images will accumulate each time you do a rebuild cycle of the network, and need to be manually removed via the `docker rmi` command.

As such, the sandbox environment is recommended for in-depth chaincode development.

**Note:** A separate chaincode exists for multinode operation. This will eventually supersede the original `nostro.go` chaincode. Below are the steps required to run the `multinode-nostro.go` chaincode.

#### 2.2.1. Build obc-peer Project and Create Docker Image
1. Navigate to the obc-peer project directory and build the binary:
```
cd $GOPATH/src/github.com/openblockchain/obc-peer
go build
```
2. In the same directory, build the `openchain-peer` Docker Image.
You can view the current docker images using `docker images`.
```
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/github.com/openblockchain/obc-peer$ docker images
REPOSITORY                 TAG                 IMAGE ID            CREATED             SIZE
openblockchain/baseimage   latest              dc9e3a20b079        8 hours ago         990.2 MB
busybox                    latest              47bcc53f74dc        11 weeks ago        1.113 MB
openblockchain/baseimage   0.0.9               a21f7691fc12        11 weeks ago        990.2 MB
```
Build the `openchain-peer` docker image. **Note:** this may take up to 1 minute.
```bash
go test github.com/openblockchain/obc-peer/openchain/container -run=BuildImage_Peer
```
You should now see an `openchain-peer` image in the list of docker images.
```bash
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/github.com/openblockchain/obc-peer$ docker images
REPOSITORY                 TAG                 IMAGE ID            CREATED             SIZE
openchain-peer             latest              aacf63e26d33        4 seconds ago       1.159 GB
openblockchain/baseimage   latest              dc9e3a20b079        9 hours ago         990.2 MB
busybox                    latest              47bcc53f74dc        11 weeks ago        1.113 MB
openblockchain/baseimage   0.0.9               a21f7691fc12        11 weeks ago        990.2 MB
```
#### 2.2.2. TERMINAL 1 - First Validating Peer (vp0)
The fabric requires each validating peer to be named in lowercase letters and ending in a numerical value, starting with 0 and increasing by 1 for every node (i.e. vp0, vp1, vp2,...).

In a new terminal, navigate to the `POC01` folder use the `1-MNStartup.sh` script to start the validating peer. 
Use the `-n` flag to specify the peer name. 
Use the `-l` flag to specify the log mode (debug, info [default], error, warning).
```bash
cd $GOPATH/src/Corres-Bank-POC/POC01
./1-MNStartup.sh -n vp0 -l debug
```
You should see the following output:
```bash
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/Corres-Bank-POC/POC01$ ./1-MNStartup.sh -n vp0 -l debug
NAME=vp0. ADDRESS=.
NAME PROVIDED. STARTING FIRST VP.
[EXEC]-> docker run --rm -it -e OPENCHAIN_VM_ENDPOINT=http://172.17.0.1:4243 -e OPENCHAIN_PEER_ID=vp0 -e OPENCHAIN_PEER_ADDRESSAUTODETECT=true openchain-peer obc-peer peer --logging-level=info
12:34:16.035 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
12:34:16.042 [main] serve -> INFO 002 Security enabled status: false
12:34:16.042 [eventhub_producer] start -> INFO 003 event processor started
12:34:16.042 [main] serve -> INFO 004 Privacy enabled status: false
12:34:16.087 [state] NewState -> INFO 005 Initializing state implementation [buckettree]
12:34:16.088 [buckettree] initConfig -> INFO 006 configs passed during initialization = map[string]interface {}{"numBuckets":10009, "maxGroupingAtEachLevel":10}
12:34:16.089 [buckettree] initConfig -> INFO 007 Initializing bucket tree state implemetation with configurations &{maxGroupingAtEachLevel:10 lowestLevel:5 levelToNumBucketsMap:map[1:2 0:1 5:10009 4:1001 3:101 2:11] hashFunc:0xa14f50}
12:34:16.092 [chaincode] NewChaincodeSupport -> INFO 008 Chaincode support using peerAddress: 172.17.0.2:30303
12:34:16.093 [rest] StartOpenchainRESTServer -> INFO 009 Initializing the REST service...
12:34:16.095 [main] serve -> INFO 00a Starting peer with id=name:"vp0" , network id=dev, address=172.17.0.2:30303, discovery.rootnode=, validator=true
12:34:16.096 [genesis] func1 -> INFO 00b Creating genesis block.
12:34:16.097 [genesis] func1 -> INFO 00c No genesis block chaincodes defined.
```

#### 2.2.3. TERMINAL 2 - Second Validating Peer (vp1)
In a new terminal, navigate to the `POC01` folder use the `1-MNStartup.sh` script to start the validating peer. 
However, we need to specify an additional flag, `-a`, which notifies this peer of the ip address of the previous peer (i.e. vp0).
THe relevant ip address can be located in terminal 1:
```bash
12:34:16.095 [main] serve -> INFO 00a Starting peer with id=name:"vp0" , network id=dev, address=172.17.0.2:30303, discovery.rootnode=, validator=true
```
Start the second validating peer.
```bash
cd $GOPATH/src/Corres-Bank-POC/POC01
./1-MNStartup.sh -n vp1 -a 172.17.0.2 -l debug
```
You should see the following output:
```bash
vagrant@obc-devenv:v0.0.9-db6d2a6:/opt/gopath/src/Corres-Bank-POC/POC01$ ./1-MNStartup.sh -n vp1 -a 172.17.0.2 -l debug
NAME=vp1. ADDRESS=172.17.0.2.
NAME AND ADDRESS PROVIDED. STARTING SUBSEQUENT PEER LINKED TO 172.17.0.2.
[EXEC]-> docker run --rm -it -e OPENCHAIN_VM_ENDPOINT=http://172.17.0.1:4243 -e OPENCHAIN_PEER_ID=vp1 -e OPENCHAIN_PEER_ADDRESSAUTODETECT=true -e OPENCHAIN_PEER_DISCOVERY_ROOTNODE=172.17.0.2:30303 openchain-peer obc-peer peer --logging-level=info
12:46:21.055 [crypto] main -> INFO 001 Log level recognized 'info', set to INFO
12:46:21.058 [main] serve -> INFO 002 Security enabled status: false
12:46:21.058 [eventhub_producer] start -> INFO 003 event processor started
12:46:21.059 [main] serve -> INFO 004 Privacy enabled status: false
12:46:21.088 [state] NewState -> INFO 005 Initializing state implementation [buckettree]
12:46:21.089 [buckettree] initConfig -> INFO 006 configs passed during initialization = map[string]interface {}{"numBuckets":10009, "maxGroupingAtEachLevel":10}
12:46:21.092 [buckettree] initConfig -> INFO 007 Initializing bucket tree state implemetation with configurations &{maxGroupingAtEachLevel:10 lowestLevel:5 levelToNumBucketsMap:map[5:10009 4:1001 3:101 2:11 1:2 0:1] hashFunc:0xa14f50}
12:46:21.094 [chaincode] NewChaincodeSupport -> INFO 008 Chaincode support using peerAddress: 172.17.0.3:30303
12:46:21.095 [rest] StartOpenchainRESTServer -> INFO 009 Initializing the REST service...
12:46:21.097 [main] serve -> INFO 00a Starting peer with id=name:"vp1" , network id=dev, address=172.17.0.3:30303, discovery.rootnode=172.17.0.2:30303, validator=true
12:46:21.098 [genesis] func1 -> INFO 00b Creating genesis block.
12:46:21.099 [genesis] func1 -> INFO 00c No genesis block chaincodes defined.
12:46:22.105 [consensus/noops] newNoops -> INFO 00d NOOPS consensus type = *noops.Noops
12:46:22.105 [consensus/noops] newNoops -> INFO 00e NOOPS block size = 500
12:46:22.106 [consensus/noops] newNoops -> INFO 00f NOOPS block timeout = 1s
```

#### 2.2.4. TERMINAL 3 - Open CLI
The `MultinodeCLI.sh` script is used to construct commands that are passed to the validating peer.
__Note:__ the following flags can be used with the `MultinodeCLI.sh`:

  -t  set the timezone to be recorded against transactions, where n is the timezone in Country/City format. UTC if blank.

  -s  configures commands to operate in a sandbox environment (i.e. where `1-VPStart.sh` is used, or the --chaincodedev flag is used with `./obc-peer`).

```bash
cd $GOPATH/src/Corres-Bank-POC/POC01/CLI
chmod u+x ./MultinodeCLI.sh
./MultinodeCLI.sh
```

##### Notes for running the chaincode
1. In the first instance, run `(I)nstall Chaincode` to deploy `multinode-nostro.go` to each validating peer. The fabric returns the chaincode name as a result. This is stored by `MultinodeCLI.sh` for use in future commands.
2. If you exit and re-enter the CLI, you can specify the chaincode name using `(S)pecify Chaincode Name`.
3. For each CLI command, you will be prompted for the ID of the validating peer you wish to communicate with. This is simply the number of the vp name (e.g. 0, 1, 2,...).
