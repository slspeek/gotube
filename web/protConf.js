'use strict';
var path = require('path');

var abs = function(relative) {
  var r = path.resolve(relative);
  console.log(r);
  return r;
};

exports.config = {
  // The address of a running selenium server.
  seleniumAddress: 'http://localhost:4444/wd/hub',


  framework: ['jasmine', 'jasmine-reporters'],
  baseUrl: 'http://localhost:8484',

  // Capabilities to be passed to the webdriver instance.
  capabilities: {
    'browserName': 'firefox'
  },

  // Spec patterns are relative to the current working directly when
  // protractor is called.
  specs: ['test/e2e/scenario.js'],

  params: {
    testMovie: abs('../test-data/BetterLife_HighQuality.ogv')
  },

 /* onPrepare: function() {*/
    //jasmine.getEnv().addReporter(new jasmine.JUnitXmlReporter(
      //'junitxml/', true, true));
  /*},*/


  // Options to be passed to Jasmine-node.
  jasmineNodeOpts: {
    showColors: true,
    defaultTimeoutInterval: 30000
  }
};
