gostApp.controller('ObservedPropertiesCtrl', function ($scope, $http) {
    $scope.Page.setTitle('OBSERVED PROPERTIES');
    $scope.Page.setHeaderIcon(iconObservedProperty);

    $http.get(getUrl() + "/v1.0/ObservedProperties").then(function (response) {
        $scope.observedpropertiesList = response.data.value;
    });

     $scope.deleteObservedPropertyClicked = function (entity) {
        var res = $http.delete(getUrl() + '/v1.0/ObservedProperty(' + entity["@iot.id"] + ')');
        res.success(function(data, status, headers, config) {
            var index = $scope.observedpropertiesList.indexOf(entity);
            $scope.observedpropertiesList.splice(index, 1);
        });
        res.error(function(data, status, headers, config) {
            alert( "failure: " + JSON.stringify({data: data}));
        });
     };
});