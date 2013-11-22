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
        .when('/', {
          templateUrl: 'views/main.html',
          controller: 'MainCtrl'
        })
        .when('/login', {
          templateUrl: 'views/login.html',
          controller: 'LoginFormCtrl'
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
        //withCredentials: true,
        //headers: {
        //'authorization': 'Xasic c3RldmVuOmdudQ=='
        //},
        chunkRetryInterval: 5000,
        simultaneousUploads: 4
      };
      flowFactoryProvider.on('catchAll', function(event) {
        console.log('catchAll', event.arguments);
      });
    })
    .service('ahttp', function($rootScope, $http, $location, $base64, authService) {
      var service = {
        code: '',
        header: function() {
          return {
            'authorization': this.code
          };
        },
        storedPath: '',
        presentLogin: function() {
          this.storedPath = $location.path();
          $location.path('/login');
        },
        goBack: function() {
          $location.path(this.storedPath);
        },
        http: $http,
        setPrincipal: function(username, password) {
          this.code = 'Xasic ' + $base64.encode(username + ':' + password);
          $http.defaults.headers.common.authorization = this.code;
        },
        loginConfirmed: function() {
          authService.loginConfirmed();
        },
      };
      $rootScope.$on('event:auth-loginRequired', function() {
        console.log('required');
        service.presentLogin();
      });
      $rootScope.$on('event:auth-loginConfirmed', function() {
        console.log('confirmed');
        service.goBack();
      });
      return service;
    });

})();
