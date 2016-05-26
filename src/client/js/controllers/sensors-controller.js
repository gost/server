gostApp.controller('SensorsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('SENSORS');
    $scope.Page.setHeaderIcon(iconSensor);

    $http.get(getUrl() + "/v1.0/Sensors").then(function (response) {
        $scope.sensorsList = response.data.value;
    });

     $scope.deleteSensorClicked = function (sensor) {
        var res = $http.delete(getUrl() + '/v1.0/Sensors(' + sensor["@iot.id"] + ')');
        res.success(function(data, status, headers, config) {
            var index = $scope.sensorsList.indexOf(sensor);
            $scope.sensorsList.splice(index, 1);
        });
        res.error(function(data, status, headers, config) {
            alert( "failure: " + JSON.stringify({data: data}));
        });
     };
});