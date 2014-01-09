'use strict';

describe('Service: userLoader', function() {

  // load the service's module
  beforeEach(module('webApp'));


  describe('success path', function() {
    // instantiate service
    var userLoader;
    var $httpBackend;
    beforeEach(inject(function(_userLoader_) {
      userLoader = _userLoader_;
    }));



    beforeEach(inject(function($injector) {
      $httpBackend = $injector.get('$httpBackend');
      $httpBackend.when('GET', '/auth').respond({
        username: 'Misko'
      });
    }));

    afterEach(function() {
      $httpBackend.verifyNoOutstandingExpectation();
      $httpBackend.verifyNoOutstandingRequest();
    });

    it('should return a Username promise', function() {
      var promise = userLoader();
      var name;
      promise.then(function(data) {
        name = data;
      });
      $httpBackend.flush();
      expect(name).toBe('Misko');
    });

  });
  
  describe('failure path', function() {
    // instantiate service
    var userLoader;
    var $httpBackend;
    beforeEach(inject(function(_userLoader_) {
      userLoader = _userLoader_;
    }));



    beforeEach(inject(function($injector) {
      $httpBackend = $injector.get('$httpBackend');
      $httpBackend.expect('GET', '/auth').respond(403, '');
     
    }));

    afterEach(function() {
      $httpBackend.verifyNoOutstandingExpectation();
      $httpBackend.verifyNoOutstandingRequest();
    });

    it('should reject the promise', function() {
      var promise = userLoader();
      var name, reason_;
      promise.then(function(data) {
        name = data;
      },
      function(reason) {
      //debugger;
      reason_ = reason;
      });
      $httpBackend.flush();
      expect(name).toBeUndefined();
      expect(reason_).toBe('Unable to fetch username, fired a loginRequired event.');
    });

  });
});
