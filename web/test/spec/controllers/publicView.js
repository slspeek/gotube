'use strict';

describe('Controller: PublicviewCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var PublicviewCtrl,
    scope, page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, Page) {
    scope = $rootScope.$new();
    page = Page;
    PublicviewCtrl = $controller('PublicviewCtrl', {
      $scope: scope,
      $routeParams: {VideoId:'345'},
      Video: {'Id':'345', 'Name':'Novecento', 'Desc':'Italian classic'}
    });
  }));



  it('should attach video to the scope', function() {
    expect(scope.video).toBeDefined();
  });
  it('should attach desc to the scope', function() {
    expect(scope.video.Desc).toBe('Italian classic');
  });
  it('should attach name to the scope', function() {
    expect(scope.video.Id).toBe('345');
  });
  it('should attach downloadURL to the scope', function() {
    expect(scope.downloadURL).toBe('/public/content/videos/345/download');
  });
  it('should set the title on page', function() {
    expect(page.title()).toBe('Novecento');
  });
});

