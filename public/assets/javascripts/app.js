angular.module('snailpark', ['ngWebSocket'])
  .factory('gameServer', function($websocket) {
    var playerId = 'player';

    var gameId = Date.now();
    var s = 'game/connect?gameId=' + gameId;
    var l = window.location;
    var url = ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;

    var dataStream = $websocket(url);

    var data = {
      cards: []
    };

    dataStream.onMessage(function(message) {
      var msg = JSON.parse(message.data)
      data.cards = msg.players["player"]["hand"];
    });

    var methods = {
      data: data,
      start: function() {
        dataStream.send(JSON.stringify({ action: 'start', playerId: playerId }));
      }
    };

    return methods;
  })
  .controller('BoardController', ['$scope', 'gameServer', function ($scope, gameServer) {
    gameServer.start();
    $scope.data = gameServer.data;
  }]);

angular.module('snailpark')
  .directive('cardList', function() {
    return {
      scope: {
        cards: '='
      },
      restrict : 'EA',
      controller: function() {},
      controllerAs: 'ctrl',
      transclude: true,
      bindToController: true,
      template: '<li ng-repeat="card in ctrl.cards"><card data-set="card"/></li>'
    }
  });

angular.module('snailpark')
  .directive('card', function() {
    return {
      scope: {
        card: '=set',
        onClick: '&'
      },
      replace: true,
      controller: function() {},
      controllerAs: 'ctrl',
      bindToController: true,
      restrict: 'EA',
      templateUrl: 'card-full.html?' + Date.now()
    }
  });

angular.element(function() {
  angular.bootstrap(document, ['snailpark']);
});
