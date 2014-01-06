(function() {
  'use strict';

  angular.module('webApp')
    .controller('ViewCtrl', function($scope, $routeParams, $sce, Video) {
      $scope.video = Video;
      $scope.name = $scope.video.Name;
      $scope.desc = $scope.video.Desc;
      $scope.id = $scope.video.Id;
      $scope.videoURL = $sce.trustAsResourceUrl('/content/videos/' + $routeParams.VideoId );
      $scope.downloadURL = '/content/videos/' + $routeParams.VideoId + '/download';

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
