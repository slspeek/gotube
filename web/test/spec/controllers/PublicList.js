'use strict';

describe('Controller: PublicListCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var PublicListCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    PublicListCtrl = $controller('PublicListCtrl', {
      $scope: scope,
      VideoList: []
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(scope.videoList).toEqual([]);
  });
});
