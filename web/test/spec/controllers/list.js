'use strict';

describe('Controller: ListCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, httpBackend, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend, _Page_) {
    Page = _Page_;
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos').respond([]);
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
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

  it('should add a list of videos to the scope', function() {
    httpBackend.flush();
    expect(scope.videoList.length).toBe(0);
  });

  it('should add username to the scope', function() {
    httpBackend.flush();
    expect(scope.username).toBe('steven');
  });

  it('should set the title to Listing', function() {
    expect(Page.title()).toBe('Listing');
    httpBackend.flush();
  });
});

describe('Controller: ListCtrl, unitialized principal', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, httpBackend, mockPath;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend) {
    httpBackend = $httpBackend;
    httpBackend.when('GET', '/api/videos').respond([]);
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
      $scope: scope,
      $rootScope: {
        $broadcast: function(event) {
          mockPath = event;
        }
      },
      principal: {
        isAuthenticated: function() {
          return false;
        }
      }
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
    expect(mockPath).toBe('event:auth-loginRequired');
  });
});
