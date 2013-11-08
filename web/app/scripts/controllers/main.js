'use strict';

angular.module('webApp')
  .controller('MainCtrl', function($scope, $resource) {
    var Video = $resource('/api/videos/:Id', {
      Id: '@id'
    });
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
    $scope.save = function() {
      $scope.videoId = Video.save({
        'title': $scope.title,
        'desc': $scope.description
      });
    };
    $scope.uploadConfig = {
      url: '/upload',
      data: {
        extradata: 123
        // Will contain the file or files when sent with the upload-button
      }
    };

    $scope.onSuccess = function(response) {
      console.log(response.data);
    };
    $scope.onError = function(response) {
      console.log(response.data);
    };
  });
