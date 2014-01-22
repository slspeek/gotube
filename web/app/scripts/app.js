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
    'com.2fdevs.videogular.plugins.poster',
    'ui.bootstrap'
  ]).run(function($rootScope, authUtil) {
    $rootScope.obj = {
      flow: ''
    };
    (function() {}(authUtil));
  })

  .config(function($routeProvider) {
    $routeProvider
      .when('/upload', {
        templateUrl: 'views/upload.html',
        controller: 'UploadCtrl',
        resolve: {
          UserName: function(userLoader) {
            return userLoader();
          }
        }
      })
      .when('/fileview', {
        templateUrl: 'views/view.html',
        controller: 'FileViewCtrl'
      })
      .when('/edit/:VideoId', {
        templateUrl: 'views/edit.html',
        controller: 'EditCtrl',
        resolve: {
          Video: function(videoLoader) {
            return videoLoader();
          },
          UserName: function(userLoader) {
            return userLoader();
          }
        }
      })
      .when('/view/:VideoId', {
        templateUrl: 'views/view.html',
        controller: 'ViewCtrl',
        resolve: {
          Video: function(videoLoader) {
            return videoLoader();
          }
        }
      })
      .when('/public', {
        templateUrl: 'views/public.html',
        controller: 'PublicListCtrl',
        resolve: {
          VideoList: function(publicLoader) {
            return publicLoader();
          }
        }
      })
      .when('/list', {
        templateUrl: 'views/list.html',
        controller: 'ListCtrl',
        resolve: {
          UserName: function(userLoader) {
            return userLoader();
          }
        }
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
      .when('/publicView/:VideoId', {
        templateUrl: 'views/view.html',
        controller: 'PublicviewCtrl',
        resolve: {
          Video: function(publicVideo) {
            return publicVideo();
          }
        }

      })
      .otherwise({
        redirectTo: '/public'
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
