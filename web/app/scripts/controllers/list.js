(function() {
  'use strict';

  angular.module('webApp')
    .controller('ListCtrl', function($scope, UserName, VideoResource, Page) {
      Page.setTitle('Listing');
      $scope.username = UserName;
      $scope.videoList = VideoResource.getAll();
    });

})();
