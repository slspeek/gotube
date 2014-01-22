
(function() {
  'use strict';
  angular.module('webApp').factory('PublicResource', function($resource) {
      return $resource('/public/api/videos/:Id', {
        Id: '@Id'
      }, { getAll: { method: 'GET', params: {}, isArray: true},
           });
    });

})();
