var urls = {
    "version": getUrl() + "/Version",
};

function getVersion(){
    $.ajax({
       type: "GET",
       dataType: "json",
       url: urls["version"],
       success: function(data){
         showVersion(data);
       }
    });
}


function showVersion(data){
    $('#content').html("<b>GOST HTTP Server</b><br/>"
    + data.gostServerVersion.version +
    "<br/><br/><b>SensorThings</b><br/>" +
     data.sensorThingsApiVersion.version);
}

function getUrl(){
    return window.location.origin;
}