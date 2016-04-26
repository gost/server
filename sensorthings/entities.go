package sensorthings

import (
	"errors"
	"fmt"
)

// EntityType holds the name and type of an SensorThings entity
type EntityType string

// List of all EntityTypes
const (
	EntityTypeThing          	EntityType = "Thing"
	EntityTypeLocation         	EntityType = "Location"
	EntityTypeHistoricalLocation 	EntityType = "HistoricalLocation"
	EntityTypeDatastream        	EntityType = "Datastream"
	EntityTypeSensor 		EntityType = "Sensor"
	EntityTypeObservedProperty      EntityType = "ObservedProperty"
	EntityTypeObservation    	EntityType = "Observation"
	EntityTypeFeatureOfInterest    	EntityType = "FeatureOfInterest"
)

// String returns the string representation of the EntityType
func (e *EntityType) String() string{
	return fmt.Sprintf("%s", e)
}


// Thing in SensorThings represents a physical object in the real world. A Thing is a good starting
// point in which to start creating the SensorThings model structure. A Thing has a Location and one or
// more Datastreams to collect Observations. A minimal Thing can be created without a Location and Datastream
// and there are options to create a Things with a nested linked Location and Datastream.
type Thing struct {
	ID                     string                `json:"@iot.id"`
	NavSelf                string                `json:"@iot.selfLink"`
	Description            string                `json:"description"`
	Properties             map[string]string     `json:"properties,omitempty"`
	NavLocations           string                `json:"Locations@iot.navigationLink,omitempty"`
	NavDatastreams         string                `json:"Datastreams@iot.navigationLink,omitempty"`
	NavHistoricalLocations string                `json:"HistoricalLocations@iot.navigationLink,omitempty"`
	Locations              []*Location           `json:"Locations,omitempty"`
	Datastreams            []*Datastream         `json:"Datastreams,omitempty"`
	HistoricalLocations    []*HistoricalLocation `json:"HistoricalLocations,omitempty"`
}

// ContainsMandatoryPostParams checks if all mandatory params are available before posting
func (t *Thing) ContainsMandatoryPostParams() (bool, []error){
	if len(t.Description) == 0 {
		return false, []error{errors.New("Missing Thing.Description")}
	}

	return true, nil
}

func (t *Thing) SetLinks(externalUrl string){
	t.NavSelf = CreateEntitySefLink(externalUrl, "Things", t.ID)

	t.NavLocations = ""
	t.NavDatastreams = ""
	t.NavHistoricalLocations = ""

	if(t.Locations == nil) {
		t.NavLocations = CreateEntityLink("Things", "Locations", t.ID)
	}
	if(t.Datastreams == nil) {
		t.NavDatastreams = CreateEntityLink("Things", "Datastreams", t.ID)
	}
	if(t.HistoricalLocations == nil){
		t.NavHistoricalLocations = CreateEntityLink("Things", "HistoricalLocations", t.ID)
	}
}

// Location entity locates the Thing or the Things it associated with. A Thing’s Location entity is
// defined as the last known location of the Thing.
// A Thing’s Location may be identical to the Thing’s Observations’ FeatureOfInterest. In the context of the IoT,
// the principle location of interest is usually associated with the location of the Thing, especially for in-situ
// sensing applications. For example, the location of interest of a wifi-connected thermostat should be the building
// or the room in which the smart thermostat is located. And the FeatureOfInterest of the Observations made by the
// thermostat (e.g., room temperature readings) should also be the building or the room. In this case, the content
// of the smart thermostat’s location should be the same as the content of the temperature readings’ feature of interest.
type Location struct {
	ID                     string                `json:"@iot.id"`
	NavSelf                string                `json:"@iot.selfLink"`
	Description            string                `json:"descritption"`
	EncodingType           string                `json:"encodingtype"`
	Location               string                `json:"location"`
	NavThings              string                `json:"Things@iot.navigationLink,omitempty"`
	NavHistoricalLocations string                `json:"HistoricalLocations@iot.navigationLink,omitempty"`
	Things                 []*Thing              `json:"things"`
	HistoricalLocations    []*HistoricalLocation `json:"HistoricalLocations,omitempty"`
}

func (l *Location) ContainsMandatoryPostParams() (bool, []error){
	err := []error{}
	if len(l.Description) == 0 {
		err = append(err, errors.New("Missing Location.Description"))
	}

	if len(l.EncodingType) == 0 {
		err = append(err, errors.New("Missing Location.EncodingType"))
	}

	if len(l.Location) == 0 {
		err = append(err, errors.New("Missing Location.Location"))
	}

	if(len(err) > 0) {
		return false, err
	}else{
		return true, nil
	}
}

// Datastream in SensorThings represents a collection of Observations from a Sensor. A physical Sensor will send its
// data (Observations) to a single Datastream. A Datastream can have only one Sensor and that Sensor can only
// sense one ObservedProperty.
type Datastream struct {
	ID                  string            `json:"@iot.id"`
	NavSelf             string            `json:"@iot.selfLink"`
	Description         string            `json:"descritption"`
	unitOfMeasurement   map[string]string `json:"unitOfMeasurement"`
	observationType     string            `json:"observationType"`
	observedArea        map[string]string `json:"observedArea"`
	NavThings           string            `json:"Things@iot.navigationLink,omitempty"`
	NavSensors          string            `json:"Sensors@iot.navigationLink,omitempty"`
	NavObservations     string            `json:"Observations@iot.navigationLink,omitempty"`
	NavObservedProperty string            `json:"ObservedProperty@iot.navigationLink,omitempty"`
	Thing               *Thing            `json:"Thing"`
	Sensor              *Sensor           `json:"Sensor"`
	Observation         *Observation      `json:"Observation"`
	ObservedProperty    *ObservedProperty `json:"ObservedProperty"`
}

// Sensor in SensorThings represents the physical device capable of observing a physical property and converting
// it to an electrical impulse and be converted to a empirical value to represent a measurement value of the physical property
type Sensor struct {
	ID             string        `json:"@iot.id"`
	NavSelf        string        `json:"@iot.selfLink"`
	Description    string        `json:"descritption"`
	EncodingType   string        `json:"encodingtype"`
	Metadata       string        `json:"metadata"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// HistoricalLocation in sensorthings represents the current and previous locations of a thing including time
type HistoricalLocation struct {
	ID           string      `json:"@iot.id"`
	NavSelf      string      `json:"@iot.selfLink"`
	Time         string      `json:"@iot.selfLink"`
	NavThing     string      `json:"Thing@iot.navigationLink,omitempty"`
	NavLocations string      `json:"Locations@iot.navigationLink,omitempty"`
	Thing        *Thing      `json:"Thing"`
	Locations    []*Location `json:"Locations,omitempty"`
}

// ObservedProperty in SensorThings represents the physical phenomenon being observed by the Sensor. An ObserveProperty is
// linked to a Datatream which can only have one ObserveProperty
type ObservedProperty struct {
	ID             string        `json:"@iot.id"`
	NavSelf        string        `json:"@iot.selfLink"`
	Description    string        `json:"descritption"`
	Name           string        `json:"name"`
	Definition     string        `json:"definition"`
	NavDatastreams string        `json:"Datastreams@iot.navigationLink,omitempty"`
	Datastreams    []*Datastream `json:"Datastreams,omitempty"`
}

// Observation in SensorThings represents a single Sensor reading of an ObservedProperty. A physical device, a Sensor, sends
// Observations to a specified Datastream. An Observation requires a FeaturOfInterest entity, if none is provided in the request,
// the Location of the Thing associated with the Datastream, will be assigned to the new Observation as the FeaturOfInterest.
type Observation struct {
	ID                   string             `json:"@iot.id"`
	NavSelf              string             `json:"@iot.selfLink"`
	PhenomenonTime       string             `json:"phenomenonTime"`
	Result               string             `json:"result"`
	ResultTime           string             `json:"resultTime"`
	ResultQuality        string             `json:"resultQuality"`
	ValidTime            string             `json:"validTime"`
	Parameters           string             `json:"parameters"`
	NavDatastreams       string             `json:"Datastreams@iot.navigationLink,omitempty"`
	NavFeatureOfInterest string             `json:"FeatureOfInterest@iot.navigationLink,omitempty"`
	Datastreams          []*Datastream      `json:"Datastreams,omitempty"`
	FeatureOfInterest    *FeatureOfInterest `json:"FeatureOfInterest,omitempty"`
}

// FeatureOfInterest in SensorThings represents the phenomena an Observation is detecting. In some cases a FeatureOfInterest
// can be the Location of the Sensor and therefore of the Observation. A FeatureOfInterest is linked to a single Observation
type FeatureOfInterest struct {
	ID              string       `json:"@iot.id"`
	NavSelf         string       `json:"@iot.selfLink"`
	Description     string       `json:"descritption"`
	EncodingType    string       `json:"encodingtype"`
	Feature         string       `json:"feature"`
	NavObservations string       `json:"Observations@iot.navigationLink,omitempty"`
	Observation     *Observation `json:"Observation"`
}
