(function() {
  'use strict';

  angular.module('webApp')
    .controller('LoginFormCtrl', function($scope, authUtil) {
      $scope.username = 'steven';
      $scope.password = 'gnu';
      $scope.submit = function() {
        authUtil.login($scope.username, $scope.password);
      };
        
    });

})();
