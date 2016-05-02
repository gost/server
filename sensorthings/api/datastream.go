// Copyright (c) 2016 The GOST Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"errors"
	gostErrors "github.com/geodan/gost/errors"
	"github.com/geodan/gost/sensorthings/entities"
	"github.com/geodan/gost/sensorthings/models"
	"github.com/geodan/gost/sensorthings/odata"
)

// GetDatastream todo
func (a *APIv1) GetDatastream(id string, qo *odata.QueryOptions) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// GetDatastreams todo
func (a *APIv1) GetDatastreams(qo *odata.QueryOptions) (*models.ArrayResponse, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// PostDatastream todo
func (a *APIv1) PostDatastream(datastream entities.Datastream, x string) (*entities.Datastream, []error) {
	return nil, []error{gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))}
}

// PatchDatastream todo
func (a *APIv1) PatchDatastream(id string, datastream entities.Datastream) (*entities.Datastream, error) {
	return nil, gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}

// DeleteDatastream todo
func (a *APIv1) DeleteDatastream(id string) error {
	return gostErrors.NewRequestNotImplemented(errors.New("not implemented yet"))
}
