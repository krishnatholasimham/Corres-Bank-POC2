var nodes, edges, network;

var DIR = 'img/';
var EDGE_LENGTH_MAIN = 300;
var EDGE_LENGTH_SUB = 50;

var GLOBAL_BLOCKCHAIN_HEIGHT = 0;
var GLOBAL_BLOCKCHAIN_HASH="";

var GREEN = "#b3ffb3";
var RED = "#ffc299";

var vp0_state = {
    id: 'vp0_state',
    height : 0,
    lastTxnID : "",
    previousBlockHash : "",
    currentBlockHash : "",
    size: 150, 
    label: "Validating Peer 0 World State\n\nAwaiting Transactions.....\n", 
    color: GREEN, 
    shape: 'box', 
    font: {'face': 'monospace', 'align': 'left', 'size':'16'}
};

var vp1_state = {
    id: 'vp1_state',
    height : 0,
    lastTxnID : "",
    previousBlockHash : "",
    currentBlockHash : "",
    size: 150, 
    label: "Validating Peer 1 World State\n\nAwaiting Transactions.....\n", 
    color: GREEN, 
    shape: 'box', 
    font: {'face': 'monospace', 'align': 'left', 'size':'16'}
};

var vp2_state = {
    id: 'vp2_state',
    height : 0,
    lastTxnID : "",
    previousBlockHash : "",
    currentBlockHash : "",
    size: 150, 
    label: "Validating Peer 2 World State\n\nAwaiting Transactions.....\n", 
    color: GREEN, 
    shape: 'box', 
    font: {'face': 'monospace', 'align': 'left', 'size':'16'}
};
var vp3_state = {
    id: 'vp3_state',
    height : 0,
    lastTxnID : "",
    previousBlockHash : "",
    currentBlockHash : "",
    size: 150, 
    label: "Validating Peer 3 World State\n\nAwaiting Transactions.....\n", 
    color: GREEN, 
    shape: 'box', 
    font: {'face': 'monospace', 'align': 'left', 'size':'16'}
};

var peerStateNodes = new Array(vp0_state,vp1_state,vp2_state,vp3_state);
    
if(typeof(EventSource) !== "undefined") {
    var source = new EventSource("api/server.php");
    source.onmessage = function(event) {
        eventHandler(event.data);
        
    };
} else {
    document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
}

function draw() {
     nodes = new vis.DataSet();
            nodes.on('*', function () {
                //document.getElementById('nodes').innerHTML = JSON.stringify(nodes.get(), null, 4);
            });
            
      nodes.add({id: 1, size:'50', label: 'Validating Peer 0', image: DIR + 'peer.png', shape: 'image'});
      nodes.add({id: 2, size:'50', label: 'Validating Peer 1', image: DIR + 'peer.png', shape: 'image'});
      nodes.add({id: 3, size:'50', label: 'Validating Peer 2', image: DIR + 'peer.png', shape: 'image'});
      nodes.add({id: 4, size:'50', label: 'Validating Peer 3', image: DIR + 'peer.png', shape: 'image'});
      
      nodes.add(vp0_state);
      nodes.add(vp1_state);
      nodes.add(vp2_state);
      nodes.add(vp3_state);
      
    edges = new vis.DataSet();
            edges.on('*', function () {
                //document.getElementById('edges').innerHTML = JSON.stringify(edges.get(), null, 4);
            });
            
      edges.add({from: 1, to: 2, length: EDGE_LENGTH_MAIN, smooth:true});
      edges.add({from: 2, to: 3, length: EDGE_LENGTH_MAIN, smooth:true});
      edges.add({from: 3, to: 4, length: EDGE_LENGTH_MAIN, smooth:true});
      edges.add({from: 4, to: 1, length: EDGE_LENGTH_MAIN, smooth:true});
      
      edges.add({from: 1, to: vp0_state.id, length: 400, smooth:true});
      edges.add({from: 2, to: vp1_state.id, length: EDGE_LENGTH_MAIN, smooth:true});
      edges.add({from: 3, to: vp2_state.id, length: 400, smooth:true});
      edges.add({from: 4, to: vp3_state.id, length: EDGE_LENGTH_MAIN, smooth:true});
      
      // create a network
      var container = document.getElementById('mynetwork');
      var data = {
        nodes: nodes,
        edges: edges
      };
      var options = {
          layout:{
                    randomSeed:2
                }
            };
      network = new vis.Network(container, data, options);
}

function eventHandler(eventData){
    var decodedEventData = decodeEvent(eventData);
    if(decodedEventData[0] !== ""){
        var peerStateNode = getPeerStateNode(decodedEventData[0]);
        var peerStateLabel = constructPeerStateLabel(peerStateNode,decodedEventData);


        updatePeerStateNode(peerStateNode.id,peerStateLabel);
        updatePeerStateNodeData(peerStateNode,decodedEventData);

        for (var i = 0; i < peerStateNodes.length; i++) {

            var peerColour = getPeerStateNodeColour(peerStateNodes[i].currentBlockHash, peerStateNodes[i].height);
            updatePeerStateNodeColour(peerStateNodes[i].id,peerColour);
        }
    }
}

function updatePeerStateNodeData(peer, eventData){
    peer.lastTxnID = eventData[2];
    peer.height = parseInt(eventData[3]);
    peer.currentBlockHash = eventData[4];
    peer.previousBlockHash = eventData[5];
    
    if (GLOBAL_BLOCKCHAIN_HEIGHT < peer.height){
        GLOBAL_BLOCKCHAIN_HEIGHT = peer.height;
        GLOBAL_BLOCKCHAIN_HASH = peer.currentBlockHash;
    }
}

function updatePeerStateNode(id, newLabel) {
            try {
                nodes.update({
                    id: id,
                    label: newLabel,
                    font: {'face': 'monospace', 'align': 'left', 'size':'10'}
                });
            }
            catch (err) {
                alert(err);
            }
}

function updatePeerStateNodeColour(id, newColour) {
            
            try {
                nodes.update({
                    id: id,
                    color: newColour
                });
            }
            catch (err) {
                alert(err);
            }
}

function decodeEvent(data){
    var dataArray = data.split("|");
    return dataArray;
}

function constructPeerStateLabel(peer, eventData){
    var peerStateLabel;
    var headerString = "Validating Peer "+peer.id.charAt(2)+" World State\n";
    
    peerStateLabel = headerString.concat("\nTransaction ID : ",eventData[1]);
    peerStateLabel = peerStateLabel.concat("\nReceived at :",eventData[2]);
    
    peerStateLabel = peerStateLabel.concat("\n\nPrevious Block Hash : ",eventData[5]);
    peerStateLabel = peerStateLabel.concat("\nCurrent Block Hash : ",eventData[4]);
    peerStateLabel = peerStateLabel.concat("\nBlockchain Height : ",eventData[3]);
    return peerStateLabel;
}

function getPeerStateNode(peer){

    var peerStateNode;
    switch (peer){
        case 'vp0':
            peerStateNode = vp0_state;
            break;
            
        case 'vp1':
            peerStateNode = vp1_state;
            break; 
        
        case 'vp2':
            peerStateNode = vp2_state;
            break;
        
        case 'vp3':
            peerStateNode = vp3_state;
            break;
            
        default: peerStateNode = vp0_state;
    }
    return peerStateNode;
}

function getPeerStateNodeColour(peerHash, height){

    if (GLOBAL_BLOCKCHAIN_HEIGHT <= height){
        return GREEN;
    }
    else return RED;
}