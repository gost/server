package api

import (
	"errors"

	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"

	gostErrors "github.com/geodan/gost/errors"
)

// GetSensor retrieves a sensor by id and given query
func (a *APIv1) GetSensor(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	s, err := a.db.GetSensor(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(s, qo)
	return s, nil
}

// GetSensorByDatastream retrieves a sensor by given datastream
func (a *APIv1) GetSensorByDatastream(id interface{}, qo *odata.QueryOptions, path string) (*entities.Sensor, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	s, err := a.db.GetSensorByDatastream(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(s, qo)
	return s, nil
}

// GetSensors retrieves an array of sensors based on the given query
func (a *APIv1) GetSensors(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.Sensor{})
	if err != nil {
		return nil, err
	}

	sensors, count, err := a.db.GetSensors(qo)
	if err != nil {
		return nil, err
	}

	for idx, item := range sensors {
		i := *item
		a.ProcessGetRequest(&i, qo)
		sensors[idx] = &i
	}

	var data interface{} = sensors
	return &models.ArrayResponse{
		Count:    count,
		NextLink: a.CreateNextLink(count, path, qo),
		Data:     &data,
	}, nil
}

// PostSensor adds a new sensor to the database
func (a *APIv1) PostSensor(sensor *entities.Sensor) (*entities.Sensor, []error) {
	_, err := sensor.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	supported, err2 := entities.CheckEncodingSupported(sensor, sensor.EncodingType)
	if !supported || err2 != nil {
		return nil, []error{err2}
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
		supported, err := entities.CheckEncodingSupported(sensor, sensor.EncodingType)
		if !supported || err != nil {
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
