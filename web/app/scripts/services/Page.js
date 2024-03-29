'use strict';

angular.module('webApp').
factory('Page', function($rootScope, $location) {
  var title = 'Gotube';
  var onUpload = false;
  var onMyVideos = false;
  $rootScope.$on('$routeChangeSuccess', function() {
    var current =  $location.path();
    if (current === '/upload') {
      onUpload = true;
    } else {
      onUpload = false;
    }
    if (current === '/list') {
      onMyVideos = true;
    } else {
      onMyVideos = false;
    }
  });
  return {
    title: function() {
      return title;
    },
    setTitle: function(newTitle) {
      title = newTitle;
    },
    onUpload: function() {
      return onUpload;
    },
    onMyVideos: function() {
      return onMyVideos;
    }
  };
});
