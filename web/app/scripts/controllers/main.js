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
      console.log('username: ' + ahttp.username);
      $scope.videoId = Video.save({
        'Owner': ahttp.username,
        'Name': $scope.title,
        'Desc': $scope.description
      });
    };
    $scope.options = function() {
      console.log('Options called');
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
