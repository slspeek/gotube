'use strict';

angular.module('webApp')
  .factory('videoLoader', function(VideoResource, $route, $q) {
    return function() {
      var delay = $q.defer();
      VideoResource.get({
        Id: $route.current.params.VideoId
      }, function(video) {
        delay.resolve(video);
      }, function() {
        delay.reject('Unable to fetch video ' + $route.current.params.VideoId);

      });
      return delay.promise;
    };
  });
