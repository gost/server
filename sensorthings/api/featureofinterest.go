package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetFeatureOfInterest todo
func (a *APIv1) GetFeatureOfInterest(id string, qo *odata.QueryOptions) (*entities.FeatureOfInterest, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetFeatureOfInterests todo
func (a *APIv1) GetFeatureOfInterests(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostFeatureOfInterest todo
func (a *APIv1) PostFeatureOfInterest(foi entities.FeatureOfInterest, x string) (*entities.FeatureOfInterest, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PatchFeatureOfInterest todo
func (a *APIv1) PatchFeatureOfInterest(id string, foi entities.FeatureOfInterest) (*entities.FeatureOfInterest, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteFeatureOfInterest todo
func (a *APIv1) DeleteFeatureOfInterest(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
