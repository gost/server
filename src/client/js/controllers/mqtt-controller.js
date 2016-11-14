gostApp.controller('MqttCtrl', function($scope) {
    $scope.Page.setTitle('MQTT');
    $scope.Page.setHeaderIcon('fa-envelope');

    $scope.$on("$destroy", function() {
        client.unsubscribe("$SYS/broker/#");
    });

    client = new Paho.MQTT.Client(location.hostname, Number(9001), guid());
    client.onConnectionLost = onConnectionLost;
    client.onMessageArrived = onMessageArrived;
    client.connect({ onSuccess: onConnect });

    function onConnect() {
        client.subscribe("$SYS/broker/#");
    }

    function onConnectionLost(responseObject) {
        if (responseObject.errorCode !== 0) {
            console.log("onConnectionLost:" + responseObject.errorMessage);
        }
    }

    // called when a message arrives
    function onMessageArrived(message) {
        switch (message.destinationName) {
            case '$SYS/broker/version':
                $scope.brokerVersion = message.payloadString;
                break;
            case '$SYS/broker/timestamp':
                $scope.brokerTimestamp = message.payloadString;
                break;
            case '$SYS/broker/uptime':
                $scope.brokerUptime = message.payloadString;
                break;
            case '$SYS/broker/subscriptions/count':
                $scope.brokerSubscriptions = message.payloadString;
                break;
            case '$SYS/broker/clients/connected':
                $scope.clientsConnected = message.payloadString;
                break;
            case '$SYS/broker/clients/disconnected':
                $scope.clientsDisconnected = message.payloadString;
                break;
            case '$SYS/broker/clients/expired':
                $scope.clientsExpired = message.payloadString;
                break;
            case '$SYS/broker/messages/sent':
                $scope.messagesSent = message.payloadString;
                break;
            case '$SYS/broker/messages/received':
                $scope.messagesReceived = message.payloadString;
                break;
            case '$SYS/broker/messages/stored':
                $scope.messagesStored = message.payloadString;
                break;
            case '$SYS/broker/bytes/sent':
                $scope.trafficSent = message.payloadString;
                break;
            case '$SYS/broker/bytes/received':
                $scope.trafficReceived = message.payloadString;
                break;
        }

        $scope.$apply();
    }
});