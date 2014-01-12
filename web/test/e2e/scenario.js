/* global by:false, element:false */
describe('Gotube proof of concept scenario', function() {
  'use strict';

  beforeEach(function() {
    browser.get('/');
  });

  var login = function() {
    expect(browser.getCurrentUrl()).toContain('#/login');
    element(by.model('username')).clear();
    element(by.model('username')).sendKeys('steven');
    element(by.model('password')).clear();
    element(by.model('password')).sendKeys('gnu');
    element(by.id('login')).click();
  };


  it(
    'should automatically redirect to login', function() {
      login();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.id('upload')).click();
      element(by.model('name')).sendKeys('Better life');
      element(by.id('desc')).sendKeys('Cartoon');
      element(by.id('flow-btn-input-id')).sendKeys(browser.params.testMovie);
      browser.sleep(3000);
      browser.get('/');
      //login();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.linkText('Better life')).click();
      expect(element(by.binding('{{video.Name}}')).getText()).toBe('Better life');
      expect(element(by.binding('{{video.Desc}}')).getText()).toBe('Cartoon');
      browser.get('/');
      //login();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.className('glyphicon-remove')).click();
      expect(browser.getCurrentUrl()).toContain('#/remove');
      element(by.linkText('Remove')).click();
      browser.sleep(1000);
      expect(browser.getCurrentUrl()).toContain('#/list');
      expect(element(by.tagName('body')).getText()).not.toContain('Better life');

    }, 20000);


});
