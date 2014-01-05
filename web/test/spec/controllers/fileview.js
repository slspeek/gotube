'use strict';

describe('Controller: FileviewCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var FileviewCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope) {
    scope = $rootScope.$new();
    FileviewCtrl = $controller('FileViewCtrl', {
      $scope: scope
    });
  }));

  it('should attach a list of awesomeThings to the scope', function () {
  });
});
