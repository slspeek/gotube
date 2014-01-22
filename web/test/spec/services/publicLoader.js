'use strict';

describe('Service: publicLoader', function () {

  var pinokio = {
    Id: 1,
    Name: 'Pinokio'
  };

  beforeEach(function() {
    this.addMatchers({
      toEqualData: function(expected) {
        return angular.equals(this.actual, expected);
      }
    });
  });
  // load the service's module
  beforeEach(module('webApp'));

  // instantiate service
  var publicLoader, $httpBackend;
  beforeEach(inject(['$httpBackend', 'publicLoader', function ($b, $p) {
    publicLoader = $p;
    $httpBackend = $b;
    $httpBackend.when('GET', '/public/api/videos').respond([pinokio]);
  }]));

  afterEach(function() {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
  });

  it('should return a Videos promise', function() {
    var promise = publicLoader();
    var videos;
    promise.then(function(data) {
      videos = data;
    });
    $httpBackend.flush();
    expect(videos[0]).toEqualData(pinokio);
  });
  it('should do something', function () {
    expect(!!publicLoader).toBe(true);
  });

});
