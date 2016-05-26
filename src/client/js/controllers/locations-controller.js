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

     $scope.deleteLocationClicked = function (entity) {
        var res = $http.delete(getUrl() + '/v1.0/Locations(' + entity["@iot.id"] + ')');
        res.success(function(data, status, headers, config) {
            var index = $scope.locationsList.indexOf(entity);
            $scope.locationsList.splice(index, 1);
        });
        res.error(function(data, status, headers, config) {
            alert( "failure: " + JSON.stringify({data: data}));
        });
     };
});