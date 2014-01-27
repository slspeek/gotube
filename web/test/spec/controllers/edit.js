'use strict';

describe('Controller: EditCtrl', function() {


  beforeEach(module('webApp', function($provide) {
    $provide.value('$route', {
      current: {
        params: {
          VideoId: 1
        }
      }
    });
  }));

  var EditCtrl,
    scope, httpBackend, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, $httpBackend, _Page_, userLoader, videoLoader) {
    Page = _Page_;
    scope = $rootScope.$new();
    httpBackend = $httpBackend;
    httpBackend.expect('GET', '/auth').respond({
      username: 'Misko'
    });
    httpBackend.expect('GET', '/api/videos/1').respond({
      Id: 1,
      Owner: 'Misko',
      Name: 'Novecento',
      Desc: 'Italian classic'
    });
    var usernamePromise = userLoader();
    var username;
    usernamePromise.then(function(name) {
      username = name;
    });
    var videoPromise = videoLoader();
    var video;
    videoPromise.then(function(result) {
      video = result;
    });
    httpBackend.flush();
    EditCtrl = $controller('EditCtrl', {
      $scope: scope,
      UserName: username,
      Video: video
    });
  }));

  afterEach(function() {
    httpBackend.verifyNoOutstandingExpectation();
    httpBackend.verifyNoOutstandingRequest();
  });

  it('should set the title to upload', function() {
    expect(Page.title()).toBe('Edit');
  });
  it('should call the server on save', function() {
    httpBackend.expect('PUT', '/api/videos/1').respond({
      Id: 1
    });
    scope.save();
    httpBackend.flush();
  });



});
