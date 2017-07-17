package configuration

const (
	// ServerVersion specifies the current GOST Server version
	ServerVersion string = "v0.5"

	// SensorThingsAPIVersion specifies the supported SensorThings API version
	SensorThingsAPIVersion string = "v1.0"

	// DefaultMaxEntries is used when config maxEntries is empty or $top exceeds this default value
	DefaultMaxEntries int = 200
)
