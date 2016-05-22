gostApp.controller('SensorsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('SENSORS');
    $scope.Page.setHeaderIcon(iconSensor);

    $http.get(getUrl() + "/v1.0/Sensors").then(function (response) {
        $scope.sensorsList = response.data.value;
    });
});