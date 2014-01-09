'use strict';

describe('Service: authUtil', function() {

  beforeEach(function() {
    this.addMatchers({
      toEqualData: function(expected) {
        return angular.equals(this.actual, expected);
      }
    });
  });

  describe('method presentLogin', function() {


    beforeEach(module('webApp'));

    var authUtil, $location, $httpBackend;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $httpBackend = $injector.get('$httpBackend');
      $location = $injector.get('$location');
    }));

    afterEach(function() {
      $httpBackend.verifyNoOutstandingExpectation();
      $httpBackend.verifyNoOutstandingRequest();
    });


    it('should store the old path', function() {
      $location.path('/foo');
      authUtil.presentLogin();
      expect(authUtil.storedPath).toBe('/foo');
    });

    it('should redirect to /login', function() {
      $location.path('/foo');
      authUtil.presentLogin();
      expect($location.path()).toBe('/login');

    });
    it('should no op on path=/login', function() {
      $location.path('/login');
      authUtil.presentLogin();
      expect($location.path()).toBe('/login');
    });

  });


  describe('method goBack', function() {
    beforeEach(module('webApp'));

    var authUtil, $location;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $location = $injector.get('$location');
    }));


    it('should go back to the old path', function() {
      authUtil.storedPath = 'Bar';
      authUtil.goBack();
      expect(authUtil.storedPath).toBe('Bar');
    });


  });

  describe('binding to loginRequired', function() {

    beforeEach(module('webApp'));

    var authUtil, $rootScope, $location;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $location = $injector.get('$location');
      $rootScope = $injector.get('$rootScope');
    }));


    it('should go to login on loginRequired', function() {
      $location.path('/foo');
      $rootScope.$broadcast('event:auth-loginRequired');
      expect($location.path()).toBe('/login');
    });

  });
  describe('method checkAuth', function() {

    beforeEach(module('webApp', function($provide) {
      var principal = jasmine.createSpyObj('principal', ['isAuthenticated', 'identity']);
      principal.isAuthenticated.andReturn('true');
      principal.identity.andReturn({ name: function() { return 'Misko';}});
      $provide.value('principal', principal);
    }));


    var authUtil, $rootScope, $location;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $location = $injector.get('$location');
      $rootScope = $injector.get('$rootScope');
    }));


    it('should return username when authenticated', function() {
      expect(authUtil.checkAuth()).toBe('Misko');
    });

  });
  describe('binding to loginConfirmed', function() {

    beforeEach(module('webApp'));

    var authUtil, $rootScope, $location;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $location = $injector.get('$location');
      $rootScope = $injector.get('$rootScope');
    }));


    it('should go to login on loginRequired', function() {
      authUtil.storedPath = '/foo';
      $rootScope.$broadcast('event:auth-loginConfirmed');
      expect($location.path()).toBe('/foo');
    });

  });
  describe('method login', function() {

    beforeEach(module('webApp'));

    var authUtil, $rootScope, $location, $httpBackend;

    beforeEach(inject(function($injector) {
      authUtil = $injector.get('authUtil');
      $location = $injector.get('$location');
      $rootScope = $injector.get('$rootScope');
      $httpBackend = $injector.get('$httpBackend');


      $httpBackend.when('GET', 'auth').respond({
        username: 'foo'
      });
    }));


    it('should goback on login success', function() {
      authUtil.storedPath = '/foo';
      authUtil.login('foo', 'bar');
      $httpBackend.flush();
      expect($location.path()).toBe('/foo');
    });

  });

});
