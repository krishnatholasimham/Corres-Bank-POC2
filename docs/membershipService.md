# Enable Membership & Confidentiality
**Table of Contents**

[toc]

The following set of instructions explains how to activate Membership & Confidentiality Services in the POC application. Once the above security services are activated, the application will require an enrolled Client and a Validating Peer to deploy, invoke and query functionalities with the chaincode.

Built on top of the Membership Service, **Confidentiality** is a feature which aims to keep sensitive business information from being disclosed to unintended parties. Within a financial transaction context, any other party other than the Payer Bank and the Beneficiary Bank of a transaction are deemed as unintended parties to access that transaction.

In the POC built, even-though every financial institution maintains a copy of all transactions by a synchronised world state, any given financial institution should only be allowed to access financial transactions which aligns to the above logic. The confidentiality implementation in the POC application attempts to achieve this goal.

## 1. Setup Instructions

### 1.1 Enable Security and Privacy - core.yaml
Enabling security will force every entity on the network to enrol with the ECA (Enrollment Certificate Authority) and maintain a valid set of certificates in order to communicate with peers. Enabling privacy of transactions (requires security to be enabled) encrypts the transaction content during transit and at rest including the state data.

To enable security and privacy:

##### Option 1 - manual update alternative

1. Open the `core.yaml` file located at `$GOPATH/src/github.com/hyperledger/fabric/peer/core.yaml`.
2. Set the `security.enabled` value to `true` **before building** the peer executable. Alternatively, you can enable security by running the peer with environment variable `CORE_SECURITY_ENABLED=true`.
3. Set the `privacy.enabled` value to `true` as well. Alternatively, you can enable privacy by running the peer with environment variable `CORE_SECURITY_PRIVACY=true`
  **Note:** Our current `1-MNStartup.sh` script sets the default value of `CORE_SECURITY_ENABLED` to `false`, unless the `-m` flag is specified to signal membership mode.
4. Save and close the file.

##### Option 2 - Just use our scripts

1.  Run the peer with a `-m` flag passed to our 1-MNStartup.sh` or `1-VPStart.sh` - that will set the flags needed for you.



### 1.2 Add Users and Roles - membersrvc.yaml

**Note :** for all updates to membersrvc.yaml, can either use the sample file as follows:

- Flip secureLocalEnv = true on utils.js locally
- Use docs/membersvcs.secure-sample.yaml to overwrite hyperledger/fabric/membersrvc/membersrvc.yaml
- Follow steps in this doc to set the env (we can exact summarize those here later, specifically those are cleanDocker.sh and startMembershipServer.sh, other env related stuff should be in setup.md)
- To test:
  - Login (enroll) users, wf1, anz1, boa1 (on CLI) pass is test.
  - To Track on CLI: 
    - Use (k) on CLI to show all keys before adding a transaction.
    - Use CLI (k), find new key and use (Q) to query it from CLI - for BOA will be garbled, for anz/wf will show properly.
  - To Add and track on UI
    - Use UI to add funding/payment.
    - See transactions and data shows when logged in at wf/anz.
    - See garbled data shows when logged in as boa.

Or follow the detailed steps below.


The `eca.users` section informs the network of the different users who are authorised to perform transactions and the roles assigned to them. There are 4 roles, defined as:

* 1 - simple client such as a wallet: CLIENT
* 2 - non-validating peer: PEER
* 4 - validating peer: VALIDATOR
* 8 - auditing client: AUDITOR

**Note :** The numeric values above indicates the assigned role.  

To test this functionality you need to add yourself as a new client who is attached to a predefined bank.

1. Open the `membersrvc.yaml` file located at `$GOPATH/src/github.com/hyperledger/fabric/membersrvc`.
2. Scroll down to `eca.affiliations` section located around line 54.
3. Add the following entries under the `banks` section.

		anz #Australia and New Zealand Banking Group Limited
		wf #Wells Fargo Bank
		ba #Bank of America
		lb #Llyods Bank

	**Note :** In order to demonstrate **Confidentiality** we require at least 3 banks to represent an unintended third party to a given transaction.

1. Scroll down to `eca.users` section located around line 66.
4. Add yourself as a `client` after the last entry of clients and set any password you like. The syntax for adding a new user is `<username: role password affiliation affiliation_role>`. E.g. `binhn: 1 7avZQLwcUe9q institution_a     00005`

	**Note :** Please ensure you use an affiliation role which is higher than *00001* as that will be reserved for the Admin Client of that institution. This will be explained further as we progress through the document.

5. After the entry which represents yourself add 3 more users as below who have affiliations with other financial institutions (the three example entries below assume that your affiliation is with the *anz* bank. If this is not the case change one of the users to a different affiliation so that no two users will have the same affiliation).


		amir: 1 amir1 wf 00002
		john: 1 john1 ba 00002
		paul: 1 paul1 lb 00002


6. Add the entries for the Admin clients of each financial institution after the above users.


		vp0Admin: 1 vp0admin_secret anz 00001
		vp1Admin: 1 vp1admin_secret wf 00001
		vp2Admin: 1 vp2admin_secret ba 00001
		vp3Admin: 1 vp3admin_secret lb 00001

	**Note :** DO NOT alter these records as they are being used in the golang code to register the Admin users upon executing the deploy command. Also note that the affiliation role *00001* of each financial institution has been used by the respective Admin client.  

4. Add the validating peers to the users section after the `vp` entry. Validating peers entries follow the same syntax as above. Each validating peer **must** have the password set as `<insert validating_peer_name>_secret`.
5. Scroll down to the `aca.attributes` section (around line 162) and add the user Enrolment ID as an attribute to each user created above. The syntax for adding the attribute is `attribute-entry-<id>: <username>;<affiliation>;enrolment;<username>;2015-01-01T00:00:00-03:00;;`

	Once completed the attribute section should look something similar this;

		attribute-entry-12: heshansp;anz;enrolment;heshansp;2015-01-01T00:00:00-03:00;;
		attribute-entry-13: amir;wf;enrolment;amir;2015-01-01T00:00:00-03:00;;
		attribute-entry-14: john;ba;enrolment;john;2015-01-01T00:00:00-03:00;;
		attribute-entry-15: paul;lb;enrolment;paul;2015-01-01T00:00:00-03:00;;


6. Scroll down to `aca.enabled` variable and set it to `true`.


5. Save an close the file.

Once all of the above is complete, the `users` section must look something similar to this.

```
eca:
        # This hierarchy is used to create the Pre-key tree, affiliations is the top of this hierarchy, 'banks_and_institutions' is used to create the key associated to auditors of both banks and
        # institutions, 'banks' is used to create a key associated to auditors of banks, 'bank_a' is used to create a key associated to auditors of bank_a, etc.
        affiliations:
           banks_and_institutions:
              banks:
                  - bank_a
                  - bank_b
                  - bank_c
                  - anz #Australia and New Zealand Banking Group Limited
                  - wf #Wells Fargo Bank
                  - ba #Bank of America
                  - lb #Llyods Bank
              institutions:
                  - institution_a
        users:
                #
                # The fields of each user are as follows:
                #    <EnrollmentID>: <system_role (1:client, 2: peer, 4: validator, 8: auditor)> <EnrollmentPWD> <Affiliation> <Affiliation_Role> <JSON_Metadata>
                #
                # The optional JSON_Metadata field is of the following format:
                #   { "registrar": { "roles": <array-of-role-names>, "delegateRoles": <array-of-role-names> } }
                # The 'registrar' section is used to control access to registration of new users directly via the ECAA.RegisterUser GRPC call.
                # (See the 'fabric/membersrvc/protos/ca.proto' file for the definition of ECAA.RegisterUser.)
                # Note that this also controls who can register users via the client SDK.
                #
                # Only users with a 'registrar' section may be a registrar to register other users.  In particular,
                # 1) the "roles" field specifies which member roles may be registered by this user, and
                # 2) the "delegateRoles" field specifies which member roles may become the "roles" field of registered users.
                # The valid role names are "client", "peer", "validator", and "auditor".
                #
                # Example1:
                #    The 'admin' user below can register clients, peers, validators, or auditors; furthermore, the 'admin' user can register other
                #    users who can then register clients only.
                #
                # Example2:
                #    The 'WebAppAdmin' user below can register clients only, but none of the users registered by this user can register other users.
                #
                admin: 1 Xurw3yU9zI0l institution_a 00001 '{"registrar":{"roles":["client","peer","validator","auditor"],"delegateRoles":["client"]}}'
                WebAppAdmin: 1 DJY27pEnl16d institution_a 00002 '{"registrar":{"roles":["client"]}}'
                lukas: 1 NPKYL39uKbkj bank_a 00001
                system_chaincode_invoker: 1 DRJ20pEql15a institution_a 00002
                diego: 1 DRJ23pEQl16a institution_a 00003
                jim: 1 6avZQLwcUe9b bank_a 00004
                binhn: 1 7avZQLwcUe9q institution_a 00005
                heshansp: 1 heshan1 anz 00002
                amir: 1 amir1 wf 00002
                john: 1 john1 ba 00002
                paul: 1 paul1 lb 00002

                # Users for asset transfer with roles test located at
                # sdk/node/test/unit/asset-mgmt-with-roles.js
                alice: 1 CMS10pEQlB16 bank_a 00006
                bob: 1 NOE63pEQbL25 bank_a 00007
                assigner: 1 Tc43PeqBl11 bank_a 00008

                vp0Admin: 1 vp0admin_secret anz 00001
                vp1Admin: 1 vp1admin_secret wf 00001
                vp2Admin: 1 vp2admin_secret ba 00001
                vp3Admin: 1 vp3admin_secret lb 00001

                vp: 4 f3489fy98ghf
                vp0: 4 vp0_secret
                vp1: 4 vp1_secret
                vp2: 4 vp2_secret
                vp3: 4 vp3_secret

                test_vp0: 4 MwYpmSRjupbT
```

### 1.3 Docker Cleanup
This step deletes old docker images and containers. **Note: This step is mandatory to demonstrate the features of confidentiality**.

1. Open a terminal and ssh into your vagrant box.
2. Clean up the unused docker images and containers by running the `cleanDocker.sh` script.

  ```bash
  cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI && ./cleanDocker.sh
  ```

Once you have cleaned the docker images and the containers;

1. ```docker ps -a``` command should output something similar to following.

		vagrant@hyperledger-devenv:v0.0.10-5651cfc:~$ docker ps -a
		CONTAINER ID        IMAGE                                        COMMAND                  CREATED             STATUS                    PORTS               NAMES
		dac0a298fd35        hyperledger/fabric-baseimage                 "go install github.co"   26 hours ago        Exited (0) 26 hours ago                       tiny_mirzakhani
		7377dbda482d        hyperledger/fabric-baseimage:x86_64-0.0.10   "/bin/bash -l -c prin"   26 hours ago        Exited (0) 26 hours ago                       goofy_kalam

2. ```docker images``` command should output something similar to following.

		vagrant@hyperledger-devenv:v0.0.10-5651cfc:~$ docker images
		REPOSITORY                     TAG                 IMAGE ID            CREATED             SIZE
		hyperledger/fabric-ccenv       latest              d51086b071b6        24 hours ago        1.777 GB
		hyperledger/fabric-src         latest              4959efc10f68        24 hours ago        1.76 GB
		hyperledger/fabric-baseimage   latest              dbeb248f74ac        26 hours ago        1.687 GB
		busybox                        latest              2b8fd9751c4c        5 weeks ago         1.093 MB
		hyperledger/fabric-baseimage   x86_64-0.0.10       a2392cc113fd        9 weeks ago         1.076 GB

### 1.4 Start Certificate Authority (CA) server (Terminal 1)
The Certificate Authority Server is responsible for enrolling users and managing the certificates allocated to each user. Following steps will allow you to start the CA Server.

1. Open a terminal and ssh into your vagrant box (alternatively use the same terminal used in section 1.3).
2. Run the following command:
  ```bash
  cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI && ./startMembershipServer.sh
  ```

Running the above command clears out the old crypto profiles and builds and runs the CA server with the default setup, which is defined in the `membersrvc.yaml` configuration file. The default configuration includes multiple users who are already registered with the CA; these users are listed in the 'users' section of the configuration file.

### 1.5 Build Chaincode, Fabric Project and Create Docker Image (Terminal 2)
The following should be executed from within Vagrant.

1. Build Chaincode

  ```bash
  cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro/ && echo && echo building chaincode: && go build multinode-nostro.go
  ```		

2. Build Fabric Project and the Peer Docker Image

  ```bash
	cd $GOPATH/src/github.com/hyperledger/fabric/ && echo && echo making peer: && make peer && echo && echo making peer-image: && make peer-image && echo
  ```

  You should get the following output:
  
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
		
		
## 2. Running with Membership Services

### 2.1 Sandbox Environment

1. Go to `POC01` directory and run `1-VPStart.sh`
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-VPStart.sh -rfm
```

2. In a new terminal, go to `POC01` directory and run `2-CCStart.sh`
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./2-CCStart.sh -rf
```

3. In a new terminal, go to `POC01/CLI` directory and run `MultinodeCLI.sh`
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI/ && ./MultinodeCLI.sh -fsm
```

### 2.2 Multinode Environment

#### 2.2.1 Start First Validating Peer - vp0 (Terminal 2)
For the following steps you may use the same terminal used in section 1.5. Provided the security is enabled in the network, you **must** ensure to start the validating peers in a secure context where the commands will carry the credentials to be authenticated at the CA Server.

1. Navigate to the POC01 folder use the `1-MNStartup.sh` script to start the validating peer.
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -fm -n vp0 -l debug
```

The `-m` flag indicates the Membership Services are active. Once the above command is executed, you should see the following output.

#### 2.2.2 Start Second Validating Peer - vp1 (Terminal 3)
In a new terminal, navigate to the POC01 folder use the ```1-MNStartup.sh``` script to start the validating peer. However, we need to specify an additional flag, -a, which notifies this peer of the ip address of the previous peer (i.e. vp0). The relevant ip address can be located in terminal 2.

Same as above the second validating peer must be started in a secure context by passing the relevant authentication credentials.

1. Start the second validating peer.
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -fm -n vp1 -a 172.17.0.2 -l debug
```

#### 2.2.3 Start Third Validating Peer - vp2 (Terminal 4)
In a new terminal, navigate to the POC01 folder use the ```1-MNStartup.sh``` script to start the validating peer. However, we need to specify an additional flag, -a, which notifies this peer of the ip address of the previous peer (i.e. vp1). The relevant ip address can be located in terminal 3.

Same as above the third validating peer must be started in a secure context by passing the relevant authentication credentials.

1. Start the third validating peer.
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -fm -n vp2 -a 172.17.0.3 -l debug
```

#### 2.2.4 Start Fourth Validating Peer - vp3 (Terminal 5)
In a new terminal, navigate to the POC01 folder use the ```1-MNStartup.sh``` script to start the validating peer. However, we need to specify an additional flag, -a, which notifies this peer of the ip address of the previous peer (i.e. vp2). The relevant ip address can be located in terminal 4.

Same as above the fourth validating peer must be started in a secure context by passing the relevant authentication credentials.

1. Start the fourth validating peer.
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/ && ./1-MNStartup.sh -fm -n vp3 -a 172.17.0.4 -l debug
```

#### 2.2.3 Open CLI (Terminal 4)  
Provided the security is enabled, only an enrolled user is allowed to deploy, invoke or query the chaincode functions. Therefore when opening the CLI, you must indicate your intension to transact in a secure environment.

1. Navigate to `CLI` folder and execute `MultinodeCLI.sh` script.
```bash
cd $GOPATH/src/github.com/hyperledger/fabric/Corres-Bank-POC/POC01/CLI && chmod u+x ./MultinodeCLI.sh && ./MultinodeCLI.sh -fm
```

The `-m` flag indicates the Membership Services are active.

### 2.3 Chaincode Instructions
As mentioned earlier given that security is enabled, only an enrolled user is allowed to deploy, invoke or query the chaincode functions. Therefore the very first thing to do next is to enroll the user you created in section 1.2.

#### 2.3.1 Enroll User

1. Press `L` and hit enter to initiate user enrollment.
2. Nominate your validating peer (ex: 0).
2. When prompted type your username.
3. Hit enter.
4. When prompted provide the same password you added in section 1.2
5. Hit enter.

In the nominated validating peer terminal you should see the following appear.

	11:48:11.833 [crypto] Login -> INFO 18ae Registering client [heshansp] with name [heshansp]...
	11:48:11.867 [crypto] register -> INFO 18af [client.heshansp] Register crypto engine...
	11:48:11.867 [crypto] register -> INFO 18b0 [client.heshansp] Register crypto engine...done.
	11:48:11.867 [crypto] Login -> INFO 18b1 Registering client [heshansp] with name [heshansp]...done!

#### 2.3.2 Deploy Chaincode
1. Press `I` and hit enter to initiate deploy process.
2. When prompted hit `y` and hit enter.

To deploy the chaincode, you **must** be enrolled at vp0 as the deployment always takes place at vp0. If everything went well, you should see an output similar to the following.

	CORE_PEER_ADDRESS=172.17.0.2:30303 peer chaincode deploy -u heshansp -p github.com/hyperledger/fabric/Corres-Bank-POC/POC01/chaincode/multinode-nostro -c '{"Function":"init", "Args": ["a","100"]}
	CHAINCODE NAME: 5fce71c0473b737ddddd8bfbe6b45e7f3c89160b8f27f2034ec65fc7d1203a4701aae503b6c2f0eb198bed6ee5e4d673dd5983ab24e38ffa3928505c5dc5efdb`

**Note :** In a privacy enabled blockchain implementation, the deployment will take longer than usual (~1-2 mins) as the entire peer database will be encrypted prior to being synchronised with other peers. So please be patient and allow this time for the deployment to be successful.

#### 2.3.3 Synchronise Crypto Material

1. Test Get Keys query

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

2. Open a new terminal and ssh into your vagrant box.
3. Type `docker ps -a` and hit enter. You should see something like the following

		vagrant@hyperledger-devenv:v0.0.10-5651cfc:~$ docker ps -a
		CONTAINER ID        IMAGE                                        COMMAND                  CREATED              STATUS                     PORTS                    NAMES
		55ee73784e41        hyperledger/fabric-peer                      "peer node start --lo"   About a minute ago   Up About a minute          0.0.0.0:5003->5000/tcp   lonely_noether
		c2a559192a8c        hyperledger/fabric-peer                      "peer node start --lo"   About a minute ago   Up About a minute          0.0.0.0:5002->5000/tcp   stoic_lumiere
		212380f781f0        hyperledger/fabric-peer                      "peer node start --lo"   About a minute ago   Up About a minute          0.0.0.0:5001->5000/tcp   nauseous_kowalevski
		9d9eb9147528        hyperledger/fabric-peer                      "peer node start --lo"   2 minutes ago        Up 2 minutes               0.0.0.0:5000->5000/tcp   zen_jennings
		ee1c3b085ddb        hyperledger/fabric-src                       "go install github.co"   6 minutes ago        Exited (0) 6 minutes ago                            big_newton
		0b47c99b7940        hyperledger/fabric-src                       "go install github.co"   6 minutes ago        Exited (0) 6 minutes ago                            dreamy_lumiere
		dac0a298fd35        hyperledger/fabric-baseimage                 "go install github.co"   28 hours ago         Exited (0) 28 hours ago                             tiny_mirzakhani
		7377dbda482d        hyperledger/fabric-baseimage:x86_64-0.0.10   "/bin/bash -l -c prin"   28 hours ago         Exited (0) 28 hours ago                             goofy_kalam

4. Find the docker container IDs of the validating peers.

	**Note :** The validating peer which was created first is `VP0`, second is `VP1`, third is `VP2` and so on. Hence in this example The Docker Container ID of `VP0` is `9d9eb9147528` and the Docker Container ID of `VP1` is `212380f781f0` and so on.

5. In the CLI press `T` and hit enter.
6. When prompted provide the Docker Container ID of `VP0`.
7. When prompted, notify the chaincode that you are running 4 validating peers.
8. Provide the Docker container of each validating peer to synchronise the Crypto Material of Admin Clients.

## 3. Restarting Validating Peers
Enrollment of a validating peer is a once-off activity, and will fail if retried.
The temporary certificate files need to be removed, and the CA server rebuilt before a network can be reinitiated.

1. Stop current CA server using `ctrl-c`.
2. Delete certificate files and remake directory:
  ```bash
rm -rf /var/hyperledger/production/ && mkdir /var/hyperledger/production/
```
3. Rebuild and start CA server:
  ```bash
cd $GOPATH/src/github.com/hyperledger/fabric && make membersrvc && membersrvc
```
4. Restart your network using [2.2 Multinode Environment](#22-multinode-environment)
