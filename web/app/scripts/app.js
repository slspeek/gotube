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
      .when('/fileview', {
        templateUrl: 'views/view.html',
        controller: 'FileViewCtrl'
      })
      .when('/view/:VideoId/:BlobId', {
        templateUrl: 'views/view.html',
        controller: 'ViewCtrl',
        resolve: {
          Video: function(videoLoader) {
            return videoLoader();
          }
        }
      })
      .when('/list', {
        templateUrl: 'views/list.html',
        controller: 'ListCtrl'
      })
      .when('/login', {
        templateUrl: 'views/login.html',
        controller: 'LoginFormCtrl'
      })
      .when('/remove/:VideoId', {
        templateUrl: 'views/remove.html',
        controller: 'RemoveCtrl',
        resolve: {
          Video: function(videoLoader) {
            return videoLoader();
          }
        }
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
  }).config(function($httpProvider) {
    $httpProvider.defaults.headers.common['Do-Not-Challenge'] = 'True';
  });
})();
