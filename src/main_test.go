package main

import ()
import "testing"

func TestVersionHandler(t *testing.T) {
	//arrange
	/*server := "localhost"
	port := 8088

	cfg := configuration.Config{}
	mqttServer := mqtt.CreateMQTTClient(configuration.MQTTConfig{})
	database := postgis.NewDatabase("", 123, "", "", "", "", false, 50, 100, 200)
	api := api.NewAPI(database, cfg, mqttServer)
	api.Start()

	gostServer := http.CreateServer(server, port, &api, false, nil, nil)
	go gostServer.Start()
	versionURL := fmt.Sprintf("%s/Version", "http://"+server+":"+strconv.Itoa(port))

	fmt.Println(versionURL)
	// act
	request, _ := net.NewRequest("GET", versionURL, nil)
	res, _ := net.DefaultClient.Do(request)

	//assert
	assert.Equal(t, 200, res.StatusCode, "result should be http 200")

	// teardown
	gostServer.Stop()*/
}
