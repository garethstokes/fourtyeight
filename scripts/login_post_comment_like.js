var http  = require('http'),
    login = require('./lib/login'),
    post  = require('./lib/post'),
    conf  = require('./lib/config');

/* 

curl -XPOST localhost:8080/library/U3lHibjx9f8TzaHyXmkLug==/document \
  -H "Content-Type: application/json" \
  -d '{ "imageUrl": "http://i.imgur.com/1dxWiPF.jpg", "text": "the iron throne" , "expiry":1000 }'

*/

login(function(result) {
  console.log(result);

  var options = {
    host: conf.host,
    port: conf.port,
    path: '/library/' + result.token + '/document',
    method: 'POST'
  };

  var message = { 
    "imageUrl": "http://i.imgur.com/1dxWiPF.jpg", 
    "expiry": 1000 
  };

  post(options, message, function(response) {
    console.log("created post OK");
    console.log(response);
 

    var commentOptions = {
      host: conf.host,
      port: conf.port,
      path: '/library/' + result.token + '/document/' + response.key + "/post",
      method: 'POST'
    }

    var commentMessage = { 
      "imageUrl": "http://i.imgur.com/1dxWiPF.jpg", 
      "text": "the iron throne" 
    }

    post(commentOptions, commentMessage, function(commentResponse){
      console.log("comment response:");
      console.log(commentResponse);

      var likeOptions = {
        host: conf.host,
        port: conf.port,
        path: '/library/' + result.token + '/document/' + response.key + "/like/1",
        method: 'POST'
      }

      post(likeOptions, null, function(likeResponse){
        console.log("like response:");
         console.log("");
        console.log(likeResponse.comments);
      });

    });

  });

});
