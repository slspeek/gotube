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
  ]).run(function($rootScope, $location, authService) {
    console.log('Auth: ' + authService);
    var service = {
      presentLogin: function() {
        var path = $location.path();
        if (path !== '/login') {
          this.storedPath = $location.path();
        }
        console.log('Stored: ' + this.storedPath);
        $location.path('/login');
      },
      goBack: function() {
        console.log('Going back to: ' + this.storedPath);
        $location.path(this.storedPath);
      }
    };
    $rootScope.$on('event:auth-loginRequired', function() {
      console.log('required');
      service.presentLogin();
    });
    $rootScope.$on('event:auth-loginConfirmed', function() {
      console.log('confirmed');
      service.goBack();
    });
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
      maxChunkRetries: 2,
      testChunks: true,
      chunkRetryInterval: 5000,
      simultaneousUploads: 4
    };
    flowFactoryProvider.on('success', function(event) {
      console.log('catchAll', event);
      console.dir(event);
    });
  });
})();
