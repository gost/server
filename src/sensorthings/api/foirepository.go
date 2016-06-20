package api

import (
	"errors"

	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
)

//ToDo: Locking?

type FoiRepository struct {
	db                *models.Database
	updatedThings     []string
	thingToFoi        map[string]string
	datastreamToThing map[string]string
}

// LoadInMemory loads the previous states into memory
// Todo implement
func (f *FoiRepository) LoadInMemory() {

}

// ThingLocationUpdated should be called when a location by thing is updated, when needed a new
// FOI will be created and inserted in the database
func (f *FoiRepository) ThingLocationUpdated(thingID string) {
	// Do not add if already exist
	for _, t := range f.updatedThings {
		if t == thingID {
			return
		}
	}

	f.updatedThings = append(f.updatedThings, thingID)
}

// GetThing retrieves a thing ID from memory by datastream ID
func (f *FoiRepository) GetFoiIDByDatastreamID(datastreamID string) (string, error) {
	db := *f.db

	var ok bool
	var err error
	var tID string

	// CHECK IF DATASTREAM IS IN LIST, IF NOT LOAD FROM DATABASE INCLUDING THING ADD TO LISTS, ERROR IF NOT FOUND
	if tID, ok = f.datastreamToThing[datastreamID]; !ok {
		// Datastream not found in list, look in database
		if ds, err := db.GetDatastream(datastreamID, nil); err != nil {
			// Datastream not found in database
			return "", errors.New("Datastream not found")
		} else {
			// Datastream found search for thing by datastream
			if thing, err := db.GetThingByDatastream(ds.ID, nil); err != nil {
				return "", errors.New("Thing by datastream not found")
			} else {
				// Thing found, setup in datastreamToThing
				f.datastreamToThing[ds.ID.(string)] = thing.ID.(string)
				tID = thing.ID.(string)
			}
		}
	}

	var foiID string

	// Check if thing to foi is added
	if foiID, ok = f.thingToFoi[tID]; !ok {
		t, err := db.GetThing(tID, nil)

		// Thing not found in database
		if err != nil {
			return "", errors.New("Thing not found")
		}

		// Create a new foi
		if foiID, err = f.insertFoi(t.ID.(string)); err != nil {
			return "", errors.New("Error adding FeatureOfInterest")
		}

		// Update ThingToFoi list
		f.thingToFoi[t.ID.(string)] = foiID
	}

	//check if foi needs to be updated
	for idx, t := range f.updatedThings {
		if t == tID {
			//should update
			if foiID, err = f.insertFoi(tID); err != nil {
				return "", errors.New("Error adding FeatureOfInterest")
			}

			f.thingToFoi[tID] = foiID
			f.updatedThings = append(f.updatedThings[:idx], f.updatedThings[idx+1:]...)
			break
		}
	}

	// ToDo: what happens when GOST is closed and some fois need to be updated
	return foiID, nil
}

func (f *FoiRepository) insertFoi(thingID string) (string, error) {
	db := *f.db
	t, err := db.GetThing(thingID, nil)
	if err != nil {
		return "", err
	}

	l, err := db.GetLocationsByThing(t.ID, nil)
	if err != nil || len(l) == 0 {
		return "", err
	}

	foi := &entities.FeatureOfInterest{}
	foi.Description = l[0].Description
	foi.EncodingType = l[0].EncodingType
	foi.Feature = l[0].Location

	nFoi, err := db.PostFeatureOfInterest(foi)
	if err != nil {
		return "", err
	}

	return nFoi.ID.(string), nil
}
