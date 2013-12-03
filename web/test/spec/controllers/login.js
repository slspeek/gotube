'use strict';

describe('Controller: LginFormCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var LoginFormCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    LoginFormCtrl = $controller('LoginFormCtrl', {
      $scope: scope
    });
  }));

  it('should set username and password', function () {
    expect(scope.username).toBe('steven');
    expect(scope.password).toBe('gnu');
  });
});
