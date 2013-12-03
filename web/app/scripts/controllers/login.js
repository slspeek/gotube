(function() {
  'use strict';

  angular.module('webApp')
    .controller('LoginFormCtrl', function($scope, ahttp) {
      $scope.username = 'steven';
      $scope.password = 'gnu';
      $scope.submit = function() {
        ahttp.setPrincipal($scope.username, $scope.password);
        ahttp.http.post('auth', {
          headers: {
            'Content-Type': 'application/json'
          }
        }).then(function(response) {
          ahttp.loginConfirmed();
          console.log('success', response.data);
        }, function(response) {
          console.log('error', response);
        });
      };
    });

})();
