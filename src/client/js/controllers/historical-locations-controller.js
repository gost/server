gostApp.controller('HistoricalLocationsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('HISTORICAL LOCATIONS');
    $scope.Page.setHeaderIcon(iconHistoricalLocation);

    $http.get(getUrl() + "/v1.0/HistoricalLocations").then(function (response) {
        $scope.historicallocationsList = response.data.value;
    });
});