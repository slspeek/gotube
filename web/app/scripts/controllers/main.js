'use strict';

angular.module('webApp')
  .controller('MainCtrl', function($scope, $resource, ahttp) {
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
    $scope.options = function() {
      return {
        headers: ahttp.header(),
        target: '/upload'
      };
    };
  })
  .controller('LoginFormCtrl', function($scope, ahttp) {
    $scope.submit = function() {
      ahttp.setPrincipal($scope.username, $scope.password);
      ahttp.http.post('auth', {
        headers: {
          'Content-Type': 'application/json'
        }
      }).then(function(response) {
        ahttp.loginConfirmed();
        console.log('success', response.data);
      }, function(response) {
        console.log('error', response);
      });
    };
  });
