(function() {
  'use strict';

  angular.module('webApp')
    .controller('EditCtrl', function($scope, $location, UserName, Video, VideoResource, Page) {
      Page.setTitle('Edit');
      $scope.UserName = UserName;
      $scope.video = Video;
      $scope.save = function() {
        VideoResource.update($scope.video, function() {
          $location.path('/list');
        });
      };
    });

})();
