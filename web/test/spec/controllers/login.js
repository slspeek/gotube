'use strict';

describe('Controller: LoginFormCtrl', function () {

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

});
