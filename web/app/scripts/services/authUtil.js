(function() {
  'use strict';
  angular.module('webApp').factory('authUtil',
    function($rootScope, $location, $http, authService, authority, principal) {
      console.log('Auth service injected: ' + authService);
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
          $http.post('auth', {
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
