/*
- load page
- submit to fourtyeight
- submit to facebook
- save cookie
*/

window.startDrop = function(config) {
  "use strict";
  /*global $ */
  /*global setTimeout */

  function center(element) {
    $(element).css({
      "top": "50%", //(config.height /2) - 150,
      "left": "50%", //((config.width /2) - $(element).width() /2) + "px"
      "margin-top": "-150px",
      "margin-left": "-" + ($(element).width() /2) + "px"
    });
  }

  function fadeMesageOut() {
    $('#message #form').fadeOut(500);

    setTimeout( function() {
      $('.logo').animate({
        top: '+=150px'
      }, 100, function() {
        // swap it with the splash
        $(this).css({
          width: "50px",
          top: "200px",
          left: "0px"
        });

        this.src = 'media/splash.png';

        $(this).animate({ 
          width: "300px", 
          opacity: 0, 
          top: "-=100px",
          left: "-=100px"
        }, 500);
      });
    }, 800);
  }

  center('#message');

  $('#sign-in-with-facebook').mousedown(function() {
    $(this).addClass('fb-down');
  }).mouseup(function() {
    $(this).removeClass('fb-down');
  });

  var inputGreeting = 'Give us your email...';

  $('input').focus();
  $('input').val('Give us your email...');
  $('input').keydown(function() {
    if (this.value == inputGreeting) {
      this.value = '';
    }
    
  }).focusout(function() {
    if (this.value === '') {
      this.value = inputGreeting;
    }
  });

  $('form').submit(function() {
    var value = $(this).serialize();
    
    if( $('#email').val().indexOf('@') == -1 ) {
      $('#error').css({
        "margin": "8px 14px 8px 32px"
      }).html( 'that email is not valid.' );
      return false;
    }

    // fade out before saving to server.
    fadeMesageOut();

    setTimeout(function() {
      $('#message').remove();
      center('#thank-you');
      $('#thank-you').fadeIn(500);
    }, 2000);


    return false;
  });
  
};
