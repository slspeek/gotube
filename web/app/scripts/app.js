'use strict';

var webApp = angular.module('webApp', [
  'ngCookies',
  'ngRoute',
  'ngFlow',
  'ngSanitize',
  'ngResource'
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
          permanentErrors: [401  ],
          maxChunkRetries: 2,
          testChunks: true,
          chunkRetryInterval: 5000,
          simultaneousUploads: 1
        };
        flowFactoryProvider.on('catchAll', function(event) {
          console.log('catchAll', arguments);
        });
        });
        // Can be used with different implementations of Flow.js
        //   // flowFactoryProvider.factory = fustyFlowFactory;
        //   }]);
