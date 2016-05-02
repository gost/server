package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetObservedProperty todo
func (a *APIv1) GetObservedProperty(id string, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetObservedProperties todo
func (a *APIv1) GetObservedProperties(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetObservedPropertiesByDatastream todo
func (a *APIv1) GetObservedPropertiesByDatastream(datastreamID string, qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostObservedProperty todo
func (a *APIv1) PostObservedProperty(op entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PatchObservedProperty todo
func (a *APIv1) PatchObservedProperty(id string, op entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservedProperty todo
func (a *APIv1) DeleteObservedProperty(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
