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
  ])
    .config(function($routeProvider) {
      $routeProvider
        .when('/upload', {
          templateUrl: 'views/edit.html',
          controller: 'EditCtrl'
        })
        .when('/view/:VideoId', {
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
          redirectTo: '/login'
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
