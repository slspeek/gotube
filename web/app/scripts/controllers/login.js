(function() {
  'use strict';

  angular.module('webApp')
    .controller('LoginFormCtrl', function($scope, $location, $http, authService, authority) {
      $scope.submit = function() {
        var xmlHttp = new XMLHttpRequest();
        xmlHttp.open( 'GET', 'auth', false, $scope.username, $scope.password);
        xmlHttp.send( null );
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
      };
    });

})();
