## Gost MQTT - Getting started

SensorThings Api has an MQTT extension allowing users/devices to publish and subscribe to updates. In 
this article two usecases are described: Receiving and Publishing data.

## 1] Receiving MQTT data

How to subscribe to Mqtt messages:

. Console tool: Mosquitto (http://mosquitto.org/)

Mosquitto is an open source message broker that implements the MQTT protocol versions 3.1 and 3.1.1. 

```sh
mosquitto_sub -h gost.geodan.nl -t "Datastreams(1)/Observations"
```

. Desktop tool: MQTT.fx (http://mqttfx.jfx4ee.org/)

Connect to: gost.geodan.nl

Subscribe to: Datastream(1)/Observations

<img src="images/gost_mqtt_fx.png" width="353">

. Java

todo

. .NET

using NuGet package M2Mqtt (https://www.nuget.org/packages/M2Mqtt/)

<script src="https://gist.github.com/bertt/e7f30456a44e4e138ebf784bcddefad9"></script>

todo

. Python

todo

. Javascript

todo

testcase: publish an observation using the following HTTP Request and check if information is received:

```sh
curl -X POST -H "Accept: application/json" -H "Content-Type: application/json" -d '{
  "phenomenonTime": "2016-05-09T11:04:15.790Z",
  "resultTime" : "2016-05-09T11:04:15.790Z",
  "result" : 38,
  "Datastream":{"@iot.id":"1"}
}' "http://gost.geodan.nl/v1.0/Observations"
```

## 2] Publishing MQTT data

. Mosquitto

todo 

. MQTT.fx

todo

. .NET

todo 

. Python

todo

. Javascript

todo



