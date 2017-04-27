angular.module('snailpark', ['ngWebSocket'])
  .factory('gameServer', function($websocket) {
    var playerId = 'player';

    var gameId = Date.now();
    var s = 'game/connect?gameId=' + gameId;
    var l = window.location;
    var url = ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;

    var dataStream = $websocket(url);

    var data = {
      currentPlayerId: null,
      state: null,
      player: {},
      enemy:  {}
    };

    dataStream.onMessage(function(message) {
      var msg = JSON.parse(message.data)

      var filterAttackers = function(msg) {
        var attackers = [];
        for(var i in msg.engagements) {
          var e = msg.engagements[i];
          attackers.push(e.attacker);
        }

        return attackers;
      }

      data.currentPlayerId = msg.currentPlayerId;
      data.state = msg.state;
      data.attackers = filterAttackers(msg);
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
      },

      playCard: function(id) {
        dataStream.send(JSON.stringify({ action: 'playCard', playerId: playerId, card: id }));
      },

      targetCard: function(id) {
        dataStream.send(JSON.stringify({ action: 'target', playerId: playerId, card: id }));
      },

      ping: function() {
        dataStream.send(JSON.stringify({ action: 'ping', playerId: playerId }));
        setTimeout(this.ping.bind(this), 45 * 1000);
      }
    };

    return methods;
  })
  .controller('BoardController', ['$scope', 'gameServer', function ($scope, gameServer) {
    gameServer.start();

    $scope.data = gameServer.data;
    $scope.next = gameServer.next
    $scope.playCard = gameServer.playCard
    $scope.targetCard = gameServer.targetCard
    $scope.newGame = function() {
      console.log('Not supported');
    }

    gameServer.ping();
  }])

angular.module('snailpark')
  .directive('cardList', function() {
    return {
      scope: {
        cards: '=',
        clickCard: '&'
      },
      restrict : 'EA',
      controller: function() {},
      controllerAs: 'ctrl',
      transclude: true,
      bindToController: true,
      template: '<li ng-repeat="card in ctrl.cards" ><card data-set="card" click-card="ctrl.clickCard({ id: id })"></card></li>'
    }
  });

angular.module('snailpark')
  .directive('card', function() {
    return {
      scope: {
        card: '=set',
        clickCard: '&'
      },
      replace: true,
      controller: function() {},
      link: function(scope) {
        scope.attacking = function() {
          var value = scope.ctrl.card.tags.attackTarget;
          console.log(value);
          return typeof value != 'undefined' && value != '';
        }
      },
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
        currentPlayerId: '=',
        state: '=',
        next: '&'
      },
      restrict : 'EA',
      controller: function() {},
      link: function(scope) {
        scope.disabled = function() {
          return scope.ctrl.currentPlayerId != "player";
        }

        scope.btnText = function() {
          if(scope.ctrl.state == "blockers" || scope.ctrl.state == "blockTarget") {
            return 'Begin combat';
          } else {
            return 'End turn';
          }
        }

        scope.help = function() {
          if(scope.ctrl.currentPlayerId === 'ai') {
            return 'Your opponent is thinking...';
          } else if(scope.ctrl.state == "main") {
            return 'Play a card, or attack with your creatures by clicking on them...';
          } else if(scope.ctrl.state == "targeting") {
            return 'Pick a target...';
          } else if(scope.ctrl.state == "blockers" || scope.ctrl.state == "blockTarget") {
            return 'Opponent is attacking! Assign blockers and end turn when you are ready';
          } else {
            return '...';
          }
        }
      },
      controllerAs: 'ctrl',
      transclude: true,
      bindToController: true,
      template: '<input type="button" id="end-turn" value="{{ btnText() }}" ng-click="ctrl.next()" ng-class="{ \'btn-disabled\': disabled() }" class="btn pull-left"/><div id="state-help" class="pull-left">{{ help() }}</div>'
    }
  });

angular.module('snailpark')
  .directive('modalDialog', function() {
    return {
      scope: {
        state: '=',
        playerAvatar: '=',
        newGame: '&'
      },
      restrict : 'EA',
      controller: function() {},
      controllerAs: 'ctrl',
      link: function(scope) {
        scope.content = function() {
          if(scope.ctrl.playerAvatar["attributes"]["toughness"] <= 0) {
            return 'You lost! :(';
          } else {
            return 'You won! :)';
          }
        }
      },
      transclude: true,
      bindToController: true,
      template: '<div id="myModal" ng-if="ctrl.state == \'finished\'" class="modal"><div class="modal-content">{{ content() }}</div></div>'
    }
  });

angular.element(function() {
  angular.bootstrap(document, ['snailpark']);
});
