var express = require('express');
var cfenv = require('cfenv');

var app = express();

app.use(express.static('public'));

app.get('/', function (req, res) {
   res.send('Blockchain POC');
})

// get the app environment from Cloud Foundry
var appEnv = cfenv.getAppEnv();

// start server on the specified port and binding host
app.listen(appEnv.port, appEnv.bind, function() {

	// print a message when the server starts listening
  console.log("server starting on " + appEnv.url);
});
