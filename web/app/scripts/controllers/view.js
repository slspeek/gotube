(function() {
  'use strict';

  angular.module('webApp')
    .controller('ViewCtrl', function($scope,  $sce, Video, Page) {
      $scope.video = Video;
      Page.setTitle(Video.Name);
      var id = Video.Id;
      $scope.videoURL = $sce.trustAsResourceUrl('/content/videos/' + id);
      $scope.downloadURL = '/content/videos/' + id + '/download';

      $scope.stretchModes = [{
        label: 'None',
        value: 'none'
      }, {
        label: 'Fit',
        value: 'fit'
      }, {
        label: 'Fill',
        value: 'fill'
      }];

      $scope.config = {
        width: 740,
        height: 380,
        autoHide: true,
        autoPlay: false,
        responsive: true,
        stretch: $scope.stretchModes[1],
        theme: {
          url: 'bower_components/videogular-themes-default/videogular.css',

          playIcon: '&#xe000;',
          pauseIcon: '&#xe001;',
          volumeLevel3Icon: '&#xe002;',
          volumeLevel2Icon: '&#xe003;',
          volumeLevel1Icon: '&#xe004;',
          volumeLevel0Icon: '&#xe005;',
          muteIcon: '&#xe006;',
          enterFullScreenIcon: '&#xe007;',
          exitFullScreenIcon: '&#xe008;'
        }
      };

    });

})();
