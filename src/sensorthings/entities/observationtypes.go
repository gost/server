package entities

import (
)

// ObservationType
type ObservationType struct {
	Code int
	Value string
}

var (
	OMCategoryObservation = ObservationType{0,"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CategoryObservation"}  // IRI
	OMCountObservation = ObservationType{1,"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_CountObservation"} // integer
	OMMeasurement = ObservationType{2,"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Measurement"} // double
	OMObservation = ObservationType{3,"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_Observation"} // any
	OMTruthObservation = ObservationType{4,"http://www.opengis.net/def/observationType/OGC-OM/2.0/OM_TruthObservation"} // boolean
)

var ObservationTypes = []ObservationType{
	OMCategoryObservation,
	OMCountObservation,
	OMMeasurement,
	OMObservation,
	OMTruthObservation}
