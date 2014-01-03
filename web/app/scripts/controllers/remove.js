'use strict';

angular.module('webApp')
  .controller('RemoveCtrl', function ($scope, Video, $timeout, $location) {
    $scope.video = Video;
    $scope.remove = function() {
      console.log('Video before delete' + $scope.video);
      $scope.video.$remove({}, function() {
        $scope.message = 'Video was removed';
        $timeout(function() {
          $location.path('/list');
          }, 1000);
      });
    };
    $scope.cancel = function() {
      $location.path('/list');
    };
  });
