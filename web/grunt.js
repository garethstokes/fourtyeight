/*jshint node:true */

module.exports = function(grunt) {
  "use strict";

  grunt.initConfig({

    meta: {
      banner: '/*  {{ the drop }}\n' +
              ' *  version 0.01, <%= grunt.template.today("yyyy-mm-dd") %>\n' +
              ' *\n' +
              ' *  By Dave Cave (caveman), Nick Skelton (shredder) and Gareth Stokes (garrydanger) */\n',
      commonFiles: [
        'assets/javascripts/rainemitter.js', 
        'assets/javascripts/rainstate.js', 
        'assets/javascripts/signup.js'
      ]
    },

    lint: {
      files: ['assets/javascripts/signup.js']
    },

    jshint: {
      options: {
        curly: false,

        eqeqeq: true,
        forin: true,
        immed: true,
        latedef: true,
        newcap: true,
        noarg: true,
        noempty: true,
        nonew: true,
        plusplus: true,
        sub: true,
        undef: true,
        unused: true,
        boss: true,
        eqnull: true,
        browser: true
      },
      globals: {}
    },

    clean: {
      build: 'static/js/drop.js'
    },
  
    concat: {
      all: {
        options: { stripBanners: true },
        src: [
          '<config:meta.commonFiles>'
        ],
        dest: 'static/js/drop.js'
      }
    },

    min: {
      all: {
        src: ['<banner:meta.banner>', '<config:concat.all.dest>'],
        dest: 'static/js/drop.min.js'
      }
    },

    wiki: {
      docs: {
        src: 'docs/',
        dest: 'docs/html/'
      }
    },

    reload: {
      port: 7777,
      proxy: {
        host: 'localhost',
        port: '8000'
      }
    },

    watch: {
      files: [
        '<config:lint.files>'
      ],
      tasks: 'default reload'
    }
  });

  grunt.loadNpmTasks('grunt-contrib-clean');
  grunt.loadNpmTasks('grunt-reload');

  // Tasks to use from command line

  grunt.registerTask('default', 'lint clean concat min');
  grunt.registerTask('w', 'server reload watch'); // you need to run reload before watch to start proxy server 
};

