gostApp.controller('HistoricalLocationsCtrl', function ($scope, $http) {
    $scope.Page.setTitle('HISTORICAL LOCATIONS');
    $scope.Page.setHeaderIcon(iconHistoricalLocation);

    $http.get(getUrl() + "/v1.0/HistoricalLocations").then(function (response) {
        $scope.historicallocationsList = response.data.value;
    });

     $scope.deleteHistoricalLocationClicked = function (entity) {
        var res = $http.delete(getUrl() + '/v1.0/HistoricalLocations(' + entity["@iot.id"] + ')');
        res.success(function(data, status, headers, config) {
            var index = $scope.historicallocationsList.indexOf(entity);
            $scope.historicallocationsList.splice(index, 1);
        });
        res.error(function(data, status, headers, config) {
            alert( "failure: " + JSON.stringify({data: data}));
        });
     };
});