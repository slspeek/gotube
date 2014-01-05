'use strict';

describe('Controller: ViewCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ViewCtrl,
    scope;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope) {
    scope = $rootScope.$new();
    ViewCtrl = $controller('ViewCtrl', {
      $scope: scope,
      $routeParams: {VideoId:'345'},
      Video: {'Id':'345', 'Name':'Novecento', 'Desc':'Italian classic'}
    });
  }));



  it('should attach name to the scope', function() {
    expect(scope.name).toBe('Novecento');
  });
  it('should attach desc to the scope', function() {
    expect(scope.desc).toBe('Italian classic');
  });
  it('should attach name to the scope', function() {
    expect(scope.id).toBe('345');
  });
});
