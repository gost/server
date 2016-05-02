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

function showThing(thing){
    var html="Id: " + thing["@iot.id"] + "<br/>"; 
    html+="Description: " + thing.description + "<br/>";
    html+="Organisation: " + thing.properties.organisation + "<br/>";
    html+="Owner: " + thing.properties.owner + "<br/>";
    html += "Link: " + "<a href = '" + thing["@iot.selfLink"] + "'>" + thing["@iot.selfLink"] + "</a><br/>";
    html += "<br/>";
    return html;
}

function showThings(data){
    var html = "Number of things: " +  data.count + "<br/><br/>";
    
    for(var i=0;i<data.count;i++){
        html += "<b>thing " + i + "</b><br/>";
        var val0=data.value[i];
        html+=showThing(val0);
    }
    
    $('#content').html(html);
}

function postThing(){
   var body= {"description": "Testsensor1","properties": {"organisation": "Geodan","owner": "Bert"}};
    $.ajax({
        type: "POST",
        data: JSON.stringify(body),
        dataType: "json",
        contentType: "application/json; charset=utf-8",
        url: getUrl() + "/v1.0/Things",
        success: function(data){
            var html = "Thing created! <br/><br/>";
            html += showThing(data);
            $('#content').html(html);
        }
    });
}


function getUrl(){
    return window.location.origin;
}