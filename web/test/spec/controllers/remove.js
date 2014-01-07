'use strict';

describe('Controller: RemoveCtrl', function () {

  // load the controller's module
  beforeEach(module('webApp'));

  var RemoveCtrl,
    scope, Page;

  // Initialize the controller and a mock scope
  beforeEach(inject(function ($controller, $rootScope, _Page_) {
    Page = _Page_;
    scope = $rootScope.$new();
    RemoveCtrl = $controller('RemoveCtrl', {
      $scope: scope,
      Video: {Name: 'Novecento'}
    });
  }));

  it('should attach video object to the scope', function () {
    expect(scope.video).toBeDefined();
  });

  it('should set the title to remove {name}', function() {
    expect(Page.title()).toBe('Remove Novecento?');
  });
});
