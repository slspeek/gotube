'use strict';

describe('Controller: ListCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, httpBackend, page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend, Page) {
    page = Page;
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos').respond([]);
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
      $scope: scope,
      UserName: 'Misko'
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
    expect(scope.username).toBe('Misko');
    httpBackend.flush();
  });

  it('should set the title to Listing', function() {
    expect(page.title()).toBe('Listing');
    httpBackend.flush();
  });
});
