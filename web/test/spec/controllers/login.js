'use strict';

describe('Controller: LoginFormCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var LoginFormCtrl,
    scope, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, _Page_) {
    Page = _Page_;
    scope = $rootScope.$new();
    LoginFormCtrl = $controller('LoginFormCtrl', {
      $scope: scope
    });
  }));

  it('should set the title to Login', function() {
    expect(Page.title()).toBe('Login');
  });

});
