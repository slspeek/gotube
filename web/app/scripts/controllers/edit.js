(function() {
  'use strict';

  angular.module('webApp')
    .controller('EditCtrl', function($scope, $location, UserName, Video, VideoResource, Page) {
      Page.setTitle('Edit');
      $scope.UserName = UserName;
      $scope.Video = Video;
      $scope.save = function() {
        VideoResource.update($scope.Video, function() {
          $location.path('/list');
        });
      };
    });

})();
