gostApp.controller('SensorsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('SENSORS');
    $scope.Page.setHeaderIcon(iconSensor);

    $http.get(getUrl() + "/v1.0/Sensors").then(function (response) {
        $scope.sensorsList = response.data.value;
    });

     $scope.deleteSensorClicked = function (sensorId) {
        var res = $http.delete(getUrl() + '/v1.0/Sensors(' + sensorId + ')');
        res.success(function(data, status, headers, config) {
            alert( "removed sensor");
        });
      };
});