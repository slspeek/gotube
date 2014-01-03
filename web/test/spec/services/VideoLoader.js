'use strict';

describe('Service: Videoloader', function() {

  beforeEach(function() {
    this.addMatchers({
      toEqualData: function(expected) {
        return angular.equals(this.actual, expected);
      }
    });
  });

  beforeEach(module('webApp', function($provide) {
    $provide.value('$route', { current: { params: { VideoId: 1}}});
  }));

  var Loader, $httpBackend;
  beforeEach(inject(function(videoLoader) {
    Loader = videoLoader;
  }));

  beforeEach(inject(function($injector) {
    $httpBackend = $injector.get('$httpBackend');
    $httpBackend.when('GET', '/api/videos/1').respond({
      Id: 1,
      Name: 'Pinokio'
    });
  }));

  afterEach(function() {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
  });

  it('should return a Video promise, determined by the route', function() {
    var promise = Loader();
    var video;
    promise.then(function(data) {
      video = data;
    });
    $httpBackend.flush();
    expect(video).toEqualData({ Id:1, Name:'Pinokio'});
  });
});
