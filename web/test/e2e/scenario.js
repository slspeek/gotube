/* global by:false, element:false */
describe('Gotube proof of concept scenario', function() {
  'use strict';

  beforeEach(function() {
    browser.get('/#/list');
  });

  var login = function() {
    expect(browser.getCurrentUrl()).toContain('#/login');
    element(by.model('username')).clear();
    element(by.model('username')).sendKeys('steven');
    element(by.model('password')).clear();
    element(by.model('password')).sendKeys('gnu');
    element(by.id('login')).click();
  };

  var newTitle = 'Life without money';
  it(
    'should upload a video and view it', function() {
      login();
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.id('upload')).click();
      element(by.model('name')).sendKeys('Better life');
      element(by.id('desc')).sendKeys('Cartoon');
      element(by.id('flow-btn-input-id')).sendKeys(browser.params.testMovie);
      //browser.sleep(3000);
      browser.get('/#/list');
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.className('glyphicon-pencil')).click();
      expect(browser.getCurrentUrl()).toContain('#/edit');
      expect(element(by.model('video.Public')).isSelected()).toBe(false);
      element(by.model('video.Name')).clear();
      element(by.model('video.Name')).sendKeys(newTitle);
      element(by.id('save-button')).click();

      //browser.get('/#/list'); 
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.linkText(newTitle)).click();
      expect(browser.getCurrentUrl()).toContain('#/view');
      expect(element(by.binding('{{video.Name}}')).getText()).toBe(newTitle);
      expect(element(by.binding('{{video.Desc}}')).getText()).toBe('Cartoon');
      browser.get('/#/');
      expect(browser.getCurrentUrl()).toContain('#/public');
      expect(element(by.tagName('body')).getText()).not.toContain(newTitle);

      browser.get('/#/list');
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.className('glyphicon-pencil')).click();
      expect(browser.getCurrentUrl()).toContain('#/edit');
      element(by.model('video.Public')).click();
      expect(element(by.model('video.Public')).isSelected()).toBe(true);
      element(by.id('save-button')).click();

      browser.get('/#/');
      expect(browser.getCurrentUrl()).toContain('#/public');
      expect(element(by.tagName('body')).getText()).toContain(newTitle);

      browser.get('/#/list');
      expect(browser.getCurrentUrl()).toContain('#/list');
      element(by.className('glyphicon-remove')).click();
      expect(browser.getCurrentUrl()).toContain('#/remove');
      element(by.linkText('Remove')).click();
      browser.sleep(1000);
      expect(browser.getCurrentUrl()).toContain('#/list');
      expect(element(by.tagName('body')).getText()).not.toContain(newTitle);

    }, 30000);


});
