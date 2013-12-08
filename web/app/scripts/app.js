'use strict';

(function() {
  angular.module('webApp', [
    'ngCookies',
    'ngRoute',
    'ngFlow',
    'ngSanitize',
    'ngResource',
    'http-auth-interceptor',
    'ngBase64',
    'com.2fdevs.videogular',
    'com.2fdevs.videogular.plugins.controls',
    'com.2fdevs.videogular.plugins.overlayplay',
    'com.2fdevs.videogular.plugins.buffering',
    'com.2fdevs.videogular.plugins.poster'
  ]).run(function(ahttp) {

  })
    .config(function($routeProvider) {
      $routeProvider
        .when('/upload', {
          templateUrl: 'views/edit.html',
          controller: 'EditCtrl'
        })
        .when('/upload/:VideoId', {
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
        target: '/upload',
        permanentErrors: [401],
        maxChunkRetries: 2,
        testChunks: true,
        //withCredentials: true,
        chunkRetryInterval: 5000,
        simultaneousUploads: 4
      };
      flowFactoryProvider.on('catchAll', function(event) {
        console.log('catchAll', event.arguments);
      });
    });
})();
