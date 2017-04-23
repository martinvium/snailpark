angular.module('snailpark', ['ngWebSocket'])
  .factory('gameServer', function($websocket) {
    var playerId = 'player';

    var gameId = Date.now();
    var s = 'game/connect?gameId=' + gameId;
    var l = window.location;
    var url = ((l.protocol === "https:") ? "wss://" : "ws://") + l.host + l.pathname + s;

    var dataStream = $websocket(url);

    var collection = [];

    dataStream.onMessage(function(message) {
      collection.push(JSON.parse(message.data));
    });

    var methods = {
      collection: collection,
      start: function() {
        dataStream.send(JSON.stringify({ action: 'start', playerId: playerId }));
      }
    };

    return methods;
  })
  .controller('BoardController', ['$scope', 'gameServer', function ($scope, gameServer) {
    gameServer.start();
    // gameServer.collection.then(function(data) {
    //   console.log(data);
    // });

    $scope.cards = [
      { id: 'my-card', title: 'bla' },
      { id: 'my-card2', title: 'bla2' }
    ];
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
      templateUrl: 'card-full.html'
    }
  });

angular.element(function() {
  angular.bootstrap(document, ['snailpark']);
});
