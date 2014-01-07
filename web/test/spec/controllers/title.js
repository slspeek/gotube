'use strict';

describe('Controller: TitleCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var TitleCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    TitleCtrl = $controller('TitleCtrl', {
      $scope: scope,
      Page: {}
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
    expect(!!scope.Page).toBe(true);
  });
});
