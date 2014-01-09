'use strict';

describe('Controller: UploadCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var UploadCtrl,
    scope, httpBackend, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, $httpBackend, _Page_, userLoader) {
    Page = _Page_;
    scope = $rootScope.$new();
    httpBackend = $httpBackend;
    httpBackend.expect('GET', '/auth').respond({username: 'Misko'});
    UploadCtrl = $controller('UploadCtrl', {
      $scope: scope,
      UserName: userLoader()
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  it('should set the title to upload', function() {
    httpBackend.flush();
    expect(Page.title()).toBe('Upload');
  });

  

});
