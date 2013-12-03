
(function() {
  'use strict';
  angular.module('webApp').factory('VideoResource', function($resource) {
      return $resource('/api/videos/:Id', {
        Id: '@id'
      }, { getAll: { method: 'GET', params: {}, isArray: true}});
    });

})();
