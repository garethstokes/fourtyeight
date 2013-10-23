var http  = require('http'),
    login = require('./lib/login'),
    get   = require('./lib/get');

login(function(result) {
  console.log(result);

  var options = {
    host: 'localhost',
    port: 8000,
    path: '/library/' + result.token
  };

  get(options, function(library) {
    console.log(library);
  });
});
