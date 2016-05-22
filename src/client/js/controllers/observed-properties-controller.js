gostApp.controller('ObservedPropertiesCtrl', function ($scope, $http) {
    $scope.Page.setTitle('OBSERVED PROPERTIES');
    $scope.Page.setHeaderIcon(iconObservedProperty);

    $http.get(getUrl() + "/v1.0/ObservedProperties").then(function (response) {
        $scope.observedpropertiesList = response.data.value;
    });
});