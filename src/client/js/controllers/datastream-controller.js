gostApp.controller('DatastreamCtrl', function ($scope, $http, $routeParams, Page) {
    $scope.id = $routeParams.id;
    $scope.Page.setTitle('DATASTREAM(' + $scope.id + ')');
    $scope.Page.setHeaderIcon(iconDatastream);

    $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")").then(function (response) {
        $scope.description = response.data["description"];
        $scope.unitOfMeasurement = response.data["unitOfMeasurement"];
        $scope.observedArea = response.data["observedArea"];
        $scope.Page.selectedDatastream = response.data;
    });

    $scope.tabPropertiesClicked = function () {

    };

    $scope.tabThingClicked = function () {
        $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")/Thing").then(function (response) {
            $scope.thingId = response.data["@iot.id"];
            $scope.thingDescription = response.data["description"];
            $scope.thingProperties = response.data["properties"];
        });
    };

    $scope.tabSensorClicked = function () {
        $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")/Sensor").then(function (response) {
            $scope.sensorId = response.data["@iot.id"];
            $scope.sensorDescription = response.data["description"];
            $scope.sensorEncoding = response.data["encodingtype"];
            $scope.sensorMetadata = response.data["metadata"];
        });
    };

    $scope.tabObservationsClicked = function () {
        $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")/Observations").then(function (response) {
            $scope.observationsList = response.data.value;
            labels = [];
            values = [];
            angular.forEach($scope.observationsList, function (value, key) {
                labels.push(value['phenomenonTime']);
                values.push(value['result']);
            });

            createObservationChart(labels, values);
        });
    };

    $scope.tabObservedPropertyClicked = function () {
        $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")/ObservedProperty").then(function (response) {
            $scope.observedPropertyId = response.data["@iot.id"];
            $scope.observedPropertyName = response.data["name"];
            $scope.observedPropertyDescription = response.data["description"];
            $scope.observedPropertyDefinition = response.data["definition"];
        });
    };
});