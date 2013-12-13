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
      principal: {
        identity: function() {
          return {
            name: function() {
              return 'steven';
            }
          };
        },
        isAuthenticated: function() {
          return true;
        }
      }
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  

});
