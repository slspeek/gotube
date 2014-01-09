'use strict';

describe('Controller: ListCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, httpBackend, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, userLoader, $httpBackend, _Page_) {
    Page = _Page_;
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos').respond([]);
    httpBackend.expect('GET', '/auth').respond({
      username: 'Misko'
    });
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
      $scope: scope,
      UserName: userLoader()
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
    var name;
    scope.username.then(function(answer) {
      name = answer;
    });
    httpBackend.flush();
    expect(name).toBe('Misko');
  });

  it('should set the title to Listing', function() {
    expect(Page.title()).toBe('Listing');
    httpBackend.flush();
  });
});
