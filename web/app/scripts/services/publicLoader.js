'use strict';

angular.module('webApp')
  .factory('publicLoader', function(PublicResource, $q) {
    return function() {
      var delay = $q.defer();
      PublicResource.getAll({}, function(videos) {
        delay.resolve(videos);
      }, function() {
        delay.reject('Unable to fetch videos ');

      });
      return delay.promise;
    };
  });
