
(function() {
  'use strict';
  angular.module('webApp').factory('VideoResource', function($resource) {
      return $resource('/api/videos/:Id', {
        Id: '@Id'
      }, { getAll: { method: 'GET', params: {}, isArray: true},
           update: { method: 'PUT', params: {Id: '@Id'}}
           });
    });

})();
