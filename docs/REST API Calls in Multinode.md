#Access the chaincode REST API within a Multinode Network

In order to access the Chaincode REST API in a Multinode Network you will need to add port forwarding in the Vagrant File and open up the corresponding ports when bringing up each Docker Container that will be used by the respective Validating Peers. 

Please follow the steps provided below.

### Step 1 : Updating the Vagrant File
In the `core.yaml` file you will see that the port which listens to REST API traffic is `5000`. However within one network environment (in our case this is the Vagrant box) you cannot have multiple Docker Containers listening to the same port. Therefore you will have expose multiple ports which will be used to communicate with different validating peers. 

For the moment let's open four ports given that we will have a four peer network. Please note that more peers you will have in the network, more ports you will have to open.
 
* If you have the Vagrant Box currently up and running, please exit it and shut it down by using the `vagrant halt` command.
* In your local file structure goto `$GOPATH/src/github.com/hyperledger/fabric/devenv` where you will find a file called `Vagrantfile`.
* Open the `Vagrantfile` using your favourite text editor.
* In the port forwarding section (located near line 38) add the following lines.

>  config.vm.network :forwarded_port, guest: 5001, host: 5001 # Openchain REST services
  
>  config.vm.network :forwarded_port, guest: 5002, host: 5002 # Openchain REST services

>  config.vm.network :forwarded_port, guest: 5003, host: 5003 # Openchain REST services
  
  * Save and close the file.

### Step 2 : Updating the Shell Script
Once the above step is complete, please get a copy of the updated `1-MNStartup.sh` from github (this is in the master). The updated `1-MNStartup.sh` would have added a `-p` parameter when starting each docker container to open up the respective port we forwarded in the Vagrant File.

Note: For the purposes of understanding the logic underneath, each Docker Container is listening to REST API traffic ONLY on port 5000. Therefore every time a REST API call is made on a port forwarded in the Vagrant file, we need to map that call back to port `5000` in the respective Docker Container. The script is modified to do this.

### Step 3 : Sending REST API Calls
Now start a four peer network using the instructions provided [here](https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/blob/master/docs/setup.md). To invoke different validating peers please use the following Host URLs.

* VP0 : `http://0.0.0.0:5000`
* VP1 : `http://0.0.0.0:5001`
* VP2 : `http://0.0.0.0:5002`
* VP3 : `http://0.0.0.0:5003`

Please have a look at [this](https://github.com/ANZ-Blockchain-Lab/Corres-Bank-POC/blob/master/docs/chaincodeAPI.md) document to see the different endpoints available for REST API calls.
