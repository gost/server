gostApp.controller('DatastreamsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('DATASTREAMS');
    $scope.Page.setHeaderIcon(iconDatastream);

    $http.get(getUrl() + "/v1.0/Datastreams").then(function (response) {
        $scope.datastreamsList = response.data.value;
    });

    $scope.datastreamClicked = function (datastreamID) {
        angular.forEach($scope.things, function (value, key) {
            if (value["@iot.id"] == thingID) {
                $scope.Page.selectedDatastream = value;
            }
        });

        $scope.Page.go("datastream/" + datastreamID);
    };
});
