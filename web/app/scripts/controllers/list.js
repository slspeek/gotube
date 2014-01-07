(function() {
  'use strict';

  angular.module('webApp')
    .controller('ListCtrl', function($rootScope, $scope, VideoResource, principal, Page) {
      Page.setTitle('Listing');
      $scope.videoList = VideoResource.getAll();
      if (!principal.isAuthenticated()) {
        $rootScope.$broadcast('event:auth-loginRequired');
      } else {
        $scope.username = principal.identity().name();
      }

    });

})();
