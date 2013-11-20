'use strict';

(function() {
  angular.module('webApp', [
    'ngCookies',
    'ngRoute',
    'ngFlow',
    'ngSanitize',
    'ngResource',
    'authentication'
  ])
    .config(function($routeProvider) {
      $routeProvider
        .when('/', {
          templateUrl: 'views/main.html',
          controller: 'MainCtrl'
        })
        .otherwise({
          redirectTo: '/'
        });
    }).config(function(flowFactoryProvider) {
      flowFactoryProvider.defaults = {
        target: '/upload',
        permanentErrors: [401],
        maxChunkRetries: 2,
        testChunks: true,
        chunkRetryInterval: 5000,
        simultaneousUploads: 4
      };
      flowFactoryProvider.on('catchAll', function(event) {
        console.log('catchAll', event.arguments);
      });
    });

})();
