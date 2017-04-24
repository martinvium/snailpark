angular.module('snailpark', ['ngWebSocket'])
  .factory('gameServer', function($websocket) {
    var playerId = 'player';

    var gameId = Date.now();
    var s = 'game/connect?gameId=' + gameId;
    var l = window.location;
    var url = ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;

    var dataStream = $websocket(url);

    var data = {
      player: {},
      enemy:  {}
    };

    dataStream.onMessage(function(message) {
      var msg = JSON.parse(message.data)
      data.player = msg.players["player"];
      data.player.hand = msg.players["player"]["hand"];
      data.player.board = msg.players["player"]["board"];
      data.enemy = msg.players["ai"];
      data.enemy.hand = msg.players["ai"]["hand"];
      data.enemy.board = msg.players["ai"]["board"];
    });

    var methods = {
      data: data,

      start: function() {
        dataStream.send(JSON.stringify({ action: 'start', playerId: playerId }));
      },

      next: function() {
        dataStream.send(JSON.stringify({ action: 'endTurn', playerId: playerId }));
      }
    };

    return methods;
  })
  .controller('BoardController', ['$scope', 'gameServer', function ($scope, gameServer) {
    gameServer.start();
    $scope.data = gameServer.data;

    $scope.next = function() {
      gameServer.next()
    }
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

angular.module('snailpark')
  .directive('energyIndicator', function() {
    return {
      scope: {
        current: '=',
        max: '=',
        title: '@'
      },
      restrict : 'EA',
      controller: function() {},
      controllerAs: 'ctrl',
      transclude: true,
      bindToController: true,
      template: '<div class="mana">{{ ctrl.title }} energy: {{ ctrl.current }} of {{ ctrl.max }}</div>'
    }
  });

angular.module('snailpark')
  .directive('nextButton', function() {
    return {
      scope: {
        state: '=',
        next: '&'
      },
      restrict : 'EA',
      controller: function() {},
      controllerAs: 'ctrl',
      transclude: true,
      bindToController: true,
      template: '<input type="button" id="end-turn" value="End turn" ng-click="ctrl.next()" class="btn pull-left"/><div id="state-help" class="pull-left"></div>'
    }
  });

angular.element(function() {
  angular.bootstrap(document, ['snailpark']);
});
