'use strict';

describe('Service: Page', function() {

  // load the service's module
  beforeEach(module('webApp'));

  // instantiate service
  var Page, $rootScope, $location;
  beforeEach(inject(['Page', '$rootScope', '$location', function(P, $rS, $l) {
    Page = P;
    $rootScope = $rS;
    $location = $l;
  }]));

  it('should default to Gotube', function() {
    expect(Page.title()).toBe('Gotube');
  });

  it('should be a simple struct', function() {
    var tt = 'Test title';
    Page.setTitle(tt);
    expect(Page.title()).toBe(tt);
  });
  
  it('should set onUpload to false', function() {
    expect(Page.onUpload()).toBe(false);
    
  });

  it('should set onMyVideos to false', function() {
    expect(Page.onMyVideos()).toBe(false);
  });
  
  it('should set onMyVideos with /list route', function() {
    $location.path('/list');
    $rootScope.$broadcast('$routeChangeSuccess', '', '/list');
    expect(Page.onMyVideos()).toBe(true);
    });
  
  
  it('should set onUpload with /upload route', function() {
    $location.path('/upload');
    $rootScope.$broadcast('$routeChangeSuccess', '', '/list');
    expect(Page.onUpload()).toBe(true);
    });
});
