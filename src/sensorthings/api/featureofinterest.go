package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetFeatureOfInterest todo
func (a *APIv1) GetFeatureOfInterest(id string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	l, err := a.db.GetFeatureOfInterest(id)
	if err != nil {
		return nil, err
	}

	l.SetLinks(a.config.GetExternalServerURI())
	return l, nil
}

// GetFeatureOfInterestByObservation todo
func (a *APIv1) GetFeatureOfInterestByObservation(id string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	l, err := a.db.GetFeatureOfInterestByObservation(id)
	if err != nil {
		return nil, err
	}

	l.SetLinks(a.config.GetExternalServerURI())
	return l, nil
}

// GetFeatureOfInterests todo
func (a *APIv1) GetFeatureOfInterests(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	fois, err := a.db.GetFeatureOfInterests()
	return processFeatureOfInterest(a, fois, err)
}

func processFeatureOfInterest(a *APIv1, fois []*entities.FeatureOfInterest, err error) (*models.ArrayResponse, error) {
	uri := a.config.GetExternalServerURI()
	for idx, item := range fois {
		i := *item
		i.SetLinks(uri)
		fois[idx] = &i
	}

	var data interface{} = fois
	return &models.ArrayResponse{
		Count: len(fois),
		Data:  &data,
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

	l.SetLinks(a.config.GetExternalServerURI())
	return l, nil
}

// PatchFeatureOfInterest todo
func (a *APIv1) PatchFeatureOfInterest(id string, foi *entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteFeatureOfInterest todo
func (a *APIv1) DeleteFeatureOfInterest(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
