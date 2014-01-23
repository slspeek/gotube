'use strict';

describe('Controller: ListCtrl', function() {

  // load the controller's module
  beforeEach(module('webApp'));

  var ListCtrl,
    scope, page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function($controller, $rootScope, Page) {
    page = Page;
    scope = $rootScope.$new();
    ListCtrl = $controller('ListCtrl', {
      $scope: scope,
      UserName: 'Misko',
      VideoList: []
    });
  }));


  it('should add a list of videos to the scope', function() {
    expect(scope.videoList.length).toBe(0);
  });

  it('should add username to the scope', function() {
    expect(scope.username).toBe('Misko');
  });

  it('should set the title to Listing', function() {
    expect(page.title()).toBe('Listing');
  });

  it('should use private view url', function() {
   expect(scope.viewUrl(0)).toBe('/#/view/0'); 
  });
  
});
