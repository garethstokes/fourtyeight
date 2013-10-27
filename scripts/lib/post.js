var http = require('http');

module.exports = function(options, data, callback) {

  var userString = JSON.stringify(data);

  var headers = {
    'Content-Type': 'application/json',
    'Content-Length': userString.length
  };

  options.headers = headers;

  // Setup the request.  The options parameter is
  // the object we defined above.
  var req = http.request(options, function(res) {
    res.setEncoding('utf-8');

    var responseString = '';

    res.on('data', function(data) {
      responseString += data;
    });

    res.on('end', function() {
      var response = JSON.parse(responseString);
      if (response.ok === false) {
        console.log('error');
        console.log(response);
        return;
      }
      callback(response.result);
    });
  });

  req.on('error', function(e) {
    // TODO: handle error.
  });

  req.write(userString);
  req.end();
};
