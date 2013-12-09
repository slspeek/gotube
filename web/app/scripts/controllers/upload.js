(function() {
  'use strict';

  angular.module('webApp')
    .controller('UploadCtrl', function($scope, $routeParams, $resource, $location, ahttp, VideoResource) {
      $scope.VideoId = $routeParams.VideoId;
      $scope.save = function() {
        console.log('username: ' + ahttp.username);
        VideoResource.save({
          'Owner': ahttp.username,
          'Name': $scope.title,
          'Desc': $scope.description
        }, function(data) {
          console.log(data.Id);
          $scope.returnedId = data.Id;
          $location.path('/upload/' + data.Id);
        });

      };
      $scope.options = function() {
        console.log('Options called');
        return {
          headers: ahttp.header(),
          target: '/upload/' + $scope.VideoId
        };
      };
    });

})();
