'use strict';

describe('Service: Page', function() {

  // load the service's module
  beforeEach(module('webApp'));

  // instantiate service
  var Page;
  beforeEach(inject(function(_Page_) {
    Page = _Page_;
  }));

  it('should default to Gotube', function() {
    expect(Page.title()).toBe('Gotube');
  });

  it('should be a simple struct', function() {
    var tt = 'Test title';
    Page.setTitle(tt);
    expect(Page.title()).toBe(tt);
  });

});
