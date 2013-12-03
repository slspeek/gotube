(function() {
  'use strict';

  angular.module('webApp')
    .controller('LoginFormCtrl', function($scope, $location, ahttp) {
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
          $location.path('/list');
          console.log('success', response.data);
        }, function(response) {
          console.log('error', response);
        });
      };
    });

})();
