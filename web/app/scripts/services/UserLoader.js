'use strict';

angular.module('webApp')
  .service('userLoader', function($q, $http, $rootScope) {
    return function() {
      var delay = $q.defer();
      $http.get('/auth').success(function(data) {
        delay.resolve(data.username);
      }).error(function() {
        delay.reject('Unable to fetch username, fired a loginRequired event.');
        $rootScope.$broadcast('auth-loginRequired');
      });
      return delay.promise;
    };
  });
