gostApp.controller('DatastreamCtrl', function ($scope, $http, $routeParams, Page) {
    $scope.id = $routeParams.id;
    $scope.Page.setTitle('DATASTREAM(' + $scope.id + ')');
    $scope.Page.setHeaderIcon(iconDatastream);

    $scope.$on("$destroy", function () {
        client.unsubscribe("Datastreams(" + $scope.id + ")/Observations");
    });

    labels = [];
    values = [];

    client = new Paho.MQTT.Client(location.hostname, Number(9001), guid());
    client.onConnectionLost = onConnectionLost;
    client.onMessageArrived = onMessageArrived;
    client.connect({ onSuccess: onConnect });

      function onConnect() {
        client.subscribe("Datastreams(" + $scope.id + ")/Observations");
        console.log("Datastreams(" + $scope.id + ")/Observations subscribed");
    }

    function onConnectionLost(responseObject) {
        if (responseObject.errorCode !== 0) {
            console.log("onConnectionLost:" + responseObject.errorMessage);
        }
    }

    function onMessageArrived(message) {
        switch (message.destinationName) {
            case 'Datastreams(' + $scope.id + ')/Observations':
				try {
					  //$scope.brokerVersion = message.payloadString;
					var r = JSON.parse(message.payloadString);
					//$scope.observationsList.push(r)

					values.push(r['result'])
					labels.push(r['phenomenonTime'])
					
					if(values.length > 50){
						values.splice(0, 1);
					}
					
					if(labels.length > 50){
						labels.splice(0, 1);
					}
					
					observationChartAddData(labels, values);
					//$scope.$apply();
					break;
				}
				catch(err) {
					
				}
              
        }
    }

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
            $scope.sensorEncoding = response.data["encodingType"];
            $scope.sensorMetadata = response.data["metadata"];
        });
    };

    $scope.tabObservationsClicked = function () {
        $http.get(getUrl() + "/v1.0/Datastreams(" + $scope.id + ")/Observations?$top=40").then(function (response) {
            response.data.value.reverse();
            $scope.observationsList = response.data.value;

            labels = []
            values = []
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