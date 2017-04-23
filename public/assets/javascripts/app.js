angular.module('snailpark', [])
  .controller('BoardController', ['$scope', function ($scope) {
    $scope.greetMe = 'World';
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
