package entities

import (
	"errors"
)

// ObservationType holds the information on a ObservationType
type ObservationType struct {
	Code  int64
	Value string
}

// List of supported ObservationTypes
var (
	OMCategoryUnknown     = ObservationType{0, ""}
	OMCategoryObservation = ObservationType{1, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"} // IRI
	OMCountObservation    = ObservationType{2, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CountObservation"}    // integer
	OMMeasurement         = ObservationType{3, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement"}         // double
	OMObservation         = ObservationType{4, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Observation"}         // any
	OMTruthObservation    = ObservationType{5, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_TruthObservation"}    // boolean
)

// ObservationTypes is a list of names mapped to their ObservationType Value
var ObservationTypes = []ObservationType{
	OMCategoryUnknown,
	OMCategoryObservation,
	OMCountObservation,
	OMMeasurement,
	OMObservation,
	OMTruthObservation}

// GetObservationTypeByValue Get the observationType based on value, returns error
func GetObservationTypeByValue(observationType string) (ObservationType, error) {
	for _, k := range ObservationTypes {
		if k.Value == observationType {
			return k, nil
		}
	}

	return OMCategoryObservation, errors.New("ObservationType not supported")
}

// GetObservationTypeByID Get the observationType based on value, returns error
func GetObservationTypeByID(observationType int64) (ObservationType, error) {
	for _, k := range ObservationTypes {
		if k.Code == observationType {
			return k, nil
		}
	}

	return OMCategoryObservation, errors.New("ObservationType not supported")
}
