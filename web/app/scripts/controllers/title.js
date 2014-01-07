'use strict';

angular.module('webApp')
  .controller('TitleCtrl', function ($scope, Page) {
    $scope.Page = Page;
  });
