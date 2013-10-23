var http = require('http');

/* 
 * example options

  var options = {
    host: 'www.google.com',
    port: 80,
    path: '/index.html'
  };

*/

module.exports = function(options, callback) {
  
  http.get(options, function(res) {
    console.log("Got response: " + res.statusCode);

    var responseData = "";

    res.on("data", function(chunk) {
      responseData += chunk;
    });

    res.on("end", function() {
      var response = JSON.parse(responseData);
      if (response.ok === false) {
        console.log('error');
        console.log(result);
        return;
      }
      callback(response.result);
    });
  }).on('error', function(e) {
    console.log("Got error: " + e.message);
  });
};
