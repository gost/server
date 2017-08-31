package api

import (
	"errors"

	"github.com/gost/server/sensorthings/entities"
	"github.com/gost/server/sensorthings/models"
	"github.com/gost/server/sensorthings/odata"

	gostErrors "github.com/gost/server/errors"
)

// GetSensor retrieves a sensor by id and given query
func (a *APIv1) GetSensor(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	s, err := a.db.GetSensor(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(s, qo)
	return s, nil
}

// GetSensorByDatastream retrieves a sensor by given datastream
func (a *APIv1) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	s, err := a.db.GetSensorByDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.SetLinks(s, qo)
	return s, nil
}

// GetSensors retrieves an array of sensors based on the given query
func (a *APIv1) GetSensors(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	sensors, count, err := a.db.GetSensors(qo)
	if err != nil {
		return nil, err
	}

	for idx, item := range sensors {
		i := *item
		a.SetLinks(&i, qo)
		sensors[idx] = &i
	}

	var data interface{} = sensors
	return a.createArrayResponse(count, path, qo, data), nil
}

// PostSensor adds a new sensor to the database
func (a *APIv1) PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error) {
	_, err := containsMandatoryParams(sensor)
	if err != nil {
		return nil, err
	}

	if (sensor.EncodingType != entities.EncodingPDF.Value) && (sensor.EncodingType != entities.EncodingSensorML.Value) {
		err := errors.New("Encoding not supported. Supported encoding: " + entities.EncodingPDF.Value + ", " + entities.EncodingSensorML.Value)
		return nil, []error{err}
	}

	ns, err2 := a.db.PostSensor(sensor)
	if err2 != nil {
		return nil, []error{err2}
	}

	ns.SetAllLinks(a.config.GetExternalServerURI())

	return ns, nil
}

// PatchSensor updates a sensor in the database
func (a *APIv1) PatchSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, error) {
	if sensor.Datastreams != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch Sensor"))
	}

	if len(sensor.EncodingType) != 0 {
		if (sensor.EncodingType != entities.EncodingPDF.Value) && (sensor.EncodingType != entities.EncodingSensorML.Value) {
			err := errors.New("Encoding not supported. Supported encoding: " + entities.EncodingPDF.Value + ", " + entities.EncodingSensorML.Value)
			return nil, err
		}
	}

	return a.db.PatchSensor(id, sensor)
}

// PutSensor updates the given thing in the database
func (a *APIv1) PutSensor(id interface{}, sensor *entities.Sensor) (*entities.Sensor, []error) {
	var err error
	putsensor, err := a.db.PutSensor(id, sensor)
	if err != nil {
		return nil, []error{err}
	}

	putsensor.SetAllLinks(a.config.GetExternalServerURI())
	return putsensor, nil
}

// DeleteSensor deletes a sensor from the database by given sensor id
func (a *APIv1) DeleteSensor(id interface{}) error {
	return a.db.DeleteSensor(id)
}
