package entities

import (
	"errors"
)

// ObservationType holds the information on a EncodingType
type ObservationType struct {
	Code  int
	Value string
}

var (
	OMCategoryObservation = ObservationType{0, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"} // IRI
	OMCountObservation    = ObservationType{1, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CountObservation"}    // integer
	OMMeasurement         = ObservationType{2, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement"}         // double
	OMObservation         = ObservationType{3, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Observation"}         // any
	OMTruthObservation    = ObservationType{4, "http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_TruthObservation"}    // boolean
)

var ObservationTypes = []ObservationType{
	OMCategoryObservation,
	OMCountObservation,
	OMMeasurement,
	OMObservation,
	OMTruthObservation}

// Get the observationType based on value, returns error
func GetObservationTypeByValue(observationType string) (ObservationType, error) {
	for _, k := range ObservationTypes {
		if k.Value == observationType {
			return k, nil
		}
	}

	return OMCategoryObservation, errors.New("ObservationType not supported")
}

// Get the observationType based on value, returns error
func GetObservationTypeById(observationType int) (ObservationType, error) {
	for _, k := range ObservationTypes {
		if k.Code == observationType {
			return k, nil
		}
	}

	return OMCategoryObservation, errors.New("ObservationType not supported")
}
