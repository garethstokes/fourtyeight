var post = require('./post');

module.exports = function(callback) {

  var user = {
    name: 'caveman',
    password: 'bobafett'
  };

  var options = {
    host: 'localhost',
    port: 8000,
    path: '/user/login',
    method: 'POST'
  };

  post(options, user, callback);

};
