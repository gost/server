var urls = {
    "version": getUrl() + "/Version",
    "things": getUrl() + "/Things",
};

function getVersion(){
    $.ajax({
       type: "GET",
       dataType: "json",
       url: getUrl() + "/Version",
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

function getThings(){
    $.ajax({
       type: "GET",
       dataType: "json",
       url: getUrl() + "/v1.0/Things",
       success: function(data){
         showThings(data);
       }
    });
}

function showThings(data){
    var html = "Number of things: " +  data.count + "<br/><br/>";
    
    for(var i=0;i<data.count;i++){
        html += "<b>thing " + i + "</b><br/>";
        var val0=data.value[i];
        html+="Id: " + val0["@iot.id"] + "<br/>"; 
        html+="Description: " + val0.description + "<br/>";
        html+="Organisation: " + val0.properties.organisation + "<br/>";
        html += "Link: " + "<a href = '" + val0["@iot.selfLink"] + "'>" + val0["@iot.selfLink"] + "</a><br/>";
        html += "<br/>";
    }
    
    $('#content').html(html);
}

function getUrl(){
    return window.location.origin;
}