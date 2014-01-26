(function() {
  'use strict';

  angular.module('webApp')
    .controller('UploadCtrl', function($rootScope, $location, $scope, VideoResource, Page, UserName) {
      Page.setTitle('Upload');
      $scope.username = UserName;
      $scope.fileAdded = function(file) {
        $scope.filename = file.file.name;
        if ($scope.name === '') {
          $scope.name = $scope.filename;
        }
        $scope.save();
      };
      $scope.save = function() {
        $scope.videoId = VideoResource.save({
          'Owner': $scope.username,
          'Name': $scope.name,
          'Desc': $scope.desc
        }, function(data) {
          $scope.setFlowTarget('/api/videos/' + data.Id + '/upload');
          $scope.obj.flow.upload();
        });
      };
      $scope.name = '';
      $scope.desc = '';
      $scope.success = function() {
        $location.replace();
        $location.path('/view/' + $scope.videoId.Id);
      };
      $scope.setFlowTarget = function(target) {
        $scope.obj.flow.opts.target = target;
      };
    });

})();
