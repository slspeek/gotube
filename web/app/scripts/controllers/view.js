(function() {
  'use strict';

  angular.module('webApp')
    .controller('ViewCtrl', function($scope, $routeParams, $sce, VideoResource, ahttp) {
      $scope.id = $routeParams.VideoId;
      var video = VideoResource.get({
        Id: $scope.id
      }, function() {
        $scope.name = video.Name;
        $scope.desc = video.Desc;
      });

      $scope.videoURL = $sce.trustAsResourceUrl('/serve/' + $routeParams.BlobId + '.ogv');


    });

})();
