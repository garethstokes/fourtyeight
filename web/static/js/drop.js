// a particle effect that emulates rain
rainEmitter = gamvas.ParticleEmitter.extend({
    // overwrite constructor
    create: function(name, x, y, img, anim) {
        // position our emitter slightly over the screen
        var dim = gamvas.getCanvasDimension();
        this._super(name, 0, -dim.h*0.5-20, img, anim);
        var st = gamvas.state.getCurrentState();

        // get our rain image and set its center to the bottom center
        var rainimg = new gamvas.Image(st.resource.getImage('media/drop.png'));
        rainimg.setCenter(1, 30);
        this.setImage(rainimg);

        // our raindrops orient themself along their path
        this.aligntoPath = true;

        // start with 100 drops per second
        this.setParticleRate(4);

        // let the particles start from a area slightly wider then the screen
        // range is always +/- from the center, so we use 0.6 time screen width
        this.setParticleStartPositionRange(new gamvas.Vector2D(dim.w*0.6, 0));

        // set the alpha of our raindrops to 0.5 over whole lifetime
        this.setAlphaTable([ [0.0, 0.5], [1.0, 0.8] ]);

        // raindrops do fall fast
        this.setParticleSpeed(dim.h*1.5);
        this.setParticleSpeedRange(dim.h*0.3);

        // the lifetime is something you have to match to your speed and
        // speed range by try and error, so they and up over the water
        this.setParticleLifeTime(0.45);
        this.setParticleLifeTimeRange(0.06);

        // emulate a bit of wind by turning the emitter slightly
        this.setRotation(0.15);

        // initialize the stuff we need for the splash emitters
        this.splashCounter = 0;
        this.splashImage = st.resource.getImage('media/splash.png');
        this.splashes = [];
    },

    draw: function(t) {
        // draw our spashes and delete them if finished
      var newsplashes = [];
      for (var i = 0; i < this.splashes.length; i++) {
        this.splashes[i].draw(t);
        if (this.splashes[i].lifeTime < 0.5) {
          newsplashes.push(this.splashes[i]);
        } else {
          delete this.splashes[i].draw(t);
        }
      }

      this.splashes = newsplashes;

      // call the super draing function to get the actual
      // particles drawn
      this._super(t);
    },

    onParticleEnd: function(pos, rot, scale, vel) {
        // if a particle dies, create a single particle emitter with a splash image
        var spi = new gamvas.Image(this.splashImage);
        spi.setCenter(16, 16);
        var sp = new gamvas.ParticleEmitter('splash'+this.splashCounter, pos.x, pos.y, spi);
        // only one splash
        sp.setParticleLimit(1);
        // fade in quickly to 0.2, then at a 1/10 of lifetime to 0.5
        // and the rest (9/10) slowly fade out
        sp.setAlphaTable([ [0.0, 0.0], [0.01, 0.2], [0.1, 0.5], [1.0, 0.0] ]);
        // scale quickly to 0.5 and then slow down until 1.0
        sp.setScaleTable([ [0.0, 0.2], [0.3, 0.5], [1.0, 1.0] ]);
        // stay in place
        sp.setParticleSpeed(0);
        sp.setParticleLifeTime(0.4);
        // save it for drawing
        this.splashes.push(sp);
        this.splashCounter++;
    }
});

rainState = gamvas.State.extend({
	init: function() {
        // create the emitter
        this.emitter = new rainEmitter('rain');
        this.addActor(this.emitter);
        this.dim = gamvas.getCanvasDimension();

        // sets the background color
        this.clearColor = "#323232";

        // disable screen clearing, as we draw a fullscreen image anyway
        this.autoClear = true;

        // get the background
        //this.bg = new gamvas.Image(this.resource.getImage('media/rainbg.jpeg'));
	},

  preDraw: function(t) {
  },

  postDraw: function(t) {
      // draw help after camera was applied
      //drawNextEffect(this.c);
  }, 
  
  loading: function(t) {
    var d = gamvas.getCanvasDimension();
    var tp = (d.h/2)-5;
    var w = parseInt(d.w*0.7, 10);
    var off = parseInt((d.w-w)/2, 10);
  }

});


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
