'use strict';

angular.module('webApp')
  .controller('PublicListCtrl', function($scope, VideoList, Page) {
      Page.setTitle('Public Listing');
      $scope.username = 'Public';
      $scope.videoList = VideoList;
    });
