'use strict';

describe('Controller: UploadCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var UploadCtrl,
    scope, httpBackend;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, $httpBackend) {
    scope = $rootScope.$new();
    httpBackend = $httpBackend;
    UploadCtrl = $controller('UploadCtrl', {
      $scope: scope,
      $routeParams: {VideoId:'128'},
      ahttp: {header: function() {return "My header";}}
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  it('should produce options function ', function () {
    var opts = scope.options();
    expect(opts.headers).toBe('My header');
    expect(opts.target).toBe('/upload/128');

  });

  it('should save', function () {
    httpBackend.when('POST', '/api/videos').respond({Id: '124'});
    scope.title = 'My video';
    scope.save();
    httpBackend.flush();
    expect(scope.returnedId).toBe('124');
  });

  it('should copy VideoId from routeParam ', function () {
    expect(scope.VideoId).toBe('128');
  });
});
