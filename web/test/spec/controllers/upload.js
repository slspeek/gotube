'use strict';

describe('Controller: UploadCtrl', function() {

  describe('Initialization', function() {

    // load the controller's module
    beforeEach(module('webApp'));

    var UploadCtrl,
      scope, Page;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope,  _Page_) {
      Page = _Page_;
      scope = $rootScope.$new();
      UploadCtrl = $controller('UploadCtrl', {
        $scope: scope,
        UserName: 'Misko'
      });
    }));


    it('should set the title to upload', function() {
      expect(Page.title()).toBe('Upload');
    });
  });

  describe('method success', function() {

    // load the controller's module
    beforeEach(module('webApp'));

    var UploadCtrl,
      scope, $location;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope,  _$location_) {
      $location = _$location_;
      scope = $rootScope.$new();
      scope.videoId = {Id:'42'};
      UploadCtrl = $controller('UploadCtrl', {
        $scope: scope,
        UserName: 'Misko'
      });
    }));


    it('should set the $location to /view/42', function() {
      scope.success();
      expect($location.path()).toBe('/view/42');
    });
  });

  describe('method save', function() {

    // load the controller's module
    beforeEach(module('webApp'));

    var UploadCtrl,
      scope, httpBackend, Page;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, $httpBackend, _Page_, userLoader) {
      Page = _Page_;
      scope = $rootScope.$new();
      scope.obj = {
        flow : {
          opts: {},
          upload: function(){}
          }
          };
      spyOn(scope.obj.flow, 'upload'); 
      httpBackend = $httpBackend;
      httpBackend.expect('POST', '/api/videos').respond({
        Id: 42
      });
      UploadCtrl = $controller('UploadCtrl', {
        $scope: scope,
        UserName: 'Misko'
      });
    }));

    afterEach(function() {
      httpBackend.verifyNoOutstandingExpectation();
      httpBackend.verifyNoOutstandingRequest();
    });

    it('should call the server', function() {
      scope.save(); 
      httpBackend.flush();
      expect(scope.obj.flow.upload).toHaveBeenCalled();
    });
  });
  describe('method fileAdded', function() {

    // load the controller's module
    beforeEach(module('webApp'));

    var UploadCtrl,
      scope, httpBackend, Page;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, $httpBackend, _Page_, userLoader) {
      Page = _Page_;
      scope = $rootScope.$new();
      scope.obj = {
        flow : {
          opts: {},
          upload: function(){}
          }
          };
      spyOn(scope.obj.flow, 'upload'); 
      httpBackend = $httpBackend;
      httpBackend.expect('POST', '/api/videos').respond({
        Id: 42
      });
      UploadCtrl = $controller('UploadCtrl', {
        $scope: scope,
        UserName: 'Misko'
      });
    }));

    afterEach(function() {
      httpBackend.verifyNoOutstandingExpectation();
      httpBackend.verifyNoOutstandingRequest();
    });

    it('should start the upload', function() {
      scope.name = 'Something';
      scope.fileAdded({file: {name: "Foo"}}); 
      httpBackend.flush();
      expect(scope.obj.flow.upload).toHaveBeenCalled();
    });
    it('should set the name to the filename if empty', function() {
      scope.fileAdded({file: {name: "Foo"}}); 
      httpBackend.flush();
      expect(scope.obj.flow.upload).toHaveBeenCalled();
    });
  });



});
