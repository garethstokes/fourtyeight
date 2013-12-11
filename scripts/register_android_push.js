var http  = require('http'),
    login = require('./lib/login'),
    post  = require('./lib/post'),
    conf  = require('./lib/config');

/* 

curl -XPOST localhost:8080/library/U3lHibjx9f8TzaHyXmkLug==/document \
  -H "Content-Type: application/json" \
  -d '{ "imageUrl": "http://i.imgur.com/1dxWiPF.jpg", "text": "the iron throne" , "expiry":1000 }'
APA91bGZKZYSrcJHzhwXHqcDNqJKBlfUVYxVxfKslGc3Sv_ZqY00d8Wu7pphpL_Y7wPEFy8Y4l8tF5-QQKE0PhE8kM363MjqfIdDVCsE49kuSRCs8PRPiuWlXYaeqRYkOzdm11l4xIESNLOldOv7pqVR1vnWggX4Jg

*/

login(function(result) {
  console.log(result);

  var options = {
    host: conf.host,
    port: conf.port,
    path: '/apns/register/android',
    method: 'POST'
  };

  var message = { 
    "Token": result.token, 
    "DeviceToken": "APA91bGZKZYSrcJHzhwXHqcDNqJKBlfUVYxVxfKslGc3Sv_ZqY00d8Wu7pphpL_Y7wPEFy8Y4l8tF5-QQKE0PhE8kM363MjqfIdDVCsE49kuSRCs8PRPiuWlXYaeqRYkOzdm11l4xIESNLOldOv7pqVR1vnWggX4Jg" 
  };

  post(options, message, function(response) {
    console.log(response);
  });

});
