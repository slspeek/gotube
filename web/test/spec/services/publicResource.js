'use strict';

describe('Service: PublicResource', function () {

  // load the service's module
  beforeEach(module('webApp'));

  // instantiate service
  var PublicResource;
  beforeEach(inject(function (_PublicResource_) {
    PublicResource = _PublicResource_;
  }));

  it('should do something', function () {
    expect(!!PublicResource).toBe(true);
  });

});
