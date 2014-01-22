
'use strict';

describe('Service: publicVideo', function() {

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

  beforeEach(module('webApp', function($provide) {
    $provide.value('$route', {
      current: {
        params: {
          VideoId: 1
        }
      }
    });
  }));

  var Loader, $httpBackend;
  beforeEach(inject(function(publicVideo) {
    Loader = publicVideo;
  }));

  beforeEach(inject(function($injector) {
    $httpBackend = $injector.get('$httpBackend');
    $httpBackend.when('GET', '/public/api/videos/1').respond(pinokio);
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
    expect(video).toEqualData(pinokio);
  });
});
