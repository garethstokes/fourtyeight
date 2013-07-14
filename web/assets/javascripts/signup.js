
window.drop = function() {
  "use strict";
  /*global $ */
  /*jshint browser:true */

  function center(element) {
    $(element).css({
      "top": "50%", 
      "left": "50%",
      "margin-top": "-150px",
      "margin-left": "-" + ($(element).width() /2) + "px"
    });
  }

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

  function fadeMesageOut() {
    $('#message #form').fadeOut(500);

    setTimeout( function() {
      $('.logo').animate({
        top: '+=150px'
      }, 200, function() {
        // swap it with the splash
        $(this).css({
          width: "50px",
          top: "200px",
          left: "0px"
        });

        this.src = 'media/splash.png?time=' + new Date().getTime();

        $(this).animate({ 
          width: "300px", 
          opacity: 0, 
          top: "-=100px",
          left: "-=100px"
        }, 500);
      });
    }, 800);
  }

  function error( message ) {
     $('#error').css({
       "margin": "8px 0px 8px 32px",
       "width": "180px"
     }).html( message );
  }

  function showMessage() {
    
    //$('#message').css({ "display": "block" });
    center('#message');

    $('#message').fadeIn(400);
    $('#sign-in-with-facebook').mousedown(function() {
       $(this).addClass('fb-down');
    }).mouseup(function() {
       $(this).removeClass('fb-down');
    });
  
    var inputGreeting = 'Give us your email...';
  
    $('input').focus();
    $('input').val('Give us your email...');
    $('input').keydown(function() {
      if (this.value === inputGreeting) {
        this.value = '';
      }

    }).focusout(function() {
      if (this.value === '') {
        this.value = inputGreeting;
      }
    }).focusout(function() {
      $('form').submit();
    });

    $('form').submit(function() {
      //var value = $(this).serialize();

      var value = $('#email').val();

      if (value === inputGreeting) { return false; }

      if( value.indexOf('@') === -1 ) {
       error( 'that email is not valid.' );
       return false;
      }

      $.ajax({
        method: "post",
        url: $(this).attr('action'), //sumbits it to the given url of the form
        data: $(this).serialize()
      }).success(function(json){
        //act on result.

        if (json.ok === false) {
          error( json.result );
          return false;
        }

        fadeMesageOut();
     
        setTimeout(function() {
          $('#message').remove();
          showThankyou();
        }, 2000);

        store.persist( value );

      }).error(function(model, status) {
        $('.error').html(status);
      });

     return false;
   });
  }


  function showThankyou() {
    center('#thank-you');
    $('#thank-you').css({ "margin-left": "+=45" });
    $('#thank-you').fadeIn(500);
    $('#thank-you').append('<br /><br /><input type="button" value="clear local storage"></input>')
                   .click(function() { store.clear(); window.location = "rain.html"; });
  }

  if (typeof store.get() === 'undefined') {
    showMessage();
  } else {
    showThankyou();
  }
};
