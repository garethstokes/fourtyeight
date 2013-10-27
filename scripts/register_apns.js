var http  = require('http'),
    login = require('./lib/login'),
    post  = require('./lib/post'),
    conf  = require('./lib/config');

// DONT LOG IN
(function() {

  var options = {
    host: conf.host,
    port: conf.port,
    path: '/apns/register',
    method: 'POST'
  };

  var message = { 
    token: "invalid_token",
    deviceToken: "abcd"
  };

  post(options, message, function(response) {
    console.log(response);
  });

})();
