(function() {
  'use strict';

  angular.module('webApp')
    .controller('ListCtrl', function($scope, UserName, VideoList, Page) {
      Page.setTitle('Listing');
      $scope.username = UserName;
      $scope.videoList = VideoList;
      $scope.loggedOn = function() {
        return $scope.username !== '';
      };
      $scope.viewUrl = function(id) {
        if ($scope.loggedOn()) {
          return '/#/view/' + id;
        } else {
          return '/#/publicView/' + id;
        }
      };
    });

})();
