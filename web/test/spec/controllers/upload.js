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
      principal: {
        identity: function() {
          return {
            name: function() {
              return 'steven';
            }
          };
        }
      }
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  it('should produce options function ', function () {
    var opts = scope.options();
    expect(opts.target).toBe('/upload/128');

  });

  it('should copy VideoId from routeParam ', function () {
    expect(scope.VideoId).toBe('128');
  });
});
