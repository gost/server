package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/geodan/gost/src/sensorthings/entities"
	"github.com/geodan/gost/src/sensorthings/models"
	"log"
)

const (
	foiStateBucketName      = "foistate"
	foiDatastreamToThingKey = "datastreamToThing"
	foiThingToFoiKey        = "thingToFoi"
	foiUpdatedThingsKey     = "updatedThings"
)

var foiRepoLocker = &sync.Mutex{}

// FoiRepository is used to get a FeatureOfInterest (FOI) id by given datastream id. When an observation
// is added without any FOI id or FOI deep insert the linked things (last) Location should be used as FOI.
// FoiRepository loads the current state in memory and keeps track of FOI ID's created for a things location
// to do fast look ups of a FOI by observation.
type FoiRepository struct {
	db                *InternalDatabase
	updatedThings     []string
	thingToFoi        map[string]string
	datastreamToThing map[string]string
}

// LoadInMemory loads the previous FOI states into memory
func (f *FoiRepository) LoadInMemory() {
	f.updatedThings = make([]string, 0)
	f.thingToFoi = make(map[string]string, 0)
	f.datastreamToThing = make(map[string]string, 0)

	if !f.db.open {
		panic("BoltDB not opened yet")
	}

	f.db.bolt.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(foiStateBucketName))
		return nil
	})

	f.db.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(foiStateBucketName))
		if b == nil {
			return nil
		}

		json.Unmarshal(b.Get([]byte(foiUpdatedThingsKey)), &f.updatedThings)
		json.Unmarshal(b.Get([]byte(foiDatastreamToThingKey)), &f.datastreamToThing)
		json.Unmarshal(b.Get([]byte(foiThingToFoiKey)), &f.thingToFoi)

		return nil
	})
}

// ThingLocationUpdated should be called when a location by thing is updated, if needed a new
// FOI will be created and inserted into the database
func (f *FoiRepository) ThingLocationUpdated(thingID string) {
	// Do not add if already exist
	for _, t := range f.updatedThings {
		if t == thingID {
			return
		}
	}

	foiRepoLocker.Lock()
	log.Printf("Location updated")
	f.updatedThings = append(f.updatedThings, thingID)
	f.SaveState(foiUpdatedThingsKey)
	foiRepoLocker.Unlock()
}

// GetFoiIDByDatastreamID retrieves a FOI ID from memory by given datastream ID
func (f *FoiRepository) GetFoiIDByDatastreamID(gdb *models.Database, datastreamID interface{}) (string, error) {
	db := *gdb

	dId := toStringID(datastreamID)

	var ok bool
	var err error
	var foiID, tID string

	// check if datastream is in list, if not load from database including thing add to list, error if not found
	if tID, ok = f.datastreamToThing[dId]; !ok {
		// Datastream not found in list, look in database
		if _, err := db.GetDatastream(dId, nil); err != nil {
			// Datastream not found in database
			return "", errors.New("Datastream not found")
		} else {
			// Datastream found search for thing by datastream
			if thing, err := db.GetThingByDatastream(dId, nil); err != nil {
				return "", errors.New("Thing by datastream not found")
			} else {
				// Thing found, setup in datastreamToThing
				tStringId := toStringID(thing.ID)
				foiRepoLocker.Lock()
				log.Printf("Datastream to thing added: datastream: %v thing: %v", dId, tStringId)
				f.datastreamToThing[dId] = tStringId
				f.SaveState(foiDatastreamToThingKey)
				foiRepoLocker.Unlock()
				tID = tStringId
			}
		}
	}

	// Check if thing to foi is added
	if foiID, ok = f.thingToFoi[tID]; !ok {
		t, err := db.GetThing(tID, nil)

		// Thing not found in database
		if err != nil {
			return "", errors.New("Thing not found")
		}

		// Create a new foi
		if foiID, err = f.insertFoi(gdb, tID); err != nil {
			return "", errors.New("Error adding FeatureOfInterest")
		}

		// Update ThingToFoi list
		foiRepoLocker.Lock()
		f.thingToFoi[tID] = foiID
		log.Printf("ThingToFoi added: thing: %v foi: %v", toStringID(t.ID), foiID)
		f.SaveState(foiThingToFoiKey)
		foiRepoLocker.Unlock()
	}

	//check if foi needs to be updated (new location added for thing but no FOI created for it yet)
	for idx, t := range f.updatedThings {
		if t == tID {
			if foiID, err = f.insertFoi(gdb, tID); err != nil {
				return "", errors.New("Error adding FeatureOfInterest")
			}

			foiRepoLocker.Lock()
			f.thingToFoi[tID] = foiID
			f.updatedThings = append(f.updatedThings[:idx], f.updatedThings[idx+1:]...)
			f.SaveState(foiUpdatedThingsKey)
			foiRepoLocker.Unlock()
			break
		}
	}

	return foiID, nil
}

// insertFoi inserts a new FOI into the database and returns it's new ID
func (f *FoiRepository) insertFoi(gdb *models.Database, thingID string) (string, error) {
	db := *gdb

	l, err := db.GetLocationsByThing(thingID, nil)
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

	log.Printf("Inserted %v \n", nFoi.Description)

	return toStringID(nFoi.ID), nil
}

func (f *FoiRepository) SaveState(key string) error {
	if !f.db.open {
		return fmt.Errorf("db must be opened before saving!")
	}

	return f.db.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(foiStateBucketName))

		var value interface{}
		switch key {
		case foiUpdatedThingsKey:
			value = f.updatedThings
		case foiDatastreamToThingKey:
			value = f.datastreamToThing
		case foiThingToFoiKey:
			value = f.thingToFoi
		}

		enc, err := json.Marshal(value)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), enc)
	})
}

func toStringID(id interface{}) string {
	return fmt.Sprintf("%v", id)
}
