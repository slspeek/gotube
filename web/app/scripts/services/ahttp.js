'use strict';

(function() {
  angular.module('webApp').service('ahttp', function($rootScope, $http, $location, $base64, authService) {
      var service = {
        header: function() {
          return {
            'authorization': this.code
          };
        },
        presentLogin: function() {
          this.storedPath = $location.path();
          console.log('Stored: ' + this.storedPath);
          $location.path('/login');
        },
        goBack: function() {
          console.log('Going back to: ' + this.storedPath);
          $location.path(this.storedPath);
        },
        http: $http,
        setPrincipal: function(username, password) {
          this.username = username;
          this.code = 'Basic ' + $base64.encode(username + ':' + password);
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
