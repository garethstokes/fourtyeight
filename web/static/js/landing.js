var application = function(configuration) {
  "use strict";
  /*jshint browser:true */
  /*global console */

  var submit = configuration.submit,
      form = configuration.form,
      result = configuration.result,
      fadeDuration = 100;

  submit.click(function() {
    form.fadeOut(fadeDuration, function() {
      result.fadeIn(fadeDuration);
    });
  });
};
