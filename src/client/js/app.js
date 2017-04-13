var gostApp = angular.module('gostApp', ['ngRoute']);

gostApp.factory('Page', function () {
    var title = 'GOST';
    return {
        title: function () { return title; },
        setTitle: function (newTitle) { title = newTitle; }
    };
});

// Set routes
gostApp.config(function ($routeProvider) {
    $routeProvider.
        when("/home", { templateUrl: "views/home.html", controller: "HomeCtrl" }).
        when('/thing/:id', { templateUrl: 'views/thing.html', controller: "ThingCtrl" }).
        when("/things", { templateUrl: "views/things.html", controller: "ThingsCtrl" }).
        when('/locations', { templateUrl: 'views/locations.html', controller: "LocationsCtrl" }).
        when('/historicallocations', { templateUrl: 'views/historicallocations.html', controller: "HistoricalLocationsCtrl" }).
        when('/sensors', { templateUrl: 'views/sensors.html', controller: "SensorsCtrl" }).
        when('/observedproperties', { templateUrl: 'views/observedproperties.html', controller: "ObservedPropertiesCtrl" }).
        when('/datastream/:id', { templateUrl: 'views/datastream.html', controller: "DatastreamCtrl" }).
        when('/datastreams', { templateUrl: 'views/datastreams.html', controller: "DatastreamsCtrl" }).
        when('/observations', { templateUrl: 'views/observations.html', controller: "ObservationsCtrl" }).
        when('/featuresofinterest', { templateUrl: 'views/featuresofinterest.html', controller: "FeaturesOfInterestCtrl" }).
        when('/mqtt', { templateUrl: 'views/mqtt.html', controller: "MqttCtrl" }).
        otherwise({redirectTo: '/home' });;
})

var iconThing = "fa-cube";
var iconDatastream = "fa-line-chart";
var iconFeatureOfInterest = "fa-map";
var iconHistoricalLocation = "fa-history";
var iconLocation = "fa-map-marker";
var iconObservation = "fa-eye";
var iconSensor = "fa-dashboard";
var iconObservedProperty = "fa-list";

function getUrl() {
    var path = window.location.pathname;
    path = path.substring(0, path.lastIndexOf('/Dashboard'));
    return window.location.origin + path;
}

var olMap;
var geoJSONLayer;
var olCollectionGeoJSON;
var defaultMarkerStyle = new ol.style.Style({
    image: new ol.style.Icon(({
        anchor: [0.5, 1],
        anchorXUnits: 'fraction',
        anchorYUnits: 'fraction',
        src: 'assets/img/map_marker.svg'
    }))
})

function createMap(target) {
    $(".map").empty();

    olCollectionGeoJSON = new ol.Collection();
    jsonSource = new ol.source.Vector({
        format: new ol.format.GeoJSON(),
        features: olCollectionGeoJSON
    });

    geoJSONLayer = new ol.layer.Vector({
        source: jsonSource,
        style: defaultMarkerStyle
    });

    olMap = new ol.Map({
        controls: ol.control.defaults({
            attributionOptions: ({
                collapsible: false
            })
        }).extend([
        ]),
        layers: [
            /*new ol.layer.Tile({
                source: new ol.source.OSM()
            }),*/
            new ol.layer.Tile({
                source: new ol.source.XYZ({
                    tileSize: [512, 512],
                    url: 'https://api.mapbox.com/styles/v1/edenhalperin/cih84uopy000a95m41htugsnm/tiles/{z}/{x}/{y}?access_token=pk.eyJ1IjoiZWRlbmhhbHBlcmluIiwiYSI6IlFRZG0zMWMifQ.QUNKx4tIMjZfwmrE8SE6Bg'
                })
            }),
            geoJSONLayer
        ],
        target: 'map',
        view: new ol.View({
            zoom: 8,
            center: ol.proj.transform([5.3833333, 52.15], 'EPSG:4326', 'EPSG:3857'),
            maxZoom: 20
        }),
        interactions: ol.interaction.defaults({ altShiftDragRotate: false, pinchRotate: false })
    });

    setTimeout(delayedUpdateMap, 300);
}

function addGeoJSONFeature(geojson) {
    var defaultGeoJSONProjection = 'EPSG:4326';
    var mapProjection = olMap.getView().getProjection();
    var geoJSONString = JSON.stringify(geojson);

    if (geoJSONString.includes("FeatureCollection")) {
        olCollectionGeoJSON.push((new ol.format.GeoJSON()).readFeatures(geojson, { dataProjection: defaultGeoJSONProjection, featureProjection: mapProjection }));
    }
    else {
        var geom = (new ol.format.GeoJSON()).readGeometry(geojson, { dataProjection: defaultGeoJSONProjection, featureProjection: mapProjection });
        var feature = new ol.Feature(geom);
        olCollectionGeoJSON.push(feature);
    }
}

function zoomToGeoJSONLayerExtent() {
    var extent = geoJSONLayer.getSource().getExtent();
    olMap.getView().fit(extent, olMap.getSize());
}

function delayedUpdateMap() {
    olMap.updateSize();
}

var myChart;

function createObservationChart(labels, values){
	if(myChart != null){
		myChart.destroy();	
	}
	
    //$("#observationChartWrapper").empty();
    Chart.defaults.global.responsive = true;
    Chart.defaults.global.maintainAspectRatio = false;
    Chart.defaults.global.legend.display = false;
    chartData = {
           labels: labels,
           datasets: [
               {
                   fill: true,
                   lineTension: 0.1,
                   backgroundColor: "rgba(115,135,156,0.15)",
                   borderColor: "rgba(115,135,156,1)",
                   borderCapStyle: 'butt',
                   borderDash: [],
                   borderDashOffset: 0.0,
                   borderJoinStyle: 'miter',
                   pointBorderColor: "rgba(115,135,156,1)",
                   pointBackgroundColor: "#fff",
                   pointBorderWidth: 1,
                   pointHoverRadius: 5,
                   pointHoverBackgroundColor: "rgba(75,192,192,1)",
                   pointHoverBorderColor: "rgba(220,220,220,1)",
                   pointHoverBorderWidth: 2,
                   pointRadius: 2,
                   pointHitRadius: 10,
                   data: values,
               }
           ]
       }

    var ctx = $("#observationChart");
	ctx.innerHTML = "";
    myChart = new Chart(ctx, {
        type: 'line',
        data: chartData,
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:false
                    }
                }]
            }
        }
    });
}

function observationChartAddData(l, d){
	myChart.data.datasets[0].data = d;
	myChart.data.labels = l;
	myChart.update();
}

$('.fs-button').on('click', function(){
      if (!document.fullscreenElement &&    // alternative standard method
          !document.mozFullScreenElement && !document.webkitFullscreenElement && !document.msFullscreenElement ) {  // current working methods
        if (document.documentElement.requestFullscreen) {
          document.documentElement.requestFullscreen();
        } else if (document.documentElement.msRequestFullscreen) {
          document.documentElement.msRequestFullscreen();
        } else if (document.documentElement.mozRequestFullScreen) {
          document.documentElement.mozRequestFullScreen();
        } else if (document.documentElement.webkitRequestFullscreen) {
          document.documentElement.webkitRequestFullscreen(Element.ALLOW_KEYBOARD_INPUT);
        }
      } else {
        if (document.exitFullscreen) {
          document.exitFullscreen();
        } else if (document.msExitFullscreen) {
          document.msExitFullscreen();
        } else if (document.mozCancelFullScreen) {
          document.mozCancelFullScreen();
        } else if (document.webkitExitFullscreen) {
          document.webkitExitFullscreen();
        }
      }
});

function guid() {
  function s4() {
    return Math.floor((1 + Math.random()) * 0x10000)
      .toString(16)
      .substring(1);
  }
  return s4() + s4() + '-' + s4() + '-' + s4() + '-' +
    s4() + '-' + s4() + s4() + s4();
}
