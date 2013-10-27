var http  = require('http'),
    login = require('./lib/login'),
    get   = require('./lib/get'),
    conf  = require('./lib/config');

login(function(result) {
  console.log(result);

  var options = {
    host: conf.host,
    port: conf.port,
    path: '/library/' + result.token
  };

  get(options, function(library) {
    console.log(library);
  });
});
