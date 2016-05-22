gostApp.controller('ThingCtrl', function ($scope, $http, $routeParams, Page) {
    $scope.id = $routeParams.id;
    $scope.Page.setTitle('THING(' + $scope.id + ')');
    $scope.Page.setHeaderIcon(iconThing);

    $http.get(getUrl() + "/v1.0/Things(" + $scope.id + ")").then(function (response) {
        $scope.description = response.data["description"];
        $scope.properties = response.data["properties"];
        $scope.Page.selectedThing = response.data;
    });

    $scope.tabPropertiesClicked = function () {
        $scope.mapVisible = false;
    };

    $scope.tabLocationsClicked = function () {
        $scope.mapVisible = true;
        createMap();

        $http.get(getUrl() + "/v1.0/Things(" + $scope.Page.selectedThing["@iot.id"] + ")/Locations").then(function (response) {
            $scope.locationsList = response.data.value;
            angular.forEach($scope.locationsList, function (value, key) {
                addGeoJSONFeature(value["location"]);
            });

            zoomToGeoJSONLayerExtent();
        });
    };

    $scope.tabHistoricalLocationsClicked = function () {
        $scope.mapVisible = false;

        $http.get(getUrl() + "/v1.0/Things(" + $scope.Page.selectedThing["@iot.id"] + ")/HistoricalLocations").then(function (response) {
            $scope.historicalLocationsList = response.data.value;
        });
    };

    $scope.tabDatastreamsClicked = function () {
        $scope.mapVisible = false;

        $http.get(getUrl() + "/v1.0/Things(" + $scope.Page.selectedThing["@iot.id"] + ")/Datastreams").then(function (response) {
            $scope.datastreamsList = response.data.value;
        });
    };

    $scope.datastreamClicked = function (datastreamID) {
        angular.forEach($scope.things, function (value, key) {
            if (value["@iot.id"] == thingID) {
                $scope.Page.selectedDatastream = value;
            }
        });

        $scope.Page.go("datastream/" + datastreamID);
    };
});