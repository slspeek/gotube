'use strict';

angular.module('webApp').
factory('Page', function() {
  var title = 'Gotube';
  return {
    title: function() {
      return title;
    },
    setTitle: function(newTitle) {
      title = newTitle;
    }
  };
});
