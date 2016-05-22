gostApp.controller('ObservationsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('OBSERVATIONS');
    $scope.Page.setHeaderIcon(iconObservation);

    $http.get(getUrl() + "/v1.0/Observations").then(function (response) {
        $scope.observationsList = response.data.value;
    });
});