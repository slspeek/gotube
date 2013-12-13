'use strict';

(function() {
  angular.module('webApp', [
    'ngCookies',
    'ngRoute',
    'ngFlow',
    'ngSanitize',
    'ngResource',
    'http-auth-interceptor',
    'authentication',
    'base64',
    'com.2fdevs.videogular',
    'com.2fdevs.videogular.plugins.controls',
    'com.2fdevs.videogular.plugins.overlayplay',
    'com.2fdevs.videogular.plugins.buffering',
    'com.2fdevs.videogular.plugins.poster'
  ]).run(function($rootScope, authUtil) {
    $rootScope.obj = {
      flow: ''
    };
    console.log('Auth service injected: ' + authUtil);
  })

  .config(function($routeProvider) {
    $routeProvider
      .when('/upload', {
        templateUrl: 'views/upload.html',
        controller: 'UploadCtrl'
      })
      .when('/view/:VideoId/:BlobId', {
        templateUrl: 'views/view.html',
        controller: 'ViewCtrl'
      })
      .when('/list', {
        templateUrl: 'views/list.html',
        controller: 'ListCtrl'
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginFormCtrl'
      })
      .otherwise({
        redirectTo: '/list'
      });
  }).config(function(flowFactoryProvider) {
    flowFactoryProvider.defaults = {
      maxChunkRetries: 2,
      testChunks: true,
      chunkRetryInterval: 5000,
      simultaneousUploads: 4
    };
  });
})();
