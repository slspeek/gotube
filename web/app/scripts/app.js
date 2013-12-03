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
        .when('/upload/:VideoId', {
          templateUrl: 'views/main.html',
          controller: 'MainCtrl'
        })
        .when('/view/:VideoId', {
          templateUrl: 'views/view.html',
          controller: 'ViewCtrl'
        })
        .when('/view', {
          templateUrl: 'views/list.html',
          controller: 'ListCtrl'
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
