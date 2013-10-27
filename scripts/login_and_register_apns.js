var http  = require('http'),
    login = require('./lib/login'),
    post  = require('./lib/post');

login(function(result) {
  console.log(result);

  var options = {
    host: 'localhost',
    port: 8000,
    path: '/apns/register',
    method: 'POST'
  };

  var message = { 
    token: result.token,
    deviceToken: "abcd"
  };

  post(options, message, function(response) {
    console.log(response);
  });

});
