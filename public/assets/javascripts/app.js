angular.module('snailpark', [])
  .controller('BoardController', ['$scope', function ($scope) {
    $scope.greetMe = 'World';
  }]);

angular.element(function() {
  angular.bootstrap(document, ['snailpark']);
});



// app.directive('card', function() {
//   return {
//     scope: {
//       item: '=set',
//       onClick: '&'
//     },
//     replace: true,
//     controller: function() {},
//     controllerAs: 'ctrl',
//     bindToController: true,
//     restrict: 'EA',
//     template: '<input type="checkbox" ng-click="ctrl.onClick({item: ctrl.item})" ng-checked="ctrl.item.active" /> '
//   }
// });
