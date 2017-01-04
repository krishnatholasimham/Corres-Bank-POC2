/*
Copyright (C) Australia and New Zealand Banking Group Limited (ANZ)
833 Collins Street, Docklands 3008, ABN 11 005 357 522.
Unauthorized copying of this file, via any medium is strictly prohibited
Proprietary and confidential
Written by Heshan Peiris <heshan.peiris@anz.com> August 2016
*/

package main

import (
	"flag"
	"fmt"
	"os"
	"bytes"
	"strconv"

	"github.com/hyperledger/fabric/events/consumer"
	pb "github.com/hyperledger/fabric/protos"
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"encoding/json"
)

var chaincodeName = ""
var peer = "jdoe"
type adapter struct {
	chaincode	   chan *pb.Event_ChaincodeEvent
	notify             chan *pb.Event_Block
}

type blockInfo struct {
	CurrentBlockHash  string `json:"currentBlockHash"`
	Height            int    `json:"height"`
	PreviousBlockHash string `json:"previousBlockHash"`
}

//GetInterestedEvents implements consumer.EventAdapter interface for registering interested events
func (a *adapter) GetInterestedEvents() ([]*pb.Interest, error) {
	return []*pb.Interest{
		&pb.Interest{EventType: pb.EventType_CHAINCODE, RegInfo: &pb.Interest_ChaincodeRegInfo{ChaincodeRegInfo: &pb.ChaincodeReg{ChaincodeID: chaincodeName, EventName: "*"}}},
		&pb.Interest{EventType: pb.EventType_BLOCK},
	}, nil
}


//Recv implements consumer.EventAdapter interface for receiving events
func (a *adapter) Recv(msg *pb.Event) (bool, error) {

	if o, e := msg.Event.(*pb.Event_ChaincodeEvent); e {
		a.chaincode <- o
		return true, nil
	}
	if o, e := msg.Event.(*pb.Event_Block); e {
		a.notify <- o
		return true, nil
	}
	return false, nil
}

//Disconnected implements consumer.EventAdapter interface for disconnecting
func (a *adapter) Disconnected(err error) {
	fmt.Printf("Disconnected...exiting\n")
	os.Exit(1)
}

func createEventClient(eventAddress string) *adapter {
	var obcEHClient *consumer.EventsClient

	chaincodeEvent := make(chan *pb.Event_ChaincodeEvent)
	BlockEvent := make(chan *pb.Event_Block)

	adapter := &adapter{chaincode: chaincodeEvent, notify: BlockEvent}
	obcEHClient = consumer.NewEventsClient(eventAddress, adapter)
	if err := obcEHClient.Start(); err != nil {
		fmt.Printf("could not start chat %s\n", err)
		obcEHClient.Stop()
		return nil
	}

	return adapter
}

func sendRequest(listner string, uuid string, seconds string, nanos string){
	apiUrl := "http://10.0.2.2"
	resource := "/blockchain/api/"

	fmt.Println("Peer : ", peer)

	height, currentBlockHash,previousBlockHash := getBlockInfo(peer)


	data := url.Values{}
	data.Set("txn_uuid", uuid)
	data.Add("peer", peer)
	data.Add("seconds",seconds)
	data.Add("nanos",nanos)
	data.Add("height",height)
	data.Add("currentBlockHash",currentBlockHash)
	data.Add("previousBlockHash",previousBlockHash)

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u) // "https://api.com/user/"

	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("\n")
	fmt.Println("Message was successfully sent", string(body))
}

func getBlockInfo (peer string) (height string, currentBlockHash string,previousBlockHash string) {
	peerRESTAddress := getPeerRESTAddress(peer)

	resp, err := http.Get(peerRESTAddress+"/chain")
	if err != nil {
		fmt.Println("Error : ",err.Error());
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error : ",err.Error());
	}

	var blockinfo blockInfo
	err = json.Unmarshal(body, &blockinfo)
	if err != nil {
		fmt.Println("Error : ",err.Error());
	}

	return strconv.Itoa(blockinfo.Height), blockinfo.CurrentBlockHash, blockinfo.PreviousBlockHash
}

func getPeerRESTAddress(peer string) string{
	var host string

	if (strings.Compare(peer,"jdoe")==0){
		host = "http://0.0.0.0:5000"
	}else{
		peerID := strings.Replace(peer,"vp","",-1)
		peerIDint,_ := strconv.Atoi(peerID)
		peerIDint = peerIDint+2

		peerID = strconv.Itoa(peerIDint)
		host = "http://172.17.0."+peerID+":5000"
	}
	return host
}

func getPeerID(eventAddress string) string{
	var peerID string

	i := strings.Index(eventAddress, ":")
	address := eventAddress[:i]
	addressArr := strings.Split(address, ".")

	switch addressArr[3] {
	case "2":
		peerID = "vp0"
	case "3":
		peerID = "vp1"
	case "4":
		peerID = "vp2"
	case "5":
		peerID = "vp3"
	default: peerID = "jdoe"
	}
	return peerID
}


func main() {

	var eventAddress string
	var sendEvents bool


	flag.StringVar(&eventAddress, "events-address", "0.0.0.0:31315", "address of events server")
	flag.StringVar(&chaincodeName, "chaincode-name", "anz", "chaincode name")
	flag.BoolVar(&sendEvents, "send-events",false,"should the events be communicated outside")

	flag.Parse()

	fmt.Printf("Event Address: %s\n", eventAddress)
	peer = getPeerID(eventAddress)

	a := createEventClient(eventAddress)
	if a == nil {
		fmt.Printf("Error creating event client\n")
		return
	}

	for {
		select {
		case c := <-a.chaincode:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Chaincode Event Received\n")
			fmt.Printf("--------------\n")
			fmt.Printf("Transaction :\n\t[%v]\n", c.ChaincodeEvent.String())

		case b := <-a.notify:
			fmt.Printf("\n")
			fmt.Printf("\n")
			fmt.Printf("Block Event Received\n")
			fmt.Printf("--------------\n")
			for _, r := range b.Block.Transactions {
				fmt.Printf("Transaction ID:\n\t[%v]\n", r.ChaincodeID)
				fmt.Printf("Transaction Type:\n\t[%v]\n", r.Type)
				fmt.Printf("Transaction UUID:\n\t[%v]\n", r.Uuid)
				fmt.Printf("Transaction Timestamp:\n\t[%v]\n", r.Timestamp)
				fmt.Printf("Transaction Payload:\n\t[%v]\n", r.Payload)
				if (sendEvents){
					sendRequest(eventAddress,
						r.Uuid,
						strconv.FormatInt(r.Timestamp.Seconds, 10),
						strconv.FormatInt(int64(r.Timestamp.Nanos), 10))
				}


			}




		}
	}
}
