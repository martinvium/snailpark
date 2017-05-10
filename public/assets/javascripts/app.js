var app = angular.module('snailpark', ['ngWebSocket']);

app.factory('gameServer', function($websocket) {
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
    var packet = JSON.parse(message.data)

    if(packet.t != 'FULL_STATE') {
      console.log('Only support FULL_STATE');
      return;
    }

    var msg = packet.m;

    var filterAttackers = function(msg) {
      var attackers = [];
      for(var i in msg.entities) {
        if (msg.entities[i].tags["attackTarget"]) {
          attackers.push(msg.entities[i]);
        }
      }

      return attackers;
    }

    var findGameEntity = function(msg) {
      for(var i in msg.entities) {
        if(msg.entities[i].tags.type == "game") {
          return msg.entities[i];
        }
      }
    }

    data.game = findGameEntity(msg);
    data.currentPlayerId = data.game.tags.currentPlayerId;
    data.state = data.game.tags.state;
    data.targeting = ['targeting', 'blockTarget'].indexOf(data.state) !== -1;
    data.attackers = filterAttackers(msg);
    data.entities = msg.entities
    data.player = msg.players["player"];
    data.enemy = msg.players["ai"];
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
});

app.controller('BoardController', ['$scope', 'gameServer', function ($scope, gameServer) {
  gameServer.start();

  $scope.data = gameServer.data;
  $scope.enemyId = 'ai';
  $scope.playerId = 'player';
  $scope.cardDetails = { card: null };
  $scope.next = gameServer.next
  $scope.playCard = gameServer.playCard
  $scope.targetCard = gameServer.targetCard

  $scope.$watch('data.entities', function(n) {
    $scope.entities = {
      ai: {
        hand: filterPlayerAndLocation(gameServer.data.entities, 'ai', 'hand'),
        board: filterPlayerAndLocation(gameServer.data.entities, 'ai', 'board'),
      },
      player: {
        hand: filterPlayerAndLocation(gameServer.data.entities, 'player', 'hand'),
        board: filterPlayerAndLocation(gameServer.data.entities, 'player', 'board')
      }
    }
  });

  function filterPlayerAndLocation(s, p, l) {
    var filtered = [];
    for(var i in s) {
      if(s[i].playerId == p && s[i].tags['location'] == l) {
        filtered.push(s[i])
      }
    }

    return filtered;
  }

  $scope.newGame = function() {
    console.log('Not supported');
  }

  $scope.showCardDetails = function(args) {
    $scope.cardDetails.card = args.card;
    $scope.$apply();
  }

  gameServer.ping();
}])

app.directive('cardList', function() {
  return {
    scope: {
      cards: '=',
      cardDetails: '&',
      clickCard: '&'
    },
    restrict : 'EA',
    controller: function() {},
    controllerAs: 'ctrl',
    transclude: true,
    bindToController: true,
    template: '<card ng-repeat="card in ctrl.cards" data-set="card" data-card-details="ctrl.cardDetails({card: card})" click-card="ctrl.clickCard({ id: id })"></card>'
  }
});

app.directive('card', function() {
  return {
    scope: {
      card: '=set',
      cardDetails: '&',
      clickCard: '&'
    },
    replace: true,
    controller: function() {},
    link: function(scope, element) {
      scope.attacking = function() {
        var value = scope.ctrl.card.tags.attackTarget;
        return typeof value != 'undefined' && value != '';
      }

      if(typeof scope.ctrl.cardDetails != 'undefined') {
        element.on('mouseover', function(e) {
          scope.ctrl.cardDetails({ card: scope.ctrl.card });
        });

        element.on('mouseout', function(e) {
          scope.ctrl.cardDetails({ card: null });
        });
      } else {
        console.log('no details for card');
      }
    },
    controllerAs: 'ctrl',
    bindToController: true,
    restrict: 'EA',
    templateUrl: 'card-full.html?' + Date.now()
  }
});

app.directive('energyIndicator', function() {
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
    link: function(scope) {
      scope.range = [1, 2, 3, 4, 5, 6, 7];
    },
    template: '<div class="energy"><ul><li ng-repeat="i in range" ng-class="{ \'current\': i <= ctrl.current, \'dead\': i > ctrl.max }"></li></ul></div>'
  }
});

app.directive('nextButton', function() {
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
    },
    controllerAs: 'ctrl',
    transclude: true,
    bindToController: true,
    template: '<input type="button" id="end-turn" value="{{ btnText() }}" ng-click="ctrl.next()" ng-class="{ \'btn-disabled\': disabled() }" class="btn pull-left"/>'
  }
});

app.directive('helpText', function() {
  return {
    scope: {
      currentPlayerId: '=',
      state: '=',
    },
    restrict : 'EA',
    controller: function() {},
    link: function(scope) {
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
    template: '<div class="state-help">{{ help() }}</div>'
  }
});


app.directive('modalDialog', function() {
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
