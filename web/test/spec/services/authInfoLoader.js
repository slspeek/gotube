'use strict';

describe('Service: authInfoLoader', function () {

  // load the service's module
  beforeEach(module('webApp'));

  // instantiate service
  var authInfoLoader;
  beforeEach(inject(['authInfoLoader', function ($ail) {
    authInfoLoader = $ail;
  }]));

  it('should do something', function () {
    expect(!!authInfoLoader).toBe(true);
  });

});
