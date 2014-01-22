'use strict';

angular.module('webApp')
  .factory('publicVideo',
    function(PublicResource, $route, $q) {
      return function() {
        var delay = $q.defer();
        PublicResource.get({
          Id: $route.current.params.VideoId
        }, function(video) {
          delay.resolve(video);
        }, function() {
          delay.reject('Unable to fetch video ' + $route.current.params.VideoId);

        });
        return delay.promise;
      };
    }
);
