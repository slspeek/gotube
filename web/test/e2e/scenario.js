/* global element:false, sleep:false, input:false */
describe('Gotube proof of concept scenario', function() {
  'use strict';

  beforeEach(function() {
    browser().navigateTo('/');
  });


  it(
    'should automatically redirect to login', function() {
      expect(browser().location().url()).toBe( '/login');
      input('username').enter('steven');
      input('password').enter('gnu');
      element('#login').click();
      expect(browser().location().url()).toBe( '/list');
      sleep(1);
    });


});
