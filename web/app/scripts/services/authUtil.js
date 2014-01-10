(function() {
  'use strict';
  angular.module('webApp').factory('authUtil',
    function($rootScope, $location, $http, authService, authority, principal) {
      var service = {
        presentLogin: function() {
          var path = $location.path();
          if (path !== '/login') {
            this.storedPath = $location.path();
          } else {
            return;
          }
          $location.path('/login');
        },
        goBack: function() {
          $location.path(this.storedPath);
        },
        checkAuth: function() {
          if (!principal.isAuthenticated()) {
            $rootScope.$broadcast('event:auth-loginRequired');
          } else {
            return principal.identity().name();
          }
        },
        login: function(username, password) {
          var xmlHttp = new XMLHttpRequest();
          xmlHttp.open('GET', 'auth', false, username, password);
          xmlHttp.send(null);
          $http.get('auth', {
            headers: {
              'Content-Type': 'application/json'
            }
          }).then(function(response) {
            authority.authorize(response.data);
            authService.loginConfirmed();
          }, function(response) {
            console.log('error', response);
          });
        }
      };
      $rootScope.$on('event:auth-loginRequired', function() {
        service.presentLogin();
      });
      $rootScope.$on('event:auth-loginConfirmed', function() {
        service.goBack();
      });
      return service;
    });

})();
