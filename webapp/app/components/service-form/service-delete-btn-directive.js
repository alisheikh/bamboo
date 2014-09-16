module.exports = function () {
  return {
    restrict: "AE",
    template:  '<button class="btn btn-danger" ng-click="showModal()" title="Delete"><i class="icon ion-android-trash"></i></button>',
    scope: {
      serviceId: "="
    },
    controller: ["$scope", "Service", "$modal", "$rootScope", function ($scope, Domain, $modal, $rootScope) {
      $scope.actionName = "Delete It!";

      console.log($scope.serviceId);

      $scope.showModal = function () {
        $scope.modal = $modal({
          title: "Are you sure?",
          template: "bamboo/modal-confirm",
          content: "Delete Marathon ID " + $scope.serviceId,
          scope: $scope,
          show: true
        });
      };

      $scope.doAction = function () {
        Domain.destroy({
            id: $scope.serviceId
          })
          .then(function () {
            $scope.modal.hide();
            $scope.modal = null;
            $rootScope.$broadcast("services.reset");
          });
      };
    }]
  }
};