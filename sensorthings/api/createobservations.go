package api

import (
	"fmt"

	entities "github.com/gost/core"
)

// PostCreateObservations checks for correctness of the datastreams and observations and calls PostcreateObservations on the database
// ToDo: use transactions
func (a *APIv1) PostCreateObservations(data *entities.CreateObservations) ([]string, []error) {
	_, err := containsMandatoryParams(data)
	if err != nil {
		return nil, err
	}

	returnList := make([]string, 0)
	for i := 0; i < len(data.Datastreams); i++ {
		for j := 0; j < len(data.Datastreams[i].Observations); j++ {
			obs, errors := a.PostObservationByDatastream(data.Datastreams[i].ID, data.Datastreams[i].Observations[j])
			if errors == nil || len(errors) == 0 {
				returnList = append(returnList, obs.GetSelfLink())
			} else {
				errorString := ""
				for k := 0; k < len(errors); k++ {
					if len(errorString) > 0 {
						errorString += ", "
					}

					errorString += fmt.Sprintf("%v", errors[k].Error())
				}
				returnList = append(returnList, errorString)
			}
		}
	}

	return returnList, nil
}
