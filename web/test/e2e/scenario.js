/* global by:false, element:false , protractor:false */
describe('Gotube proof of concept scenario', function() {
  'use strict';

  var protractor;

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
      element(by.id('flow-btn-input-id')).sendKeys(browser.params.testMovie);
      //element(by.id('flow-btn-input-id')).sendKeys('test-data/BetterLife_HighQuality.ogv');
      browser.sleep(3000);
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
