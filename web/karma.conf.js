// Karma configuration
// http://karma-runner.github.io/0.10/config/configuration-file.html
'use strict';

module.exports = function(config) {
  config.set({
    // base path, that will be used to resolve files and exclude
    basePath: '',

    // testing framework to use (jasmine/mocha/qunit/...)
    frameworks: ['jasmine'],

    // list of files / patterns to load in the browser
    files: [
      'app/bower_components/angular/angular.js',
      'app/bower_components/angular-mocks/angular-mocks.js',
      'app/bower_components/angular-cookies/angular-cookies.js',
      'app/bower_components/angular-sanitize/angular-sanitize.js',

      'app/bower_components/angular-resource/angular-resource.js',
      'app/bower_components/angular-route/angular-route.js',
      'app/bower_components/flow.js/src/flow.js',
      'app/bower_components/angular-digest-interceptor/digest.js',
      'app/bower_components/angular-localstorage/angular-local-storage.js',
      'app/bower_components/js-md5/js/md5.js',
      'app/bower_components/ng-flow/src/directives/btn.js',
      'app/bower_components/ng-flow/src/directives/drop.js',
      'app/bower_components/ng-flow/src/directives/img.js',
      'app/bower_components/ng-flow/src/directives/init.js',
      'app/bower_components/ng-flow/src/directives/transfers.js',
      'app/bower_components/ng-flow/src/ng-flow.js',
      'app/bower_components/ng-flow/src/provider.js',
      'app/scripts/*.js',
      'app/scripts/**/*.js',
      'test/mock/**/*.js',
      'test/spec/**/*.js'
    ],

    // list of files / patterns to exclude
    exclude: [],

    // web server port
    port: 8080,

    // level of logging
    // possible values: LOG_DISABLE || LOG_ERROR || LOG_WARN || LOG_INFO || LOG_DEBUG
    logLevel: config.LOG_INFO,


    // enable / disable watching file and executing tests whenever any file changes
    autoWatch: false,


    // Start these browsers, currently available:
    // - Chrome
    // - ChromeCanary
    // - Firefox
    // - Opera
    // - Safari (only Mac)
    // - PhantomJS
    // - IE (only Windows)
    browsers: ['Chrome'],


    // Continuous Integration mode
    // if true, it capture browsers, run tests and exit
    singleRun: false
  });
};
