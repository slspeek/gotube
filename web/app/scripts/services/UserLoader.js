'use strict';

angular.module('webApp')
  .service('userLoader', function($q, $http, $rootScope) {
    return function() {
      var delay = $q.defer();
      $http.get('/auth', {config: {ignoreAuthModule:true}}).success(function(data) {
        delay.resolve(data.username);
      }).error(function() {
        console.log('Error entered');
        delay.reject('Unable to fetch username, fired a loginRequired event.');
        $rootScope.$broadcast('event:auth-loginRequired');
        console.log('after broeadcat');
      });
      return delay.promise;
    };
  });
