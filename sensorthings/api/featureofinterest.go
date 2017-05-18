package api

import (
	"errors"

	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetFeatureOfInterest returns a FeatureOfInterest by id
func (a *APIv1) GetFeatureOfInterest(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.FeatureOfInterest{})
	if err != nil {
		return nil, err
	}

	l, err := a.db.GetFeatureOfInterest(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(l, qo)
	return l, nil
}

// GetFeatureOfInterestByObservation retrieves a FeatureOfInterest by given Observation id
func (a *APIv1) GetFeatureOfInterestByObservation(id interface{}, qo *odata.QueryOptions, path string) (*entities.FeatureOfInterest, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.FeatureOfInterest{})
	if err != nil {
		return nil, err
	}

	l, err := a.db.GetFeatureOfInterestByObservation(id, qo)
	if err != nil {
		return nil, err
	}

	a.ProcessGetRequest(l, qo)
	return l, nil
}

// GetFeatureOfInterests return FeaturesOfInterest based on the given QueryOptions
func (a *APIv1) GetFeatureOfInterests(qo *odata.QueryOptions, path string) (*models.ArrayResponse, error) {
	_, err := a.QueryOptionsSupported(qo, &entities.FeatureOfInterest{})
	if err != nil {
		return nil, err
	}

	fois, count, err := a.db.GetFeatureOfInterests(qo)
	return processFeatureOfInterest(a, fois, qo, path, count, err)
}

func processFeatureOfInterest(a *APIv1, fois []*entities.FeatureOfInterest, qo *odata.QueryOptions, path string, count int, err error) (*models.ArrayResponse, error) {
	for idx, item := range fois {
		i := *item
		a.ProcessGetRequest(&i, qo)
		fois[idx] = &i
	}

	var data interface{} = fois
	return &models.ArrayResponse{
		Count:    count,
		NextLink: a.CreateNextLink(count, path, qo),
		Data:     &data,
	}, nil
}

// PostFeatureOfInterest adds a FeatureOfInterest to the database
func (a *APIv1) PostFeatureOfInterest(foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	_, err := foi.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	supported, err2 := entities.CheckEncodingSupported(foi, foi.EncodingType)
	if !supported || err2 != nil {
		return nil, []error{err2}
	}

	l, err2 := a.db.PostFeatureOfInterest(foi)
	if err2 != nil {
		return nil, []error{err2}
	}

	l.SetAllLinks(a.config.GetExternalServerURI())
	return l, nil
}

// PutFeatureOfInterest adds a FeatureOfInterest to the database
func (a *APIv1) PutFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, []error) {
	supported, err2 := entities.CheckEncodingSupported(foi, foi.EncodingType)
	if !supported || err2 != nil {
		return nil, []error{err2}
	}

	l, err2 := a.db.PutFeatureOfInterest(id, foi)
	if err2 != nil {
		return nil, []error{err2}
	}

	l.SetAllLinks(a.config.GetExternalServerURI())
	return l, nil
}

// PatchFeatureOfInterest updates the given FeatureOfInterest in the database
func (a *APIv1) PatchFeatureOfInterest(id interface{}, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	if foi.Observations != nil {
		return nil, gostErrors.NewBadRequestError(errors.New("Unable to deep patch FeatureOfInterest"))
	}

	if len(foi.EncodingType) != 0 {
		supported, err := entities.CheckEncodingSupported(foi, foi.EncodingType)
		if !supported || err != nil {
			return nil, err
		}
	}

	return a.db.PatchFeatureOfInterest(id, foi)
}

// DeleteFeatureOfInterest deletes a given FeatureOfInterest from the database
func (a *APIv1) DeleteFeatureOfInterest(id interface{}) error {
	return a.db.DeleteFeatureOfInterest(id)
}
