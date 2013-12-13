/* global by:false, element:false */
describe('Gotube proof of concept scenario', function() {
  'use strict';

  beforeEach(function() {
    browser.get('/');
  });


  it(
    'should automatically redirect to login', function() {
      expect(browser.getCurrentUrl()).toContain( '#/login');
      element(by.model('username')).clear();
      element(by.model('username')).sendKeys('steven');
      element(by.model('password')).clear();
      element(by.model('password')).sendKeys('gnu');
      element(by.id('login')).click();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.id('upload')).click();
      element(by.model('name')).sendKeys('Better life');
      element(by.id('desc')).sendKeys('Cartoon');
      element(by.id('flow-btn-input-id')).sendKeys('/home/steven/projs/nog/src/github.com/slspeek/gotube/test-data/BetterLife_HighQuality.ogv');
      //element(by.id('upload')).click();
      //browser.sleep(3000);
      browser.get('/');
      expect(browser.getCurrentUrl()).toContain( '#/login');
      element(by.model('username')).clear();
      element(by.model('username')).sendKeys('steven');
      element(by.model('password')).clear();
      element(by.model('password')).sendKeys('gnu');
      element(by.id('login')).click();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.linkText('Better life')).click();
      expect(element(by.binding('{{name}}')).getText()).toBe('Better life');
      expect(element(by.binding('{{desc}}')).getText()).toBe('Cartoon');

    }, 10000);


});