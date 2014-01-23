'use strict';

angular.module('webApp')
  .factory('authInfoLoader', function($q, $http) {
    return function() {
      var delay = $q.defer();
      $http.get('/auth').success(function(data) {
        delay.resolve(data.username);
      }).error(function() {
        delay.resolve('');
      });
      return delay.promise;
    };
  });
