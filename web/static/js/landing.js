var application = function(configuration) {
  "use strict";
  /*jshint browser:true */
  /*global console */

  var submit = configuration.submit,
      form = configuration.form,
      result = configuration.result,
      fadeDuration = 100;

  var STORAGE_KEY = 'drop-registered';
  var store = {};
  if (typeof window.localStorage === 'object') {
    store = {
      get: function() {
        return window.localStorage[STORAGE_KEY];
      },
      persist: function(value) {
        window.localStorage[STORAGE_KEY] = value;
      },
      clear: function() {
        delete window.localStorage[STORAGE_KEY];
      }
    };
  }

  function error( message ) {
     $('#error').css({
       "margin": "8px 0px 8px 32px",
       "width": "180px"
     }).html( message );
  }

  submit.click(function() {
      //var value = $(this).serialize();

      $('#error').html('');
      var value = $('#email').val();

      //if (value === inputGreeting) { return false; }

      if( value.indexOf('@') === -1 ) {
       error( 'that email is not valid.' );
       return false;
      }

      $.ajax({
        method: "post",
        url: "/waitinglist", //sumbits it to the given url of the form
        data: { email: value }
      }).success(function(json){
        //act on result.

        if (json.ok === false) {
          error( json.result );
          return false;
        }

        form.fadeOut(fadeDuration, function() {
          result.fadeIn(fadeDuration);
          store.persist( value );
        });

      }).error(function(model, status) {
        $('.error').html(status);
      });

     return false;
  });

  submit.click(function() {
  });
};
