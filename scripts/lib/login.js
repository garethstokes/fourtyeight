var post = require('./post'),
    conf = require('./config');

module.exports = function(callback) {

  var user = {
    name: 'caveman',
    password: 'bobafett'
  };

  var options = {
    host: conf.host,
    port: conf.port,
    path: '/user/login',
    method: 'POST'
  };

  post(options, user, callback);

};
