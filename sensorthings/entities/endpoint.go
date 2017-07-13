package entities

// Endpoint describes the name and url of an SensorThings endpoint
// an array of endpoints is returned by the server by requesting the base path host/v1.0
type Endpoint struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}
