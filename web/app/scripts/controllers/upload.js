(function() {
  'use strict';

  angular.module('webApp')
    .controller('UploadCtrl', function($rootScope, $location, $scope, principal, VideoResource, Page) {
      Page.setTitle('Upload');
      if (!principal.isAuthenticated()) {
        $rootScope.$broadcast('event:auth-loginRequired');
      }
      $scope.fileAdded = function(file) {
        $scope.filename = file.file.name;
        if ($scope.name === '') {
          $scope.name = $scope.filename;
        }
        $scope.save();
      };
      $scope.save = function() {
        $scope.videoId = VideoResource.save({
          'Owner': principal.identity().name(),
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
        $location.path('list');
      };
      $scope.setFlowTarget = function(target) {
        $scope.obj.flow.opts.target = target;
      };
    });

})();
