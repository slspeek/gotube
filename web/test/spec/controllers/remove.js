'use strict';

describe('Controller: RemoveCtrl', function() {


  describe('Initialization', function() {

    beforeEach(module('webApp'));

    var RemoveCtrl,
      scope, Page;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, _Page_) {
      Page = _Page_;
      scope = $rootScope.$new();
      RemoveCtrl = $controller('RemoveCtrl', {
        $scope: scope,
        Video: {
          Name: 'Novecento'
        }
      });
    }));

    it('should attach video object to the scope', function() {
      expect(scope.video).toBeDefined();
    });

    it('should set the title to remove {name}', function() {
      expect(Page.title()).toBe('Remove Novecento?');
    });
  });

  describe('method remove with success on the server', function() {

    // load the controller's module
    beforeEach(module('webApp', function($provide) {
      $provide.value('$route', {
        current: {
          params: {
            VideoId: 27
          }
        }
      });
    }));

    var RemoveCtrl,
      scope, Page, $httpBackend, $timeout, $location;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, _Page_, _$httpBackend_, videoLoader, _$timeout_, _$location_) {
      Page = _Page_;
      $timeout = _$timeout_;
      $location = _$location_;

      $httpBackend = _$httpBackend_;
      $httpBackend.expect('GET', '/api/videos/27').respond({
        Id: 27,
        Owner: 'Misko',
        Name: 'Novecento',
        Desc: 'Italian classic'
      });
      var Video;
      var VideoPromise = videoLoader();
      VideoPromise.then(function(video) {
        Video = video;
      });
      $httpBackend.flush();
      $httpBackend.expect('DELETE', '/api/videos/27').respond(200, '');
      scope = $rootScope.$new();
      RemoveCtrl = $controller('RemoveCtrl', {
        $scope: scope,
        Video: Video
      });
    }));

    afterEach(function() {
      $httpBackend.verifyNoOutstandingExpectation();
      $httpBackend.verifyNoOutstandingRequest();
    });
    it('should call the server ', function() {

      scope.remove();
      $httpBackend.flush();
      expect(scope.message).toBe('Video was removed');

      $timeout.flush();
      expect($location.path()).toBe('/list');
    });


  });
  describe('method remove with failure on the server', function() {

    // load the controller's module
    beforeEach(module('webApp', function($provide) {
      $provide.value('$route', {
        current: {
          params: {
            VideoId: 27
          }
        }
      });
    }));

    var RemoveCtrl,
      scope, Page, $httpBackend;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, _Page_, _$httpBackend_, videoLoader) {
      Page = _Page_;

      $httpBackend = _$httpBackend_;
      $httpBackend.expect('GET', '/api/videos/27').respond({
        Id: 27,
        Owner: 'Misko',
        Name: 'Novecento',
        Desc: 'Italian classic'
      })
      var Video;
      var VideoPromise = videoLoader();
      VideoPromise.then(function(video) {
        Video = video;
      });
      $httpBackend.flush();
      $httpBackend.expect('DELETE', '/api/videos/27').respond(501, '');
      scope = $rootScope.$new();
      RemoveCtrl = $controller('RemoveCtrl', {
        $scope: scope,
        Video: Video
      });
    }));

    afterEach(function() {
      $httpBackend.verifyNoOutstandingExpectation();
      $httpBackend.verifyNoOutstandingRequest();
    });
    it('should call the server ', function() {

      scope.remove();
      $httpBackend.flush();
    });


  });

  describe('method cancel()', function() {

    beforeEach(module('webApp'));

    var RemoveCtrl,
      scope, $location;

    // Initialize the controller and a mock scope
    beforeEach(inject(function($controller, $rootScope, _$location_) {
      $location = _$location_;
      scope = $rootScope.$new();
      RemoveCtrl = $controller('RemoveCtrl', {
        $scope: scope,
        Video: {
          Name: 'Novecento'
        }
      });
    }));

    it('should set the location to /list', function() {
      scope.cancel()
      expect($location.path()).toBe('/list');
    });
  });

});
