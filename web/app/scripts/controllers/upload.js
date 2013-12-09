(function() {
  'use strict';

  angular.module('webApp')
    .controller('UploadCtrl', function($location, $scope, $routeParams) {
      $scope.VideoId = $routeParams.VideoId;
      $scope.options = function() {
        console.log('Options called');
        return {
          target: '/upload/' + $scope.VideoId
        };
      };
      $scope.success = function() {
        console.log('it reached UploadCtrl');
        $location.replace();
        $location.path('list');
      };
    });

})();
