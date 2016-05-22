gostApp.controller('LocationsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('LOCATIONS');
    $scope.Page.setHeaderIcon(iconLocation);

    createMap();
    $http.get(getUrl() + "/v1.0/Locations").then(function (response) {
        $scope.locationsList = response.data.value;
        angular.forEach($scope.locationsList, function (value, key) {
            addGeoJSONFeature(value["location"]);
        });

        zoomToGeoJSONLayerExtent();
    });
});