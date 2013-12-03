'use strict';

describe('Controller: ViewCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ViewCtrl,
    scope, httpBackend;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend) {
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos/345').respond({'Id':'345', 'Name':'Novecento', 'Desc':'Italian classic'});
    scope = $rootScope.$new();
    ViewCtrl = $controller('ViewCtrl', {
      $scope: scope,
      $routeParams: {VideoId:'345'}
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });


  it('should attach name to the scope', function() {
    httpBackend.flush();
    expect(scope.name).toBe('Novecento');
  });
  it('should attach desc to the scope', function() {
    httpBackend.flush();
    expect(scope.desc).toBe('Italian classic');
  });
  it('should attach name to the scope', function() {
    httpBackend.flush();
    expect(scope.id).toBe('345');
  });
});
