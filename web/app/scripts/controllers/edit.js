(function() {
  'use strict';

  angular.module('webApp')
    .controller('EditCtrl', function($scope, $location, ahttp, VideoResource) {
      $scope.save = function() {
        $scope.videoId = VideoResource.save({
          'Owner': ahttp.username,
          'Name': $scope.title,
          'Desc': $scope.description
        }, function(data) {
          console.log(data.Id);
          $location.path('/upload/' + data.Id);
        });
      };
    });

})();
