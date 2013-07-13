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
