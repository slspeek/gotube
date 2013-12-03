'use strict';

describe('Controller: ListCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, httpBackend;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend) {
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos').respond([]);
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
      $scope: scope,
      ahttp: { username: 'steven'}
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  it('should add a list of videos to the scope', function() {
    httpBackend.flush();
    expect(scope.videoList.length).toBe(0);
  });

  it('should add username to the scope', function() {
    httpBackend.flush();
    expect(scope.username).toBe('steven');
  });
});
