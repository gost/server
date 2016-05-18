package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/src/errors"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"github.com/geodan/gost/src/sensorthings/odata"
)

// GetObservedProperty todo
func (a *APIv1) GetObservedProperty(id string, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	op, err := a.db.GetObservedProperty(id)
	if err != nil {
		return nil, err
	}

	op.SetLinks(a.config.GetExternalServerURI())
	return op, nil
}

// GetObservedPropertyByDatastream todo
func (a *APIv1) GetObservedPropertyByDatastream(datastreamID string, qo *odata.QueryOptions) (*entities.ObservedProperty, error) {
	op, err := a.db.GetObservedPropertyByDatastream(datastreamID)
	if err != nil {
		return nil, err
	}

	op.SetLinks(a.config.GetExternalServerURI())
	return op, nil
}

// GetObservedProperties todo
func (a *APIv1) GetObservedProperties(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	ops, err := a.db.GetObservedProperties()
	if err != nil {
		return nil, err
	}

	uri := a.config.GetExternalServerURI()
	for idx, item := range ops {
		i := *item
		i.SetLinks(uri)
		ops[idx] = &i
	}

	var data interface{} = ops
	response := models.ArrayResponse{
		Count: len(ops),
		Data:  &data,
	}

	return &response, nil
}

// PostObservedProperty todo
func (a *APIv1) PostObservedProperty(op *entities.ObservedProperty) (*entities.ObservedProperty, []error) {
	_, err := op.ContainsMandatoryParams()
	if err != nil {
		return nil, err
	}

	nop, err2 := a.db.PostObservedProperty(op)
	if err2 != nil {
		return nil, []error{err2}
	}

	return nop, nil
}

// PatchObservedProperty todo
func (a *APIv1) PatchObservedProperty(id string, op *entities.ObservedProperty) (*entities.ObservedProperty, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteObservedProperty todo
func (a *APIv1) DeleteObservedProperty(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
