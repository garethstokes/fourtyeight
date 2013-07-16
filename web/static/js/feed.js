drop = (function() {
  "use strict";
  /*jshint browser:true */
  /*global console */

  var events = {};
  var containers = {
    feed: "feed-content"
  };

  function bind( post ) {
    console.log( post.mainPost );

    var template = '';
    
    // text
    if (typeof post.mainPost.imageUrl === 'undefined') {
      template = document.getElementById( containers.text );
    } else if ( typeof post.mainPost.text === 'undefined') {
      template = document.getElementById( containers.image );
    } else { // hybrid
      template = document.getElementById( containers.hybrid );
    }

    var clone = template.cloneNode(true);
    clone.id = '';

    var html = clone.outerHTML;
    html = html.replace( '{{feed-author}}', post.mainPost.ownerId );
    
    // image
    if (typeof post.mainPost.imageUrl !== 'undefined') {
      var img = "<img src='" + post.mainPost.imageUrl + "' />";
      html = html.replace( '{{feed-image}}', img );
    }

    // text
    if (typeof post.mainPost.text !== 'undefined') {
      var text = "<h3>" + post.mainPost.text + "</h3>";
      html = html.replace( '{{feed-text}}', text );
    }

    // comments
    var comments = '<img class="no-posts" src="media/no_posts_badge.png" alt="have your say">';
    html = html.replace( '{{feed-comments}}', comments );

    var content = document.getElementById( containers.content );
    content.innerHTML = content.innerHTML + html;
  }

  function on(name, cb) {
    if (typeof events[name] === 'undefined') {
      events[name] = [];
    }

    events[name].push(cb);
  }
  
  function fire(name, value) {
    var callbacks = events[name];
    for( var i = 0; i < callbacks.length; i = i + 1 ) {
      callbacks[i](value);
    }
  }

  return {
    on: on,
    fire: fire, 

    feed: function(config) {
      containers = config.containers;

      on('data', function(data) {
        for( var i = 0; i < data.length; i = i +1 ) {
          bind( data[i] );
        }
      });

      config.datasource();
    }
  };
}());
