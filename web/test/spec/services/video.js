'use strict';
describe('Service: VideoResource', function() {

  var $httpBackend;
  var video;
  // load the controller's module
  beforeEach(module('webApp'));

  beforeEach(inject(function($injector) {
    $httpBackend = $injector.get('$httpBackend');
    // backend definition common for all tests
    $httpBackend.when('GET', '/api/videos').respond([], {
      'A-Token': 'xxx'
    });

    video = $injector.get('VideoResource');
  }));

  afterEach(function() {
    $httpBackend.verifyNoOutstandingExpectation();
    $httpBackend.verifyNoOutstandingRequest();
  });

  it('should have method getAll', function() {
    var videos = video.getAll();
    $httpBackend.flush();

    expect(videos.length).toBe(0);
  });
});
