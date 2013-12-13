(function() {
  'use strict';

  angular.module('webApp')
    .controller('EditCtrl', function($scope, $location, principal, VideoResource) {
      $scope.save = function() {
        $scope.videoId = VideoResource.save({
          'Owner': principal.identity().name(),
          'Name': $scope.title,
          'Desc': $scope.description
        }, function(data) {
          console.log(data.Id);
          $location.path('/upload/' + data.Id);
        });
      };
      $scope.title = '';
      $scope.description = '';
    });

})();
